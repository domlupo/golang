// A position represents a two dimensional point (PixelPosition) or a three dimension point (position).

package main

import (
	"fmt"
	"log/slog"
	"math"
)

// Two dimensional point.
type PixelPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Three dimensional point.
type position struct {
	layer
	PixelPosition
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

type line struct {
	start PixelPosition
	end   PixelPosition
}

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

func (d direction) String() string {
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

func (p PixelPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", p.X))
	a = append(a, slog.Int("y", p.Y))
	return a
}

// Returns the angle in radians from p1 -> p2 clockwise from the north axis.
func (p1 PixelPosition) angle(p2 PixelPosition) float64 {
	return math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X)) + math.Pi/2
}

func (p PixelPosition) chebyshevDistance(other PixelPosition) int {
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

func (p PixelPosition) equals(other PixelPosition) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p position) equals(other position) bool {
	return p.layer == other.layer && p.PixelPosition.equals(other.PixelPosition)
}

func (p PixelPosition) adjacent(d direction) (PixelPosition, bool) {
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

	a := PixelPosition{X: p.X + x, Y: p.Y + y}
	if !newWindowConfig().inbounds(a) {
		return PixelPosition{}, false
	}
	return a, true
}

func (p PixelPosition) allAdjacent() []PixelPosition {
	a := []PixelPosition{}
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
	a = append(a, p.PixelPosition.log()...)
	return a
}

func (p1 PixelPosition) add(p2 PixelPosition) (PixelPosition, bool) {
	p := PixelPosition{X: p1.X + p2.X, Y: p1.Y + p2.Y}
	if !newWindowConfig().inbounds(p) {
		return PixelPosition{}, false
	} else {
		return p, true
	}
}

// The point directions are p -> q -> r.
func orientation(p1 PixelPosition, p2 PixelPosition, p3 PixelPosition) pointOrientation {
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

func onSegment(p1 PixelPosition, p2 PixelPosition, p3 PixelPosition) bool {
	return p2.X <= max(p1.X, p3.X) && p2.X >= min(p1.X, p3.X) && p2.Y <= max(p1.Y, p3.Y) && p2.Y >= min(p1.Y, p3.Y)
}

func (l1 line) intersects(l2 line) bool {
	o1 := orientation(l1.start, l1.end, l2.start)
	o2 := orientation(l1.start, l1.end, l2.end)
	o3 := orientation(l2.start, l2.end, l1.start)
	o4 := orientation(l2.start, l2.end, l1.end)

	switch {
	case o1 != o2 && o3 != o4:
		return true
	case o1 == collinear && onSegment(l1.start, l2.start, l1.end):
		return true
	case o2 == collinear && onSegment(l1.start, l2.end, l1.end):
		return true
	case o3 == collinear && onSegment(l2.start, l1.start, l2.end):
		return true
	case o4 == collinear && onSegment(l2.start, l1.end, l2.end):
		return true
	default:
		return false
	}
}

func (l line) inclusivePoints() []PixelPosition {
	var inner []PixelPosition
	var start int
	var end int
	switch {
	case l.start.X == l.end.X && l.start.Y == l.end.Y:
		panic("inclusivePoints: cannot find inclusive points for a point")
	case l.start.X == l.end.X:
		start = min(l.start.Y, l.end.Y)
		end = max(l.start.Y, l.end.Y)
		y := start
		for y <= end {
			inner = append(inner, PixelPosition{l.start.X, y})
			y += 1
		}
	case l.start.Y == l.end.Y:
		start = min(l.start.Y, l.end.Y)
		end = max(l.start.Y, l.end.Y)
		x := start
		for x <= end {
			inner = append(inner, PixelPosition{x, l.start.Y})
			x += 1
		}
	default:
		panic("inclusivePoints: cannot find inclusive points for line that is not either horizontal or vertical")
	}
	return inner
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
func collides(p1, p2 PixelPosition, r1, r2 int) (int, bool) {
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
