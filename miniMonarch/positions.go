// general position utility functions
package main

import (
	"fmt"
	"log/slog"
)

// This global should only be used for function signatures, structs, etc. It can
// be calculated during runtime via hexConfig fields. It is calculated from hexagon width
// * (hexagon height/4). The 4 is from the hexagon position context being 1/4 of the total
// hexagon area.
const hexagonPostitionContextPixelCount = 460

type hexConfig struct {
	hexPositionContext [hexagonPostitionContextPixelCount]hexagonTriangle
	hexWidth           int
	hexHeight          int
	xOffset            int
	yOffset            int
	windowWidth        int
	windowHeight       int
}

func newHexConfig(windowWidth int, windowHeight int) hexConfig {
	hexPositionContext := newHexPositionContext()
	hc := hexConfig{
		hexPositionContext: hexPositionContext,
		hexWidth:           46,
		hexHeight:          40,
		xOffset:            19,
		yOffset:            10,
		windowWidth:        windowWidth,
		windowHeight:       windowHeight,
	}
	return hc
}

// The hexagon position context is a hardcoded mapping to aid with converting
// pixel position to hexagon grid position. It avoids the the need of CPU intensive math.
// It helps CPU performance by removing the need for floating point or sqrt math logic.
//
// Specific CPU profiling did not prove this optimizaton was needed but prior CPU profiling
// for this game showed that integer math is much more efficient. Other pixel position to
// hexagon algorithms appeared as complex or even more complex than this method so this
// more optimized solution was implemented first.
//
// It is small portion of the hexagonal grid that is repeated multiple times.
// It covers 460 (46 cols, 10 rows) of pixels.
// It has a:
//   - right triangle on the left side.
//   - non right triangle in the middle.
//   - right triangle on the right side.
//
// right triangle on left side
//
// ***/ \
// **/   \
// */     \
//
// non right triangle in middle
//
//    /*\
//   /***\
//  /*****\
//
// right triangle on right side
//
//    / \***
//   /   \**
//  /     \*
//
// The hexagon position context is also used for another section of the hexagonal grid
// that is repeated multiple times. It is usally referenced as "inverted"
//
//  \     /
//   \   /
//    \ /
//
// This is possible because this section becomes the prior section by following these steps:
//   - horizontally flip the entire section.
//   - split entire section vertically in half.
//   - swap the two vertical halfs.
//
// An example highlighting the middle section is shown below:
//
// non right triangle in middle
//
//    /*\
//   /***\
//  /*****\
//
// horizontally flip the entire section.
//
//  \*****/
//   \***/
//    \*/
//
// split the entir section vertically in half
//
//  \** **/
//   \* */
//    \ /
//
// swap the two vertical halfs.
//
//  **/ \**
//  */   \*
//  /     \
//
// As shown above, the prior sections middile triangle is now converted into
// this sections left and right triangle.

type hexagonTriangle int

const (
	unknownHexTriangle = iota
	leftHexTriangle
	middleHexTriangle
	rightHexTriangle
)

func newHexPositionContext() [hexagonPostitionContextPixelCount]hexagonTriangle {
	var pc [hexagonPostitionContextPixelCount]hexagonTriangle

	rowPixelCount := 46

	// row 1
	rowStart := 0
	rowMiddleOffset := 1
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 2
	rowStart = rowPixelCount
	rowMiddleOffset = 3
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 3
	rowStart = rowPixelCount * 2
	rowMiddleOffset = 6
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 4
	rowStart = rowPixelCount * 3
	rowMiddleOffset = 8
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 5
	rowStart = rowPixelCount * 4
	rowMiddleOffset = 10
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 6
	rowStart = rowPixelCount * 5
	rowMiddleOffset = 13
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 7
	rowStart = rowPixelCount * 6
	rowMiddleOffset = 15
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 8
	rowStart = rowPixelCount * 7
	rowMiddleOffset = 17
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 9
	rowStart = rowPixelCount * 8
	rowMiddleOffset = 20
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 10
	rowStart = rowPixelCount * 9
	rowMiddleOffset = 22
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)

	return pc
}

