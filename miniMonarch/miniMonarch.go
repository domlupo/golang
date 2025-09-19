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

// computer config
var blueComputerActive = true
var redComputerActive = true

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

// game graphic config
var selectectionBoxLineWidth = 1
var missileBlastRadiusLineWidth = 1
var healthLineWidth = 2
var backgroundColor = color.RGBA{200, 200, 200, 255}
var highlightColor = color.RGBA{70, 255, 0, 200}
var unitHealthColor = highlightColor
var nodeHealthColor = highlightColor
var bluePlayerNodeColor = color.RGBA{0, 0, 255, 100}
var bluePlayerUnitColor = color.RGBA{0, 0, 255, 100}
var redPlayerNodeColor = color.RGBA{255, 0, 0, 100}
var redPlayerUnitColor = color.RGBA{255, 0, 0, 100}
var neutralNodeColor = color.RGBA{100, 60, 100, 200}

// debug only configs
var drawFluxes = false
var drawUnitHealth = false
var drawBlueTargets = false
var drawRedTargets = false
var fluxLineWidth = 2
var fluxDirectionRadius = 3
var fluxColor = highlightColor
var targetNodeLineWidth = 1

// debug grid config
// The debug grid helps debugging display issues. Usually only altering
// offset is needed, width can remain 1.
// Offset values should only be 0 - 20 (half of grid size without window multipler)
// Some helpful values are: 0, 0.5, 1, 19, 19.5, 20
var drawDebugGrid = false
var debugGridColor = color.Black

const debugGridSquareCount = 32

// other configs
var runProfiler = false

const restart = true

// general
var slowTickTime time.Time = time.Now()
var priorTime time.Time = time.Now()
var currentTime time.Time = time.Now()

// general debugging
var paused = false

// globals
var produceWait time.Duration = 1_000_000 * time.Second
var produceDelay time.Duration = 3 * time.Second
var leftMouseClicked = false
var leftMouseHeld = false
var priorLeftMouseClicked = false
var leftMouseReleased = false
var middleMouseScrollClicked = false
var rightMouseClicked = false
var initialSelectedX = float32(-1.0)
var initialSelectedY = float32(-1.0)
var currentSelectedX = float32(-1.0)
var currentSelectedY = float32(-1.0)

type Game struct {
	// float32 rather than float64 since alot of Ebiten draw APIs use float32
	*drawContext
	backgroundColor color.RGBA
}

func newGame() Game {
	g := Game{
		drawContext:     newDrawContext(),
		backgroundColor: color.RGBA{200, 200, 200, 255},
	}
	return g
}

var logger l
var gameState state
var compOne computer
var compTwo computer
var toggles map[string]*toggle = make(map[string]*toggle)

// TODO: Should this be a method rather than a function?
func restartGame(g *Game) {
	// TODO: Need to restart timers, etc
	gameState = *newState(g.windowWidth/g.gridSize, g.windowHeight/g.gridSize)
	compOne = *newComputer(bluePlayer)
	compTwo = *newComputer(redPlayer)
	setupState(gameState, complexMissile)
	updateComputerTargetNodes(gameState)
}

