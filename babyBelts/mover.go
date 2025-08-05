package main

import "fmt"

type moverDirection int

const (
	mNone moverDirection = iota
	mDown
	mUp
)

type moverLength int

const (
	none moverLength = iota
	short
	long
)

type moverState int

const (
	waiting moverState = iota
	moving
	dropping
)

type mover struct {
	length          moverLength
	direction       moverDirection
	moverState      moverState
	lockedShapeKind kind
}

func newMover(length moverLength, direction moverDirection) mover {
	m := mover{
		length:    length,
		direction: direction,
	}
	return m
}

func (m mover) update() mover {
	switch m.moverState {
	case waiting:
		return m
	case moving:
		m.moverState = dropping
		return m
	case dropping:
		return m
	}
	panic(fmt.Sprintf("cannot update mover main.mover: %#v", m))
}

func (m *mover) getMoverShapePixelPos(gridPos gridPosition) pixelPosition {
	var x = (gridSize / 2) - (shapeSize / 2)

	var smallOffset = (gridSize / 2) - shapeSize - (beltMiddleLineWidth / 2)
	var largeOffset = (gridSize / 2) + (beltMiddleLineWidth / 2)

	// TODO: The below y inner grid positions are quite confusing, but I could not think of a better way to represent them.
	// Most likely there is no obvious fix because the mover y offsets are arbitrary. I suspect these will change in the future.
	// If true, I will reconsider this issue then.
	if (m.length == short && m.direction == mDown && m.moverState == waiting) || (m.length == short && m.direction == mUp && m.moverState == dropping) {
		return getPixelPosition(gridPosition{gridPos.x, gridPos.y + (beltMiddleLineWidth / 2)}, innerGridPosition{x, smallOffset})
	} else if (m.length == long && m.direction == mDown && m.moverState == waiting) || (m.length == long && m.direction == mUp && m.moverState == dropping) {
		return getPixelPosition(gridPosition{gridPos.x, gridPos.y + (beltMiddleLineWidth / 2)}, innerGridPosition{x, largeOffset})
	} else if (m.length == short && m.direction == mUp && m.moverState == waiting) || (m.length == short && m.direction == mDown && m.moverState == dropping) {
		return getPixelPosition(gridPosition{gridPos.x, gridPos.y - 1}, innerGridPosition{x, largeOffset})
	} else if (m.length == long && m.direction == mUp && m.moverState == waiting) || (m.length == long && m.direction == mDown && m.moverState == dropping) {
		return getPixelPosition(gridPosition{gridPos.x, gridPos.y - (beltMiddleLineWidth / 2)}, innerGridPosition{x, smallOffset})
	}
	panic(fmt.Sprintf("cannot get mover shape pixel position, unexpected main.mover: %#v", m))
}
