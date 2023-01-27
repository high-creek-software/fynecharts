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
	defaultDotDiameter = 18
)

type TimeSeriesChart struct {
	widget.BaseWidget
	canvas fyne.Canvas

	title  string
	yTitle string
	xTitle string
	labels []string
	data   []float64

	dotDiameter        float32
	minHeight          float32
	suggestedTickCount int
}

func (t *TimeSeriesChart) CreateRenderer() fyne.WidgetRenderer {
	titleLbl := widget.NewLabel(t.title)
	titleLbl.Hide()
	yLbl := widget.NewLabel(t.yTitle)
	yLbl.Hide()
	ySep := canvas.NewLine(theme.ForegroundColor())
	ySep.StrokeWidth = 2
	xLbl := widget.NewLabel(t.xTitle)
	xLbl.Hide()
	xSep := canvas.NewLine(theme.ForegroundColor())
	xSep.StrokeWidth = 2

	return &timeSeriesChartRenderer{timeSeriesChart: t,
		titleLbl:        titleLbl,
		yLbl:            yLbl,
		ySeparator:      ySep,
		xLbl:            xLbl,
		xSeparator:      xSep,
		yLabelPositions: make(map[*widget.Label]float64),
	}
}

func NewTimeSeriesChart(canvas fyne.Canvas, title string, labels []string, data []float64) *TimeSeriesChart {
	tc := &TimeSeriesChart{canvas: canvas, title: title, labels: labels, data: data, dotDiameter: defaultDotDiameter, minHeight: defaultMinHeight, suggestedTickCount: defaultSuggestedTickCount}
	tc.ExtendBaseWidget(tc)
	tc.Refresh()

	return tc
}

type timeSeriesChartRenderer struct {
	timeSeriesChart *TimeSeriesChart
	titleLbl        *widget.Label
	yLbl            *widget.Label
	ySeparator      *canvas.Line
	xLbl            *widget.Label
	xSeparator      *canvas.Line

	yLabels         []*widget.Label
	yLabelPositions map[*widget.Label]float64
	xLabels         []*widget.Label
	data            []*dot

	yLblMax fyne.Size
	xLblMax fyne.Size

	yAxis axis
}

func (t *timeSeriesChartRenderer) Destroy() {

}

func (t *timeSeriesChartRenderer) Layout(size fyne.Size) {
	titleSize := t.titleLbl.MinSize()
	//ySize := b.yLbl.MinSize()
	xSize := fyne.NewSize(0, 0)
	if t.xLbl.Visible() {
		xSize = t.xLbl.MinSize()
	}

	titleX := size.Width/2 - titleSize.Width/2
	titlePos := fyne.NewPos(titleX, theme.Padding())
	t.titleLbl.Move(titlePos)
	t.titleLbl.Resize(titleSize)

	xLblX := size.Width/2 - xSize.Width/2
	xPos := fyne.NewPos(xLblX, size.Height-xSize.Height-theme.Padding())
	t.xLbl.Move(xPos)

	xOffset := theme.Padding() + t.yLblMax.Width

	xSepY := size.Height - xSize.Height - t.xLblMax.Height - 5
	t.xSeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	t.xSeparator.Position2 = fyne.NewPos(size.Width-theme.Padding(), xSepY)

	t.ySeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	t.ySeparator.Position2 = fyne.NewPos(xOffset, titleSize.Height+2*theme.Padding())

	availableHeight := size.Height - titleSize.Height - theme.Padding() - xSize.Height - t.xLblMax.Height - theme.Padding()
	columnWidth := (size.Width - xOffset - theme.Padding()) / float32(len(t.timeSeriesChart.labels))

	if len(t.yLabelPositions) > 0 {
		for lbl, y := range t.yLabelPositions {
			lblSize := lbl.MinSize()
			scale := t.yAxis.normalize(y)
			pos := fyne.NewPos(0, size.Height-xSize.Height-t.xLblMax.Height-theme.Padding()-(availableHeight*scale)-lblSize.Height/2)
			if pos.Y < titleSize.Height+2*theme.Padding() {
				lbl.Hide()
				continue
			}
			lbl.Move(pos)
		}
	}

	if len(t.timeSeriesChart.labels) > 0 {
		for idx := range t.timeSeriesChart.labels {
			lbl := t.xLabels[idx]
			lblSize := lbl.MinSize()
			xCellOffset := float32(idx) * columnWidth
			lblPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-lblSize.Width/2,
				size.Height-2*theme.Padding()-xSize.Height-lblSize.Height)
			lbl.Move(lblPos)
		}
	}

	if len(t.data) > 0 {
		for idx, d := range t.timeSeriesChart.data {
			dt := t.data[idx]
			scale := t.yAxis.normalize(d)
			dt.Resize(fyne.NewSize(t.timeSeriesChart.dotDiameter, availableHeight*scale))
			xCellOffset := float32(idx) * columnWidth
			rectPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-dt.Size().Width/2,
				size.Height-xSize.Height-t.xLblMax.Height-theme.Padding()-(availableHeight*scale))
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
	return fyne.NewSize(float32(len(t.timeSeriesChart.labels))*xCellWidth,
		titleSize.Height+xLblSize.Height+t.xLblMax.Height+float32(paddingCount)*theme.Padding()+t.timeSeriesChart.minHeight)
}

