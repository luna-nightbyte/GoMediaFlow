package settings

import (
	"goStreamer/modules/db"
	"log"
)

type settings struct {
	Server server `json:"server"`
	Client client `json:"client"`
}
type server struct {
	Net network `json:"network"`
}
type client struct {
	Net    network `json:"network"`
	Webcam webcam  `json:"webcam"`
	Dir    files   `json:"files"`
}

type network struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type webcam struct {
	Enable bool   `json:"enable"`
	Target string `json:"target"`
}
type files struct {
	Source folder `json:"source"`
	Target folder `json:"target"`
	Output folder `json:"output"`
}
type folder struct {
	Folder string `json:"folder"`
	Last   string `json:"last"`
}

var (
	Settings settings
	Path     string
)

func (c *settings) Init(path string) {
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
func (c *settings) Update() {
	var tmpConfig settings
	tmpConfig.read()
	c.write()
	if !c.verify() {
		tmpConfig.write()
	}
}

func (c *settings) verify() bool {
	ok, err := db.Check(Path, &c)
	if !ok {
		log.Println("Config error: ", err)
		return false
	}
	return true
}

func (c *settings) read() bool {
	err := db.Read(Path, &c)
	if err != nil {
		log.Println("Error reading config: ", err)
		return false
	}
	return true
}
 
func (c *settings) write() bool {
	err := db.Write(Path, &c)
	if err != nil {
		log.Println("Error writing config: ", err)
		return false
	}
	return true
}

func (c *settings) UpdateLastFiles(source, target, swapped string) {
	c.Client.Dir.Source.Last = source
	c.Client.Dir.Target.Last = target
	c.Client.Dir.Output.Last = swapped
	c.Update()
}
