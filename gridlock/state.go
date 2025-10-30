// High level game state and interactions between game components.

package main

import (
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
)

// not enforced but must be a singleton
type state struct {
	// TODO: The next largest optimization might be to split hexagon linked list
	// into respective layers.
	g                  grid
	oLUT               obstacleLUT
	nLUT               nodeLUT
	hexagons           hexagons
	obstacles          obstacles
	nodes              nodes
	units              units
	developerMode      bool
	unitsProducedTotal int
	unitsDeletedTotal  int
}

func newState(n nodes, o obstacles, h hexagons, oLUT obstacleLUT, nLUT nodeLUT, g grid, developerMode bool) *state {
	s := state{
		g:             g,
		obstacles:     o,
		oLUT:          oLUT,
		hexagons:      h,
		nodes:         n,
		nLUT:          nLUT,
		units:         make(units),
		developerMode: developerMode,
	}

	return &s
}

func (s *state) refresh() {
	for _, u := range s.units {
		u.refresh()
	}
}

func (s *state) analyzeTargets() {
	for _, u := range s.units {
		if u.staticTarget {
			u.localTargets = []position{u.globalTarget}
			u.prioritizeTargets()
			continue
		}

		steer := s.steer(u)
		if !u.position.equals(steer.PixelPosition) {
			u.updateTarget(steerTarget, steer)
		}

		u.canReachNode = node{}
		if s.nLUT[u.PixelPosition].ChebyshevDistance < u.reach {
			if len(s.nLUT[u.PixelPosition].NearestNodes) > 1 {
				logger.log(errorLevel, undecidedBehavior, "Undecided behavior. Unit can reach multiple nodes. Arbitrarily picking one")
			}
			p := s.nLUT[u.PixelPosition].NearestNodes[0]
			n := s.nodes[p]
			if u.player != n.player {
				u.canReachNode = *n
			}
		}

		nearest, _ := s.nearestUnitWithinVision(*u)
		u.updateTarget(enemyTarget, nearest)

		// TODO: Remove this hack once guided missile is fixed
		if u.layer == ground {
			if !u.globalTarget.isZero() && u.pathGlobalTarget {
				u.pathfindGlobalTarget(s.obstacles, s.oLUT, s.g)
			} else if len(u.localTargets) > 0 && u.atTerminalPosition {
				u.localTargets = u.localTargets[:len(u.localTargets)-1]
			}
		}

		u.prioritizeTargets()
	}
}

func (s *state) steer(u *unit) position {
	x1 := u.position.X + u.radius
	y1 := u.position.Y + u.radius

	gridPos, _ := s.g.getGridPos(u.PixelPosition)
	hex := s.hexagons[gridPos]
	var neHex, eHex, seHex, swHex, wHex, nwHex *hexagon

	if pos, ok := s.g.adjacent(gridPos, hsNortheast); ok {
		neHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsEast); ok {
		eHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSoutheast); ok {
		seHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSouthwest); ok {
		swHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsWest); ok {
		wHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsNorthwest); ok {
		nwHex = s.hexagons[pos]
	}

	hexagons := [7]*hexagon{hex, neHex, eHex, seHex, swHex, wHex, nwHex}

	var xSteer int
	var ySteer int
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

			if u.player != current.player {
				current = current.next
				continue
			}

			x2 := current.X + current.radius
			y2 := current.Y + current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			rSum := u.radius + current.radius + 1
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				magnitude := 1

				switch {
				case u.X > current.X:
					xSteer += magnitude
				case u.X < current.X:
					xSteer -= magnitude
				}

				switch {
				case u.Y > current.Y:
					ySteer += magnitude
				case u.Y < current.Y:
					ySteer -= magnitude
				}
			}

			current = current.next
		}
	}

	return position{layer: u.layer, PixelPosition: PixelPosition{X: u.X + xSteer, Y: u.Y + ySteer}}
}

