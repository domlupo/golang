// Currently contains most game state and game logic. May need to split in future.
package main

import (
	"log/slog"
	"math"
)

// not enforced but must be a singleton
type state struct {
	*hexConfig
	// TODO: The next largest optimization might be to split hexagon linked list
	// into respective layers.
	pixelPathing       pixelPathing
	hexagons           map[gridPosition]*hexagon
	obstacles          map[gridPosition]*obstacle
	nodes              map[gridPosition]*node
	units              map[position]*unit
	enableFlux         bool
	unitsProducedTotal int
	unitsDeletedTotal  int
}

func newState(nodes map[gridPosition]*node, obstacles map[gridPosition]*obstacle, pixelPathing pixelPathing, hexConfig *hexConfig, enableFlux bool) *state {
	s := state{
		hexConfig:    hexConfig,
		obstacles:    obstacles,
		pixelPathing: pixelPathing,
		hexagons:     setupHexagons(hexConfig),
		nodes:        nodes,
		units:        make(map[position]*unit),
		enableFlux:   enableFlux,
	}

	return &s
}

func (s *state) refresh() {
	for _, u := range s.units {
		u.refresh()
	}

	for _, h := range s.hexagons {
		h.refresh(s.hexConfig)
	}
}

func (s *state) analyzeTargets() {
	for _, u := range s.units {
		if u.staticTarget {
			continue
		}

		uGridPos, uGridPosOk := u.getGridPos(s.hexConfig)
		if !uGridPosOk {
			panic("unit grid position must exist")
		}
		hex := s.hexagons[uGridPos]
		if flux, fluxOk := u.getFlux(hex, s.hexConfig); fluxOk {
			u.updateTarget(fluxTarget, flux)
		}

		u.reachesNodes(s.nodes, s.hexConfig)

		nearest, _ := s.nearestUnitWithinVision(*u)
		u.updateTarget(enemyTarget, nearest)

		if !u.globalTarget.isZero() && u.pathGlobalTarget {
			u.pathfindGlobalTarget(s.obstacles, s.pixelPathing, s.hexConfig)
		} else if len(u.localTargets) > 0 && u.atTerminalPosition {
			u.localTargets = u.localTargets[:len(u.localTargets)-1]
		}

		u.prioritizeTargets()
	}
}

func (s *state) analyzeMoves() {
	for _, h := range s.hexagons {
		if h.unitCount != 0 {
			h.avgUnit = pixelPosition{h.sumUnitX / h.unitCount, h.sumUnitY / h.unitCount}
		}
	}
}

