package fynecharts

import (
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

	tickFormat func(input float64) string
}

func (b *BaseChart) CreateRenderer() fyne.WidgetRenderer {
	titleLbl := canvas.NewText(b.title, theme.ForegroundColor())
	titleLbl.TextSize = theme.TextSize() + 6
	titleLbl.Hide()
	yLbl := canvas.NewText(b.yTitle, theme.ForegroundColor())
	yLbl.TextSize = theme.TextSize() + 3
	yLbl.Hide()
	ySep := canvas.NewLine(theme.ForegroundColor())
	ySep.StrokeWidth = 2
	xLbl := canvas.NewText(b.xTitle, theme.ForegroundColor())
	xLbl.TextSize = theme.TextSize() + 3
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
	bc := &BaseChart{title: title, xLabels: xLabels, minHeight: minHeight, suggestedTickCount: suggestedTickCount, tickFormat: defaultTickFormat}

	return bc
}

func (b *BaseChart) SetXLabel(xLbl string) {
	b.xTitle = xLbl
	b.Refresh()
}

func (b *BaseChart) SetYLabel(yLbl string) {
	b.yTitle = yLbl
	b.Refresh()
}

func (b *BaseChart) SetMinHeight(h float32) {
	b.minHeight = h
}

func (b *BaseChart) UpdateSuggestedTickCount(count int) {
	b.suggestedTickCount = count
}

func (b *BaseChart) UpdateTickFormat(f func(input float64) string) {
	b.tickFormat = f
}

type baseChartRenderer struct {
	baseChart *BaseChart

	titleLbl   *canvas.Text
	yLbl       *canvas.Text
	ySeparator *canvas.Line
	xLbl       *canvas.Text
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

	yPos := fyne.NewPos(theme.Padding(), titleSize.Height+2*theme.Padding())
	b.yLbl.Move(yPos)

	xSize := b.xLabelSize()
	xLblX := size.Width/2 - xSize.Width/2
	xPos := fyne.NewPos(xLblX, size.Height-xSize.Height-theme.Padding())
	b.xLbl.Move(xPos)

	xOffset := b.xOffset()

	reqBottom := b.requiredBottomHeight()
	xSepY := size.Height - reqBottom
	b.xSeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.xSeparator.Position2 = fyne.NewPos(size.Width-theme.Padding(), xSepY)

	b.ySeparator.Position1 = fyne.NewPos(xOffset, xSepY)
	b.ySeparator.Position2 = fyne.NewPos(xOffset, b.requiredTopHeight())

	availableHeight := b.availableHeight(size)
	columnWidth := b.columnWidth(size, xOffset)

	if len(b.yLabelPositions) > 0 {
		for lbl, y := range b.yLabelPositions {
			lblSize := lbl.MinSize()
			scale := b.yAxis.normalize(y)
			pos := fyne.NewPos(0, size.Height-reqBottom-(availableHeight*scale)-lblSize.Height/2)
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

func (b *baseChartRenderer) yLabelSize() fyne.Size {
	ySize := fyne.NewSize(0, 0)
	if b.yLbl.Visible() {
		ySize = b.yLbl.MinSize()
	}
	return ySize
}

func (b *baseChartRenderer) titleLabelSize() fyne.Size {
	tSize := fyne.NewSize(0, 0)
	if b.titleLbl.Visible() {
		tSize = b.titleLbl.MinSize()
	}
	return tSize
}

func (b *baseChartRenderer) xOffset() float32 {
	//ySize := b.yLabelSize()
	//return fyne.Max(b.yLblMax.Width, ySize.Width+2*theme.Padding())
	return b.yLblMax.Width + theme.Padding()
}

func (b *baseChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (b *baseChartRenderer) columnWidth(size fyne.Size, xOffset float32) float32 {
	return (size.Width - xOffset - theme.Padding()) / float32(len(b.baseChart.xLabels))
}

func (b *baseChartRenderer) requiredTopHeight() float32 {
	paddingCount := 0
	titleSize := b.titleLabelSize()
	if titleSize.Height > 0 {
		paddingCount += 3
	}
	ySize := b.yLabelSize()
	if ySize.Height > 0 {
		paddingCount += 2
	}

	return titleSize.Height + ySize.Height + float32(paddingCount)*theme.Padding()
}

func (b *baseChartRenderer) requiredBottomHeight() float32 {
	paddingCount := 0
	xSize := b.xLabelSize()
	if xSize.Height > 0 {
		paddingCount += 2
	}

	if b.xLblMax.Height > 0 {
		paddingCount++
	}

	return xSize.Height + b.xLblMax.Height + float32(paddingCount)*theme.Padding()
}

func (b *baseChartRenderer) availableHeight(size fyne.Size) float32 {
	topHeight := b.requiredTopHeight()
	bottomHeight := b.requiredBottomHeight()
	return size.Height - topHeight - bottomHeight
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
		b.yLbl.Text = b.baseChart.yTitle
		b.yLbl.Refresh()
		b.yLbl.Show()
	} else {
		b.yLbl.Hide()
	}

	if b.baseChart.xTitle != "" {
		b.xLbl.Text = b.baseChart.xTitle
		b.xLbl.Refresh()
		b.xLbl.Show()
	} else {
		b.xLbl.Hide()
	}

	//for _, lbl := range b.xLabels {
	//	lbl.Hide()
	//}
	b.xLabels = nil
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
	//for _, lbl := range b.yLabels {
	//	lbl.Hide()
	//}
	b.yLabels = nil
	tickLabels, _, _, _, err := generateTicks(b.yAxis.min, b.yAxis.max, b.baseChart.suggestedTickCount, containmentContainData, defaultQ(), defaultWeights(), defaultLegibility)
	if err != nil {
		log.Println("error generating ticks")
		return
	}

	for idx, tl := range tickLabels {
		var lbl *widget.Label
		if idx >= len(b.yLabels) {
			lbl = widget.NewLabel(b.baseChart.tickFormat(tl))
			b.yLabels = append(b.yLabels, lbl)
		} else {
			lbl = b.yLabels[idx]
			lbl.SetText(b.baseChart.tickFormat(tl))
		}
		lbl.Show()
		b.yLblMax = b.yLblMax.Max(lbl.MinSize())
		b.yLabelPositions[lbl] = tl
	}

}
