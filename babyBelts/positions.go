package main

import "fmt"

// readable low resolution position used for grid items like belts, movers, consumers, etc
type gridPosition struct {
	x int // 0 - 32
	y int // 0 - 18
}

// readable high resolution position used for non grid items like shapes
type innerGridPosition struct {
	x int
	y int
}

// unreadable high resolution position used for non grid items like shape
type pixelPosition struct {
	x int
	y int
}

func getPixelPosition(gridPos gridPosition, innerGridPos innerGridPosition) pixelPosition {
	var x = (gridPos.x * gridSize) + innerGridPos.x
	var y = (gridPos.y * gridSize) + innerGridPos.y

	return pixelPosition{x, y}
}

func getPixelPositionFromDirection(direction direction) pixelPosition {
	switch direction {
	case down:
		return pixelPosition{0, 1}
	case left:
		return pixelPosition{-1, 0}
	case right:
		return pixelPosition{1, 0}
	case up:
		return pixelPosition{0, -1}
	default:
		panic(fmt.Sprintf("unexpected main.direction: %#v", direction))
	}
}

func getBeltGridPosFromShapePixelPos(pixelPos pixelPosition) gridPosition {
	var x = ((pixelPos.x + (shapeSize / 2)) / gridSize)
	var y = ((pixelPos.y + (shapeSize / 2)) / gridSize)

	return gridPosition{x, y}
}

func getDirectionForShape(bd beltDirection, pixelPos pixelPosition) direction {
	var innerGridPos = innerGridPosition{pixelPos.x % gridSize, pixelPos.y % gridSize}
	var smallOffset = (gridSize / 2) - shapeSize - (beltMiddleLineWidth / 2)
	var largeOffset = (gridSize / 2) + (beltMiddleLineWidth / 2)

	switch bd {
	case bDown:
		return down
	case bDownLeft:
		if innerGridPos.x == largeOffset && innerGridPos.y != largeOffset {
			// outer piece
			return down
		} else if innerGridPos.x == smallOffset && innerGridPos.y != smallOffset && (innerGridPos.y <= 4 || innerGridPos.y > 32) {
			// inner piece
			return down
		} else {
			return left
		}
	case bDownRight:
		if innerGridPos.x == smallOffset && innerGridPos.y != largeOffset {
			// outer piece
			return down
		} else if innerGridPos.x == largeOffset && innerGridPos.y != largeOffset && (innerGridPos.y <= 4 || innerGridPos.y > 32) {
			// inner piece
			return down
		} else {
			return right
		}
	case bLeft:
		return left
	case bLeftDown:
		if innerGridPos.y == smallOffset && innerGridPos.x != smallOffset {
			// outer piece
			return left
		} else if innerGridPos.y == largeOffset && innerGridPos.x != largeOffset && (innerGridPos.x > smallOffset) {
			// inner piece
			return left
		} else {
			return down
		}
	case bLeftUp:
		if innerGridPos.y == largeOffset && innerGridPos.x != smallOffset {
			// outer piece
			return left
		} else if innerGridPos.y == smallOffset && innerGridPos.x != largeOffset && (innerGridPos.x > largeOffset) {
			// inner piece
			return left
		} else {
			return up
		}
	case bRight:
		return right
	case bRightDown:
		if innerGridPos.y == smallOffset && innerGridPos.x != largeOffset {
			// inner piece
			return right
		} else if innerGridPos.y == largeOffset && (innerGridPos.x <= 4 || innerGridPos.x > 32) {
			// outer piece
			return right
		} else {
			return down
		}
	case bRightUp:
		if innerGridPos.y == smallOffset && innerGridPos.x != smallOffset && (innerGridPos.x <= 4 || innerGridPos.x > 32) {
			// inner piece
			return right
		} else if innerGridPos.y == largeOffset && innerGridPos.x != largeOffset {
			// outer piece
			return right
		} else {
			return up
		}
	case bUp:
		return up
	case bUpLeft:
		if innerGridPos.x == largeOffset && innerGridPos.y != smallOffset {
			// outer piece
			return up
		} else if innerGridPos.x == smallOffset && innerGridPos.y != largeOffset && innerGridPos.y > smallOffset {
			// inner piece
			return up
		} else {
			return left
		}
	case bUpRight:
		if innerGridPos.x == smallOffset && innerGridPos.y != smallOffset {
			// outer piece
			return up
		} else if innerGridPos.x == largeOffset && innerGridPos.y != largeOffset && innerGridPos.y > smallOffset {
			// inner piece
			return up
		} else {
			return right
		}
	}
	panic(fmt.Sprintf("unexpected main.beltDirection: %#v", bd))
}
