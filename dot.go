package fynecharts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type dot struct {
	widget.BaseWidget
	canvas fyne.Canvas

	value float64
	pop   *widget.PopUp
}

func (d *dot) CreateRenderer() fyne.WidgetRenderer {
	return &dotRenderer{
		d:      d,
		circle: canvas.NewCircle(theme.PrimaryColor()),
	}
}

func newDot(canvas fyne.Canvas, value float64) *dot {
	d := &dot{canvas: canvas, value: value}
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
