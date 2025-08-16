package main

import "fmt"

// readable low resolution position used for grid items like belts, consumers, etc
type gridPosition struct {
	x int // 0 - 32
	y int // 0 - 18
}

// readable high resolution position used for non grid items like offers
type innerGridPosition struct {
	x int
	y int
}

// unreadable high resolution position used for non grid items like offers
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

func splitPixelPositionForOffer(pixelPos pixelPosition) (gridPosition, innerGridPosition) {
	var x = pixelPos.x + (offerSize / 2)
	var y = pixelPos.y + (offerSize / 2)

	var gridPosX = x / gridSize
	var gridPosY = y / gridSize
	var innerGridPosX = x % gridSize
	var innerGridPosY = y % gridSize

	return gridPosition{gridPosX, gridPosY}, innerGridPosition{innerGridPosX, innerGridPosY}
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

func getBeltGridPosFromOfferPixelPos(pixelPos pixelPosition) gridPosition {
	var x = ((pixelPos.x + (offerSize / 2)) / gridSize)
	var y = ((pixelPos.y + (offerSize / 2)) / gridSize)

	return gridPosition{x, y}
}

func getDirectionForOffer(bd beltDirection, pixelPos pixelPosition) direction {
	var innerGridPos = innerGridPosition{(pixelPos.x + 7) % gridSize, (pixelPos.y + 7) % gridSize}

	switch bd {
	case bDown:
		return down
	case bDownLeft:
		if innerGridPos.y >= innerGridPos.x {
			return left
		} else {
			return down
		}
	case bDownRight:
		if (gridSize - innerGridPos.x - 1) <= innerGridPos.y {
			return right
		} else {
			return down
		}
	case bLeft:
		return left
	case bLeftDown:
		if innerGridPos.x > innerGridPos.y {
			return left
		} else {
			return down
		}
	case bLeftUp:
		if (gridSize - innerGridPos.x - 1) < innerGridPos.y {
			return left
		} else {
			return up
		}
	case bRight:
		return right
	case bRightDown:
		if (innerGridPos.x + innerGridPos.y) < gridSize-1 {
			return right
		} else {
			return down
		}
	case bRightUp:
		if innerGridPos.y > innerGridPos.x {
			return right
		} else {
			return up
		}
	case bUp:
		return up
	case bUpLeft:
		if (innerGridPos.x + innerGridPos.y) < gridSize {
			return left
		} else {
			return up
		}
	case bUpRight:
		if innerGridPos.x >= innerGridPos.y {
			return right
		} else {
			return up
		}
	}
	panic(fmt.Sprintf("unexpected main.beltDirection: %#v", bd))
}
