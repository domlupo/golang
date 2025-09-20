// Currently contains most game state and game logic. May need to split in future.
package main

import (
	"math"
)

// not enforced but must be a singleton
type state struct {
	windowWidth  int
	windowHeight int
	gridSize     int
	nodes        map[gridPosition]*node
	fluxes       map[gridPosition]*flux
	units        map[position]*unit
}

func newState(windowWidth int, windowHeight int, gridSize int) *state {
	s := state{
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		gridSize:     gridSize,
		nodes:        make(map[gridPosition]*node),
		fluxes:       make(map[gridPosition]*flux),
		units:        make(map[position]*unit),
	}

	for i := range s.windowWidth / s.gridSize {
		for j := range s.windowHeight / s.gridSize {
			gridPos := gridPosition{i, j}
			flux := newFlux(gridPos)
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

	for uPos, u := range s.units {
		u.refresh()

		if !u.movedLastTurn {
			gridPos, _ := splitPixelPosition(uPos.pixelPosition, s.gridSize)
			s.fluxes[gridPos].staticUnits++
		}
	}

	for _, flux := range s.fluxes {
		xMax := s.windowWidth/s.gridSize - 1
		yMax := s.windowHeight/s.gridSize - 1
		flux.setFluxVector(s, xMax, yMax)
	}
}

func (s *state) analyzeTargets() {
	for uPos1, u1 := range s.units {
		if u1.staticTarget {
			continue
		}

		if u1.fluxTurnCount > u1.fluxTurnCountDelay {
			gridPos, _ := splitPixelPosition(uPos1.pixelPosition, s.gridSize)
			fluxDelta := s.fluxes[gridPos].getFluxDelta(u1.fluxMagnitude)
			u1.updateTarget(fluxTarget, &position{ground, pixelPosition{uPos1.x + fluxDelta.x, uPos1.y + fluxDelta.y}})
			u1.fluxTurnCount = 0
		}

		var reachesNode, isTargetNode = s.collidesWithNode(
			uPos1,
			*u1,
			u1.reach,
			s.nodes)
		if isTargetNode {
			u1.canReachTargetNode = reachesNode
		}

		shortestDistance := math.MaxInt64
		var closestUnitPos *position
		for uPos2, u2 := range s.units {
			if u1.player == u2.player {
				// unit cannot see itself or comrades
				continue
			}

			if !u1.attacksGround && uPos2.layer == ground ||
				!u1.attacksLowAltitude && uPos2.layer == lowAltitude ||
				!u1.attacksHighAltitude && uPos2.layer == highAltitude {
				continue
			}

			xDelta := uPos1.x - uPos2.x
			yDelta := uPos1.y - uPos2.y
			distance := xDelta*xDelta + yDelta*yDelta
			if distance < shortestDistance && distance < u1.vision*u1.vision {
				shortestDistance = distance
				closestUnitPos = &uPos2
			}
		}

		u1.updateTarget(localTarget, closestUnitPos)

	}
}

func (s *state) moveUnits() {
	logger.log(infoLevel, stateLogType, "moveUnits", "unitCount", len(s.units))

	var unitsMovedCount int
	var unitsAlreadyAttacked int
	var unitsWithoutTarget int
	var unitsAtTargetNode int
	var unitsAtNonTargetNode int
	var unitsCollideCount int
	var unitsBorderCollideCount int

	for uPos, u := range s.units {
		if u.movedThisTurn {
			continue
		}

		if u.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		// TODO: This updating target seems a bit wrong. Check if there is a better way to handle this.
		// My initial thoughts are to update unit targets only in analyze targets then in move units
		// I simply get the move. May need to split some logic up.
		newPos := u.getNewPos(uPos)

		if u.currentTarget == nil {
			unitsWithoutTarget++
		}

		if uPos.x == newPos.x && uPos.y == newPos.y {
			continue
		}

		collidesWithBorder := s.collidesWithBorder(newPos, u.radius)

		var collidesWithUnit = s.collidesWithUnit(
			newPos,
			*u,
			s.units)

		var collidesWithNode, isTargetNode = s.collidesWithNode(
			newPos,
			*u,
			u.radius,
			s.nodes)

		if collidesWithBorder && collidesWithUnit != nil {
			// TODO: Low priority but eventually need to handle air collision here.
			// Most likely can refactor to avoid updating this code directly.
			unitsBorderCollideCount++
			unitsCollideCount++
		} else if collidesWithBorder {
			unitsBorderCollideCount++
		} else if collidesWithUnit != nil {
			unitsCollideCount++
			if uPos.layer == lowAltitude || uPos.layer == highAltitude {
				otherUnit := s.units[*collidesWithUnit]
				otherUnit.health = 0
				s.units[*collidesWithUnit] = s.units[position{collidesWithUnit.layer, pixelPosition{collidesWithUnit.x, collidesWithUnit.y}}]
				u.health = 0
				s.units[uPos] = s.units[position{uPos.layer, pixelPosition{uPos.x, uPos.y}}]
			}
		} else if collidesWithNode != nil {
			if isTargetNode {
				unitsAtTargetNode++
			} else {
				// TODO: determine behavior. Prob don't want us to pass thru enemy nodes
				unitsAtNonTargetNode++
				unitsMovedCount++
				s.units[newPos] = s.units[position{uPos.layer, pixelPosition{uPos.x, uPos.y}}]
				s.units[newPos].completeMove()
				delete(s.units, position{uPos.layer, pixelPosition{uPos.x, uPos.y}})
			}
		} else {
			unitsMovedCount++
			s.units[newPos] = s.units[position{uPos.layer, pixelPosition{uPos.x, uPos.y}}]
			s.units[newPos].completeMove()
			delete(s.units, position{uPos.layer, pixelPosition{uPos.x, uPos.y}})
		}
	}

	logger.log(
		debugLevel,
		stateLogType,
		"moveUnits",
		"units moved", unitsMovedCount,
		"units without target", unitsWithoutTarget,
		"units at target node", unitsAtTargetNode,
		"units at non target node", unitsAtNonTargetNode,
		"units collided", unitsCollideCount,
		"units already attacked", unitsAlreadyAttacked,
	)
}

// TODO: This does not handle a unit colliding with multiple nodes which will probably occur in the future.
func (s *state) collidesWithNode(pos position, unit unit, unitCollisionRadius int, nodes map[gridPosition]*node) (*node, bool) {
	if pos.layer == lowAltitude || pos.layer == highAltitude {
		return nil, false
	}

	var collidesWithTargetNode bool
	x1 := pos.x
	y1 := pos.y
	r1 := unitCollisionRadius

	for nPos, n := range nodes {
		if unit.player == n.player {
			// unit cannot collide with own nodes
			// HACK: This will probably change in the future, and it will most likely break attackNodes logic
			// since it assumes a unit never collides with its own node
			continue
		}

		nPixelPos := getPixelPosition(nPos, topLeft, s.gridSize)

		x2 := nPixelPos.x
		y2 := nPixelPos.y
		r2 := n.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2

		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			if unit.globalTarget != nil {
				x1 := unit.globalTarget.x
				y1 := unit.globalTarget.y
				r1 := 1 // TODO: why is this 1? Is this a bug?

				xDelta := x2 - x1
				yDelta := y2 - y1
				rSum := r1 + r2
				if (xDelta*xDelta)+(yDelta*yDelta) <= (rSum * rSum) {
					collidesWithTargetNode = true
				}
			}
			return n, collidesWithTargetNode
		}
	}

	return nil, collidesWithTargetNode
}

func (s *state) explode(pos position, blastRadius int, units map[position]*unit) {
	x1 := pos.x
	y1 := pos.y
	r1 := blastRadius

	for otherPos, otherUnit := range units {
		if otherPos.layer != ground {
			continue
		}

		x2 := otherPos.x + otherUnit.radius
		y2 := otherPos.y + otherUnit.radius
		r2 := otherUnit.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			otherUnit.health = 0
			s.units[otherPos] = otherUnit
		}
	}
}

