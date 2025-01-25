package web

import (
	"goStreamer/modules/config"
	"goStreamer/modules/db"
)

type output struct {
	source file
	target file
	output file
}
type file struct {
	file string
}

func (o *output) Update(source, target, output string) {
	o.source.file = source
	o.target.file = target
	o.output.file = output
	config.Config.InputSource = o.source.file
	config.Config.InputTarget = o.target.file
	config.Config.OutputFile = o.output.file
	db.Write("config.json", config.Config)
}

func (o *output) Recieve() {

}
func (o *output) Send() {

}
