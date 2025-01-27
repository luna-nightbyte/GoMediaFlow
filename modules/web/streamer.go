package web

import (
	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"log"
	"net"
	"sync"

	"gocv.io/x/gocv"
)

type FrameFeeder struct {
	Running bool
}

func (f *FrameFeeder) Start(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer f.cleanup(conn)

	var window *gocv.Window
	if config.Config.Local.Webcam.Enable {
		window = gocv.NewWindow("Camera Feed")
		defer window.Close()
	}

	log.Println("Frame feeder started.")

	for frame := range webcam.FrameChan {
		// Encode frame to JPEG
		buf, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, frame.Mat, []int{gocv.IMWriteJpegQuality, 90})
		if err != nil {
			log.Printf("Error encoding frame: %v", err)
			continue
		}

		// Send the encoded frame over the connection
		if _, err = conn.Write(buf.GetBytes()); err != nil {
			log.Printf("Error sending frame: %v", err)
			break
		}

		// Optionally display the frame
		if config.Config.Local.Webcam.Enable && window != nil {
			window.IMShow(frame.Mat)
			if window.WaitKey(1) >= 0 {
				break
			}
		}

		frame.Mat.Close()
	}
	log.Println("Frame feeder stopped.")
}

func (f *FrameFeeder) cleanup(conn net.Conn) {
	close(webcam.FrameChan)
	//if err := conn.Close(); err != nil {
	//	log.Printf("Error closing connection: %v", err)
	//}
}
