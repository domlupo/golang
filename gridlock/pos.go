// A unique 3D point for game physics.

package main

import (
	"fmt"
	"log/slog"
	"math"
)

// A unique 3D point for game physics. Units rotate which require high resolution for smooth
// movement. floats are used for high resolution but are more complex than ints.
type Pos struct {
	X float64
	Y float64
	L layer // z direction where each plane is named.
}

// Return the equivalent Pixel for a Pos. Precision is lost by rounding to the nearest integer.
func (p Pos) pixel() Pixel {
	return Pixel{X: int(math.Round(p.X)), Y: int(math.Round(p.Y))}
}

// Set the layer for the Pos. I suspect setting layers will be a sharp and error prone situation.
// This method provides some basic safety and documents where layers are set.
func (p Pos) layer(l layer) Pos {
	if p.L != lNone {
		panic(fmt.Sprintf("Pos.layer: cannot set non zero value layer %s to %s", p.L.log(), l.log()))
	}
	return Pos{X: p.X, Y: p.Y, L: l}
}

// Remove the layer for the Pos. I suspect removing layers will be a sharp and error prone
// situation. This method helps document where layers are removed.
func (p Pos) noLayer() Pos {
	return Pos{X: p.X, Y: p.Y}
}

// The relative 2D (x, y) orientation among 3 Pos.
type posOrientation int

const (
	noOrientation posOrientation = iota
	counterClockwise
	clockwise
	collinear
)

// Typical log method.
func (p Pos) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Float64("X", float64(p.X)))
	a = append(a, slog.Float64("Y", float64(p.Y)))
	a = append(a, p.L.log())
	return a
}

// Returns the angle in radians from p1 -> p2 clockwise from the north axis. When p1 and p2 have
// different layers and it is safe to find angle between them, callers should zero value both
// layers. Callers are responsible because I prefer the explicitness and safety trade off to avoid
// hard to debug issues rather than simpler caller code.
func (p1 Pos) angle(p2 Pos) float64 {
	if p1.L != p2.L {
		panic(fmt.Sprintf("Pos.angle: p1 layer %s != p2 layer %s", p1.L.log(), p2.L.log()))
	}
	return math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X)) + math.Pi/2
}

// Atypical equals method. Pos internally uses floats which causes approximate equals.
func (p1 Pos) equals(p2 Pos) bool {
	return p1.L == p2.L && equalFloats(p1.X, p2.X) && equalFloats(p1.Y, p2.Y)
}

// Returns the relative 2D (x, y) orientation among 3 Pos. The Pos directions are p -> q -> r.
func orientation(p1 Pos, p2 Pos, p3 Pos) posOrientation {
	if p1.L != p2.L || p2.L != p3.L || p1.L != p3.L {
		panic(fmt.Sprintf("Pos.orientation: mismatched layers p1 layer %s, p2 layer %s, p3 layer %s",
			p1.L.log(), p2.L.log(), p3.L.log()))
	}

	o := (p2.Y-p1.Y)*(p3.X-p2.X) - (p2.X-p1.X)*(p3.Y-p2.Y)
	switch {
	case o == 0:
		return collinear
	case o > 0:
		return clockwise
	case o < 0:
		return counterClockwise
	default:
		panic("Pos.orientation unknown orientation")

	}
}

// Returns the squared distance if passed values collide.
// TODO: This method should probably be removed once I add a collides method to dynamic that works
// for both circle and polygon?
func collides(p1, p2 Pos, r1, r2 int) (float64, bool) {
	if p1.L != p2.L {
		panic(fmt.Sprintf("Pos.collides: p1 layer %s != p2 layer %s", p1.L.log(), p2.L.log()))
	}

	xDelta := p2.X - p1.X
	yDelta := p2.Y - p1.Y
	rSum := float64(r1 + r2)
	squaredDistance := xDelta*xDelta + yDelta*yDelta
	if squaredDistance <= rSum*rSum {
		return squaredDistance, true
	} else {
		return 0, false
	}
}
