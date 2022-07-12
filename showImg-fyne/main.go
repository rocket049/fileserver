package main

import (
	"flag"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	t := flag.String("t", "", "title of window")
	f := flag.String("f", "", "picture file name")
	flag.Parse()
	p := flag.Arg(0)
	if p == "" {
		p = *f
		if p == "" {
			panic("need picture file")
		}
	}
	app := app.New()

	icon1 := canvas.NewImageFromFile(p)

	icon1.FillMode = canvas.ImageFillOriginal

	box1 := container.NewVBox(icon1)

	var title string
	if *t != "" {
		title = *t
	} else {
		title = filepath.Base(p)
	}
	w := app.NewWindow(title)
	w.SetContent(box1)

	w.SetOnClosed(func() {
		app.Quit()
	})
	w.ShowAndRun()
}

func ErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
