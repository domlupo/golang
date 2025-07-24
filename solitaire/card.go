package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type card struct {
	x, y, unselectedX, unselectedY int
	width, height                  int
	hidden                         bool
	symbol                         symbol
	value                          value
	image                          ebiten.Image
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
	}
	return &c
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

func (c *card) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	if c.hidden {
		screen.DrawImage(cardBackImage, op)
	} else {
		screen.DrawImage(&c.image, op)
	}
}
