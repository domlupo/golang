package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type stock struct {
	cardPile cardPile
}

func newStock(x, y int) *stock {
	s := stock{
		cardPile: *newCardPile(x, y),
	}
	return &s
}

func (f *stock) draw(screen *ebiten.Image) {
	f.cardPile.draw(screen)
}
