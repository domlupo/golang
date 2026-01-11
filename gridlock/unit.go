// A unit is a game component which moves and attacks other game components.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Units path around obstacles in cw or ccw motion
type pathDirection int

const (
	noPathDirection pathDirection = iota
	cw
	ccw
)

type unitTarget int

const (
	noTarget unitTarget = iota
	activeTarget
	steerTarget
	enemyTarget
	localTarget
	globalTarget
)

func (ut unitTarget) String() string {
	switch ut {
	case noTarget:
		return "noTarget"
	case activeTarget:
		return "activeTarget"
	case steerTarget:
		return "steerTarget"
	case enemyTarget:
		return "enemyTarget"
	case localTarget:
		return "localTarget"
	case globalTarget:
		return "globalTarget"
	default:
		panic(fmt.Sprintf("Cannot convert unit target to string: %d", ut))
	}
}

// All units within a 2 dimensional layer. It may be wrong to use Pixel as key when unit positions
// are the higher resolution Pos but I think this is simpler and can switch if needed.
type units map[Pixel]*unit

type unit struct {
	unitConfig
	Pos
	subPos              subPos
	id                  id
	player              player
	activeTarget        Pixel
	activeTargetType    unitTarget
	steerTarget         Pixel
	enemyTarget         Pixel
	localTarget         Pixel
	localTargets        []Pixel
	globalTargets       []Pixel
	nearest             *unit
	skipTargets         []id
	actions             []action
	pathGlobalTarget    bool
	terminalPosition    Pos
	moveDone            bool
	bresenhamError      int
	bresenhamErrorDy    int
	bresenhamErrorDx    int
	steerTurnCount      int
	steerTurnCountLimit int
	health              int
	movedThisTurn       bool
	attackedThisTurn    bool
	collidedThisturn    bool
	canReachNode        node
	selected            bool
	prior               *unit
	next                *unit
	spawnOrder          int
	angle               float64
	sourceNode          GridPosition
	sourcePosition      Pos
	sourceUnit          id
	showVision          bool
}

// High resolution 2D point.
type subPos struct {
	x int
	y int
}

func (s subPos) log() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, slog.Int("x", s.x))
	a = append(a, slog.Int("y", s.y))
	return a
}

func (u unit) newSubPos(p Pos) subPos {
	return subPos{
		x: int(math.Round(float64(p.X) * u.speedFactor)),
		y: int(math.Round(float64(p.Y) * u.speedFactor)),
	}
}

func (u unit) downscale(s subPos) Pos {
	return Pos{
		X: math.Round(float64(s.x) / float64(u.speedFactor)),
		Y: math.Round(float64(s.y) / float64(u.speedFactor)),
	}
}

func newUnit(
	uc unitConfig,
	pos Pos,
	p player,
	globalTargets []Pixel,
	actions []action,
	sourceNode GridPosition,
	selected bool) *unit {

	u := unit{
		id:                  generateId(),
		Pos:                 pos,
		player:              p,
		globalTargets:       globalTargets,
		pathGlobalTarget:    true,
		actions:             actions,
		unitConfig:          uc,
		health:              uc.maxHealth,
		selected:            selected,
		steerTurnCountLimit: 20,
		sourceNode:          sourceNode,
		sourcePosition:      pos,
	}

	u.subPos = u.newSubPos(u.Pos)
	for i, e := range u.equipments {
		e.attached = &u
		u.equipments[i] = e
	}

	if u.vision < u.radius {
		panic(fmt.Sprintf("newUnit: vision should not be less than radius for %s", uc.name.string()))
	}
	return &u
}

func (u unit) terminal(step uint64) bool {
	if u.moveDone && u.isTerminalWhenMoveDone {
		return true
	}
	if step >= u.terminalStep {
		return true
	}
	return false
}

func (u *unit) refresh() {
	u.showVision = false
	u.movedThisTurn = false
	u.attackedThisTurn = false
	u.collidedThisturn = false

	// TODO: Fix once target contract / structure is better understood. Needed for human players.
	if len(u.globalTargets) == 0 {
		u.angle = 0
	} else {
		u.angle = u.Pos.noLayer().angle(u.globalTargets[0].pos())
	}

	for i, e := range u.equipments {
		e.rotate()
		u.equipments[i] = e
	}
}

