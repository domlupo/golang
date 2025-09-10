package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// not enforced but must be a singleton
type state struct {
	nodes  map[gridPosition]*node
	fluxes map[gridPosition]*flux
	units  map[pixelPosition]*unit
}

func newState() *state {
	s := state{
		nodes:  make(map[gridPosition]*node),
		fluxes: make(map[gridPosition]*flux),
		units:  make(map[pixelPosition]*unit),
	}

	for i := range windowWidth / gridSize {
		for j := range windowHeight / gridSize {
			var gridPos = gridPosition{i, j}
			var flux = newFlux(gridPos)
			s.fluxes[gridPos] = &flux
		}
	}

	return &s
}

func (s *state) addNode(n node) {
	s.nodes[n.gridPos] = &n
}

func (s *state) refresh() {
	for _, flux := range s.fluxes {
		flux.refresh()
	}

	for unitPos, unit := range s.units {
		unit.refresh()

		if !unit.movedLastTurn {
			var gridPos, _ = splitPixelPosition(unitPos)
			s.fluxes[gridPos].staticUnits++
		}
	}

	for _, flux := range s.fluxes {
		var xMax = (windowWidth / gridSize) - 1
		var yMax = (windowHeight / gridSize) - 1
		flux.setFluxVector(s, xMax, yMax)
	}
}

func (s *state) analyzeTargets() {
	for unitPos1, unit1 := range s.units {
		var collidesWithNode, isTargetNode = s.collidesWithNode(
			unitPos1,
			*unit1,
			unit1.reach,
			s.nodes)
		if isTargetNode {
			// TODO: HACK, unit can reach further than collide
			unit1.canReachTargetNode = collidesWithNode
		}

		// check if unit 1 can see unit 2
		for unitPos2, unit2 := range s.units {
			if unit1.player == unit2.player {
				// unit cannot see itself or comrades
				continue
			}

			var xDelta = unitPos1.x - unitPos2.x
			var yDelta = unitPos1.y - unitPos2.y
			// TODO: Does this need changed so local target is always what is nearest
			if xDelta*xDelta+yDelta*yDelta <= unit1.vision*unit1.vision {
				unit1.localTarget = &unitPos2
				s.units[unitPos1] = unit1
			}
		}
	}
}

func (s *state) moveUnits() {
	var unitsMovedCount = 0
	var unitsAlreadyAttacked = 0
	var unitsWithoutTarget = 0
	var unitsAtTargetNode = 0
	var unitsAtNonTargetNode = 0
	var unitsCollideCount = 0
	var unitsBorderCollideCount = 0

	for unitPos, unit := range s.units {
		if unit.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		var newPos pixelPosition
		if unit.fluxTurnCount > unitFluxCountLimit {
			var gridPos, _ = splitPixelPosition(unitPos)
			var fluxDelta = s.fluxes[gridPos].getFluxDelta()
			newPos = pixelPosition{unitPos.x + fluxDelta.x, unitPos.y + fluxDelta.y}
			unit.fluxTurnCount = 0
		} else {
			newPos = unit.getNewPos(unitPos)
		}

		if unit.currentTarget == noTarget {
			unitsWithoutTarget++
		}

		if unitPos.x == newPos.x && unitPos.y == newPos.y {
			continue
		}

		var collidesWithBorder = collidesWithBorder(newPos, unit.radius)

		var collidesWithUnit = s.collidesWithUnit(
			newPos,
			*unit,
			s.units)

		var collidesWithNode, isTargetNode = s.collidesWithNode(
			newPos,
			*unit,
			unit.radius,
			s.nodes)

		if collidesWithBorder && collidesWithUnit {
			unitsBorderCollideCount++
			unitsCollideCount++
		} else if collidesWithBorder {
			unitsBorderCollideCount++
		} else if collidesWithUnit {
			unitsCollideCount++
		} else if collidesWithNode != nil {
			if isTargetNode {
				unitsAtTargetNode++
			} else {
				// TODO: determine behavior. Prob don't want units to pass thru enemy nodes
				unitsAtNonTargetNode++
				unitsMovedCount++
				s.units[newPos] = s.units[pixelPosition{unitPos.x, unitPos.y}]
				s.units[newPos].completeMove()
				delete(s.units, pixelPosition{unitPos.x, unitPos.y})
			}
		} else {
			unitsMovedCount++
			s.units[newPos] = s.units[pixelPosition{unitPos.x, unitPos.y}]
			s.units[newPos].completeMove()
			delete(s.units, pixelPosition{unitPos.x, unitPos.y})
		}
	}

	logger.Debug("moveUnits",
		"units moved", unitsMovedCount,
		"units without target", unitsWithoutTarget,
		"units at target node", unitsAtTargetNode,
		"units at non target node", unitsAtNonTargetNode,
		"units collided", unitsCollideCount,
		"units already attacked", unitsAlreadyAttacked,
	)
}

