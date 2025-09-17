// A game component which moves and attacks other game components
package main

import (
	"image/color"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
	currentTarget      unitTarget
	globalTarget       *position
	localTarget        *position
	atTerminalPosition bool
	bresenhamError     *int
	nextBresenhamError *int
	unitConfig
	health             int
	movedLastTurn      bool
	movedThisTurn      bool
	attackedLastTurn   bool
	attackedThisTurn   bool
	fluxTurnCount      int
	canReachTargetNode *node
	selected           bool
}

func newUnit(p player, globalTarget *position, config unitConfig) unit {
	id, err := uuid.NewV4()
	if err != nil {
		logger.log(errorLevel, selectedUnitLogType, "newUnit", "failed to generate UUID", err)
	}
	u := unit{
		id:                 id,
		player:             p,
		currentTarget:      noTarget,
		globalTarget:       globalTarget,
		localTarget:        nil,
		atTerminalPosition: false,
		bresenhamError:     nil,
		nextBresenhamError: nil,
		unitConfig:         config,
		health:             config.maxHealth,
		movedLastTurn:      false,
		movedThisTurn:      false,
		attackedLastTurn:   false,
		attackedThisTurn:   false,
		fluxTurnCount:      0,
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
	u.localTarget = nil
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
	if u.localTarget != nil {
		u.updateTarget(localTarget)
		x1 = u.localTarget.x
		y1 = u.localTarget.y
	} else if u.globalTarget != nil {
		u.updateTarget(globalTarget)
		x1 = u.globalTarget.x
		y1 = u.globalTarget.y
	} else {
		u.updateTarget(noTarget)
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

func (u *unit) updateTarget(target unitTarget) bool {
	var t string
	switch target {
	case noTarget:
		t = "noTarget"
	case localTarget:
		t = "localTarget"
	case globalTarget:
		t = "globalTarget"
	}
	logger.logSelected(u.selected, infoLevel, selectedUnitLogType, "updateTarget", "target", t)
	if target != u.currentTarget {
		u.currentTarget = target
		u.bresenhamError = nil
		u.nextBresenhamError = nil
		return true
	}
	return false
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
	var target string
	switch u.currentTarget {
	case localTarget:
		target = "localTarget"
	case globalTarget:
		target = "globalTarget"
	case noTarget:
		target = "noTarget"
	case neutralPlayer:
		player = "neutralPlayer"
	}
	a = append(a, slog.String("currentTarget", target))
	if u.globalTarget != nil {
		a = append(a, slog.Any("globalTarget", u.globalTarget.log()))
	} else {
		a = append(a, slog.String("globalTarget", "nil"))
	}
	if u.localTarget != nil {
		a = append(a, slog.Any("localTarget", u.localTarget.log()))
	} else {
		a = append(a, slog.String("localTarget", "nil"))
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

func (u *unit) draw(screen *ebiten.Image, unitPos position, selected bool) {
	var unitColor color.Color
	switch u.player {
	case bluePlayer:
		unitColor = bluePlayerUnitColor
	case redPlayer:
		unitColor = redPlayerUnitColor
	default:
		panic("unit needs a color")
	}

	vector.DrawFilledCircle(
		screen,
		float32(unitPos.x)*windowMultiplier,
		float32(unitPos.y)*windowMultiplier,
		float32(u.radius)*windowMultiplier,
		unitColor,
		false,
	)

	if selected {
		drawUnfilledBoxInt(
			screen,
			highlightColor,
			unitPos.x-u.radius,
			unitPos.y-u.radius,
			unitPos.x+u.radius,
			unitPos.y+u.radius,
		)
	}

	var xTarget *int
	var yTarget *int
	if selected && u.currentTarget != noTarget {
		if u.currentTarget == localTarget && u.localTarget != nil {
			xTarget = &u.localTarget.x
			yTarget = &u.localTarget.y
		} else if u.currentTarget == globalTarget && u.globalTarget != nil {
			xTarget = &u.globalTarget.x
			yTarget = &u.globalTarget.y
		}
	}

	if xTarget != nil && yTarget != nil {
		vector.StrokeLine(
			screen,
			float32(unitPos.x)*windowMultiplier,
			float32(unitPos.y)*windowMultiplier,
			float32(*xTarget)*windowMultiplier,
			float32(*yTarget)*windowMultiplier,
			float32(targetNodeLineWidth)*windowMultiplier,
			highlightColor,
			false,
		)
	}

	if drawUnitHealth && u.health != u.maxHealth {
		var healthLineLength = (float32(u.health) / float32(u.maxHealth)) * float32(u.radius*2)

		vector.DrawFilledRect(
			screen,
			float32(unitPos.x-u.radius)*windowMultiplier,
			float32(unitPos.y-u.radius)*windowMultiplier,
			healthLineLength*windowMultiplier,
			float32(u.radius)*windowMultiplier,
			unitHealthColor,
			false,
		)
	}
}