func (s *state) collidesWithUnit(pos position, unit unit, units map[position]*unit) *position {
	x1 := pos.x + unit.radius
	y1 := pos.y + unit.radius
	r1 := unit.radius

	for otherPos, otherUnit := range units {
		if unit.id == otherUnit.id {
			// unit cannot collide with itself
			continue
		}

		if pos.layer != otherPos.layer {
			continue
		}

		x2 := otherPos.x + otherUnit.radius
		y2 := otherPos.y + otherUnit.radius
		r2 := otherUnit.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			return &otherPos
		}
	}

	return nil
}

func (s *state) processTerminalPositions() {
	for uPos, u := range s.units {
		if !u.explodeAtTerminalPos {
			continue
		}

		if u.attacksPosition && u.atTerminalPosition {
			u.health = 0
			s.units[uPos] = u
			s.explode(uPos, u.explodeRadius, s.units)
		} else if u.atTerminalPosition {
			u.health = 0
			s.units[uPos] = u
		}

	}
}

func (s *state) attackUnits() {
	var unitsDeletedCount int
	var unitAttackedOtherUnitCount int

	for uPos, u := range s.units {
		// TODO: This may be wrong behavior but I am not sure if the issue is occurring
		// e.g. the target was killed by another unit this unit should probably find a potential new target.
		// Then attack that rather than not attack for a turn
		if u.localTarget == nil {
			continue
		}

		xTargetPixelPos := u.localTarget.x
		yTargetPixelPos := u.localTarget.y

		var xDelta int
		var yDelta int
		if uPos.x < (xTargetPixelPos) {
			xDelta = xTargetPixelPos - uPos.x
		} else if uPos.x > (xTargetPixelPos) {
			xDelta = uPos.x - xTargetPixelPos
		}
		if uPos.y < (yTargetPixelPos) {
			yDelta = yTargetPixelPos - uPos.y
		} else if uPos.y > (yTargetPixelPos) {
			yDelta = uPos.y - yTargetPixelPos
		}

		var nearestVisibleUnit = s.units[*u.localTarget]
		var rSum = u.reach + nearestVisibleUnit.radius
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[*u.localTarget] = nearestVisibleUnit
			s.units[uPos] = u
		}
	}

	for unitPos, unit := range s.units {
		if unit.health <= 0 {
			unitsDeletedCount++
			delete(s.units, unitPos)
		}
	}

	logger.log(
		debugLevel,
		stateLogType,
		"attackUnits",
		"units attacked", unitAttackedOtherUnitCount,
		"units deleted", unitsDeletedCount,
	)
}

