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

var runProfiler = false

// This must only be used for drawing images and windows.
// Never use for game logic.
const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = 1280
const windowHeight = 720
const gridSize = 40

// graphics setup
var backgroundColor = color.RGBA{200, 200, 200, 255}
var highlightColor = color.RGBA{70, 255, 0, 200}
var selectectionBoxLineWidth = 1
var targetNodeLineWidth = 1
var healthLineWidth = 2

// computer setup
var blueComputerActive = true
var redComputerActive = true
var bluePlayerActive = !blueComputerActive
var redPlayerActive = !redComputerActive

// player setup
var bluePlayerNodeColor = color.RGBA{0, 0, 255, 100}
var bluePlayerUnitColor = color.RGBA{0, 0, 255, 100}
var redPlayerNodeColor = color.RGBA{255, 0, 0, 100}
var redPlayerUnitColor = color.RGBA{255, 0, 0, 100}
var neutralNodeColor = color.RGBA{100, 60, 100, 200}

// unit and node setup
// Note: unitRadius set to an odd number is better for testing due to division bugs
var unitRadius = 5
var nodeRadius = gridSize / 2
var initialNodeHealth = 100

// The debug grid helps debugging display issues. Usually only altering
// offset is needed, width can remain 1.
// Offset values should only be 0 - 20 (half of grid size without window multipler)
// Some helpful values are: 0, 0.5, 1, 19, 19.5, 20
const logLevel = slog.LevelInfo
const filterTimeAndLevelFromLogs = true
const debugOutput = true
const debugGridOn = true
const debugGridOffsetWithoutWindowMultipler = 0
const debugGridOffset = debugGridOffsetWithoutWindowMultipler * windowMultiplier
const debugGridWidth = 1 * windowMultiplier
const debugGridSize = 32

var debugGridColor = color.Black

type game struct{}

var logger *slog.Logger
var gameState state
var compOne computer
var compTwo computer

func init() {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	logOptions := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, attribute slog.Attr) slog.Attr {
			if filterTimeAndLevelFromLogs && (attribute.Key == "time" || attribute.Key == "level") {
				return slog.Attr{}
			} else {
				return attribute
			}
		},
	}

	logger = slog.New(slog.NewJSONHandler(os.Stderr, logOptions))
	gameState = *newState()
	compOne = *newComputer(bluePlayer)
	compTwo = *newComputer(redPlayer)
	setupState(gameState, complex)
}

// globals
var currentTime time.Time
var priorTime time.Time
var playerSelected player = noPlayer
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

// debug constants
var debugOutputWait = 1_000_000 * time.Second
var debugOutputDelay = 1 * time.Second

func (g *game) Update() error {
	// set global update defaults
	currentTime = time.Now()
	leftMouseReleased = false
	leftMouseHeld = false

	// update game state prior to player actions
	gameState.refreshUnits()
	gameState.moveUnits()
	gameState.attackUnits()
	gameState.attackNodes()

	if produceWait < produceDelay {
		produceWait = produceWait + currentTime.Sub(priorTime)
	} else {
		produceWait = 0
		gameState.produceUnits()
	}

	// player actions
	var graphicalX, graphicalY = ebiten.CursorPosition()
	var x = float32(graphicalX) / windowMultiplier
	var y = float32(graphicalY) / windowMultiplier
	var pixelPos = pixelPosition{int(x), int(y)}
	var gridPos, _ = splitPixelPosition(pixelPos)

	priorLeftMouseClicked = leftMouseClicked
	// TODO: Fix player selected. Reuse collision algorithm or is that overkill?
	if gameState.nodes[gridPos] != nil {
		playerSelected = gameState.nodes[gridPos].player
	}
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
		logger.Info("middle click")
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
	gameState.updateNodeOwnership()

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if blueComputerActive {
		compOne.updateActions(gameState)
	}
	if redComputerActive {
		compTwo.updateActions(gameState)
	}

	// debug output
	if debugOutputWait < debugOutputDelay {
		debugOutputWait = debugOutputWait + currentTime.Sub(priorTime)
	} else if debugOutput == false {
		// log nothing
	} else {
		logger.Info(
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
		return errors.New("regular termination")
	}

	priorTime = currentTime
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	drawDebugGrid(screen, debugGridOn)
	gameState.draw(screen, playerSelected)

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
			logger.Error("could not create CPU profile")
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Error("could not start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	ebiten.SetWindowSize(windowWidth*windowMultiplier, windowHeight*windowMultiplier)
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
