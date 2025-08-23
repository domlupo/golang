package main

import (
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// This must only be used for drawing images and windows.
// Never use for game logic.
const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = 1280
const windowHeight = 720
const gridSize = 40

var backgroundColor = color.RGBA{200, 200, 200, 255}

// player setup
var playerOneNodeColor = color.RGBA{0, 0, 255, 100}
var playerOneUnitColor = color.RGBA{0, 0, 255, 255}
var playerTwoNodeColor = color.RGBA{255, 0, 0, 100}
var playerTwoUnitColor = color.RGBA{255, 0, 0, 255}

// player actions
var selectectionBoxLineWidth = 2
var targetNodeLineWidth = 1
var highlightColor = color.RGBA{70, 255, 0, 200}

var unitWidth = 5
var unitHeight = 5

const debugOutput = true

// The debug grid helps debugging display issues. Usually only altering
// offset is needed, width can remain 1.
// Offset values should only be 0 - 20 (half of grid size without window multipler)
// Some helpful values are: 0, 0.5, 1, 19, 19.5, 20
const debugGridOn = true
const debugGridOffsetWithoutWindowMultipler = 0
const debugGridOffset = debugGridOffsetWithoutWindowMultipler * windowMultiplier
const debugGridWidth = 1 * windowMultiplier
const debugGridSize = 32

var debugGridColor = color.Black

type game struct{}

var gameState state

func init() {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	gameState = *newState()

	// setup player one
	var playerOneNodeOne = newNode(playerOne)
	var playerOneNodeTwo = newNode(playerOne)

	// setup player two
	var playerTwoNodeOne = newNode(playerTwo)
	var playerTwoNodeTwo = newNode(playerTwo)

	//setup game state
	gameState.addNode(gridPosition{1, 1}, playerOneNodeOne)
	gameState.addNode(gridPosition{1, 10}, playerOneNodeTwo)
	gameState.addNode(gridPosition{10, 1}, playerTwoNodeOne)
	gameState.addNode(gridPosition{10, 10}, playerTwoNodeTwo)
}

// globals
var currentTime time.Time
var priorTime time.Time
var playerSelected player = noPlayer
var produceWait time.Duration = 1_000_000 * time.Second
var produceDelay time.Duration = 1 * time.Second

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
	gameState.moveUnits()
	gameState.attackUnits()

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
	var gridPos, _ = splitPixelPosition(pixelPosition{int(x), int(y)})

	priorLeftMouseClicked = leftMouseClicked
	playerSelected = gameState.nodes[gridPos].player
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
		println("middle click")
	}
	if rightMouseClicked {
		gameState.updateTargets(&gridPos)
	}

	// update game state after player actions
	if leftMouseReleased {
		gameState.updateSelected(
			pixelPosition{int(initialSelectedX), int(initialSelectedY)},
			pixelPosition{int(currentSelectedX), int(currentSelectedY)})
	}

	// debug output
	if debugOutputWait < debugOutputDelay {
		debugOutputWait = debugOutputWait + currentTime.Sub(priorTime)
	} else if debugOutput == false {
		// print nothing
	} else {
		println(strings.Repeat("-", 50))
		println("debug output")
		println("total units: ", len(gameState.units))
		println("")
		debugOutputWait = 0
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
	ebiten.SetWindowSize(windowWidth*windowMultiplier, windowHeight*windowMultiplier)
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
