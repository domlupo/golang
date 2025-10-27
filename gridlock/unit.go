// A unit is a game component which moves and attacks other game components.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"math"
	"slices"

	"github.com/gofrs/uuid/v5"
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

// All units on map
type units map[position]*unit

type unit struct {
	position
	unitConfig
	id                 uuid.UUID
	player             player
	activeTarget       position
	steerTarget        position
	enemyTarget        position
	localTarget        position
	localTargets       []position
	globalTarget       position
	pathGlobalTarget   bool
	atTerminalPosition bool
	bresenhamError     *int
	nextBresenhamError *int
	health             int
	movedThisTurn      bool
	attackedThisTurn   bool
	AtCollision        bool
	canReachNode       node
	selected           bool
	prior              *unit
	next               *unit
	producedAt         int
}

func newUnit(p player, globalTarget position, uc unitConfig, pos position, producedOrder int, selected bool) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.log(errorLevel, selectedUnit, "newUnit", slog.String("failed to generate UUID", err.Error()))
		panic("cannot newUnit since uuid generation failed")
	}

	u := unit{
		id:               id,
		position:         pos,
		player:           p,
		globalTarget:     globalTarget,
		pathGlobalTarget: true,
		unitConfig:       uc,
		health:           uc.maxHealth,
		selected:         selected,
		producedAt:       producedOrder,
	}

	return u
}

func (u *unit) refresh() {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "beforeRefresh", u.log()...)
	u.movedThisTurn = false
	u.attackedThisTurn = false
	u.AtCollision = false
}

func (u *unit) prioritizeTargets() {
	if !u.steerTarget.isZero() {
		if u.selected {
			fmt.Println("STEER")
		}
		//u.localTargets = []position{}
		u.updateTarget(activeTarget, u.steerTarget)
	} else if !u.enemyTarget.isZero() {
		u.steerTarget = position{}
		u.localTargets = []position{}
		u.updateTarget(activeTarget, u.enemyTarget)
	} else if len(u.localTargets) > 0 {
		u.steerTarget = position{}
		u.updateTarget(activeTarget, u.localTargets[len(u.localTargets)-1])
	} else {
		u.steerTarget = position{}
		// TODO: determine behavior. Right now this path may occur for multiple reasons or only one.
		// The reason I am sure of is it occurs after a targeted unit dies off since there is no path
		// for the unit to continue. The issue is that that this might cause units to pause for an entire
		// turn since no target is set which seems wrong.
		u.pathGlobalTarget = true
	}
}

// Uses Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
// Specifically, the integer incremental error form where error is stored over mulitple moves
// This algorithm prioritizes straight line movement. It is not guaranteed the final position is reached.
func (u *unit) getNewPos() (position, func()) {
	logger.logSelected(u.selected, infoLevel, selectedUnit, "getNewPos", u.position.log()...)

	if u.activeTarget.isZero() {
		return u.position, func() {}
	}

	x1 := u.activeTarget.X
	y1 := u.activeTarget.Y
	x0 := u.X
	y0 := u.Y

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
	if u.bresenhamError == nil {
		bError := dx + dy
		bErrorCopy := dx + dy
		u.bresenhamError = &bError
		u.nextBresenhamError = &bErrorCopy
	}
	if u.selected && x0 == x1 {
		fmt.Println("X0 == X1")
	}
	if u.selected && y0 == y1 {
		fmt.Println("Y0 == Y1")
	}
	doubleError := *u.bresenhamError * 2
	u.atTerminalPosition = true
	if doubleError >= dy && x0 != x1 {
		u.atTerminalPosition = false
		*u.nextBresenhamError += dy
		x0 += sx
	}
	if doubleError <= dx && y0 != y1 {
		u.atTerminalPosition = false
		*u.nextBresenhamError += dx
		y0 += sy
	}

	newPos := position{u.layer, PixelPosition{x0, y0}}
	return newPos, func() { u.completeMove(newPos) }
}

