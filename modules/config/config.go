package config

import (
	"goStreamer/modules/db"
	"log"
)

type config struct {
	IP              string `json:"ip"`
	PORT            int    `json:"port"`
	InputSource     string `json:"input_source"`
	InputTarget     string `json:"input_target"`
	OutputFile      string `json:"output_file"`
	ViewLocalStream bool   `json:"show_local_stream"`
	UseWebcam       bool   `json:"use_webcam"`
}

var (
	Config config
	Path   string
)

func (c *config) Init(path string) {
	Path = path
	ok, err := db.Check(Path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return
	}
	err = db.Read(Path, &c)
	if err != nil {
		log.Fatal("Error reading config: ", err)
		return
	}
}
