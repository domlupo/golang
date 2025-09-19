// A game component which moves and attacks other game components
package main

import (
	"image/color"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:                  infantry,
	radius:                5,
	maxHealth:             1000,
	attack:                1,
	vision:                30,
	reach:                 6,
	produceLayer:          ground,
	staticTarget:          false,
	attacksGround:         true,
	attacksLowAltitude:    false,
	attacksHighAltitude:   false,
	attacksPosition:       false,
	detonateAtTerminalPos: false,
}
var sniperConfig = unitConfig{
	name:                  sniper,
	radius:                7,
	maxHealth:             5,
	attack:                1,
	vision:                100,
	reach:                 100,
	produceLayer:          ground,
	staticTarget:          false,
	attacksGround:         true,
	attacksLowAltitude:    false,
	attacksHighAltitude:   false,
	attacksPosition:       false,
	detonateAtTerminalPos: false,
}

// TODO: Figure out better naming for various missiles
// TODO: Figure out better attack conventions e.g. direct, vs explosion, etc
var unguidedGroundMissileConfig = unitConfig{
	name:                  unguidedGroundMissile,
	radius:                3,
	maxHealth:             100,
	attack:                1,
	vision:                100,
	reach:                 0,
	produceLayer:          lowAltitude,
	staticTarget:          true,
	attacksGround:         false,
	attacksLowAltitude:    false,
	attacksHighAltitude:   false,
	attacksPosition:       true,
	detonateAtTerminalPos: true,
}
var guidedAirMissileConfig = unitConfig{
	name:                  guidedAirMissile,
	radius:                3,
	maxHealth:             100,
	attack:                1000,
	vision:                3000,
	reach:                 0,
	produceLayer:          lowAltitude,
	staticTarget:          false,
	attacksGround:         false,
	attacksLowAltitude:    true,
	attacksHighAltitude:   true,
	attacksPosition:       false,
	detonateAtTerminalPos: true,
}

type name int

const (
	noName = iota
	infantry
	sniper
	unguidedGroundMissile
	guidedAirMissile
	mobileSAM
)

type unitTarget int

const (
	noTarget = iota
	currentTarget
	fluxTarget
	localTarget
	globalTarget
)

type unitConfig struct {
	name                  name
	radius                int
	maxHealth             int
	attack                int
	vision                int
	reach                 int
	produceLayer          layer
	staticTarget          bool
	attacksGround         bool
	attacksLowAltitude    bool
	attacksHighAltitude   bool
	attacksPosition       bool
	detonateAtTerminalPos bool
}

func (config unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	var name string
	switch config.name {
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
	a = append(a, slog.Int("radius", config.radius))
	a = append(a, slog.Int("maxHealth", config.maxHealth))
	a = append(a, slog.Int("attack", config.attack))
	a = append(a, slog.Int("vision", config.vision))
	a = append(a, slog.Int("reach", config.reach))
	var layer string
	switch config.produceLayer {
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
	a = append(a, slog.Bool("staticTarget", config.staticTarget))
	a = append(a, slog.Bool("attacksGround", config.attacksGround))
	a = append(a, slog.Bool("attacksLowAltitude", config.attacksLowAltitude))
	a = append(a, slog.Bool("attacksHighAltitude", config.attacksHighAltitude))
	a = append(a, slog.Bool("attacksPosition", config.attacksPosition))
	return a
}

type unit struct {
	id                 uuid.UUID
	player             player
	currentTarget      *position
	fluxTarget         *position
	localTarget        *position
	globalTarget       *position
	fluxExpire         int
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
	canReachTargetNode *node
	selected           bool
}

func newUnit(p player, globalTarget *position, config unitConfig) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.log(errorLevel, selectedUnitLogType, "newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:            id,
		player:        p,
		currentTarget: nil,
		fluxTarget:    nil,
		localTarget:   nil,
		globalTarget:  globalTarget,
		// TODO: Figure out better naming for these two flux variables
		// They are a bit confusing but not sure what else to call them at the moment.
		// Most likely they will change significantly.
		fluxExpire:         0,
		fluxTurnCount:      0,
		atTerminalPosition: false,
		bresenhamError:     nil,
		nextBresenhamError: nil,
		unitConfig:         config,
		health:             config.maxHealth,
		movedLastTurn:      false,
		movedThisTurn:      false,
		attackedLastTurn:   false,
		attackedThisTurn:   false,
		canReachTargetNode: nil,
		selected:           false,
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
	u.canReachTargetNode = nil

	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "afterRefresh", "unit", u.log())
}

func (u *unit) getNewPos(currentPos position) position {
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "getNewPos", "currentPos", currentPos.log())

	// Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	// Specifically, the integer incremental error form where error is stored over mulitple moves
	// This algorithm prioritizes straight line movement. It is not guaranteed the final position is reached.
	var x0 = currentPos.x
	var y0 = currentPos.y
	var x1 = 0
	var y1 = 0

	if u.fluxTarget != nil {
		u.fluxExpire = max(u.fluxExpire-1, 0)
		u.updateTarget(currentTarget, u.fluxTarget)
		x1 = u.fluxTarget.x
		y1 = u.fluxTarget.y
		if u.fluxExpire == 0 {
			u.fluxTarget = nil
		}
	} else if u.localTarget != nil {
		u.updateTarget(currentTarget, u.localTarget)
		x1 = u.localTarget.x
		y1 = u.localTarget.y
	} else if u.globalTarget != nil {
		u.updateTarget(currentTarget, u.globalTarget)
		x1 = u.globalTarget.x
		y1 = u.globalTarget.y
	} else {
		u.updateTarget(currentTarget, nil)
		var desiredPos = currentPos
		return desiredPos
	}

	var sx int
	var sy int
	var dx = x1 - x0
	var dy = y1 - y0
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
		var error = dx + dy
		var errorCopy = dx + dy
		u.bresenhamError = &error
		u.nextBresenhamError = &errorCopy
	}
	var doubleError = *u.bresenhamError * 2
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

	var newPos = position{currentPos.layer, pixelPosition{x0, y0}}
	return newPos
}

