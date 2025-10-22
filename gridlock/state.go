// High level game state and interactions between game components.

package main

import (
	"context"
	"log/slog"
	"math"
	"os"
	"os/signal"

	"github.com/gofrs/uuid/v5"
)

// not enforced but must be a singleton
type state struct {
	*hexConfig
	// TODO: The next largest optimization might be to split hexagon linked list
	// into respective layers.
	pixelPath                   pixelPath
	hexagons                    hexagons
	obstacles                   obstacles
	nodes                       nodes
	units                       units
	developerMode               bool
	unitsProducedTotal          int
	unitsDeletedTotal           int
	unitPathfind                heatmap
	unitPathfindHeatmapComplete bool
}

func newState(n nodes, o obstacles, h hexagons, pixelPath pixelPath, hexConfig *hexConfig, developerMode bool) *state {
	s := state{
		hexConfig:     hexConfig,
		obstacles:     o,
		pixelPath:     pixelPath,
		hexagons:      h,
		nodes:         n,
		units:         make(units),
		developerMode: developerMode,
		unitPathfind:  make(heatmap),
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
			continue
		}

		steer := s.steer(u)
		if !u.position.equals(steer.pixelPosition) {
			u.updateTarget(steerTarget, steer)
		}

		u.reachesNodes(s.nodes, s.hexConfig)

		nearest, _ := s.nearestUnitWithinVision(*u)
		u.updateTarget(enemyTarget, nearest)

		if !u.globalTarget.isZero() && u.pathGlobalTarget {
			u.pathfindGlobalTarget(s.obstacles, s.pixelPath, s.hexConfig)
		} else if len(u.localTargets) > 0 && u.atTerminalPosition {
			u.localTargets = u.localTargets[:len(u.localTargets)-1]
		}

		u.prioritizeTargets()
	}
}

