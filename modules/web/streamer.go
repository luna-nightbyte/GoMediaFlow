package web

import (
	"fmt"
	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"log"

	"gocv.io/x/gocv"
)

func (s *Server) FrameFeeder() {
	defer s.WG.Done()

	var window *gocv.Window
	if config.Config.ViewLocalStream {
		window := gocv.NewWindow("Camera Feed")
		defer window.Close()
	}
	for frame := range webcam.FrameChan {
		// Encode frame to JPEG
		buf, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, frame.Mat, []int{gocv.IMWriteJpegQuality, 90})
		if err != nil {
			log.Printf("Error encoding frame: %v", err)
			continue
		}

		// Send the encoded frame over the connection
		_, err = s.Conn.Write(buf.GetBytes())
		if err != nil {
			log.Printf("Error sending frame: %v", err)
			break
		}
		if config.Config.ViewLocalStream {

			// Display the frame locally (optional)
			window.IMShow(frame.Mat)
			if window.WaitKey(1) >= 0 {
				break
			}
		}
		fmt.Println("Recieved frame")
		frame.Mat.Close()
	}
	close(webcam.FrameChan)
}
