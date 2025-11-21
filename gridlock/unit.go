// A unit is a game component which moves and attacks other game components.

package main

import (
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

// An action which should occur when a condition is met.
type reaction struct {
	when  func(*unit) bool
	do    func(*unit)
	next  uint64
	delay uint64
}

// All units within a 2 dimensional layer
type units map[PixelPosition]*unit

type unit struct {
	unitConfig
	position
	subPixelPosition    subPixelPosition
	id                  id
	player              player
	activeTarget        position
	steerTarget         position
	enemyTarget         position
	localTarget         position
	localTargets        []position
	globalTargets       []position
	nearest             target
	skipTargets         []id
	actions             []action
	pathGlobalTarget    bool
	atTerminalPosition  bool
	bresenhamError      int
	bresenhamErrorDy    int
	bresenhamErrorDx    int
	steerTurnCount      int
	steerTurnCountLimit int
	health              int
	movedThisTurn       bool
	attackedThisTurn    bool
	AtCollision         bool
	canReachNode        node
	selected            bool
	prior               *unit
	next                *unit
	spawnOrder          int
	angle               float64
	sourceNode          GridPosition
	sourcePosition      PixelPosition
	sourceUnit          id
}

// High resolution two dimensional point.
type subPixelPosition struct {
	x int
	y int
}

func (s subPixelPosition) log() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, slog.Int("x", s.x))
	a = append(a, slog.Int("y", s.y))
	return a
}

func (u unit) newSubPixelPosition(p position) subPixelPosition {
	return subPixelPosition{
		x: int(math.Round(float64(p.X) * u.speedFactor)),
		y: int(math.Round(float64(p.Y) * u.speedFactor)),
	}
}

func (u unit) downscale(s subPixelPosition) PixelPosition {
	return PixelPosition{
		X: int(math.Round(float64(s.x) / float64(u.speedFactor))),
		Y: int(math.Round(float64(s.y) / float64(u.speedFactor))),
	}
}

func newUnit(
	uc unitConfig,
	pos position,
	p player,
	globalTargets []position,
	actions []action,
	sourceNode GridPosition,
	selected bool) unit {
	u := unit{
		id:                  generateId(),
		position:            pos,
		player:              p,
		globalTargets:       globalTargets,
		pathGlobalTarget:    true,
		actions:             actions,
		unitConfig:          uc,
		health:              uc.maxHealth,
		selected:            selected,
		steerTurnCountLimit: 20,
		sourceNode:          sourceNode,
		sourcePosition:      pos.PixelPosition,
	}
	u.subPixelPosition = u.newSubPixelPosition(u.position)

	if u.name == fighter {
		u.reactions = []reaction{
			{
				when: func(unit *unit) bool {
					return unit.PixelPosition.equals(unit.sourcePosition) && unit.missileCount != unit.missileCountMax
				},
				do: func(unit *unit) {
					unit.missileCount = unit.missileCountMax
					unit.skipTargets = []id{}
				},
				delay: u.missileFireDelay,
			},
			{
				when: func(unit *unit) bool { return !unit.nearest.equals(position{}) && unit.missileCount != 0 },
				do: func(unit *unit) {
					newUnit := newUnit(
						missileConfig,
						position{layer: missileConfig.produceLayer, PixelPosition: unit.PixelPosition},
						unit.player,
						[]position{unit.nearest.position},
						[]action{},
						GridPosition{},
						false)
					newUnit.sourceUnit = unit.id
					unit.skipTargets = append(unit.skipTargets, unit.nearest.id)
					unit.requests = append(unit.requests, request{&newUnit, false})
					unit.missileCount--
				},
				delay: u.missileFireDelay,
			}}
	}

	return u
}

func (u unit) attackLayers() []layer {
	var layers []layer
	if u.attacksLand {
		layers = append(layers, land)
	}
	if u.attacksLowAir {
		layers = append(layers, lowAir)
	}
	if u.attacksHighAir {
		layers = append(layers, highAir)
	}
	return layers
}

