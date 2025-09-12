package main

import (
	"fmt"
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

type unitConfig struct {
	radius    int
	maxHealth int
	attack    int
	vision    int
	reach     int
}

type unit struct {
	id                 uuid.UUID
	player             player
	currentTarget      unitTarget
	globalTarget       *pixelPosition
	localTarget        *pixelPosition
	bresenhamError     *int
	nextBresenhamError *int
	unitConfig
	health           int
	movedLastTurn    bool
	movedThisTurn    bool
	attackedLastTurn bool
	attackedThisTurn bool
	fluxTurnCount    int
	// TODO: determine behavior here. Right now this may actually be direct collision with target node
	// and does not factor in reach. Reach should most likely be used here.
	canReachTargetNode *node
	selected           bool
}

func newUnit(p player, globalTarget *pixelPosition, config unitConfig) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.Error("newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:                 id,
		player:             p,
		currentTarget:      noTarget,
		globalTarget:       globalTarget,
		localTarget:        nil,
		bresenhamError:     nil,
		nextBresenhamError: nil,
		unitConfig:         config,
		health:             config.maxHealth,
		movedLastTurn:      false,
		movedThisTurn:      false,
		attackedLastTurn:   false,
		attackedThisTurn:   false,
		fluxTurnCount:      0,
		canReachTargetNode: nil,
		selected:           false,
	}

	return u
}

func (u *unit) refresh() {
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:refresh:before: %+v\n", u))

	if !u.attackedLastTurn && !u.movedLastTurn {
		u.fluxTurnCount += 1
	} else {
		u.fluxTurnCount = 0
	}

	u.movedLastTurn = u.movedThisTurn
	u.movedThisTurn = false
	u.attackedLastTurn = u.attackedThisTurn
	u.attackedThisTurn = false
	u.localTarget = nil
	u.canReachTargetNode = nil

	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:refresh:after: %+v", u))
}

func (u *unit) getNewPos(currentPos pixelPosition) pixelPosition {
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:getNewPos:currentPos %+v", u))

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
		logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "unit:getNewPos:localTarget")
	} else if u.globalTarget != nil {
		u.updateTarget(globalTarget)
		x1 = u.globalTarget.x
		y1 = u.globalTarget.y
		logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "unit:getNewPos:globalTarget")
	} else {
		u.updateTarget(noTarget)
		logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "unit:getNewPos:noTarget")
		var desiredPos = currentPos
		logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:getNewPos:desiredPos: %+v", desiredPos))
		return desiredPos
	}
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:getNewPos:targetPos: {%d, %d}", x1, y1))

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

	var desiredPos = pixelPosition{x0, y0}
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, fmt.Sprintf("unit:getNewPos:desiredPos %+v", desiredPos))
	return desiredPos
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
	u.movedThisTurn = true
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
		float32(u.radius)*windowMultiplier,
		unitColor,
		false,
	)

	if selected {
		drawUnfilledBoxInt(
			screen,
			highlightColor,
			unitPos.x-u.radius,
			unitPos.y-u.radius,
			unitPos.x+u.radius,
			unitPos.y+u.radius,
		)
	}

	if drawUnitHealth && u.health != u.maxHealth {
		var healthLineLength = (float32(u.health) / float32(u.maxHealth)) * float32(u.radius*2)

		vector.DrawFilledRect(
			screen,
			float32(unitPos.x-u.radius)*windowMultiplier,
			float32(unitPos.y-u.radius)*windowMultiplier,
			healthLineLength*windowMultiplier,
			float32(u.radius)*windowMultiplier,
			unitHealthColor,
			false,
		)
	}
}
