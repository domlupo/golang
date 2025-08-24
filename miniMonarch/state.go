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
	// TODO: determine game design behavior when target is not set
	if productionNode.target == nil {
		return
	}

	var unitPos = getPixelPosition(nodeGridPos, innerGridPosition{gridSize / 2, gridSize / 2})
	s.units[unitPos] = newUnit(productionNode.player, *productionNode.target)
}

func (s *state) refreshUnits() {
	for unitPos, u := range s.units {
		u.attackedUnitThisTurn = false
		s.units[unitPos] = u
	}
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

		_, alreadyOccupiedByUnit := newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}]
		var gridPos, _ = splitPixelPosition(pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta})
		// TODO: Determine behavior wanted when units travel through different player nodes
		var alreadyOccupiedByOpposingNode = true
		if s.nodes[gridPos].player == noPlayer || s.nodes[gridPos].player == unit1.player {
			alreadyOccupiedByOpposingNode = false
		}

		if alreadyOccupiedByUnit {
			newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
		} else if alreadyOccupiedByOpposingNode {
			unit1.nextToTargetNode = true
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
			// TODO: setting this <= rather than < caused obvious bugs but it seems more correct
			// to use <. Need to determine right behavior in the future.
			if distanceBetweenUnits < float64(unit1.vision) {
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
			unit1.attackedUnitThisTurn = true
			s.units[unitPos1] = unit1
		}
	}

	for attackedUnitPos, attackedUnit := range s.units {
		if !attackedUnit.alive {
			delete(s.units, attackedUnitPos)
		}
	}
}

func (s *state) attackNodes() {
	for _, unit := range s.units {
		if unit.attackedUnitThisTurn {
			continue
		}

		// TODO: Fix this hack.
		if unit.nextToTargetNode && s.nodes[unit.target].player != unit.player {
			var targetNode = s.nodes[unit.target]
			targetNode.health = max(targetNode.health-1, 0)
			if targetNode.health == 0 {
				targetNode.capturedBy = &unit.player
			}
			s.nodes[unit.target] = targetNode
		}
	}
}

func (s *state) updateNodeOwnership() {
	for nodeGridPos, node := range s.nodes {
		if node.health == 0 {
			node.health = initialNodeHealth
			node.player = *node.capturedBy
			node.capturedBy = nil
			s.nodes[nodeGridPos] = node
		}
	}
}

func (s *state) updateSelected(clickPosition pixelPosition, releasePosition pixelPosition) {
	// TODO: generalize implementation, for prototype this assumes bottom right release position.
	// The future implementation needs to handle all four cases for release position.

	var topLeftGridPos, _ = splitPixelPosition(clickPosition)
	var bottomRightGridPos, _ = splitPixelPosition(releasePosition)

	for nodeGridPos, node := range s.nodes {
		if (topLeftGridPos.x < nodeGridPos.x) &&
			(nodeGridPos.x < bottomRightGridPos.x) &&
			(topLeftGridPos.y < nodeGridPos.y) &&
			(nodeGridPos.y < bottomRightGridPos.y) {
			node.selected = true
		} else {
			node.selected = false
		}
		s.nodes[nodeGridPos] = node
	}
}

func (s *state) updateTargets(target *gridPosition) {
	for nodeGridPos, node := range s.nodes {
		if node.player == neutralPlayer {
			continue
		}

		if node.selected {
			node.target = target
		}
		s.nodes[nodeGridPos] = node
	}
}

func (s *state) produceUnits() {
	for nodePos, n := range s.nodes {
		if n.player == neutralPlayer {
			continue
		}
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
