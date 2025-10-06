// A game component which moves and attacks other game components
package main

import (
	"image/color"
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:                 infantry,
	radius:               3,
	maxHealth:            2000,
	attack:               1,
	vision:               30,
	reach:                4,
	produceLayer:         ground,
	staticTarget:         false,
	attacksGround:        true,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	attacksPosition:      false,
	explodeAtTerminalPos: false,
	fluxTurnCountDelay:   10,
	fluxMagnitude:        5,
	fluxExpireMax:        5,
	explodeRadius:        0,
}
var sniperConfig = unitConfig{
	name:                 sniper,
	radius:               7,
	maxHealth:            5,
	attack:               1,
	vision:               100,
	reach:                100,
	produceLayer:         ground,
	staticTarget:         false,
	attacksGround:        true,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	attacksPosition:      false,
	explodeAtTerminalPos: false,
	fluxTurnCountDelay:   10,
	fluxMagnitude:        10,
	fluxExpireMax:        10,
	explodeRadius:        0,
}

// TODO: Figure out better naming for various missiles
// TODO: Figure out better attack conventions e.g. direct, vs explosion, etc
var unguidedGroundMissileConfig = unitConfig{
	name:                 unguidedGroundMissile,
	radius:               3,
	maxHealth:            100,
	attack:               1,
	vision:               100,
	reach:                0,
	produceLayer:         lowAltitude,
	staticTarget:         true,
	attacksGround:        false,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	attacksPosition:      true,
	explodeAtTerminalPos: true,
	fluxTurnCountDelay:   10,
	fluxMagnitude:        10,
	fluxExpireMax:        10,
	explodeRadius:        100,
}
var guidedAirMissileConfig = unitConfig{
	name:                 guidedAirMissile,
	radius:               3,
	maxHealth:            100,
	attack:               1000,
	vision:               250,
	reach:                0,
	produceLayer:         lowAltitude,
	staticTarget:         false,
	attacksGround:        false,
	attacksLowAltitude:   true,
	attacksHighAltitude:  true,
	attacksPosition:      false,
	explodeAtTerminalPos: true,
	fluxTurnCountDelay:   10,
	fluxMagnitude:        10,
	fluxExpireMax:        10,
	explodeRadius:        100,
}

type name int

const (
	noName name = iota
	infantry
	sniper
	unguidedGroundMissile
	guidedAirMissile
	mobileSAM
)

type unitTarget int

const (
	noTarget unitTarget = iota
	currentTarget
	fluxTarget
	enemyTarget
	localTarget
	globalTarget
)

type unitConfig struct {
	name                 name
	radius               int
	maxHealth            int
	attack               int
	vision               int
	reach                int
	produceLayer         layer
	staticTarget         bool
	attacksGround        bool
	attacksLowAltitude   bool
	attacksHighAltitude  bool
	attacksPosition      bool
	explodeAtTerminalPos bool
	// TODO: Figure out better flux names
	fluxTurnCountDelay int
	fluxMagnitude      int
	fluxExpireMax      int
	explodeRadius      int
}

func (uc unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	var name string
	switch uc.name {
	case noName:
		name = "no name"
	case infantry:
		name = "infantry"
	case sniper:
		name = "sniper"
	case unguidedGroundMissile:
		name = "unguidedGroundMissile"
	case guidedAirMissile:
		name = "guidedAirMissile"
	default:
		name = "unknown"
	}
	a = append(a, slog.String("name", name))
	a = append(a, slog.Int("radius", uc.radius))
	a = append(a, slog.Int("maxHealth", uc.maxHealth))
	a = append(a, slog.Int("attack", uc.attack))
	a = append(a, slog.Int("vision", uc.vision))
	a = append(a, slog.Int("reach", uc.reach))
	var layer string
	switch uc.produceLayer {
	case noLayer:
		layer = "noLayer"
	case ground:
		layer = "ground"
	case lowAltitude:
		layer = "lowAltitude"
	case highAltitude:
		layer = "highAltitude"
	}
	a = append(a, slog.String("produceLayer", layer))
	a = append(a, slog.Bool("staticTarget", uc.staticTarget))
	a = append(a, slog.Bool("attacksGround", uc.attacksGround))
	a = append(a, slog.Bool("attacksLowAltitude", uc.attacksLowAltitude))
	a = append(a, slog.Bool("attacksHighAltitude", uc.attacksHighAltitude))
	a = append(a, slog.Bool("attacksPosition", uc.attacksPosition))
	return a
}

