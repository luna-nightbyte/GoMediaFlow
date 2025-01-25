package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type UICreator struct {
	App    fyne.App
	Window fyne.Window
}

func New(title string) *UICreator {
	app := app.New()
	window := app.NewWindow(title)
	return &UICreator{
		App:    app,
		Window: window,
	}
}

func (ui *UICreator) AddFileSelector(label string, placeholder string) (*widget.Entry, *widget.Button) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	button := widget.NewButton(label, func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			if reader == nil {
				return
			}
			entry.SetText(reader.URI().Path())
		}, ui.Window)
	})
	return entry, button
}

func (ui *UICreator) AddFolderSelector(label string, placeholder string) (*widget.Entry, *widget.Button) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	button := widget.NewButton(label, func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			if uri == nil {
				return
			}
			entry.SetText(uri.Path())
		}, ui.Window)
	})
	return entry, button
}

func (ui *UICreator) AddOutputSelector(label string, placeholder string) (*widget.Entry, *widget.Button) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	button := widget.NewButton(label, func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			if uri == nil {
				return
			}
			entry.SetText(uri.Path())
		}, ui.Window)
	})
	return entry, button
}

func (ui *UICreator) AddOutputFilename(label string, placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	return entry
}

func (ui *UICreator) AddSubmitButton(label string, callback func()) *widget.Button {
	return widget.NewButton(label, callback)
}

func (ui *UICreator) Run(content fyne.CanvasObject) {
	ui.Window.SetContent(content)
	ui.Window.Resize(fyne.NewSize(500, 400))
	ui.Window.ShowAndRun()
}
