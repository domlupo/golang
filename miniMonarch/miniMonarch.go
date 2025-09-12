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

// logging config
const outputSummaryLogs = true

// This must only be used for drawing images and windows.
// Never use for game logic.
const windowMultiplier = 2.5

// aspect ratio config (16:9 aspect ratio)
const windowWidth = 1280
const windowHeight = 720
const gridSize = 40

// game balance config
var fluxMagnitude = 5
var unitFluxCountLimit = 10
var nodeRadius = gridSize / 2
var initialNodeHealth = 100

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	radius:    5,
	maxHealth: 500,
	attack:    1,
	vision:    40,
	reach:     13,
}
var sniperConfig = unitConfig{
	radius:    7,
	maxHealth: 5,
	attack:    1,
	vision:    100,
	reach:     100,
}

// game graphic config
var selectectionBoxLineWidth = 1
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

const debugGridOffsetWithoutWindowMultipler = 0
const debugGridOffset = debugGridOffsetWithoutWindowMultipler * windowMultiplier
const debugGridWidth = 1 * windowMultiplier
const debugGridSquareCount = 32

// other configs
var runProfiler = false

const restart = true

// general
var currentTime time.Time
var priorTime time.Time = time.Now()

// general debugging
var paused = false

// logger debugging
var debugOutputWait = 1_000_000 * time.Second
var debugOutputDelay = 1 * time.Second

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

type game struct{}

var logger l
var gameState state
var compOne computer
var compTwo computer
var toggles map[string]*toggle = make(map[string]*toggle)

func restartGame() {
	// TODO: Need to restart timers, etc
	gameState = *newState()
	compOne = *newComputer(bluePlayer)
	compTwo = *newComputer(redPlayer)
	setupState(gameState, complexSniper)
	updateComputerTargetNodes(gameState)
}

func init() {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	logger = newLogger()

	toggles["pause"] = newToggle(ebiten.KeyP, func() { paused = !paused })
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

	restartGame()
}

func (g *game) Update() error {
	currentTime = time.Now()

	for _, toggle := range toggles {
		toggle.checkToggle(currentTime)
	}

	// set global update defaults
	leftMouseReleased = false
	leftMouseHeld = false

	// update game state prior to player actions
	if !paused {
		gameState.refresh()
		gameState.analyzeTargets()
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
	var x = float32(graphicalX) / windowMultiplier
	var y = float32(graphicalY) / windowMultiplier
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
		gameState.updateTargets(&pixelPos)
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
		updateComputerTargetNodes(gameState)
	}

	// debug output
	if debugOutputWait < debugOutputDelay {
		debugOutputWait = debugOutputWait + currentTime.Sub(priorTime)
	} else if outputSummaryLogs {
		logger.log(
			infoLevel,
			mainLogType,
			"game_state",
			slog.Int("total_units",
				len(gameState.units)))
		debugOutputWait = 0
	}

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

			restartGame()
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

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	drawDebug(screen, drawDebugGrid)
	gameState.draw(screen, drawBlueTargets, drawRedTargets, drawFluxes)

	if leftMouseHeld {
		drawUnfilledBoxFloat32(
			screen,
			highlightColor,
			initialSelectedX,
			initialSelectedY,
			currentSelectedX,
			currentSelectedY,
		)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth * windowMultiplier, windowHeight * windowMultiplier
}

func main() {
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

	ebiten.SetWindowSize(windowWidth*windowMultiplier, windowHeight*windowMultiplier)
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
