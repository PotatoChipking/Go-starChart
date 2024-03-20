package chart

import (
	"io"
	"math"

	"github.com/caarlos0/starcharts/internal/chart/svg"
)

type XAxis struct {
	Name        string
	StrokeWidth float64
	Color       string
}

func (xa *XAxis) Measure(canvas *Box, ra *Range, ticks []Tick) *Box {
	var ltx, rtx int
	var tx, ty float64
	left, right, bottom := math.MaxInt32, 0, 0
	for _, t := range ticks {
		v := t.Value
		tb := measureText(t.Label, AxisFontSize)

		tx = canvas.Left + ra.Translate(v)
		ty = canvas.Bottom + XAxisMargin + tb.Height()
		ltx = int(tx) - int(tb.Width())>>1
		rtx = int(tx) + int(tb.Width())>>1

		left = int(math.Min(float64(left), float64(ltx)))
		right = int(math.Max(float64(right), float64(rtx)))
		bottom = int(math.Max(float64(bottom), ty))
	}

	tb := measureText(xa.Name, AxisFontSize)
	bottom += XAxisMargin + int(tb.Height())

	return &Box{
		Top:    canvas.Bottom,
		Left:   float64(left),
		Right:  float64(right),
		Bottom: float64(bottom),
	}
}

func (xa *XAxis) Render(w io.Writer, canvasBox *Box, ra *Range, ticks []Tick) {
	strokeWidth := normaliseStrokeWidth(xa.StrokeWidth)
	strokeStyle := styles("stroke", xa.Color)
	fillStyle := styles("fill", xa.Color)

	svg.Path().
		Attr("stroke-width", strokeWidth).
		Attr("style", strokeStyle).
		MoveToF(float64(canvasBox.Left)-xa.StrokeWidth/2, float64(canvasBox.Bottom)).
		LineTo(int(canvasBox.Right), int(canvasBox.Bottom)).
		Render(w)

	var tx, ty int
	var maxTextHeight int
	for _, t := range ticks {
		v := t.Value
		lx := ra.Translate(v)

		tx = int(canvasBox.Left + lx)

		svg.Path().
			Attr("stroke-width", strokeWidth).
			Attr("style", strokeStyle).
			MoveTo(tx, int(canvasBox.Bottom)).
			LineTo(tx, int(canvasBox.Bottom+VerticalTickHeight)).
			Render(w)

		tb := measureText(t.Label, AxisFontSize)

		tx = tx - int(tb.Width())>>1
		ty = int(canvasBox.Bottom + XAxisMargin + tb.Height())

		svg.Text().
			Content(t.Label).
			Attr("style", fillStyle).
			Attr("x", svg.Point(tx)).
			Attr("y", svg.Point(ty)).
			Render(w)

		maxTextHeight = int(math.Max(float64(maxTextHeight), tb.Height()))
	}

	tb := measureText(xa.Name, AxisFontSize)
	tx = int(canvasBox.Right - float64((int(canvasBox.Width())>>1 + int(tb.Width())>>1)))
	ty = int(canvasBox.Bottom + float64(XAxisMargin) + float64(maxTextHeight) + float64(XAxisMargin) + tb.Height())

	svg.Text().
		Content(xa.Name).
		Attr("style", fillStyle).
		Attr("x", svg.Point(tx)).
		Attr("y", svg.Point(ty)).
		Render(w)
}