func (u *unit) prioritizeTargets() {
	u.steerTurnCount++
	if !u.steerTarget.equals(Pixel{}) && u.steerTurnCount < u.steerTurnCountLimit {
		u.updateTarget(activeTarget, u.steerTarget)
		u.activeTargetType = steerTarget
	} else if !u.enemyTarget.equals(Pixel{}) {
		u.steerTurnCount = 0
		u.steerTarget = Pixel{}
		u.localTargets = []Pixel{}
		u.updateTarget(activeTarget, u.enemyTarget)
		u.activeTargetType = enemyTarget
	} else if len(u.localTargets) > 0 {
		u.steerTurnCount = 0
		u.steerTarget = Pixel{}
		u.updateTarget(activeTarget, u.localTargets[len(u.localTargets)-1])
		u.activeTargetType = localTarget
	} else {
		u.steerTurnCount = 0
		u.steerTarget = Pixel{}
		// TODO: determine behavior. Right now this path may occur for multiple reasons or only one.
		// The reason I am sure of is it occurs after a targeted unit dies off since there is no path
		// for the unit to continue. The issue is that that this might cause units to pause for an entire
		// turn since no target is set which seems wrong.
		u.pathGlobalTarget = true
	}
}

// Uses Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham's_line_algorithm
// Specifically, the integer incremental error form where error is stored over multiple moves
// This algorithm prioritizes straight line movement. It is not guaranteed the final position is reached.
func (u *unit) getMoves() []move {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "getMoves: ", u.Pos.log()...)

	moves := make([]move, 0)
	u.moveDone = true
	u.bresenhamErrorDx = 0
	u.bresenhamErrorDy = 0

	s := u.newSubPos(u.activeTarget.pos())
	x1 := s.x
	y1 := s.y
	x0 := u.subPos.x
	y0 := u.subPos.y

	var sx int
	var sy int
	dx := x1 - x0
	dy := y1 - y0
	if dx < 0 {
		dx = dx * -1
	}
	if dy > 0 {
		dy = dy * -1
	}
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	if u.bresenhamError == 0 {
		u.bresenhamError = dx + dy
	}
	doubleError := u.bresenhamError * 2
	if doubleError >= dy && x0 != x1 {
		u.moveDone = false
		u.bresenhamErrorDy = dy
		x0 += sx
	}
	if doubleError <= dx && y0 != y1 {
		u.moveDone = false
		u.bresenhamErrorDx = dx
		y0 += sy
	}

	if u.moveDone {
		return moves
	}

	s = subPos{x0, y0}
	moves = append(moves, u.newMove(s))

	// TODO: instead of arbitrarily picking second and third option, prioritize to minimize error
	x0 = u.subPos.x
	y0 = u.subPos.y
	if u.bresenhamErrorDy != 0 && u.bresenhamErrorDx != 0 {
		x0 += sx
		s := subPos{x0, y0}
		moves = append(moves, u.newMove(s))
		x0 -= sx
		y0 += sy
		s = subPos{x0, y0}
		moves = append(moves, u.newMove(s))
	} else if u.bresenhamErrorDy != 0 {
		y0 += sy
		s := subPos{x0, y0}
		moves = append(moves, u.newMove(s))
		x0 += sx
		s = subPos{x0, y0}
		moves = append(moves, u.newMove(s))
	} else if u.bresenhamErrorDx != 0 {
		x0 += sx
		s := subPos{x0, y0}
		moves = append(moves, u.newMove(s))
		y0 += sy
		s = subPos{x0, y0}
		moves = append(moves, u.newMove(s))
	}

	u.bresenhamErrorDy = dy
	u.bresenhamErrorDx = dx
	return moves
}

func (u *unit) newMove(s subPos) move {
	return move{
		Pos:     u.downscale(s).layer(u.L),
		confirm: func() { u.confirmMove(s) }}
}

func (u *unit) confirmMove(s subPos) {
	u.Pos = u.downscale(s).layer(u.L)
	u.movedThisTurn = true
	if u.subPos.x != s.x {
		u.bresenhamError += u.bresenhamErrorDy
	}
	if u.subPos.y != s.y {
		u.bresenhamError += u.bresenhamErrorDx
	}
	u.subPos = s
}

func (u *unit) updateTargetUnit(ut unitTarget, unit *unit) {
	if unit == nil {
		u.updateTarget(ut, Pixel{})
	} else {
		u.updateTarget(ut, unit.pixel())
	}
}

