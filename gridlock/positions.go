// A position represents a two dimensional (pixel) or a three dimension point (position).

package main

import (
	"fmt"
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

type direction int

const (
	noDirecton direction = iota
	north
	northeast
	east
	southeast
	south
	southwest
	west
	northwest
)

type pointOrientation int

const (
	noOrientation pointOrientation = iota
	counterClockwise
	clockwise
	collinear
)

func directions() []direction {
	return []direction{north, northeast, east, southeast, south, southwest, west, northwest}
}

func (d direction) string() string {
	switch d {
	case noDirecton:
		return "noDirection"
	case north:
		return "north"
	case northeast:
		return "northeast"
	case east:
		return "east"
	case southeast:
		return "southeast"
	case south:
		return "south"
	case southwest:
		return "southwest"
	case west:
		return "west"
	case northwest:
		return "northwest"
	default:
		panic(fmt.Sprintf("Cannot convert direction to string: %d", d))
	}
}

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

func (p Pixel) adjacent(d direction) (Pixel, bool) {
	var x int
	var y int

	switch d {
	case north:
		y = -1
	case northeast:
		x = 1
		y = -1
	case east:
		x = 1
	case southeast:
		x = 1
		y = 1
	case south:
		y = 1
	case southwest:
		x = -1
		y = 1
	case west:
		x = -1
		x = -1
	case northwest:
		x = -1
		y = -1
	default:
		panic(fmt.Sprintf("cannot get adjacent pixel position for %s unknown direction: %d", p.log(), d))
	}

	a := Pixel{X: p.X + x, Y: p.Y + y}
	if !newWindowConfig().contains(a) {
		return Pixel{}, false
	}
	return a, true
}

func (p Pixel) allAdjacent() []Pixel {
	a := []Pixel{}
	for _, d := range directions() {
		if pos, ok := p.adjacent(d); ok {
			a = append(a, pos)
		}
	}
	return a
}

func (p position) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, p.layer.log())
	a = append(a, p.Pixel.log()...)
	return a
}

func (p1 Pixel) add(p2 Pixel) (Pixel, bool) {
	p := Pixel{X: p1.X + p2.X, Y: p1.Y + p2.Y}
	if !newWindowConfig().contains(p) {
		return Pixel{}, false
	} else {
		return p, true
	}
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

func collidesWithBorder(p position, offset int) bool {
	wc := newWindowConfig()
	if p.X+offset > wc.width ||
		p.X-offset < 0 ||
		p.Y+offset > wc.height ||
		p.Y-offset < 0 {
		return true
	} else {
		return false
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
