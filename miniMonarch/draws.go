// basic drawing utility functions
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type drawContext struct {
	screen *ebiten.Image
	// windowMultipler must only be used for drawing images and windows.
	// Never use for game logic. float32 rather than float64 because
	// alot of Ebiten draw APIs use float32
	windowMultiplier float32
	windowWidth      int
	windowHeight     int
	gridSize         int
	// TODO: Determine if debugGridOffset can be simplified into one var
	debugGridOffsetWithoutWindowMultipler int
	debugGridOffset                       float32
	debugGridWidth                        float32
}

func newDrawContext() *drawContext {
	var windowMultiplier float32 = 2.5
	debugGridOffsetWithoutWindowMultipler := 0
	dc := drawContext{
		screen:           nil,
		windowMultiplier: windowMultiplier,
		// aspect ratio config (16:9 aspect ratio)
		windowWidth:                           1280,
		windowHeight:                          720,
		gridSize:                              40,
		debugGridOffsetWithoutWindowMultipler: 0,
		// TODO: determine if windowMultipler should be removed from here. I think yes
		debugGridOffset: float32(debugGridOffsetWithoutWindowMultipler) * windowMultiplier,
		debugGridWidth:  1 * windowMultiplier,
	}
	return &dc
}

func drawUnfilledBoxFloat32(dc drawContext, color color.Color, initialX float32, initialY float32, finalX float32, finalY float32) {
	var initialXGraphical = float32(initialX) * dc.windowMultiplier
	var initialYGraphical = float32(initialY) * dc.windowMultiplier
	var finalXGraphical = float32(finalX) * dc.windowMultiplier
	var finalYGraphical = float32(finalY) * dc.windowMultiplier

	// TODO: fix this to not be global
	var selectectionBoxLineWidth = float32(selectectionBoxLineWidth) * dc.windowMultiplier

	// draw top line
	vector.StrokeLine(
		dc.screen,
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
		dc.screen,
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
		dc.screen,
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
		dc.screen,
		initialXGraphical,
		finalYGraphical,
		finalXGraphical,
		finalYGraphical,
		selectectionBoxLineWidth,
		color,
		false,
	)
}

func drawUnfilledBoxInt(dc drawContext, color color.Color, initialX int, initialY int, finalX int, finalY int) {
	drawUnfilledBoxFloat32(dc, color, float32(initialX), float32(initialY), float32(finalX), float32(finalY))
}

func drawDebug(dc drawContext, drawDebugGrid bool) {
	if !drawDebugGrid {
		return
	}

	for i := range debugGridSquareCount {
		var x = float32(dc.gridSize*i) * dc.windowMultiplier
		vector.StrokeLine(
			dc.screen,
			x-dc.debugGridOffset,
			0, x-dc.debugGridOffset,
			float32(dc.windowHeight)*dc.windowMultiplier,
			dc.debugGridWidth,
			debugGridColor,
			false)
		vector.StrokeLine(
			dc.screen,
			x+dc.debugGridOffset,
			0, x+dc.debugGridOffset,
			float32(dc.windowHeight)*dc.windowMultiplier,
			dc.debugGridWidth,
			debugGridColor,
			false)
		var y = float32(dc.gridSize*i) * dc.windowMultiplier
		vector.StrokeLine(
			dc.screen,
			0, y-dc.debugGridOffset,
			float32(dc.windowWidth)*dc.windowMultiplier,
			y-dc.debugGridOffset,
			dc.debugGridWidth,
			debugGridColor,
			false)
		vector.StrokeLine(
			dc.screen,
			0,
			y+dc.debugGridOffset,
			float32(dc.windowWidth)*dc.windowMultiplier,
			y+dc.debugGridOffset,
			dc.debugGridWidth,
			debugGridColor,
			false)
	}
}
