package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// not enforced but must be a singleton
type state struct {
	nodes map[gridPosition]node
}

func newState() *state {
	s := state{
		nodes: make(map[gridPosition]node),
	}
	return &s
}

func (s *state) draw(screen *ebiten.Image, selectedPlayer player) {
	s.drawNodes(screen)
	s.drawTargets(screen, selectedPlayer)
}

func (s *state) drawNodes(screen *ebiten.Image) {
	for nodePos, n := range s.nodes {
		var nodeColor color.Color
		switch n.player {
		case playerOne:
			nodeColor = playerOneColor
		case playerTwo:
			nodeColor = playerTwoColor
		default:
			panic("node needs a color")
		}

		vector.DrawFilledRect(
			screen,
			float32(nodePos.x*gridSize)*windowMultiplier,
			float32(nodePos.y*gridSize)*windowMultiplier,
			gridSize*windowMultiplier,
			gridSize*windowMultiplier,
			nodeColor,
			false,
		)
	}
}

func (s *state) drawTargets(screen *ebiten.Image, selectedPlayer player) {
	for nodePos, n := range s.nodes {
		var nodeColor color.Color
		switch n.player {
		case playerOne:
			nodeColor = playerOneColor
		case playerTwo:
			nodeColor = playerTwoColor
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
}

func (s *state) addNode(gridPos gridPosition, n node) {
	s.nodes[gridPos] = n
}
