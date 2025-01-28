package web

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"net"
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
	Conn  net.Conn
	Mutex sync.Mutex
}

func (s *Server) GetFile(ctx context.Context) error {
	if config.Config.Local.Webcam.Enable {
		return fmt.Errorf("%s", "Webcam enabled..")
	}

	// Receive the output file
	fmt.Fprintln(s.Conn, "REQUEST_FILE")
	time.Sleep(500 * time.Microsecond) // Small delay between messages.
	_, err := s.ReceiveFile()
	if err != nil {
		fmt.Println(err)
		return err
	}
	ok, _ := s.WaitForDone(ctx, make([]byte, 4096))
	if !ok {
		return fmt.Errorf("%s", "Did not reacieve done msg..")
	}
	// fmt.Fprintln(s.Conn, "EXIT")
	// time.Sleep(500 * time.Microsecond)
	return nil
}
func (s *Server) HandleIncomingCommands(ctx context.Context, webcam_source int) {
	defer s.Conn.Close()

	scanner := bufio.NewScanner(s.Conn)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Printf("Received command: %s\n", command)

		switch command {
		// Server asks for the source again if it doesn't have it
		case CommandSendSource:
			if err := s.SendFile(command, config.Config.LastSource()); err != nil {
				fmt.Printf("Error receiving source: %v\n", err)
			}

		case CommandSendTarget:
			// Server asks for the target again if it doesn't have it
			if err := s.SendFile(command, config.Config.LastTarget()); err != nil {
				fmt.Printf("Error receiving target: %v\n", err)
			}

		case CommandRequestFile:
			// Server is ready to send the processed file
			if err := s.GetFile(ctx); err != nil {
				fmt.Printf("Error receiving file: %v\n", err)
			}

		case CommandStartFrames:
			// Server needs frames to start again
			s.mux.Lock()
			s.Ready = true
			s.mux.Unlock()
			go webcam.StartFrameChannel(ctx, webcam_source)
			s.WG.Add(1)
			go s.Frames.Start(&s.WG, s.Conn)
			s.WG.Wait()

		case CommandStopFrames:
			// Server stops the frames
			s.mux.Lock()
			s.Ready = false
			s.mux.Unlock()

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

// Optionally v2?
func startWebcamStream(client *Server) {
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
