// A line is two unordered connected Positions.

package main

import (
	"log/slog"
)

// Two unordered connected Positions.
type line struct {
	p1 Pos
	p2 Pos
}

// Typical log method.
func (l line) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("p1", l.p1.log()))
	a = append(a, nest("p2", l.p2.log()))
	return a
}

// Typical equals method.
func (l1 line) equals(l2 line) bool {
	return (l1.p1.equals(l2.p1) && l1.p2.equals(l2.p2)) ||
		(l1.p2.equals(l2.p1) && l1.p1.equals(l2.p2))
}

func (l line) isVertical() bool {
	return l.p1.X == l.p2.X
}

func (l line) contains(p Pos) bool {
	return p.X <= max(l.p1.X, l.p2.X) &&
		p.X >= min(l.p1.X, l.p2.X) && p.Y <= max(l.p1.Y, l.p2.Y) && p.Y >= min(l.p1.Y, l.p2.Y)
}

func (l line) perpendicular() line {
	return line{l.p1, Pos{X: -l.p2.Y, Y: l.p2.X}}
}

func (l1 line) parallel(l2 line) bool {
	switch {
	case l1.isVertical() && l2.isVertical():
		return true
	case l1.isVertical() && !l2.isVertical():
		return false
	case !l1.isVertical() && l2.isVertical():
		return false
	default:
		slope1 := (l1.p2.Y - l1.p1.Y) / (l1.p2.X - l1.p1.X)
		slope2 := (l2.p2.Y - l2.p1.Y) / (l2.p2.X - l2.p1.X)
		return slope1 == slope2
	}
}

func (l1 line) intersects(l2 line) bool {
	o1 := orientation(l1.p1, l1.p2, l2.p1)
	o2 := orientation(l1.p1, l1.p2, l2.p2)
	o3 := orientation(l2.p1, l2.p2, l1.p1)
	o4 := orientation(l2.p1, l2.p2, l1.p2)

	switch {
	case o1 != o2 && o3 != o4:
		return true
	case o1 == collinear && l1.contains(l2.p1):
		return true
	case o2 == collinear && l1.contains(l2.p2):
		return true
	case o3 == collinear && l2.contains(l1.p1):
		return true
	case o4 == collinear && l2.contains(l1.p2):
		return true
	default:
		return false
	}
}
