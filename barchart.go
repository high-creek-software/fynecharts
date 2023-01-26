package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
)

const (
	defaultBarWidth  = 25
	defaultMinHeight = 100
)

type BarChart struct {
	widget.BaseWidget
	canvas fyne.Canvas

	title  string
	yTitle string
	xTitle string
	labels []string
	data   []float64

	barWidth  float32
	minHeight float32
}

func (b *BarChart) CreateRenderer() fyne.WidgetRenderer {
	titleLbl := widget.NewLabel(b.title)
	titleLbl.Hide()
	yLbl := widget.NewLabel(b.yTitle)
	yLbl.Hide()
	ySep := canvas.NewLine(theme.ForegroundColor())
	ySep.StrokeWidth = 1
	xLbl := widget.NewLabel(b.xTitle)
	xLbl.Hide()
	xSep := canvas.NewLine(theme.ForegroundColor())
	xSep.StrokeWidth = 1

	return &barChartRenderer{barChart: b, titleLbl: titleLbl, yLbl: yLbl, xLbl: xLbl, xSeparator: xSep, ySeparator: ySep}
}

func (b *BarChart) UpdateData(labels []string, data []float64) {
	b.labels = labels
	b.data = data
	b.Refresh()
}

func (b *BarChart) SetBarWidth(w float32) {
	b.barWidth = w
}

func (b *BarChart) SetMinHeight(h float32) {
	b.minHeight = h
}

func (b *BarChart) SetXLabel(xlbl string) {
	b.xTitle = xlbl
	b.Refresh()
}

func NewBarChart(canvas fyne.Canvas, title string, labels []string, data []float64) *BarChart {
	bc := &BarChart{canvas: canvas, title: title, labels: labels, data: data, barWidth: defaultBarWidth, minHeight: defaultMinHeight}
	bc.ExtendBaseWidget(bc)
	bc.Refresh()

	return bc
}

type barChartRenderer struct {
	barChart   *BarChart
	titleLbl   *widget.Label
	yLbl       *widget.Label
	ySeparator *canvas.Line
	xLbl       *widget.Label
	xSeparator *canvas.Line

	labels []*widget.Label
	data   []*bar

	xLblMax   fyne.Size
	dataMax   float64
	dataMin   float64
	dataRange float64

	barWidth float32
	barScale float32
}

func (b *barChartRenderer) Destroy() {

}

func (b *barChartRenderer) Layout(size fyne.Size) {

	titleSize := b.titleLbl.MinSize()
	//ySize := b.yLbl.MinSize()
	xSize := fyne.NewSize(0, 0)
	if b.xLbl.Visible() {
		xSize = b.xLbl.MinSize()
	}

	titleX := size.Width/2 - titleSize.Width/2
	titlePos := fyne.NewPos(titleX, theme.Padding())
	b.titleLbl.Move(titlePos)
	b.titleLbl.Resize(titleSize)

	xLblX := size.Width/2 - xSize.Width/2
	xPos := fyne.NewPos(xLblX, size.Height-xSize.Height-theme.Padding())
	b.xLbl.Move(xPos)

	xOffset := float32(0)

	xSepY := size.Height - xSize.Height - b.xLblMax.Height - 5
	b.xSeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.xSeparator.Position2 = fyne.NewPos(size.Width-theme.Padding(), xSepY)

	b.ySeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.ySeparator.Position2 = fyne.NewPos(xOffset, titleSize.Height+2*theme.Padding())

	availableHeight := size.Height - titleSize.Height - theme.Padding() - xSize.Height - b.xLblMax.Height - theme.Padding()
	columnWidth := (size.Width - (2 * theme.Padding())) / float32(len(b.barChart.labels))

	if len(b.labels) > 0 {
		for idx := range b.barChart.labels {
			lbl := b.labels[idx]
			lblSize := lbl.MinSize()
			xCellOffset := float32(idx) * columnWidth
			lblPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-lblSize.Width/2,
				size.Height-2*theme.Padding()-xSize.Height-float32(lblSize.Height))
			lbl.Move(lblPos)
		}
	}

	if len(b.data) > 0 {
		for idx, d := range b.barChart.data {
			rect := b.data[idx]
			scale := float32((d - b.dataMin) / (b.dataMax - b.dataMin))
			rect.Resize(fyne.NewSize(b.barChart.barWidth, availableHeight*scale))
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-rect.Size().Width/2,
				size.Height-xSize.Height-(availableHeight*scale)-b.xLblMax.Height-theme.Padding())
			rect.Move(rectPos)
		}
	}
}

func (b *barChartRenderer) MinSize() fyne.Size {
	b.barScale = 10
	titleSize := fyne.NewSize(0, 0)
	paddingCount := 0
	if b.titleLbl.Visible() {
		titleSize = b.titleLbl.MinSize()
		paddingCount++
		paddingCount++
	}
	xLblSize := fyne.NewSize(0, 0)
	if b.xLbl.Visible() {
		xLblSize = b.xLbl.MinSize()
		paddingCount++
	}

	xCellWidth := b.xLblMax.Width + 2
	return fyne.NewSize(float32(len(b.barChart.labels))*xCellWidth,
		titleSize.Height+xLblSize.Height+b.xLblMax.Height+float32(paddingCount)*theme.Padding()+b.barChart.minHeight)
}

func (b *barChartRenderer) Objects() []fyne.CanvasObject {
	cos := []fyne.CanvasObject{b.titleLbl, b.yLbl, b.xLbl}
	for _, lbl := range b.labels {
		cos = append(cos, lbl)
	}

	for _, d := range b.data {
		cos = append(cos, d)
	}

	if b.xSeparator != nil {
		cos = append(cos, b.xSeparator)
	}

	if b.ySeparator != nil {
		cos = append(cos, b.ySeparator)
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

	for idx, datum := range b.barChart.data {
		if idx >= len(b.data) {
			b.dataMax = math.Max(b.dataMax, datum)
			b.dataMin = math.Min(b.dataMin, datum)
			b.data = append(b.data, newBar(b.barChart.canvas, datum))
		} else {
			b.data[idx].value = datum
		}
	}

	b.dataRange = b.dataMax - b.dataMin

	for idx := range b.barChart.labels {
		var lbl *widget.Label
		if idx >= len(b.labels) {
			lbl = widget.NewLabel(b.barChart.labels[idx])
			b.xLblMax = b.xLblMax.Max(lbl.MinSize())
			b.labels = append(b.labels, lbl)
		} else {
			lbl = b.labels[idx]
			lbl.SetText(b.barChart.labels[idx])
			b.xLblMax = b.xLblMax.Max(lbl.MinSize())
		}
	}

}
