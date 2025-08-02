package main

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
	var x = (gridPos.x * gridSizeWithoutMultiplier) + innerGridPos.x
	var y = (gridPos.y * gridSizeWithoutMultiplier) + innerGridPos.y

	return pixelPosition{x, y}
}

func getBeltGridPosFromShapePixelPos(pixelPos pixelPosition) gridPosition {
	var x = ((pixelPos.x + (shapeSize / 2)) / gridSizeWithoutMultiplier)
	var y = ((pixelPos.y + (shapeSize / 2)) / gridSizeWithoutMultiplier)

	return gridPosition{x, y}
}
