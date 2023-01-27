package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/high-creek-software/fynecharts"
)

func main() {
	app := app.New()
	window := app.NewWindow("Example Bar")
	window.Resize(fyne.NewSize(300, 150))

	chart := fynecharts.NewBarChart(window.Canvas(),
		"Simple Bar Chart",
		[]string{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"},
		[]float64{25, 34, 45, 10, 20, 32, 56, 10, 2, 42},
	)

	window.SetContent(chart)

	window.ShowAndRun()
}
