package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"math"
)

type BarChart struct {
	*BaseChart
	canvas fyne.Canvas

	data []float64

	barWidth float32

	hoverFormat func(float64) string
}

func (b *BarChart) CreateRenderer() fyne.WidgetRenderer {
	bcr := b.BaseChart.CreateRenderer().(*baseChartRenderer)

	return &barChartRenderer{barChart: b, baseChartRenderer: bcr}
}

func (b *BarChart) UpdateData(labels []string, data []float64) {
	b.xLabels = labels
	b.data = data
	b.Refresh()
}

func (b *BarChart) SetBarWidth(w float32) {
	b.barWidth = w
}

func (b *BarChart) SetMinHeight(h float32) {
	b.minHeight = h
}

func (b *BarChart) UpdateSuggestedTickCount(count int) {
	b.suggestedTickCount = count
}

func (b *BarChart) UpdateHoverFormat(f func(float642 float64) string) {
	b.hoverFormat = f
}

func NewBarChart(canvas fyne.Canvas, title string, labels []string, data []float64) *BarChart {
	bc := &BarChart{BaseChart: newBaseChart(title, labels, defaultMinHeight, defaultSuggestedTickCount),
		canvas:      canvas,
		data:        data,
		barWidth:    defaultBarWidth,
		hoverFormat: defaultHoverFormat,
	}
	bc.ExtendBaseWidget(bc)
	bc.Refresh()

	return bc
}

type barChartRenderer struct {
	*baseChartRenderer
	barChart *BarChart

	data []*bar
}

func (b *barChartRenderer) Destroy() {

}

func (b *barChartRenderer) Layout(size fyne.Size) {
	b.baseChartRenderer.Layout(size)

	xSize := b.xLabelSize()
	xOffset := b.xOffset()

	availableHeight := b.availableHeight(size)
	columnWidth := b.columnWidth(size, xOffset)

	if len(b.data) > 0 {
		for idx, d := range b.barChart.data {
			bar := b.data[idx]
			scale := b.yAxis.normalize(d)
			bar.Resize(fyne.NewSize(b.barChart.barWidth, availableHeight*scale))
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-bar.Size().Width/2,
				size.Height-xSize.Height-b.xLblMax.Height-theme.Padding()-(availableHeight*scale))
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
	return fyne.NewSize(float32(len(b.barChart.xLabels))*xCellWidth,
		titleSize.Height+xLblSize.Height+b.xLblMax.Height+float32(paddingCount)*theme.Padding()+b.barChart.minHeight)
}

func (b *barChartRenderer) Objects() []fyne.CanvasObject {
	cos := b.baseChartRenderer.Objects()
	for _, d := range b.data {
		cos = append(cos, d)
	}
	return cos
}

func (b *barChartRenderer) Refresh() {
	b.yAxis = axis{normalizer: linearNormalizer{}}
	for idx, datum := range b.barChart.data {
		b.yAxis.max = math.Max(b.yAxis.max, datum)
		b.yAxis.min = math.Min(b.yAxis.min, datum)
		if idx >= len(b.data) {
			b.data = append(b.data, newBar(b.barChart.canvas, b.barChart.hoverFormat(datum)))
		} else {
			b.data[idx].displayValue = b.barChart.hoverFormat(datum)
		}
	}
	b.yAxis.dataRange = b.yAxis.max - b.yAxis.min

	b.baseChartRenderer.Refresh()
}