func (u unit) collideLayers() []layer {
	var layers []layer
	if u.collidesLand {
		layers = append(layers, land)
	}
	if u.collidesLowAir {
		layers = append(layers, lowAir)
	}
	if u.collidesHighAir {
		layers = append(layers, highAir)
	}
	return layers
}

func (u *unit) refresh() {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "beforeRefresh", u.log()...)
	u.movedThisTurn = false
	u.attackedThisTurn = false
	u.AtCollision = false
	u.angle = u.PixelPosition.angle(u.globalTargets[0].PixelPosition)
}

func (u *unit) prioritizeTargets() {
	u.steerTurnCount++
	if !u.steerTarget.equals(position{}) && u.steerTurnCount < u.steerTurnCountLimit {
		u.updateTarget(activeTarget, u.steerTarget)
	} else if !u.enemyTarget.equals(position{}) {
		u.steerTurnCount = 0
		u.steerTarget = position{}
		u.localTargets = []position{}
		u.updateTarget(activeTarget, u.enemyTarget)
	} else if len(u.localTargets) > 0 {
		u.steerTurnCount = 0
		u.steerTarget = position{}
		u.updateTarget(activeTarget, u.localTargets[len(u.localTargets)-1])
	} else {
		u.steerTurnCount = 0
		u.steerTarget = position{}
		// TODO: determine behavior. Right now this path may occur for multiple reasons or only one.
		// The reason I am sure of is it occurs after a targeted unit dies off since there is no path
		// for the unit to continue. The issue is that that this might cause units to pause for an entire
		// turn since no target is set which seems wrong.
		u.pathGlobalTarget = true
	}
}

// Uses Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham's_line_algorithm
// Specifically, the integer incremental error form where error is stored over mulitple moves
// This algorithm prioritizes straight line movement. It is not guaranteed the final position is reached.
func (u *unit) getMoves() []Move {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "getMoves: ", u.position.log()...)

	moves := make([]Move, 0)
	u.atTerminalPosition = true
	u.bresenhamErrorDx = 0
	u.bresenhamErrorDy = 0

	s := u.newSubPixelPosition(u.activeTarget)
	x1 := s.x
	y1 := s.y
	x0 := u.subPixelPosition.x
	y0 := u.subPixelPosition.y

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
		u.atTerminalPosition = false
		u.bresenhamErrorDy = dy
		x0 += sx
	}
	if doubleError <= dx && y0 != y1 {
		u.atTerminalPosition = false
		u.bresenhamErrorDx = dx
		y0 += sy
	}

	if u.atTerminalPosition {
		return moves
	}

	s = subPixelPosition{x0, y0}
	moves = append(moves, u.newMove(s))

	// TODO: instead of arbitrarily picking second and third option, prioritize to minimize error
	x0 = u.subPixelPosition.x
	y0 = u.subPixelPosition.y
	if u.bresenhamErrorDy != 0 && u.bresenhamErrorDx != 0 {
		x0 += sx
		s := subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
		x0 -= sx
		y0 += sy
		s = subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
	} else if u.bresenhamErrorDy != 0 {
		y0 += sy
		s := subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
		x0 += sx
		s = subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
	} else if u.bresenhamErrorDx != 0 {
		x0 += sx
		s := subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
		y0 += sy
		s = subPixelPosition{x0, y0}
		moves = append(moves, u.newMove(s))
	}

	u.bresenhamErrorDy = dy
	u.bresenhamErrorDx = dx
	return moves
}

func (u *unit) newMove(s subPixelPosition) Move {
	return Move{
		position: position{u.layer, u.downscale(s)},
		confirm:  func() { u.confirmMove(s) }}
}

