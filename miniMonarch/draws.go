package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// file contains general drawing utiliy functions

func drawUnfilledBoxFloat32(screen *ebiten.Image, color color.Color, initialX float32, initialY float32, finalX float32, finalY float32) {
	var initialXGraphical = float32(initialX) * windowMultiplier
	var initialYGraphical = float32(initialY) * windowMultiplier
	var finalXGraphical = float32(finalX) * windowMultiplier
	var finalYGraphical = float32(finalY) * windowMultiplier

	// TODO: fix this to not be global
	var selectectionBoxLineWidth = float32(selectectionBoxLineWidth) * windowMultiplier

	// draw top line
	vector.StrokeLine(
		screen,
		initialXGraphical,
		initialYGraphical,
		finalXGraphical,
		initialYGraphical,
		selectectionBoxLineWidth,
		color,
		false,
	)
	//  draw left line
	vector.StrokeLine(
		screen,
		initialXGraphical,
		initialYGraphical,
		initialXGraphical,
		finalYGraphical,
		selectectionBoxLineWidth,
		color,
		false,
	)
	//  draw right line
	vector.StrokeLine(
		screen,
		finalXGraphical,
		initialYGraphical,
		finalXGraphical,
		finalYGraphical,
		selectectionBoxLineWidth,
		color,
		false,
	)
	//  draw bottom line
	vector.StrokeLine(
		screen,
		initialXGraphical,
		finalYGraphical,
		finalXGraphical,
		finalYGraphical,
		selectectionBoxLineWidth,
		color,
		false,
	)
}

func drawUnfilledBoxInt(screen *ebiten.Image, color color.Color, initialX int, initialY int, finalX int, finalY int) {
	drawUnfilledBoxFloat32(screen, color, float32(initialX), float32(initialY), float32(finalX), float32(finalY))
}

func drawDebugGrid(screen *ebiten.Image, drawDebugGrid bool) {
	if !drawDebugGrid {
		return
	}

	for i := range debugGridSize {
		var x = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(
			screen,
			x-debugGridOffset,
			0, x-debugGridOffset,
			float32(windowHeight*windowMultiplier),
			debugGridWidth,
			debugGridColor,
			false)
		vector.StrokeLine(
			screen,
			x+debugGridOffset,
			0, x+debugGridOffset,
			float32(windowHeight*windowMultiplier),
			debugGridWidth,
			debugGridColor,
			false)
		var y = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(
			screen,
			0, y-debugGridOffset,
			float32(windowWidth*windowMultiplier),
			y-debugGridOffset,
			debugGridWidth,
			debugGridColor,
			false)
		vector.StrokeLine(
			screen,
			0,
			y+debugGridOffset,
			float32(windowWidth*windowMultiplier),
			y+debugGridOffset,
			debugGridWidth,
			debugGridColor,
			false)
	}
}
