package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// This solitaire code was ported. It was primarily tested manually rather than unit tests.

func TestTableauToFoundation(t *testing.T) {
	// setup
	var tableau = newTableau(0, 0)
	var card = newCard(Ace, Heart, ebiten.Image{})
	tableau.cardPile.addCard(*card)
	var foundation = newFoundation(0, 0)
	if len(tableau.cardPile.cards) != 1 {
		t.Errorf("test case needs tableau to start with one card")
	}
	if len(foundation.cardPile.cards) != 0 {
		t.Errorf("test case needs foundation to start with zero cards")
	}

	// execute
	var cardPile = tableau.cardPile.splitCardPile(*card)
	foundation.addCardPileToFoundation(*cardPile)

	// assert
	if len(tableau.cardPile.cards) != 0 {
		t.Errorf("expected tableau to have zero cards")
	}
	if len(foundation.cardPile.cards) != 1 {
		t.Errorf("expected foundation to have one card")
	}
}

func TestRemoveCard(t *testing.T) {
	// setup
	var stock = newStock(0, 0)
	var card = newCard(Ace, Heart, ebiten.Image{})
	stock.cardPile.addCard(*card)
	if len(stock.cardPile.cards) != 1 {
		t.Errorf("test case needs stock to start with one card")
	}

	// execute
	stock.cardPile.removeCard()

	// assert
	if len(stock.cardPile.cards) != 0 {
		t.Errorf("expected stock to have zero cards")
	}
}

func TestFoundationAddCardPile(t *testing.T) {
	// setup
	var foundation = newFoundation(0, 0)
	if len(foundation.cardPile.cards) != 0 {
		t.Errorf("test case needs foundation to start with zero cards")
	}

	var AceOfHeartsCard = newCard(Ace, Heart, ebiten.Image{})
	var AceOfHeartsCardPile = newCardPile(0, 0, true)
	AceOfHeartsCardPile.addCard(*AceOfHeartsCard)

	var TwoOfHeartsCard = newCard(Two, Heart, ebiten.Image{})
	var TwoOfHeartsCardPile = newCardPile(0, 0, true)
	TwoOfHeartsCardPile.addCard(*TwoOfHeartsCard)

	var ThreeOfHeartsCard = newCard(Three, Heart, ebiten.Image{})
	var ThreeOfHeartsCardPile = newCardPile(0, 0, true)
	ThreeOfHeartsCardPile.addCard(*ThreeOfHeartsCard)

	// execute
	foundation.addCardPileToFoundation(*TwoOfHeartsCardPile)
	foundation.addCardPileToFoundation(*ThreeOfHeartsCardPile)
	if len(foundation.cardPile.cards) != 0 {
		t.Errorf("expected foundation to not add two or three of hearts when foundation is empty")
	}

	foundation.addCardPileToFoundation(*AceOfHeartsCardPile)
	if len(foundation.cardPile.cards) != 1 {
		t.Errorf("expected foundation to add ace of hearts when foundation is empty")
	}

	foundation.addCardPileToFoundation(*ThreeOfHeartsCardPile)
	if len(foundation.cardPile.cards) != 1 {
		t.Errorf("expected foundation to not add three of hearts when foundation only has ace of hearts")
	}

	foundation.addCardPileToFoundation(*TwoOfHeartsCardPile)
	if len(foundation.cardPile.cards) != 2 {
		t.Errorf("expected foundation to add two of hearts when foundation is has ace of hearts")
	}

	foundation.addCardPileToFoundation(*ThreeOfHeartsCardPile)
	if len(foundation.cardPile.cards) != 3 {
		t.Errorf("expected foundation to add three of hearts when foundation has two of hearts")
	}
}

func TestTableauAddCardPile(t *testing.T) {
	// setup
	var twoOfHeartsCard = newCard(Two, Heart, ebiten.Image{})
	var tableau = newTableau(0, 0)
	tableau.cardPile.addCard(*twoOfHeartsCard)

	if len(tableau.cardPile.cards) != 1 {
		t.Errorf("test case needs tableau to start with only two of hearts")
	}

	// execute
	var aceOfSpadesCard = newCard(Ace, Spade, ebiten.Image{})
	var aceOfSpadesCardPile = newCardPile(0, 0, true)
	aceOfSpadesCardPile.addCard(*aceOfSpadesCard)
	tableau.addCardPileToTableau(*aceOfSpadesCardPile)

	// assert
	if len(tableau.cardPile.cards) != 2 {
		t.Errorf("expected tableau to add ace of spades when only two of hearts in tableau")
	}
}
