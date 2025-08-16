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
)

type producer struct {
	offerType offerType
	size      int
	shape     producerShape
}

func newProducer(ot offerType) *producer {
	p := producer{
		offerType: ot,
		size:      2,
		shape:     squareProducerShape,
	}
	return &p
}

func (p *producer) getOfferPosition() gridPosition {
	var xGridPos int = 0
	var yGridPos int = 0

	for xGridPos == 0 && yGridPos == 0 {
		switch p.shape {
		case squareProducerShape:
			xGridPos = rand.IntN(p.size) - (p.size / 2)
			yGridPos = rand.IntN(p.size) - (p.size / 2)
		default:
			panic(fmt.Sprintf("unexpected main.producerShape: %#v", p.shape))
		}
	}

	return gridPosition{xGridPos, yGridPos}
}
