package main

import "errors"

type stock struct {
	cardPile cardPile
}

func newStock(x, y int) *stock {
	s := stock{
		cardPile: *newCardPile(x, y, false),
	}
	return &s
}

func (s *stock) drawCard() (*card, error) {
	if len(s.cardPile.cards) == 0 {
		return nil, errors.New("cannot draw card from empty stock pile")
	}

	return s.cardPile.removeCard()
}
