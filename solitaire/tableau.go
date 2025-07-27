package main

type tableau struct {
	cardPile cardPile
}

func newTableau(x, y int) *tableau {
	t := tableau{
		cardPile: *newCardPile(x, y, true),
	}
	return &t
}

func (t *tableau) addCardPileToTableau(cp cardPile) bool {
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

	var card, _ = t.cardPile.peek()
	var symbol = card.symbol
	if (symbol == Club || symbol == Spade) && (cardToAdd.symbol == Club || cardToAdd.symbol == Spade) {
		return false
	}
	if (symbol == Heart || symbol == Diamond) && (cardToAdd.symbol == Heart || cardToAdd.symbol == Diamond) {
		return false
	}

	if t.cardPile.cards[len(t.cardPile.cards)-1].value-1 == cardToAdd.value {
		t.cardPile.addCardPile(cp)

		for ci, c := range t.cardPile.cards[:len(t.cardPile.cards)-1] {
			if !c.hidden {
				t.cardPile.cards[ci].active = true
				t.cardPile.cards[ci].height = cardNumberHeight
			}
		}

		return true
	}

	return false
}
