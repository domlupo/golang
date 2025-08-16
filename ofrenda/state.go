package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type beltDirection int

const (
	bNone beltDirection = iota
	bDown
	bDownLeft
	bDownRight
	bLeft
	bLeftDown
	bLeftUp
	bRight
	bRightDown
	bRightUp
	bUp
	bUpLeft
	bUpRight
)

// not enforced, but must be be a singleton
type state struct {
	offers    map[pixelPosition]offer
	belts     map[gridPosition]beltDirection
	consumers map[gridPosition]kind
	producers map[gridPosition]kind
}

func newState() *state {
	s := state{
		offers:    make(map[pixelPosition]offer),
		belts:     make(map[gridPosition]beltDirection),
		consumers: make(map[gridPosition]kind),
		producers: make(map[gridPosition]kind)}
	return &s
}

func (s *state) draw(screen *ebiten.Image) {
	s.drawBelts(screen)
	s.drawProducers(screen)
	s.drawConsumers(screen)
	// offers need to be draw last since they appear above everything else
	s.drawoffers(screen)
}

// TODO: This needs thought out more due to complexity. It is only a simple case.
// The full implementation needs to verify all related offer position are possible
// before moving forward. It also needs to verify no offers are being overwritten
// due to occupying the same pixel position. There is probably also determinism
// concerns as well as other complexities.
func (s *state) moveoffers() {
	var newOffers = make(map[pixelPosition]offer)

	var newPixelPos pixelPosition
	for oldPixelPos, offer := range s.offers {
		var movement pixelPosition
		var beltGridPos = getBeltGridPosFromofferPixelPos(oldPixelPos)
		var beltDirection = s.belts[beltGridPos]
		var offerDirection = getDirectionForShape(beltDirection, oldPixelPos)
		movement = getPixelPositionFromDirection(offerDirection)
		newPixelPos = pixelPosition{oldPixelPos.x + movement.x, oldPixelPos.y + movement.y}
		newOffers[newPixelPos] = offer
	}
	s.offers = newOffers
}

func (s *state) produceOffers() {
	for gridPos, kind := range s.producers {
		var offerPixelPos = getPixelPosition(gridPos, innerGridPosition{13, 5})
		s.offers[offerPixelPos] = offer{kind}
	}
}

func (s *state) consumeoffers() {
	for pixelPos, _ := range s.offers {
		var gridPos, _ = splitPixelPosition(pixelPos)
		if s.consumers[gridPos] != noKind {
			delete(s.offers, pixelPos)
		}
	}
}

func (s *state) addBelt(gridPosition gridPosition, direction beltDirection) {
	s.belts[gridPosition] = direction
}

func (s *state) drawBelts(screen *ebiten.Image) {
	for k, v := range s.belts {
		var xLeft = float32(gridSize * windowMultiplier * k.x)
		var yUp = float32(gridSize * windowMultiplier * k.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch v {
		case bDown:
			image = beltDownImage
		case bDownLeft:
			image = beltDownLeftImage
		case bDownRight:
			image = beltDownRightImage
		case bLeft:
			image = beltLeftImage
		case bLeftDown:
			image = beltLeftDownImage
		case bLeftUp:
			image = beltLeftUpImage
		case bRight:
			image = beltRightImage
		case bRightDown:
			image = beltRightDownImage
		case bRightUp:
			image = beltRightUpImage
		case bUp:
			image = beltUpImage
		case bUpLeft:
			image = beltUpLeftImage
		case bUpRight:
			image = beltUpRightImage
		default:
			panic(fmt.Sprintf("unexpected main.direction: %#v", v))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addConsumer(gridPosition gridPosition, kind kind) {
	s.consumers[gridPosition] = kind
}

func (s *state) drawConsumers(screen *ebiten.Image) {
	for gridPos, kind := range s.consumers {
		var xLeft = float32(gridSize * windowMultiplier * gridPos.x)
		var yUp = float32(gridSize * windowMultiplier * gridPos.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch kind {
		case skull:
			image = consumerSkullImage
		case candy:
			image = consumerCandyImage
		case flower:
			image = consumerFlowerImage
		default:
			panic(fmt.Sprintf("unexpected main.offer: %#v", kind))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addProducer(gridPosition gridPosition, kind kind) {
	s.producers[gridPosition] = kind
}

func (s *state) drawProducers(screen *ebiten.Image) {
	for gridPos, kind := range s.producers {
		var xLeft = float32(gridSize * windowMultiplier * gridPos.x)
		var yUp = float32(gridSize * windowMultiplier * gridPos.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch kind {
		case skull:
			image = producerImage
		case candy:
			image = producerImage
		case flower:
			image = producerImage
		default:
			panic(fmt.Sprintf("unexpected main.offer: %#v", kind))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addOffer(pixelPos pixelPosition, offer offer) {
	s.offers[pixelPos] = offer
}

func (s *state) drawoffers(screen *ebiten.Image) {
	for pixelPosition, offer := range s.offers {
		var xLeft = float64(pixelPosition.x) * windowMultiplier
		var yUp = float64(pixelPosition.y) * windowMultiplier

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(xLeft, float64(yUp))

		var image *ebiten.Image
		switch offer.kind {
		case skull:
			image = skullImage
		case candy:
			image = candyImage
		case flower:
			image = flowerImage
		default:
			panic(fmt.Sprintf("cannot draw offer, unexpected main.offer: %#v", offer))
		}
		screen.DrawImage(image, op)
	}
}
