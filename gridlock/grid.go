// Represents a low level hexagonal grid.

package main

import (
	"fmt"
	"log/slog"
)

// A hexagonal grid.
type grid struct {
	*hexConfig
}

// Denotes position for multi pixel game components such as nodes, obstacles, hexagons, etc
type GridPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type hexPoint int

const (
	hpCenter hexPoint = iota
	hpNorth
	hpNortheast
	hpSoutheast
	hpSouth
	hpSouthwest
	hpNorthwest
)

type hexSide int

const (
	hsNone hexSide = iota
	hsNortheast
	hsEast
	hsSoutheast
	hsSouthwest
	hsWest
	hsNorthwest
)

func hexPoints() []hexPoint {
	return []hexPoint{hpCenter, hpNorth, hpNortheast, hpSoutheast, hpSouth, hpSouthwest, hpNorthwest}
}

func hexSides() []hexSide {
	return []hexSide{hsNortheast, hsEast, hsSoutheast, hsSouthwest, hsWest, hsNorthwest}
}

func (gridPos GridPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", gridPos.X))
	a = append(a, slog.Int("y", gridPos.Y))
	return a
}

func (gridPos GridPosition) equals(other GridPosition) bool {
	return gridPos.X == other.X && gridPos.Y == other.Y
}

// Returns conservative partitions which collide with the passed position and distance.
// The caller is responsible for verifying the analysis performed with the partitions is valid.
// For example, analyzing partitions for collision detection between two circles would be invalid
// if the circle with the position and distance not passed has a radius which exceeds the
// conservative hex radius. Unoptimized because the hexagon collision uses conservative circles.
func (g grid) partitions(p PixelPosition, d int) []GridPosition {
	gridPos, ok := g.getGridPos(p)
	if !ok {
		panic(fmt.Sprintf("partitions: cannot find grid pos for %s", p.log()))
	}
	visited := make(map[GridPosition]struct{})
	visited[gridPos] = struct{}{}
	q := []GridPosition{gridPos}
	nextQ := []GridPosition{}
	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		for _, a := range g.allAdjacent(current) {
			if _, ok := visited[a]; !ok {
				visited[a] = struct{}{}
				// TODO: Fix this. This isn't correct but will probably work fine for awhile.
				// It doesn't handle border partition collisions 100% correct but this should be rare.
				// Also, I suspect this partitions method will most likey be refactored before issue
				// blocks development.
				pa, okA := g.getCenterWithFallback(a)
				if !okA {
					panic(fmt.Sprintf("partitions: cannot find center for pos %s", a.log()))
				}
				_, collides := collides(p, pa, d, g.conservativeHexRadius())
				if collides {
					nextQ = append(nextQ, a)
				}
			}
		}

		if len(q) == 0 && len(nextQ) == 0 {
			break
		} else if len(q) == 0 {
			q = nextQ
			nextQ = []GridPosition{}
		}
	}

	partitions := make([]GridPosition, len(visited))
	i := 0
	for p := range visited {
		partitions[i] = p
		i++
	}
	return partitions
}

func (g grid) allAdjacent(pos GridPosition) []GridPosition {
	a := []GridPosition{}
	for _, hs := range hexSides() {
		if p, ok := g.adjacent(pos, hs); ok {
			a = append(a, p)
		}
	}
	return a
}

func (g grid) adjacent(p GridPosition, side hexSide) (GridPosition, bool) {
	x := p.X
	y := p.Y
	switch side {
	case hsNone:
		panic("need a direction to find adjacent grid position")
	case hsNortheast:
		y--
		if p.Y%2 == 0 {
			x++
		}
	case hsEast:
		x++
	case hsSoutheast:
		y++
		if p.Y%2 == 0 {
			x++
		}
	case hsWest:
		x--
	case hsSouthwest:
		y++
		if p.Y%2 == 1 {
			x--
		}
	case hsNorthwest:
		y--
		if p.Y%2 == 1 {
			x--
		}
	default:
		panic("Cannot find adjacent grid position. Unknown hex side.")
	}

	if x < 0 || y < 0 || y >= g.yIndexMax || (y%2 == 0 && x >= g.xIndexMax-1) || (y%2 == 1 && x >= g.xIndexMax) {
		return GridPosition{}, false
	} else {
		return GridPosition{x, y}, true
	}
}

func (g grid) isAdjacent(p1 GridPosition, p2 GridPosition) bool {
	ne, neOk := g.adjacent(p1, hsNortheast)
	e, eOk := g.adjacent(p1, hsEast)
	se, seOk := g.adjacent(p1, hsSoutheast)
	sw, swOk := g.adjacent(p1, hsSouthwest)
	w, wOk := g.adjacent(p1, hsWest)
	nw, nwOk := g.adjacent(p1, hsNorthwest)

	return (neOk && ne.equals(p2)) ||
		(eOk && e.equals(p2)) ||
		(seOk && se.equals(p2)) ||
		(swOk && sw.equals(p2)) ||
		(wOk && w.equals(p2)) ||
		(nwOk && nw.equals(p2))
}

