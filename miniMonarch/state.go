package main

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

// not enforced but must be a singleton
type state struct {
	nodes map[gridPosition]node
	units map[pixelPosition]unit
}

func newState() *state {
	s := state{
		nodes: make(map[gridPosition]node),
		units: make(map[pixelPosition]unit),
	}
	return &s
}

func (s *state) addNode(gridPos gridPosition, n node) {
	s.nodes[gridPos] = n
}

func (s *state) addUnit(nodeGridPos gridPosition, productionNode node) {
	var unitPos = getPixelPosition(nodeGridPos, innerGridPosition{gridSize / 2, gridSize / 2})
	s.units[unitPos] = newUnit(productionNode.player, *productionNode.target)
}

// TODO: Need to implement a move algorithm similar to below.
// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (s *state) moveUnits() {
	var newUnits = make(map[pixelPosition]unit)

	for unitPos1, unit1 := range s.units {
		// check if unit 1 can see unit 2
		var nearestVisibleUnitPos *pixelPosition
		for unitPos2, unit2 := range s.units {
			if unitPos1.x == unitPos2.x && unitPos1.y == unitPos2.y {
				// unit cannot see itself
				continue
			}

			if unit1.player == unit2.player {
				// unit cannot see comarades
				continue
			}

			var distanceBetweenUnits = math.Sqrt(math.Pow(float64(unitPos1.x-unitPos2.x), 2) + math.Pow(float64(unitPos1.y-unitPos2.y), 2))
			if distanceBetweenUnits <= float64(unit1.vision) {
				nearestVisibleUnitPos = &unitPos2
			}
		}

		var xTargetPixelPos = 0
		var yTargetPixelPos = 0

		if nearestVisibleUnitPos != nil {
			xTargetPixelPos = nearestVisibleUnitPos.x
			yTargetPixelPos = nearestVisibleUnitPos.y
		} else {
			xTargetPixelPos = (unit1.target.x * gridSize) + gridSize/2
			yTargetPixelPos = (unit1.target.y * gridSize) + gridSize/2
		}

		var xDelta = 0
		var xMagnitude = 0.0
		var yMagnitude = 0.0

		if unitPos1.x < (xTargetPixelPos) {
			xDelta = 1
			xMagnitude = math.Abs(float64(xTargetPixelPos - unitPos1.x))
		} else if unitPos1.x > (xTargetPixelPos) {
			xDelta = -1
			xMagnitude = math.Abs(float64(unitPos1.x - xTargetPixelPos))
		}

		var yDelta = 0
		if unitPos1.y < (yTargetPixelPos) {
			yDelta = 1
			yMagnitude = math.Abs(float64(yTargetPixelPos - unitPos1.y))
		} else if unitPos1.y > (yTargetPixelPos) {
			yDelta = -1
			yMagnitude = math.Abs(float64(unitPos1.y - yTargetPixelPos))
		}

		if xMagnitude > yMagnitude {
			yDelta = 0
		} else if yMagnitude > xMagnitude {
			xDelta = 0
		} else {
			// TODO: This is indeterministic so may lead to issues
			var coinFlip = rand.Int()
			switch coinFlip {
			case 0:
				xDelta = 0
			case 1:
				yDelta = 0
			}
		}

		if _, alreadyOccupied := newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}]; alreadyOccupied {
			newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
		} else {
			newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}] = unit1
		}
	}

	s.units = newUnits
}

func (s *state) attackUnits() {
	for unitPos1, unit1 := range s.units {
		// check if unit 1 can attack unit 2
		var nearestVisibleUnitPos *pixelPosition = nil
		var nearestVisibleUnit *unit = nil
		for unitPos2, unit2 := range s.units {
			if unitPos1.x == unitPos2.x && unitPos1.y == unitPos2.y {
				// unit cannot attack itself
				continue
			}

			if unit1.player == unit2.player {
				// unit cannot attack comarades
				continue
			}

			var distanceBetweenUnits = math.Sqrt(math.Pow(float64(unitPos1.x-unitPos2.x), 2) + math.Pow(float64(unitPos1.y-unitPos2.y), 2))
			if distanceBetweenUnits <= float64(unit1.vision) {
				nearestVisibleUnitPos = &unitPos2
				nearestVisibleUnit = &unit2
			}
		}

		if nearestVisibleUnitPos == nil {
			continue
		}

		var xTargetPixelPos = nearestVisibleUnitPos.x
		var yTargetPixelPos = nearestVisibleUnitPos.y

		var xMagnitude = 0.0
		var yMagnitude = 0.0
		if unitPos1.x < (xTargetPixelPos) {
			xMagnitude = math.Abs(float64(xTargetPixelPos - unitPos1.x))
		} else if unitPos1.x > (xTargetPixelPos) {
			xMagnitude = math.Abs(float64(unitPos1.x - xTargetPixelPos))
		}
		if unitPos1.y < (yTargetPixelPos) {
			yMagnitude = math.Abs(float64(yTargetPixelPos - unitPos1.y))
		} else if unitPos1.y > (yTargetPixelPos) {
			yMagnitude = math.Abs(float64(unitPos1.y - yTargetPixelPos))
		}

		if float64(unit1.reach) >= math.Sqrt(math.Pow(xMagnitude, 2)+math.Pow(yMagnitude, 2)) {
			nearestVisibleUnit.alive = false
			s.units[*nearestVisibleUnitPos] = *nearestVisibleUnit
		}
	}

	for attackedUnitPos, attackedUnit := range s.units {
		if !attackedUnit.alive {
			delete(s.units, attackedUnitPos)
		}
	}
}

func (s *state) produceUnits() {
	for nodePos, n := range s.nodes {
		gameState.addUnit(nodePos, n)
	}
}

func (s *state) draw(screen *ebiten.Image, selectedPlayer player) {
	for nodePos, node := range s.nodes {
		node.drawNode(screen, nodePos)
	}
	for nodePos, node := range s.nodes {
		node.drawTargets(screen, nodePos, selectedPlayer)
	}
	for unitPos, unit := range s.units {
		unit.draw(screen, unitPos)
	}
}