// TODO: This does not handle a unit colliding with multiple nodes which will probably occur in the future.
func (s *state) collidesWithNode(newPos pixelPosition, unit unit, unitCollisionRadius int, nodes map[gridPosition]*node) (*node, bool) {
	var collidesWithTargetNode = false
	var x1 = newPos.x
	var y1 = newPos.y
	var r1 = unitCollisionRadius

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

		var xDelta = x2 - x1
		var yDelta = y2 - y1
		var rSum = r1 + r2

		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			if unit.globalTarget != nil {
				var x1 = unit.globalTarget.x
				var y1 = unit.globalTarget.y
				var r1 = 1

				var xDelta = x2 - x1
				var yDelta = y2 - y1
				var rSum = r1 + r2
				if (xDelta*xDelta)+(yDelta*yDelta) <= (rSum * rSum) {
					collidesWithTargetNode = true
				}
			}
			return node, collidesWithTargetNode
		}
	}

	return nil, collidesWithTargetNode
}

func (s *state) collidesWithUnit(newPos pixelPosition, unit unit, units map[pixelPosition]*unit) bool {
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

		var xDelta = x2 - x1
		var yDelta = y2 - y1
		var rSum = r1 + r2
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
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

			var xDelta = unitPos1.x - unitPos2.x
			var yDelta = unitPos1.y - unitPos2.y
			// TODO: Does this need changed so local target is always what is nearest
			if xDelta*xDelta+yDelta*yDelta < unit1.vision*unit1.vision {
				nearestVisibleUnitPos = &unitPos2
				nearestVisibleUnit = *unit2
			}
		}

		if nearestVisibleUnitPos == nil {
			continue
		}

		var xTargetPixelPos = nearestVisibleUnitPos.x
		var yTargetPixelPos = nearestVisibleUnitPos.y

		var xDelta = 0
		var yDelta = 0
		if unitPos1.x < (xTargetPixelPos) {
			xDelta = xTargetPixelPos - unitPos1.x
		} else if unitPos1.x > (xTargetPixelPos) {
			xDelta = unitPos1.x - xTargetPixelPos
		}
		if unitPos1.y < (yTargetPixelPos) {
			yDelta = yTargetPixelPos - unitPos1.y
		} else if unitPos1.y > (yTargetPixelPos) {
			yDelta = unitPos1.y - yTargetPixelPos
		}
		if unit1.reach*unit1.reach >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - unit1.attack
			unit1.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[*nearestVisibleUnitPos] = &nearestVisibleUnit
			s.units[unitPos1] = unit1
		}
	}

	for attackedUnitPos, attackedUnit := range s.units {
		if attackedUnit.health <= 0 {
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
		if unit.globalTarget == nil {
			noTargetNodeCount++
			continue
		}

		if unit.attackedThisTurn {
			attackedUnitThisTurnCount++
			continue
		}

		if unit.canReachTargetNode != nil {
			attackedNodeThisTurnCount++
			var targetNode = unit.canReachTargetNode
			targetNode.health = max(targetNode.health-unit.attack, 0)
			unit.attackedThisTurn = true
			if targetNode.health == 0 {
				targetNode.capturedBy = &unit.player
			}
			s.nodes[targetNode.gridPos] = targetNode
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

func (s *state) updateNodeOwnership() bool {
	var nodeOwnershipChanged = false

	for nodeGridPos, node := range s.nodes {
		if node.health == 0 {
			newNode := newNode(*node.capturedBy, node.gridPos, node.unitConfig)
			s.nodes[nodeGridPos] = &newNode
			nodeOwnershipChanged = true
		}
	}

	return nodeOwnershipChanged
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
			unit.globalTarget = target
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
		if !s.collidesWithUnit(unitPos, newUnit(n.player, n.target, n.unitConfig), s.units) {
			var newUnit = newUnit(n.player, n.target, n.unitConfig)
			s.units[unitPos] = &newUnit
		}
	}
}

func (s *state) draw(screen *ebiten.Image, drawBlueTargets bool, drawRedTargets, drawFluxes bool) {
	for nodePos, node := range s.nodes {
		node.drawNode(screen, nodePos)
	}
	for nodePos, node := range s.nodes {
		node.drawTargets(screen, nodePos, drawBlueTargets, drawRedTargets)
	}
	for unitPos, unit := range s.units {
		unit.draw(screen, unitPos, unit.selected)
	}
	if drawFluxes {
		for fluxPos, flux := range s.fluxes {
			flux.draw(screen, fluxPos)
		}
	}
}
