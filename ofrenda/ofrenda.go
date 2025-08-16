package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// This must only be used for drawing images and windows.
// Never use for game logic.
const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = 1280
const windowHeight = 720
const gridSize = 40
const offerSize = 14
const beltMiddleLineWidth = 2

var backgroundColor = color.RGBA{200, 200, 200, 0xff}

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

var beltDownImage *ebiten.Image
var beltDownLeftImage *ebiten.Image
var beltDownRightImage *ebiten.Image
var beltLeftImage *ebiten.Image
var beltLeftDownImage *ebiten.Image
var beltLeftUpImage *ebiten.Image
var beltRightImage *ebiten.Image
var beltRightDownImage *ebiten.Image
var beltRightUpImage *ebiten.Image
var beltUpImage *ebiten.Image
var beltUpLeftImage *ebiten.Image
var beltUpRightImage *ebiten.Image
var producerImage *ebiten.Image

var consumerSkullImage *ebiten.Image
var consumerCandyImage *ebiten.Image
var consumerFlowerImage *ebiten.Image

var skullImage *ebiten.Image
var candyImage *ebiten.Image
var flowerImage *ebiten.Image

var moverImage *ebiten.Image
var moverLongDownImage *ebiten.Image
var moverLongUpImage *ebiten.Image
var moverShortDownImage *ebiten.Image
var moverShortUpImage *ebiten.Image

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

	// belt images
	beltDownImage, _, _ = ebitenutil.NewImageFromFile("assets/beltDown.png")
	beltDownLeftImage, _, _ = ebitenutil.NewImageFromFile("assets/beltDownLeft.png")
	beltDownRightImage, _, _ = ebitenutil.NewImageFromFile("assets/beltDownRight.png")
	beltLeftImage, _, _ = ebitenutil.NewImageFromFile("assets/beltLeft.png")
	beltLeftDownImage, _, _ = ebitenutil.NewImageFromFile("assets/beltLeftDown.png")
	beltLeftUpImage, _, _ = ebitenutil.NewImageFromFile("assets/beltLeftUp.png")
	beltRightImage, _, _ = ebitenutil.NewImageFromFile("assets/beltRight.png")
	beltRightDownImage, _, _ = ebitenutil.NewImageFromFile("assets/beltRightDown.png")
	beltRightUpImage, _, _ = ebitenutil.NewImageFromFile("assets/beltRightUp.png")
	beltUpImage, _, _ = ebitenutil.NewImageFromFile("assets/beltUp.png")
	beltUpLeftImage, _, _ = ebitenutil.NewImageFromFile("assets/beltUpLeft.png")
	beltUpRightImage, _, _ = ebitenutil.NewImageFromFile("assets/beltUpRight.png")
	// consumer images
	consumerSkullImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerCircle.png")
	consumerCandyImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerCross.png")
	consumerFlowerImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerSquare.png")
	// producer images
	producerImage, _, _ = ebitenutil.NewImageFromFile("assets/producer.png")
	// offer images
	skullImage, _, _ = ebitenutil.NewImageFromFile("assets/circle.png")
	candyImage, _, _ = ebitenutil.NewImageFromFile("assets/cross.png")
	flowerImage, _, _ = ebitenutil.NewImageFromFile("assets/square.png")

	// game setup
	gameState = *newState()
}

var currentTime time.Time
var priorTime time.Time
var wait time.Duration = 1_000_000 * time.Second
var produceDelay time.Duration = 1 * time.Second

func (g *game) Update() error {
	currentTime = time.Now()

	if wait < produceDelay {
		wait = wait + currentTime.Sub(priorTime)
	} else {
		wait = 0
		gameState.produceOffers()
	}

	priorTime = currentTime
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(backgroundColor)
	gameState.draw(screen)
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