type unit struct {
	id uuid.UUID
	position
	player             player
	currentTarget      position
	fluxTarget         position
	enemyTarget        position
	localTarget        position
	localTargets       []position
	globalTarget       position
	pathGlobalTarget   bool
	fluxTurnCount      int
	atTerminalPosition bool
	bresenhamError     *int
	nextBresenhamError *int
	unitConfig
	health             int
	movedLastTurn      bool
	movedThisTurn      bool
	attackedLastTurn   bool
	attackedThisTurn   bool
	canReachTargetNode node
	selected           bool
	fluxExpireCount    int
	prior              *unit
	next               *unit
	producedAt         int
}

func newUnit(p player, globalTarget position, uc unitConfig, pos position, producedOrder int) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.log(errorLevel, selectedUnitLogType, "newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:                 id,
		position:           pos,
		player:             p,
		currentTarget:      position{},
		fluxTarget:         position{},
		enemyTarget:        position{},
		localTarget:        position{},
		localTargets:       []position{},
		globalTarget:       globalTarget,
		pathGlobalTarget:   true,
		fluxTurnCount:      0,
		atTerminalPosition: false,
		bresenhamError:     nil,
		nextBresenhamError: nil,
		unitConfig:         uc,
		health:             uc.maxHealth,
		movedLastTurn:      false,
		movedThisTurn:      false,
		attackedLastTurn:   false,
		attackedThisTurn:   false,
		canReachTargetNode: node{},
		selected:           false,
		fluxExpireCount:    uc.fluxExpireMax,
		producedAt:         producedOrder,
	}

	return u
}

func (u *unit) refresh() {
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "beforeRefresh", "unit", u.log())

	if !u.attackedLastTurn && !u.movedLastTurn {
		u.fluxTurnCount += 1
	} else {
		u.fluxTurnCount = 0
	}

	u.movedLastTurn = u.movedThisTurn
	u.movedThisTurn = false
	u.attackedLastTurn = u.attackedThisTurn
	u.attackedThisTurn = false
	u.canReachTargetNode = node{}

	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "afterRefresh", "unit", u.log())
}

func (u *unit) analyzeTargets() {
	if !u.fluxTarget.isZero() {
		u.localTargets = []position{}
		u.fluxExpireCount = max(u.fluxExpireCount-1, 0)
		u.updateTarget(currentTarget, u.fluxTarget)
		if u.fluxExpireCount == 0 {
			u.fluxTarget = position{}
			u.fluxExpireCount = u.fluxExpireMax
		}
	} else if !u.enemyTarget.isZero() {
		u.localTargets = []position{}
		u.updateTarget(currentTarget, u.enemyTarget)
	} else if !u.localTarget.isZero() {
		u.updateTarget(currentTarget, u.localTarget)
	} else {
		u.localTargets = []position{}
		u.updateTarget(currentTarget, position{})
	}
}

// Uses Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
// Specifically, the integer incremental error form where error is stored over mulitple moves
// This algorithm prioritizes straight line movement. It is not guaranteed the final position is reached.
func (u *unit) getNewPos() (position, func()) {
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "getNewPos", "currentPos", u.position.log())

	if u.currentTarget.isZero() {
		return u.position, func() {}
	}

	x1 := u.currentTarget.x
	y1 := u.currentTarget.y
	x0 := u.x
	y0 := u.y

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

	newPos := position{u.layer, pixelPosition{x0, y0}}
	return newPos, func() { u.completeMove(newPos) }
}

func (u *unit) updateTarget(ut unitTarget, target position) {
	var unitTarget string
	switch ut {
	case currentTarget:
		if target.isZero() || u.currentTarget.isZero() {
			// TODO: Determine if bresenham needs zeroed more often.
			u.bresenhamError = nil
			u.nextBresenhamError = nil
		}
		u.currentTarget = target
		unitTarget = "currentTarget"
	case fluxTarget:
		u.fluxTarget = target
		u.fluxTurnCount = 0
		unitTarget = "fluxTarget"
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

	var targetLog []slog.Attr
	if target.isZero() {
		targetLog = nil
	} else {
		targetLog = target.log()
	}
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "updateTarget", "unitTarget", unitTarget, "target", targetLog)
}

