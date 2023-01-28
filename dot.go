package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type dot struct {
	widget.BaseWidget
	canvas fyne.Canvas

	displayValue string
	pop          *widget.PopUp
}

func (d *dot) MouseIn(event *desktop.MouseEvent) {
	if d.canvas != nil {
		lbl := widget.NewLabel(d.displayValue)
		d.pop = widget.NewPopUp(lbl, d.canvas)
		d.pop.ShowAtPosition(event.AbsolutePosition.SubtractXY(0, lbl.MinSize().Height+10))
	}
}

func (d *dot) MouseMoved(event *desktop.MouseEvent) {

}

func (d *dot) MouseOut() {
	if d.pop == nil {
		return
	}
	d.pop.Hide()
	d.pop = nil
}

func (d *dot) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (d *dot) CreateRenderer() fyne.WidgetRenderer {
	return &dotRenderer{
		d:      d,
		circle: canvas.NewCircle(theme.PrimaryColor()),
	}
}

func newDot(canvas fyne.Canvas, value string) *dot {
	d := &dot{canvas: canvas, displayValue: value}
	d.ExtendBaseWidget(d)

	return d
}

type dotRenderer struct {
	d *dot

	circle *canvas.Circle
}

func (d *dotRenderer) Destroy() {

}

func (d *dotRenderer) Layout(size fyne.Size) {
	d.circle.Resize(size)
}

func (d *dotRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (d *dotRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.circle}
}

func (d *dotRenderer) Refresh() {

}