func (s *state) moveUnits() {
	logger.log(infoLevel, stateLogType, "moveUnits", slog.Int("unitCount", len(s.units)))

	var unitsMovedCount int
	var unitsAlreadyAttacked int
	var unitsWithoutTarget int
	var unitsAtUnwalkable int
	var unitsAtNode int
	var unitsAtFinalPosition int
	var unitsCollideCount int
	var unitsBorderCollideCount int

	for _, u := range s.units {
		if u.movedThisTurn {
			continue
		}

		gridPos, _ := u.getGridPos(s.hexConfig)
		hex := s.hexagons[gridPos]

		// TODO: fix this isn't accurate since it assumes the unit does not move
		hex.sumUnitX += u.x
		hex.sumUnitY += u.y

		if u.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		newPos, completeMove := u.getNewPos()

		if u.currentTarget.isZero() {
			unitsWithoutTarget++
		}

		if u.x == newPos.x && u.y == newPos.y {
			unitsAtFinalPosition++
			continue
		}

		if !s.pixelPathing[newPos.pixelPosition].walkable {
			unitsAtUnwalkable++
			u.AtUnwalkable = true
			switch u.player {
			case bluePlayer:
				hex.blueUnitsUnwalkable++
			case redPlayer:
				hex.redUnitsUnwalkable++
			}
			continue
		}

		collidesWithBorder := s.collidesWithBorder(newPos, u.radius)

		var collidedUnitPos, collidesWithUnit = s.collidesWithUnit(
			newPos,
			*u)

		switch {
		case collidesWithUnit && u.player == bluePlayer:
			u.AtCollision = true
			hex.blueUnitsColliding++
		case collidesWithUnit && u.player == redPlayer:
			u.AtCollision = true
			hex.redUnitsColliding++
		}

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
		} else if !u.canReachNode.isZero() {
			unitsAtNode++
		} else {
			unitsMovedCount++
			s.units[newPos] = s.units[position{u.layer, pixelPosition{u.x, u.y}}]

			oldGridPos, _ := u.pixelPosition.getGridPos(s.hexConfig)
			newGridPos, _ := newPos.pixelPosition.getGridPos(s.hexConfig)
			if oldGridPos.x != newGridPos.x || oldGridPos.y != newGridPos.y {
				oldHex := s.hexagons[oldGridPos]
				newHex := s.hexagons[newGridPos]
				oldHex.unitCount--
				// TODO: This can panic, debug why.
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
		slog.Int("units moved", unitsMovedCount),
		slog.Int("units without target", unitsWithoutTarget),
		slog.Int("units at unwalkable", unitsAtUnwalkable),
		slog.Int("units at node", unitsAtNode),
		slog.Int("units collided", unitsCollideCount),
		slog.Int("untis already attacked", unitsAlreadyAttacked),
		slog.Int("units at final position", unitsAtFinalPosition),
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

func (s *state) nearestUnitWithinVision(u unit) (position, bool) {
	x1 := u.position.x + u.radius
	y1 := u.position.y + u.radius

	gridPos, _ := u.pixelPosition.getGridPos(s.hexConfig)
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

		if linkedList == nil || (linkedList.id == u.id && linkedList == linkedList.next) {
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

			if current.id == u.id {
				current = current.next
				continue
			}

			checkOccurred = true

			if u.player == current.player {
				current = current.next
				continue
			}

			if !u.attacksGround && current.layer == ground ||
				!u.attacksLowAltitude && current.layer == lowAltitude ||
				!u.attacksHighAltitude && current.layer == highAltitude {
				current = current.next
				continue
			}

			x2 := current.x + current.radius
			y2 := current.y + current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			distance := xDelta*xDelta + yDelta*yDelta
			if distance < shortestDistance && distance <= u.vision*u.vision {
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

	gridPos, _ := newPos.pixelPosition.getGridPos(s.hexConfig)
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
		if u.enemyTarget.isZero() {
			continue
		}

		xTargetPixelPos := u.enemyTarget.x
		yTargetPixelPos := u.enemyTarget.y

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

		var nearestVisibleUnit = s.units[u.enemyTarget]
		var rSum = u.reach + nearestVisibleUnit.radius
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[u.enemyTarget] = nearestVisibleUnit
			s.units[u.position] = u

			gridPos, ok := u.pixelPosition.getGridPos(s.hexConfig)
			if !ok {
				panic("expected grid position to exist")
			}
			hex := s.hexagons[gridPos]
			switch u.player {
			case bluePlayer:
				hex.blueUnitsAttacking++
			case redPlayer:
				hex.redUnitsAttacking++
			default:
				panic("expected blue or red player")
			}
		}
	}

	for _, u := range s.units {
		if u.health <= 0 {

			gridPos, _ := u.pixelPosition.getGridPos(s.hexConfig)
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
		slog.Int("units attacked", unitAttackedOtherUnitCount),
		slog.Int("units deleted", unitsDeletedCount),
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

		if !u.canReachNode.isZero() {
			attackedNodeThisTurnCount++
			n := s.nodes[u.canReachNode.gridPos]
			n.currentHealth = max(n.currentHealth-u.attack, 0)
			u.attackedThisTurn = true
			if n.currentHealth == 0 {
				n.capturedBy = u.player
			}
		} else {
			cannotReachNodeCount++
		}
	}

	logger.log(
		debugLevel,
		stateLogType,
		"attackNodes",
		slog.Int("no target node", noTargetNodeCount),
		slog.Int("attacked unit this turn", attackedNodeThisTurnCount),
		slog.Int("attacked node this turn", attackedNodeThisTurnCount),
		slog.Int("cannot reach node", cannotReachNodeCount),
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

	topLeftGridPos, _ := topLeftPixelPos.getGridPos(s.hexConfig)
	bottomRightGridPos, _ := bottomRightPixelPos.getGridPosOutbounds(s.hexConfig)

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
		pixelPos, ok := n.gridPos.getPixelPos(center, s.hexConfig)
		if !ok {
			panic("could not get pixel pos for produce units")
		}
		uPos := position{n.unitConfig.produceLayer, pixelPos}
		newUnit := newUnit(n.player, n.target, n.unitConfig, uPos, s.unitsProducedTotal, n.selected)
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

func (s *state) draw(dc drawConfig, drawBlueTargets bool, drawRedTargets, drawFluxes bool, drawVisibilityGrids bool, drawPixelPathing bool) {
	for _, o := range s.obstacles {
		o.drawObstacle(dc)
		if drawVisibilityGrids {
			o.drawVisibilityGrid(dc)
		}
	}

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
		for _, h := range s.hexagons {
			h.drawFlux(dc)
		}
	}
	if drawPixelPathing {
		s.pixelPathing.draw(dc)
	}
}