func (s *state) attackNodes() {
	var noTargetNodeCount int
	var attackedUnitThisTurnCount int
	var attackedNodeThisTurnCount int
	var cannotReachNodeCount int

	for _, u := range s.units {
		if u.globalTarget == nil {
			noTargetNodeCount++
			continue
		}

		if u.attackedThisTurn {
			attackedUnitThisTurnCount++
			continue
		}

		if u.canReachTargetNode != nil {
			attackedNodeThisTurnCount++
			targetNode := u.canReachTargetNode
			targetNode.currentHealth = max(targetNode.currentHealth-u.attack, 0)
			u.attackedThisTurn = true
			if targetNode.currentHealth == 0 {
				targetNode.capturedBy = &u.player
			}
			s.nodes[targetNode.gridPos] = targetNode
		} else {
			cannotReachNodeCount++
		}
	}

	logger.log(
		debugLevel,
		stateLogType,
		"attackNodes",
		"no target node", noTargetNodeCount,
		"attacked unit this turn", attackedUnitThisTurnCount,
		"attacked node this turn", attackedNodeThisTurnCount,
		"cannot reach node", cannotReachNodeCount,
	)
}

func (s *state) updateNodeOwnership() bool {
	var nodeOwnershipChanged bool

	for nGridPos, n := range s.nodes {
		if n.currentHealth == 0 {
			newNode := newNode(*n.capturedBy, n.gridPos, n.maxHealth, n.radius, n.unitConfig)
			s.nodes[nGridPos] = &newNode
			nodeOwnershipChanged = true
		}
	}

	return nodeOwnershipChanged
}

