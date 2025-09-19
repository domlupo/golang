// main entry for high level game logic
package main

import (
	"errors"
	"image/color"
	"log"
	"log/slog"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var slowTick = false

const slowTickDuration = 10 * time.Second

// game balance config
var fluxMagnitude = 10
var unitFluxCountLimit = 10
var nodeRadius = 40 / 2 // TODO: Fix hardcoding
var initialNodeHealth = 100
var missileBlastRadius = 100

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

var fluxLineWidth = 2
var fluxDirectionRadius = 3
var targetNodeLineWidth = 1

// debug grid config
// The debug grid helps debugging display issues. Usually only altering
// offset is needed, width can remain 1.
// Offset values should only be 0 - 20 (half of grid size without window multipler)
// Some helpful values are: 0, 0.5, 1, 19, 19.5, 20
var debugGridColor = color.Black

const debugGridSquareCount = 32

// other configs
var runProfiler = false

type Game struct {
	// float32 rather than float64 since alot of Ebiten draw APIs use float32
	*drawContext
	backgroundColor          color.RGBA
	produceWait              time.Duration
	produceDelay             time.Duration
	slowTickTime             time.Time
	priorTime                time.Time
	currentTime              time.Time
	leftMouseClicked         bool
	leftMouseHeld            bool
	priorLeftMouseClicked    bool
	leftMouseReleased        bool
	middleMouseScrollClicked bool
	rightMouseClicked        bool
	paused                   bool
	blueComputerActive       bool
	redComputerActive        bool
	hotRestart               bool
	drawFluxes               bool
	drawBlueTargets          bool
	drawRedTargets           bool
	drawDebugGrid            bool
	initialSelectedX         float32
	initialSelectedY         float32
	currentSelectedX         float32
	currentSelectedY         float32
}

func newGame() Game {
	g := Game{
		drawContext:        newDrawContext(),
		backgroundColor:    color.RGBA{200, 200, 200, 255},
		produceWait:        1_000_000 * time.Second,
		produceDelay:       3 * time.Second,
		slowTickTime:       time.Now(),
		priorTime:          time.Now(),
		currentTime:        time.Now(),
		initialSelectedX:   float32(-1.0),
		initialSelectedY:   float32(-1.0),
		currentSelectedX:   float32(-1.0),
		currentSelectedY:   float32(-1.0),
		blueComputerActive: true,
		redComputerActive:  true,
		hotRestart:         true,
	}
	return g
}

var logger l
var gameState state
var compOne computer
var compTwo computer
var toggles map[string]*toggle = make(map[string]*toggle)

// TODO: Should this be a method rather than a function?
func (g *Game) restart() {
	// TODO: Need to restart timers, etc
	gameState = *newState(g.windowWidth/g.gridSize, g.windowHeight/g.gridSize)
	compOne = *newComputer(bluePlayer)
	compTwo = *newComputer(redPlayer)
	setupState(gameState, complexMissile)
	g.updateComputerTargetNodes(gameState)
}

func (g *Game) Update() error {
	g.currentTime = time.Now()

	for _, toggle := range toggles {
		toggle.checkToggle(g.currentTime)
	}

	// refresh
	g.leftMouseReleased = false
	g.leftMouseHeld = false

	// update game state prior to player actions
	if !g.paused && !(slowTick && (g.currentTime.Sub(g.slowTickTime) < slowTickDuration)) {
		if slowTick {
			logger.log(infoLevel, alwaysLogType, "slowTick")
		}
		g.slowTickTime = time.Now()

		gameState.refresh()
		gameState.analyzeTargets()
		gameState.processTerminalPositions()
		gameState.attackUnits()
		gameState.attackNodes()
		gameState.moveUnits()

		if g.produceWait < g.produceDelay {
			g.produceWait = g.produceWait + g.currentTime.Sub(g.priorTime)
		} else {
			g.produceWait = 0
			gameState.produceUnits()
		}
	}

	// player actions
	var graphicalX, graphicalY = ebiten.CursorPosition()
	var x = float32(graphicalX) / g.windowMultiplier
	var y = float32(graphicalY) / g.windowMultiplier
	var pixelPos = pixelPosition{int(x), int(y)}

	g.priorLeftMouseClicked = g.leftMouseClicked
	g.leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	g.middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	g.rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	g.currentSelectedX = x
	g.currentSelectedY = y
	if g.priorLeftMouseClicked && !g.leftMouseClicked {
		g.leftMouseReleased = true
	} else if !g.priorLeftMouseClicked && g.leftMouseClicked {
		g.initialSelectedX = x
		g.initialSelectedY = y
	} else if g.leftMouseClicked {
		g.leftMouseHeld = true
	}
	if g.middleMouseScrollClicked {
		logger.log(infoLevel, mainLogType, "middle click")
	}
	if g.rightMouseClicked {
		// TODO: Determine behavior. Should player choose which layer they target or handle it via game logic?
		// TODO: It seems a bit hacky that I need to pass computer active here. Determine if this is right behavior.
		gameState.updateTargets(&position{ground, pixelPos}, g.blueComputerActive, g.redComputerActive)
	}

	// update game state after player actions
	if g.leftMouseReleased {
		gameState.updateSelected(
			pixelPosition{int(g.initialSelectedX), int(g.initialSelectedY)},
			pixelPosition{int(g.currentSelectedX), int(g.currentSelectedY)})
	}
	var nodeOwnershipChanged = gameState.updateNodeOwnership()

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if nodeOwnershipChanged {
		logger.log(infoLevel, mainLogType, "nodeOwnershipChanged")
		g.updateComputerTargetNodes(gameState)
	}

	var infantryCount = 0
	var sniperCount = 0
	var unguidedGroundMissileCount = 0
	var guidedAirMissileCount = 0
	var unknownCount = 0
	for _, u := range gameState.units {
		switch u.name {
		case infantry:
			infantryCount++
		case sniper:
			sniperCount++
		case unguidedGroundMissile:
			unguidedGroundMissileCount++
		case guidedAirMissile:
			guidedAirMissileCount++
		default:
			unknownCount++
		}

	}
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("total_units", len(gameState.units)))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("infantryCount", infantryCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("sniperCount", sniperCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("unguidedGroundMissile", unguidedGroundMissileCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("guidedAirMissile", guidedAirMissileCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("unknownCount", unknownCount))

	var blueNodeCount = 0
	var redNodeCount = 0

	for _, node := range gameState.nodes {
		switch node.player {
		case bluePlayer:
			blueNodeCount = blueNodeCount + 1
		case redPlayer:
			redNodeCount = redNodeCount + 1
		}
	}

	if blueNodeCount == 0 || redNodeCount == 0 {
		if runProfiler {
			pprof.StopCPUProfile()
		}
		if g.hotRestart {
			logger.log(infoLevel, mainLogType, "restartGame")

			g.restart()
		} else {
			return errors.New("regular termination")
		}
	}

	g.priorTime = g.currentTime
	return nil
}

func (g *Game) updateComputerTargetNodes(s state) {
	if g.blueComputerActive {
		compOne.updateTargetNodes(s)
	}
	if g.redComputerActive {
		compTwo.updateTargetNodes(s)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawContext.screen == nil {
		g.drawContext.screen = screen
	}

	screen.Fill(g.drawContext.backgroundColor)
	drawDebug(*g.drawContext, g.drawDebugGrid)
	gameState.draw(*g.drawContext, g.drawBlueTargets, g.drawRedTargets, g.drawFluxes)

	if g.leftMouseHeld {
		drawUnfilledBoxFloat32(
			*g.drawContext,
			g.drawContext.highlightColor,
			g.initialSelectedX,
			g.initialSelectedY,
			g.currentSelectedX,
			g.currentSelectedY,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(g.windowWidth) * g.windowMultiplier), int(float32(g.windowHeight) * g.windowMultiplier)
}

func main() {
	g := newGame()

	if runProfiler {
		f, err := os.Create("cpu.prof")
		if err != nil {
			logger.log(errorLevel, mainLogType, "could not create CPU profile")
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			logger.log(errorLevel, mainLogType, "could not start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	ebiten.SetWindowSize(int(float32(g.windowWidth)*g.windowMultiplier), int(float32(g.windowHeight)*g.windowMultiplier))

	logger = newLogger()

	toggles["pause"] = newToggle(ebiten.KeyP, func() { g.paused = !g.paused })
	toggles["slowTick"] = newToggle(ebiten.KeyT, func() { slowTick = !slowTick })
	toggles["logLevel"] = newToggle(ebiten.KeyL, func() { logger.toggleLevel() })
	toggles["mainLogType"] = newToggle(ebiten.KeyM, func() { logger.toggleType(mainLogType) })
	toggles["stateLogType"] = newToggle(ebiten.KeyS, func() { logger.toggleType(stateLogType) })
	toggles["selectedUnitLogType"] = newToggle(ebiten.KeyU, func() { logger.toggleType(selectedUnitLogType) })
	toggles["drawAll"] = newToggle(ebiten.KeyA, func() {
		var toggle = !(g.drawFluxes && g.drawBlueTargets && g.drawRedTargets && g.drawDebugGrid)
		g.drawFluxes = toggle
		g.drawBlueTargets = toggle
		g.drawRedTargets = toggle
		g.drawDebugGrid = toggle
	})
	toggles["drawFlux"] = newToggle(ebiten.KeyF, func() { g.drawFluxes = !g.drawFluxes })
	toggles["drawBlueTargets"] = newToggle(ebiten.KeyB, func() { g.drawBlueTargets = !g.drawBlueTargets })
	toggles["drawRedTargets"] = newToggle(ebiten.KeyR, func() { g.drawRedTargets = !g.drawRedTargets })
	toggles["drawDebugGrid"] = newToggle(ebiten.KeyG, func() { g.drawDebugGrid = !g.drawDebugGrid })

	g.restart()
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
