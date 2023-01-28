package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/high-creek-software/fynecharts"
)

func main() {
	app := app.New()
	window := app.NewWindow("Time Chart")
	window.Resize(fyne.NewSize(300, 150))

	chart := fynecharts.NewTimeSeriesChart(window.Canvas(),
		"Simple Time Series",
		[]string{"Jan. 12, 2023", "Jan. 13, 2023", "Jan. 14, 2023", "Jan. 15, 2023", "Jan. 16, 2023", "Jan. 17, 2023"},
		[]float64{12.3, 19.8, 9.8, 13.5, 56, 4},
	)
	//chart.SetXLabel("Days with stuff")
	window.SetContent(chart)
	window.ShowAndRun()
}