func (g *Game) Update() error {
	currentTime = time.Now()

	for _, toggle := range toggles {
		toggle.checkToggle(currentTime)
	}

	// set global update defaults
	leftMouseReleased = false
	leftMouseHeld = false

	// update game state prior to player actions
	if !paused && !(slowTick && (currentTime.Sub(slowTickTime) < slowTickDuration)) {
		if slowTick {
			logger.log(infoLevel, alwaysLogType, "slowTick")
		}
		slowTickTime = time.Now()

		gameState.refresh()
		gameState.analyzeTargets()
		gameState.processTerminalPositions()
		gameState.attackUnits()
		gameState.attackNodes()
		gameState.moveUnits()

		if produceWait < produceDelay {
			produceWait = produceWait + currentTime.Sub(priorTime)
		} else {
			produceWait = 0
			gameState.produceUnits()
		}
	}

	// player actions
	var graphicalX, graphicalY = ebiten.CursorPosition()
	var x = float32(graphicalX) / g.windowMultiplier
	var y = float32(graphicalY) / g.windowMultiplier
	var pixelPos = pixelPosition{int(x), int(y)}

	priorLeftMouseClicked = leftMouseClicked
	leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	currentSelectedX = x
	currentSelectedY = y
	if priorLeftMouseClicked && !leftMouseClicked {
		leftMouseReleased = true
	} else if !priorLeftMouseClicked && leftMouseClicked {
		initialSelectedX = x
		initialSelectedY = y
	} else if leftMouseClicked {
		leftMouseHeld = true
	}
	if middleMouseScrollClicked {
		logger.log(infoLevel, mainLogType, "middle click")
	}
	if rightMouseClicked {
		// TODO: Determine behavior. Should player choose which layer they target or handle it via game logic?
		gameState.updateTargets(&position{ground, pixelPos})
	}

	// update game state after player actions
	if leftMouseReleased {
		gameState.updateSelected(
			pixelPosition{int(initialSelectedX), int(initialSelectedY)},
			pixelPosition{int(currentSelectedX), int(currentSelectedY)})
	}
	var nodeOwnershipChanged = gameState.updateNodeOwnership()

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if nodeOwnershipChanged {
		logger.log(infoLevel, mainLogType, "nodeOwnershipChanged")
		updateComputerTargetNodes(gameState)
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
		if restart {
			logger.log(infoLevel, mainLogType, "restartGame")

			restartGame(g)
		} else {
			return errors.New("regular termination")
		}
	}

	priorTime = currentTime
	return nil
}

func updateComputerTargetNodes(s state) {
	if blueComputerActive {
		compOne.updateTargetNodes(s)
	}
	if redComputerActive {
		compTwo.updateTargetNodes(s)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawContext.screen == nil {
		g.drawContext.screen = screen
	}

	screen.Fill(backgroundColor)
	drawDebug(*g.drawContext, drawDebugGrid)
	gameState.draw(*g.drawContext, drawBlueTargets, drawRedTargets, drawFluxes)

	if leftMouseHeld {
		drawUnfilledBoxFloat32(
			*g.drawContext,
			highlightColor,
			initialSelectedX,
			initialSelectedY,
			currentSelectedX,
			currentSelectedY,
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

	toggles["pause"] = newToggle(ebiten.KeyP, func() { paused = !paused })
	toggles["slowTick"] = newToggle(ebiten.KeyT, func() { slowTick = !slowTick })
	toggles["logLevel"] = newToggle(ebiten.KeyL, func() { logger.toggleLevel() })
	toggles["mainLogType"] = newToggle(ebiten.KeyM, func() { logger.toggleType(mainLogType) })
	toggles["stateLogType"] = newToggle(ebiten.KeyS, func() { logger.toggleType(stateLogType) })
	toggles["selectedUnitLogType"] = newToggle(ebiten.KeyU, func() { logger.toggleType(selectedUnitLogType) })
	toggles["drawAll"] = newToggle(ebiten.KeyA, func() {
		var toggle = !(drawFluxes && drawUnitHealth && drawBlueTargets && drawRedTargets && drawDebugGrid)
		drawFluxes = toggle
		drawUnitHealth = toggle
		drawBlueTargets = toggle
		drawRedTargets = toggle
		drawDebugGrid = toggle
	})
	toggles["drawFlux"] = newToggle(ebiten.KeyF, func() { drawFluxes = !drawFluxes })
	toggles["drawUnitHealth"] = newToggle(ebiten.KeyH, func() { drawUnitHealth = !drawUnitHealth })
	toggles["drawBlueTargets"] = newToggle(ebiten.KeyB, func() { drawBlueTargets = !drawBlueTargets })
	toggles["drawRedTargets"] = newToggle(ebiten.KeyR, func() { drawRedTargets = !drawRedTargets })
	toggles["drawDebugGrid"] = newToggle(ebiten.KeyG, func() { drawDebugGrid = !drawDebugGrid })

	restartGame(&g)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
