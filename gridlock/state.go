// High level game state and interactions between game components.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

type stateConfig struct {
	moveUnits                bool
	produceUnits             bool
	produceUnitsStepInterval uint64
	spawnUnitRadiusLimit     int
}

func newStateConfig(g grid) stateConfig {
	c := stateConfig{
		moveUnits:                true,
		produceUnits:             true,
		produceUnitsStepInterval: 180,
		spawnUnitRadiusLimit:     g.minimumDistanceBetweenNonAdjacentHexagons(),
	}
	return c
}

// A request to modify state.
type request struct {
	spawn    *unit
	complete bool
}

// The type of analysis which will process the partitions.
type partitionAnalysis int

const (
	paNone partitionAnalysis = iota
	paCollision
	paVision
	paSteer
)

func (pa partitionAnalysis) String() string {
	switch pa {
	case paNone:
		return "paNone"
	case paCollision:
		return "paCollision"
	case paVision:
		return "paVision"
	case paSteer:
		return "paSteer"
	default:
		panic(fmt.Sprintf("Cannot convert partition analysis to string: %d", pa))
	}
}

func (pa partitionAnalysis) log() slog.Attr {
	return slog.String("partitionAnalysis", pa.String())
}

// not enforced but must be a singleton
type state struct {
	c                stateConfig
	g                grid
	oLUT             obstacleLUT
	nLUT             nodeLUT
	hexagons         hexagons
	obstacles        obstacles
	nodes            nodes
	landUnits        units
	ordanceUnits     units
	aircraftUnits    units
	ballisticUnits   units
	bombGravityUnits units
	chaffUnits       units
	reactiveUnits    map[id]*unit
	requesterUnits   map[id]*unit
	step             uint64
	developerMode    bool
	unitSpawnTotal   int
	unitDespawnTotal int
}

func newState(n nodes, o obstacles, h hexagons, oLUT obstacleLUT, nLUT nodeLUT, g grid, developerMode bool) *state {
	s := state{
		c:                newStateConfig(g),
		g:                g,
		obstacles:        o,
		oLUT:             oLUT,
		hexagons:         h,
		nodes:            n,
		nLUT:             nLUT,
		landUnits:        make(units),
		ordanceUnits:     make(units),
		aircraftUnits:    make(units),
		ballisticUnits:   make(units),
		bombGravityUnits: make(units),
		chaffUnits:       make(units),
		reactiveUnits:    make(map[id]*unit),
		requesterUnits:   make(map[id]*unit),
		developerMode:    developerMode,
	}

	return &s
}

// Progress game state by a single step
func (s *state) tick() {
	logger.log(infoLevel, stateLogic, "tick: start", slog.Uint64("step", s.step))
	logger.log(infoLevel, stateLogic, "tick: refresh")
	s.refresh()
	logger.log(infoLevel, stateLogic, "tick: analyzeTargets")
	s.analyzeTargets()
	logger.log(infoLevel, stateLogic, "tick: processReactions")
	s.processReactions()
	logger.log(infoLevel, stateLogic, "tick: processRequests")
	s.processRequests()
	logger.log(infoLevel, stateLogic, "tick: processTerminalPositions")
	s.processTerminalUnits()
	logger.log(infoLevel, stateLogic, "tick: attackUnits")
	s.attackUnits()
	logger.log(infoLevel, stateLogic, "tick: attackNodes")
	s.attackNodes()
	logger.log(infoLevel, stateLogic, "tick: moveUnits")
	if s.c.moveUnits {
		s.moveUnits()
	}
	logger.log(infoLevel, stateLogic, "tick: despawnUnits")
	s.despawnUnits()
	logger.log(infoLevel, stateLogic, "tick: produceUnits")
	if s.c.produceUnits && s.step%s.c.produceUnitsStepInterval == 0 {
		s.produceUnits()
	}
	s.step++
}

func (s state) units() func(yield func(u *unit) bool) {
	return func(yield func(u *unit) bool) {
		for _, u := range s.landUnits {
			if !yield(u) {
				return
			}
		}
		for u := range s.unitsAir() {
			if !yield(u) {
				return
			}
		}
	}
}

