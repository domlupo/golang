package main

type foundation struct {
	cardPile cardPile
}

func newFoundation(x, y int) *foundation {
	f := foundation{
		cardPile: *newCardPile(x, y, false),
	}
	return &f
}

func (f *foundation) addCardPileToFoundation(cp cardPile) bool {
	if len(cp.cards) != 1 {
		return false
	}

	var cardToAdd = cp.cards[0]

	if len(f.cardPile.cards) == 0 {
		if cardToAdd.value == Ace {
			f.cardPile.addCardPile(cp)
			return true
		} else {
			return false
		}
	}

	if f.cardPile.cards[0].symbol != cardToAdd.symbol {
		return false
	}

	if f.cardPile.cards[len(f.cardPile.cards)-1].value+1 == cardToAdd.value {
		f.cardPile.addCardPile(cp)
		return true
	}

	return false
}
