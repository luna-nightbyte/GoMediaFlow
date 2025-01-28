package web

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// Server represents the server side of the web module.
// This is when we initially sends the startup arguments and files.
// Optionally, the server will then send the webcam feed to the client if enabled.
type Server struct {
	Conn     net.Conn
	Listener net.Listener
	WG       sync.WaitGroup
	Frames   FrameFeeder
	mux      sync.Mutex
	Ready    bool
}

// Connect establishes a connection to the remote server.
func (s *Server) Connect(ip string, port int) {
	fmt.Println("Connecting to server...")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalf("Failed to connect to remote computer: %v", err)
	}
	s.Conn = conn

	log.Println("Connected to server:", ip, port)
}

func (s *Server) ListenAndAccept(port int) {

	s.CloseConnection() // Ensure last connection is closed

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	fmt.Printf("Server is running on port %s...", port)
	s.Conn, err = ln.Accept()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

}

func (s *Server) isClosed() bool {
	return s.Conn == nil
}
func (s *Server) CloseConnection() {
	if s.isClosed() {
		return
	}
	s.Conn.Close()
}

func (c *Server) SendMessage(message string) error {
	_, err := c.Conn.Write([]byte(message + "\n"))
	return err
}
func (s *Server) WaitForDone(ctx context.Context, buffer []byte) (bool, string) {

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping fileHandler.")
			return false, ""
		default:
			n, err := s.Conn.Read(buffer)
			if err != nil {
				log.Printf("Error reading from connection: %v\n", err)
				return false, ""
			}
			content := string(buffer[:n])
			if content == "DONE" {
				return true, content
			}
			return false, content
		}
	}
}
