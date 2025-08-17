package main

import (
	"fmt"
	"math/rand/v2"
)

type producerShape int

const (
	noProducerShape producerShape = iota
	diamondProducerShape
	squareProducerShape
	plusProducerShape
	crossProducerShape
)

type producer struct {
	offerType offerType
	size      int
	shape     producerShape
}

func newProducer(ot offerType, size int) *producer {
	p := producer{
		offerType: ot,
		size:      size,
		shape:     squareProducerShape,
	}
	return &p
}

func (p *producer) getOfferPosition() gridPosition {
	var gridPos gridPosition

	switch p.shape {
	case squareProducerShape:
		var positions = p.getOfferPositions()
		gridPos = positions[rand.IntN(len(positions))]
	default:
		panic(fmt.Sprintf("unexpected main.producerShape: %#v", p.shape))
	}

	return gridPos
}

func (p *producer) getOfferPositions() []gridPosition {
	var offerPositions []gridPosition

	switch p.shape {
	case squareProducerShape:
		if p.size >= 1 {
			offerPositions = append(offerPositions, gridPosition{0, 1})
			offerPositions = append(offerPositions, gridPosition{1, 1})
			offerPositions = append(offerPositions, gridPosition{1, 0})
			offerPositions = append(offerPositions, gridPosition{1, -1})
			offerPositions = append(offerPositions, gridPosition{0, -1})
			offerPositions = append(offerPositions, gridPosition{-1, -1})
			offerPositions = append(offerPositions, gridPosition{-1, 0})
			offerPositions = append(offerPositions, gridPosition{-1, 1})
		}
		if p.size >= 2 {
			offerPositions = append(offerPositions, gridPosition{0, 2})
			offerPositions = append(offerPositions, gridPosition{1, 2})
			offerPositions = append(offerPositions, gridPosition{2, 2})
			offerPositions = append(offerPositions, gridPosition{2, 1})
			offerPositions = append(offerPositions, gridPosition{2, 0})
			offerPositions = append(offerPositions, gridPosition{2, -1})
			offerPositions = append(offerPositions, gridPosition{2, -2})
			offerPositions = append(offerPositions, gridPosition{1, -2})
			offerPositions = append(offerPositions, gridPosition{0, -2})
			offerPositions = append(offerPositions, gridPosition{-1, -2})
			offerPositions = append(offerPositions, gridPosition{-2, -2})
			offerPositions = append(offerPositions, gridPosition{-2, -1})
			offerPositions = append(offerPositions, gridPosition{-2, 0})
			offerPositions = append(offerPositions, gridPosition{-2, 1})
			offerPositions = append(offerPositions, gridPosition{-2, 2})
			offerPositions = append(offerPositions, gridPosition{-1, 2})
		}
	default:
		panic(fmt.Sprintf("unexpected main.producerShape: %#v", p.shape))
	}

	return offerPositions
}
