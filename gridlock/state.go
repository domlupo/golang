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
	// helps isolate issues during debugging. For example, if toggling spatial partitioning off
	// prevents bug from occurring, the bug is most likely within partitioner.go logic.
	optimizeSpatialPartitions bool
	spawnUnitRadiusLimit      int // TODO: move this into partitioner. Not sure why I added here.
}

func newStateConfig(g grid) stateConfig {
	c := stateConfig{
		moveUnits:                 true,
		produceUnits:              true,
		produceUnitsStepInterval:  180,
		optimizeSpatialPartitions: true,
		spawnUnitRadiusLimit:      g.minimumPixelsBetweenNonAdjacentHexagons,
	}
	return c
}

// A counter for game ticks
type step = uint64

// A request to modify state.
type request struct {
	spawn    *unit
	complete bool
}

// not enforced but must be a singleton
type state struct {
	config           stateConfig
	grid             grid
	partitioner      partitioner
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
	step             step
	developerMode    bool
	unitSpawnTotal   int
	unitDespawnTotal int
}

func newState(n nodes, o obstacles, h hexagons, oLUT obstacleLUT, nLUT nodeLUT, g grid,
	p partitioner, developerMode bool) *state {
	s := state{
		config:           newStateConfig(g),
		grid:             g,
		partitioner:      p,
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
	logger.log(infoLevel, stateLogic, "state.tick: start", slog.Uint64("step", s.step))
	logger.log(infoLevel, stateLogic, "state.tick: refresh")
	s.refresh()
	logger.log(infoLevel, stateLogic, "state.tick: analyzeTargets")
	s.analyzeTargets()
	logger.log(infoLevel, stateLogic, "state.tick: processReactions")
	s.processReactions()
	logger.log(infoLevel, stateLogic, "state.tick: processRequests")
	s.processRequests()
	logger.log(infoLevel, stateLogic, "state.tick: processTerminalPositions")
	s.processTerminalUnits()
	logger.log(infoLevel, stateLogic, "state.tick: attackUnits")
	s.attackUnits()
	logger.log(infoLevel, stateLogic, "state.tick: attackNodes")
	s.attackNodes()
	logger.log(infoLevel, stateLogic, "state.tick: moveUnits")
	if s.config.moveUnits {
		s.moveUnits()
	}
	logger.log(infoLevel, stateLogic, "state.tick: despawnUnits")
	s.despawnUnits()
	logger.log(infoLevel, stateLogic, "state.tick: produceUnits")
	if s.config.produceUnits && s.step%s.config.produceUnitsStepInterval == 0 {
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

func (s state) unitsLayer(p Pos) *units {
	var units units
	switch p.L {
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
		panic(fmt.Sprintf("state.unitsLayer: unknown layer %s", p.L.log()))
	}
	return &units
}

func (s *state) spawn(u *unit) {
	// TODO: Verify if this should be >= rather than more conservative >
	if u.radius > s.config.spawnUnitRadiusLimit && !u.isSkipSpawnRadiusCheck {
		panic(fmt.Sprintf(
			"state.spawn: cannot spawn unit radius %d exceeds spawn radius limit %d",
			u.radius,
			s.config.spawnUnitRadiusLimit))
	}
	t, ok := s.grid.getTile(u.pixel())
	if !ok {
		panic(fmt.Sprintf("state.spawn: cannot get tile for Pixel %s", u.pixel().log()))
	}
	hex := s.hexagons[t]
	hex.add(u)

	if u.isReactive {
		s.reactiveUnits[u.id] = u
	}
	if u.isRequester {
		s.requesterUnits[u.id] = u
	}

	units := *s.unitsLayer(u.Pos)
	if _, ok := units[u.pixel()]; !ok {
		units[u.pixel()] = u
	} else {
		panic(fmt.Sprintf(
			"state.spawn cannot spawn unit at %s since it is already occupied.",
			u.pixel().log()))
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
	units := *s.unitsLayer(r.spawn.Pos)
	if _, ok := units[r.spawn.pixel()]; ok {
		return false
	}
	s.spawn(r.spawn)
	return true
}

func (s *state) moveUnit(u *unit, m move) {
	old, oOk := s.grid.getTile(u.pixel())
	new, nOk := s.grid.getTile(m.pixel())
	if !oOk || !nOk {
		panic(fmt.Sprintf(
			"state.moveUnit: cannot move unit %s with move %s",
			u.pixel().log(),
			m.pixel().log()))
	}
	if !old.equals(new) {
		oldHex := s.hexagons[old]
		oldHex.subtract(*u)
		newHex := s.hexagons[new]
		newHex.add(u)
	}
	s.deleteUnit(u.Pos)
	m.confirm()
	units := *s.unitsLayer(m.Pos)
	if o, ok := units[u.pixel()]; ok {
		panic(fmt.Sprintf(
			"state.moveUnit: Pixel for unit occupied by other unit.\n unit: %s \n other: %s \n",
			u.log(),
			o.log()))
	}
	units[u.pixel()] = u
}

func (s *state) deleteUnit(p Pos) {
	var units units
	switch p.L {
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
		panic(fmt.Sprintf("state.deleteUnit: unknown layer %s", p.L.log()))
	}
	delete(units, p.pixel())
}

func (s *state) refresh() {
	for u := range s.units() {
		u.refresh()
	}
}

func (s *state) processReactions() {
	logger.log(infoLevel,
		stateLogic, "state.processReactions", slog.Int("reactiveUnitCount", len(s.reactiveUnits)))
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
	logger.log(infoLevel,
		stateLogic, "state.processRequests", slog.Int("requesterUnitCounter", len(s.requesterUnits)))
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

func (s *state) partitions(u unit, pa partitionAnalysis) []*hexagon {
	if !s.config.optimizeSpatialPartitions {
		parts := []*hexagon{}
		for _, part := range s.hexagons {
			parts = append(parts, part)
		}
		return parts
	}
	return s.partitioner.partitions(u, pa, s.config.spawnUnitRadiusLimit)
}

func (s *state) analyzeTargets() {
	for u := range s.units() {
		if u.name == chaff || u.name == bombGravity {
			continue
		}

		if u.name == missileHeatSeeking ||
			u.name == missileAntiBallistic || u.name == missileRadarActive {
			s.nearestTargetWithinVision(u)
			u.updateTargetUnit(enemyTarget, u.nearest)
		}

		if u.name == missileHeatSeeking ||
			u.name == flare ||
			u.name == missileCruise ||
			u.name == missileBallistic ||
			u.name == missileAntiBallistic ||
			u.name == missileRadarActive {
			u.localTargets = []Pixel{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		if u.name == fighter {
			s.nearestTargetWithinVision(u)
			if u.moveDone && len(u.globalTargets) > 1 {
				if len(u.actions) > 0 && u.actions[0].equals(u.globalTargets[0]) {
					u.actions[0].do()
					u.actions = u.actions[1:]
				}
				u.globalTargets = u.globalTargets[1:]
			} else if u.moveDone && len(u.globalTargets) == 1 {
				u.globalTargets = s.nodes[u.sourceNode].targets
				u.actions = s.nodes[u.sourceNode].actions
			}
			u.localTargets = []Pixel{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		if u.name == bomber {
			s.nearestTargetWithinVision(u)
			if u.moveDone && len(u.globalTargets) > 1 {
				if len(u.actions) > 0 && u.actions[0].equals(u.globalTargets[0]) {
					u.actions[0].do()
					u.actions = u.actions[1:]
				}
				u.globalTargets = u.globalTargets[1:]
			} else if u.moveDone && len(u.globalTargets) == 1 {
				u.globalTargets = s.nodes[u.sourceNode].targets
				u.actions = s.nodes[u.sourceNode].actions
			}
			u.localTargets = []Pixel{u.globalTargets[0]}
			u.prioritizeTargets()
			continue
		}

		steer := s.steer(u)
		if !u.pixel().equals(steer) {
			u.updateTarget(steerTarget, steer)
		}

		u.canReachNode = node{}
		if s.nLUT[u.pixel()].ChebyshevDistance < u.reach {
			if len(s.nLUT[u.pixel()].NearestNodes) > 1 {
				logger.log(
					errorLevel,
					undecidedBehavior,
					"Undecided behavior. Unit can reach multiple nodes. Arbitrarily picking one")
			}
			p := s.nLUT[u.pixel()].NearestNodes[0]
			n := s.nodes[p]
			if u.player != n.player {
				u.canReachNode = *n
			}
		}

		s.nearestTargetWithinVision(u)
		u.updateTargetUnit(enemyTarget, u.nearest)
		if len(u.globalTargets) > 0 && u.pathGlobalTarget {
			u.pathfindGlobalTarget(s.obstacles, s.oLUT, s.grid)
		} else if len(u.localTargets) > 0 && u.moveDone {
			u.localTargets = u.localTargets[:len(u.localTargets)-1]
		}

		u.prioritizeTargets()
	}
}

func (s *state) steer(u *unit) Pixel {
	partitions := s.partitions(*u, paSteer)

	var xSteer int
	var ySteer int
	for _, partition := range partitions {
		if partition == nil {
			continue
		}

		linkedList := partition.get(u.Pos.L).linkedList

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
			if _, ok := collides(u.Pos, current.Pos, u.radius+1, current.radius+1); ok {
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

	return Pixel{X: u.pixel().X + xSteer, Y: u.pixel().Y + ySteer}
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
			if u.Pos.L != lLand {
				move = m
				break
			}
			if s.oLUT[m.pixel()].UnwalkableChebyshevDistance > u.radius {
				move = m
				break
			}
		}

		var collidedBorder bool
		var collidedUnitPos Pos
		var collidesWithUnit bool
		if move.equals(Pos{}) {
			// TODO: Should this be optimized to only perform this check when the unit
			// does not typically move e.g. chaff
			collidedBorder = s.grid.collidesWithBorder(u.pixel(), u.radius)
			collidedUnitPos, collidesWithUnit = s.collidesWithUnit(u.Pos, *u)
		} else {
			collidedBorder = s.grid.collidesWithBorder(move.pixel(), u.radius)
			collidedUnitPos, collidesWithUnit = s.collidesWithUnit(move.Pos, *u)
		}

		if collidedBorder && collidesWithUnit {
			unitsBorderCollideCount++
			unitsCollideCount++
		} else if collidedBorder {
			unitsBorderCollideCount++
		} else if collidesWithUnit {
			unitsCollideCount++
			otherUnits := *s.unitsLayer(collidedUnitPos)
			otherUnit := otherUnits[collidedUnitPos.pixel()]
			otherUnit.collidedThisturn = true
			u.collidedThisturn = true
		} else if !u.canReachNode.isZero() {
			unitsAtNode++
		} else if move.equals(Pos{}) {
			unitsNotMovedCount++
		} else {
			unitsMovedCount++
			s.moveUnit(u, move)
		}
	}

	logger.log(
		infoLevel,
		stateLogic,
		"state.moveUnits",
		slog.Int("units moved", unitsMovedCount),
		slog.Int("units not moved", unitsNotMovedCount),
		slog.Int("units at node", unitsAtNode),
		slog.Int("units collided", unitsCollideCount),
		slog.Int("units already attacked", unitsAlreadyAttacked),
		slog.Int("units at final position", unitsAtFinalPosition),
	)
}

func (s *state) explode(u *unit) {
	u.health = 0
	x1 := u.X
	y1 := u.Y
	r1 := u.explodeRadius

	for _, u := range s.landUnits {
		x2 := u.X + float64(u.radius)
		y2 := u.Y + float64(u.radius)
		r2 := u.radius

		xDelta := x2 - x1
		yDelta := y2 - y1
		rSum := float64(r1 + r2)
		if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
			u.health = 0
			s.landUnits[u.pixel()] = u
		}
	}
}

func (s state) nearestTargetWithinVision(u *unit) {
	partitions := s.partitions(*u, paVision)
	d := math.MaxFloat64
	var t *unit
	for _, partition := range partitions {
		if partition == nil {
			continue
		}

		for _, l := range u.trackLayers {
			linkedList := partition.get(l).linkedList
			if target, distance := s.nearestTargetWithinLayer(*u, linkedList); distance < d {
				t = target
				d = distance
			}
		}
	}
	u.nearest = t
}

func (s state) nearestTargetWithinLayer(u unit, linkedList *unit) (*unit, float64) {
	d := math.MaxFloat64
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

		squaredDistance, collides := collides(u.Pos, current.Pos, u.vision, current.vision)
		if squaredDistance < d && collides {
			if _, los := u.pixel().lineOfSight(current.pixel(), s.oLUT, u.radius); los {
				d = squaredDistance
				t = current
			}
			if current.L != lLand {
				d = squaredDistance
				t = current
			}
		}

		current = current.next
	}
	return t, d
}

func (s *state) collidesWithUnit(p Pos, u unit) (Pos, bool) {
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

	return Pos{}, false
}

func (s state) collidesWithUnitLayer(u unit, p Pos, linkedList *unit) (Pos, bool) {
	if linkedList == nil || (linkedList.id == u.id && linkedList == linkedList.next) {
		return Pos{}, false
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

		if u.ignoreCollision != nil && u.ignoreCollision(current) && u.L == current.L {
			panic(fmt.Sprintf(
				"state.collidesWithUnitLayer: Cannot ignore collisions between units in same layer "+
					"unit: %s, current: %s",
				u.log(),
				current.log()))
		}

		if u.ignoreCollision != nil && u.ignoreCollision(current) ||
			current.id == u.sourceUnit ||
			current.sourceUnit == u.id {
			current = current.next
			continue
		}

		if _, ok := collides(p, current.Pos, u.radius, current.radius); ok {
			return current.Pos, true
		}
		current = current.next
	}
	return Pos{}, false
}
func (s *state) processTerminalUnits() {
	for u := range s.unitsAir() {
		if !u.terminal(s.step) {
			continue
		}
		if u.isExplodeWhenTerminal {
			s.explode(u)
		}
		if u.isDespawnWhenTerminal {
			u.health = 0
		}
	}
}

func (s *state) attackUnits() {
	var unitAttackedOtherUnitCount int

	for _, u := range s.landUnits {
		// TODO: This may be wrong behavior but I am not sure if the issue is occurring e.g. the target
		// was killed by another unit this unit should probably find a potential new target. Then attack
		// that rather than not attack for a turn
		if u.enemyTarget.equals(Pixel{}) {
			continue
		}

		xTargetPixel := u.enemyTarget.X
		yTargetPixel := u.enemyTarget.Y

		var xDelta int
		var yDelta int
		if u.pixel().X < (xTargetPixel) {
			xDelta = xTargetPixel - u.pixel().X
		} else if u.pixel().X > (xTargetPixel) {
			xDelta = u.pixel().X - xTargetPixel
		}
		if u.pixel().Y < (yTargetPixel) {
			yDelta = yTargetPixel - u.pixel().Y
		} else if u.pixel().Y > (yTargetPixel) {
			yDelta = u.pixel().Y - yTargetPixel
		}

		var nearestVisibleUnit = s.landUnits[u.enemyTarget]
		var rSum = u.reach + nearestVisibleUnit.reach
		if rSum*rSum >= xDelta*xDelta+yDelta*yDelta {
			nearestVisibleUnit.health = nearestVisibleUnit.health - u.attack
			u.attackedThisTurn = true
			unitAttackedOtherUnitCount++
			s.landUnits[u.enemyTarget] = nearestVisibleUnit
			s.landUnits[u.pixel()] = u
		}
	}

	logger.log(
		debugLevel,
		stateLogic,
		"state.attackUnits",
		slog.Int("units attacked", unitAttackedOtherUnitCount),
	)
}

func (s *state) despawnUnits() {
	var despawnUnitTurnCount int

	for u := range s.units() {
		if u.health <= 0 || (u.collidedThisturn && !u.isSurvivesCollision) {
			t, ok := s.grid.getTile(u.pixel())
			if !ok {
				panic(fmt.Sprintf("state.despawnUnits: expected tile to exist for %s", u.pixel().log()))
			}
			hex := s.hexagons[t]
			hex.subtract(*u)
			s.deleteUnit(u.Pos)
			delete(s.reactiveUnits, u.id)
			delete(s.requesterUnits, u.id)
			s.unitDespawnTotal++
			despawnUnitTurnCount++

			if !u.sourceNode.equals(Tile{}) {
				n := s.nodes[u.sourceNode]
				switch u.player {
				case pOne:
					n.pOneUnitAliveCount--
				case pTwo:
					n.pTwoUnitAliveCount--
				default:
					panic(fmt.Sprintf("state.despawnUnits: unknown player %d", u.player))
				}
			}
		}
	}

	logger.log(
		debugLevel,
		stateLogic,
		"state.despawnUnits",
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
			n := s.nodes[u.canReachNode.tile]
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
		"state.attackNodes",
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
			newNode := newNode(n.capturedBy, n.tile, n.maxHealth, n.unitConfig)
			newNode.control = n.control
			newNode.singleSelect = n.singleSelect
			newNode.multiSelect = n.multiSelect
			newNode.pOneUnitAliveCount = n.pOneUnitAliveCount
			newNode.pTwoUnitAliveCount = n.pTwoUnitAliveCount
			s.nodes[n.tile] = &newNode
			nodeOwnershipChanged = true
		}
	}

	return nodeOwnershipChanged
}

func (s *state) updateTargets(target Pixel, filter []player) {
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
			// TODO: Fix this so air units are selected, need a multi layer iterator
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
			// TODO: Fix this so air units are selected, need a multi layer iterator
			end := s.hexagons[n.tile].get(lLand).linkedList
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
			panic(fmt.Sprintf("state.produceUnits: unknown player %d", n.player))
		}

		if unitAliveCount == n.unitConfig.maxAlivePerNode {
			continue
		}

		// TODO: determine behavior for producing units during collision
		pixel, ok := s.grid.getPixel(n.tile, hpCenter)
		if !ok {
			panic(fmt.Sprintf("state.produceUnits: could not get Pixel for tile %s", n.tile.log()))
		}
		p := pixel.posLayer(n.unitConfig.produceLayer)
		u := newUnit(
			n.unitConfig,
			p,
			n.player,
			n.targets,
			n.actions,
			n.tile,
			n.selectActive() && s.developerMode)
		if _, collides := s.collidesWithUnit(p, *u); !collides {
			s.spawn(u)
			incrementAliveCount()
		}
	}
}
