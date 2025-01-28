package ui

import (
	"context"
	"fmt"
	"goStreamer/modules/config"
	"goStreamer/modules/hardware/webcam"
	"goStreamer/modules/web"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (ui UICreator) HandleUI(server *web.Server, ctx context.Context, webcam_source int) {

	files, err := os.ReadDir(config.Config.Local.Targetfolder)
	if err != nil {
		fmt.Printf("Error reading target folder: %v\n", err)
		return
	}
	for index := range files {
		// Save target as latest in config
		target := filepath.Join(config.Config.Local.Targetfolder, files[index].Name())
		config.Config.UpdateLastFiles(config.Config.LastSource(), target, config.Config.LastSwapped())

		server.SendFileWithRetry(ctx, "SEND_TARGET", config.Config.LastTarget())
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
			// Save source as latest in config
			source := filepath.Join(config.Config.Local.SourceFolder, files[index].Name())
			config.Config.UpdateLastFiles(source, config.Config.LastTarget(), config.Config.LastSwapped())

			server.SendFileWithRetry(ctx, "SEND_SOURCE", config.Config.LastSource())
		}
	}
	time.Sleep(500 * time.Microsecond)
	server.Conn.Close()

}
