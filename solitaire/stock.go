package main

type stock struct {
	cardPile cardPile
}

func newStock(x, y int) *stock {
	s := stock{
		cardPile: *newCardPile(x, y, false),
	}
	return &s
}
