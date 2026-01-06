// Responsible for rendering pixels and images to the screen.

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type drawer struct {
	drawConfig
	*state
	*controller
	debugger *debugger
}

func newDrawer(c drawConfig, s *state, controller *controller, d *debugger) *drawer {
	drawer := drawer{
		drawConfig: c,
		state:      s,
		controller: controller,
		debugger:   d,
	}
	return &drawer
}

func (d drawer) draw(screen *ebiten.Image) {
	if d.drawConfig.screen == nil {
		d.screen = screen
	}
	d.screen.Fill(d.backgroundColor)
	d.drawHexGrid()
	for _, n := range d.nodes {
		n.drawNode(d.drawConfig, d.grid)
	}
	for _, o := range d.obstacles {
		o.drawObstacle(d.drawConfig, d.grid)
	}
	for _, u := range d.landUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}
	for _, u := range d.chaffUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}
	for _, u := range d.bombGravityUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}
	for _, u := range d.ordanceUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}
	for _, u := range d.aircraftUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}
	for _, u := range d.ballisticUnits {
		u.draw(d.drawConfig, u.selected, d.grid.pixelMultiplier)
	}

	d.debugger.draw(d.drawConfig)

	if d.leftMouseHeld {
		drawUnfilledBoxFloat32(
			d.drawConfig,
			d.drawConfig.highlightColor,
			d.initialSelectedX,
			d.initialSelectedY,
			d.currentSelectedX,
			d.currentSelectedY,
			d.grid.pixelMultiplier,
		)
	}
}

func (d drawer) drawHexGrid() {
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

	for i := 0; i < d.grid.xIndexMax; i++ {
		for j := 0; j < d.grid.yIndexMax; j++ {
			if j%2 == 0 && i == d.grid.xIndexMax-1 {
				continue
			}

			p := GridPosition{i, j}

			n := d.grid.getPixelOutbounds(p, hpNorth)
			nX := float32(n.X) * d.grid.pixelMultiplier
			nY := float32(n.Y) * d.grid.pixelMultiplier

			nw := d.grid.getPixelOutbounds(p, hpNorthwest)
			nwX := float32(nw.X) * d.grid.pixelMultiplier
			nwY := float32(nw.Y) * d.grid.pixelMultiplier

			ne := d.grid.getPixelOutbounds(p, hpNortheast)
			neX := float32(ne.X+1) * d.grid.pixelMultiplier
			neY := float32(ne.Y) * d.grid.pixelMultiplier

			sw := d.grid.getPixelOutbounds(p, hpSouthwest)
			swX := float32(sw.X) * d.grid.pixelMultiplier
			swY := float32(sw.Y+1) * d.grid.pixelMultiplier

			// draw west side
			//    / \
			//  *|*
			//  *|*
			vector.StrokeLine(
				d.screen,
				nwX,
				nwY,
				swX,
				swY,
				float32(d.hexGridLineWidth)*d.grid.pixelMultiplier,
				d.hexGridColor,
				false)

			// draw northwest
			//   */*\
			//   |
			//   |
			vector.StrokeLine(
				d.screen,
				nX,
				nY,
				nwX,
				nwY,
				float32(d.hexGridLineWidth)*d.grid.pixelMultiplier,
				d.hexGridColor,
				false)

			// draw northeast
			//    /*\*
			//   |
			//   |
			vector.StrokeLine(
				d.screen,
				nX,
				nY,
				neX,
				neY,
				float32(d.hexGridLineWidth)*d.grid.pixelMultiplier,
				d.hexGridColor,
				false)
		}
	}
}

