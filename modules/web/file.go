package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"goStreamer/modules/settings"
)

// Header represents the metadata for a file transfer.
type Header struct {
	ClientIP   string `json:"client_ip"`
	ClientPort int    `json:"client_port"`
	Command    string `json:"command"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
}

// SendFile streams a file to the remote server.
func (s *Server) SendFile(command, filePath string) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}
	fileSize := stat.Size()

	// Prepare header as JSON
	header := Header{
		ClientIP:   settings.Settings.Client.Net.IP,
		ClientPort: settings.Settings.Client.Net.Port,
		Command:    command,
		FileName:   stat.Name(),
		FileSize:   fileSize,
	}
	headerData, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal header: %v", err)
	}

	// Send header length and data
	if _, err = s.Conn.Write(headerData); err != nil {
		return fmt.Errorf("failed to send header: %v", err)
	}
	time.Sleep(1 * time.Second)
	// Stream the file
	_, err = io.Copy(s.Conn, file)
	if err != nil {
		return fmt.Errorf("failed to send file: %v", err)
	}

	return nil
}

// ReceiveFile receives a file from the remote connection and saves it locally.
func (s *Server) ReceiveFile() (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Read header data
	headerData := make([]byte, 1024)
	n, err := s.Conn.Read(headerData)
	if err != nil {
		return "", fmt.Errorf("failed to read header: %v", err)
	}

	var header Header
	if err := json.Unmarshal(headerData[:n], &header); err != nil {
		return "", fmt.Errorf("failed to unmarshal header: %v", err)
	}
	if header.FileName == "" {
		return "", fmt.Errorf("unknown filetype in header: %v", header)
	}
	// Create a new file locally to save the received data
	outFilePath := filepath.Join(settings.Settings.Client.Swapped(), header.FileName)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Receive file data
	_, err = io.CopyN(outFile, s.Conn, header.FileSize)
	if err != nil {
		return "", fmt.Errorf("failed to receive file: %v", err)
	}

	log.Println("File received successfully:", outFilePath)
	return outFilePath, nil
}

// Close cleans up resources.
func (s *Server) Close() {
	if s.Conn != nil {
		s.Conn.Close()
	}
	fmt.Println("Connection closed.")
}

func (s *Server) SendFileWithRetry(ctx context.Context, sendtype string, file string) {
	ok := false
	msg := "RETRY"
	for msg == "RETRY" {
		if err := s.SendFile(sendtype, file); err != nil {
			log.Println("Error sending file:", err)
		}
		ok, msg = s.WaitForDone(ctx, make([]byte, 4096))
		if ok {
			fmt.Println("Finished sending file!")
		}
	}
	if !ok {
		fmt.Println("Error sending file")
	}
}
