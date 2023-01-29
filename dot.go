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
	showValue    bool
	pos          fyne.Position
}

func (d *dot) MouseIn(event *desktop.MouseEvent) {
	d.showValue = true
	d.pos = event.Position
	d.Refresh()
	d.canvas.Refresh(d)
}

func (d *dot) MouseMoved(event *desktop.MouseEvent) {
	d.pos = event.Position
	d.Refresh()
	d.canvas.Refresh(d)
}

func (d *dot) MouseOut() {
	d.showValue = false
	d.Refresh()
	d.canvas.Refresh(d)
}

func (d *dot) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (d *dot) updateDisplayValue(v string) {
	d.displayValue = v
	d.Refresh()
}

func (d *dot) CreateRenderer() fyne.WidgetRenderer {
	return &dotRenderer{
		d:       d,
		circle:  canvas.NewCircle(theme.PrimaryColor()),
		wrapper: canvas.NewRectangle(theme.BackgroundColor()),
		display: widget.NewLabel(d.displayValue),
	}
}

func newDot(canvas fyne.Canvas, value string) *dot {
	d := &dot{canvas: canvas, displayValue: value}
	d.ExtendBaseWidget(d)
	d.Refresh()

	return d
}

type dotRenderer struct {
	d *dot

	circle  *canvas.Circle
	wrapper *canvas.Rectangle
	display *widget.Label
}

func (d *dotRenderer) Destroy() {

}

func (d *dotRenderer) Layout(size fyne.Size) {
	d.circle.Resize(size)
	d.wrapper.Resize(d.display.MinSize())
	d.display.Resize(d.display.MinSize())
}

func (d *dotRenderer) MinSize() fyne.Size {
	return d.circle.MinSize()
}

func (d *dotRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.circle, d.wrapper, d.display}
}

func (d *dotRenderer) Refresh() {
	if d.d.showValue {
		d.display.SetText(d.d.displayValue)
		d.display.Move(d.d.pos.SubtractXY(0, d.display.MinSize().Height))
		d.wrapper.Move(d.d.pos.SubtractXY(0, d.display.MinSize().Height))

		d.display.Show()
		d.wrapper.Show()
	} else {
		d.display.Hide()
		d.wrapper.Hide()
	}
}