func drawUnfilledBoxFloat32(
	dc drawConfig,
	color color.Color,
	initialX float32,
	initialY float32,
	finalX float32,
	finalY float32,
	pixelMultiplier float32) {

	initialXGraphical := float32(initialX) * pixelMultiplier
	initialYGraphical := float32(initialY) * pixelMultiplier
	finalXGraphical := float32(finalX) * pixelMultiplier
	finalYGraphical := float32(finalY) * pixelMultiplier
	selectectionBoxLineWidth := float32(dc.selectectionBoxLineWidth) * pixelMultiplier

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

func drawUnfilledBoxInt(
	dc drawConfig,
	color color.Color,
	initialX int,
	initialY int,
	finalX int,
	finalY int,
	pixelMultiplier float32) {

	drawUnfilledBoxFloat32(
		dc,
		color,
		float32(initialX),
		float32(initialY),
		float32(finalX),
		float32(finalY),
		pixelMultiplier)
}

func drawHex(dc drawConfig, g grid, p GridPosition, color drawTriangleColor, selected bool) {
	center := g.getPixelOutbounds(p, hpCenter)
	north := g.getPixelOutbounds(p, hpNorth)
	northeast := g.getPixelOutbounds(p, hpNortheast)
	southeast := g.getPixelOutbounds(p, hpSoutheast)
	southwest := g.getPixelOutbounds(p, hpSouthwest)
	south := g.getPixelOutbounds(p, hpSouth)
	northwest := g.getPixelOutbounds(p, hpNorthwest)

	// The one pixel offsets for inclusive or exclusive hex points are handled at the game logic
	// layer. The pixelMultiplier offsets are handled here since pixelMultiplier is a
	// draw concept.
	xCenter := float32(center.X) * g.pixelMultiplier
	yCenter := float32(center.Y) * g.pixelMultiplier
	xNorth := float32(north.X) * g.pixelMultiplier
	yNorth := float32(north.Y) * g.pixelMultiplier
	xNortheast := float32(northeast.X+1) * g.pixelMultiplier
	yNortheast := float32(northeast.Y) * g.pixelMultiplier
	xSoutheast := float32(southeast.X+1) * g.pixelMultiplier
	ySoutheast := float32(southeast.Y+1) * g.pixelMultiplier
	xSouthwest := float32(southwest.X) * g.pixelMultiplier
	ySouthwest := float32(southwest.Y+1) * g.pixelMultiplier
	xSouth := float32(south.X) * g.pixelMultiplier
	ySouth := float32(south.Y+1) * g.pixelMultiplier
	xNorthwest := float32(northwest.X) * g.pixelMultiplier
	yNorthwest := float32(northwest.Y) * g.pixelMultiplier

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5, 0, 5, 6, 0, 6, 1}
	vertices := []ebiten.Vertex{
		{DstX: xCenter, DstY: yCenter, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNorth, DstY: yNorth, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNortheast, DstY: yNortheast, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSoutheast, DstY: ySoutheast, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSouth, DstY: ySouth, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xSouthwest, DstY: ySouthwest, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
		{DstX: xNorthwest, DstY: yNorthwest, SrcX: 0, SrcY: 0,
			ColorR: color.red, ColorG: color.green, ColorB: color.blue, ColorA: color.alpha},
	}

	dc.screen.DrawTriangles(
		vertices, indices, dc.drawTriangleDefaultImage, &ebiten.DrawTrianglesOptions{})

	if selected {
		drawHexHighlight(dc, g, p)
	}
}

func drawHexHighlight(dc drawConfig, g grid, p GridPosition) {
	lineWidth := float32(dc.hexLineWidth) * g.pixelMultiplier

	north := g.getPixelOutbounds(p, hpNorth)
	nX := float32(north.X) * g.pixelMultiplier
	nY := float32(north.Y) * g.pixelMultiplier

	nw := g.getPixelOutbounds(p, hpNorthwest)
	nwX := float32(nw.X) * g.pixelMultiplier
	nwY := float32(nw.Y) * g.pixelMultiplier

	ne := g.getPixelOutbounds(p, hpNortheast)
	neX := float32(ne.X+1) * g.pixelMultiplier
	neY := float32(ne.Y) * g.pixelMultiplier

	sw := g.getPixelOutbounds(p, hpSouthwest)
	swX := float32(sw.X) * g.pixelMultiplier
	swY := float32(sw.Y+1) * g.pixelMultiplier

	se := g.getPixelOutbounds(p, hpSoutheast)
	seX := float32(se.X+1) * g.pixelMultiplier
	seY := float32(se.Y+1) * g.pixelMultiplier

	s := g.getPixelOutbounds(p, hpSouth)
	sX := float32(s.X) * g.pixelMultiplier
	sY := float32(s.Y+1) * g.pixelMultiplier

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
