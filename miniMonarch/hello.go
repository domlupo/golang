package main

import (
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// This must only be used for drawing images and windows.
// Never use for game logic.
const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = 1280
const windowHeight = 720
const gridSize = 40

var backgroundColor = color.RGBA{200, 200, 200, 0xff}

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

type direction int

const (
	right direction = iota
	left
	up
	down
)

type game struct{}

func init() {
	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}
}

var currentTime time.Time
var priorTime time.Time

// debug constants
var debugOutputWait = 1_000_000 * time.Second
var debugOutputDelay = 1 * time.Second

func (g *game) Update() error {
	currentTime = time.Now()

	// player actions
	var mouseX, mouseY = ebiten.CursorPosition()
	var x = float64(mouseX) / windowMultiplier
	var y = float64(mouseY) / windowMultiplier
	var _, _ = splitPixelPosition(pixelPosition{int(x), int(y)})
	var leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	var middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	var rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	if leftMouseClicked {
		println("left click")
	} else if middleMouseScrollClicked {
		println("middle click")
	} else if rightMouseClicked {
		println("right click")
	}

	// debug output
	if debugOutputWait < debugOutputDelay {
		debugOutputWait = debugOutputWait + currentTime.Sub(priorTime)
	} else if debugOutput == false {
		// print nothing
	} else {
		println(strings.Repeat("-", 50))
		println("debug output")
		println("")
		debugOutputWait = 0
	}

	priorTime = currentTime
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	drawDebugGrid(screen, debugGridOn)
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

func drawDebugGrid(screen *ebiten.Image, drawDebugGrid bool) {
	if !drawDebugGrid {
		return
	}

	for i := range debugGridSize {
		var x = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(screen, x-debugGridOffset, 0, x-debugGridOffset, float32(windowHeight*windowMultiplier), debugGridWidth, debugGridColor, false)
		vector.StrokeLine(screen, x+debugGridOffset, 0, x+debugGridOffset, float32(windowHeight*windowMultiplier), debugGridWidth, debugGridColor, false)
		var y = float32(gridSize * windowMultiplier * i)
		vector.StrokeLine(screen, 0, y-debugGridOffset, float32(windowWidth*windowMultiplier), y-debugGridOffset, debugGridWidth, debugGridColor, false)
		vector.StrokeLine(screen, 0, y+debugGridOffset, float32(windowWidth*windowMultiplier), y+debugGridOffset, debugGridWidth, debugGridColor, false)
	}
}