func (s state) unitsAir() func(yield func(u *unit) bool) {
	return func(yield func(u *unit) bool) {
		for _, u := range s.ordanceUnits {
			if !yield(u) {
				return
			}
		}
		for _, u := range s.aircraftUnits {
			if !yield(u) {
				return
			}
		}
		for _, u := range s.ballisticUnits {
			if !yield(u) {
				return
			}
		}
		for _, u := range s.bombGravityUnits {
			if !yield(u) {
				return
			}
		}
		for _, u := range s.chaffUnits {
			if !yield(u) {
				return
			}
		}
	}
}

func (s state) unitsLayer(l layer) *units {
	var units units
	switch l {
	case lLand:
		units = s.landUnits
	case lOrdnance:
		units = s.ordanceUnits
	case lAircraft:
		units = s.aircraftUnits
	case lBallistic:
		units = s.ballisticUnits
	case lBombGravity:
		units = s.bombGravityUnits
	case lChaff:
		units = s.chaffUnits
	default:
		panic(fmt.Sprintf("layerUnits: unknown layer %d", l))
	}
	return &units
}

func (s *state) spawn(u *unit) {
	// TODO: Verify if this should be >= rather than more conservative >
	if u.radius > s.c.spawnUnitRadiusLimit && !u.skipSpawnRadiusCheck {
		panic(fmt.Sprintf("spawn: cannot spawn unit radius %d exceeds spawn radius limit %d", u.radius, s.c.spawnUnitRadiusLimit))
	}
	p, ok := s.g.getGridPos(u.PixelPosition)
	if !ok {
		panic(fmt.Sprintf("spawnUnit: cannot get get grid position %s", u.PixelPosition.log()))
	}
	hex := s.hexagons[p]
	hex.add(u)

	if u.reactive {
		s.reactiveUnits[u.id] = u
	}
	if u.requester {
		s.requesterUnits[u.id] = u
	}

	units := *s.unitsLayer(u.layer)
	if _, ok := units[u.PixelPosition]; !ok {
		units[u.PixelPosition] = u
	} else {
		panic(fmt.Sprintf("Cannot spawn unit at %s since it is already occupied.", u.PixelPosition.log()))
	}

	if u.terminalStepDelay != 0 {
		u.terminalStep = s.step + u.terminalStepDelay
	} else {
		u.terminalStep = math.MaxUint64
	}
	u.spawnOrder = s.unitSpawnTotal
	s.unitSpawnTotal++
}

func (s *state) spawnRequest(r request) bool {
	if r.spawn.terminalStepDelay != 0 {
		r.spawn.terminalStep = s.step + r.spawn.terminalStepDelay
	}
	units := *s.unitsLayer(r.spawn.layer)
	if _, ok := units[r.spawn.PixelPosition]; ok {
		return false
	}
	s.spawn(r.spawn)
	return true
}

func (s *state) moveUnit(u *unit, m move) {
	old, oOk := s.g.getGridPos(u.PixelPosition)
	new, nOk := s.g.getGridPos(m.PixelPosition)
	if !oOk || !nOk {
		panic(fmt.Sprintf("moveUnit: cannot move unit %s %s", u.PixelPosition.log(), m.PixelPosition.log()))
	}
	if !old.equals(new) {
		oldHex := s.hexagons[old]
		oldHex.subtract(*u)
		newHex := s.hexagons[new]
		newHex.add(u)
	}
	s.deleteUnit(u.position)
	m.confirm()
	units := *s.unitsLayer(m.layer)
	if o, ok := units[u.PixelPosition]; ok {
		panic(fmt.Sprintf("moveUnit: position for unit occupied by other unit. unit: %s other: %s", u.log(), o.log()))
	}
	units[u.PixelPosition] = u
}

func (s *state) deleteUnit(p position) {
	var units units
	switch p.layer {
	case lLand:
		units = s.landUnits
	case lOrdnance:
		units = s.ordanceUnits
	case lAircraft:
		units = s.aircraftUnits
	case lBallistic:
		units = s.ballisticUnits
	case lBombGravity:
		units = s.bombGravityUnits
	case lChaff:
		units = s.chaffUnits
	default:
		panic(fmt.Sprintf("deleteUnit: unknown layer %d", p.layer))
	}
	delete(units, p.PixelPosition)
}

