// General position utility functions. Most are pure functions. Used across many different files.

package main

import (
	"fmt"
	"log/slog"
)

// low resolution position used for grid items such as nodes, obstacles, hexagons, etc
type gridPosition struct {
	x int
	y int
}

// high resolution position.
type PixelPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type layer int

const (
	noLayer layer = iota
	ground
	lowAltitude
	highAltitude
)

// globally unique position
type position struct {
	layer
	PixelPosition
}

type hexPoint int

const (
	center hexPoint = iota
	northPoint
	northeastPoint
	southeastPoint
	southPoint
	southwestPoint
	northwestPoint
)

type hexSide int

const (
	noSide hexSide = iota
	northeastSide
	eastSide
	southeastSide
	southwestSide
	westSide
	northwestSide
)

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

func (gridPos gridPosition) equals(other gridPosition) bool {
	return gridPos.x == other.x && gridPos.y == other.y
}

func (pixelPos PixelPosition) equals(other PixelPosition) bool {
	return pixelPos.X == other.X && pixelPos.Y == other.Y
}

func (pixelPos PixelPosition) adjacent(d direction, hc *hexConfig) (PixelPosition, bool) {
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
		panic("cannot get adjacent pixel position unknown direction")
	}

	p := PixelPosition{X: pixelPos.X + x, Y: pixelPos.Y + y}
	if p.X < 0 || p.X > hc.windowWidth-1 || p.Y < 0 || p.Y >= hc.windowHeight-1 {
		return PixelPosition{}, false
	}
	return p, true
}

func (pixelPos PixelPosition) allAdjacent(hc *hexConfig) []PixelPosition {
	a := []PixelPosition{}
	for _, d := range directions() {
		if p, ok := pixelPos.adjacent(d, hc); ok {
			a = append(a, p)
		}
	}
	return a
}

func (gridPos gridPosition) adjacent(side hexSide, hc *hexConfig) (gridPosition, bool) {
	x := gridPos.x
	y := gridPos.y
	switch side {
	case noSide:
		panic("need a direction to find adjacent grid position")
	case northeastSide:
		y--
		if gridPos.y%2 == 0 {
			x++
		}
	case eastSide:
		x++
	case southeastSide:
		y++
		if gridPos.y%2 == 0 {
			x++
		}
	case westSide:
		x--
	case southwestSide:
		y++
		if gridPos.y%2 == 1 {
			x--
		}
	case northwestSide:
		y--
		if gridPos.y%2 == 1 {
			x--
		}
	default:
		panic("Cannot find adjacent grid position. Unknown hex side.")
	}

	if x < 0 || y < 0 || y >= hc.yIndexMax || (y%2 == 0 && x >= hc.xIndexMax-1) || (y%2 == 1 && x >= hc.xIndexMax) {
		return gridPosition{}, false
	} else {
		return gridPosition{x, y}, true
	}
}

func (gridPos gridPosition) isAdjacent(other gridPosition, hc *hexConfig) bool {
	ne, neOk := gridPos.adjacent(northeastSide, hc)
	e, eOk := gridPos.adjacent(eastSide, hc)
	se, seOk := gridPos.adjacent(southeastSide, hc)
	sw, swOk := gridPos.adjacent(southwestSide, hc)
	w, wOk := gridPos.adjacent(westSide, hc)
	nw, nwOk := gridPos.adjacent(northwestSide, hc)

	return (neOk && ne.equals(other)) ||
		(eOk && e.equals(other)) ||
		(seOk && se.equals(other)) ||
		(swOk && sw.equals(other)) ||
		(wOk && w.equals(other)) ||
		(nwOk && nw.equals(other))
}

func (pos position) log() []slog.Attr {
	var a = []slog.Attr{}
	var layer string
	switch pos.layer {
	case noLayer:
		layer = "noLayer"
	case ground:
		layer = "ground"
	case lowAltitude:
		layer = "lowAltitude"
	case highAltitude:
		layer = "highAltitude"
	}
	a = append(a, slog.String("layer", layer))
	a = append(a, pos.PixelPosition.log()...)
	return a
}

func (pos PixelPosition) add(other PixelPosition, hc *hexConfig) (PixelPosition, bool) {
	x := pos.X + other.X
	y := pos.Y + other.Y

	if x < 0 || x >= hc.windowWidth || y < 0 || y >= hc.windowHeight {
		return PixelPosition{}, false
	} else {
		return PixelPosition{x, y}, true
	}
}

func (gridPos gridPosition) getPixelPos(hexPoint hexPoint, hc *hexConfig) (PixelPosition, bool) {
	pixelPos := gridPos.getPixelPosOutbounds(hexPoint, hc)
	x, y := pixelPos.X, pixelPos.Y

	if x < 0 || x >= hc.windowWidth || y < 0 || y >= hc.windowHeight {
		return PixelPosition{}, false
	} else {
		return PixelPosition{x, y}, true
	}
}

