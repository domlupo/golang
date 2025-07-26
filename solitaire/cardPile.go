package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type cardPile struct {
	x, y                   int
	width, height          int
	cards                  []card
	cardCountAffectsHeight bool
	selected               bool
	emptyImage             ebiten.Image
}

func newCardPile(x, y int) *cardPile {
	cp := cardPile{
		x:                      x,
		y:                      y,
		width:                  cardWidth,
		height:                 cardHeight,
		cards:                  nil,
		cardCountAffectsHeight: false,
		selected:               false,
		emptyImage:             *cardPlaceholder,
	}
	return &cp
}

func (cp *cardPile) addCard(c card) {
	var newY int

	if len(cp.cards) == 0 {
		newY = cp.y
	} else {
		newY = cp.cards[len(cp.cards)-1].y + cardNumberHeight
	}

	c.x, c.unselectedX = cp.x, cp.x
	c.y, c.unselectedY = newY, newY

	if len(cp.cards) > 0 {
		var priorCard = cp.cards[len(cp.cards)-1]
		if cp.cardCountAffectsHeight {
			priorCard.height = cardNumberHeight
			cp.height = cp.height + cardNumberHeight
		} else {
			priorCard.height = 0
			priorCard.width = 0
		}
	}

	cp.cards = append(cp.cards, c)
}

func (cp *cardPile) removeCard() card {
	if len(cp.cards) == 0 {
		panic("cannot remove cards from an empty card pile")
	}

	if cp.cardCountAffectsHeight {
		cp.height = max(cp.height-cardNumberHeight, cardHeight)
	} else {
		cp.height = cardHeight
	}

	var removedCard = cp.cards[len(cp.cards)-1]
	cp.cards = cp.cards[:len(cp.cards)-1]

	if len(cp.cards) > 1 {
		var newLastCard = cp.cards[len(cp.cards)-1]
		newLastCard.height = cardHeight
		newLastCard.width = cardWidth
	}

	return removedCard
}

func (cp *cardPile) xMid() int {
	return (cp.x + cp.width) / 2
}

func (cp *cardPile) yMid() int {
	return (cp.y + cp.height) / 2
}

func (cp *cardPile) hasCard(card card) bool {
	for _, c := range cp.cards {
		if card.Equal(&c) {
			return true
		}
	}

	return false
}

// TODO: delete this if unused
/*
func (cp *cardPile) hasCardPile(other cardPile) bool {
	if len(cp.cards) == 0 {
		return false
	}

	if len(other.cards) == 0 {
		return false
	}

	for _, c := range other.cards {
		if cp.hasCard(c) == false {
			return false
		}
	}

	return true
}
*/

func (cp *cardPile) addCardPile(other cardPile) {
	for _, c := range other.cards {
		cp.addCard(c)
	}
}

func (cp *cardPile) containsPoint(x, y int) bool {
	// does not handle corners on cards
	return x >= cp.x && x <= cp.x+cp.width && y >= cp.y && y <= cp.y+cp.height
}

func (cp *cardPile) splitCardPile(card card) *cardPile {
	if len(cp.cards) == 0 {
		panic("cannot split an empty card pile")
	}

	var newCardPile = newCardPile(card.x, card.y)
	for ci, c := range cp.cards {
		if card.Equal(&c) {
			newCardPile.cards = cp.cards[ci:]
			cp.cards = cp.cards[:ci]
			break
		}
	}

	return newCardPile
}

// remove card pile
func (cp *cardPile) removeCardPile(other cardPile) {
	if cp == nil {
		return
	}
}

// move
func (cp *cardPile) move(x, y int) {
	cp.x = cp.x + x
	cp.y = cp.y + y
	for _, c := range cp.cards {
		c.x = c.x + x
		c.y = c.y + y
	}
}

// peek
func (cp *cardPile) peek() card {
	return cp.cards[len(cp.cards)-1]
}

func (cp *cardPile) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(windowMultiplier, windowMultiplier)
	op.GeoM.Translate(float64(cp.x), float64(cp.y))
	screen.DrawImage(&cp.emptyImage, op)
	for _, c := range cp.cards {
		c.draw(screen)
	}
}
