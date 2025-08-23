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
)

type node struct {
	player   player
	target   *gridPosition
	selected bool
}

func newNode(p player) node {
	n := node{
		player:   p,
		target:   nil,
		selected: false,
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
	default:
		panic("node needs a color")
	}

	var topLeftPixelPosition = getPixelPosition(nodeGridPos, innerGridPosition{0, 0})
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
		drawUnfilledBoxInt(screen, highlightColor, topLeftPixelPosition.x, topLeftPixelPosition.y, bottomRightPixelPosition.x, bottomRightPixelPosition.y)
	}
}

func (n *node) drawTargets(screen *ebiten.Image, nodePos gridPosition, selectedPlayer player) {
	var nodeColor color.Color
	switch n.player {
	case playerOne:
		nodeColor = playerOneNodeColor
	case playerTwo:
		nodeColor = playerTwoNodeColor
	default:
		panic("node needs a color")
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
