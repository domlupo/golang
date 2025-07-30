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

var grid [32][18][]gridImage

// TODO: should this be string instead of int?
type gridImage int

const (
	NoImage  gridImage = iota
	BeltDown           = iota
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
var moverDownImage *ebiten.Image
var moverUpImage *ebiten.Image
var producerImage *ebiten.Image
var squareImage *ebiten.Image

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
	// mover images
	moverDownImage, _, _ = ebitenutil.NewImageFromFile("assets/moverDown.png")
	moverUpImage, _, _ = ebitenutil.NewImageFromFile("assets/moverUp.png")
	// producer images
	producerImage, _, _ = ebitenutil.NewImageFromFile("assets/producer.png")
	// shape images
	circleImage, _, _ = ebitenutil.NewImageFromFile("assets/circle.png")
	crossImage, _, _ = ebitenutil.NewImageFromFile("assets/cross.png")
	squareImage, _, _ = ebitenutil.NewImageFromFile("assets/square.png")
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// counter clockwise circle
	grid[1][1] = []gridImage{BeltLeftDown}
	grid[1][2] = []gridImage{BeltDown}
	grid[1][3] = []gridImage{BeltDown}
	grid[1][4] = []gridImage{BeltDownRight}
	grid[2][4] = []gridImage{BeltRight}
	grid[3][4] = []gridImage{BeltRight}
	grid[4][4] = []gridImage{BeltRightUp}
	grid[4][3] = []gridImage{BeltUp}
	grid[4][2] = []gridImage{BeltUp}
	grid[4][1] = []gridImage{BeltUpLeft}
	grid[3][1] = []gridImage{BeltLeft}
	grid[2][1] = []gridImage{BeltLeft}

	// counter clockwise circle
	grid[1][6] = []gridImage{BeltUpRight}
	grid[1][7] = []gridImage{BeltUp}
	grid[1][8] = []gridImage{BeltUp}
	grid[1][9] = []gridImage{BeltLeftUp}
	grid[2][9] = []gridImage{BeltLeft}
	grid[3][9] = []gridImage{BeltLeft}
	grid[4][9] = []gridImage{BeltDownLeft}
	grid[4][8] = []gridImage{BeltDown}
	grid[4][7] = []gridImage{BeltDown}
	grid[4][6] = []gridImage{BeltRightDown}
	grid[3][6] = []gridImage{BeltRight}
	grid[2][6] = []gridImage{BeltRight}

	grid[6][1] = []gridImage{BeltRight}
	grid[7][1] = []gridImage{BeltRight}
	grid[8][1] = []gridImage{BeltRight}
	grid[9][1] = []gridImage{BeltRight, MoverDown}
	grid[10][1] = []gridImage{BeltRight, MoverUp}
	grid[11][1] = []gridImage{BeltRight}
	grid[12][1] = []gridImage{BeltRight}
	grid[13][1] = []gridImage{BeltRight}
	grid[14][1] = []gridImage{BeltRight}
	grid[15][1] = []gridImage{BeltRight}

	grid[7][2] = []gridImage{MoverDown}
	grid[8][2] = []gridImage{MoverUp}
	// skip two
	grid[11][2] = []gridImage{ConsumerCircle}
	grid[12][2] = []gridImage{ConsumerCross}
	grid[13][2] = []gridImage{ConsumerSquare}
	grid[14][2] = []gridImage{Producer}

	for i := range len(grid) {
		for j := range len(grid[i]) {
			for k := range len(grid[i][j]) {
				var xLeft = float32(gridSize * i)
				var yUp = float32(gridSize * j)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(windowMultiplier, windowMultiplier)
				op.GeoM.Translate(float64(xLeft), float64(yUp))

				var image *ebiten.Image
				switch grid[i][j][k] {
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
				case MoverUp:
					image = moverUpImage
				case MoverDown:
					image = moverDownImage
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
