package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// the utils package has various functions which help with manual debugging and testing

func drawDebugGrid(screen *ebiten.Image, drawDebugGrid bool) {
	if !drawDebugGrid {
		return
	}

	for i := range debugGridSize {
		var x = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(screen, x-debugGridOffset, 0, x-debugGridOffset, float32(windowHeight*windowMultiplier), debugGridWidth, debugGridColor, false)
		vector.StrokeLine(screen, x+debugGridOffset, 0, x+debugGridOffset, float32(windowHeight*windowMultiplier), debugGridWidth, debugGridColor, false)
		var y = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(screen, 0, y-debugGridOffset, float32(windowWidth*windowMultiplier), y-debugGridOffset, debugGridWidth, debugGridColor, false)
		vector.StrokeLine(screen, 0, y+debugGridOffset, float32(windowWidth*windowMultiplier), y+debugGridOffset, debugGridWidth, debugGridColor, false)
	}
}

func (s state) addCCWBelt(topLeftCorner gridPosition) *state {
	s.addBelt(topLeftCorner, bLeftDown)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 1}, bDown)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 2}, bDown)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 3}, bDownRight)
	gameState.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 3}, bRight)
	gameState.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 3}, bRight)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 3}, bRightUp)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 2}, bUp)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 1}, bUp)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y}, bUpLeft)
	gameState.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y}, bLeft)
	gameState.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y}, bLeft)
	return &s
}

func (s state) addCWBelt(topLeftCorner gridPosition) *state {
	s.addBelt(topLeftCorner, bUpRight)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 1}, bUp)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 2}, bUp)
	gameState.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 3}, bLeftUp)
	gameState.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 3}, bLeft)
	gameState.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 3}, bLeft)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 3}, bDownLeft)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 2}, bDown)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 1}, bDown)
	gameState.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y}, bRightDown)
	gameState.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y}, bRight)
	gameState.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y}, bRight)
	return &s
}

func (s state) addMoversAndBelts(topLeftCorner gridPosition) *state {
	// first row movers
	s.addMover(gridPosition{topLeftCorner.x, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 1, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 2, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 3, topLeftCorner.y}, mover{long, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 4, topLeftCorner.y}, mover{long, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 5, topLeftCorner.y}, mover{long, mDown})

	// second row belts
	s.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 1}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 1}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 1}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 1}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 1}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 1}, bRight)

	// third row belts
	s.addBelt(gridPosition{topLeftCorner.x, topLeftCorner.y + 2}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 2}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 2}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 2}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 2}, bRight)
	s.addBelt(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 2}, bRight)

	// fourth row movers
	s.addMover(gridPosition{topLeftCorner.x, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 3}, mover{long, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 3}, mover{long, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 3}, mover{long, mUp})
	return &s
}

func (s state) addMoversAndConsumers(topLeftCorner gridPosition) *state {
	// first row movers
	s.addMover(gridPosition{topLeftCorner.x, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 1, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 2, topLeftCorner.y}, mover{short, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 3, topLeftCorner.y}, mover{long, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 4, topLeftCorner.y}, mover{long, mDown})
	s.addMover(gridPosition{topLeftCorner.x + 5, topLeftCorner.y}, mover{long, mDown})

	// second row consumers
	s.addConsumer(gridPosition{topLeftCorner.x, topLeftCorner.y + 1}, circle)
	s.addConsumer(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 1}, cross)
	s.addConsumer(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 1}, square)
	s.addConsumer(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 1}, circle)
	s.addConsumer(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 1}, cross)
	s.addConsumer(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 1}, square)

	// third row consumers
	s.addConsumer(gridPosition{topLeftCorner.x, topLeftCorner.y + 2}, circle)
	s.addConsumer(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 2}, cross)
	s.addConsumer(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 2}, square)
	s.addConsumer(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 2}, circle)
	s.addConsumer(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 2}, cross)
	s.addConsumer(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 2}, square)

	// fourth row movers
	s.addMover(gridPosition{topLeftCorner.x, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 1, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 2, topLeftCorner.y + 3}, mover{short, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 3, topLeftCorner.y + 3}, mover{long, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 4, topLeftCorner.y + 3}, mover{long, mUp})
	s.addMover(gridPosition{topLeftCorner.x + 5, topLeftCorner.y + 3}, mover{long, mUp})

	return &s
}
