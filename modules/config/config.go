package config

import (
	"goStreamer/modules/db"
	"log"
)

type config struct {
	Server server `json:"server"`
	Local  local  `json:"local"`
}

type server struct {
	IP       string `json:"ip"`
	DialPort int    `json:"port"`
}
type local struct {
	Webcam       webcam `json:"webcam"`
	SourceFolder string `json:"source_folder"`
	Targetfolder string `json:"target_folder"`
	OutputFolder string `json:"output_folder"`
	LastSource   string `json:"last_source"`
	LastTarget   string `json:"last_target"`
	LastSwapped  string `json:"last_swapped"`
}
type webcam struct {
	Enable bool   `json:"enable"`
	Target string `json:"target"`
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
func (c *config) Update() {
	var tmpConfig config
	tmpConfig.read()
	c.write()
	if !c.verify() {
		tmpConfig.write()
	}
}

func (c *config) verify() bool {
	ok, err := db.Check(Path, &c)
	if !ok {
		log.Fatal("Config error: ", err)
		return false
	}
	return true
}

func (c *config) read() bool {
	err := db.Read(Path, &c)
	if err != nil {
		log.Println("Error reading config: ", err)
		return false
	}
	return true
}

func (c *config) write() bool {
	err := db.Write(Path, &c)
	if err != nil {
		log.Println("Error writing config: ", err)
		return false
	}
	return true
}

func (c *config) UpdateLastFiles(source, target, swapped string) {
	c.Local.LastSource = source
	c.Local.LastTarget = target
	c.Local.LastSwapped = swapped
	c.Update()
}
func (c *config) LastSource() string {
	return c.Local.LastSource
}
func (c *config) LastTarget() string {
	return c.Local.LastTarget
}
func (c *config) LastSwapped() string {
	return c.Local.LastSwapped
}