func (s *state) moveUnits() {
	logger.log(infoLevel, stateLogic, "moveUnits", slog.Int("unitCount", len(s.units)))

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

		if u.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		moves := u.getMoves()
		if len(moves) == 0 {
			continue
		}

		if u.activeTarget.isZero() {
			unitsWithoutTarget++
		}

		var newPos position
		var confirmMove = func() {}
		for _, m := range moves {
			p := m.suggest

			if u.layer != ground {
				newPos = p
				confirmMove = m.confirm
				break
			}

			if s.oLUT[p.PixelPosition].UnwalkableChebyshevDistance > u.radius {
				newPos = p
				confirmMove = m.confirm
				break
			}
		}

		if newPos.isZero() {
			unitsAtUnwalkable++
			continue
		}

		collidesWithBorder := collidesWithBorder(newPos, u.radius)

		var collidedUnitPos, collidesWithUnit = s.collidesWithUnit(
			newPos,
			*u)

		switch {
		case collidesWithUnit && u.player == bluePlayer:
			u.AtCollision = true
		case collidesWithUnit && u.player == redPlayer:
			u.AtCollision = true
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
				s.units[collidedUnitPos] = s.units[position{collidedUnitPos.layer, PixelPosition{collidedUnitPos.X, collidedUnitPos.Y}}]
				u.health = 0
				s.units[u.position] = s.units[position{u.layer, PixelPosition{u.X, u.Y}}]
			}
		} else if !u.canReachNode.isZero() {
			unitsAtNode++
		} else {
			unitsMovedCount++
			s.units[newPos] = s.units[position{u.layer, PixelPosition{u.X, u.Y}}]

			oldGridPos, _ := s.g.getGridPos(u.PixelPosition)
			newGridPos, _ := s.g.getGridPos(newPos.PixelPosition)
			if !oldGridPos.equals(newGridPos) {
				oldHex := s.hexagons[oldGridPos]
				newHex := s.hexagons[newGridPos]
				oldHex.subtract(u.player)
				newHex.add(u.player)

				if u.next.id == u.id && u.prior.id == u.id {
					oldHex.linkedList = nil
				} else if u.next.id == u.id || u.prior.id == u.id {
					panic("unit is in an invalid state")
				} else {
					oldHex.linkedList = u.next
				}

				u.next.prior = u.prior
				u.prior.next = u.next

				newHex.updatePointers(u)
			}

			delete(s.units, u.position)
			confirmMove()
		}
	}

	logger.log(
		debugLevel,
		stateLogic,
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

func (s *state) explode(pos position, blastRadius int, units units) {
	x1 := pos.X
	y1 := pos.Y
	r1 := blastRadius

	for _, otherUnit := range units {
		if otherUnit.layer != ground {
			continue
		}

		x2 := otherUnit.X + otherUnit.radius
		y2 := otherUnit.Y + otherUnit.radius
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
	x1 := u.position.X + u.radius
	y1 := u.position.Y + u.radius

	gridPos, _ := s.g.getGridPos(u.PixelPosition)
	hex := s.hexagons[gridPos]
	var neHex, eHex, seHex, swHex, wHex, nwHex *hexagon

	if pos, ok := s.g.adjacent(gridPos, hsNortheast); ok {
		neHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsEast); ok {
		eHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSoutheast); ok {
		seHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSouthwest); ok {
		swHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsWest); ok {
		wHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsNorthwest); ok {
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

			x2 := current.X + current.radius
			y2 := current.Y + current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			distance := xDelta*xDelta + yDelta*yDelta
			if distance < shortestDistance && distance <= u.vision*u.vision {
				if _, los := u.lineOfSight(current.PixelPosition, s.oLUT, u.radius); los {
					shortestDistance = distance
					nearestUnitPos = current.position
				}
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
	x1 := newPos.X + unit.radius
	y1 := newPos.Y + unit.radius
	r1 := unit.radius

	gridPos, _ := s.g.getGridPos(newPos.PixelPosition)
	hex := s.hexagons[gridPos]
	var neHex, eHex, seHex, swHex, wHex, nwHex *hexagon

	if pos, ok := s.g.adjacent(gridPos, hsNortheast); ok {
		neHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsEast); ok {
		eHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSoutheast); ok {
		seHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsSouthwest); ok {
		swHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsWest); ok {
		wHex = s.hexagons[pos]
	}
	if pos, ok := s.g.adjacent(gridPos, hsNorthwest); ok {
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

			x2 := current.X + current.radius
			y2 := current.Y + current.radius
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

		xTargetPixelPos := u.enemyTarget.X
		yTargetPixelPos := u.enemyTarget.Y

		var xDelta int
		var yDelta int
		if u.X < (xTargetPixelPos) {
			xDelta = xTargetPixelPos - u.X
		} else if u.X > (xTargetPixelPos) {
			xDelta = u.X - xTargetPixelPos
		}
		if u.Y < (yTargetPixelPos) {
			yDelta = yTargetPixelPos - u.Y
		} else if u.Y > (yTargetPixelPos) {
			yDelta = u.Y - yTargetPixelPos
		}

		var nearestVisibleUnit = s.units[u.enemyTarget]
		var rSum = u.reach + nearestVisibleUnit.radius
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.units[u.enemyTarget] = nearestVisibleUnit
			s.units[u.position] = u
		}
	}

	for _, u := range s.units {
		if u.health <= 0 {

			gridPos, _ := s.g.getGridPos(u.PixelPosition)
			hex := s.hexagons[gridPos]
			hex.subtract(u.player)
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
		stateLogic,
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
		stateLogic,
		"attackNodes",
		slog.Int("no target node", noTargetNodeCount),
		slog.Int("attacked unit this turn", attackedNodeThisTurnCount),
		slog.Int("attacked node this turn", attackedNodeThisTurnCount),
		slog.Int("cannot reach node", cannotReachNodeCount),
	)
}

func (s *state) updateForNodeCapture() bool {
	var nodeOwnershipChanged bool

	for _, n := range s.nodes {
		if n.currentHealth == 0 {
			newNode := newNode(n.capturedBy, n.gridPos, n.maxHealth, n.unitConfig)
			newNode.control = s.nodes[n.gridPos].control
			newNode.singleSelect = s.nodes[n.gridPos].singleSelect
			newNode.multiSelect = s.nodes[n.gridPos].multiSelect
			s.nodes[n.gridPos] = &newNode
			nodeOwnershipChanged = true
		}
	}

	return nodeOwnershipChanged
}

func (s *state) updateSelected(clickPosition PixelPosition, releasePosition PixelPosition) {
	var priorSingleSelect GridPosition
	for _, n := range s.nodes {
		if n.singleSelect {
			priorSingleSelect = n.gridPos
		}
		n.singleSelect = false
		n.multiSelect = false
	}

	clickedGridPos, _ := s.g.getGridPos(clickPosition)
	clickedNode, nodeClicked := s.nodes[clickedGridPos]
	if nodeClicked && clickPosition.equals(releasePosition) && clickedGridPos.equals(priorSingleSelect) {
		for _, n := range s.nodes {
			if n.player == clickedNode.player {
				n.multiSelect = true
			}
		}
		return
	} else if nodeClicked && clickPosition.equals(releasePosition) {
		clickedNode.singleSelect = true
		return
	}

	var topLeftPixelPos PixelPosition
	var bottomRightPixelPos PixelPosition

	switch {
	case clickPosition.X < releasePosition.X && clickPosition.Y < releasePosition.Y:
		topLeftPixelPos = clickPosition
		bottomRightPixelPos = releasePosition
	case clickPosition.X > releasePosition.X && clickPosition.Y > releasePosition.Y:
		topLeftPixelPos = releasePosition
		bottomRightPixelPos = clickPosition
	case clickPosition.X > releasePosition.X && clickPosition.Y < releasePosition.Y:
		topLeftPixelPos = PixelPosition{releasePosition.X, clickPosition.Y}
		bottomRightPixelPos = PixelPosition{clickPosition.X, releasePosition.Y}
	case clickPosition.X < releasePosition.X && clickPosition.Y > releasePosition.Y:
		topLeftPixelPos = PixelPosition{clickPosition.X, releasePosition.Y}
		bottomRightPixelPos = PixelPosition{releasePosition.X, clickPosition.Y}
	}

	// TODO: Fix bad UI behavior when cursor ends off screen. Right now, when a cursor
	// ends off screen the behavior differs based on which side it ends. If the cursor
	// ends left or above window, the position is incorrectly set to {0, 0} which
	// usually causes too much selection to occur. The behavior is correct when
	// the cursor ends right or below window. It is tricky and error prone to fix this
	// due to hexagon LUT repetition, so I am punting it.
	topLeftGridPos, _ := s.g.getGridPos(topLeftPixelPos)
	bottomRightGridPos, _ := s.g.getGridPosOutbounds(bottomRightPixelPos)

	var firstSelectedNode *node
	for _, n := range s.nodes {
		// TODO: Fix broken code or remove it? Not sure if this is still needed.
		if (topLeftGridPos.X < n.gridPos.X) &&
			(n.gridPos.X < bottomRightGridPos.X) &&
			(topLeftGridPos.Y < n.gridPos.Y) &&
			(n.gridPos.Y < bottomRightGridPos.Y) {

			if firstSelectedNode == nil {
				firstSelectedNode = n
				firstSelectedNode.singleSelect = true
				firstSelectedNode.multiSelect = false
			} else {
				firstSelectedNode.singleSelect = false
				firstSelectedNode.multiSelect = true
				n.multiSelect = true
			}
		}
	}

	if s.developerMode {
		for _, u := range s.units {
			if (topLeftPixelPos.X < u.X) &&
				(u.X < bottomRightPixelPos.X) &&
				(topLeftPixelPos.Y < u.Y) &&
				(u.Y < bottomRightPixelPos.Y) {
				u.selected = true
			} else {
				u.selected = false
			}
		}

		for _, h := range s.hexagons {
			p := s.g.getPixelPosOutbounds(h.gridPos, hpCenter)
			if (topLeftPixelPos.X < p.X) &&
				(p.X < bottomRightPixelPos.X) &&
				(topLeftPixelPos.Y < p.Y) &&
				(p.Y < bottomRightPixelPos.Y) {
				h.selected = true
			} else {
				h.selected = false
			}
		}

		for _, h := range s.hexagons {
			p := s.g.getPixelPosOutbounds(h.gridPos, hpCenter)
			if (topLeftPixelPos.X < p.X) &&
				(p.X < bottomRightPixelPos.X) &&
				(topLeftPixelPos.Y < p.Y) &&
				(p.Y < bottomRightPixelPos.Y) {
				h.selected = true
			} else {
				h.selected = false
			}
		}

		visited := make(map[uuid.UUID]struct{})
		for _, o := range s.obstacles {
			g := *o.visibilityGrid
			obstaclePoint := g[0]
			if _, ok := visited[obstaclePoint.visibilityGridId]; ok {
				continue
			}
			visited[obstaclePoint.visibilityGridId] = struct{}{}
			for i, op := range *o.visibilityGrid {
				if (topLeftPixelPos.X < op.X) &&
					(op.X < bottomRightPixelPos.X) &&
					(topLeftPixelPos.Y < op.Y) &&
					(op.Y < bottomRightPixelPos.Y) {
					op.selected = true
				} else {
					op.selected = false
				}
				ops := *o.visibilityGrid
				ops[i] = op
			}
		}
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

		if n.selected() {
			n.target = target
		}

		if s.developerMode {
			for _, u := range s.units {
				if u.selected {
					u.updateTarget(globalTarget, target)
				}
			}
		} else if n.singleSelect {
			for _, h := range n.control {
				end := h.linkedList
				current := end
				var checkOccurred bool
				for {

					if current == nil {
						break
					}

					if current.id == end.id && checkOccurred {
						break
					}

					checkOccurred = true
					current.updateTarget(globalTarget, target)
					current = current.next
				}
			}
		} else if n.multiSelect {
			end := s.hexagons[n.gridPos].linkedList
			current := end
			var checkOccurred bool
			for {
				if current == nil {
					break
				}

				if current.id == end.id && checkOccurred {
					break
				}

				checkOccurred = true

				if !current.movedThisTurn {
					current.updateTarget(globalTarget, target)
				}

				current = current.next
			}
		}
	}
}

func (s *state) produceUnits() {
	for _, n := range s.nodes {
		if n.player == neutralPlayer {
			continue
		}

		// TODO: determine behavior for producing units during collision
		pixelPos, ok := s.g.getPixelPos(n.gridPos, hpCenter)
		if !ok {
			panic("could not get pixel pos for produce units")
		}
		uPos := position{n.unitConfig.produceLayer, pixelPos}
		newUnit := newUnit(n.player, n.target, n.unitConfig, uPos, s.unitsProducedTotal, n.selected() && s.developerMode)
		if _, collides := s.collidesWithUnit(uPos, newUnit); !collides {
			hex := s.hexagons[n.gridPos]
			if hex.linkedList == nil && hex.unitCount != 0 {
				panic("if hex linked list is nil, unit count should be 0")
			}

			hex.updatePointers(&newUnit)

			s.units[uPos] = &newUnit
			hex.add(newUnit.player)
			s.unitsProducedTotal++
		}
	}
}

func (s *state) draw(dc drawConfig) {

	for _, n := range s.nodes {
		n.drawNode(dc)
	}

	for _, o := range s.obstacles {
		o.drawObstacle(dc)
	}

	// TODO: Store units by layer for better performance
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
}
