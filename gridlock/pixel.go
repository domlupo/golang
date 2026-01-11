// A unique 2D point for game visuals.

package main

import (
	"log/slog"
)

// A unique 2D point for game visuals. The number of Pixels represents the native game resolution.
// Pixels are used to show game visuals to players. A less obvious use is representing screen input
// from player e.g. mouse click.
type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// typical log method
func (p Pixel) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("X", p.X))
	a = append(a, slog.Int("Y", p.Y))
	return a
}

// Returns the equivalent Pos for a Pixel. If a comparison between a Pixel and Pos requires a
// conversion, it is typically best to convert to Pos over Pixel because the conversion has easier
// to predict behavior for players. Converting from Pixel to Pos requires less predictable rounding.
func (p Pixel) pos() Pos {
	return Pos{X: float64(p.X), Y: float64(p.Y)}
}

// Converts a 2D point to a 3D point. The zero layer aka "lNone" may be used if methods provided by
// the Pos struct need called but the point has no layer. e.g. Pos.angle().
func (p Pixel) posLayer(l layer) Pos {
	return Pos{X: float64(p.X), Y: float64(p.Y), L: l}
}

// typical equals method
func (p1 Pixel) equals(p2 Pixel) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

// Chebyshev distance treats diagonal, horizontal and vertical movement as equal. It is also called
// chessboard distance because the king piece moves 1 square in any direction.
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