func setHexPositionContext(pc *[hexagonPostitionContextPixelCount]hexagonTriangle, rowPixelCount int, rowStart int, rowMiddleOffset int) {
	pixelCount := rowPixelCount / 2

	leftHexStart := rowStart
	middleHexStart := rowStart + pixelCount - rowMiddleOffset
	rightHexStart := middleHexStart + 2*rowMiddleOffset
	for i := leftHexStart; i < middleHexStart; i++ {
		pc[i] = leftHexTriangle
	}
	for i := middleHexStart; i < rightHexStart; i++ {
		pc[i] = middleHexTriangle
	}
	for i := rightHexStart; i < rowStart+rowPixelCount; i++ {
		pc[i] = rightHexTriangle
	}
}

type direction int

const (
	noDirection = iota
	northEast
	east
	southEast
	southWest
	west
	northWest
)

// readable low resolution position used for grid items e.g. nodes
type gridPosition struct {
	x int
	y int
}

func (pos gridPosition) adjacent(d direction, xMax int, yMax int) (gridPosition, bool) {
	x := pos.x
	y := pos.y
	switch d {
	case noDirection:
		panic("need a direction to find adjacent grid position")
	case northEast:
		y--
		if pos.y%2 == 0 {
			x++
		}
	case east:
		x++
	case southEast:
		y++
		if pos.y%2 == 0 {
			x++
		}
	case west:
		x--
	case southWest:
		y++
		if pos.y%2 == 1 {
			x--
		}
	case northWest:
		y--
		if pos.y%2 == 1 {
			x--
		}
	default:
		panic("No adjacent")
	}

	if x < 0 || x > xMax || y < 0 || y > yMax {
		return gridPosition{}, false
	} else {
		return gridPosition{x, y}, true
	}
}

type layer int

const (
	noLayer = iota
	ground
	lowAltitude
	highAltitude
)

type position struct {
	layer
	pixelPosition
}

func (pos position) isZero() bool {
	return pos.layer == noLayer
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

// unreadable high resolution position
type pixelPosition struct {
	x int
	y int
}

func (pixelPos pixelPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", pixelPos.x))
	a = append(a, slog.Int("y", pixelPos.y))
	return a
}

type hexOffset int

const (
	center = iota
	topOffset
	topLeftOffset
	topRightOffset
	bottomOffset
	bottomLeftOffset
	bottomRightOffset
)

// TODO: write unit tests, did verify manually it worked with a few values
func getPixelPosition(gridPos gridPosition, offset hexOffset, hc *hexConfig) (pixelPosition, bool) {
	xPadding := (hc.windowWidth % hc.hexWidth) / 2
	x := xPadding + hc.hexWidth*gridPos.x
	if gridPos.y%2 == 1 {
		x = x - hc.hexWidth/2
	}
	y := (hc.hexHeight * 3 / 4) * gridPos.y

	switch offset {
	case center:
		// do nothing
	case topOffset:
		y -= hc.hexHeight / 2
	case topLeftOffset:
		x -= hc.hexWidth / 2
		y -= hc.hexHeight / 4
	case topRightOffset:
		x += hc.hexWidth / 2
		y -= hc.hexHeight / 4
	case bottomOffset:
		y += hc.hexHeight / 2
	case bottomLeftOffset:
		x -= hc.hexWidth / 2
		y += hc.hexHeight / 4
	case bottomRightOffset:
		x += hc.hexWidth / 2
		y += hc.hexHeight / 4
	}

	if x < 0 || x >= hc.windowWidth || y < 0 || y >= hc.windowHeight {
		return pixelPosition{}, false
	} else {
		return pixelPosition{x, y}, true
	}
}

func getGridPos(pixelPos pixelPosition, hc *hexConfig) (gridPosition, bool) {
	if pixelPos.x < 0 || pixelPos.x >= hc.windowWidth || pixelPos.y < 0 || pixelPos.y >= hc.windowHeight {
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
		triangle := hc.hexPositionContext[innerX+innerY*hc.hexWidth]
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
		hexPositionContextYOffset := (hc.hexWidth * hc.hexHeight / 4) * 3
		var hexPositionContextXOffset int
		if innerX >= 0 && innerX < (hc.hexWidth/2) {
			hexPositionContextXOffset = (hc.hexWidth / 2)
		} else {
			hexPositionContextXOffset = -(hc.hexWidth / 2)
		}

		triangle := hc.hexPositionContext[innerX+(innerY*hc.hexWidth)+hexPositionContextXOffset-hexPositionContextYOffset]
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
