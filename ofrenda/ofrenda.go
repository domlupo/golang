package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"strings"
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

var backgroundColor = color.RGBA{200, 200, 200, 0xff}
var producerOutputRangeColor = color.RGBA{245, 40, 145, 200}
var producerOutputRangeLineWidth = 1 * windowMultiplier

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

var ofrendaImage *ebiten.Image

var bottleImage *ebiten.Image
var candyImage *ebiten.Image
var flowerImage *ebiten.Image
var skullImage *ebiten.Image

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
	// TODO: update this to ofrenda image
	ofrendaImage, _, _ = ebitenutil.NewImageFromFile("assets/consumerCircle.png")
	producerImage, _, _ = ebitenutil.NewImageFromFile("assets/producer.png")
	// offer images
	bottleImage, _, _ = ebitenutil.NewImageFromFile("assets/bottle.png")
	candyImage, _, _ = ebitenutil.NewImageFromFile("assets/candy.png")
	flowerImage, _, _ = ebitenutil.NewImageFromFile("assets/flower.png")
	skullImage, _, _ = ebitenutil.NewImageFromFile("assets/skull.png")

	// game setup
	gameState = *newState()

	// first producer
	var x = rand.IntN((windowWidth/gridSize)-4) + 2
	var y = rand.IntN((windowHeight/gridSize)-4) + 2
	gameState.addProducer(gridPosition{x, y}, bottle, rand.IntN(2)+1)

	// second producer
	x = rand.IntN((windowWidth/gridSize)-4) + 2
	y = rand.IntN((windowHeight/gridSize)-4) + 2
	gameState.addProducer(gridPosition{x, y}, candy, rand.IntN(2)+1)

	// third producer
	x = rand.IntN((windowWidth/gridSize)-4) + 2
	y = rand.IntN((windowHeight/gridSize)-4) + 2
	gameState.addProducer(gridPosition{x, y}, flower, rand.IntN(2)+1)

	// fourth producer
	x = rand.IntN((windowWidth/gridSize)-4) + 2
	y = rand.IntN((windowHeight/gridSize)-4) + 2
	gameState.addProducer(gridPosition{x, y}, skull, rand.IntN(2)+1)
}

var currentTime time.Time
var priorTime time.Time

// game vars
var produceWait time.Duration = 1_000_000 * time.Second
var produceDelay time.Duration = 1 * time.Second
var currentBeltDir = bDown
var currentRound = 0
var roundInProgress = false
var roundWait time.Duration = 1_000_000 * time.Second
var roundDelay time.Duration = 30 * time.Second
var gameOver = false
var drawProducerOutput = false
var consumerPlaced = false

// mouse button constants
var rightMouseInputWait time.Duration = 1_000_000 * time.Second
var rightMouseInputDelay time.Duration = 300 * time.Millisecond
var middleMouseInputWait time.Duration = 1_000_000 * time.Second
var middleMouseInputDelay time.Duration = 1 * time.Second
var leftMouseInputWait time.Duration = 1_000_000 * time.Second
var leftMouseInputDelay time.Duration = 10 * time.Millisecond

// debug constants
var debugOutputWait = 1_000_000 * time.Second
var debugOutputDelay = 1 * time.Second

func (g *game) Update() error {
	currentTime = time.Now()

	// game actions
	if roundInProgress {
		gameState.consumeoffers()
		gameState.moveoffers()

		if produceWait < produceDelay {
			produceWait = produceWait + currentTime.Sub(priorTime)
		} else {
			produceWait = 0
			gameState.produceOffers()
		}
	}

	// player actions
	var mouseX, mouseY = ebiten.CursorPosition()
	var x = float64(mouseX) / windowMultiplier
	var y = float64(mouseY) / windowMultiplier
	var gridPos, _ = splitPixelPosition(pixelPosition{int(x), int(y)})
	var leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	var middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	var rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	if leftMouseClicked && !consumerPlaced {
		gameState.addConsumer(&gridPos)
		consumerPlaced = true
	}

	if leftMouseInputWait < leftMouseInputDelay {
		leftMouseInputWait = leftMouseInputWait + currentTime.Sub(priorTime)
	} else if leftMouseClicked {
		leftMouseInputWait = 0
		gameState.addBelt(gridPos, currentBeltDir)
	}

	if rightMouseInputWait < rightMouseInputDelay {
		rightMouseInputWait = rightMouseInputWait + currentTime.Sub(priorTime)
		delete(gameState.belts, gridPos)
	} else if rightMouseClicked {
		drawProducerOutput = true
		rightMouseInputWait = 0
		currentBeltDir += 1
		if currentBeltDir == 13 {
			currentBeltDir = 1
		}
		gameState.addBelt(gridPosition{31, 17}, currentBeltDir)
	} else {
		drawProducerOutput = false
	}

	if middleMouseInputWait < middleMouseInputDelay {
		middleMouseInputWait = middleMouseInputWait + currentTime.Sub(priorTime)
	} else if middleMouseScrollClicked && !roundInProgress {
		currentRound += 1
		roundInProgress = true
		middleMouseInputWait = 0
	}

	if roundWait < roundDelay {
		roundWait = roundWait + currentTime.Sub(priorTime)
	} else {

		if currentRound == 1 && gameState.metrics[bottle] < 5 {
			gameOver = true
		} else if currentRound == 2 && (gameState.metrics[bottle] < 10 || gameState.metrics[candy] < 5) {
			gameOver = true
		}

		roundInProgress = false
		roundWait = 0
		clear(gameState.offers)
	}

	// debug output
	if debugOutputWait < debugOutputDelay {
		debugOutputWait = debugOutputWait + currentTime.Sub(priorTime)
	} else if debugOutput == false {
		// print nothing
	} else {
		println(strings.Repeat("-", 50))
		if gameOver {
			println("game is over")
		} else {
			println("game in progress")
		}
		if roundInProgress {
			println("round in progress")
		} else {
			println("round not in progress")
		}
		println("round: ", currentRound)
		println("round 1 goal:  5 total bottles")
		println("round 2 goal: 10 total bottles, 5 total candy")
		println("offers outstanding: ", len(gameState.offers))
		println("offers consumed:    ", gameState.metrics[bottle]+gameState.metrics[candy]+gameState.metrics[flower]+gameState.metrics[skull])
		fmt.Println("  bottles: ", gameState.metrics[bottle])
		fmt.Println("  candy:   ", gameState.metrics[candy])
		fmt.Println("  flower:  ", gameState.metrics[flower])
		fmt.Println("  skulls:  ", gameState.metrics[skull])
		println("")
		debugOutputWait = 0
	}

	priorTime = currentTime
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	gameState.draw(screen)
	if drawProducerOutput {
		gameState.drawProducerOutputRange(screen)
	}
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
