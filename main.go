package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"goStreamer/modules/local"
	"goStreamer/modules/settings"
	"goStreamer/modules/ui"
	"goStreamer/modules/web"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var server web.Server

func init() {
	settings.Settings.Init("settings.json")
	os.MkdirAll(settings.Settings.Client.Source(), os.ModePerm)
	os.MkdirAll(settings.Settings.Client.Target(), os.ModePerm)
	os.MkdirAll(settings.Settings.Client.Swapped(), os.ModePerm)

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := ui.New("GoStreamer")
	var content *fyne.Container

	if !settings.Settings.Client.Webcam.Enable { // We expect files then
		sourceEntry, sourceButton := ui.AddFolderSelector("Select a source folder", "Choose a folder...")
		targetEntry, targetButton := ui.AddFolderSelector("Select a target folder", "Choose a folder...")
		outputEntry, outputButton := ui.AddFolderSelector("Select output Folder", "Choose an output folder...")
		sourceEntry.Text = settings.Settings.Client.Source()
		targetEntry.Text = settings.Settings.Client.Target()
		outputEntry.Text = settings.Settings.Client.Swapped()
		submitButton := ui.AddSubmitButton("Submit", func() {

			// Update files and config
			local.Files.Update(sourceEntry.Text, targetEntry.Text, outputEntry.Text)

			server.Connect(settings.Settings.Server.Net.IP, settings.Settings.Server.Net.Port)
			defer server.Conn.Close()
			ui.HandleUI(&server, ctx, -1)
		})

		getFileButton := ui.AddSubmitButton("Get swapped", func() {

			server.Connect(settings.Settings.Server.Net.IP, settings.Settings.Server.Net.Port)
			defer server.Conn.Close()
			err := server.GetFile(ctx)
			if err != nil {
				log.Println("Error getting file: ", err)
			}
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
			ui.HandleUI(&server, ctx, source)
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
