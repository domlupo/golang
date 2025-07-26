package main

type waste struct {
	cardPile cardPile
}

func newWaste(x, y int) *waste {
	w := waste{
		cardPile: *newCardPile(x, y, false),
	}
	return &w
}
