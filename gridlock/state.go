// Currently contains most game state and game logic. May need to split in future.
package main

import (
	"math"
)

// not enforced but must be a singleton
type state struct {
	*hexConfig
	// TODO: The next largest optimization might be to split hexagon linked list
	// into respective layers.
	hexagons           map[gridPosition]*hexagon
	nodes              map[gridPosition]*node
	fluxes             map[gridPosition]*flux
	units              map[position]*unit
	unitsProducedTotal int
	unitsDeletedTotal  int
}

func newState(hexConfig *hexConfig) *state {
	s := state{
		hexConfig: hexConfig,
		hexagons:  make(map[gridPosition]*hexagon),
		nodes:     make(map[gridPosition]*node),
		fluxes:    make(map[gridPosition]*flux),
		units:     make(map[position]*unit),
	}

	for i := range s.hexConfig.xIndexMax {
		for j := range s.hexConfig.yIndexMax {
			gridPos := gridPosition{i, j}
			flux := newFlux(gridPos)
			s.fluxes[gridPos] = &flux
			hexagon := hexagon{}
			s.hexagons[gridPos] = &hexagon
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

	for _, u := range s.units {
		u.refresh()

		if !u.movedLastTurn {
			gridPos, _ := getGridPos(u.pixelPosition, s.hexConfig)
			s.fluxes[gridPos].staticUnits++
		}
	}

	for _, flux := range s.fluxes {
		flux.setFluxVector(s)
	}
}

func (s *state) analyzeTargets() {
	for _, u := range s.units {
		if u.staticTarget {
			continue
		}

		if u.fluxTurnCount > u.fluxTurnCountDelay {
			gridPos, _ := getGridPos(u.pixelPosition, s.hexConfig)
			fluxDelta := s.fluxes[gridPos].getFluxDelta(u.fluxMagnitude)
			u.updateTarget(fluxTarget, position{ground, pixelPosition{u.x + fluxDelta.x, u.y + fluxDelta.y}})
		}

		var node, isTargetNode, _ = s.collidesWithNode(
			u.position,
			*u,
			u.reach,
			s.nodes)

		if isTargetNode {
			u.canReachTargetNode = node
		}

		nearest, _ := s.nearestUnitWithinVision(*u)

		u.updateTarget(localTarget, nearest)
		u.analyzeTargets()
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

	for _, u := range s.units {
		if u.movedThisTurn {
			continue
		}

		if u.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		// TODO: Consider refactoring this function into a separate moveUnit. Can use named
		// returns to defer completeMove call on successful moves. Punting this because
		// the function will probably change a lot in the future and the refactor can happen then.
		newPos, completeMove := u.getNewPos()

		if u.currentTarget.isZero() {
			unitsWithoutTarget++
		}

		if u.x == newPos.x && u.y == newPos.y {
			continue
		}

		collidesWithBorder := s.collidesWithBorder(newPos, u.radius)

		var collidedUnitPos, collidesWithUnit = s.collidesWithUnit(
			newPos,
			*u)

		var _, isTargetNode, collidesWithNode = s.collidesWithNode(
			newPos,
			*u,
			u.radius,
			s.nodes)

		if collidesWithBorder && collidesWithUnit {
			// TODO: Low priority but eventually need to handle air collision here.
			// Most likely can refactor to avoid updating this code directly.
			unitsBorderCollideCount++
			unitsCollideCount++
		} else if collidesWithBorder {
			unitsBorderCollideCount++
		} else if collidesWithUnit {
			unitsCollideCount++
			if u.layer == lowAltitude || u.layer == highAltitude {
				otherUnit := s.units[collidedUnitPos]
				otherUnit.health = 0
				s.units[collidedUnitPos] = s.units[position{collidedUnitPos.layer, pixelPosition{collidedUnitPos.x, collidedUnitPos.y}}]
				u.health = 0
				s.units[u.position] = s.units[position{u.layer, pixelPosition{u.x, u.y}}]
			}
		} else if collidesWithNode {
			if isTargetNode {
				unitsAtTargetNode++
			} else {
				panic("This path doesn't really happen yet. Handle it when it starts to occur if it ever does...")
			}
		} else {
			unitsMovedCount++
			s.units[newPos] = s.units[position{u.layer, pixelPosition{u.x, u.y}}]

			oldGridPos, _ := getGridPos(u.pixelPosition, s.hexConfig)
			newGridPos, _ := getGridPos(newPos.pixelPosition, s.hexConfig)
			if oldGridPos.x != newGridPos.x || oldGridPos.y != newGridPos.y {
				oldHex := s.hexagons[oldGridPos]
				newHex := s.hexagons[newGridPos]
				oldHex.unitCount--
				newHex.unitCount++

				if u.next.id == u.id && u.prior.id == u.id {
					oldHex.linkedList = nil
				} else if u.next.id == u.id || u.prior.id == u.id {
					panic("unit is in an invalid state")
				} else {
					oldHex.linkedList = u.next
				}

				u.next.prior = u.prior
				u.prior.next = u.next

				updatePointers(newHex, u)
			}

			delete(s.units, u.position)
			completeMove()
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

func updatePointers(h *hexagon, u *unit) {
	if h.linkedList == nil {
		h.linkedList = u
		h.linkedList.prior = u
		h.linkedList.next = u
	}
	u.next = h.linkedList
	u.prior = h.linkedList.prior
	h.linkedList.prior.next = u
	h.linkedList.prior = u
}

// TODO: This does not handle a unit colliding with multiple nodes which will probably occur in the future.
func (s *state) collidesWithNode(pos position, u unit, unitCollisionRadius int, nodes map[gridPosition]*node) (node, bool, bool) {
	if pos.layer == lowAltitude || pos.layer == highAltitude {
		return node{}, false, false
	}

	var collidesWithTargetNode bool
	x1 := pos.x
	y1 := pos.y
	r1 := unitCollisionRadius

	for _, n := range nodes {
		if u.player == n.player {
			// unit cannot collide with own nodes
			// HACK: This will probably change in the future, and it will most likely break attackNodes logic
			// since it assumes a unit never collides with its own node
			continue
		}

		nPixelPos, _ := getPixelPosition(n.gridPos, center, s.hexConfig)

		x2 := nPixelPos.x
		y2 := nPixelPos.y
		r2 := n.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2

		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			var globalTargetGridPos gridPosition
			var globalTargetGridPosOk bool
			if !u.globalTarget.isZero() {
				globalTargetGridPos, globalTargetGridPosOk = getGridPos(u.globalTarget.pixelPosition, s.hexConfig)
			}
			collidesWithTargetNode := globalTargetGridPosOk && globalTargetGridPos.x == n.gridPos.x && globalTargetGridPos.y == n.gridPos.y
			return *n, collidesWithTargetNode, true
		}
	}

	return node{}, collidesWithTargetNode, false
}

func (s *state) explode(pos position, blastRadius int, units map[position]*unit) {
	x1 := pos.x
	y1 := pos.y
	r1 := blastRadius

	for _, otherUnit := range units {
		if otherUnit.layer != ground {
			continue
		}

		x2 := otherUnit.x + otherUnit.radius
		y2 := otherUnit.y + otherUnit.radius
		r2 := otherUnit.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			otherUnit.health = 0
			s.units[otherUnit.position] = otherUnit
		}
	}
}

func (s *state) nearestUnitWithinVision(unit unit) (position, bool) {
	x1 := unit.position.x + unit.radius
	y1 := unit.position.y + unit.radius

	gridPos, _ := getGridPos(unit.pixelPosition, s.hexConfig)
	hex := s.hexagons[gridPos]
	var neHex, eHex, seHex, swHex, wHex, nwHex *hexagon

	if pos, ok := gridPos.adjacent(northeastSide, s.hexConfig); ok {
		neHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(eastSide, s.hexConfig); ok {
		eHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(southeastSide, s.hexConfig); ok {
		seHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(southwestSide, s.hexConfig); ok {
		swHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(westSide, s.hexConfig); ok {
		wHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(northwestSide, s.hexConfig); ok {
		nwHex = s.hexagons[pos]
	}

	hexagons := [7]*hexagon{hex, neHex, eHex, seHex, swHex, wHex, nwHex}

	shortestDistance := math.MaxInt64
	var nearestUnitPos position

	for _, hex := range hexagons {
		if hex == nil {
			continue
		}

		linkedList := hex.linkedList

		if linkedList == nil || (linkedList.id == unit.id && linkedList == linkedList.next) {
			continue
		}

		end := linkedList
		current := linkedList
		// checkOccurred handles the edge case where a hexagon has a single
		// unit within in. This unit should still be checked since it
		// is not the unit passed into the function but don't self cycle
		// since it will never progress to a different current unit.
		var checkOccurred bool
		for {
			if current.id == end.id && checkOccurred {
				break
			}

			if current.id == unit.id {
				current = current.next
				continue
			}

			checkOccurred = true

			if unit.player == current.player {
				current = current.next
				continue
			}

			if !unit.attacksGround && current.layer == ground ||
				!unit.attacksLowAltitude && current.layer == lowAltitude ||
				!unit.attacksHighAltitude && current.layer == highAltitude {
				current = current.next
				continue
			}

			x2 := current.x + current.radius
			y2 := current.y + current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			distance := xDelta*xDelta + yDelta*yDelta
			if distance < shortestDistance && distance <= unit.vision*unit.vision {
				shortestDistance = distance
				nearestUnitPos = current.position
			}

			current = current.next
		}
	}

	return nearestUnitPos, shortestDistance == math.MaxInt64
}

// checks the passed newPos hexagon and its adjacent hexagons for collisions between
// passed unit and unit midpoints within those hexagons. Returns the position
// of the unit which collided with the passed unit.
func (s *state) collidesWithUnit(newPos position, unit unit) (position, bool) {
	x1 := newPos.x + unit.radius
	y1 := newPos.y + unit.radius
	r1 := unit.radius

	gridPos, _ := getGridPos(newPos.pixelPosition, s.hexConfig)
	hex := s.hexagons[gridPos]
	var neHex, eHex, seHex, swHex, wHex, nwHex *hexagon

	if pos, ok := gridPos.adjacent(northeastSide, s.hexConfig); ok {
		neHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(eastSide, s.hexConfig); ok {
		eHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(southeastSide, s.hexConfig); ok {
		seHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(southwestSide, s.hexConfig); ok {
		swHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(westSide, s.hexConfig); ok {
		wHex = s.hexagons[pos]
	}
	if pos, ok := gridPos.adjacent(northwestSide, s.hexConfig); ok {
		nwHex = s.hexagons[pos]
	}

	hexagons := [7]*hexagon{hex, neHex, eHex, seHex, swHex, wHex, nwHex}

	for _, hex := range hexagons {
		if hex == nil {
			continue
		}

		linkedList := hex.linkedList

		if linkedList == nil || (linkedList.id == unit.id && linkedList == linkedList.next) {
			continue
		}

		end := linkedList
		current := linkedList
		// checkOccurred handles the edge case where a hexagon has a single
		// unit within in. This unit should still be checked since it
		// is not the unit passed into the function but don't self cycle
		// since it will never progress to a different current unit.
		var checkOccurred bool
		for {
			if current.id == end.id && checkOccurred {
				break
			}

			if current.id == unit.id {
				current = current.next
				continue
			}

			checkOccurred = true

			if newPos.layer != current.layer {
				current = current.next
				continue
			}

			x2 := current.x + current.radius
			y2 := current.y + current.radius
			r2 := current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			rSum := r1 + r2
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				return current.position, true
			}

			current = current.next
		}
	}

	return position{}, false
}

func (s *state) processTerminalPositions() {
	for _, u := range s.units {
		if !u.explodeAtTerminalPos {
			continue
		}

		if u.attacksPosition && u.atTerminalPosition {
			u.health = 0
			s.units[u.position] = u
			s.explode(u.position, u.explodeRadius, s.units)
		} else if u.atTerminalPosition {
			u.health = 0
			s.units[u.position] = u
		}

	}
}

func (s *state) attackUnits() {
	var unitsDeletedCount int
	var unitAttackedOtherUnitCount int

	for _, u := range s.units {
		// TODO: This may be wrong behavior but I am not sure if the issue is occurring
		// e.g. the target was killed by another unit this unit should probably find a potential new target.
		// Then attack that rather than not attack for a turn
		if u.localTarget.isZero() {
			continue
		}

		xTargetPixelPos := u.localTarget.x
		yTargetPixelPos := u.localTarget.y

		var xDelta int
		var yDelta int
		if u.x < (xTargetPixelPos) {
			xDelta = xTargetPixelPos - u.x
		} else if u.x > (xTargetPixelPos) {
			xDelta = u.x - xTargetPixelPos
		}
		if u.y < (yTargetPixelPos) {
			yDelta = yTargetPixelPos - u.y
		} else if u.y > (yTargetPixelPos) {
			yDelta = u.y - yTargetPixelPos
		}

		var nearestVisibleUnit = s.units[u.localTarget]
		var rSum = u.reach + nearestVisibleUnit.radius
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[u.localTarget] = nearestVisibleUnit
			s.units[u.position] = u
		}
	}

	for _, u := range s.units {
		if u.health <= 0 {

			gridPos, _ := getGridPos(u.pixelPosition, s.hexConfig)
			hex := s.hexagons[gridPos]
			hex.unitCount--
			if hex.linkedList.id == u.id && u.next.id != u.id {
				hex.linkedList = u.next
			} else if hex.linkedList.id == u.id && u.next.id == u.id {
				hex.linkedList = nil
			}

			u.next.prior = u.prior
			u.prior.next = u.next
			delete(s.units, u.position)
			s.unitsDeletedTotal++
			unitsDeletedCount++
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
		if u.globalTarget.isZero() {
			noTargetNodeCount++
			continue
		}

		if u.attackedThisTurn {
			attackedUnitThisTurnCount++
			continue
		}

		if !u.canReachTargetNode.isZero() {
			attackedNodeThisTurnCount++
			targetNode := s.nodes[u.canReachTargetNode.gridPos]
			targetNode.currentHealth = max(targetNode.currentHealth-u.attack, 0)
			u.attackedThisTurn = true
			if targetNode.currentHealth == 0 {
				targetNode.capturedBy = u.player
			}
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

	for _, n := range s.nodes {
		if n.currentHealth == 0 {
			newNode := newNode(n.capturedBy, n.gridPos, n.maxHealth, s.hexConfig, n.unitConfig)
			s.nodes[n.gridPos] = &newNode
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

	topLeftGridPos, _ := getGridPos(topLeftPixelPos, s.hexConfig)
	bottomRightGridPos, _ := getGridPos(bottomRightPixelPos, s.hexConfig)

	for _, n := range s.nodes {
		if (topLeftGridPos.x < n.gridPos.x) &&
			(n.gridPos.x < bottomRightGridPos.x) &&
			(topLeftGridPos.y < n.gridPos.y) &&
			(n.gridPos.y < bottomRightGridPos.y) {
			n.selected = true
		} else {
			n.selected = false
		}
		s.nodes[n.gridPos] = n
	}

	for _, u := range s.units {
		if (topLeftPixelPos.x < u.x) &&
			(u.x < bottomRightPixelPos.x) &&
			(topLeftPixelPos.y < u.y) &&
			(u.y < bottomRightPixelPos.y) {
			u.selected = true
		} else {
			u.selected = false
		}
		s.units[u.position] = u
	}

}

func (s *state) updateTargets(target position, blueComputerActive bool, redComputerActive bool) {
	for _, n := range s.nodes {
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
		s.nodes[n.gridPos] = n
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
	for _, n := range s.nodes {
		if n.player == neutralPlayer {
			continue
		}

		// TODO: determine behavior for producing units during collision
		pixelPos, _ := getPixelPosition(n.gridPos, center, s.hexConfig)
		uPos := position{n.unitConfig.produceLayer, pixelPos}
		newUnit := newUnit(n.player, n.target, n.unitConfig, uPos, s.unitsProducedTotal)
		if _, collides := s.collidesWithUnit(uPos, newUnit); !collides {
			hex := s.hexagons[n.gridPos]
			if hex.linkedList == nil && hex.unitCount != 0 {
				panic("if hex linked list is nil, unit count should be 0")
			}

			updatePointers(hex, &newUnit)

			s.units[uPos] = &newUnit
			hex.unitCount++
			s.unitsProducedTotal++
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
	for _, n := range s.nodes {
		n.drawNode(dc)
	}
	for _, n := range s.nodes {
		n.drawTargets(dc, drawBlueTargets, drawRedTargets)
	}
	for _, u := range s.units {
		if u.layer == highAltitude {
			u.draw(dc, u.selected)
		}
	}
	for _, u := range s.units {
		if u.layer == lowAltitude {
			u.draw(dc, u.selected)
		}
	}
	for _, u := range s.units {
		if u.layer == ground {
			u.draw(dc, u.selected)
		}
	}
	if drawFluxes {
		for fPos, flux := range s.fluxes {
			flux.draw(dc, fPos)
		}
	}
}
