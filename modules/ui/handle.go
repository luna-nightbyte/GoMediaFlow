package ui

import (
	"context"
	"fmt"
	"goStreamer/modules/settings"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/web"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (ui UICreator) HandleUI(server *web.Server, ctx context.Context, webcam_source int) {

	files, err := os.ReadDir(settings.Settings.Client.Target())
	if err != nil {
		fmt.Printf("Error reading target folder: %v\n", err)
		return
	}
	for index := range files {
		// Save target as latest in config
		target := filepath.Join(settings.Settings.Client.Target(), files[index].Name())
		settings.Settings.UpdateLastFiles(settings.Settings.Client.LastSource(), target, settings.Settings.Client.LastSwapped())

		server.SendFileWithRetry(ctx, "SEND_TARGET", settings.Settings.Client.LastTarget())
	}

	if settings.Settings.Client.Webcam.Enable {
		if webcam_source == -1 && settings.Settings.Client.Webcam.Target == "-1" {
			log.Println("No source selected for webcam")
			return
		}
		if webcam_source == -1 {
			var err error
			webcam_source, err = strconv.Atoi(settings.Settings.Client.Webcam.Target)
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

		files, err := os.ReadDir(settings.Settings.Client.Source())
		if err != nil {
			fmt.Printf("Error reading source folder: %v\n", err)
			return
		}
		for index := range files {
			// Save source as latest in config
			source := filepath.Join(settings.Settings.Client.Source(), files[index].Name())
			settings.Settings.UpdateLastFiles(source, settings.Settings.Client.LastTarget(), settings.Settings.Client.LastSwapped())

			server.SendFileWithRetry(ctx, "SEND_SOURCE", settings.Settings.Client.LastSource())
		}
	}
	time.Sleep(500 * time.Microsecond)
	server.Conn.Close()

}