func (s *state) refresh() {
	for u := range s.units() {
		u.refresh()
	}
}

func (s *state) processReactions() {
	logger.log(infoLevel, stateLogic, "processReactions", slog.Int("reactiveUnitCount", len(s.reactiveUnits)))
	for _, u := range s.reactiveUnits {
		for i, r := range u.reactions {
			if r.next > s.step {
				continue
			}

			if r.when(u) {
				r.do(u)
				r.next = s.step + r.delay
				u.reactions[i] = r
			}
		}
	}
}

func (s *state) processRequests() {
	logger.log(infoLevel, stateLogic, "processRequests", slog.Int("requesterUnitCounter", len(s.requesterUnits)))
	for _, u := range s.requesterUnits {
		for i, r := range u.requests {
			if s.spawnRequest(r) {
				r.complete = true
				u.requests[i] = r
			}
		}
	}

	for _, u := range s.requesterUnits {
		incompleteReqs := []request{}
		for _, r := range u.requests {
			if !r.complete {
				incompleteReqs = append(incompleteReqs, r)
			}
		}
		u.requests = incompleteReqs
	}
}

func (s *state) analyzeTargets() {
	for u := range s.units() {
		if u.name == chaff || u.name == bombGravity {
			continue
		}

		if u.name == missileHeatSeeking || u.name == missileAntiBallistic || u.name == missileRadarActive {
			s.nearestTargetWithinVision(u)
			u.updateTargetUnit(enemyTarget, u.nearest)
		}

		if u.name == missileHeatSeeking ||
			u.name == flare ||
			u.name == missileCruise ||
			u.name == missileBallistic ||
			u.name == missileAntiBallistic ||
			u.name == missileRadarActive {
			u.localTargets = []position{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		if u.name == fighter {
			s.nearestTargetWithinVision(u)
			if u.moveDone && len(u.globalTargets) > 1 {
				if len(u.actions) > 0 && u.actions[0].pos.equals(u.globalTargets[0]) {
					u.actions[0].do()
					u.actions = u.actions[1:]
				}
				u.globalTargets = u.globalTargets[1:]
			} else if u.moveDone && len(u.globalTargets) == 1 {
				u.globalTargets = s.nodes[u.sourceNode].targets
				u.actions = s.nodes[u.sourceNode].actions
			}
			u.localTargets = []position{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		if u.name == bomber {
			s.nearestTargetWithinVision(u)
			if u.moveDone && len(u.globalTargets) > 1 {
				if len(u.actions) > 0 && u.actions[0].pos.equals(u.globalTargets[0]) {
					u.actions[0].do()
					u.actions = u.actions[1:]
				}
				u.globalTargets = u.globalTargets[1:]
			} else if u.moveDone && len(u.globalTargets) == 1 {
				u.globalTargets = s.nodes[u.sourceNode].targets
				u.actions = s.nodes[u.sourceNode].actions
			}
			u.localTargets = []position{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		steer := s.steer(u)
		if !u.position.equals(steer) {
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

		s.nearestTargetWithinVision(u)
		u.updateTargetUnit(enemyTarget, u.nearest)

		if len(u.globalTargets) > 0 && u.pathGlobalTarget {
			u.pathfindGlobalTarget(s.obstacles, s.oLUT, s.g)
		} else if len(u.localTargets) > 0 && u.moveDone {
			u.localTargets = u.localTargets[:len(u.localTargets)-1]
		}

		u.prioritizeTargets()
	}
}

func (s *state) steer(u *unit) position {
	partitions := s.partitions(*u, paSteer)

	var xSteer int
	var ySteer int
	for _, partition := range partitions {
		if partition == nil {
			continue
		}

		linkedList := partition.get(u.layer).linkedList

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

			// TODO: Improve steering so I don't have to do weird increments here? Or are these
			// increments fine since units should steer before colliding? I also enabled infantry
			// to survive collisions which seems correct long term but that may be wrong.
			// Also, a similar but not equivalent collision check is done during moveUnits. This
			// logic can maybe be moved there to save CPU.
			// TODO: I think technically there should be a check on collides so the second radius
			// is never greater than state.spawnUnitRadiusLimit. This should be a minor issue until later.
			if _, ok := collides(u.PixelPosition, current.PixelPosition, u.radius+1, current.radius+1); ok {
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
	var unitsMovedCount int
	var unitsNotMovedCount int
	var unitsAlreadyAttacked int
	var unitsAtNode int
	var unitsAtFinalPosition int
	var unitsCollideCount int
	var unitsBorderCollideCount int

	for u := range s.units() {
		if u.movedThisTurn {
			continue
		}

		if u.attackedThisTurn {
			unitsAlreadyAttacked++
			continue
		}

		var move move
		moves := u.getMoves()
		for _, m := range moves {
			if u.layer != lLand {
				move = m
				break
			}
			if s.oLUT[m.PixelPosition].UnwalkableChebyshevDistance > u.radius {
				move = m
				break
			}
		}

		var collidedBorder bool
		var collidedUnitPos position
		var collidesWithUnit bool
		if move.equals(position{}) {
			// TODO: Should this be optimized to only perform this check when the unit
			// does not typically move e.g. chaff
			collidedBorder = collidesWithBorder(u.position, u.radius)
			collidedUnitPos, collidesWithUnit = s.collidesWithUnit(u.position, *u)
		} else {
			collidedBorder = collidesWithBorder(move.position, u.radius)
			collidedUnitPos, collidesWithUnit = s.collidesWithUnit(move.position, *u)
		}

		if collidedBorder && collidesWithUnit {
			unitsBorderCollideCount++
			unitsCollideCount++
		} else if collidedBorder {
			unitsBorderCollideCount++
		} else if collidesWithUnit {
			unitsCollideCount++
			otherUnits := *s.unitsLayer(collidedUnitPos.layer)
			otherUnit := otherUnits[collidedUnitPos.PixelPosition]
			otherUnit.collidedThisturn = true
			u.collidedThisturn = true
		} else if !u.canReachNode.isZero() {
			unitsAtNode++
		} else if move.equals(position{}) {
			unitsNotMovedCount++
		} else {
			unitsMovedCount++
			s.moveUnit(u, move)
		}
	}

	logger.log(
		infoLevel,
		stateLogic,
		"moveUnits",
		slog.Int("units moved", unitsMovedCount),
		slog.Int("units not moved", unitsNotMovedCount),
		slog.Int("units at node", unitsAtNode),
		slog.Int("units collided", unitsCollideCount),
		slog.Int("untis already attacked", unitsAlreadyAttacked),
		slog.Int("units at final position", unitsAtFinalPosition),
	)
}

func (s *state) explode(u *unit) {
	u.health = 0
	x1 := u.X
	y1 := u.Y
	r1 := u.explodeRadius

	for _, u := range s.landUnits {
		x2 := u.X + u.radius
		y2 := u.Y + u.radius
		r2 := u.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := r1 + r2
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			u.health = 0
			s.landUnits[u.PixelPosition] = u
		}
	}
}

func (s state) partitions(u unit, pa partitionAnalysis) []*hexagon {
	var partitions []GridPosition
	switch {
	case u.name == infantry && pa == paCollision:
		partitions = s.g.partitions(u.PixelPosition, u.radius+s.c.spawnUnitRadiusLimit)
	case u.name == infantry && pa == paVision:
		partitions = s.g.partitions(u.PixelPosition, u.vision+s.c.spawnUnitRadiusLimit)
	case u.name == infantry && pa == paSteer:
		partitions = s.g.partitions(u.PixelPosition, s.c.spawnUnitRadiusLimit)
	case pa == paCollision:
		partitions = s.g.partitions(u.PixelPosition, u.radius+s.c.spawnUnitRadiusLimit)
	case pa == paVision:
		partitions = s.g.partitions(u.PixelPosition, u.vision+s.c.spawnUnitRadiusLimit)
	default:
		panic(fmt.Sprintf("partitions: cannot get partitions for unit %s and analysis %s", u.name, pa.String()))
	}
	hexs := []*hexagon{}
	for _, p := range partitions {
		hexs = append(hexs, s.hexagons[p])
	}
	return hexs
}

func (s state) nearestTargetWithinVision(u *unit) {
	partitions := s.partitions(*u, paVision)
	d := math.MaxInt64
	var t *unit
	for _, partition := range partitions {
		if partition == nil {
			continue
		}

		for _, l := range u.targetLayers {
			linkedList := partition.get(l).linkedList
			if target, distance := s.nearestTargetWithinLayer(*u, linkedList); distance < d {
				t = target
				d = distance
			}
		}
	}
	u.nearest = t
}

func (s state) nearestTargetWithinLayer(u unit, linkedList *unit) (*unit, int) {
	d := math.MaxInt64
	var t *unit
	if linkedList == nil || (linkedList.id == u.id && linkedList == linkedList.next) {
		return t, d
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

		if u.player == current.player ||
			u.ignoreTarget != nil && u.ignoreTarget(current) ||
			u.ignore != nil && u.ignore(&u, current) {
			current = current.next
			continue
		}

		squaredDistance, collides := collides(u.PixelPosition, current.PixelPosition, u.vision, current.radius)
		if squaredDistance < d && collides {
			if _, los := u.lineOfSight(current.PixelPosition, s.oLUT, u.radius); los {
				d = squaredDistance
				t = current
			}
			if current.layer != lLand {
				d = squaredDistance
				t = current
			}
		}

		current = current.next
	}

	return t, d
}

func (s *state) collidesWithUnit(p position, u unit) (position, bool) {
	partitions := s.partitions(u, paCollision)
	for _, partition := range partitions {
		if partition == nil {
			continue
		}

		for _, l := range u.collisionLayers {
			linkedList := partition.get(l).linkedList
			if p, collides := s.collidesWithUnitLayer(u, p, linkedList); collides {
				return p, collides
			}
		}
	}

	return position{}, false
}

func (s state) collidesWithUnitLayer(u unit, p position, linkedList *unit) (position, bool) {
	if linkedList == nil || (linkedList.id == u.id && linkedList == linkedList.next) {
		return position{}, false
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

		if u.ignoreCollision != nil && u.ignoreCollision(current) && u.layer == current.layer {
			panic(fmt.Sprintf("collidesWithUnitLayer: Cannot ignore collisions between units in same layer. unit: %s, current: %s",
				u.log(), current.log()))
		}

		if u.ignoreCollision != nil && u.ignoreCollision(current) ||
			current.id == u.sourceUnit ||
			current.sourceUnit == u.id {
			current = current.next
			continue
		}

		if _, ok := collides(p.PixelPosition, current.PixelPosition, u.radius, current.radius); ok {
			return current.position, true
		}
		current = current.next
	}
	return position{}, false
}
func (s *state) processTerminalUnits() {
	for u := range s.unitsAir() {
		if !u.terminal(s.step) {
			continue
		}
		if u.explodeWhenTerminal {
			s.explode(u)
		}
		if u.despawnWhenTerminal {
			u.health = 0
		}
	}
}

func (s *state) attackUnits() {
	var unitAttackedOtherUnitCount int

	for _, u := range s.landUnits {
		// TODO: This may be wrong behavior but I am not sure if the issue is occurring
		// e.g. the target was killed by another unit this unit should probably find a potential new target.
		// Then attack that rather than not attack for a turn
		if u.enemyTarget.equals(position{}) {
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

		var nearestVisibleUnit = s.landUnits[u.enemyTarget.PixelPosition]
		var rSum = u.reach + nearestVisibleUnit.radius
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.landUnits[u.enemyTarget.PixelPosition] = nearestVisibleUnit
			s.landUnits[u.position.PixelPosition] = u
		}
	}

	logger.log(
		debugLevel,
		stateLogic,
		"attackUnits",
		slog.Int("units attacked", unitAttackedOtherUnitCount),
	)
}

func (s *state) despawnUnits() {
	var despawnUnitTurnCount int

	for u := range s.units() {
		if u.health <= 0 || (u.collidedThisturn && !u.survivesCollision) {
			gridPos, ok := s.g.getGridPos(u.PixelPosition)
			if !ok {
				panic(fmt.Sprintf("despawnUnits: expected grid pos to exist for %s", u.PixelPosition.log()))
			}
			hex := s.hexagons[gridPos]
			hex.subtract(*u)
			s.deleteUnit(u.position)
			delete(s.reactiveUnits, u.id)
			delete(s.requesterUnits, u.id)
			s.unitDespawnTotal++
			despawnUnitTurnCount++

			if !u.sourceNode.equals(GridPosition{}) {
				n := s.nodes[u.sourceNode]
				switch u.player {
				case pOne:
					n.pOneUnitAliveCount--
				case pTwo:
					n.pTwoUnitAliveCount--
				default:
					panic(fmt.Sprintf("attackUnits: unknown player %d", u.player))
				}
			}
		}
	}

	logger.log(
		debugLevel,
		stateLogic,
		"despawnUnits",
		[]slog.Attr{
			slog.Int("despawn unit turn count", despawnUnitTurnCount),
			slog.Int("despawn unit total", s.unitDespawnTotal)}...,
	)
}

func (s *state) attackNodes() {
	var noTargetNodeCount int
	var attackedUnitThisTurnCount int
	var attackedNodeThisTurnCount int
	var cannotReachNodeCount int

	for _, u := range s.landUnits {
		if len(u.globalTargets) == 0 {
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
			newNode.control = n.control
			newNode.singleSelect = n.singleSelect
			newNode.multiSelect = n.multiSelect
			newNode.pOneUnitAliveCount = n.pOneUnitAliveCount
			newNode.pTwoUnitAliveCount = n.pTwoUnitAliveCount
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
			if n.player == clickedNode.player && n.unitConfig.name == clickedNode.unitConfig.name {
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
		for u := range s.units() {
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

		for _, n := range s.nodes {
			p := s.g.getPixelPosOutbounds(n.gridPos, hpCenter)
			if (topLeftPixelPos.X < p.X) &&
				(p.X < bottomRightPixelPos.X) &&
				(topLeftPixelPos.Y < p.Y) &&
				(p.Y < bottomRightPixelPos.Y) {
				n.selected = true
			} else {
				n.selected = false
			}
		}

		visited := make(map[id]struct{})
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

func (s *state) updateTargets(target position, filter []player) {
	for _, n := range s.nodes {
		if slices.Contains(filter, n.player) {
			continue
		}

		if n.selectActive() {
			n.targets = append(n.targets, target)
		}

		if s.developerMode {
			for _, u := range s.landUnits {
				if u.selected {
					u.updateTarget(globalTarget, target)
				}
			}
		} else if n.singleSelect {
			// TODO: Fix this so air units are selected, need a multi layer iter
			for _, h := range n.control {
				end := h.get(lLand).linkedList
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
			// TODO: Fix this so air units are selected, need a multi layer iter
			end := s.hexagons[n.gridPos].get(lLand).linkedList
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
		if n.player == pNeutral {
			continue
		}

		var unitAliveCount int
		var incrementAliveCount func()
		switch n.player {
		case pOne:
			unitAliveCount = n.pOneUnitAliveCount
			incrementAliveCount = func() { n.pOneUnitAliveCount++ }
		case pTwo:
			unitAliveCount = n.pTwoUnitAliveCount
			incrementAliveCount = func() { n.pTwoUnitAliveCount++ }
		default:
			panic(fmt.Sprintf("produceUnits: unknown player %d", n.player))
		}

		if unitAliveCount == n.unitConfig.maxAlivePerNode {
			continue
		}

		// TODO: determine behavior for producing units during collision
		pixelPos, ok := s.g.getPixelPos(n.gridPos, hpCenter)
		if !ok {
			panic("could not get pixel pos for produce units")
		}
		p := position{n.unitConfig.produceLayer, pixelPos}
		u := newUnit(
			n.unitConfig,
			p,
			n.player,
			n.targets,
			n.actions,
			n.gridPos,
			n.selectActive() && s.developerMode)
		if !ok {
			panic("Could not get node center")
		}
		if _, collides := s.collidesWithUnit(p, u); !collides {
			s.spawn(&u)
			incrementAliveCount()
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
	for _, u := range s.landUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range s.chaffUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range s.bombGravityUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range s.ordanceUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range s.aircraftUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range s.ballisticUnits {
		u.draw(dc, u.selected)
	}
}
