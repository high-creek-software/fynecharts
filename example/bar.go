package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/high-creek-software/fynecharts"
)

func main() {
	app := app.New()
	window := app.NewWindow("Example Bar")

	chart := fynecharts.NewBarChart(window.Canvas(), "Simple Bar Chart", []string{"One", "Two", "Three", "Four"}, []float64{25, 34, 45, 10})

	window.SetContent(chart)

	window.ShowAndRun()
}
