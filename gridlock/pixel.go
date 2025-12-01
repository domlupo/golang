// A Pixel represents a two dimensional point. A position represents a three dimension point.

package main

import (
	"log/slog"
	"math"
)

// Two dimensional point.
type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Three dimensional point.
type position struct {
	layer
	Pixel
}

type pointOrientation int

const (
	noOrientation pointOrientation = iota
	counterClockwise
	clockwise
	collinear
)

func (p Pixel) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", p.X))
	a = append(a, slog.Int("y", p.Y))
	return a
}

// Returns the angle in radians from p1 -> p2 clockwise from the north axis.
func (p1 Pixel) angle(p2 Pixel) float64 {
	return math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X)) + math.Pi/2
}

func (p Pixel) chebyshevDistance(other Pixel) int {
	xDelta := p.X - other.X
	if xDelta < 0 {
		xDelta *= -1
	}
	yDelta := p.Y - other.Y
	if yDelta < 0 {
		yDelta *= -1
	}
	return max(xDelta, yDelta)
}

func (p Pixel) equals(other Pixel) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p position) equals(other position) bool {
	return p.layer == other.layer && p.Pixel.equals(other.Pixel)
}

func (p position) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, p.layer.log())
	a = append(a, p.Pixel.log()...)
	return a
}

// The point directions are p -> q -> r.
func orientation(p1 Pixel, p2 Pixel, p3 Pixel) pointOrientation {
	o := (p2.Y-p1.Y)*(p3.X-p2.X) - (p2.X-p1.X)*(p3.Y-p2.Y)
	switch {
	case o == 0:
		return collinear
	case o > 0:
		return clockwise
	case o < 0:
		return counterClockwise
	default:
		panic("cannot find pixel position orientation")

	}
}

// Returns the squared distance if passed values collide.
func collides(p1, p2 Pixel, r1, r2 int) (int, bool) {
	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y
	xDelta := x2 - x1
	yDelta := y2 - y1
	rSum := r1 + r2
	squaredDistance := xDelta*xDelta + yDelta*yDelta
	if squaredDistance <= rSum*rSum {
		return squaredDistance, true
	} else {
		return 0, false
	}
}
