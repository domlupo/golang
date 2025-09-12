// general position utility functions
package main

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

// unreadable high resolution position
type pixelPosition struct {
	x int
	y int
}

func getPixelPosition(gridPos gridPosition, innerGridPos innerGridPosition) pixelPosition {
	var x = (gridPos.x * gridSize) + innerGridPos.x
	var y = (gridPos.y * gridSize) + innerGridPos.y

	return pixelPosition{x, y}
}

func splitPixelPosition(pixelPos pixelPosition) (gridPosition, innerGridPosition) {
	var gridPosX = pixelPos.x / gridSize
	var gridPosY = pixelPos.y / gridSize
	var innerGridPosX = pixelPos.x % gridSize
	var innerGridPosY = pixelPos.y % gridSize
	return gridPosition{gridPosX, gridPosY}, innerGridPosition{innerGridPosX, innerGridPosY}
}

func collidesWithBorder(pixelPos pixelPosition, offset int) bool {
	if pixelPos.x+offset > windowWidth {
		return true
	} else if pixelPos.x-offset < 0 {
		return true
	} else if pixelPos.y+offset > windowHeight {
		return true
	} else if pixelPos.y-offset < 0 {
		return true
	}

	return false
}
