package fynecharts

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
	"log"
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
	ySep.StrokeWidth = 2
	xLbl := widget.NewLabel(b.xTitle)
	xLbl.Hide()
	xSep := canvas.NewLine(theme.ForegroundColor())
	xSep.StrokeWidth = 2

	return &barChartRenderer{barChart: b, titleLbl: titleLbl, yLbl: yLbl, xLbl: xLbl, xSeparator: xSep, ySeparator: ySep, yLabelPositions: make(map[*widget.Label]float64)}
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

	yLabels         []*widget.Label
	yLabelPositions map[*widget.Label]float64
	xLabels         []*widget.Label
	data            []*bar

	yLblMax fyne.Size
	xLblMax fyne.Size

	yAxis axis

	barWidth float32
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

	xOffset := float32(theme.Padding() + b.yLblMax.Width)

	xSepY := size.Height - xSize.Height - b.xLblMax.Height - 5
	b.xSeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.xSeparator.Position2 = fyne.NewPos(size.Width-theme.Padding(), xSepY)

	b.ySeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.ySeparator.Position2 = fyne.NewPos(xOffset, titleSize.Height+2*theme.Padding())

	availableHeight := size.Height - titleSize.Height - theme.Padding() - xSize.Height - b.xLblMax.Height - theme.Padding()
	columnWidth := (size.Width - xOffset - theme.Padding()) / float32(len(b.barChart.labels))

	if len(b.barChart.labels) > 0 {
		for idx := range b.barChart.labels {
			lbl := b.xLabels[idx]
			lblSize := lbl.MinSize()
			xCellOffset := float32(idx) * columnWidth
			lblPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-lblSize.Width/2,
				size.Height-2*theme.Padding()-xSize.Height-lblSize.Height)
			lbl.Move(lblPos)
		}
	}

	if len(b.yLabelPositions) > 0 {
		for lbl, y := range b.yLabelPositions {
			lblSize := lbl.MinSize()
			scale := b.yAxis.normalize(y)
			pos := fyne.NewPos(0, size.Height-xSize.Height-b.xLblMax.Height-theme.Padding()-(availableHeight*scale)-lblSize.Height/2)
			if pos.Y < titleSize.Height+2*theme.Padding() {
				lbl.Hide()
				continue
			}
			lbl.Move(pos)
		}
	}

	if len(b.data) > 0 {
		for idx, d := range b.barChart.data {
			bar := b.data[idx]
			scale := b.yAxis.normalize(d)
			bar.Resize(fyne.NewSize(b.barChart.barWidth, availableHeight*scale))
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-bar.Size().Width/2,
				size.Height-xSize.Height-(availableHeight*scale)-b.xLblMax.Height-theme.Padding())
			bar.Move(rectPos)
		}
	}
}

func (b *barChartRenderer) MinSize() fyne.Size {
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
	for _, lbl := range b.yLabels {
		cos = append(cos, lbl)
	}

	for _, lbl := range b.xLabels {
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

	b.yAxis = axis{normalizer: linearNormalizer{}}

	for idx, datum := range b.barChart.data {
		b.yAxis.max = math.Max(b.yAxis.max, datum)
		b.yAxis.min = math.Min(b.yAxis.min, datum)
		if idx >= len(b.data) {
			b.data = append(b.data, newBar(b.barChart.canvas, datum))
		} else {
			b.data[idx].value = datum
		}
	}
	b.yAxis.dataRange = b.yAxis.max - b.yAxis.min

	for idx := range b.barChart.labels {
		var lbl *widget.Label
		if idx >= len(b.xLabels) {
			lbl = widget.NewLabel(b.barChart.labels[idx])
			b.xLblMax = b.xLblMax.Max(lbl.MinSize())
			b.xLabels = append(b.xLabels, lbl)
		} else {
			lbl = b.xLabels[idx]
			lbl.SetText(b.barChart.labels[idx])
			b.xLblMax = b.xLblMax.Max(lbl.MinSize())
		}
	}

	maps.Clear(b.yLabelPositions)
	for _, lbl := range b.yLabels {
		lbl.Hide()
	}
	tickLabels, _, _, _, err := generateTicks(b.yAxis.min, b.yAxis.max, 4, containmentWithinData, defaultQ(), defaultWeights(), defaultLegibility)
	if err != nil {
		log.Println("error generating ticks")
		return
	}
	log.Println("Tick Labels:", tickLabels)
	for idx, tl := range tickLabels {
		var lbl *widget.Label
		if idx >= len(b.yLabels) {
			lbl = widget.NewLabel(fmt.Sprintf("%.1f", tl))
			b.yLabels = append(b.yLabels, lbl)
		} else {
			lbl = b.yLabels[idx]
			lbl.SetText(fmt.Sprintf("%.1f", tl))
		}
		lbl.Show()
		b.yLblMax = b.xLblMax.Max(lbl.MinSize())
		b.yLabelPositions[lbl] = tl
	}

}
