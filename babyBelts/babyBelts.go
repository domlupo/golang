package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultiplier)
const windowHeight = int(720 * windowMultiplier)

// HACK: and TODO: Is this better to put elsewere?
// The windowMultiplier should maybe not be visible here
// and instead only be seen in draw layer.
const gridSizeWithoutMultiplier = 40
const gridSize = (40 * windowMultiplier)
const shapeSize = 14
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
var circleImage *ebiten.Image
var consumerCircleImage *ebiten.Image
var consumerCrossImage *ebiten.Image
var consumerSquareImage *ebiten.Image
var crossImage *ebiten.Image
var producerImage *ebiten.Image
var squareImage *ebiten.Image

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
	consumerCircleImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerCircle.png")
	consumerCrossImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerCross.png")
	consumerSquareImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerSquare.png")
	// producer images
	producerImage, _, _ = ebitenutil.NewImageFromFile("assets/producer.png")
	// shape images
	circleImage, _, _ = ebitenutil.NewImageFromFile("assets/circle.png")
	crossImage, _, _ = ebitenutil.NewImageFromFile("assets/cross.png")
	squareImage, _, _ = ebitenutil.NewImageFromFile("assets/square.png")

	// mover images
	moverImage, _, _ = ebitenutil.NewImageFromFile("assets/mover.png")
	moverLongDownImage, _, _ = ebitenutil.NewImageFromFile("assets/moverLongDown.png")
	moverLongUpImage, _, _ = ebitenutil.NewImageFromFile("assets/moverLongUp.png")
	moverShortDownImage, _, _ = ebitenutil.NewImageFromFile("assets/moverShortDown.png")
	moverShortUpImage, _, _ = ebitenutil.NewImageFromFile("assets/moverShortUp.png")

	gameState = *newState()

	gameState.addMover(gridPosition{2, 0}, mover{short, mDown})
	gameState.addMover(gridPosition{3, 0}, mover{long, mDown})
	gameState.addMover(gridPosition{4, 2}, mover{short, mUp})
	gameState.addMover(gridPosition{5, 2}, mover{long, mUp})

	gameState.addBelt(gridPosition{2, 1}, bRight)
	gameState.addBelt(gridPosition{3, 1}, bRight)
	gameState.addBelt(gridPosition{4, 1}, bRight)
	gameState.addBelt(gridPosition{5, 1}, bRight)

	gameState.addShapeToMover(gridPosition{2, 0}, circle)
	gameState.addShapeToMover(gridPosition{3, 0}, circle)
	gameState.addShapeToMover(gridPosition{4, 2}, circle)
	gameState.addShapeToMover(gridPosition{5, 2}, circle)

	gameState.addMoversAndConsumers(gridPosition{10, 10})
	gameState.addMoversAndBelts(gridPosition{10, 2})
	gameState.addProducer(gridPosition{6, 5}, circle)
}

func (g *game) Update() error {
	gameState.moveShapes()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	gameState.draw(screen)
	drawDebugGrid(screen, debugGridOn)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
