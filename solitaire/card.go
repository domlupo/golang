package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type card struct {
	x, y, unselectedX, unselectedY int
	width, height                  int
	hidden                         bool
	symbol                         symbol
	value                          value
	image                          ebiten.Image
	active                         bool
}

func newCard(value value, symbol symbol, image ebiten.Image) *card {
	c := card{
		x:           0,
		y:           0,
		unselectedX: 0,
		unselectedY: 0,
		width:       cardWidth,
		height:      cardHeight,
		hidden:      true,
		symbol:      symbol,
		value:       value,
		image:       image,
		active:      true,
	}
	return &c
}

func (c *card) containsPoint(x, y int) bool {
	if !c.active {
		return false
	}

	// does not handle corners cut on cards
	return x >= c.x && x <= c.x+c.width && y >= c.y && y <= c.y+c.height
}

func (c *card) Equal(other *card) bool {
	if c == nil || other == nil {
		return c == other
	}
	if c.value != other.value {
		return false
	}
	if c.symbol != other.symbol {
		return false
	}
	return true
}

func (c *card) String() string {
	var value string
	switch c.value {
	case 0:
		value = "ace"
	case 1:
		value = "two"
	case 2:
		value = "three"
	case 3:
		value = "four"
	case 4:
		value = "five"
	case 5:
		value = "six"
	case 6:
		value = "seven"
	case 7:
		value = "eight"
	case 8:
		value = "nine"
	case 9:
		value = "ten"
	case 10:
		value = "jack"
	case 11:
		value = "queen"
	case 12:
		value = "king"
	default:
		value = "failed to interpret value"
	}

	var symbol string
	switch c.symbol {
	case 0:
		symbol = "diamond"
	case 1:
		symbol = "heart"
	case 2:
		symbol = "spade"
	case 3:
		symbol = "club"
	default:
		symbol = "failed to interpret symbol"
	}

	return fmt.Sprintf("%s %s x:%d y:%d", value, symbol, c.x, c.y)
}

func (c *card) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(windowMultiplier, windowMultiplier)
	op.GeoM.Translate(float64(c.x), float64(c.y))
	if c.hidden {
		screen.DrawImage(cardBackImage, op)
	} else {
		screen.DrawImage(&c.image, op)
	}
}