func (t *timeSeriesChartRenderer) Objects() []fyne.CanvasObject {
	cos := []fyne.CanvasObject{t.titleLbl, t.yLbl, t.xLbl}
	for _, lbl := range t.yLabels {
		cos = append(cos, lbl)
	}

	for _, lbl := range t.xLabels {
		cos = append(cos, lbl)
	}

	for _, d := range t.data {
		cos = append(cos, d)
	}

	if t.xSeparator != nil {
		cos = append(cos, t.xSeparator)
	}

	if t.ySeparator != nil {
		cos = append(cos, t.ySeparator)
	}

	return cos
}

func (t *timeSeriesChartRenderer) Refresh() {
	if t.timeSeriesChart.title != "" {
		t.titleLbl.SetText(t.timeSeriesChart.title)
		t.titleLbl.Show()
	} else {
		t.titleLbl.Hide()
	}

	if t.timeSeriesChart.yTitle != "" {
		t.yLbl.SetText(t.timeSeriesChart.yTitle)
		t.yLbl.Show()
	} else {
		t.yLbl.Hide()
	}

	if t.timeSeriesChart.xTitle != "" {
		t.xLbl.SetText(t.timeSeriesChart.xTitle)
		t.xLbl.Show()
	} else {
		t.xLbl.Hide()
	}

	t.yAxis = axis{normalizer: linearNormalizer{}}

	for idx, datum := range t.timeSeriesChart.data {
		t.yAxis.max = math.Max(t.yAxis.max, datum)
		t.yAxis.min = math.Min(t.yAxis.min, datum)
		if idx >= len(t.data) {
			t.data = append(t.data, newDot(t.timeSeriesChart.canvas, datum))
		} else {
			t.data[idx].value = datum
		}
	}
	t.yAxis.dataRange = t.yAxis.max - t.yAxis.min

	for _, lbl := range t.xLabels {
		lbl.Hide()
	}
	for idx := range t.timeSeriesChart.labels {
		var lbl *widget.Label
		if idx >= len(t.xLabels) {
			lbl = widget.NewLabel(t.timeSeriesChart.labels[idx])
			t.xLabels = append(t.xLabels, lbl)
		} else {
			lbl = t.xLabels[idx]
			lbl.SetText(t.timeSeriesChart.labels[idx])
		}
		lbl.Show()
		t.xLblMax = t.xLblMax.Max(lbl.MinSize())
	}

	maps.Clear(t.yLabelPositions)
	for _, lbl := range t.yLabels {
		lbl.Hide()
	}
	tickLabels, _, _, _, err := generateTicks(t.yAxis.min, t.yAxis.max, t.timeSeriesChart.suggestedTickCount, containmentWithinData, defaultQ(), defaultWeights(), defaultLegibility)
	if err != nil {
		log.Println("error generating ticks")
		return
	}
	log.Println("Tick Labels:", tickLabels)
	for idx, tl := range tickLabels {
		var lbl *widget.Label
		if idx >= len(t.yLabels) {
			lbl = widget.NewLabel(fmt.Sprintf("%.2f", tl))
			t.yLabels = append(t.yLabels, lbl)
		} else {
			lbl = t.yLabels[idx]
			lbl.SetText(fmt.Sprintf("%.2f", tl))
		}
		lbl.Show()
		t.yLblMax = t.xLblMax.Max(lbl.MinSize())
		t.yLabelPositions[lbl] = tl
	}
}
