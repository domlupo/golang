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

func newCard(symbol symbol, value value, image ebiten.Image) *card {
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
