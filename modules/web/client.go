package web

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"goStreamer/modules/hardware/webcam"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

const (
	CommandSendSource  = "SEND_SOURCE"
	CommandSendTarget  = "SEND_TARGET"
	CommandRequestFile = "REQUEST_FILE"
	CommandStartFrames = "START_FRAMES"
	CommandStopFrames  = "STOP_FRAMES"
	CommandExit        = "EXIT"
	BufferSize         = 1024
	FrameRate          = 30
)

type Client struct {
	Conn       net.Conn
	Ready      bool
	HaveSource bool
	HaveTarget bool
	Mutex      sync.Mutex
}

func (c *Client) SendMessage(message string) error {
	_, err := c.Conn.Write([]byte(message + "\n"))
	return err
}

func (c *Client) ReceiveFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(c.Conn)
	var fileSize int64
	if err := binary.Read(reader, binary.LittleEndian, &fileSize); err != nil {
		return err
	}

	written, err := io.CopyN(file, reader, fileSize)
	if err != nil {
		return err
	}

	fmt.Printf("Received file %s (%d bytes)\n", filename, written)
	return nil
}

func (c *Client) SendFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := info.Size()
	buffer := make([]byte, BufferSize)

	// Send file size
	if err := binary.Write(c.Conn, binary.LittleEndian, fileSize); err != nil {
		return err
	}

	// Send file content
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if _, err := c.Conn.Write(buffer[:n]); err != nil {
			return err
		}
	}

	fmt.Printf("Sent file %s (%d bytes)\n", filename, fileSize)
	return nil
}

func (s *Server) HandleClient(ctx context.Context, webcam_source int) {
	defer s.Conn.Close()

	client := &Client{
		Conn: s.Conn,
	}

	scanner := bufio.NewScanner(s.Conn)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Printf("Received command: %s\n", command)

		switch command {
		case CommandSendSource:
			client.Mutex.Lock()
			client.HaveSource = true
			client.Mutex.Unlock()
			if err := client.ReceiveFile("source.txt"); err != nil {
				fmt.Printf("Error receiving source: %v\n", err)
			}

		case CommandSendTarget:
			client.Mutex.Lock()
			client.HaveTarget = true
			client.Mutex.Unlock()
			if err := client.ReceiveFile("target.txt"); err != nil {
				fmt.Printf("Error receiving target: %v\n", err)
			}

		case CommandRequestFile:
			filename := "processed.txt"
			if err := client.SendFile(filename); err != nil {
				fmt.Printf("Error sending file: %v\n", err)
			}

		case CommandStartFrames:
			client.Mutex.Lock()

			go webcam.StartFrameChannel(ctx, webcam_source)
			client.Ready = true
			client.Mutex.Unlock()
			go startWebcamStream(client)

		case CommandStopFrames:
			client.Mutex.Lock()
			client.Ready = false
			client.Mutex.Unlock()

		case CommandExit:
			fmt.Println("Client disconnected.")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from client: %v\n", err)
	}
}

func startWebcamStream(client *Client) {
	stream := mjpeg.NewStream()

	// Simulated webcam feed using generated images
	fmt.Println("Starting webcam stream...")
	for client.Ready {
		// Receive frames from the channel
		for frame := range webcam.FrameChan {
			buffer := new(bytes.Buffer)

			// Convert gocv.Mat to JPEG
			buf, err := gocv.IMEncode(".jpg", frame.Mat)
			if err != nil {
				fmt.Printf("Error encoding frame to JPEG: %v\n", err)
				continue
			}

			// Update the MJPEG stream
			stream.UpdateJPEG(buf.GetBytes())

			buffer.Reset()
		}

		// Throttle the frame rate
		time.Sleep(time.Second / FrameRate)
	}
	fmt.Println("Stopped webcam stream.")
}

func generateFakeFrame() *bytes.Buffer {
	// Placeholder for a function to capture webcam frames
	return bytes.NewBuffer([]byte{})
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on port 8080...")

}
