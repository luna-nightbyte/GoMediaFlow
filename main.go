package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"goStreamer/modules/config"
	"goStreamer/modules/files"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/ui"
	"goStreamer/modules/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func init() {
	config.Config.Init("config.json")

}
func main() {

	var server web.Server

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := ui.New("GoStreamer")
	var content *fyne.Container

	// Default webcam source
	webcam_source := 0

	if !config.Config.UseWebcam { // We expect files then

		sourceEntry, sourceButton := ui.AddFileSelector("Select a source face", "Choose a file...")
		targetEntry, targetButton := ui.AddFileSelector("Select a target face", "Choose a file...")
		outputEntry, outputButton := ui.AddOutputSelector("Select Output Folder", "Choose an output folder...")
		outputNameEntry := ui.AddOutputFilename("Filename", "Enter filename...")

		submitButton := ui.AddSubmitButton("Submit", func() {

			output_path := filepath.Join(outputEntry.Text, outputNameEntry.Text)
			if !files.IsFileAndExist(sourceEntry.Text, "image") {
				log.Fatal("Wrong input source type..")
			}
			if !files.IsFileAndExist(targetEntry.Text, "image") && !files.IsFileAndExist(targetEntry.Text, "video") {
				log.Fatal("Wrong input target type..")
			}
			if !files.IsVideoOrImageFileName(output_path) && !files.IsVideoOrImageFileName(output_path) {
				log.Fatal("Wrong output type..")
			}

			// Update files and config
			server.Files.Update(sourceEntry.Text, targetEntry.Text, filepath.Join(outputNameEntry.Text, outputNameEntry.Text))
		})
		content = container.NewVBox(
			sourceEntry, sourceButton,
			targetEntry, targetButton,
			outputEntry, outputButton,
			outputNameEntry,
			submitButton,
		)
	} else { // We got webcam
		sourceEntry, sourceButton := ui.AddFileSelector("Select a source face", "Choose a file...")

		webcamTarget := ui.AddOutputFilename("Filename", "Enter webgam target (default is usually 0)")

		submitButton := ui.AddSubmitButton("Submit", func() {

			source, err := strconv.Atoi(webcamTarget.Text)
			if err != nil {
				log.Fatal("Wrong webcam type!")
			}
			server.Files.UpdateSingle(sourceEntry.Text, webcamTarget.Text)
			webcam_source = source
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
	server.Connect(config.Config.IP, config.Config.PORT)
	defer server.Conn.Close()

	fmt.Fprintln(server.Conn, "SEND_SOURCE")
	server.Send(server.Files.Source())

	if config.Config.UseWebcam {

		fmt.Fprintln(server.Conn, "START_FRAMES")
		go webcam.StartFrameChannel(ctx, webcam_source)
		server.WG.Add(1)
		go server.FrameFeeder()
		server.WG.Wait()
		fmt.Fprintln(server.Conn, "STOP_FRAMES") // Stop processing frames on server
	} else {
		// Send target file if no webcam is used.
		fmt.Fprintln(server.Conn, "SEND_TARGET")
		server.Send(server.Files.Target())
	}
	// Receive the output file
	fmt.Fprintln(server.Conn, "REQUEST_FILE")
	server.Recieve()

	fmt.Fprintln(server.Conn, "EXIT")
}
