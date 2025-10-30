// A position represents a point in 2d (PixelPosition) or 3d space (position)

package main

import (
	"fmt"
	"log/slog"
)

// high resolution position.
type PixelPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type layer int

const (
	noLayer layer = iota
	land
	lowAir
	highAir
)

// globally unique position
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

func (pos position) isZero() bool {
	return pos.layer == noLayer
}

func (pixelPos PixelPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", pixelPos.X))
	a = append(a, slog.Int("y", pixelPos.Y))
	return a
}

func (l layer) log() slog.Attr {
	var layer string
	switch l {
	case noLayer:
		layer = "noLayer"
	case land:
		layer = "land"
	case lowAir:
		layer = "lowAir"
	case highAir:
		layer = "HighAir"
	}
	return slog.String("layer", layer)
}

func (pos PixelPosition) chebyshevDistance(otherPos PixelPosition) int {
	xDelta := pos.X - otherPos.X
	if xDelta < 0 {
		xDelta *= -1
	}
	yDelta := pos.Y - otherPos.Y
	if yDelta < 0 {
		yDelta *= -1
	}
	return max(xDelta, yDelta)
}

func (pixelPos PixelPosition) equals(other PixelPosition) bool {
	return pixelPos.X == other.X && pixelPos.Y == other.Y
}

func (pos PixelPosition) adjacent(d direction) (PixelPosition, bool) {
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
		panic(fmt.Sprintf("cannot get adjacent pixel position for %s unknown direction: %d", pos.log(), d))
	}

	p := PixelPosition{X: pos.X + x, Y: pos.Y + y}
	if !newWindowConfig().inbounds(p) {
		return PixelPosition{}, false
	}
	return p, true
}

func (pixelPos PixelPosition) allAdjacent() []PixelPosition {
	a := []PixelPosition{}
	for _, d := range directions() {
		if p, ok := pixelPos.adjacent(d); ok {
			a = append(a, p)
		}
	}
	return a
}

func (pos position) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, pos.layer.log())
	a = append(a, pos.PixelPosition.log()...)
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
func orientation(p PixelPosition, q PixelPosition, r PixelPosition) pointOrientation {
	v := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	switch {
	case v == 0:
		return collinear
	case v > 0:
		return clockwise
	case v < 0:
		return counterClockwise
	default:
		panic("cannot find pixel position orientation")

	}
}

func onSegment(p PixelPosition, q PixelPosition, r PixelPosition) bool {
	return q.X <= max(p.X, r.X) && q.X >= min(p.X, r.X) && q.Y <= max(p.Y, r.Y) && q.Y >= min(p.Y, r.Y)
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

func collidesWithBorder(pos position, offset int) bool {
	wc := newWindowConfig()
	if pos.X+offset > wc.width ||
		pos.X-offset < 0 ||
		pos.Y+offset > wc.height ||
		pos.Y-offset < 0 {
		return true
	} else {
		return false
	}
}
