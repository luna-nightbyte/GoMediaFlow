package main

// /home/thomas/Pictures/FRS/PXL_20240305_090139715.MP.jpg
// /home/thomas/Pictures/FRS/PXL_20240305_090144477.MP.jpg
// /home/thomas/Desktop
//

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/local"
	"goStreamer/modules/ui"
	"goStreamer/modules/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var server web.Server

func init() {
	config.Config.Init("config.json")
	os.MkdirAll(config.Config.Local.SourceFolder, os.ModePerm)
	os.MkdirAll(config.Config.Local.Targetfolder, os.ModePerm)
	os.MkdirAll(config.Config.Local.OutputFolder, os.ModePerm)

}
func waitForDone(ctx context.Context, buffer []byte) bool {

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping fileHandler.")
			return false
		default:
			n, err := server.Conn.Read(buffer)
			if err != nil {
				log.Printf("Error reading from connection: %v\n", err)
				return false
			}
			content := string(buffer[:n])
			if content == "DONE" {
				return true
			}
		}
	}
}
func fileHandler(ctx context.Context, webcam_source int) {
	files, err := os.ReadDir(config.Config.Local.Targetfolder)
	if err != nil {
		fmt.Printf("Error reading target folder: %v\n", err)
		return
	}
	for index := range files {
		if err := server.SendFile("SEND_TARGET", filepath.Join(config.Config.Local.Targetfolder, files[index].Name())); err != nil {
			log.Println("Error sending file:", err)
		}
		ok := waitForDone(ctx, make([]byte, 4096))
		if ok {
			fmt.Println("Finished sending file!")
		}
	}

	if config.Config.Local.Webcam.Enable {
		if webcam_source == -1 && config.Config.Local.Webcam.Target == "-1" {
			log.Println("No source selected for webcam")
			return
		}
		if webcam_source == -1 {
			var err error
			webcam_source, err = strconv.Atoi(config.Config.Local.Webcam.Target)
			if err != nil {
				log.Println("Error setting webcam source from config")
				return
			}
		}
		fmt.Fprintln(server.Conn, "START_FRAMES")
		go webcam.StartFrameChannel(ctx, webcam_source)
		server.WG.Add(1)
		go server.Frames.Start(&server.WG, server.Conn)
		server.WG.Wait()
		fmt.Fprintln(server.Conn, "STOP_FRAMES") // Stop processing frames on server

	} else {
		// Send source file if no webcam is used.

		files, err := os.ReadDir(config.Config.Local.SourceFolder)
		if err != nil {
			fmt.Printf("Error reading source folder: %v\n", err)
			return
		}
		for index := range files {
			if err := server.SendFile("SEND_SOURCE", filepath.Join(config.Config.Local.SourceFolder, files[index].Name())); err != nil {
				log.Println("Error sending file:", err)
			}
			ok := waitForDone(ctx, make([]byte, 4096))
			if ok {
				fmt.Println("Finished sending file!")
			}
		}
	}
	time.Sleep(500 * time.Microsecond)
	server.Conn.Close()

}
func getFile(ctx context.Context) {
	if config.Config.Local.Webcam.Enable {
		return
	}

	// Receive the output file
	fmt.Fprintln(server.Conn, "REQUEST_FILE")
	time.Sleep(500 * time.Microsecond)
	_, err := server.ReceiveFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	ok := waitForDone(ctx, make([]byte, 4096))
	if !ok {
		log.Println("Did not reacieve done msg..")
	}
	fmt.Fprintln(server.Conn, "EXIT")
	time.Sleep(500 * time.Microsecond)
}
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := ui.New("GoStreamer")
	var content *fyne.Container

	if !config.Config.Local.Webcam.Enable { // We expect files then
		sourceEntry, sourceButton := ui.AddFolderSelector("Select a source folder", "Choose a folder...")
		targetEntry, targetButton := ui.AddFolderSelector("Select a target folder", "Choose a folder...")
		outputEntry, outputButton := ui.AddFolderSelector("Select output Folder", "Choose an output folder...")
		sourceEntry.Text = config.Config.Local.SourceFolder
		targetEntry.Text = config.Config.Local.Targetfolder
		outputEntry.Text = config.Config.Local.OutputFolder
		submitButton := ui.AddSubmitButton("Submit", func() {

			// Update files and config
			local.Files.Update(sourceEntry.Text, targetEntry.Text, outputEntry.Text)

			server.Connect(config.Config.Server.IP, config.Config.Server.DialPort)
			defer server.Conn.Close()
			fileHandler(ctx, -1)
		})

		getFileButton := ui.AddSubmitButton("Get swapped", func() {

			server.Connect(config.Config.Server.IP, config.Config.Server.DialPort)
			defer server.Conn.Close()
			getFile(ctx)
		})
		content = container.NewVBox(
			sourceEntry, sourceButton,
			targetEntry, targetButton,
			outputEntry, outputButton,
			submitButton,
			getFileButton,
		)
	} else { // We got webcam
		sourceEntry, sourceButton := ui.AddFileSelector("Select a source face", "Choose a file...")

		webcamTarget := ui.AddOutputFilename("Filename", "Enter webgam target (default is usually 0)")

		submitButton := ui.AddSubmitButton("Submit", func() {

			source, err := strconv.Atoi(webcamTarget.Text)
			if err != nil {
				log.Println("Wrong webcam type!")
			}
			local.Files.UpdateSingle(sourceEntry.Text, webcamTarget.Text)
			fileHandler(ctx, source)
		})

		content = container.NewVBox(
			webcamTarget,
			sourceEntry, sourceButton,
			submitButton,
		)
	}
	// Start UI
	ui.Run(content)

	// Connection to the server running face swapper

}
