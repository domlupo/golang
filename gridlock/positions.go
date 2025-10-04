// general position utility functions. Most are pure functions. Used across many different files.
package main

import (
	"fmt"
	"log/slog"
)

type hexagon struct {
	linkedList *unit
	unitCount  int
}

// developer readable low resolution position used for grid items e.g. nodes
type gridPosition struct {
	x int
	y int
}

// developer unreadable high resolution position
type pixelPosition struct {
	x int
	y int
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
	pixelPosition
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

type line struct {
	start pixelPosition
	end   pixelPosition
}

type pointOrientation int

const (
	noOrientation pointOrientation = iota
	counterClockwise
	clockwise
	collinear
)

func (pos position) isZero() bool {
	return pos.layer == noLayer
}

func (pixelPos pixelPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", pixelPos.x))
	a = append(a, slog.Int("y", pixelPos.y))
	return a
}

func (pos gridPosition) adjacent(side hexSide, hc *hexConfig) (gridPosition, bool) {
	x := pos.x
	y := pos.y
	switch side {
	case noSide:
		panic("need a direction to find adjacent grid position")
	case northeastSide:
		y--
		if pos.y%2 == 0 {
			x++
		}
	case eastSide:
		x++
	case southeastSide:
		y++
		if pos.y%2 == 0 {
			x++
		}
	case westSide:
		x--
	case southwestSide:
		y++
		if pos.y%2 == 1 {
			x--
		}
	case northwestSide:
		y--
		if pos.y%2 == 1 {
			x--
		}
	default:
		panic("No adjacent")
	}

	if x < 0 || x >= hc.xIndexMax || y < 0 || y >= hc.yIndexMax {
		return gridPosition{}, false
	} else {
		return gridPosition{x, y}, true
	}
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
	a = append(a, pos.pixelPosition.log()...)
	return a
}

func (pos pixelPosition) add(other pixelPosition, hc *hexConfig) (pixelPosition, bool) {
	x := pos.x + other.x
	y := pos.y + other.y

	if x < 0 || x >= hc.windowWidth || y < 0 || y >= hc.windowHeight {
		return pixelPosition{}, false
	} else {
		return pixelPosition{x, y}, true
	}
}

func (gridPos gridPosition) getPixelPos(hexPoint hexPoint, hc *hexConfig) (pixelPosition, bool) {
	pixelPos := gridPos.getPixelPosOutbounds(hexPoint, hc)
	x, y := pixelPos.x, pixelPos.y

	if x < 0 || x >= hc.windowWidth || y < 0 || y >= hc.windowHeight {
		return pixelPosition{}, false
	} else {
		return pixelPosition{x, y}, true
	}
}

func (gridPos gridPosition) getPixelPosOutbounds(hexPoint hexPoint, hc *hexConfig) pixelPosition {
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

	return pixelPosition{x + hc.xOffset, y + hc.yOffset}
}

func (pixelPos pixelPosition) getGridPos(hc *hexConfig) (gridPosition, bool) {
	if pixelPos.x < 0 || pixelPos.x > hc.windowWidth-1 || pixelPos.y < 0 || pixelPos.y > hc.windowHeight-1 {
		return gridPosition{}, false
	}

	xIndex := (pixelPos.x+hc.hexWidth-hc.xOffset)/hc.hexWidth - 1
	yIndex := (pixelPos.y+(hc.hexHeight/4*3)-hc.yOffset)/(hc.hexHeight/4*3) - 2
	var xIndexOffset int
	var yIndexOffset int
	innerX := (pixelPos.x + hc.hexWidth - hc.xOffset) % hc.hexWidth
	innerY := (pixelPos.y + (hc.hexHeight * 3 / 2) - hc.yOffset) % (hc.hexHeight * 3 / 2)

	firstYSectionStart := 0
	secondYSectionStart := firstYSectionStart + hc.hexHeight/4
	thirdYSectionStart := secondYSectionStart + hc.hexHeight/2
	fourthYSectionStart := thirdYSectionStart + hc.hexHeight/4

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

func orientation(p pixelPosition, q pixelPosition, r pixelPosition) pointOrientation {
	v := (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
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

func onSegment(p pixelPosition, q pixelPosition, r pixelPosition) bool {
	return q.x <= max(p.x, r.x) && q.x >= min(p.x, r.x) && q.y <= max(p.y, r.y) && q.y >= min(p.y, r.y)
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
