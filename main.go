package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"tidy/modules/config"
	"tidy/modules/hardware/webcam"

	"gocv.io/x/gocv"
)

func init() {
	config.Config.Init("config.json")

}
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	inputsourceInt, err := strconv.Atoi(config.Config.InputSource)
	if err != nil {
		log.Fatal("Wrong device type in config..")
	}
	go webcam.StartFrameChannel(ctx, inputsourceInt)

	// Create a window to display the video feed (optional)
	window := gocv.NewWindow("Camera Feed")
	defer window.Close()

	// Prepare a connection to the remote computer
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Config.IP, config.Config.PORT))
	if err != nil {
		log.Fatalf("Failed to connect to remote computer: %v", err)
	}
	defer conn.Close()

	for frame := range webcam.FrameChan {
		// Encode frame to JPEG
		buf, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, frame.Mat, []int{gocv.IMWriteJpegQuality, 90})
		if err != nil {
			log.Printf("Error encoding frame: %v", err)
			continue
		}

		// Send the encoded frame over the connection
		_, err = conn.Write(buf.GetBytes())
		if err != nil {
			log.Printf("Error sending frame: %v", err)
			break
		}

		// Display the frame locally (optional)
		window.IMShow(frame.Mat)
		if window.WaitKey(1) >= 0 {
			break
		}
		fmt.Println("Recieved frame")
		frame.Mat.Close()
	}

	close(webcam.FrameChan) // Close the channel explicitly

}
