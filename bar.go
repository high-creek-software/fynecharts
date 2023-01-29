package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type bar struct {
	widget.BaseWidget
	canvas fyne.Canvas

	displayValue string
	showValue    bool
	pos          fyne.Position
	idx          int
	onTouched    func(idx int)
}

func (b *bar) Tapped(event *fyne.PointEvent) {
	if b.onTouched != nil {
		b.onTouched(b.idx)
	}
}

func (b *bar) updateOnTouched(f func(idx int), idx int) {
	b.onTouched = f
	b.idx = idx
}

func (b *bar) MouseIn(event *desktop.MouseEvent) {
	b.showValue = true
	b.pos = event.Position
	b.Refresh()
	b.canvas.Refresh(b)
}

func (b *bar) MouseMoved(event *desktop.MouseEvent) {
	b.pos = event.Position
	b.Refresh()
	b.canvas.Refresh(b)
}

func (b *bar) MouseOut() {
	b.showValue = false
	b.Refresh()
	b.canvas.Refresh(b)
}

func (b *bar) updateDisplayValue(v string) {
	b.displayValue = v
	b.Refresh()
}

func (b *bar) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (b *bar) CreateRenderer() fyne.WidgetRenderer {
	return &barRenderer{
		b:       b,
		rect:    canvas.NewRectangle(theme.PrimaryColor()),
		wrapper: canvas.NewRectangle(theme.BackgroundColor()),
		display: widget.NewLabel(b.displayValue),
	}
}

func newBar(canvas fyne.Canvas, value string) *bar {
	b := &bar{canvas: canvas, displayValue: value}
	b.ExtendBaseWidget(b)
	b.Refresh()

	return b
}

type barRenderer struct {
	b *bar

	rect    *canvas.Rectangle
	wrapper *canvas.Rectangle
	display *widget.Label
}

func (b *barRenderer) Destroy() {

}

func (b *barRenderer) Layout(size fyne.Size) {
	b.rect.Resize(size)
	b.wrapper.Resize(b.display.MinSize())
	b.display.Resize(b.display.MinSize())
}

func (b *barRenderer) MinSize() fyne.Size {
	return b.rect.MinSize()
}

func (b *barRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{b.rect, b.wrapper, b.display}
}

func (b *barRenderer) Refresh() {
	if b.b.showValue {
		b.display.SetText(b.b.displayValue)
		b.display.Move(b.b.pos.SubtractXY(0, b.display.MinSize().Height))
		b.wrapper.Move(b.b.pos.SubtractXY(0, b.display.MinSize().Height))

		b.display.Show()
		b.wrapper.Show()
	} else {
		b.display.Hide()
		b.wrapper.Hide()
	}
}