func (u *unit) pathfindGlobalTarget(visibilityGrid []obstacleTarget, pixelPathing pixelPathing) {
	obstaclePixelPos, los := u.globalTarget.pixelPosition.lineOfSight(u.pixelPosition, pixelPathing)
	switch los {
	case true:
		u.updateTarget(localTarget, u.globalTarget)
	case false:
		nearest := math.MaxInt64
		var nearestIndex int
		for index, ot := range visibilityGrid {
			distance := obstaclePixelPos.chebyshevDistance(ot.pixelPosition)
			if distance < nearest {
				nearest = distance
				nearestIndex = index
			}
		}

		// TODO: Fix incorrect behavior. Right now cw and ccw pathing assumes that every visibility grid point
		// should be visited until a line of sight is reached. This assumption is not always true. For example,
		// if an obstacle layout is similar to an H where the middle crevice has multiple points, * -> H ->*
		// the unit will incorrectly go into the crevice rather than skip those points. I don't think
		// it is as simple as comparing n and n+ 1, but this needs determined.
		current := visibilityGrid[nearestIndex]
		cwPath := []position{u.globalTarget}
		cwPathDistance := current.pixelPosition.chebyshevDistance(u.globalTarget.pixelPosition)
		end := current.ccw
		for current != *end {
			cwPath = append(cwPath, position{layer: u.layer, pixelPosition: current.pixelPosition})
			_, los := current.lineOfSight(u.pixelPosition, pixelPathing)
			if los {
				break
			}
			current = *current.cw
			cwPathDistance += current.ccwDistance
		}
		cwPathDistance += u.chebyshevDistance(cwPath[len(cwPath)-1].pixelPosition)

		current = visibilityGrid[nearestIndex]
		ccwPath := []position{u.globalTarget}
		ccwPathDistance := current.pixelPosition.chebyshevDistance(u.globalTarget.pixelPosition)
		end = current.cw
		for current != *end {
			ccwPath = append(ccwPath, position{layer: u.layer, pixelPosition: current.pixelPosition})
			_, los := current.lineOfSight(u.pixelPosition, pixelPathing)
			if los {
				break
			}
			current = *current.ccw
			ccwPathDistance += current.cwDistance
		}
		ccwPathDistance += u.chebyshevDistance(ccwPath[len(ccwPath)-1].pixelPosition)

		shortestPath := cwPath
		if ccwPathDistance < cwPathDistance {
			shortestPath = ccwPath
		}

		u.localTargets = shortestPath
		u.pathGlobalTarget = false
		u.updateTarget(localTarget, shortestPath[len(shortestPath)-1])
	}
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
	if !u.currentTarget.isZero() {
		a = append(a, slog.Any("currentTarget", u.currentTarget.log()))
	} else {
		a = append(a, slog.String("currentTarget", "nil"))
	}
	if !u.fluxTarget.isZero() {
		a = append(a, slog.Any("fluxTarget", u.fluxTarget.log()))
	} else {
		a = append(a, slog.String("fluxTarget", "nil"))
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
	a = append(a, slog.Bool("movedLastTurn", u.movedLastTurn))
	a = append(a, slog.Bool("movedLastTurn", u.movedLastTurn))
	a = append(a, slog.Bool("attackedLastTurn", u.attackedLastTurn))
	a = append(a, slog.Bool("attackedThisTurn", u.attackedThisTurn))
	a = append(a, slog.Int("fluxTurnCount", u.fluxTurnCount))
	a = append(a, slog.Any("canReachTargetNode", u.canReachTargetNode))
	a = append(a, slog.Bool("selected", u.selected))
	return a
}

func (u *unit) draw(dc drawContext, selected bool) {
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
		float32(u.x)*dc.windowMultiplier,
		float32(u.y)*dc.windowMultiplier,
		float32(u.radius)*dc.windowMultiplier,
		unitColor,
		false,
	)

	if selected {
		drawUnfilledBoxInt(
			dc,
			dc.highlightColor,
			u.x-u.radius,
			u.y-u.radius,
			u.x+u.radius,
			u.y+u.radius,
		)
	}

	if selected && !u.currentTarget.isZero() {
		vector.StrokeLine(
			dc.screen,
			float32(u.x)*dc.windowMultiplier,
			float32(u.y)*dc.windowMultiplier,
			float32(u.currentTarget.x)*dc.windowMultiplier,
			float32(u.currentTarget.y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			dc.highlightColor,
			false,
		)

		for _, lt := range u.localTargets {
			vector.StrokeLine(
				dc.screen,
				float32(u.x)*dc.windowMultiplier,
				float32(u.y)*dc.windowMultiplier,
				float32(lt.x)*dc.windowMultiplier,
				float32(lt.y)*dc.windowMultiplier,
				float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}

		if u.unitConfig.name == unguidedGroundMissile {
			vector.StrokeCircle(
				dc.screen,
				float32(u.currentTarget.x)*dc.windowMultiplier,
				float32(u.currentTarget.y)*dc.windowMultiplier,
				float32(u.explodeRadius)*dc.windowMultiplier,
				float32(dc.missileExplodeRadiusLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}

		if u.unitConfig.name == guidedAirMissile {
			vector.StrokeCircle(
				dc.screen,
				float32(u.x)*dc.windowMultiplier,
				float32(u.y)*dc.windowMultiplier,
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
			float32(u.x-u.radius)*dc.windowMultiplier,
			float32(u.y-u.radius)*dc.windowMultiplier,
			healthLineLength*dc.windowMultiplier,
			float32(u.radius)*dc.windowMultiplier,
			dc.unitHealthColor,
			false,
		)
	}
}
