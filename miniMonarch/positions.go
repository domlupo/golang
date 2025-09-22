// general position utility functions
package main

import (
	"log/slog"
)

type hexagonTriangle int

// The hexagon position context is a hardcoded mapping to aid with pixel position : relative hexagon.
// It avoids the the need of CPU intensive math to convert from pixel position -> hexagon.
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
// ***/ \
// **/   \
// */     \
//
// non right triangle in middle
//    /*\
//   /***\
//  /*****\
//
// right triangle on right side
//    / \***
//   /   \**
//  /     \*

const (
	unknownHexTriangle = iota
	leftHexTriangle
	middleHexTriangle
	rightHexTriangle
)

// TODO: Remove hard coded hexagon configs

func newHexPositionContext() [460]hexagonTriangle {
	var pc [460]hexagonTriangle

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

func setHexPositionContext(pc *[460]hexagonTriangle, rowPixelCount int, rowStart int, rowMiddleOffset int) {
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

func hexPosContextIndex(innerGridPos innerGridPosition) int {
	return innerGridPos.x + innerGridPos.y*46
}

type direction int

// typically, iterate CW north to northWest throughout codebase
const (
	noDirection = iota
	north
	northEast
	east
	southEast
	south
	southWest
	west
	northWest
)

// readable low resolution position used for grid items e.g. nodes
type gridPosition struct {
	x int
	y int
}

func (pos gridPosition) adjacent(d direction, maxX int, maxY int) *gridPosition {
	var gridPos *gridPosition
	switch d {
	case noDirection:
		panic("need a direction to find adjacent grid position")
	case north:
		if pos.y == 0 {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x, pos.y - 1}
		}
	case northEast:
		if pos.x == maxX || pos.y == 0 {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x + 1, pos.y - 1}
		}
	case east:
		if pos.x == maxX {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x + 1, pos.y}
		}
	case southEast:
		if pos.x == maxX || pos.y == maxY {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x + 1, pos.y + 1}
		}
	case south:
		if pos.y == maxY {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x, pos.y + 1}
		}
	case southWest:
		if pos.x == 0 || pos.y == maxY {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x - 1, pos.y + 1}
		}
	case west:
		if pos.x == 0 {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x - 1, pos.y}
		}
	case northWest:
		if pos.x == 0 || pos.y == 0 {
			gridPos = nil
		} else {
			gridPos = &gridPosition{pos.x - 1, pos.y - 1}
		}
	}

	return gridPos
}

// readable high resolution position used for non grid items e.g. units
type innerGridPosition struct {
	x int
	y int
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

type offset int

const (
	topLeft = iota
	center
)

func getPixelPosition(gridPos gridPosition, offset offset, gridSize int) pixelPosition {
	var innerGridPos innerGridPosition
	switch offset {
	case topLeft:
		innerGridPos = innerGridPosition{0, 0}
	case center:
		innerGridPos = innerGridPosition{gridSize / 2, gridSize / 2}
	}

	x := (gridPos.x * gridSize) + innerGridPos.x
	y := (gridPos.y * gridSize) + innerGridPos.y
	return pixelPosition{x, y}
}

func splitPixelPosition(pixelPos pixelPosition, gridSize int) (gridPosition, innerGridPosition) {
	gridPosX := pixelPos.x / gridSize
	gridPosY := pixelPos.y / gridSize
	innerGridPosX := pixelPos.x % gridSize
	innerGridPosY := pixelPos.y % gridSize
	return gridPosition{gridPosX, gridPosY}, innerGridPosition{innerGridPosX, innerGridPosY}
}
