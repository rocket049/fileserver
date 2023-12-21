package main

import (
	"flag"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	cmd1 := exec.Command("sync", p)
	cmd1.Run()
	cmd1.Wait()

	app := app.New()

	icon1 := canvas.NewImageFromFile(p)

	icon1.FillMode = canvas.ImageFillOriginal

	viewer := widget.NewCard("", "", icon1)

	btReload := widget.NewButton("Reload", func() {
		icon2 := canvas.NewImageFromFile(p)
		icon2.FillMode = canvas.ImageFillOriginal
		viewer.SetContent(icon2)
		viewer.Refresh()
	})

	box1 := container.NewVBox(btReload, viewer)

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
