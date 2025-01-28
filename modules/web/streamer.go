package web

import (
	"log"
	"net"
	"sync"

	"gocv.io/x/gocv"

	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/settings"
)

type FrameFeeder struct {
	Running bool
}

func (f *FrameFeeder) Start(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer f.cleanup()

	var window *gocv.Window
	if settings.Settings.Client.Webcam.Enable {
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
		if settings.Settings.Client.Webcam.Enable && window != nil {
			window.IMShow(frame.Mat)
			if window.WaitKey(1) >= 0 {
				break
			}
		}

		frame.Mat.Close()
	}
	log.Println("Frame feeder stopped.")
}

func (f *FrameFeeder) cleanup() {
	close(webcam.FrameChan)
}
