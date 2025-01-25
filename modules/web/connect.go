package web

import (
	"fmt"
	"goStreamer/modules/files"
	"log"
	"net"
	"sync"
)

type Server struct {
	Conn  net.Conn
	WG    sync.WaitGroup
	Files files.Output
}

// Connect establishes a connection to the remote server.
// Blocks until connection is established
func (s *Server) Connect(ip string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalf("Failed to connect to remote computer: %v", err)
	}
	s.Conn = conn
	log.Println("Connected to server:", ip, port)
}