func (u *unit) updateTarget(ut unitTarget, target Pixel) {
	// TODO: Determine if this should be refactored to protect itself from zero value updates.
	// I think it should, but it breaks other code.

	var unitTarget string
	switch ut {
	case activeTarget:
		if !u.activeTarget.equals(target) {
			u.bresenhamError = 0
		}
		u.activeTarget = target
		unitTarget = "activeTarget"
	case steerTarget:
		u.steerTarget = target
		unitTarget = "steerTarget"
	case enemyTarget:
		u.enemyTarget = target
		unitTarget = "enemyTarget"
	case localTarget:
		u.localTarget = target
		unitTarget = "localTarget"
	case globalTarget:
		u.globalTargets = []Pixel{target}
		u.pathGlobalTarget = true
		unitTarget = "globalTargets"
	default:
		panic("Need to set a unit target")
	}

	var targetLog = []slog.Attr{slog.String("unit target", unitTarget)}
	if !target.equals(Pixel{}) {
		targetLog = append(targetLog, target.log()...)
	}
	logger.logSelected(u.selected, infoLevel, selectedUnit, "updateTarget", targetLog...)
}

// Not optimal but close enough to unblock future development. A hack is used which is not guaranteed
// to find the best path. I suspect future path finding will need to be optimal and cache results at startup.
//
// The hack is when creating a path, los checks are performed from the current path position to the unit. The
// path only considers jumping to obstacles which directly block the los but certain paths are more efficient
// when jumping to obstacles which are not in the los.
//
// A smaller issue is when the desired final position is within a visibility grid, but the los check from
// the position to unit immediately jumps to a different obstacle visibility grid when the first
// visibility grid should have been traveled.
func (u *unit) pathfindGlobalTarget(obstacles obstacles, oLUT obstacleLUT, g grid) int {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "pathfindGlobalTarget")

	var foundPath bool
	tries := [][]Pixel{{u.globalTargets[0]}}
	nextTries := [][]Pixel{}

	for !foundPath {
		for _, path := range tries {
			start := path[len(path)-1]
			obstaclePoint, los := start.obstaclePoint(u.pixel(), oLUT, obstacles, g)

			if los {
				foundPath = true
				nextTries = append(nextTries, path)
				continue
			} else if !los && obstaclePoint.isZero() {
				// TODO: HACK which handles cursor off screen. Need to determine
				// what is desired player experience here.
				foundPath = true
				break
			}

			cwPath := u.getPath(obstaclePoint, cw, oLUT, obstacles, g)
			ccwPath := u.getPath(obstaclePoint, ccw, oLUT, obstacles, g)

			nextTries = append(nextTries, slices.Concat(path, cwPath))
			nextTries = append(nextTries, slices.Concat(path, ccwPath))
		}
		tries = nextTries
		nextTries = [][]Pixel{}
	}

	paths := [][]Pixel{}
	for _, try := range tries {
		end := try[len(try)-1]
		if _, los := u.pixel().lineOfSight(end, oLUT, u.radius); los {
			paths = append(paths, try)
		}
	}

	for i, path := range paths {
		trimPath := []Pixel{path[0]}
		for j := 0; j < len(path); j++ {
			pos := path[j]
			for k := len(path) - 1; k > j; k-- {
				otherPixel := path[k]
				if _, los := pos.lineOfSight(otherPixel, oLUT, u.radius); los {
					j = k
					trimPath = append(trimPath, path[k])
					j--
					break
				}
			}
		}
		paths[i] = trimPath
	}

	shortestDistance := math.MaxInt64
	shortestPath := 0
	for i, path := range paths {
		distance := u.pixel().chebyshevDistance(path[len(path)-1])
		for i := 0; i < len(path)-1; i++ {
			distance += path[i].chebyshevDistance(path[i+1])
		}
		if distance < shortestDistance {
			shortestDistance = distance
			shortestPath = i
		}
	}

	if len(paths) == 0 {
		// HACK which handles two issues. The first is when cursor is off screen.
		// The second is a rare bug where a unit cannot find a path.
		return 0
	}

	u.pathGlobalTarget = false
	u.localTargets = paths[shortestPath]
	return shortestDistance
}

func (u *unit) getPath(
	current obstaclePoint,
	direction pathDirection,
	oLUT obstacleLUT,
	obstacles obstacles,
	g grid) []Pixel {

	path := []Pixel{}

	var end obstaclePoint
	switch direction {
	case cw:
		end = *current.ccw
	case ccw:
		end = *current.cw
	default:
		panic("To set end, path direction must be clockwise or counter clockwise")
	}

	for current != end {
		path = append(path, current.Pixel)
		obstaclePoint, los := current.obstaclePoint(u.pixel(), oLUT, obstacles, g)

		if los || current.visibilityGridId != obstaclePoint.visibilityGridId {
			break
		}

		switch direction {
		case cw:
			current = *current.cw
		case ccw:
			current = *current.ccw
		default:
			panic("To travel, path direction must be clockwise or counter clockwise")
		}
	}

	return path
}