func (g grid) getPixelPos(pos GridPosition, hexPoint hexPoint) (PixelPosition, bool) {
	p := g.getPixelPosOutbounds(pos, hexPoint)
	if !newWindowConfig().inbounds(p) {
		return PixelPosition{}, false
	} else {
		return p, true
	}
}

func (g grid) getCenterWithFallback(pos GridPosition) (PixelPosition, bool) {
	if p, ok := g.getPixelPos(pos, hpCenter); ok {
		return p, true
	}

	for _, hp := range hexPoints() {
		if p, ok := g.getPixelPos(pos, hp); ok {
			return p, true
		}
	}
	return PixelPosition{}, false
}

func (g grid) getPixelPosOutbounds(p GridPosition, hexPoint hexPoint) PixelPosition {
	x := p.X * g.hexWidth
	y := p.Y * (g.hexHeight / 4 * 3)
	if p.Y%2 == 1 {
		x -= (g.hexWidth) / 2
	}

	// The values are found from the 1. hex look up table defined in hexConfig.go
	// and 2. arbitrary rule of hexagons sides are left inclusive and right exclusive.
	// These offsets are found from the hex look up table start point which is
	// located at (xOffset, yOffset).
	switch hexPoint {
	case hpCenter:
		y -= g.hexHeight / 4
	case hpNorth:
		y -= g.hexHeight / 4 * 3
	case hpNorthwest:
		x -= g.hexWidth / 2
		y -= g.hexHeight / 2
	case hpNortheast:
		x += g.hexWidth/2 - 1
		y -= g.hexHeight / 2
	case hpSouth:
		y += g.hexHeight/4 - 1
	case hpSouthwest:
		x -= g.hexWidth / 2
		y -= 1
	case hpSoutheast:
		x += g.hexWidth/2 - 1
		y -= 1
	}

	return PixelPosition{x + g.xOffset, y + g.yOffset}
}

func (g grid) getGridPos(p PixelPosition) (GridPosition, bool) {
	wc := newWindowConfig()
	if !wc.inbounds(p) {
		return GridPosition{}, false
	}
	return g.getGridPosOutbounds(p)
}

func (g grid) getGridPosOutbounds(p PixelPosition) (GridPosition, bool) {
	xIndex := (p.X+g.hexWidth-g.xOffset)/g.hexWidth - 1
	yIndex := (p.Y+(g.hexHeight/4*3)-g.yOffset)/(g.hexHeight/4*3) - 2
	innerX := (p.X + g.hexWidth - g.xOffset) % g.hexWidth
	innerY := (p.Y + (g.hexHeight * 3 / 2) - g.yOffset) % (g.hexHeight * 3 / 2)

	firstYSectionStart := 0
	secondYSectionStart := firstYSectionStart + g.hexHeight/4
	thirdYSectionStart := secondYSectionStart + g.hexHeight/2
	fourthYSectionStart := thirdYSectionStart + g.hexHeight/4

	var xIndexOffset int
	var yIndexOffset int
	switch {
	case innerY >= firstYSectionStart && innerY < secondYSectionStart:
		triangle := g.lut[innerX+innerY*g.hexWidth]
		switch triangle {
		case leftHexTriangle:
			xIndexOffset = 0
			yIndexOffset = 1
		case middleHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 2
		case rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= secondYSectionStart && innerY < thirdYSectionStart:
		xIndexOffset = 1
		yIndexOffset = 2
	case innerY >= thirdYSectionStart && innerY < fourthYSectionStart:
		hexLookUpYOffset := (g.hexWidth * g.hexHeight / 4) * 3
		var hexLookUpXOffset int
		if innerX >= 0 && innerX < (g.hexWidth/2) {
			hexLookUpXOffset = (g.hexWidth / 2)
		} else {
			hexLookUpXOffset = -(g.hexWidth / 2)
		}

		triangle := g.lut[innerX+(innerY*g.hexWidth)+hexLookUpXOffset-hexLookUpYOffset]
		switch {
		case triangle == leftHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		case triangle == middleHexTriangle && innerX < (g.hexWidth/2):
			xIndexOffset = 0
			yIndexOffset = 2
		case triangle == middleHexTriangle && innerX >= (g.hexWidth/2):
			xIndexOffset = 1
			yIndexOffset = 2
		case triangle == rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= fourthYSectionStart && innerY < (g.hexHeight*3/2):
		if innerX >= 0 && innerX < (g.hexWidth/2) {
			xIndexOffset = 0
			yIndexOffset = 2
		} else {
			xIndexOffset = 1
			yIndexOffset = 2
		}
	default:
		panic(fmt.Sprintf("getHexGridPostion could not find a y section for the inner y: %d for pixel position %s", innerY, p.log()))
	}

	return GridPosition{xIndex + xIndexOffset, yIndex + yIndexOffset}, true
}
