package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"math"
)

const (
	defaultDotDiameter = 18
)

type TimeSeriesChart struct {
	*BaseChart
	canvas fyne.Canvas

	data []float64

	dotDiameter float32

	hoverFormat func(float64) string
}

func (t *TimeSeriesChart) CreateRenderer() fyne.WidgetRenderer {
	bcr := t.BaseChart.CreateRenderer().(*baseChartRenderer)

	return &timeSeriesChartRenderer{
		baseChartRenderer: bcr,
		timeSeriesChart:   t,
	}
}

func NewTimeSeriesChart(canvas fyne.Canvas, title string, labels []string, data []float64) *TimeSeriesChart {
	tc := &TimeSeriesChart{BaseChart: newBaseChart(title, labels, defaultMinHeight, defaultSuggestedTickCount),
		canvas:      canvas,
		data:        data,
		dotDiameter: defaultDotDiameter,
		hoverFormat: defaultHoverFormat,
	}
	tc.ExtendBaseWidget(tc)
	tc.Refresh()

	return tc
}

func (t *TimeSeriesChart) UpdateHoverFormat(f func(float642 float64) string) {
	t.hoverFormat = f
	t.Refresh()
}

func (t *TimeSeriesChart) UpdateData(lbls []string, data []float64) {
	t.xLabels = lbls
	t.data = data
	t.Refresh()
}

func (t *TimeSeriesChart) UpdateDotDiameter(diameter float32) {
	t.dotDiameter = diameter
	t.Refresh()
}

type timeSeriesChartRenderer struct {
	*baseChartRenderer
	timeSeriesChart *TimeSeriesChart

	data         []*dot
	connectLines []*canvas.Line
}

func (t *timeSeriesChartRenderer) Destroy() {

}

func (t *timeSeriesChartRenderer) Layout(size fyne.Size) {
	t.baseChartRenderer.Layout(size)

	xOffset := t.xOffset()

	availableHeight := t.availableHeight(size)
	columnWidth := t.columnWidth(size, xOffset)

	reqBottom := t.requiredBottomHeight()
	if len(t.data) > 0 {
		var previousPos *fyne.Position
		for idx, d := range t.timeSeriesChart.data {
			dt := t.data[idx]
			scale := t.yAxis.normalize(d)
			dt.Resize(fyne.NewSize(t.timeSeriesChart.dotDiameter, t.timeSeriesChart.dotDiameter))
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-dt.Size().Width/2,
				size.Height-reqBottom-(availableHeight*scale)-dt.Size().Height/2)
			if previousPos != nil {
				l := t.connectLines[idx-1]
				l.Position1 = (*previousPos).AddXY(dt.Size().Width/2, dt.Size().Height/2)
				l.Position2 = rectPos.AddXY(dt.Size().Width/2, dt.Size().Height/2)
			}
			previousPos = &rectPos
			dt.Move(rectPos)
		}
	}
}

func (t *timeSeriesChartRenderer) MinSize() fyne.Size {
	titleSize := fyne.NewSize(0, 0)
	paddingCount := 0
	if t.titleLbl.Visible() {
		titleSize = t.titleLbl.MinSize()
		paddingCount++
		paddingCount++
	}
	xLblSize := fyne.NewSize(0, 0)
	if t.xLbl.Visible() {
		xLblSize = t.xLbl.MinSize()
		paddingCount++
	}

	xCellWidth := t.xLblMax.Width + 2
	return fyne.NewSize(float32(len(t.timeSeriesChart.xLabels))*xCellWidth,
		titleSize.Height+xLblSize.Height+t.xLblMax.Height+float32(paddingCount)*theme.Padding()+t.timeSeriesChart.minHeight)
}

func (t *timeSeriesChartRenderer) Objects() []fyne.CanvasObject {
	cos := t.baseChartRenderer.Objects()
	for _, l := range t.connectLines {
		cos = append(cos, l)
	}
	for _, d := range t.data {
		cos = append(cos, d)
	}
	return cos
}

func (t *timeSeriesChartRenderer) Refresh() {
	t.yAxis = axis{normalizer: linearNormalizer{}}
	/*** Commenting this out here for now, as reuse was keeping the layout from updating on data change. ***/
	//for _, ln := range t.connectLines {
	//	ln.Hide()
	//}
	//for _, dt := range t.data {
	//	dt.Hide()
	//}
	t.connectLines = nil
	t.data = nil
	for _, datum := range t.timeSeriesChart.data {
		t.yAxis.max = math.Max(t.yAxis.max, datum)
		t.yAxis.min = math.Min(t.yAxis.min, datum)
		t.data = append(t.data, newDot(t.timeSeriesChart.canvas, t.timeSeriesChart.hoverFormat(datum)))
		/*** Commenting this out here for now, as reuse was keeping the layout from updating on data change. ***/
		//if idx >= len(t.data) {
		//	t.data = append(t.data, newDot(t.timeSeriesChart.canvas, t.timeSeriesChart.hoverFormat(datum)))
		//} else {
		//	t.data[idx].updateDisplayValue(t.timeSeriesChart.hoverFormat(datum))
		//	t.data[idx].Show()
		//}

		l := canvas.NewLine(theme.PrimaryColor())
		l.StrokeWidth = 2
		t.connectLines = append(t.connectLines, l)
		/*** Commenting this out here for now, as reuse was keeping the layout from updating on data change. ***/
		//if idx >= len(t.connectLines) {
		//	l := canvas.NewLine(theme.PrimaryColor())
		//	l.StrokeWidth = 2
		//	t.connectLines = append(t.connectLines, l)
		//}
	}
	t.yAxis.dataRange = t.yAxis.max - t.yAxis.min

	t.baseChartRenderer.Refresh()
}
