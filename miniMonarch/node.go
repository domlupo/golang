// A game component which spawns units and captured by opposing units
package main

import (
	"image/color"

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
	capturedBy    *player
	target        position
	unitConfig    unitConfig
}

func newNode(p player, gridPos gridPosition, initialNodeHealth int, hexConfig hexConfig, uc unitConfig) node {
	n := node{
		player:        p,
		gridPos:       gridPos,
		radius:        (hexConfig.hexHeight + hexConfig.hexWidth) / 4,
		selected:      false,
		maxHealth:     initialNodeHealth,
		currentHealth: initialNodeHealth,
		capturedBy:    nil,
		target:        position{},
		unitConfig:    uc,
	}
	return n
}

func (n *node) drawNode(dc drawContext, gridPos gridPosition) {
	var nodeColor color.Color
	switch n.player {
	case bluePlayer:
		nodeColor = dc.bluePlayerNodeColor
	case redPlayer:
		nodeColor = dc.redPlayerNodeColor
	case neutralPlayer:
		nodeColor = dc.neutralNodeColor
	default:
		panic("unknown player cannot drawNode")
	}

	var nodePixelPos, _ = getPixelPosition(gridPos, center, dc.hexConfig)

	vector.DrawFilledCircle(
		dc.screen,
		float32(nodePixelPos.x)*dc.windowMultiplier,
		float32(nodePixelPos.y)*dc.windowMultiplier,
		float32(n.radius)*dc.windowMultiplier,
		nodeColor,
		false,
	)

	if n.selected {
		drawUnfilledBoxInt(
			dc,
			dc.highlightColor,
			nodePixelPos.x-n.radius,
			nodePixelPos.y-n.radius,
			nodePixelPos.x+n.radius,
			nodePixelPos.y+n.radius)
	}

	if n.currentHealth != n.maxHealth {
		var healthLineLength = (float32(n.currentHealth) / float32(n.maxHealth)) * float32(n.radius*2)
		vector.StrokeLine(
			dc.screen,
			float32(nodePixelPos.x-n.radius)*dc.windowMultiplier,
			float32(nodePixelPos.y)*dc.windowMultiplier,
			(float32(nodePixelPos.x-n.radius)+healthLineLength)*dc.windowMultiplier,
			float32(nodePixelPos.y)*dc.windowMultiplier,
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
		var nodePixelPos, _ = getPixelPosition(gridPos, center, dc.hexConfig)
		vector.StrokeLine(
			dc.screen,
			float32(nodePixelPos.x)*dc.windowMultiplier,
			float32(nodePixelPos.y)*dc.windowMultiplier,
			float32(n.target.x)*dc.windowMultiplier,
			float32(n.target.y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			nodeColor,
			false,
		)
	}
}
