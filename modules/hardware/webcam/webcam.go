package webcam

import (
	"context"
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

type Frame struct {
	Mat gocv.Mat
}

var FrameChan = make(chan Frame)

// StartFrameChannel starts streaming frames into FrameChan
func StartFrameChannel(ctx context.Context, deviceID int) {
	camera, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}

	buf := gocv.NewMat()

	// Goroutine to capture frames
	go func() {
		for {
			select {
			case <-ctx.Done(): // Stop capturing on context cancellation
				fmt.Println("Stopping frame capture")

				camera.Close()
				buf.Close()
				return
			default:
				if ok := camera.Read(&buf); !ok {
					fmt.Printf("Device closed or read error: %v\n", deviceID)

					continue
				}

				if buf.Empty() {
					continue
				}

				FrameChan <- Frame{Mat: buf.Clone()}
			}
		}

	}()
}

// Example usage
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the frame capture
	go StartFrameChannel(ctx, 0)

	// Simulate processing frames
	go func() {
		for frame := range FrameChan {
			fmt.Println("Received frame!")
			frame.Mat.Close() // Release memory after processing
		}
	}()

	// Run for 10 seconds
	time.Sleep(10 * time.Second)
	cancel()         // Stop frame capture
	close(FrameChan) // Close the channel explicitly
}
