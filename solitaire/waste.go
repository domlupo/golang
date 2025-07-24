package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type waste struct {
	cardPile cardPile
}

func newWaste(x, y int) *waste {
	w := waste{
		cardPile: *newCardPile(x, y),
	}
	return &w
}

func (f *waste) draw(screen *ebiten.Image) {
	f.cardPile.draw(screen)
}
