package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const windowMultipler = 1

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultipler)
const windowHeight = int(720 * windowMultipler)

// standard card ratio is 5:7
const cardMultipler = windowMultipler
const cardHeight = 140 * cardMultipler
const cardWidth = 100 * cardMultipler
const cardBorder = 10 * cardMultipler
const cardNumberHeight = 30 * cardMultipler

type symbol int

const (
	Diamond symbol = iota
	Heart
	Spade
	Club
)

type value int

const (
	Ace value = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var img *ebiten.Image
var err error
var c *card

func init() {
	img, _, err = ebitenutil.NewImageFromFile("assets/solitaire_diamond_one.png")
	c = newCard(Diamond, Ace, *img)
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(100, 100)
	screen.DrawImage(&c.image, op)
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
