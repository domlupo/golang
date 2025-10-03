// A game component which spawns units and captured by opposing units
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

type player int

const (
	noPlayer player = iota
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
	var red float32
	var blue float32
	var green float32
	var alpha float32
	switch n.player {
	case bluePlayer:
		red = 0.5
		blue = 1
		green = 0.5
		alpha = 1
	case redPlayer:
		red = 1
		blue = 0.5
		green = 0.5
		alpha = 1
	case neutralPlayer:
		red = 0.5
		blue = 0.5
		green = 0.5
		alpha = 1
	default:
		panic("unknown player cannot drawNode")
	}

	drawHex(dc, n.gridPos, drawTriangleColor{red: red, blue: blue, green: green, alpha: alpha}, n.selected)

	c := getPixelPosOutbounds(n.gridPos, center, &dc.hexConfig)

	if n.currentHealth != n.maxHealth {
		var healthLineLength = (float32(n.currentHealth) / float32(n.maxHealth)) * float32(n.radius*2)
		vector.StrokeLine(
			dc.screen,
			float32(c.x-n.radius)*dc.windowMultiplier,
			float32(c.y)*dc.windowMultiplier,
			(float32(c.x-n.radius)+healthLineLength)*dc.windowMultiplier,
			float32(c.y)*dc.windowMultiplier,
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
		var center, ok = getPixelPos(n.gridPos, center, &dc.hexConfig)
		if !ok {
			panic("could not get pixel pos for node draw targets")
		}
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
