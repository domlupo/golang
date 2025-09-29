// A game component which spawns units and captured by opposing units
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type player int

const (
	noPlayer = iota
	bluePlayer
	redPlayer
	neutralPlayer
)

type node struct {
	player        player
	gridPos       gridPosition
	radius        int
	selected      bool
	maxHealth     int
	currentHealth int
	capturedBy    player
	target        position
	unitConfig    unitConfig
}

func newNode(p player, gridPos gridPosition, initialNodeHealth int, hexConfig *hexConfig, uc unitConfig) node {
	n := node{
		player:  p,
		gridPos: gridPos,
		// HACK: "/ 9 * 2 is temporarily set until hexagon collision is implemented
		// This value was chosen since it is large but still fully within hexagon.
		radius:        (hexConfig.hexHeight + hexConfig.hexWidth) / 9 * 2,
		selected:      false,
		maxHealth:     initialNodeHealth,
		currentHealth: initialNodeHealth,
		capturedBy:    noPlayer,
		target:        position{},
		unitConfig:    uc,
	}
	return n
}

func (n node) isZero() bool {
	// TODO: fix implementation
	return n.player == noPlayer && n.gridPos.x == 0 && n.gridPos.y == 0 && n.target.layer == noLayer
}

func (n *node) drawNode(dc drawContext) {
	var cR float32
	var cB float32
	var cG float32
	var cA float32
	switch n.player {
	case bluePlayer:
		cR = 0.5
		cB = 1
		cG = 0.5
		cA = 1
	case redPlayer:
		cR = 1
		cB = 0.5
		cG = 0.5
		cA = 1
	case neutralPlayer:
		cR = 0.5
		cB = 0.5
		cG = 0.5
		cA = 1
	default:
		panic("unknown player cannot drawNode")
	}

	center := getPixelPosOutbounds(n.gridPos, center, &dc.hexConfig)
	north := getPixelPosOutbounds(n.gridPos, northPoint, &dc.hexConfig)
	northeast := getPixelPosOutbounds(n.gridPos, northeastPoint, &dc.hexConfig)
	southeast := getPixelPosOutbounds(n.gridPos, southeastPoint, &dc.hexConfig)
	southwest := getPixelPosOutbounds(n.gridPos, southwestPoint, &dc.hexConfig)
	south := getPixelPosOutbounds(n.gridPos, southPoint, &dc.hexConfig)
	northwest := getPixelPosOutbounds(n.gridPos, northwestPoint, &dc.hexConfig)

	xCenter, yCenter := float32(center.x+dc.xOffset)*dc.windowMultiplier, float32(center.y+dc.yOffset)*dc.windowMultiplier
	xNorth, yNorth := float32(north.x+dc.xOffset)*dc.windowMultiplier, float32(north.y+dc.yOffset)*dc.windowMultiplier
	xNortheast, yNortheast := float32(northeast.x+dc.xOffset+1)*dc.windowMultiplier-1, float32(northeast.y+dc.yOffset)*dc.windowMultiplier
	xSoutheast, ySoutheast := float32(southeast.x+dc.xOffset+1)*dc.windowMultiplier-1, float32(southeast.y+dc.yOffset+1)*dc.windowMultiplier-1
	xSouthwest, ySouthwest := float32(southwest.x+dc.xOffset)*dc.windowMultiplier, float32(southwest.y+dc.yOffset+1)*dc.windowMultiplier-1
	xSouth, ySouth := float32(south.x+dc.xOffset)*dc.windowMultiplier, float32(south.y+dc.yOffset+1)*dc.windowMultiplier-1
	xNorthwest, yNorthwest := float32(northwest.x+dc.xOffset)*dc.windowMultiplier, float32(northwest.y+dc.yOffset)*dc.windowMultiplier

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5, 0, 5, 6, 0, 6, 1}
	vertices := []ebiten.Vertex{
		{DstX: xCenter, DstY: yCenter, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xNorth, DstY: yNorth, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xNortheast, DstY: yNortheast, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xSoutheast, DstY: ySoutheast, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xSouth, DstY: ySouth, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xSouthwest, DstY: ySouthwest, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xNorthwest, DstY: yNorthwest, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
	}

	dc.screen.DrawTriangles(vertices, indices, dc.drawTriangleDefaultImage, &ebiten.DrawTrianglesOptions{})

	if n.selected {
		drawUnfilledBoxInt(
			dc,
			dc.highlightColor,
			center.x-n.radius,
			center.y-n.radius,
			center.x+n.radius,
			center.y+n.radius)
	}

	if n.currentHealth != n.maxHealth {
		var healthLineLength = (float32(n.currentHealth) / float32(n.maxHealth)) * float32(n.radius*2)
		vector.StrokeLine(
			dc.screen,
			float32(center.x-n.radius)*dc.windowMultiplier,
			float32(center.y)*dc.windowMultiplier,
			(float32(center.x-n.radius)+healthLineLength)*dc.windowMultiplier,
			float32(center.y)*dc.windowMultiplier,
			float32(dc.healthLineWidth)*dc.windowMultiplier,
			dc.nodeHealthColor,
			false)
	}
}

func (n *node) drawTargets(dc drawContext, drawBlueTargets bool, drawRedTargets bool) {
	var nodeColor color.Color
	switch n.player {
	case bluePlayer:
		nodeColor = dc.bluePlayerNodeColor
	case redPlayer:
		nodeColor = dc.redPlayerNodeColor
	case neutralPlayer:
		nodeColor = dc.neutralNodeColor
	default:
		panic("unknown player for drawTargets")
	}

	if !n.target.isZero() && ((n.player == bluePlayer && drawBlueTargets) || (n.player == redPlayer && drawRedTargets)) {
		var center, _ = getPixelPos(n.gridPos, center, &dc.hexConfig)
		vector.StrokeLine(
			dc.screen,
			float32(center.x)*dc.windowMultiplier,
			float32(center.y)*dc.windowMultiplier,
			float32(n.target.x)*dc.windowMultiplier,
			float32(n.target.y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			nodeColor,
			false,
		)
	}
}