func (u *unit) updateTarget(ut unitTarget, target position) {
	var unitTarget string
	switch ut {
	case activeTarget:
		if !u.activeTarget.equals(target.PixelPosition) {
			u.bresenhamError = nil
			u.nextBresenhamError = nil
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
		u.globalTarget = target
		u.pathGlobalTarget = true
		unitTarget = "globalTarget"
	default:
		panic("Need to set a unit target")
	}

	var targetLog = []slog.Attr{slog.String("unit target", unitTarget)}
	if !target.isZero() {
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
	var foundPath bool
	tries := [][]position{{u.globalTarget}}
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

func (u *unit) completeMove(pos position) {
	u.position = pos
	u.movedThisTurn = true
	intCopy := &u.nextBresenhamError
	u.bresenhamError = *intCopy
}

func (u unit) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Any("id", u.id))
	var player string
	switch u.player {
	case noPlayer:
		player = "noPlayer"
	case bluePlayer:
		player = "bluePlayer"
	case redPlayer:
		player = "redPlayer"
	case neutralPlayer:
		player = "neutralPlayer"
	}
	a = append(a, slog.String("player", player))
	a = append(a, slog.Any("position", position.log(u.position)))
	a = append(a, slog.Any("next", u.next.id))
	a = append(a, slog.Any("prior", u.prior.id))
	a = append(a, slog.Int("producedOrder", u.producedAt))
	if !u.activeTarget.isZero() {
		a = append(a, slog.Any("activeTarget", u.activeTarget.log()))
	} else {
		a = append(a, slog.String("activeTarget", "nil"))
	}
	if !u.steerTarget.isZero() {
		a = append(a, slog.Any("steerTarget", u.steerTarget.log()))
	} else {
		a = append(a, slog.String("steerTarget", "nil"))
	}
	if !u.enemyTarget.isZero() {
		a = append(a, slog.Any("enemyTarget", u.enemyTarget.log()))
	} else {
		a = append(a, slog.String("enemyTarget", "nil"))
	}
	if !u.localTarget.isZero() {
		a = append(a, slog.Any("localTarget", u.localTarget.log()))
	} else {
		a = append(a, slog.String("localTarget", "nil"))
	}
	if !u.globalTarget.isZero() {
		a = append(a, slog.Any("globalTarget", u.globalTarget.log()))
	} else {
		a = append(a, slog.String("globalTarget", "nil"))
	}
	a = append(a, slog.Bool("atTerminalPosition", u.atTerminalPosition))
	if u.bresenhamError != nil {
		a = append(a, slog.Int("bresenhamError", *u.bresenhamError))
	} else {
		a = append(a, slog.String("bresenhamError", "nil"))
	}
	if u.nextBresenhamError != nil {
		a = append(a, slog.Int("nextBresenhamError", *u.nextBresenhamError))
	} else {
		a = append(a, slog.String("nextBresenhamError", "nil"))
	}
	a = append(a, u.unitConfig.log()...)
	a = append(a, slog.Int("health", u.health))
	a = append(a, slog.Bool("movedThisTurn", u.movedThisTurn))
	a = append(a, slog.Bool("attackedThisTurn", u.attackedThisTurn))
	a = append(a, slog.Any("canReachNode", u.canReachNode))
	a = append(a, slog.Bool("selected", u.selected))
	return a
}

func (u *unit) draw(dc drawConfig, selected bool) {
	var unitColor color.Color
	switch u.player {
	case bluePlayer:
		unitColor = dc.bluePlayerUnitColor
	case redPlayer:
		unitColor = dc.redPlayerUnitColor
	default:
		panic("unit needs a color")
	}

	vector.DrawFilledCircle(
		dc.screen,
		float32(u.X)*dc.windowMultiplier,
		float32(u.Y)*dc.windowMultiplier,
		float32(u.radius)*dc.windowMultiplier,
		unitColor,
		false,
	)

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

	if selected && !u.activeTarget.isZero() {
		if u.unitConfig.name == unguidedGroundMissile {
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

		if u.unitConfig.name == guidedAirMissile {
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

		vector.DrawFilledRect(
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
