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
	// Never use for game logic. float32 rather than idiomatic float64
	// because alot of Ebiten draw APIs use float32
	windowMultiplier float32
	windowWidth      int
	windowHeight     int
	gridSize         int
	// The debug grid helps debugging display issues. Usually only altering
	// offset is needed, width can remain 1.
	// Offset values should only be 0 - 20 (half of grid size without window multipler)
	// Some helpful values are: 0, 0.5, 1, 19, 19.5, 20
	debugGridOffset             int
	debugGridWidth              int
	debugGridIndexCount         int
	debugGridColor              color.Color
	selectectionBoxLineWidth    int
	missileBlastRadiusLineWidth int
	healthLineWidth             int
	fluxLineWidth               int
	fluxDirectionRadius         int
	targetNodeLineWidth         int
	backgroundColor             color.RGBA
	highlightColor              color.RGBA
	unitHealthColor             color.RGBA
	nodeHealthColor             color.RGBA
	fluxColor                   color.RGBA
	bluePlayerNodeColor         color.RGBA
	bluePlayerUnitColor         color.RGBA
	redPlayerNodeColor          color.RGBA
	redPlayerUnitColor          color.RGBA
	neutralNodeColor            color.RGBA
}

func newDrawContext() *drawContext {
	highlightColor := color.RGBA{70, 255, 0, 200}
	dc := drawContext{
		screen:           nil,
		windowMultiplier: 2.5,
		// aspect ratio config (16:9 aspect ratio)
		windowWidth:                 1280,
		windowHeight:                720,
		gridSize:                    40,
		debugGridOffset:             0,
		debugGridWidth:              1,
		debugGridIndexCount:         32,
		debugGridColor:              color.Black,
		selectectionBoxLineWidth:    1,
		missileBlastRadiusLineWidth: 1,
		healthLineWidth:             2,
		fluxLineWidth:               2,
		fluxDirectionRadius:         3,
		targetNodeLineWidth:         1,
		backgroundColor:             color.RGBA{200, 200, 200, 255},
		highlightColor:              highlightColor,
		unitHealthColor:             highlightColor,
		nodeHealthColor:             highlightColor,
		fluxColor:                   highlightColor,
		bluePlayerNodeColor:         color.RGBA{0, 0, 255, 100},
		bluePlayerUnitColor:         color.RGBA{0, 0, 255, 100},
		redPlayerNodeColor:          color.RGBA{255, 0, 0, 100},
		redPlayerUnitColor:          color.RGBA{255, 0, 0, 100},
		neutralNodeColor:            color.RGBA{100, 60, 100, 200},
	}
	return &dc
}

func drawUnfilledBoxFloat32(dc drawContext, color color.Color, initialX float32, initialY float32, finalX float32, finalY float32) {
	initialXGraphical := float32(initialX) * dc.windowMultiplier
	initialYGraphical := float32(initialY) * dc.windowMultiplier
	finalXGraphical := float32(finalX) * dc.windowMultiplier
	finalYGraphical := float32(finalY) * dc.windowMultiplier
	selectectionBoxLineWidth := float32(dc.selectectionBoxLineWidth) * dc.windowMultiplier

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

	for i := range dc.debugGridIndexCount {
		var x = dc.gridSize * i
		vector.StrokeLine(
			dc.screen,
			float32(x-dc.debugGridOffset)*dc.windowMultiplier,
			0, float32(x-dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.windowHeight)*dc.windowMultiplier,
			float32(dc.debugGridWidth)*dc.windowMultiplier,
			dc.debugGridColor,
			false)
		vector.StrokeLine(
			dc.screen,
			float32(x+dc.debugGridOffset)*dc.windowMultiplier,
			0, float32(x+dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.windowHeight)*dc.windowMultiplier,
			float32(dc.debugGridWidth)*dc.windowMultiplier,
			dc.debugGridColor,
			false)
		var y = dc.gridSize * i
		vector.StrokeLine(
			dc.screen,
			0, float32(y-dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.windowWidth)*dc.windowMultiplier,
			float32(y-dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.debugGridWidth)*dc.windowMultiplier,
			dc.debugGridColor,
			false)
		vector.StrokeLine(
			dc.screen,
			0,
			float32(y+dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.windowWidth)*dc.windowMultiplier,
			float32(y+dc.debugGridOffset)*dc.windowMultiplier,
			float32(dc.debugGridWidth)*dc.windowMultiplier,
			dc.debugGridColor,
			false)
	}
}
