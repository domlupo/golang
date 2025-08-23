package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type unit struct {
	player               player
	target               gridPosition
	vision               int
	reach                int
	alive                bool
	attackedUnitThisTurn bool
	// TODO: remove this nextToTargetNode HACK
	nextToTargetNode bool
}

func newUnit(p player, target gridPosition) unit {
	u := unit{
		player:               p,
		target:               target,
		vision:               40,
		reach:                1,
		alive:                true,
		attackedUnitThisTurn: false,
		nextToTargetNode:     false,
	}

	return u
}

func (u *unit) draw(screen *ebiten.Image, unitPos pixelPosition) {
	var unitColor color.Color
	switch u.player {
	case playerOne:
		unitColor = playerOneUnitColor
	case playerTwo:
		unitColor = playerTwoUnitColor
	default:
		panic("node needs a color")
	}

	vector.DrawFilledRect(
		screen,
		float32(unitPos.x-(unitWidth/2)-1)*windowMultiplier,
		float32(unitPos.y-(unitHeight/2)-1)*windowMultiplier,
		float32(unitWidth)*windowMultiplier,
		float32(unitHeight)*windowMultiplier,
		unitColor,
		false,
	)
}
