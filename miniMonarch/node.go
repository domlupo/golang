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
	player player
	target *gridPosition
}

func newNode(p player) node {
	n := node{
		player: p,
		target: nil,
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

	vector.DrawFilledRect(
		screen,
		float32(nodeGridPos.x*gridSize)*windowMultiplier,
		float32(nodeGridPos.y*gridSize)*windowMultiplier,
		gridSize*windowMultiplier,
		gridSize*windowMultiplier,
		nodeColor,
		false,
	)
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
