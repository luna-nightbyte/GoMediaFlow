package main

import (
	"context"
	"strconv"

	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/ui"
	"goStreamer/modules/web"

	"fyne.io/fyne"
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

	sourceEntry, sourceButton := ui.AddFileSelector("Select a source face", "Choose a file...")

	// Check if input is webcam int
	inputsourceInt, err := strconv.Atoi(config.Config.InputSource)
	if err != nil { // We expect files then
		targetEntry, targetButton := ui.AddFileSelector("Select a target face", "Choose a file...")
		outputEntry, outputButton := ui.AddOutputSelector("Select Output Folder", "Choose an output folder...")
		outputNameEntry := ui.AddOutputFilename("Filename", "Enter filename...")

		submitButton := ui.AddSubmitButton("Submit", func() {
			println("Source:", sourceEntry.Text)
			println("Target:", targetEntry.Text)
			println("Output folder:", outputEntry.Text)
			println("Output filename:", outputNameEntry.Text)
		})
		content = container.NewVBox(
			sourceEntry, sourceButton,
			targetEntry, targetButton,
			outputEntry, outputButton,
			outputNameEntry,
			submitButton,
		)
	} else { // We got webcam
		submitButton := ui.AddSubmitButton("Submit", func() {
			println("Source:", sourceEntry.Text)
		})
		content = container.NewVBox(
			sourceEntry, sourceButton,
			submitButton,
		)
		go webcam.StartFrameChannel(ctx, inputsourceInt)
		// ready up handling target and output files

	}

	// Start UI
	ui.Run(content)

	server.Files.Update("", "", "")
	// Connection to the server running face swapper
	server.Connect()
	defer server.Conn.Close()
	server.WG.Add(1)
	go server.FrameFeeder()
	server.WG.Wait() // Wait for connection to close
}
