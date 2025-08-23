package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type player int

const (
	noPlayer = iota
	playerOne
	playerTwo
	playerNeutral
)

type node struct {
	player     player
	target     *gridPosition
	selected   bool
	health     int
	capturedBy *player
}

func newNode(p player) node {
	n := node{
		player:     p,
		target:     nil,
		selected:   false,
		health:     initialNodeHealth,
		capturedBy: nil,
	}
	return n
}

func (n *node) drawNode(screen *ebiten.Image, nodeGridPos gridPosition) {
	var nodeColor color.Color
	switch n.player {
	case playerOne:
		nodeColor = playerOneNodeColor
	case playerTwo:
		nodeColor = playerTwoNodeColor
	case playerNeutral:
		nodeColor = neutralNodeColor
	default:
		panic("unknown player cannot drawNode")
	}

	var topLeftPixelPosition = getPixelPosition(nodeGridPos, innerGridPosition{0, 0})
	var topRightPixelPosition = getPixelPosition(nodeGridPos, innerGridPosition{gridSize, 0})
	var bottomRightPixelPosition = getPixelPosition(nodeGridPos, innerGridPosition{gridSize, gridSize})

	vector.DrawFilledRect(
		screen,
		float32(topLeftPixelPosition.x)*windowMultiplier,
		float32(topLeftPixelPosition.y)*windowMultiplier,
		gridSize*windowMultiplier,
		gridSize*windowMultiplier,
		nodeColor,
		false,
	)

	if n.selected {
		drawUnfilledBoxInt(
			screen,
			highlightColor,
			topLeftPixelPosition.x,
			topLeftPixelPosition.y,
			bottomRightPixelPosition.x,
			bottomRightPixelPosition.y)
	}

	if n.health != initialNodeHealth {
		var healthLineLength = (float32(n.health) / float32(initialNodeHealth)) * gridSize
		vector.StrokeLine(
			screen,
			float32(topLeftPixelPosition.x)*windowMultiplier,
			float32(topLeftPixelPosition.y+(gridSize/2))*windowMultiplier,
			(float32(topLeftPixelPosition.x)+healthLineLength)*windowMultiplier,
			float32(topRightPixelPosition.y+(gridSize/2))*windowMultiplier,
			float32(healthLineWidth)*windowMultiplier,
			highlightColor,
			false)
	}
}

func (n *node) drawTargets(screen *ebiten.Image, nodePos gridPosition, selectedPlayer player) {
	var nodeColor color.Color
	switch n.player {
	case playerOne:
		nodeColor = playerOneNodeColor
	case playerTwo:
		nodeColor = playerTwoNodeColor
	case playerNeutral:
		nodeColor = neutralNodeColor
	default:
		panic("unknown player for drawTargets")
	}

	if n.target != nil && n.player == selectedPlayer {
		vector.StrokeLine(
			screen,
			float32((nodePos.x*gridSize)+(gridSize/2))*windowMultiplier,
			float32((nodePos.y*gridSize)+(gridSize/2))*windowMultiplier,
			float32((n.target.x*gridSize)+(gridSize/2))*windowMultiplier,
			float32((n.target.y*gridSize)+(gridSize/2))*windowMultiplier,
			float32(targetNodeLineWidth)*windowMultiplier,
			nodeColor,
			false,
		)
	}
}
