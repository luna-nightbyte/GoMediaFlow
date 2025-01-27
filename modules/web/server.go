package web

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Server struct {
	Conn     net.Conn
	Listener net.Listener
	WG       sync.WaitGroup
	Frames   FrameFeeder
	mux      sync.Mutex
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
