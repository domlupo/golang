package main

import (
	"image/color"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type unit struct {
	id                   uuid.UUID
	player               player
	radius               int
	target               *pixelPosition
	vision               int
	reach                int
	alive                bool
	attackedUnitThisTurn bool
	// TODO: determine behavior here. Right now this may actually be direct collision with target node
	// and does not factor in reach. Reach should most likely be used here.
	canReachTargetNode *node
	selected           bool
}

func newUnit(p player, target *pixelPosition) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.Error("newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:                   id,
		player:               p,
		radius:               unitRadius,
		target:               target,
		vision:               40,
		reach:                13,
		alive:                true,
		attackedUnitThisTurn: false,
		canReachTargetNode:   nil,
		selected:             false,
	}

	return u
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