func (u *unit) updateTarget(ut unitTarget, target *position) {
	var unitTarget string
	switch ut {
	case currentTarget:
		// TODO: Look into if there is a better way to handle nil equality
		if target == nil || u.currentTarget == nil || u.currentTarget != target {
			u.bresenhamError = nil
			u.nextBresenhamError = nil
		}
		u.currentTarget = target
		unitTarget = "currentTarget"
	case fluxTarget:
		u.fluxTarget = target
		unitTarget = "fluxTarget"
	case localTarget:
		u.localTarget = target
		unitTarget = "localTarget"
	case globalTarget:
		u.globalTarget = target
		unitTarget = "globalTarget"
	default:
		panic("Need to set a unit target")
	}

	var targetLog []slog.Attr
	if target == nil {
		targetLog = nil
	} else {
		targetLog = target.log()
	}
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "updateTarget", "unitTarget", unitTarget, "target", targetLog)
}

func (u *unit) completeMove() {
	u.movedThisTurn = true
	var intCopy = &u.nextBresenhamError
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
	if u.currentTarget != nil {
		a = append(a, slog.Any("currentTarget", u.currentTarget.log()))
	} else {
		a = append(a, slog.String("currentTarget", "nil"))
	}
	if u.fluxTarget != nil {
		a = append(a, slog.Any("fluxTarget", u.fluxTarget.log()))
	} else {
		a = append(a, slog.String("fluxTarget", "nil"))
	}
	if u.localTarget != nil {
		a = append(a, slog.Any("localTarget", u.localTarget.log()))
	} else {
		a = append(a, slog.String("localTarget", "nil"))
	}
	if u.globalTarget != nil {
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

func (u *unit) draw(dc drawContext, unitPos position, selected bool) {
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
		float32(unitPos.x)*dc.windowMultiplier,
		float32(unitPos.y)*dc.windowMultiplier,
		float32(u.radius)*dc.windowMultiplier,
		unitColor,
		false,
	)

	if selected {
		drawUnfilledBoxInt(
			dc,
			dc.highlightColor,
			unitPos.x-u.radius,
			unitPos.y-u.radius,
			unitPos.x+u.radius,
			unitPos.y+u.radius,
		)
	}

	if selected && u.currentTarget != nil {
		vector.StrokeLine(
			dc.screen,
			float32(unitPos.x)*dc.windowMultiplier,
			float32(unitPos.y)*dc.windowMultiplier,
			float32(u.currentTarget.x)*dc.windowMultiplier,
			float32(u.currentTarget.y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			dc.highlightColor,
			false,
		)

		if u.unitConfig.name == unguidedGroundMissile {
			vector.StrokeCircle(
				dc.screen,
				float32(u.currentTarget.x)*dc.windowMultiplier,
				float32(u.currentTarget.y)*dc.windowMultiplier,
				float32(missileBlastRadius)*dc.windowMultiplier,
				float32(dc.missileBlastRadiusLineWidth)*dc.windowMultiplier,
				dc.highlightColor,
				false,
			)
		}
	}

	if u.health != u.maxHealth {
		var healthLineLength = (float32(u.health) / float32(u.maxHealth)) * float32(u.radius*2)

		vector.DrawFilledRect(
			dc.screen,
			float32(unitPos.x-u.radius)*dc.windowMultiplier,
			float32(unitPos.y-u.radius)*dc.windowMultiplier,
			healthLineLength*dc.windowMultiplier,
			float32(u.radius)*dc.windowMultiplier,
			dc.unitHealthColor,
			false,
		)
	}
}
