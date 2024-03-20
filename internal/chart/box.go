package chart

import "math"

type Box struct {
	//Top    int
	//Left   int
	//Right  int
	//Bottom int
	Top    float64
	Left   float64
	Right  float64
	Bottom float64
}

func (b *Box) Width() float64 {
	return abs(b.Right - b.Left)
}

func (b *Box) Height() float64 {
	return abs(b.Bottom - b.Top)
}

func (b *Box) Center() (x, y float64) {
	w2, h2 := int(b.Width())>>1, int(b.Height())>>1
	return b.Left + float64(w2), b.Top + float64(h2)
}

func (b *Box) Clone() *Box {
	return &Box{
		Top:    b.Top,
		Left:   b.Left,
		Right:  b.Right,
		Bottom: b.Bottom,
	}
}

func (b *Box) Grow(other *Box) *Box {
	return &Box{
		Top:    math.Min(float64(b.Top), float64(other.Top)),
		Left:   math.Min(float64(b.Left), float64(other.Left)),
		Right:  math.Max(float64(b.Right), float64(other.Right)),
		Bottom: math.Max(float64(b.Bottom), float64(other.Bottom)),
	}
}

func (b *Box) Corners() *BoxCorners {
	return &BoxCorners{
		TopLeft:     Point{b.Left, b.Top},
		TopRight:    Point{b.Right, b.Top},
		BottomRight: Point{b.Right, b.Bottom},
		BottomLeft:  Point{b.Left, b.Bottom},
	}
}

func (b *Box) OuterConstrain(bounds, other *Box) *Box {
	newBox := b.Clone()
	if other.Top < bounds.Top {
		delta := bounds.Top - other.Top
		newBox.Top = b.Top + delta
	}

	if other.Left < bounds.Left {
		delta := bounds.Left - other.Left
		newBox.Left = b.Left + delta
	}

	if other.Right > bounds.Right {
		delta := other.Right - bounds.Right
		newBox.Right = b.Right - delta
	}

	if other.Bottom > bounds.Bottom {
		delta := other.Bottom - bounds.Bottom
		newBox.Bottom = b.Bottom - delta
	}
	return newBox
}

type BoxCorners struct {
	TopLeft, TopRight, BottomRight, BottomLeft Point
}

func (bc *BoxCorners) Box() *Box {
	return &Box{
		Top:    math.Min(bc.TopLeft.Y, bc.TopRight.Y),
		Left:   math.Min(bc.TopLeft.X, bc.BottomLeft.X),
		Right:  math.Max(bc.TopRight.X, bc.BottomRight.X),
		Bottom: math.Max(bc.BottomLeft.Y, bc.BottomRight.Y),
	}
}

func (bc *BoxCorners) Center() (x, y int) {
	left := mean(bc.TopLeft.X, bc.BottomLeft.X)
	right := mean(bc.TopRight.X, bc.BottomRight.X)
	x = ((int(right) - int(left)) >> 1) + int(left)

	top := mean(bc.TopLeft.Y, bc.TopRight.Y)
	bottom := mean(bc.BottomLeft.Y, bc.BottomRight.Y)
	y = ((int(bottom) - int(top)) >> 1) + int(top)

	return
}

func (bc *BoxCorners) Rotate(thetaDegrees float64) *BoxCorners {
	cx, cy := bc.Center()

	thetaRadians := degreesToRadians(thetaDegrees)

	tlx, tly := rotateCoordinate(cx, cy, bc.TopLeft.X, bc.TopLeft.Y, thetaRadians)
	trx, try := rotateCoordinate(cx, cy, bc.TopRight.X, bc.TopRight.Y, thetaRadians)
	brx, bry := rotateCoordinate(cx, cy, bc.BottomRight.X, bc.BottomRight.Y, thetaRadians)
	blx, bly := rotateCoordinate(cx, cy, bc.BottomLeft.X, bc.BottomLeft.Y, thetaRadians)

	return &BoxCorners{
		TopLeft:     Point{tlx, tly},
		TopRight:    Point{trx, try},
		BottomRight: Point{brx, bry},
		BottomLeft:  Point{blx, bly},
	}
}

type Point struct {
	X, Y float64
}
