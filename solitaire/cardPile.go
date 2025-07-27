package main

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type cardPile struct {
	x, y                   int
	width, height          int
	cards                  []card
	cardCountAffectsHeight bool
	emptyImage             ebiten.Image
}

func newCardPile(x, y int, cardCountAffectsHeight bool) *cardPile {
	cp := cardPile{
		x:                      x,
		y:                      y,
		width:                  cardWidth,
		height:                 cardHeight,
		cards:                  nil,
		cardCountAffectsHeight: cardCountAffectsHeight,
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
		var priorCard = &cp.cards[len(cp.cards)-1]
		if cp.cardCountAffectsHeight {
			priorCard.height = cardNumberHeight
			cp.height = cp.height + cardNumberHeight
		} else {
			c.x = priorCard.x
			c.y = priorCard.y
		}
	}

	cp.cards = append(cp.cards, c)
}

func (cp *cardPile) removeCard() (*card, error) {
	if len(cp.cards) == 0 {
		return nil, errors.New("cannot remove cards from an empty card pile")
	}

	if cp.cardCountAffectsHeight {
		cp.height = max(cp.height-cardNumberHeight, cardHeight)
	} else {
		cp.height = cardHeight
	}

	var removedCard = cp.cards[len(cp.cards)-1]
	cp.cards = cp.cards[:len(cp.cards)-1]

	if len(cp.cards) > 1 {
		var newLastCard = &cp.cards[len(cp.cards)-1]
		newLastCard.height = cardHeight
		newLastCard.width = cardWidth
	}

	return &removedCard, nil
}

func (cp *cardPile) xMid() int {
	return cp.x + (cp.width / 2)
}

func (cp *cardPile) yMid() int {
	return cp.y + (cp.height / 2)
}

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

	var newCardPile = newCardPile(card.x, card.y, true)
	for ci, c := range cp.cards {
		if card.Equal(&c) {
			newCardPile.cards = cp.cards[ci:]
			cp.cards = cp.cards[:ci]
			if len(cp.cards) > 0 {
				cp.cards[len(cp.cards)-1].height = cardHeight
			}
			break
		}
	}

	return newCardPile
}

func (cp *cardPile) move(x, y int) *cardPile {
	cp.x = cp.x + x
	cp.y = cp.y + y
	for i := range cp.cards {
		cp.cards[i].x = cp.cards[i].x + x
		cp.cards[i].y = cp.cards[i].y + y
	}
	return cp
}

func (cp *cardPile) peek() (*card, error) {
	if len(cp.cards) == 0 {
		return nil, errors.New("card pile is empty")
	}

	return &cp.cards[len(cp.cards)-1], nil
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
