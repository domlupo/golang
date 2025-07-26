package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type tableau struct {
	cardPile cardPile
}

func newTableau(x, y int) *tableau {
	t := tableau{
		cardPile: *newCardPile(x, y, true),
	}
	return &t
}

func (t tableau) addCardPile(cp cardPile) bool {
	if len(cp.cards) == 0 {
		return false
	}

	var cardToAdd = cp.cards[0]

	if len(t.cardPile.cards) == 0 {
		if cardToAdd.value == King {
			t.cardPile.addCardPile(cp)
			return true
		} else {
			return false
		}
	}

	var symbol = t.cardPile.peek().symbol
	if symbol == Club || symbol == Spade && cardToAdd.symbol == Club || cardToAdd.symbol == Spade {
		return false
	}
	if symbol == Heart || symbol == Diamond && cardToAdd.symbol == Heart || cardToAdd.symbol == Diamond {
		return false
	}

	if t.cardPile.cards[len(t.cardPile.cards)-1].value-1 == cardToAdd.value {
		t.cardPile.addCardPile(cp)
		return true
	}

	return false
}

func (f *tableau) draw(screen *ebiten.Image) {
	f.cardPile.draw(screen)
}
