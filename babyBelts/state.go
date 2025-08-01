package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type beltDirection int

const (
	bDown beltDirection = iota
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

type shape int

const (
	circle shape = iota
	cross
	square
)

// not enforced, but must be be a singleton
type state struct {
	shapes    map[totalPosition]shape
	belts     map[gridPosition]beltDirection
	consumers map[gridPosition]shape
	producers map[gridPosition]shape
	movers    map[gridPosition]mover
}

func newState() *state {
	s := state{
		shapes:    make(map[totalPosition]shape),
		belts:     make(map[gridPosition]beltDirection),
		consumers: make(map[gridPosition]shape),
		producers: make(map[gridPosition]shape),
		movers:    make(map[gridPosition]mover),
	}
	return &s
}

func (s *state) draw(screen *ebiten.Image) {
	s.drawBelts(screen)
	s.drawProducers(screen)
	s.drawConsumers(screen)
	// movers need to be drawn after belts, producers and consumsers since they appear above them
	s.drawMovers(screen)
	// shapes need to be draw last since they appear above everything else
	s.drawShapes(screen)
}

func (s *state) addBelt(gridPosition gridPosition, direction beltDirection) {
	s.belts[gridPosition] = direction
}

func (s *state) drawBelts(screen *ebiten.Image) {
	for k, v := range s.belts {
		var xLeft = float32(gridSize * k.x)
		var yUp = float32(gridSize * k.y)

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

func (s *state) addConsumer(gridPosition gridPosition, shape shape) {
	s.consumers[gridPosition] = shape
}

func (s *state) drawConsumers(screen *ebiten.Image) {
	for k, v := range s.consumers {
		var xLeft = float32(gridSize * k.x)
		var yUp = float32(gridSize * k.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch v {
		case circle:
			image = consumerCircleImage
		case cross:
			image = consumerCrossImage
		case square:
			image = consumerSquareImage
		default:
			panic(fmt.Sprintf("unexpected main.shape: %#v", v))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addProducer(gridPosition gridPosition, shape shape) {
	s.producers[gridPosition] = shape
}

func (s *state) drawProducers(screen *ebiten.Image) {
	for k, v := range s.producers {
		var xLeft = float32(gridSize * k.x)
		var yUp = float32(gridSize * k.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch v {
		case circle:
			image = producerImage
		case cross:
			image = producerImage
		case square:
			image = producerImage
		default:
			panic(fmt.Sprintf("unexpected main.shape: %#v", v))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addMover(gridPosition gridPosition, mover mover) {
	s.movers[gridPosition] = mover
}

func (s *state) drawMovers(screen *ebiten.Image) {
	for k, v := range s.movers {
		var xLeft = float32(gridSize * k.x)
		var yUp = float32(gridSize * k.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		switch v {
		case mover{none, mNone}:
			image = moverImage
		case mover{short, mDown}:
			image = moverShortDownImage
		case mover{short, mUp}:
			// up mover needs to be in space above itself, so shift it up one grid
			op.GeoM.Translate(0, -float64(gridSize))
			image = moverShortUpImage
		case mover{long, mDown}:
			image = moverLongDownImage
		case mover{long, mUp}:
			image = moverLongUpImage
			// up mover needs to be in space above itself, so shift it up one grid
			op.GeoM.Translate(0, -float64(gridSize))
		default:
			panic(fmt.Sprintf("unexpected main.mover: %#v", v))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addShape(gridPosition gridPosition, innerGridPosition innerGridPosition, shape shape) {
	// TODO: verify position. This is very important since no validation means that multiple shapes can
	// exist on a single position since a innerGridPosition > GridSize would appear on the next gridPosition

	s.shapes[totalPosition{gridPosition, innerGridPosition}] = shape
}

func (s *state) drawShapes(screen *ebiten.Image) {
	for k, v := range s.shapes {
		var xLeft = float64(gridSize*k.gridPosition.x) + (float64(k.innerGridPosition.x) * windowMultiplier)
		var yUp = float64(gridSize*k.gridPosition.y) + (float64(k.innerGridPosition.y) * windowMultiplier)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(xLeft, float64(yUp))

		var image *ebiten.Image
		switch v {
		case circle:
			image = circleImage
		case cross:
			image = crossImage
		case square:
			image = squareImage
		default:
			panic(fmt.Sprintf("unexpected main.shape: %#v", v))
		}
		screen.DrawImage(image, op)
	}
}
