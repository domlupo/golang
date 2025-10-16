// basic drawing configuration and utility functions
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type drawConfig struct {
	screen *ebiten.Image
	// windowMultiplier must only be used for drawing images and windows.
	// Never use for game logic. float32 rather than idiomatic float64
	// because alot of Ebiten draw APIs use float32
	windowMultiplier float32
	hexConfig
	debugGridWidth                int
	debugGridIndexCount           int
	debugGridColor                color.Color
	selectectionBoxLineWidth      int
	missileExplodeRadiusLineWidth int
	healthLineWidth               int
	fluxLineWidth                 int
	fluxDirectionRadius           int
	targetNodeLineWidth           int
	visibilityGridLineWidth       int
	roygbivColors                 []color.RGBA
	backgroundColor               color.RGBA
	highlightColor                color.RGBA
	unitHealthColor               color.RGBA
	nodeHealthColor               color.RGBA
	fluxColor                     color.RGBA
	bluePlayerNodeColor           color.RGBA
	bluePlayerUnitColor           color.RGBA
	redPlayerNodeColor            color.RGBA
	redPlayerUnitColor            color.RGBA
	neutralNodeColor              color.RGBA
	avgUnitPosRadius              int
	avgUnitPosColor               color.Color
	drawTriangleDefaultImage      *ebiten.Image
}

func newDrawConfig() *drawConfig {
	// TODO: check if this should use a subimage if final game uses draw triangle.
	// Not investigating now since draw triangle may not be used in final game.
	drawTriangleDefaultImage := ebiten.NewImage(1, 1)
	drawTriangleDefaultImage.Fill(color.White)
	highlightColor := color.RGBA{70, 255, 0, 200}
	dc := drawConfig{
		screen:           nil,
		windowMultiplier: 2.5,
		// aspect ratio config (16:9 aspect ratio)
		hexConfig:                     newHexConfig(1280, 720),
		debugGridWidth:                1,
		debugGridIndexCount:           32,
		debugGridColor:                color.Black,
		selectectionBoxLineWidth:      1,
		missileExplodeRadiusLineWidth: 1,
		healthLineWidth:               2,
		fluxLineWidth:                 2,
		fluxDirectionRadius:           3,
		targetNodeLineWidth:           2,
		visibilityGridLineWidth:       1,
		backgroundColor:               color.RGBA{200, 200, 200, 255},
		highlightColor:                highlightColor,
		unitHealthColor:               highlightColor,
		nodeHealthColor:               highlightColor,
		fluxColor:                     highlightColor,
		bluePlayerNodeColor:           color.RGBA{0, 0, 255, 100},
		bluePlayerUnitColor:           color.RGBA{0, 0, 255, 100},
		redPlayerNodeColor:            color.RGBA{255, 0, 0, 100},
		redPlayerUnitColor:            color.RGBA{255, 0, 0, 100},
		neutralNodeColor:              color.RGBA{100, 60, 100, 200},
		avgUnitPosRadius:              1,
		avgUnitPosColor:               color.Black,
		drawTriangleDefaultImage:      drawTriangleDefaultImage,
		roygbivColors: []color.RGBA{
			{200, 0, 0, 255},
			{250, 150, 0, 255},
			{250, 200, 0, 255},
			{0, 200, 150, 255},
			{0, 0, 250, 255},
			{0, 0, 150, 255},
			{150, 150, 250, 255},
		},
	}
	return &dc
}

