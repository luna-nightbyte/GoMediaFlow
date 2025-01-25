package web

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Send uploads the source and target files to the server.
func (s *Server) Send() {
	filesToSend := []string{s.Files.Source(), s.Files.Target()}

	for _, filePath := range filesToSend {
		log.Printf("Sending file: %s\n", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		// Write the filename first
		filename := filepath.Base(filePath)
		_, err = s.Conn.Write([]byte(fmt.Sprintf("%s\n", filename)))
		if err != nil {
			log.Fatalf("Failed to send filename: %v", err)
		}

		// Send file content
		_, err = io.Copy(s.Conn, file)
		if err != nil {
			log.Fatalf("Failed to send file content: %v", err)
		}
		log.Printf("File %s sent successfully.\n", filename)
	}
}

// Receive downloads the output file from the server.
func (s *Server) Recieve() {
	outputPath := s.Files.Output()
	log.Printf("Receiving output file: %s\n", outputPath)

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// Receive file content
	_, err = io.Copy(file, s.Conn)
	if err != nil {
		log.Fatalf("Failed to receive file: %v", err)
	}
	log.Printf("Output file received successfully: %s\n", outputPath)
}
