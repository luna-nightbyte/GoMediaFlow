package config

import (
	"log"
	"tidy/modules/db"
)

type config struct {
	IP          string `json:"ip"`
	PORT        int    `json:"port"`
	InputSource string `json:"input_source"`
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