func (u unit) isPayloadNotFull() bool {
	for _, v := range u.consumables {
		if v.count != v.countMax {
			return true
		}
	}
	return false
}

func (u *unit) refillPayload() {
	for k := range u.consumables {
		c := u.consumables[k]
		c.count = c.countMax
		u.consumables[k] = c
	}
}

// This log method is frequently used for debugging. Log frequently used info at start
// or end for quicker debugging. Less used info can be parsed via tooling e.g. grep
func (u unit) log() []slog.Attr {
	var a = u.logShort()
	a = append(a, slog.String("player", u.player.string()))
	a = append(a, nest("subPos", u.subPos.log()))
	if u.next == nil {
		a = append(a, slog.String("next", "none"))
	} else {
		a = append(a, slog.Any("next", u.next.logShort()))
	}
	if u.prior == nil {
		a = append(a, slog.String("prior", "none"))
	} else {
		a = append(a, slog.Any("prior", u.prior.logShort()))
	}
	a = append(a, slog.Int("spawnOrder", u.spawnOrder))
	a = append(a, nest("sourceNode", u.sourceNode.log()))
	a = append(a, nest("sourcePosition", u.sourcePosition.log()))
	if !u.activeTarget.equals(Pixel{}) {
		a = append(a, slog.Any("activeTarget", u.activeTarget.log()))
	} else {
		a = append(a, slog.String("activeTarget", "nil"))
	}
	if !u.steerTarget.equals(Pixel{}) {
		a = append(a, slog.Any("steerTarget", u.steerTarget.log()))
	} else {
		a = append(a, slog.String("steerTarget", "nil"))
	}
	if !u.enemyTarget.equals(Pixel{}) {
		a = append(a, slog.Any("enemyTarget", u.enemyTarget.log()))
	} else {
		a = append(a, slog.String("enemyTarget", "nil"))
	}
	if !u.localTarget.equals(Pixel{}) {
		a = append(a, slog.Any("localTarget", u.localTarget.log()))
	} else {
		a = append(a, slog.String("localTarget", "nil"))
	}
	if len(u.localTargets) == 0 {
		a = append(a, slog.String("localTargets", "none"))
	} else {
		localtargets := []slog.Attr{}
		for _, p := range u.localTargets {
			localtargets = append(localtargets, p.log()...)
		}
		a = append(a, nest("localTargets", localtargets))
	}
	if len(u.globalTargets) == 0 {
		a = append(a, slog.String("globalTargets", "none"))
	} else {
		globalTargets := []slog.Attr{}
		for _, p := range u.globalTargets {
			globalTargets = append(globalTargets, p.log()...)
		}

		a = append(a, nest("globalTargets", globalTargets))
	}
	if len(u.actions) == 0 {
		a = append(a, slog.String("actions", "none"))
	} else {
		actions := []slog.Attr{}
		for _, a := range u.actions {
			actions = append(actions, a.log()...)
		}
		a = append(a, nest("actions", actions))
	}
	a = append(a, nest("terminalPosition", u.terminalPosition.log()))
	a = append(a, slog.Bool("moveDone", u.moveDone))
	a = append(a, slog.Int("bresenhamError", u.bresenhamError))
	a = append(a, slog.Int("bresenhamErrorDx", u.bresenhamErrorDx))
	a = append(a, slog.Int("bresenhamErrorDy", u.bresenhamErrorDy))
	a = append(a, slog.Int("health", u.health))
	a = append(a, slog.Bool("movedThisTurn", u.movedThisTurn))
	a = append(a, slog.Bool("attackedThisTurn", u.attackedThisTurn))
	a = append(a, slog.Any("canReachNode", u.canReachNode))
	a = append(a, slog.Bool("selected", u.selected))
	a = append(a, slog.String("activeTargetType", u.activeTargetType.String()))
	if u.nearest == nil {
		a = append(a, slog.String("nearest", "none"))
	} else {
		a = append(a, slog.Any("nearest", u.nearest.logShort()))
	}
	a = append(a, slog.Bool("collidedThisTurn", u.collidedThisturn))
	a = append(a, u.unitConfig.log()...)
	return a
}

func (u unit) logShort() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, slog.Any("id", u.id))
	a = append(a, nest("Pos", u.Pos.log()))
	return a
}

func (u unit) angleDraw(pixelMultiplier float32) *ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	xMult := math.Sin(u.angle) - math.Cos(u.angle)
	yMult := -math.Sqrt2 * math.Sin(u.angle+(math.Pi/4))
	op.GeoM.Rotate(u.angle)
	op.GeoM.Translate(float64(u.pixel().X)+float64(u.radius)*xMult, float64(u.pixel().Y)+float64(u.radius)*yMult)
	op.GeoM.Scale(float64(pixelMultiplier), float64(pixelMultiplier))
	return &op
}

