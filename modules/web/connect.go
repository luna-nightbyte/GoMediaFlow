package web

import (
	"fmt"
	"goStreamer/modules/config"
	"log"
	"net"
	"sync"
)

type Server struct {
	Conn  net.Conn
	WG    sync.WaitGroup
	Files output
}

func (s *Server) Connect() {
	// Prepare a connection to the remote computer
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Config.IP, config.Config.PORT))
	if err != nil {
		log.Fatalf("Failed to connect to remote computer: %v", err)
	}
	s.Conn = conn

}
