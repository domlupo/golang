// Basic drawing configuration and utility functions.

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Drawing Configurations for game components which are shown to player
type drawConfig struct {
	screen *ebiten.Image
	grid   grid
	// windowMultiplier must only be used for drawing images and windows.
	// Never use for game logic. float32 rather than idiomatic float64
	// because alot of Ebiten draw APIs use float32
	windowMultiplier              float32
	selectectionBoxLineWidth      int
	missileExplodeRadiusLineWidth int
	healthLineWidth               int
	targetNodeLineWidth           int
	hexLineWidth                  int
	backgroundColor               color.RGBA
	highlightColor                color.RGBA
	unitHealthColor               color.RGBA
	nodeHealthColor               color.RGBA
	bluePlayerNodeColor           color.RGBA
	bluePlayerUnitColor           color.RGBA
	redPlayerNodeColor            color.RGBA
	redPlayerUnitColor            color.RGBA
	neutralNodeColor              color.RGBA
	drawTriangleDefaultImage      *ebiten.Image
}

func newDrawConfig() *drawConfig {
	// TODO: check if this should use a subimage if final game uses draw triangle.
	// Not investigating now since draw triangle may not be used in final game.
	drawTriangleDefaultImage := ebiten.NewImage(1, 1)
	drawTriangleDefaultImage.Fill(color.White)
	highlightColor := color.RGBA{70, 255, 0, 200}
	// For 1280p * 720p resolution, a hexagon height of 40p is ideally 80/sqrt(3) wide.
	// 80/sqrt(3) is 46.188 which is approximately 46p.
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)
	dc := drawConfig{
		screen:                        nil,
		windowMultiplier:              2.5,
		grid:                          grid{&hc},
		selectectionBoxLineWidth:      1,
		missileExplodeRadiusLineWidth: 1,
		healthLineWidth:               2,
		targetNodeLineWidth:           2,
		hexLineWidth:                  2,
		backgroundColor:               color.RGBA{200, 200, 200, 255},
		highlightColor:                highlightColor,
		unitHealthColor:               highlightColor,
		nodeHealthColor:               highlightColor,
		bluePlayerNodeColor:           color.RGBA{0, 0, 255, 100},
		bluePlayerUnitColor:           color.RGBA{0, 0, 255, 100},
		redPlayerNodeColor:            color.RGBA{255, 0, 0, 100},
		redPlayerUnitColor:            color.RGBA{255, 0, 0, 100},
		neutralNodeColor:              color.RGBA{100, 60, 100, 200},
		drawTriangleDefaultImage:      drawTriangleDefaultImage,
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

type drawTriangleColor struct {
	red   float32
	green float32
	blue  float32
	alpha float32
}

func drawHex(dc drawConfig, p GridPosition, color drawTriangleColor, selected bool) {
	center := dc.grid.getPixelPosOutbounds(p, hpCenter)
	north := dc.grid.getPixelPosOutbounds(p, hpNorth)
	northeast := dc.grid.getPixelPosOutbounds(p, hpNortheast)
	southeast := dc.grid.getPixelPosOutbounds(p, hpSoutheast)
	southwest := dc.grid.getPixelPosOutbounds(p, hpSouthwest)
	south := dc.grid.getPixelPosOutbounds(p, hpSouth)
	northwest := dc.grid.getPixelPosOutbounds(p, hpNorthwest)

	// The one pixel offsets for inclusive or exclusive hex points are handled at the game logic layer.
	// The window multiplier offsets are handled here since window multiplier is a draw concept.
	xCenter, yCenter := float32(center.X)*dc.windowMultiplier, float32(center.Y)*dc.windowMultiplier
	xNorth, yNorth := float32(north.X)*dc.windowMultiplier, float32(north.Y)*dc.windowMultiplier
	xNortheast, yNortheast := float32(northeast.X+1)*dc.windowMultiplier, float32(northeast.Y)*dc.windowMultiplier
	xSoutheast, ySoutheast := float32(southeast.X+1)*dc.windowMultiplier, float32(southeast.Y+1)*dc.windowMultiplier
	xSouthwest, ySouthwest := float32(southwest.X)*dc.windowMultiplier, float32(southwest.Y+1)*dc.windowMultiplier
	xSouth, ySouth := float32(south.X)*dc.windowMultiplier, float32(south.Y+1)*dc.windowMultiplier
	xNorthwest, yNorthwest := float32(northwest.X)*dc.windowMultiplier, float32(northwest.Y)*dc.windowMultiplier

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
		lineWidth := float32(dc.hexLineWidth) * dc.windowMultiplier

		north := dc.grid.getPixelPosOutbounds(p, hpNorth)
		nX := float32(north.X) * dc.windowMultiplier
		nY := float32(north.Y) * dc.windowMultiplier

		nw := dc.grid.getPixelPosOutbounds(p, hpNorthwest)
		nwX := float32(nw.X) * dc.windowMultiplier
		nwY := float32(nw.Y) * dc.windowMultiplier

		ne := dc.grid.getPixelPosOutbounds(p, hpNortheast)
		neX := float32(ne.X+1) * dc.windowMultiplier
		neY := float32(ne.Y) * dc.windowMultiplier

		sw := dc.grid.getPixelPosOutbounds(p, hpSouthwest)
		swX := float32(sw.X) * dc.windowMultiplier
		swY := float32(sw.Y+1) * dc.windowMultiplier

		se := dc.grid.getPixelPosOutbounds(p, hpSoutheast)
		seX := float32(se.X+1) * dc.windowMultiplier
		seY := float32(se.Y+1) * dc.windowMultiplier

		s := dc.grid.getPixelPosOutbounds(p, hpSouth)
		sX := float32(s.X) * dc.windowMultiplier
		sY := float32(s.Y+1) * dc.windowMultiplier

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
