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
	p1 PixelPosition
	p2 PixelPosition
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

func (p PixelPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", p.X))
	a = append(a, slog.Int("y", p.Y))
	return a
}

func (l line) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("p1", l.p1.log()))
	a = append(a, nest("p2", l.p2.log()))
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
	if !newWindowConfig().contains(a) {
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
	if !newWindowConfig().contains(p) {
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
	o1 := orientation(l1.p1, l1.p2, l2.p1)
	o2 := orientation(l1.p1, l1.p2, l2.p2)
	o3 := orientation(l2.p1, l2.p2, l1.p1)
	o4 := orientation(l2.p1, l2.p2, l1.p2)

	switch {
	case o1 != o2 && o3 != o4:
		return true
	case o1 == collinear && onSegment(l1.p1, l2.p1, l1.p2):
		return true
	case o2 == collinear && onSegment(l1.p1, l2.p2, l1.p2):
		return true
	case o3 == collinear && onSegment(l2.p1, l1.p1, l2.p2):
		return true
	case o4 == collinear && onSegment(l2.p1, l1.p2, l2.p2):
		return true
	default:
		return false
	}
}

func (l line) inclusivePoints() []PixelPosition {
	var inner []PixelPosition
	switch {
	case l.p1.X == l.p2.X && l.p1.Y == l.p2.Y:
		panic("inclusivePoints: cannot find inclusive points for a point")
	case l.p1.X == l.p2.X:
		start := min(l.p1.Y, l.p2.Y)
		end := max(l.p1.Y, l.p2.Y)
		y := start
		for y <= end {
			inner = append(inner, PixelPosition{l.p1.X, y})
			y += 1
		}
	case l.p1.Y == l.p2.Y:
		start := min(l.p1.X, l.p2.X)
		end := max(l.p1.X, l.p2.X)
		x := start
		for x <= end {
			inner = append(inner, PixelPosition{x, l.p1.Y})
			x += 1
		}
	default:
		panic("inclusivePoints: cannot find inclusive points for line that is not either horizontal or vertical")
	}
	return inner
}

// Returns the adjacent points in the passed direction. The function has multiple assumptions.
// Assumes a square grid with typical compass directions. Assumes a square grid where east is
// positive x and south is positive y. x is typical and y is opposite of typical for catesian graph.
// Assumes the line should be drawn in the top left corner of the grid square. Sloped lines
// determine point direction of split grid squares via which point has more area.
// Trims any points outside the passed bounds. Accepts passed lines outside the passed bounds.
func (l line) adjacentPoints(d direction, bounds bounds) []PixelPosition {
	switch d {
	case north, south, west, east:
	default:
		panic(fmt.Sprintf("adjacentPoints: cannot get points in direction %s since it is unimplemented", d.string()))
	}

	var w PixelPosition
	var e PixelPosition
	switch {
	case l.p1.X < l.p2.X:
		w = l.p1
		e = l.p2
	default:
		w = l.p2
		e = l.p1
	}
	nY := max(l.p1.Y, l.p2.Y)
	sY := min(l.p1.Y, l.p2.Y)

	lineIsVertical := w.X == e.X
	var points []PixelPosition
	switch {
	case lineIsVertical && d == west:
		y := sY
		for y <= nY {
			points = append(points, PixelPosition{w.X - 1, y})
			y += 1
		}
	case lineIsVertical && d == east:
		y := sY
		for y <= nY {
			points = append(points, PixelPosition{w.X, y})
			y += 1
		}
	case !lineIsVertical && d == west, !lineIsVertical && d == east:
		panic(fmt.Sprintf("adjacentPoints: cannot get points in direction %s for non vertical lines since it is unimplemented", d.string()))
	case lineIsVertical && d == north:
		points = append(points, PixelPosition{w.X, sY - 1})
	case lineIsVertical && d == south:
		points = append(points, PixelPosition{w.X, nY})
	default:
		// Ebitengine y is opposite of typical cartesian grid, slope equation is not typical.
		// e.Y - w.Y instead of typical w.Y - e.Y
		slope := float64(e.Y-w.Y) / float64(w.X-e.X)
		x := float64(w.X) + 0.5
		for x <= float64(e.X+1) {
			yIntercept := float64(w.Y) + (slope * float64(w.X))
			// Ebitegine y is opposite of typical cartesian grid, y intercept equation is not typical.
			// y := yIntercept - slope*x instead of y := slope*x + yIntercept.
			y := yIntercept - slope*x
			xInteger, _ := math.Modf(x)
			yInteger, yDecimal := math.Modf(y)
			xInt := int(xInteger)
			yInt := int(yInteger)
			switch {
			case yDecimal == 0.5:
				panic("adjacentPoints: undefined behavior when y midpoint occurs")
			case yDecimal < 0.5:
				yInt -= 1
			}
			var yOffset int
			if d == south {
				yOffset++
			}
			points = append(points, PixelPosition{xInt, yInt + yOffset})
			x += 1
		}
	}

	var trimmedPoints []PixelPosition
	for _, p := range points {
		if bounds.contains(p) {
			trimmedPoints = append(trimmedPoints, p)
		}
	}
	return trimmedPoints
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