func drawUnfilledBoxFloat32(dc drawConfig, color color.Color, initialX float32, initialY float32, finalX float32, finalY float32) {
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

func drawUnfilledBoxInt(dc drawConfig, color color.Color, initialX int, initialY int, finalX int, finalY int) {
	drawUnfilledBoxFloat32(dc, color, float32(initialX), float32(initialY), float32(finalX), float32(finalY))
}

func drawDebug(dc drawConfig, drawDebugGrid bool) {
	if !drawDebugGrid {
		return
	}

	// This is the pattern drawn.
	//    / \
	//   |
	//   |
	//
	// when repeated, it forms a hexagonal grid.
	//    / \ / \
	//   |   |
	//   |   |
	//  / \ / \ / \
	// |   |   |
	// |   |   |
	//    / \ / \
	//   |   |
	//   |   |

	for i := 0; i < dc.xIndexMax+1; i++ {
		for j := 0; j < dc.yIndexMax; j++ {
			lineWidth := float32(dc.debugGridWidth) * dc.windowMultiplier
			pos := gridPosition{i, j}

			n := pos.getPixelPosOutbounds(northPoint, &dc.hexConfig)
			nX := float32(n.x) * dc.windowMultiplier
			nY := float32(n.y) * dc.windowMultiplier

			nw := pos.getPixelPosOutbounds(northwestPoint, &dc.hexConfig)
			nwX := float32(nw.x) * dc.windowMultiplier
			nwY := float32(nw.y) * dc.windowMultiplier

			ne := pos.getPixelPosOutbounds(northeastPoint, &dc.hexConfig)
			neX := float32(ne.x+1) * dc.windowMultiplier
			neY := float32(ne.y) * dc.windowMultiplier

			sw := pos.getPixelPosOutbounds(southwestPoint, &dc.hexConfig)
			swX := float32(sw.x) * dc.windowMultiplier
			swY := float32(sw.y+1) * dc.windowMultiplier

			// draw west side
			//    / \
			//  *|*
			//  *|*
			vector.StrokeLine(
				dc.screen,
				nwX,
				nwY,
				swX,
				swY,
				lineWidth,
				dc.debugGridColor,
				false)

			// draw northwest
			//   */*\
			//   |
			//   |
			vector.StrokeLine(
				dc.screen,
				nX,
				nY,
				nwX,
				nwY,
				lineWidth,
				dc.debugGridColor,
				false)

			// draw northeast
			//    /*\*
			//   |
			//   |
			vector.StrokeLine(
				dc.screen,
				nX,
				nY,
				neX,
				neY,
				lineWidth,
				dc.debugGridColor,
				false)

		}
	}
}

type drawTriangleColor struct {
	red   float32
	green float32
	blue  float32
	alpha float32
}

func drawHex(dc drawConfig, gridPos gridPosition, color drawTriangleColor, selected bool) {
	center := gridPos.getPixelPosOutbounds(center, &dc.hexConfig)
	north := gridPos.getPixelPosOutbounds(northPoint, &dc.hexConfig)
	northeast := gridPos.getPixelPosOutbounds(northeastPoint, &dc.hexConfig)
	southeast := gridPos.getPixelPosOutbounds(southeastPoint, &dc.hexConfig)
	southwest := gridPos.getPixelPosOutbounds(southwestPoint, &dc.hexConfig)
	south := gridPos.getPixelPosOutbounds(southPoint, &dc.hexConfig)
	northwest := gridPos.getPixelPosOutbounds(northwestPoint, &dc.hexConfig)

	// The one pixel offsets for inclusive or exclusive hex points are handled at the game logic layer.
	// The window multiplier offsets are handled here since window multiplier is a draw concept.
	xCenter, yCenter := float32(center.x)*dc.windowMultiplier, float32(center.y)*dc.windowMultiplier
	xNorth, yNorth := float32(north.x)*dc.windowMultiplier, float32(north.y)*dc.windowMultiplier
	xNortheast, yNortheast := float32(northeast.x+1)*dc.windowMultiplier, float32(northeast.y)*dc.windowMultiplier
	xSoutheast, ySoutheast := float32(southeast.x+1)*dc.windowMultiplier, float32(southeast.y+1)*dc.windowMultiplier
	xSouthwest, ySouthwest := float32(southwest.x)*dc.windowMultiplier, float32(southwest.y+1)*dc.windowMultiplier
	xSouth, ySouth := float32(south.x)*dc.windowMultiplier, float32(south.y+1)*dc.windowMultiplier
	xNorthwest, yNorthwest := float32(northwest.x)*dc.windowMultiplier, float32(northwest.y)*dc.windowMultiplier

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5, 0, 5, 6, 0, 6, 1}
	vertices := []ebiten.Vertex{
		{DstX: xCenter, DstY: yCenter, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNorth, DstY: yNorth, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNortheast, DstY: yNortheast, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSoutheast, DstY: ySoutheast, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSouth, DstY: ySouth, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSouthwest, DstY: ySouthwest, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNorthwest, DstY: yNorthwest, SrcX: 0, SrcY: 0, ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
	}

	dc.screen.DrawTriangles(vertices, indices, dc.drawTriangleDefaultImage, &ebiten.DrawTrianglesOptions{})

	if selected {
		lineWidth := float32(dc.debugGridWidth) * dc.windowMultiplier

		north := gridPos.getPixelPosOutbounds(northPoint, &dc.hexConfig)
		nX := float32(north.x) * dc.windowMultiplier
		nY := float32(north.y) * dc.windowMultiplier

		nw := gridPos.getPixelPosOutbounds(northwestPoint, &dc.hexConfig)
		nwX := float32(nw.x) * dc.windowMultiplier
		nwY := float32(nw.y) * dc.windowMultiplier

		ne := gridPos.getPixelPosOutbounds(northeastPoint, &dc.hexConfig)
		neX := float32(ne.x+1) * dc.windowMultiplier
		neY := float32(ne.y) * dc.windowMultiplier

		sw := gridPos.getPixelPosOutbounds(southwestPoint, &dc.hexConfig)
		swX := float32(sw.x) * dc.windowMultiplier
		swY := float32(sw.y+1) * dc.windowMultiplier

		se := gridPos.getPixelPosOutbounds(southeastPoint, &dc.hexConfig)
		seX := float32(se.x+1) * dc.windowMultiplier
		seY := float32(se.y+1) * dc.windowMultiplier

		s := gridPos.getPixelPosOutbounds(southPoint, &dc.hexConfig)
		sX := float32(s.x) * dc.windowMultiplier
		sY := float32(s.y+1) * dc.windowMultiplier

		vector.StrokeLine(
			dc.screen,
			nwX,
			nwY,
			swX,
			swY,
			lineWidth,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			nX,
			nY,
			nwX,
			nwY,
			lineWidth,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			nX,
			nY,
			neX,
			neY,
			lineWidth,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			neX,
			neY,
			seX,
			seY,
			lineWidth,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			seX,
			seY,
			sX,
			sY,
			lineWidth,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			swX,
			swY,
			sX,
			sY,
			lineWidth,
			dc.highlightColor,
			false)

	}

}
