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
		player:        p,
		gridPos:       gridPos,
		radius:        hexConfig.hexWidth / 2,
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

func (n *node) drawNode(dc drawContext, gridPos gridPosition) {
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

	center, _ := getPixelPosition(gridPos, center, &dc.hexConfig)
	top, _ := getPixelPosition(gridPos, topOffset, &dc.hexConfig)
	topRight, _ := getPixelPosition(gridPos, topRightOffset, &dc.hexConfig)
	bottomRight, _ := getPixelPosition(gridPos, bottomRightOffset, &dc.hexConfig)
	bottomLeft, _ := getPixelPosition(gridPos, bottomLeftOffset, &dc.hexConfig)
	bottom, _ := getPixelPosition(gridPos, bottomOffset, &dc.hexConfig)
	topLeft, _ := getPixelPosition(gridPos, topLeftOffset, &dc.hexConfig)

	xCenter, yCenter := float32(center.x)*dc.windowMultiplier, float32(center.y)*dc.windowMultiplier
	xTop, yTop := float32(top.x)*dc.windowMultiplier, float32(top.y)*dc.windowMultiplier
	xTopRight, yTopRight := float32(topRight.x)*dc.windowMultiplier, float32(topRight.y)*dc.windowMultiplier
	xBottomRight, yBottomRight := float32(bottomRight.x)*dc.windowMultiplier, float32(bottomRight.y)*dc.windowMultiplier
	xBottom, yBottom := float32(bottom.x)*dc.windowMultiplier, float32(bottom.y)*dc.windowMultiplier
	xBottomLeft, yBottomLeft := float32(bottomLeft.x)*dc.windowMultiplier, float32(bottomLeft.y)*dc.windowMultiplier
	xTopLeft, yTopLeft := float32(topLeft.x)*dc.windowMultiplier, float32(topLeft.y)*dc.windowMultiplier

	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5, 0, 5, 6, 0, 6, 1}
	vertices := []ebiten.Vertex{
		{DstX: xCenter, DstY: yCenter, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xTop, DstY: yTop, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xTopRight, DstY: yTopRight, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xBottomRight, DstY: yBottomRight, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xBottom, DstY: yBottom, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xBottomLeft, DstY: yBottomLeft, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
		{DstX: xTopLeft, DstY: yTopLeft, SrcX: 0, SrcY: 0, ColorR: cR, ColorG: cG, ColorB: cB, ColorA: cA},
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

func (n *node) drawTargets(dc drawContext, gridPos gridPosition, drawBlueTargets bool, drawRedTargets bool) {
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
		var center, _ = getPixelPosition(gridPos, center, &dc.hexConfig)
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
