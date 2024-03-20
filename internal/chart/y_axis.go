package chart

import (
	"io"
	"math"

	"github.com/caarlos0/starcharts/internal/chart/svg"
)

type YAxis struct {
	Name        string
	StrokeWidth float64
	Color       string
}

func (ya *YAxis) Measure(canvas *Box, ra *Range, ticks []Tick) *Box {
	tx := canvas.Right + YAxisMargin

	minX, maxX, minY, maxY := math.MaxInt32, 0, math.MaxInt32, 0
	maxTextHeight := 0
	for _, t := range ticks {
		ly := canvas.Bottom - ra.Translate(t.Value)

		tb := measureText(t.Label, AxisFontSize)
		maxTextHeight = int(math.Max(tb.Height(), float64(maxTextHeight)))

		minX = int(canvas.Right)
		maxX = int(math.Max(float64(maxX), tx+tb.Width()))

		tbh2 := int(tb.Height()) >> 1
		minY = int(math.Min(float64(minY), ly-float64(tbh2)))
		maxY = int(math.Max(float64(maxY), ly+float64(tbh2)))
	}

	maxX += YAxisMargin + maxTextHeight

	return &Box{
		Top:    float64(minY),
		Left:   float64(minX),
		Right:  float64(maxX),
		Bottom: float64(maxY),
	}
}

func (ya *YAxis) Render(w io.Writer, canvasBox *Box, ra *Range, ticks []Tick) {
	lx := canvasBox.Right
	tx := lx + YAxisMargin
	strokeStyle := styles("stroke", ya.Color)
	fillStyle := styles("fill", ya.Color)

	strokeWidth := normaliseStrokeWidth(ya.StrokeWidth)

	svg.Path().
		Attr("stroke-width", strokeWidth).
		Attr("style", strokeStyle).
		MoveTo(int(lx), int(canvasBox.Bottom)).
		LineToF(float64(lx), float64(canvasBox.Top)-ya.StrokeWidth/2).
		Render(w)

	var maxTextWidth int
	var finalTextY int
	for _, t := range ticks {
		ly := canvasBox.Bottom - ra.Translate(t.Value)
		tb := measureText(t.Label, AxisFontSize)

		if tb.Width() > float64(maxTextWidth) {
			maxTextWidth = int(tb.Width())
		}

		finalTextY = int(ly) + int(tb.Height())>>1

		svg.Path().
			Attr("stroke-width", strokeWidth).
			Attr("style", strokeStyle).
			MoveTo(int(lx), int(ly)).
			LineTo(int(lx+HorizontalTickWidth), int(ly)).
			Render(w)

		svg.Text().
			Content(t.Label).
			Attr("style", fillStyle).
			Attr("x", svg.Point(tx)).
			Attr("y", svg.Point(finalTextY)).
			Render(w)
	}

	tb := measureText(ya.Name, AxisFontSize)
	tx = canvasBox.Right + float64(YAxisMargin) + float64(maxTextWidth) + float64(YAxisMargin)
	ty := int(canvasBox.Top) + (int(canvasBox.Height())>>1 - int(tb.Height())>>1)

	svg.Text().
		Content(ya.Name).
		Attr("x", svg.Point(tx)).
		Attr("y", svg.Point(ty)).
		Attr("style", fillStyle).
		Attr("transform", rotate(90, int(tx), ty)).
		Render(w)
}
