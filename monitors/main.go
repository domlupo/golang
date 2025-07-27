package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func init() {
	var monitors []*ebiten.MonitorType
	fmt.Println(monitors)
	fmt.Println(len(monitors))
	ebiten.AppendMonitors(monitors)
	fmt.Println(monitors)
	fmt.Println(len(monitors))
	fmt.Println(ebiten.Monitor())
}

func (g *Game) Update() error {
	return nil
}
