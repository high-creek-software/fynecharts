package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/fynecharts"
	"math/rand"
)

var labels = []string{}
var data = []float64{}

func main() {
	app := app.New()
	window := app.NewWindow("Example Bar")
	window.Resize(fyne.NewSize(300, 150))

	chart := fynecharts.NewBarChart(window.Canvas(),
		"Simple Bar Chart",
		labels,
		data,
	)
	chart.UpdateSuggestedTickCount(8)
	chart.SetXLabel("Days with stuff")
	chart.SetYLabel("Amount of stuff")

	removeBtn := widget.NewButton("Remove", func() {
		if len(data) == 0 {
			return
		}
		data = data[1:]
		labels = labels[1:]
		chart.UpdateData(labels, data)
	})
	addBtn := widget.NewButton("Add", func() {
		r := genRandom()
		data = append(data, r)
		labels = append(labels, fmt.Sprintf("%d", len(data)))
		chart.UpdateData(labels, data)
	})
	grid := container.NewGridWithColumns(2, removeBtn, addBtn)
	recomputeBtn := widget.NewButton("Recompute", func() {
		//chart.UpdateData([]string{}, []float64{})
		for idx := range data {
			data[idx] = genRandom()
		}
		chart.UpdateData(labels, data)
	})
	window.SetContent(container.NewBorder(grid, recomputeBtn, nil, nil, chart))

	window.ShowAndRun()
}

func genRandom() float64 {
	return 10 + rand.Float64()*(100-10)
}
