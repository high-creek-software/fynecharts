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
	pop          *widget.PopUp
}

func (b *bar) MouseIn(event *desktop.MouseEvent) {
	if b.pop != nil {
		return
	}
	if b.canvas != nil {
		lbl := widget.NewLabel(b.displayValue)
		b.pop = widget.NewPopUp(lbl, b.canvas)
		b.pop.ShowAtPosition(event.AbsolutePosition.SubtractXY(0, lbl.MinSize().Height+10))
	}
}

func (b *bar) MouseMoved(event *desktop.MouseEvent) {

}

func (b *bar) MouseOut() {
	if b.pop == nil {
		return
	}
	b.pop.Hide()
	b.pop = nil
}

func (b *bar) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (b *bar) CreateRenderer() fyne.WidgetRenderer {
	return &barRenderer{
		b:    b,
		rect: canvas.NewRectangle(theme.PrimaryColor()),
	}
}

func newBar(canvas fyne.Canvas, value string) *bar {
	b := &bar{canvas: canvas, displayValue: value}
	b.ExtendBaseWidget(b)

	return b
}

type barRenderer struct {
	b *bar

	rect *canvas.Rectangle
}

func (b *barRenderer) Destroy() {

}

func (b *barRenderer) Layout(size fyne.Size) {
	b.rect.Resize(size)
}

func (b *barRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (b *barRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{b.rect}
}

func (b *barRenderer) Refresh() {

}
