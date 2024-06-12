package ui

import (
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
}

func NewWindow() *Window {
	return &Window{}
}

func (obj Window) ShowAndRun() {
	app := fyneApp.New()
	w := app.NewWindow("Hello")
	size := fyne.Size{Width: 500}
	w.Resize(size)

	hello := widget.NewLabel("Hello Fyne!")

	dialogModal := dialog.NewFileOpen(func(c fyne.URIReadCloser, err error) {
		//path := c.URI().Path()
		//keyReader := CreateKeyReader()
		//if keyReader.isZip(path) {
		//	key, _ := keyReader.OpenZip(path, "pass")
		//	hello.SetText("Welcome zip :)" + key.path)
		//} else {
		//	key, _ := keyReader.OpenDir(path)
		//	hello.SetText("Welcome dir :)" + key.path)
		//}

		hello.SetText("Welcome :)")
	}, w)

	button := widget.NewButton("Выбрать файл ключа", func() {
		dialogModal.Show()
	})

	w.SetContent(container.NewVBox(hello, button))

	w.ShowAndRun()
}
