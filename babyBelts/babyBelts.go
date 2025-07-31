package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const windowMultiplier = 2.5

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultiplier)
const windowHeight = int(720 * windowMultiplier)

const gridSize = int(40 * windowMultiplier)

var gridBottom [32][18]gridBottomImage
var gridTop [32][18]gridTopImage

// TODO: should these be string instead of int?
type gridBottomImage int
type gridTopImage int

const (
	NoBottomImage gridBottomImage = iota
	BeltDown
	BeltDownLeft
	BeltDownRight
	BeltLeft
	BeltLeftDown
	BeltLeftUp
	BeltRight
	BeltRightDown
	BeltRightUp
	BeltUp
	BeltUpLeft
	BeltUpRight
	ConsumerCircle
	ConsumerCross
	ConsumerSquare
	MoverDown
	MoverUp
	Producer
)

const (
	NoTopImage gridTopImage = iota
	Mover
	MoverLongDown
	MoverLongUp
	MoverShortDown
	MoverShortUp
)

type Game struct{}

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
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// counter clockwise circle
	gridBottom[1][1] = BeltLeftDown
	gridBottom[1][2] = BeltDown
	gridBottom[1][3] = BeltDown
	gridBottom[1][4] = BeltDownRight
	gridBottom[2][4] = BeltRight
	gridBottom[3][4] = BeltRight
	gridBottom[4][4] = BeltRightUp
	gridBottom[4][3] = BeltUp
	gridBottom[4][2] = BeltUp
	gridBottom[4][1] = BeltUpLeft
	gridBottom[3][1] = BeltLeft
	gridBottom[2][1] = BeltLeft

	// counter clockwise circle
	gridBottom[1][6] = BeltUpRight
	gridBottom[1][7] = BeltUp
	gridBottom[1][8] = BeltUp
	gridBottom[1][9] = BeltLeftUp
	gridBottom[2][9] = BeltLeft
	gridBottom[3][9] = BeltLeft
	gridBottom[4][9] = BeltDownLeft
	gridBottom[4][8] = BeltDown
	gridBottom[4][7] = BeltDown
	gridBottom[4][6] = BeltRightDown
	gridBottom[3][6] = BeltRight
	gridBottom[2][6] = BeltRight

	gridBottom[6][1] = BeltRight
	gridBottom[7][1] = BeltRight
	gridBottom[8][1] = BeltRight
	gridBottom[9][1] = BeltRight
	gridBottom[10][1] = BeltRight
	gridBottom[11][1] = BeltRight
	gridBottom[12][1] = BeltRight
	gridBottom[13][1] = BeltRight
	gridBottom[14][1] = BeltRight
	gridBottom[15][1] = BeltRight

	gridTop[6][2] = Mover
	gridTop[7][2] = MoverLongDown
	gridTop[8][1] = MoverLongUp
	gridTop[9][2] = MoverShortDown
	gridTop[10][1] = MoverShortUp
	gridBottom[11][2] = ConsumerCircle
	gridBottom[12][2] = ConsumerCross
	gridBottom[13][2] = ConsumerSquare
	gridBottom[14][2] = Producer

	gridBottom[6][3] = BeltRight
	gridBottom[7][3] = BeltRight
	gridBottom[8][3] = BeltRight
	gridBottom[9][3] = BeltRight
	gridBottom[10][3] = BeltRight
	gridBottom[11][3] = BeltRight
	gridBottom[12][3] = BeltRight
	gridBottom[13][3] = BeltRight
	gridBottom[14][3] = BeltRight
	gridBottom[15][3] = BeltRight

	for i := range len(gridBottom) {
		for j := range len(gridBottom[i]) {
			var xLeft = float32(gridSize * i)
			var yUp = float32(gridSize * j)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(windowMultiplier, windowMultiplier)
			op.GeoM.Translate(float64(xLeft), float64(yUp))

			var image *ebiten.Image
			switch gridBottom[i][j] {
			case BeltDown:
				image = beltDownImage
			case BeltDownLeft:
				image = beltDownLeftImage
			case BeltDownRight:
				image = beltDownRightImage
			case BeltLeft:
				image = beltLeftImage
			case BeltLeftDown:
				image = beltLeftDownImage
			case BeltLeftUp:
				image = beltLeftUpImage
			case BeltRight:
				image = beltRightImage
			case BeltRightDown:
				image = beltRightDownImage
			case BeltRightUp:
				image = beltRightUpImage
			case BeltUp:
				image = beltUpImage
			case BeltUpLeft:
				image = beltUpLeftImage
			case BeltUpRight:
				image = beltUpRightImage
			case ConsumerCircle:
				image = consumerCircleImage
			case ConsumerCross:
				image = consumerCrossImage
			case ConsumerSquare:
				image = consumerSquareImage
			case Producer:
				image = producerImage
			default:
				image = nil
			}

			if image != nil {
				screen.DrawImage(image, op)
			}
		}
	}

	for i := range len(gridTop) {
		for j := range len(gridTop[i]) {
			var xLeft = float32(gridSize * i)
			var yUp = float32(gridSize * j)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(windowMultiplier, windowMultiplier)
			op.GeoM.Translate(float64(xLeft), float64(yUp))

			var image *ebiten.Image
			switch gridTop[i][j] {
			case Mover:
				image = moverImage
			case MoverLongDown:
				image = moverLongDownImage
			case MoverLongUp:
				image = moverLongUpImage
			case MoverShortDown:
				image = moverShortDownImage
			case MoverShortUp:
				image = moverShortUpImage
			default:
				image = nil
			}

			if image != nil {
				screen.DrawImage(image, op)
			}
		}
	}

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(285*windowMultiplier), float64(45*windowMultiplier))
		screen.DrawImage(circleImage, op)
	}

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(265*windowMultiplier), float64(45*windowMultiplier))
		screen.DrawImage(squareImage, op)
	}

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(windowMultiplier, windowMultiplier)
		op.GeoM.Translate(float64(245*windowMultiplier), float64(45*windowMultiplier))
		screen.DrawImage(crossImage, op)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
