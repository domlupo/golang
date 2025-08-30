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

func (s *state) addNode(n node) {
	s.nodes[n.gridPos] = n
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
	var unitsMovedCount = 0
	var unitsWithoutTarget = 0
	var unitsAtTargetNode = 0
	var unitsAtNonTargetNode = 0
	var unitsCollidedCount = 0

	var newUnits = make(map[pixelPosition]unit)
	for unitPos1, unit1 := range s.units {
		unit1.canReachTargetNode = nil

		// check if unit 1 can see unit 2
		var nearestVisibleUnitPos *pixelPosition
		for unitPos2, unit2 := range s.units {
			if unit1.player == unit2.player {
				// unit cannot see comrades
				continue
			}

			var distanceBetweenUnits = math.Sqrt(
				math.Pow(float64(unitPos1.x-unitPos2.x), 2) +
					math.Pow(float64(unitPos1.y-unitPos2.y), 2))
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
			xTargetPixelPos = unit1.target.x
			yTargetPixelPos = unit1.target.y
		} else {
			unitsWithoutTarget++
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

		var collidesWithUnit = s.collidesWithUnit(
			pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta},
			unit1,
			s.units)

		var collidesWithNode, isTargetNode = s.collidesWithNode(
			pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta},
			unit1,
			s.nodes)

		if collidesWithUnit {
			unitsCollidedCount++
			newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
		} else if collidesWithNode != nil {
			if isTargetNode {
				unitsAtTargetNode++
				// TODO: HACK, unit can reach further than collide
				unit1.canReachTargetNode = collidesWithNode
				newUnits[pixelPosition{unitPos1.x, unitPos1.y}] = unit1
			} else {
				// TODO: determine behavior. Prob don't want units to pass thru enemy nodes
				unitsAtNonTargetNode++
				unitsMovedCount++
				newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}] = unit1
			}
		} else {
			unitsMovedCount++
			newUnits[pixelPosition{unitPos1.x + xDelta, unitPos1.y + yDelta}] = unit1
		}
	}

	logger.Debug("moveUnits",
		"units moved", unitsMovedCount,
		"units without target", unitsWithoutTarget,
		"units at target node", unitsAtTargetNode,
		"units at non target node", unitsAtNonTargetNode,
		"units collided", unitsCollidedCount,
	)

	s.units = newUnits
}

func (s *state) collidesWithNode(newPos pixelPosition, unit unit, nodes map[gridPosition]node) (*node, bool) {
	var collidesWithTargetNode = false
	var x1 = newPos.x
	var y1 = newPos.y
	var r1 = unit.radius

	for nodePos, node := range nodes {
		if unit.player == node.player {
			// unit cannot collide with own nodes
			// HACK: This will probably change in the future, and it will most likely break attackNodes logic
			// since it assumes a unit never collides with its own node
			continue
		}

		var nodePixelPos = getPixelPosition(nodePos, innerGridPosition{0, 0})

		var x2 = nodePixelPos.x
		var y2 = nodePixelPos.y
		var r2 = node.radius

		if math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2) <= math.Pow((float64(r1)+float64(r2)), 2) {
			var x1 = unit.target.x
			var y1 = unit.target.y
			var r1 = 1

			if math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2) <= math.Pow(float64(r1)+float64(r2), 2) {
				collidesWithTargetNode = true
			}
			return &node, collidesWithTargetNode
		}
	}

	return nil, collidesWithTargetNode
}

func (s *state) collidesWithUnit(newPos pixelPosition, unit unit, units map[pixelPosition]unit) bool {
	var x1 = newPos.x + unit.radius
	var y1 = newPos.y + unit.radius
	var r1 = unit.radius

	for otherPos, otherUnit := range units {
		if unit.id == otherUnit.id {
			// unit cannot collide with itself
			continue
		}

		var x2 = otherPos.x + otherUnit.radius
		var y2 = otherPos.y + otherUnit.radius
		var r2 = otherUnit.radius
		if math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2) <= math.Pow((float64(r1)+float64(r2)), 2) {
			return true
		}
	}

	return false
}

func (s *state) attackUnits() {
	var unitsDeletedCount = 0
	var unitAttackedOtherUnitCount = 0

	for unitPos1, unit1 := range s.units {
		// check if unit 1 can attack unit 2
		var nearestVisibleUnitPos *pixelPosition = nil
		var nearestVisibleUnit unit
		for unitPos2, unit2 := range s.units {
			if unit1.id == unit2.id {
				// unit cannot attack itself
				continue
			}

			if unit1.player == unit2.player {
				// unit cannot attack comarades
				continue
			}

			var distanceBetweenUnits = math.Sqrt(
				math.Pow(float64(unitPos1.x-unitPos2.x), 2) +
					math.Pow(float64(unitPos1.y-unitPos2.y), 2))
			// TODO: setting this <= rather than < caused obvious bugs but it seems more correct
			// to use <. Need to determine right behavior in the future.
			if distanceBetweenUnits < float64(unit1.vision) {
				nearestVisibleUnitPos = &unitPos2
				nearestVisibleUnit = unit2
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
		var distanceToNearestVisibleUnit = math.Sqrt(math.Pow(xMagnitude, 2) + math.Pow(yMagnitude, 2))
		if float64(unit1.reach) >= distanceToNearestVisibleUnit {
			nearestVisibleUnit.alive = false
			s.units[*nearestVisibleUnitPos] = nearestVisibleUnit
			unit1.attackedUnitThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[unitPos1] = unit1
		}
	}

	for attackedUnitPos, attackedUnit := range s.units {
		if !attackedUnit.alive {
			unitsDeletedCount++
			delete(s.units, attackedUnitPos)
		}
	}

	logger.Debug("attackUnits",
		"units attacked", unitAttackedOtherUnitCount,
		"units deleted", unitsDeletedCount,
	)
}

func (s *state) attackNodes() {
	var noTargetNodeCount = 0
	var attackedUnitThisTurnCount = 0
	var attackedNodeThisTurnCount = 0
	var cannotReachNodeCount = 0

	for _, unit := range s.units {
		if unit.target == nil {
			noTargetNodeCount++
			continue
		}

		if unit.attackedUnitThisTurn {
			attackedUnitThisTurnCount++
			continue
		}

		if unit.canReachTargetNode != nil {
			attackedNodeThisTurnCount++
			var targetNode = unit.canReachTargetNode
			targetNode.health = max(targetNode.health-1, 0)
			if targetNode.health == 0 {
				targetNode.capturedBy = &unit.player
			}
			s.nodes[targetNode.gridPos] = *targetNode
		} else {
			cannotReachNodeCount++
		}
	}

	logger.Debug("attackNodes",
		"no target node", noTargetNodeCount,
		"attacked unit this turn", attackedUnitThisTurnCount,
		"attacked node this turn", attackedNodeThisTurnCount,
		"cannot reach node", cannotReachNodeCount,
	)
}

func (s *state) updateNodeOwnership() {
	for nodeGridPos, node := range s.nodes {
		if node.health == 0 {
			newNode := newNode(*node.capturedBy, node.gridPos)
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

func (s *state) updateTargets(target *pixelPosition) {
	for nodeGridPos, node := range s.nodes {
		if node.player == neutralPlayer {
			continue
		}

		if node.selected {
			node.target = target
		}
		s.nodes[nodeGridPos] = node
	}
	for unitPixelPos, unit := range s.units {
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
		var unitPos = getPixelPosition(nodePos, innerGridPosition{0, 0})
		if !s.collidesWithUnit(unitPos, newUnit(n.player, n.target), s.units) {
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
