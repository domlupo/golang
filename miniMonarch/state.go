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
		} else if unit1.target != nil {
			xTargetPixelPos = (unit1.target.x * gridSize) + gridSize/2
			yTargetPixelPos = (unit1.target.y * gridSize) + gridSize/2
		} else {
			newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
			continue
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

		var collides = s.collides(
			pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta},
			unit1,
			s.units)

		// TODO: Determine behavior wanted when units travel through different player nodes
		// TODO: Add nodes to collision detection
		var gridPos, _ = splitPixelPosition(pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta})

		if unit1.target != nil && gridPos.x == unit1.target.x && gridPos.y == unit1.target.y {
			unit1.nextToTargetNode = true
		} else {
			unit1.nextToTargetNode = false
		}

		if collides || unit1.nextToTargetNode {
			newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
		} else {
			newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}] = unit1
		}
	}

	s.units = newUnits
}

func (s *state) collides(newPos pixelPosition, unit unit, units map[pixelPosition]unit) bool {
	var left = newPos.x
	var right = newPos.x + unit.width
	var top = newPos.y
	var bottom = newPos.y + unit.height

	for otherPos, otherUnit := range units {
		if unit.id == otherUnit.id {
			// unit cannot collide with itself
			continue
		}

		// square AABB collision detection
		if (left <= otherPos.x+otherUnit.width) &&
			(right >= otherPos.x) &&
			(top <= otherPos.y+otherUnit.height) &&
			(bottom >= otherPos.y) {
			return true
		}
	}

	return false
}

func (s *state) attackUnits() {
	for unitPos1, unit1 := range s.units {
		// check if unit 1 can attack unit 2
		var nearestVisibleUnitPos *pixelPosition = nil
		var nearestVisibleUnit *unit = nil
		for unitPos2, unit2 := range s.units {
			if unit1.id == unit2.id {
				// unit cannot attack itself
				continue
			}

			if unit1.player == unit2.player {
				// unit cannot attack comarades
				continue
			}

			var unit1XCenter = unitPos1.x + unit1.width/2
			var unit2XCenter = unitPos2.x + unit2.width/2
			var unit1YCenter = unitPos1.y + unit1.height/2
			var unit2YCenter = unitPos2.y + unit2.height/2

			var distanceBetweenUnits = math.Sqrt(math.Pow(float64(unit1XCenter-unit2XCenter), 2) + math.Pow(float64(unit1YCenter-unit2YCenter), 2))
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
		if unit.target == nil || unit.attackedUnitThisTurn {
			continue
		}

		var attackablePlayer player
		var targetedPlayer = s.nodes[*unit.target].player
		switch unit.player {
		case bluePlayer:
			attackablePlayer = redPlayer
		case redPlayer:
			attackablePlayer = bluePlayer
		}

		if unit.nextToTargetNode && (targetedPlayer == neutralPlayer || targetedPlayer == attackablePlayer) {
			var targetNode = s.nodes[*unit.target]
			targetNode.health = max(targetNode.health-1, 0)
			if targetNode.health == 0 {
				targetNode.capturedBy = &unit.player
			}
			s.nodes[*unit.target] = targetNode
		}
	}
}

func (s *state) updateNodeOwnership() {
	for nodeGridPos, node := range s.nodes {
		if node.health == 0 {
			newNode := newNode(*node.capturedBy)
			s.nodes[nodeGridPos] = newNode
		}
	}
}

func (s *state) updateSelected(clickPosition pixelPosition, releasePosition pixelPosition) {
	// TODO: generalize implementation, for prototype this assumes bottom right release position.
	// The future implementation needs to handle all four cases for release position.

	var topLeftGridPos, _ = splitPixelPosition(clickPosition)
	var bottomRightGridPos, _ = splitPixelPosition(releasePosition)

	for nodeGridPos, node := range s.nodes {
		if node.player == bluePlayer && !bluePlayerActive {
			continue
		} else if node.player == redPlayer && !redPlayerActive {
			continue
		} else if node.player == neutralPlayer {
			continue
		}

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

	for unitPixelPos, unit := range s.units {
		if unit.player == bluePlayer && !bluePlayerActive {
			continue
		} else if unit.player == redPlayer && !redPlayerActive {
			continue
		}

		if (clickPosition.x < unitPixelPos.x) &&
			(unitPixelPos.x < releasePosition.x) &&
			(clickPosition.y < unitPixelPos.y) &&
			(unitPixelPos.y < releasePosition.y) {
			unit.selected = true
		} else {
			unit.selected = false
		}
		s.units[unitPixelPos] = unit
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

	var targetPlayer = s.nodes[*target].player
	for unitPixelPos, unit := range s.units {
		if unit.player == targetPlayer {
			continue
		}

		if unit.selected {
			unit.target = target
		}
		s.units[unitPixelPos] = unit
	}
}

func (s *state) produceUnits() {
	for nodePos, n := range s.nodes {
		if n.player == neutralPlayer {
			continue
		}

		// TODO: determine how to produce units during collision
		var unitPos = getPixelPosition(nodePos, innerGridPosition{gridSize / 2, gridSize / 2})
		if !s.collides(unitPos, newUnit(n.player, n.target), s.units) {
			s.units[unitPos] = newUnit(n.player, n.target)
		}
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
		unit.draw(screen, unitPos, unit.selected)
	}
}
