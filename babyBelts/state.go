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
	shapes    map[pixelPosition]shape
	belts     map[gridPosition]beltDirection
	consumers map[gridPosition]kind
	producers map[gridPosition]kind
	movers    map[gridPosition]mover
}

func newState() *state {
	s := state{
		shapes:    make(map[pixelPosition]shape),
		belts:     make(map[gridPosition]beltDirection),
		consumers: make(map[gridPosition]kind),
		producers: make(map[gridPosition]kind),
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

// TODO: This needs thought out more due to complexity. It is only a simple case.
// The full implementation needs to verify all related shape position are possible
// before moving forward. It also needs to verify no shapes are being overwritten
// due to occupying the same pixel position. There is probably also determinism
// concerns as well as other complexities.
func (s *state) moveShapes() {
	var newShapes = make(map[pixelPosition]shape)

	var newPixelPos pixelPosition
	for oldPixelPos, shape := range s.shapes {
		var movement pixelPosition
		if shape.shapeState == locked {
			movement = pixelPosition{0, 0}
		} else {
			var beltGridPos = getBeltGridPosFromShapePixelPos(oldPixelPos)
			var beltDirection = s.belts[beltGridPos]
			var shapeDirection = getDirectionForShape(beltDirection, oldPixelPos)
			movement = getPixelPositionFromDirection(shapeDirection)
		}
		newPixelPos = pixelPosition{oldPixelPos.x + movement.x, oldPixelPos.y + movement.y}
		newShapes[newPixelPos] = shape
	}
	s.shapes = newShapes
}

func (s *state) moveMovers() {
	for gridPos, m := range s.movers {
		s.movers[gridPos] = m.update()
	}
}

func (s *state) lockShapes() {
	for gridPos, m := range s.movers {
		var shapePixelPos = m.getMoverShapePixelPos(gridPos)
		var shape = s.shapes[shapePixelPos]
		if shape.shapeState == free {
			delete(s.shapes, shapePixelPos)
			s.movers[gridPos] = mover{m.length, m.direction, moving, shape.kind}
		}
	}
}

func (s *state) unlockShapes() {
	for gridPos, m := range s.movers {
		if m.moverState == dropping {
			var shapeKind = m.lockedShapeKind
			var shapePixelPos = m.getMoverShapePixelPos(gridPos)
			s.shapes[shapePixelPos] = shape{shapeKind, free}
			s.movers[gridPos] = mover{m.length, m.direction, waiting, noKind}
		}
	}
}

func (s *state) produceShapes() {
	for gridPos, kind := range s.producers {
		var shapePixelPos = getPixelPosition(gridPos, innerGridPosition{13, 5})
		s.shapes[shapePixelPos] = shape{kind, free}
	}
}

func (s *state) consumeShapes() {
	for pixelPos, _ := range s.shapes {
		var gridPos, _ = splitPixelPosition(pixelPos)
		if s.consumers[gridPos] != noKind {
			delete(s.shapes, pixelPos)
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
		case circle:
			image = consumerCircleImage
		case cross:
			image = consumerCrossImage
		case square:
			image = consumerSquareImage
		default:
			panic(fmt.Sprintf("unexpected main.shape: %#v", kind))
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
		case circle:
			image = producerImage
		case cross:
			image = producerImage
		case square:
			image = producerImage
		default:
			panic(fmt.Sprintf("unexpected main.shape: %#v", kind))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addMover(gridPosition gridPosition, mover mover) {
	s.movers[gridPosition] = mover
}

func (s *state) drawMovers(screen *ebiten.Image) {
	for gridPos, m := range s.movers {
		var xLeft = float32(gridSize * windowMultiplier * gridPos.x)
		var yUp = float32(gridSize * windowMultiplier * gridPos.y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(xLeft), float64(yUp))

		var image *ebiten.Image
		if (m.length == none && m.direction == mNone) || m.moverState == moving {
			image = moverImage
		} else if (m.length == short && m.direction == mDown && m.moverState == waiting) || (m.length == short && m.direction == mUp && m.moverState == dropping) {
			image = moverShortDownImage
		} else if (m.length == long && m.direction == mDown && m.moverState == waiting) || (m.length == long && m.direction == mUp && m.moverState == dropping) {
			image = moverLongDownImage
		} else if (m.length == short && m.direction == mUp && m.moverState == waiting) || (m.length == short && m.direction == mDown && m.moverState == dropping) {
			// up mover needs to be in space above itself, so shift it up one grid
			op.GeoM.Translate(0, -float64(gridSize*windowMultiplier))
			image = moverShortUpImage
		} else if (m.length == long && m.direction == mUp && m.moverState == waiting) || (m.length == long && m.direction == mDown && m.moverState == dropping) {
			// up mover needs to be in space above itself, so shift it up one grid
			op.GeoM.Translate(0, -float64(gridSize*windowMultiplier))
			image = moverLongUpImage
		} else {
			panic(fmt.Sprintf("cannot draw mover, unexpected main.mover: %#v", m))
		}
		screen.DrawImage(image, op)
	}
}

func (s *state) addShape(pixelPos pixelPosition, shape shape) {
	s.shapes[pixelPos] = shape
}

func (s *state) addShapeToMover(gridPos gridPosition, shape shape) {
	var m mover
	if s.movers[gridPos] == m {
		panic(fmt.Sprintf("A belt must exist to add a shape at: %#v", gridPos))
	}
	m = s.movers[gridPos]
	s.addShape(m.getMoverShapePixelPos(gridPos), shape)
}

// TODO: Not finished, prioritized other code before handling all cases for this
/*
func (s *state) addShapeToBeltStart(gridPosition gridPosition, shape shape, directionPriority direction) {
	var beltDirection beltDirection
	if s.belts[gridPosition] == 0 {
		panic(fmt.Sprintf("A belt must exist to add a shape at: %#v", gridPosition))
	}

	beltDirection = s.belts[gridPosition]
	var smallerOffset = (gridSize / 2) - (beltMiddleLineWidth / 2) - shapeSize
	var largerOffset = (gridSize / 2) + (beltMiddleLineWidth / 2)

	if beltDirection == bRight && directionPriority == up {
		s.addShape(gridPosition, innerGridPosition{0, smallerOffset}, shape)
	} else if beltDirection == bRight && directionPriority == down {
		s.addShape(gridPosition, innerGridPosition{0, largerOffset}, shape)
	} else if beltDirection == bUpRight && directionPriority == left {
		s.addShape(gridPosition, innerGridPosition{smallerOffset, gridSize - shapeSize}, shape)
	} else if beltDirection == bUpRight && directionPriority == right {
		s.addShape(gridPosition, innerGridPosition{largerOffset, gridSize - shapeSize}, shape)
	} else {
		panic(fmt.Sprintf("Need to implement for: %#v", beltDirection))
	}
}
*/

func (s *state) drawShapes(screen *ebiten.Image) {
	for pixelPosition, shape := range s.shapes {
		var xLeft = float64(pixelPosition.x) * windowMultiplier
		var yUp = float64(pixelPosition.y) * windowMultiplier

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(xLeft, float64(yUp))

		var image *ebiten.Image
		switch shape.kind {
		case circle:
			image = circleImage
		case cross:
			image = crossImage
		case square:
			image = squareImage
		default:
			panic(fmt.Sprintf("cannot draw shape, unexpected main.shape: %#v", shape))
		}
		screen.DrawImage(image, op)
	}
}