func (u unit) draw(dc drawConfig, selected bool, pixelMultiplier float32) {
	var unitColor color.Color
	switch u.player {
	case pOne:
		unitColor = dc.pOneUnitColor
	case pTwo:
		unitColor = dc.pTwoUnitColor
	default:
		panic("unit needs a color")
	}

	switch u.name {
	case fighter:
		dc.screen.DrawImage(dc.pngFighter, u.angleDraw(pixelMultiplier))
		for _, e := range u.equipments {
			e.draw(dc, pixelMultiplier)
		}
	case bomber:
		dc.screen.DrawImage(dc.pngBomber, u.angleDraw(pixelMultiplier))

		if u.selected {
			vector.StrokeCircle(
				dc.screen,
				float32(u.X)*pixelMultiplier,
				float32(u.Y)*pixelMultiplier,
				float32(u.vision)*pixelMultiplier,
				float32(dc.defaultLineWidth)*pixelMultiplier,
				dc.highlightColor,
				false,
			)
		}

		for _, a := range u.actions {
			vector.StrokeCircle(
				dc.screen,
				float32(a.X)*pixelMultiplier,
				float32(a.Y)*pixelMultiplier,
				// TODO: hard coded value. Needed because this is the first instance of a unit
				// drawing another unit config. It is Low priority to fix this. Also this pattern
				// will occur in the future. Easier to fix once more examples to understand structure.
				20.0*pixelMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*pixelMultiplier,
				dc.highlightColor,
				false,
			)
		}
	case missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, missileRadarActive:
		dc.screen.DrawImage(dc.pngMissile, u.angleDraw(pixelMultiplier))

		if u.selected && (u.name == missileHeatSeeking || u.name == missileAntiBallistic || u.name == missileRadarActive) {
			vector.StrokeCircle(
				dc.screen,
				float32(u.X)*pixelMultiplier,
				float32(u.Y)*pixelMultiplier,
				float32(u.vision)*pixelMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*pixelMultiplier,
				dc.highlightColor,
				false,
			)
		}
	case flare:
		dc.screen.DrawImage(dc.pngFlare, u.angleDraw(pixelMultiplier))
	case chaff:
		vector.StrokeCircle(
			dc.screen,
			float32(u.X)*pixelMultiplier,
			float32(u.Y)*pixelMultiplier,
			float32(u.radius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			unitColor,
			false,
		)
	case infantry:
		var color color.Color
		switch u.player {
		case pOne:
			color = dc.pOneInfantryColor
		case pTwo:
			color = dc.pTwoInfantryColor
		default:
			panic("infantry needs a color")
		}
		vector.FillCircle(
			dc.screen,
			float32(u.pixel().X)*pixelMultiplier,
			float32(u.pixel().Y)*pixelMultiplier,
			float32(u.radius)*pixelMultiplier,
			color,
			false,
		)
	default:
		vector.FillCircle(
			dc.screen,
			float32(u.pixel().X)*pixelMultiplier,
			float32(u.pixel().Y)*pixelMultiplier,
			float32(u.radius)*pixelMultiplier,
			unitColor,
			false,
		)
	}

	if selected && !u.activeTarget.equals(Pixel{}) {
		if u.name == missileCruise {
			vector.StrokeCircle(
				dc.screen,
				float32(u.activeTarget.X)*pixelMultiplier,
				float32(u.activeTarget.Y)*pixelMultiplier,
				float32(u.explodeRadius)*pixelMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*pixelMultiplier,
				dc.highlightColor,
				false,
			)
		}
	}

	if u.health != u.maxHealth {
		var healthLineLength = (float32(u.health) / float32(u.maxHealth)) * float32(u.radius*2)

		vector.FillRect(
			dc.screen,
			float32(u.pixel().X-u.radius)*pixelMultiplier,
			float32(u.pixel().Y-u.radius)*pixelMultiplier,
			healthLineLength*pixelMultiplier,
			float32(u.radius)*pixelMultiplier,
			dc.unitHealthColor,
			false,
		)
	}
}

func (u unit) drawVision(dc drawConfig, pixelMultiplier float32) {
	vector.StrokeCircle(
		dc.screen,
		float32(u.pixel().X)*pixelMultiplier,
		float32(u.pixel().Y)*pixelMultiplier,
		float32(u.vision)*pixelMultiplier,
		float32(dc.defaultLineWidth)*pixelMultiplier,
		dc.highlightColor,
		false,
	)
}
