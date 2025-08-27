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
	width                int
	height               int
	target               *gridPosition
	vision               int
	reach                int
	alive                bool
	attackedUnitThisTurn bool
	// TODO: remove this nextToTargetNode HACK
	nextToTargetNode bool
	selected         bool
}

func newUnit(p player, target *gridPosition) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.Error("failed to generate UUID: %v", err)
	}
	u := unit{
		id:                   id,
		player:               p,
		width:                unitWidth,
		height:               unitHeight,
		target:               target,
		vision:               40,
		reach:                8,
		alive:                true,
		attackedUnitThisTurn: false,
		nextToTargetNode:     false,
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

	var x1 = unitPos.x - (unitWidth / 2) - 1
	var y1 = unitPos.y - (unitWidth / 2) - 1

	vector.DrawFilledRect(
		screen,
		float32(x1)*windowMultiplier,
		float32(y1)*windowMultiplier,
		float32(unitWidth)*windowMultiplier,
		float32(unitHeight)*windowMultiplier,
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
		x1,
		y1,
		x1+unitWidth,
		y1+unitHeight,
	)
}
