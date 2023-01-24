package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BarChart struct {
	widget.BaseWidget

	title  string
	yTitle string
	xTitle string
	labels []string
	data   []float64
}

func (b *BarChart) CreateRenderer() fyne.WidgetRenderer {
	titleLbl := widget.NewLabel(b.title)
	titleLbl.Hide()
	yLbl := widget.NewLabel(b.yTitle)
	yLbl.Hide()
	xLbl := widget.NewLabel(b.xTitle)
	xLbl.Hide()

	return &barChartRenderer{barChart: b, titleLbl: titleLbl, yLbl: yLbl, xLbl: xLbl}
}

func (b *BarChart) UpdateData(labels []string, data []float64) {
	b.labels = labels
	b.data = data
	b.Refresh()
}

func NewBarChart(title string, labels []string, data []float64) *BarChart {
	bc := &BarChart{title: title, labels: labels, data: data}
	bc.ExtendBaseWidget(bc)
	bc.Refresh()

	return bc
}

type barChartRenderer struct {
	barChart *BarChart
	titleLbl *widget.Label
	yLbl     *widget.Label
	xLbl     *widget.Label

	labels []*widget.Label
	data   []*canvas.Rectangle
}

func (b *barChartRenderer) Destroy() {

}

func (b *barChartRenderer) Layout(size fyne.Size) {
	titleSize := b.titleLbl.MinSize()
	//ySize := b.yLbl.MinSize()
	xSize := b.xLbl.MinSize()

	titleX := size.Width/2 - titleSize.Width/2
	titlePos := fyne.NewPos(titleX, theme.Padding())
	b.titleLbl.Move(titlePos)
	b.titleLbl.Resize(titleSize)

	xLblX := size.Width/2 - xSize.Width/2
	xPos := fyne.NewPos(xLblX, size.Height-xSize.Height-theme.Padding())
	b.xLbl.Move(xPos)

	lblMaxSize := fyne.NewSize(0, 0)
	if len(b.labels) > 0 {
		for idx := range b.barChart.labels {
			lbl := b.labels[idx]
			lblSize := lbl.MinSize()
			lblMaxSize = lblMaxSize.Max(lblSize)
			lblPos := fyne.NewPos(float32(idx*35), size.Height-2*theme.Padding()-xSize.Height-float32(lblSize.Height))
			lbl.Move(lblPos)
		}
	}

	if len(b.data) > 0 {
		for idx, d := range b.barChart.data {
			rect := b.data[idx]
			rect.Resize(fyne.NewSize(25, float32(d*2)))
			rectPos := fyne.NewPos(float32(idx*35), size.Height-2*theme.Padding()-xSize.Height-float32(d*2)-lblMaxSize.Height-theme.Padding())
			rect.Move(rectPos)
		}
	}
}

func (b *barChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 200)
}

func (b *barChartRenderer) Objects() []fyne.CanvasObject {
	cos := []fyne.CanvasObject{b.titleLbl, b.yLbl, b.xLbl}
	for _, lbl := range b.labels {
		cos = append(cos, lbl)
	}

	for _, d := range b.data {
		cos = append(cos, d)
	}

	return cos
}

func (b *barChartRenderer) Refresh() {
	if b.barChart.title != "" {
		b.titleLbl.SetText(b.barChart.title)
		b.titleLbl.Show()
	} else {
		b.titleLbl.Hide()
	}

	if b.barChart.yTitle != "" {
		b.yLbl.SetText(b.barChart.yTitle)
		b.yLbl.Show()
	} else {
		b.yLbl.Hide()
	}

	if b.barChart.xTitle != "" {
		b.xLbl.SetText(b.barChart.xTitle)
		b.xLbl.Show()
	} else {
		b.xLbl.Hide()
	}

	for idx := range b.barChart.data {
		if idx >= len(b.data) {
			b.data = append(b.data, canvas.NewRectangle(theme.PrimaryColor()))
		}
	}

	for idx := range b.barChart.labels {
		var lbl *widget.Label
		if idx >= len(b.labels) {
			lbl = widget.NewLabel(b.barChart.labels[idx])
			b.labels = append(b.labels, lbl)
		} else {
			lbl = b.labels[idx]
			lbl.SetText(b.barChart.labels[idx])
		}
	}

}
