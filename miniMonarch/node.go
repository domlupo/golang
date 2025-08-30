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
	target     *gridPosition
	selected   bool
	health     int
	capturedBy *player
	radius     int
}

func newNode(p player) node {
	n := node{
		player:     p,
		target:     nil,
		selected:   false,
		health:     initialNodeHealth,
		capturedBy: nil,
		radius:     nodeRadius,
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
			nodePixelPos.x,
			nodePixelPos.y,
			nodePixelPos.x,
			nodePixelPos.y)
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
			highlightColor,
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
		var targetPixelPos = getPixelPosition(*n.target, innerGridPosition{0, 0})
		vector.StrokeLine(
			screen,
			float32(nodePixelPos.x)*windowMultiplier,
			float32(nodePixelPos.y)*windowMultiplier,
			float32(targetPixelPos.x)*windowMultiplier,
			float32(targetPixelPos.y)*windowMultiplier,
			float32(targetNodeLineWidth)*windowMultiplier,
			nodeColor,
			false,
		)
	}
}