func (s *state) steer(u *unit) position {
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

			x2 := current.x + current.radius
			y2 := current.y + current.radius

			xDelta := x2 - x1
			yDelta := y2 - y1
			rSum := u.radius + current.radius + 1
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				magnitude := 1

				switch {
				case u.x > current.x:
					xSteer += magnitude
				case u.x < current.x:
					xSteer -= magnitude
				}

				switch {
				case u.y > current.y:
					ySteer += magnitude
				case u.y < current.y:
					ySteer -= magnitude
				}
			}

			current = current.next
		}
	}

	return position{layer: u.layer, pixelPosition: pixelPosition{x: u.x + xSteer, y: u.y + ySteer}}
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

		newPos, completeMove := u.getNewPos()

		if u.activeTarget.isZero() {
			unitsWithoutTarget++
		}

		if u.equals(newPos.pixelPosition) {
			unitsAtFinalPosition++
			continue
		}

		if !s.pixelPath[newPos.pixelPosition].walkable || s.pixelPath[newPos.pixelPosition].unwalkableChebyshevDistance < u.radius+1 {
			unitsAtUnwalkable++
			u.AtUnwalkable = true
			continue
		}

		collidesWithBorder := collidesWithBorder(newPos, u.radius, s.hexConfig)

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
			if !oldGridPos.equals(newGridPos) {
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

				newHex.updatePointers(u)
			}

			delete(s.units, u.position)
			completeMove()
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
				if _, los := u.lineOfSight(current.pixelPosition, s.pixelPath, u.radius); los {
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

func (s *state) updateNodeOwnership() bool {
	var nodeOwnershipChanged bool

	for _, n := range s.nodes {
		if n.currentHealth == 0 {
			newNode := newNode(n.capturedBy, n.gridPos, n.maxHealth, s.hexConfig, n.unitConfig)
			newNode.control = s.nodes[n.gridPos].control
			newNode.singleSelect = s.nodes[n.gridPos].singleSelect
			newNode.multiSelect = s.nodes[n.gridPos].multiSelect
			s.nodes[n.gridPos] = &newNode
			nodeOwnershipChanged = true
		}
	}

	return nodeOwnershipChanged
}

func (s *state) updateSelected(clickPosition pixelPosition, releasePosition pixelPosition) {
	var priorSingleSelect gridPosition
	for _, n := range s.nodes {
		if n.singleSelect {
			priorSingleSelect = n.gridPos
		}
		n.singleSelect = false
		n.multiSelect = false
	}

	clickedGridPos, _ := clickPosition.getGridPos(s.hexConfig)
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

	// TODO: Fix bad UI behavior when cursor ends off screen. Right now, when a cursor
	// ends off screen the behavior differs based on which side it ends. If the cursor
	// ends left or above window, the position is incorrectly set to {0, 0} which
	// usually causes too much selection to occur. The behavior is correct when
	// the cursor ends right or below window. It is tricky and error prone to fix this
	// due to hexagon lookup repetition, so I am punting it.
	topLeftGridPos, _ := topLeftPixelPos.getGridPos(s.hexConfig)
	bottomRightGridPos, _ := bottomRightPixelPos.getGridPosOutbounds(s.hexConfig)

	var firstSelectedNode *node
	for _, n := range s.nodes {
		if (topLeftGridPos.x < n.gridPos.x) &&
			(n.gridPos.x < bottomRightGridPos.x) &&
			(topLeftGridPos.y < n.gridPos.y) &&
			(n.gridPos.y < bottomRightGridPos.y) {

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
		pixelPos, ok := n.gridPos.getPixelPos(center, s.hexConfig)
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
			hex.unitCount++
			s.unitsProducedTotal++
		}
	}
}

func (s *state) createUnitPathfindHeatmap() {
	logger.onlyEnable(devTool)

	if len(s.unitPathfind) != 0 {
		logger.log(infoLevel, devTool, "cannot create heatmap since it was already attempted")
		return
	}

	logger.log(infoLevel, devTool, "start creating pathfinding heatmap")

	var unit *unit
	for _, u := range s.units {
		if u.selected && unit != nil {
			logger.log(infoLevel, devTool, "can only select one unit when creating heatmap")
			return
		} else if u.selected {
			unit = u
		}
	}

	if unit == nil {
		logger.log(infoLevel, devTool, "must select one unit when creating heat map")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pixelCount := 0
	for pixelPos, path := range s.pixelPath {
		select {
		case <-ctx.Done():
			return
		default:
			pixelCount++
			if pixelCount%(s.windowWidth*s.windowHeight/100) == 0 {
				percentProgess := pixelCount / (s.windowWidth * s.windowHeight / 100)
				logger.log(infoLevel, devTool, "pathfind heatmap", slog.Int("pixel positions percent progress", percentProgess))
			}

			if path.walkable {
				unit.updateTarget(globalTarget, position{layer: unit.layer, pixelPosition: pixelPos})
				distance := unit.pathfindGlobalTarget(s.obstacles, s.pixelPath, s.hexConfig)
				s.unitPathfind[pixelPos] = distance
			}
		}
	}

	s.unitPathfindHeatmapComplete = true
	logger.log(infoLevel, devTool, "done creating pathfinding heatmap")
}

func (s *state) draw(
	dc drawConfig,
	drawBlueTargets bool,
	drawRedTargets,
	drawVisibilityGrids bool,
	drawPixelPath bool) {

	for _, n := range s.nodes {
		n.drawNode(dc)
	}
	for _, n := range s.nodes {
		n.drawTargets(dc, drawBlueTargets, drawRedTargets)
	}

	visited := make(map[uuid.UUID]struct{})
	for _, o := range s.obstacles {
		o.drawObstacle(dc)

		if drawVisibilityGrids {
			g := *o.visibilityGrid
			op := g[0]
			if _, ok := visited[op.visibilityGridId]; ok {
				continue
			}
			visited[op.visibilityGridId] = struct{}{}
			o.drawVisibilityGrid(dc)
		}
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

	if drawPixelPath {
		s.pixelPath.draw(dc)
	}

	if s.unitPathfindHeatmapComplete {
		s.unitPathfind.draw(dc)
	}
}
