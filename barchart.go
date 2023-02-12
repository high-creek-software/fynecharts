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
	onTouched   func(idx int)
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
	b.Refresh()
}

func (b *BarChart) UpdateHoverFormat(f func(float642 float64) string) {
	b.hoverFormat = f
}

func (b *BarChart) UpdateOnTouched(f func(idx int)) {
	b.onTouched = f
	b.Refresh()
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

	xOffset := b.xOffset()

	availableHeight := b.availableHeight(size)
	columnWidth := b.columnWidth(size, xOffset)

	reqBottom := b.requiredBottomHeight()
	if len(b.data) > 0 {
		for idx, d := range b.barChart.data {
			br := b.data[idx]
			scale := b.yAxis.normalize(d)
			brSize := fyne.NewSize(b.barChart.barWidth, availableHeight*scale)
			br.Resize(brSize)
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-br.Size().Width/2,
				size.Height-reqBottom-(availableHeight*scale))
			br.Move(rectPos)
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

	xLblWidth := b.xLblMax.Width + 2
	xCellWidth := b.barChart.barWidth + 2
	return fyne.NewSize(float32(len(b.barChart.xLabels))*fyne.Max(xLblWidth, xCellWidth),
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
	/*** Commenting this out here for now, as reuse was keeping the layout from updating on data change. ***/
	//for _, br := range b.data {
	//	br.Hide()
	//}
	b.data = nil
	for idx, datum := range b.barChart.data {
		b.yAxis.max = math.Max(b.yAxis.max, datum)
		b.yAxis.min = math.Min(b.yAxis.min, datum)

		br := newBar(b.barChart.canvas, b.barChart.hoverFormat(datum))
		br.updateOnTouched(b.barChart.onTouched, idx)
		b.data = append(b.data, br)

		/*** Commenting this out here for now, as reuse was keeping the layout from updating on data change. ***/
		//if idx >= len(b.data) {
		//	br := newBar(b.barChart.canvas, b.barChart.hoverFormat(datum))
		//	br.updateOnTouched(b.barChart.onTouched, idx)
		//	b.data = append(b.data, br)
		//} else {
		//	b.data[idx].updateDisplayValue(b.barChart.hoverFormat(datum))
		//	b.data[idx].updateOnTouched(b.barChart.onTouched, idx)
		//	b.data[idx].Show()
		//}
	}
	b.yAxis.dataRange = b.yAxis.max - b.yAxis.min

	b.baseChartRenderer.Refresh()
}