func (s *state) updateSelected(clickPosition pixelPosition, releasePosition pixelPosition) {
	var topLeftPixelPos pixelPosition
	var bottomRightPixelPos pixelPosition

	switch {
	case clickPosition.x < releasePosition.x && clickPosition.y < releasePosition.y:
		topLeftPixelPos = clickPosition
		bottomRightPixelPos = releasePosition
	case clickPosition.x > releasePosition.x && clickPosition.y > releasePosition.y:
		topLeftPixelPos = releasePosition
		bottomRightPixelPos = clickPosition
	case clickPosition.x > releasePosition.x && clickPosition.y < releasePosition.y:
		topLeftPixelPos = pixelPosition{releasePosition.x, clickPosition.y}
		bottomRightPixelPos = pixelPosition{clickPosition.x, releasePosition.y}
	case clickPosition.x < releasePosition.x && clickPosition.y > releasePosition.y:
		topLeftPixelPos = pixelPosition{clickPosition.x, releasePosition.y}
		bottomRightPixelPos = pixelPosition{releasePosition.x, clickPosition.y}
	}

	topLeftGridPos, _ := splitPixelPosition(topLeftPixelPos, s.gridSize)
	bottomRightGridPos, _ := splitPixelPosition(bottomRightPixelPos, s.gridSize)

	for nGridPos, n := range s.nodes {
		if (topLeftGridPos.x < nGridPos.x) &&
			(nGridPos.x < bottomRightGridPos.x) &&
			(topLeftGridPos.y < nGridPos.y) &&
			(nGridPos.y < bottomRightGridPos.y) {
			n.selected = true
		} else {
			n.selected = false
		}
		s.nodes[nGridPos] = n
	}

	for unitPixelPos, unit := range s.units {
		if (topLeftPixelPos.x < unitPixelPos.x) &&
			(unitPixelPos.x < bottomRightPixelPos.x) &&
			(topLeftPixelPos.y < unitPixelPos.y) &&
			(unitPixelPos.y < bottomRightPixelPos.y) {
			unit.selected = true
		} else {
			unit.selected = false
		}
		s.units[unitPixelPos] = unit
	}

}

func (s *state) updateTargets(target *position, blueComputerActive bool, redComputerActive bool) {
	for nGridPos, n := range s.nodes {
		switch {
		case n.player == neutralPlayer:
			continue
		case n.player == bluePlayer && blueComputerActive:
			continue
		case n.player == redPlayer && redComputerActive:
			continue
		}

		if n.selected {
			n.target = target
		}
		s.nodes[nGridPos] = n
	}
	for _, u := range s.units {
		if u.player == bluePlayer && blueComputerActive {
			continue
		} else if u.player == redPlayer && redComputerActive {
			continue
		}

		if u.selected {
			u.updateTarget(globalTarget, target)
		}
	}
}

func (s *state) produceUnits() {
	for nPos, n := range s.nodes {
		if n.player == neutralPlayer {
			continue
		}

		// TODO: determine how to produce units during collision
		uPos := position{n.unitConfig.produceLayer, getPixelPosition(nPos, topLeft, s.gridSize)}
		if s.collidesWithUnit(uPos, newUnit(n.player, n.target, n.unitConfig), s.units) == nil {
			newUnit := newUnit(n.player, n.target, n.unitConfig)
			s.units[uPos] = &newUnit
		}
	}
}

func (s *state) collidesWithBorder(pos position, offset int) bool {
	if pos.pixelPosition.x+offset > s.windowWidth ||
		pos.pixelPosition.x-offset < 0 ||
		pos.pixelPosition.y+offset > s.windowHeight ||
		pos.pixelPosition.y-offset < 0 {
		return true
	} else {
		return false
	}
}

func (s *state) draw(dc drawContext, drawBlueTargets bool, drawRedTargets, drawFluxes bool) {
	for nodePos, node := range s.nodes {
		node.drawNode(dc, nodePos)
	}
	for nodePos, node := range s.nodes {
		node.drawTargets(dc, nodePos, drawBlueTargets, drawRedTargets)
	}
	for unitPos, unit := range s.units {
		if unitPos.layer == highAltitude {
			unit.draw(dc, unitPos, unit.selected)
		}
	}
	for unitPos, unit := range s.units {
		if unitPos.layer == lowAltitude {
			unit.draw(dc, unitPos, unit.selected)
		}
	}
	for unitPos, unit := range s.units {
		if unitPos.layer == ground {
			unit.draw(dc, unitPos, unit.selected)
		}
	}
	if drawFluxes {
		for fluxPos, flux := range s.fluxes {
			flux.draw(dc, fluxPos)
		}
	}
}