func (gridPos gridPosition) getPixelPosOutbounds(hexPoint hexPoint, hc *hexConfig) PixelPosition {
	x := gridPos.x * hc.hexWidth
	y := gridPos.y * (hc.hexHeight / 4 * 3)
	if gridPos.y%2 == 1 {
		x -= (hc.hexWidth) / 2
	}

	// The values are found from the 1. hex look up table defined in hexConfig.go
	// and 2. arbitrary rule of hexagons sides are left inclusive and right exclusive.
	// These offsets are found from the hex look up table start point which is
	// located at (xOffset, yOffset).
	switch hexPoint {
	case center:
		y -= hc.hexHeight / 4
	case northPoint:
		y -= hc.hexHeight / 4 * 3
	case northwestPoint:
		x -= hc.hexWidth / 2
		y -= hc.hexHeight / 2
	case northeastPoint:
		x += hc.hexWidth/2 - 1
		y -= hc.hexHeight / 2
	case southPoint:
		y += hc.hexHeight/4 - 1
	case southwestPoint:
		x -= hc.hexWidth / 2
		y -= 1
	case southeastPoint:
		x += hc.hexWidth/2 - 1
		y -= 1
	}

	return PixelPosition{x + hc.xOffset, y + hc.yOffset}
}

func (pixelPos PixelPosition) getGridPos(hc *hexConfig) (gridPosition, bool) {
	if pixelPos.X < 0 || pixelPos.X > hc.windowWidth-1 || pixelPos.Y < 0 || pixelPos.Y > hc.windowHeight-1 {
		return gridPosition{}, false
	}
	return pixelPos.getGridPosOutbounds(hc)
}

func (pixelPos PixelPosition) getGridPosOutbounds(hc *hexConfig) (gridPosition, bool) {
	xIndex := (pixelPos.X+hc.hexWidth-hc.xOffset)/hc.hexWidth - 1
	yIndex := (pixelPos.Y+(hc.hexHeight/4*3)-hc.yOffset)/(hc.hexHeight/4*3) - 2
	innerX := (pixelPos.X + hc.hexWidth - hc.xOffset) % hc.hexWidth
	innerY := (pixelPos.Y + (hc.hexHeight * 3 / 2) - hc.yOffset) % (hc.hexHeight * 3 / 2)

	firstYSectionStart := 0
	secondYSectionStart := firstYSectionStart + hc.hexHeight/4
	thirdYSectionStart := secondYSectionStart + hc.hexHeight/2
	fourthYSectionStart := thirdYSectionStart + hc.hexHeight/4

	var xIndexOffset int
	var yIndexOffset int
	switch {
	case innerY >= firstYSectionStart && innerY < secondYSectionStart:
		triangle := hc.lookUpTable[innerX+innerY*hc.hexWidth]
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
		hexLookUpYOffset := (hc.hexWidth * hc.hexHeight / 4) * 3
		var hexLookUpXOffset int
		if innerX >= 0 && innerX < (hc.hexWidth/2) {
			hexLookUpXOffset = (hc.hexWidth / 2)
		} else {
			hexLookUpXOffset = -(hc.hexWidth / 2)
		}

		triangle := hc.lookUpTable[innerX+(innerY*hc.hexWidth)+hexLookUpXOffset-hexLookUpYOffset]
		switch {
		case triangle == leftHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		case triangle == middleHexTriangle && innerX < (hc.hexWidth/2):
			xIndexOffset = 0
			yIndexOffset = 2
		case triangle == middleHexTriangle && innerX >= (hc.hexWidth/2):
			xIndexOffset = 1
			yIndexOffset = 2
		case triangle == rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= fourthYSectionStart && innerY < (hc.hexHeight*3/2):
		if innerX >= 0 && innerX < (hc.hexWidth/2) {
			xIndexOffset = 0
			yIndexOffset = 2
		} else {
			xIndexOffset = 1
			yIndexOffset = 2
		}
	default:
		panic(fmt.Sprintf("getHexGridPostion could not find a y section for the inner y: %d for pixel position %s", innerY, pixelPos.log()))
	}

	return gridPosition{xIndex + xIndexOffset, yIndex + yIndexOffset}, true
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

func collidesWithBorder(pos position, offset int, hc *hexConfig) bool {
	if pos.X+offset > hc.windowWidth ||
		pos.X-offset < 0 ||
		pos.Y+offset > hc.windowHeight ||
		pos.Y-offset < 0 {
		return true
	} else {
		return false
	}
}
