package main

// file contains general position utility functions

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
