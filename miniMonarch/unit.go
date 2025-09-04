package main

import (
	"image/color"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type unitTarget int

const (
	noTarget = iota
	localTarget
	globalTarget
)

type unit struct {
	id                   uuid.UUID
	player               player
	radius               int
	currentTarget        unitTarget
	globalTarget         *pixelPosition
	localTarget          *pixelPosition
	bresenhamError       *int
	nextBresenhamError   *int
	vision               int
	reach                int
	alive                bool
	attackedUnitThisTurn bool
	// TODO: determine behavior here. Right now this may actually be direct collision with target node
	// and does not factor in reach. Reach should most likely be used here.
	canReachTargetNode *node
	selected           bool
}

func newUnit(p player, globalTarget *pixelPosition) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.Error("newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:                   id,
		player:               p,
		radius:               unitRadius,
		currentTarget:        noTarget,
		globalTarget:         globalTarget,
		localTarget:          nil,
		bresenhamError:       nil,
		nextBresenhamError:   nil,
		vision:               40,
		reach:                13,
		alive:                true,
		attackedUnitThisTurn: false,
		canReachTargetNode:   nil,
		selected:             false,
	}

	return u
}

func (u *unit) getNewPos(currentPos pixelPosition) pixelPosition {
	// Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	// Specifically, the integer incremental error form where error is stored over mulitple moves
	var x0 = currentPos.x
	var y0 = currentPos.y
	var x1 = 0
	var y1 = 0
	if u.localTarget != nil {
		u.updateTarget(localTarget)
		x1 = u.localTarget.x
		y1 = u.localTarget.y
	} else if u.globalTarget != nil {
		u.updateTarget(globalTarget)
		x1 = u.globalTarget.x
		y1 = u.globalTarget.y
	} else {
		u.updateTarget(noTarget)
		return currentPos
	}

	var sx int
	var sy int
	var dx = x1 - x0
	var dy = y1 - y0
	if dx < 0 {
		dx = dx * -1
	}
	if dy > 0 {
		dy = dy * -1
	}
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	if u.bresenhamError == nil {
		var error = dx + dy
		var errorCopy = dx + dy
		u.bresenhamError = &error
		u.nextBresenhamError = &errorCopy
	}
	var doubleError = *u.bresenhamError * 2
	if doubleError >= dy && x0 != x1 {
		*u.nextBresenhamError += dy
		x0 += sx
	}
	if doubleError <= dx && y0 != y1 {
		*u.nextBresenhamError += dx
		y0 += sy
	}

	return pixelPosition{x0, y0}
}

func (u *unit) updateTarget(target unitTarget) bool {
	if target != u.currentTarget {
		u.currentTarget = target
		u.bresenhamError = nil
		u.nextBresenhamError = nil
		return true
	}
	return false
}

func (u *unit) completeMove() {
	var intCopy = &u.nextBresenhamError
	u.bresenhamError = *intCopy
}

func (u *unit) draw(screen *ebiten.Image, unitPos pixelPosition, selected bool) {
	var unitColor color.Color
	switch u.player {
	case bluePlayer:
		unitColor = bluePlayerUnitColor
	case redPlayer:
		unitColor = redPlayerUnitColor
	default:
		panic("unit needs a color")
	}

	vector.DrawFilledCircle(
		screen,
		float32(unitPos.x)*windowMultiplier,
		float32(unitPos.y)*windowMultiplier,
		float32(unitRadius)*windowMultiplier,
		unitColor,
		false,
	)

	var outlineColor color.Color = color.Black
	if selected {
		outlineColor = highlightColor
	}
	drawUnfilledBoxInt(
		screen,
		outlineColor,
		unitPos.x-unitRadius,
		unitPos.y-unitRadius,
		unitPos.x+unitRadius,
		unitPos.y+unitRadius,
	)
}
