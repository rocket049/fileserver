// main.go
package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func main() {
	var title = flag.String("t", "", "title")
	var fn = flag.String("f", "", "file name")
	flag.Parse()
	if *fn == "" {
		flag.Usage()
		return
	}

	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, core.Qt__Window)
	img := gui.NewQPixmap3(*fn, "", core.Qt__NoFormatConversion)
	label := widgets.NewQLabel(window, core.Qt__Widget)
	label.SetPixmap(img)
	window.SetCentralWidget(label)

	if *title == "" {
		window.SetWindowTitle(filepath.Base(*fn))
	} else {
		window.SetWindowTitle(*title)
	}

	app.SetActiveWindow(window)
	window.Show()
	app.Exec()
}