func (u *unit) confirmMove(s subPixelPosition) {
	u.position = position{u.layer, u.downscale(s)}
	u.movedThisTurn = true
	if u.subPixelPosition.x != s.x {
		u.bresenhamError += u.bresenhamErrorDy
	}
	if u.subPixelPosition.y != s.y {
		u.bresenhamError += u.bresenhamErrorDx
	}
	u.subPixelPosition = s
}

func (u *unit) updateTarget(ut unitTarget, target position) {
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
		u.globalTargets = []position{target}
		u.pathGlobalTarget = true
		unitTarget = "globalTargets"
	default:
		panic("Need to set a unit target")
	}

	var targetLog = []slog.Attr{slog.String("unit target", unitTarget)}
	if !target.equals(position{}) {
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
// A smaller issue is when the desired final position is within a visbility grid, but the los check from
// the position to unit immediately jumps to a different obstacle visibility grid when the first
// visibility grid should have been traveled.
func (u *unit) pathfindGlobalTarget(obstacles obstacles, oLUT obstacleLUT, g grid) int {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "pathfindGlobalTarget")

	var foundPath bool
	tries := [][]position{{u.globalTargets[0]}}
	nextTries := [][]position{}

	for !foundPath {
		for _, path := range tries {
			start := path[len(path)-1]
			obstaclePoint, los := start.obstaclePoint(u.PixelPosition, oLUT, obstacles, g)

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
		nextTries = [][]position{}
	}

	paths := [][]position{}
	for _, try := range tries {
		end := try[len(try)-1]
		if _, los := u.lineOfSight(end.PixelPosition, oLUT, u.radius); los {
			paths = append(paths, try)
		}
	}

	for i, path := range paths {
		trimPath := []position{path[0]}
		for j := 0; j < len(path); j++ {
			pos := path[j]
			for k := len(path) - 1; k > j; k-- {
				otherPos := path[k]
				if _, los := pos.lineOfSight(otherPos.PixelPosition, oLUT, u.radius); los {
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
		distance := u.chebyshevDistance(path[len(path)-1].PixelPosition)
		for i := 0; i < len(path)-1; i++ {
			distance += path[i].chebyshevDistance(path[i+1].PixelPosition)
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
	g grid) []position {

	path := []position{}

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
		path = append(path, position{layer: u.layer, PixelPosition: current.PixelPosition})
		obstaclePoint, los := current.obstaclePoint(u.PixelPosition, oLUT, obstacles, g)

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

func (u unit) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Any("id", u.id))
	a = append(a, slog.String("player", u.player.String()))
	a = append(a, nest("position", u.position.log()))
	a = append(a, nest("subPixelPosition", u.subPixelPosition.log()))
	a = append(a, slog.Any("next", u.next.id))
	a = append(a, slog.Any("prior", u.prior.id))
	a = append(a, slog.Int("spawnOrder", u.spawnOrder))
	a = append(a, nest("sourceNode", u.sourceNode.log()))
	a = append(a, nest("sourcePosition", u.sourcePosition.log()))
	a = append(a, slog.Int("missileCount", u.missileCount))
	if !u.activeTarget.equals(position{}) {
		a = append(a, slog.Any("activeTarget", u.activeTarget.log()))
	} else {
		a = append(a, slog.String("activeTarget", "nil"))
	}
	if !u.steerTarget.equals(position{}) {
		a = append(a, slog.Any("steerTarget", u.steerTarget.log()))
	} else {
		a = append(a, slog.String("steerTarget", "nil"))
	}
	if !u.enemyTarget.equals(position{}) {
		a = append(a, slog.Any("enemyTarget", u.enemyTarget.log()))
	} else {
		a = append(a, slog.String("enemyTarget", "nil"))
	}
	if !u.localTarget.equals(position{}) {
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
	a = append(a, slog.Bool("atTerminalPosition", u.atTerminalPosition))
	a = append(a, slog.Int("bresenhamError", u.bresenhamError))
	a = append(a, slog.Int("bresenhamErrorDx", u.bresenhamErrorDx))
	a = append(a, slog.Int("bresenhamErrorDy", u.bresenhamErrorDy))
	a = append(a, u.unitConfig.log()...)
	a = append(a, slog.Int("health", u.health))
	a = append(a, slog.Bool("movedThisTurn", u.movedThisTurn))
	a = append(a, slog.Bool("attackedThisTurn", u.attackedThisTurn))
	a = append(a, slog.Any("canReachNode", u.canReachNode))
	a = append(a, slog.Bool("selected", u.selected))
	return a
}

func (u unit) angleDraw(dc drawConfig) *ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	xMult := math.Sin(u.angle) - math.Cos(u.angle)
	yMult := -math.Sqrt2 * math.Sin(u.angle+(math.Pi/4))
	op.GeoM.Rotate(u.angle)
	op.GeoM.Translate(float64(u.X)+float64(u.radius)*xMult, float64(u.Y)+float64(u.radius)*yMult)
	op.GeoM.Scale(float64(dc.windowMultiplier), float64(dc.windowMultiplier))
	return &op
}

func (u unit) draw(dc drawConfig, selected bool) {
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
		dc.screen.DrawImage(dc.pngFighter, u.angleDraw(dc))

		if u.selected {
			vector.StrokeCircle(
				dc.screen,
				float32(u.X)*dc.windowMultiplier,
				float32(u.Y)*dc.windowMultiplier,
				float32(u.vision)*dc.windowMultiplier,
				float32(dc.defaultLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}
	case bomber:
		dc.screen.DrawImage(dc.pngBomber, u.angleDraw(dc))

		for _, a := range u.actions {
			vector.StrokeCircle(
				dc.screen,
				float32(a.pos.X)*dc.windowMultiplier,
				float32(a.pos.Y)*dc.windowMultiplier,
				// HACK: hard coded value. Needed because this is the first instance of a unit
				// drawing another unit config. It is Low priority to fix this. Also this pattern
				// will occur in the future. Easier to fix once more examples to understand structure.
				20.0*dc.windowMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}
	case missile:
		dc.screen.DrawImage(dc.pngMissile, u.angleDraw(dc))
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
			float32(u.X)*dc.windowMultiplier,
			float32(u.Y)*dc.windowMultiplier,
			float32(u.radius)*dc.windowMultiplier,
			color,
			false,
		)
	default:
		vector.FillCircle(
			dc.screen,
			float32(u.X)*dc.windowMultiplier,
			float32(u.Y)*dc.windowMultiplier,
			float32(u.radius)*dc.windowMultiplier,
			unitColor,
			false,
		)
	}

	if selected {
		drawUnfilledBoxInt(
			dc,
			dc.highlightColor,
			u.X-u.radius,
			u.Y-u.radius,
			u.X+u.radius,
			u.Y+u.radius,
		)
	}

	if selected && !u.activeTarget.equals(position{}) {
		if u.name == unguidedGroundMissile {
			vector.StrokeCircle(
				dc.screen,
				float32(u.activeTarget.X)*dc.windowMultiplier,
				float32(u.activeTarget.Y)*dc.windowMultiplier,
				float32(u.explodeRadius)*dc.windowMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}

		if u.name == guidedAirMissile {
			vector.StrokeCircle(
				dc.screen,
				float32(u.X)*dc.windowMultiplier,
				float32(u.Y)*dc.windowMultiplier,
				float32(u.vision)*dc.windowMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}
	}

	if u.health != u.maxHealth {
		var healthLineLength = (float32(u.health) / float32(u.maxHealth)) * float32(u.radius*2)

		vector.FillRect(
			dc.screen,
			float32(u.X-u.radius)*dc.windowMultiplier,
			float32(u.Y-u.radius)*dc.windowMultiplier,
			healthLineLength*dc.windowMultiplier,
			float32(u.radius)*dc.windowMultiplier,
			dc.unitHealthColor,
			false,
		)
	}
}
