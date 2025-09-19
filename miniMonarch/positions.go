// general position utility functions
package main

import "log/slog"

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
	x int // 0 - 32
	y int // 0 - 18
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

func getPixelPosition(gridPos gridPosition, innerGridPos innerGridPosition) pixelPosition {
	// TODO: fix hardcoding
	var x = (gridPos.x * 40) + innerGridPos.x
	var y = (gridPos.y * 40) + innerGridPos.y

	return pixelPosition{x, y}
}

func splitPixelPosition(pixelPos pixelPosition) (gridPosition, innerGridPosition) {
	// TODO: fix hardcoding
	var gridPosX = pixelPos.x / 40
	var gridPosY = pixelPos.y / 40
	var innerGridPosX = pixelPos.x % 40
	var innerGridPosY = pixelPos.y % 40
	return gridPosition{gridPosX, gridPosY}, innerGridPosition{innerGridPosX, innerGridPosY}
}

func collidesWithBorder(pos position, offset int) bool {
	if pos.pixelPosition.x+offset > 1280 { // TODO: Fix hardcoding
		return true
	} else if pos.pixelPosition.x-offset < 0 {
		return true
	} else if pos.pixelPosition.y+offset > 720 { // TODO: Fix hardcoding
		return true
	} else if pos.pixelPosition.y-offset < 0 {
		return true
	}

	return false
}
