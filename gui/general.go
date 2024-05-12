package gui

import (
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const HASH_LOG string = "hash_sum.log"

var file string = filepath.Join(".", HASH_LOG)
var folder string = filepath.Join(".")

func StartGui(hashFile *string, pathFolder *string) {
	a := app.New()
	w := a.NewWindow("Hello")

	clock := widget.NewLabel("Timer")
	start := widget.NewButton("Start", func() {
		*hashFile = file
		*pathFolder = folder
		w.Close()
	})

	w.SetContent(container.NewVBox(
		clock,
		formHash(&w),
		formFolder(&w),
		start,
	))

	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	w.ShowAndRun()

}

func updateTime(clock *widget.Label) {
	formater := time.Now().Format("Time: 03:04:05")
	clock.SetText(formater)
}

func formFolder(w *fyne.Window) *widget.Form {
	entryPathProjectFolder := widget.NewEntry()
	entryPathProjectFolder.SetText("<Path\\to\\project\\folder>")

	selectProjectFolder := func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err == nil && lu != nil {
				entryPathProjectFolder.SetText(lu.Path())
				folder = lu.Path()
			}
		}, *w)
	}

	formSelectProjectFolder := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Path to scan folder",
				Widget: entryPathProjectFolder,
			},
		},
		OnSubmit: selectProjectFolder,
	}

	return formSelectProjectFolder
}

func formHash(w *fyne.Window) *widget.Form {
	entryPathHash := widget.NewEntry()
	entryPathHash.SetText("<Path\\to\\logs>")

	selectFile := func() {
		dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
			if err == nil && uri != nil {
				entryPathHash.SetText(uri.URI().Path())
				file = uri.URI().Path()
			}
		}, *w)
	}

	formSelectHashFile := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Path to hash summ file",
				Widget: entryPathHash,
			},
		},
		OnSubmit: selectFile,
	}

	return formSelectHashFile
}
