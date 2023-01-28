package fynecharts

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
	"log"
)

type BaseChart struct {
	widget.BaseWidget

	title   string
	yTitle  string
	xTitle  string
	xLabels []string

	suggestedTickCount int

	minHeight float32
}

func (b *BaseChart) CreateRenderer() fyne.WidgetRenderer {
	titleLbl := canvas.NewText(b.title, theme.ForegroundColor())
	titleLbl.TextSize = theme.TextSize() + 6
	titleLbl.Hide()
	yLbl := widget.NewLabel(b.yTitle)
	yLbl.Hide()
	ySep := canvas.NewLine(theme.ForegroundColor())
	ySep.StrokeWidth = 2
	xLbl := widget.NewLabel(b.xTitle)
	xLbl.Hide()
	xSep := canvas.NewLine(theme.ForegroundColor())
	xSep.StrokeWidth = 2

	return &baseChartRenderer{baseChart: b,
		titleLbl:        titleLbl,
		yLbl:            yLbl,
		ySeparator:      ySep,
		xLbl:            xLbl,
		xSeparator:      xSep,
		yLabelPositions: make(map[*widget.Label]float64),
	}
}

func newBaseChart(title string, xLabels []string, minHeight float32, suggestedTickCount int) *BaseChart {
	bc := &BaseChart{title: title, xLabels: xLabels, minHeight: minHeight, suggestedTickCount: suggestedTickCount}

	return bc
}

func (b *BaseChart) SetXLabel(xlbl string) {
	b.xTitle = xlbl
	b.Refresh()
}

type baseChartRenderer struct {
	baseChart *BaseChart

	titleLbl   *canvas.Text
	yLbl       *widget.Label
	ySeparator *canvas.Line
	xLbl       *widget.Label
	xSeparator *canvas.Line

	yLabels         []*widget.Label
	yLabelPositions map[*widget.Label]float64
	xLabels         []*widget.Label

	yLblMax fyne.Size
	xLblMax fyne.Size

	yAxis axis
}

func (b *baseChartRenderer) Destroy() {

}

func (b *baseChartRenderer) Layout(size fyne.Size) {
	titleSize := b.titleLbl.MinSize()
	titleX := size.Width/2 - titleSize.Width/2
	titlePos := fyne.NewPos(titleX, theme.Padding())
	b.titleLbl.Move(titlePos)
	b.titleLbl.Resize(titleSize)

	xSize := b.xLabelSize()

	xLblX := size.Width/2 - xSize.Width/2
	xPos := fyne.NewPos(xLblX, size.Height-xSize.Height-theme.Padding())
	b.xLbl.Move(xPos)

	xOffset := b.xOffset()

	xSepY := size.Height - xSize.Height - b.xLblMax.Height - 5
	b.xSeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.xSeparator.Position2 = fyne.NewPos(size.Width-theme.Padding(), xSepY)

	b.ySeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.ySeparator.Position2 = fyne.NewPos(xOffset, titleSize.Height+2*theme.Padding())

	availableHeight := b.availableHeight(size)
	columnWidth := b.columnWidth(size, xOffset)

	if len(b.yLabelPositions) > 0 {
		for lbl, y := range b.yLabelPositions {
			lblSize := lbl.MinSize()
			scale := b.yAxis.normalize(y)
			pos := fyne.NewPos(0, size.Height-xSize.Height-b.xLblMax.Height-theme.Padding()-(availableHeight*scale)-lblSize.Height/2)
			lbl.Move(pos)
		}
	}

	if len(b.baseChart.xLabels) > 0 {
		for idx := range b.baseChart.xLabels {
			lbl := b.xLabels[idx]
			lblSize := lbl.MinSize()
			xCellOffset := float32(idx) * columnWidth
			lblPos := fyne.NewPos(xOffset+xCellOffset+columnWidth/2-lblSize.Width/2,
				size.Height-2*theme.Padding()-xSize.Height-lblSize.Height)
			lbl.Move(lblPos)
		}
	}
}

func (b *baseChartRenderer) xLabelSize() fyne.Size {
	xSize := fyne.NewSize(0, 0)
	if b.xLbl.Visible() {
		xSize = b.xLbl.MinSize()
	}
	return xSize
}

func (b *baseChartRenderer) xOffset() float32 {
	return b.yLblMax.Width
}

func (b *baseChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (b *baseChartRenderer) columnWidth(size fyne.Size, xOffset float32) float32 {
	return (size.Width - xOffset - theme.Padding()) / float32(len(b.baseChart.xLabels))
}

func (b *baseChartRenderer) availableHeight(size fyne.Size) float32 {
	return size.Height - b.titleLbl.MinSize().Height - theme.Padding() - b.xLbl.MinSize().Height - b.xLblMax.Height - theme.Padding()
}

func (b *baseChartRenderer) Objects() []fyne.CanvasObject {
	cos := []fyne.CanvasObject{b.titleLbl, b.yLbl, b.xLbl}
	for _, lbl := range b.yLabels {
		cos = append(cos, lbl)
	}

	for _, lbl := range b.xLabels {
		cos = append(cos, lbl)
	}

	if b.xSeparator != nil {
		cos = append(cos, b.xSeparator)
	}

	if b.ySeparator != nil {
		cos = append(cos, b.ySeparator)
	}

	return cos
}

func (b *baseChartRenderer) Refresh() {
	if b.baseChart.title != "" {
		b.titleLbl.Text = b.baseChart.title
		b.titleLbl.Refresh()
		b.titleLbl.Show()
	} else {
		b.titleLbl.Hide()
	}

	if b.baseChart.yTitle != "" {
		b.yLbl.SetText(b.baseChart.yTitle)
		b.yLbl.Show()
	} else {
		b.yLbl.Hide()
	}

	if b.baseChart.xTitle != "" {
		b.xLbl.SetText(b.baseChart.xTitle)
		b.xLbl.Show()
	} else {
		b.xLbl.Hide()
	}

	for _, lbl := range b.xLabels {
		lbl.Hide()
	}
	for idx := range b.baseChart.xLabels {
		var lbl *widget.Label
		if idx >= len(b.xLabels) {
			lbl = widget.NewLabel(b.baseChart.xLabels[idx])
			b.xLabels = append(b.xLabels, lbl)
		} else {
			lbl = b.xLabels[idx]
			lbl.SetText(b.baseChart.xLabels[idx])
		}
		lbl.Show()
		b.xLblMax = b.xLblMax.Max(lbl.MinSize())
	}

	maps.Clear(b.yLabelPositions)
	for _, lbl := range b.yLabels {
		lbl.Hide()
	}
	tickLabels, _, _, _, err := generateTicks(b.yAxis.min, b.yAxis.max, b.baseChart.suggestedTickCount, containmentContainData, defaultQ(), defaultWeights(), defaultLegibility)
	if err != nil {
		log.Println("error generating ticks")
		return
	}

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
