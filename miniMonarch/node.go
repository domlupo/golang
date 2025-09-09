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
	player     player
	gridPos    gridPosition
	radius     int
	selected   bool
	health     int
	capturedBy *player
	target     *pixelPosition
	unitConfig unitConfig
}

func newNode(p player, gridPos gridPosition, unit unitConfig) node {
	n := node{
		player:     p,
		gridPos:    gridPos,
		radius:     nodeRadius,
		selected:   false,
		health:     initialNodeHealth,
		capturedBy: nil,
		target:     nil,
		unitConfig: unit,
	}
	return n
}

func (n *node) drawNode(screen *ebiten.Image, nodeGridPos gridPosition) {
	var nodeColor color.Color
	switch n.player {
	case bluePlayer:
		nodeColor = bluePlayerNodeColor
	case redPlayer:
		nodeColor = redPlayerNodeColor
	case neutralPlayer:
		nodeColor = neutralNodeColor
	default:
		panic("unknown player cannot drawNode")
	}

	var nodePixelPos = getPixelPosition(nodeGridPos, innerGridPosition{0, 0})

	vector.DrawFilledCircle(
		screen,
		float32(nodePixelPos.x)*windowMultiplier,
		float32(nodePixelPos.y)*windowMultiplier,
		float32(nodeRadius)*windowMultiplier,
		nodeColor,
		false,
	)

	if n.selected {
		drawUnfilledBoxInt(
			screen,
			highlightColor,
			nodePixelPos.x-n.radius,
			nodePixelPos.y-n.radius,
			nodePixelPos.x+n.radius,
			nodePixelPos.y+n.radius)
	}

	if n.health != initialNodeHealth {
		var healthLineLength = (float32(n.health) / float32(initialNodeHealth)) * float32(nodeRadius*2)
		vector.StrokeLine(
			screen,
			float32(nodePixelPos.x-nodeRadius)*windowMultiplier,
			float32(nodePixelPos.y)*windowMultiplier,
			(float32(nodePixelPos.x-nodeRadius)+healthLineLength)*windowMultiplier,
			float32(nodePixelPos.y)*windowMultiplier,
			float32(healthLineWidth)*windowMultiplier,
			nodeHealthColor,
			false)
	}
}

func (n *node) drawTargets(screen *ebiten.Image, nodeGridPos gridPosition, selectedPlayer player) {
	var nodeColor color.Color
	switch n.player {
	case bluePlayer:
		nodeColor = bluePlayerNodeColor
	case redPlayer:
		nodeColor = redPlayerNodeColor
	case neutralPlayer:
		nodeColor = neutralNodeColor
	default:
		panic("unknown player for drawTargets")
	}

	if n.target != nil && n.player == selectedPlayer {
		var nodePixelPos = getPixelPosition(nodeGridPos, innerGridPosition{0, 0})
		vector.StrokeLine(
			screen,
			float32(nodePixelPos.x)*windowMultiplier,
			float32(nodePixelPos.y)*windowMultiplier,
			float32(n.target.x)*windowMultiplier,
			float32(n.target.y)*windowMultiplier,
			float32(targetNodeLineWidth)*windowMultiplier,
			nodeColor,
			false,
		)
	}
}
