package fynecharts

import "fmt"

const (
	defaultBarWidth           = 25
	defaultMinHeight          = 100
	defaultSuggestedTickCount = 4
)

func defaultHoverFormat(input float64) string {
	return fmt.Sprintf("%.2f", input)
}

func defaultTickFormat(input float64) string {
	return fmt.Sprintf("%.1f", input)
}
