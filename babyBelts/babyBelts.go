package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const windowMultipler = 2.5

// 16:9 aspect ratio
const windowWidth = int(1280 * windowMultipler)
const windowHeight = int(720 * windowMultipler)

const gridSize = int(40 * windowMultipler)

var grid [32][18]int

type Game struct{}

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

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	grid[0][0] = 1
	grid[1][1] = 1
	grid[1][2] = 1
	grid[1][17] = 1
	grid[31][17] = 1
	for i := range len(grid) {
		for j := range len(grid[i]) {
			var xLeft = float32(gridSize * i)
			var xRight = float32(gridSize * (i + 1))
			var yUp = float32(gridSize * j)
			var yDown = float32(gridSize * (j + 1))

			if grid[i][j] == 1 {
				vector.DrawFilledRect(screen, xLeft, yDown, xRight-xLeft, yUp-yDown, color.White, false)
			}
		}
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
