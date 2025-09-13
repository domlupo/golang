package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const windowMultiplier = 2.5

// aspect ratio config (16:9 aspect ratio)
const windowWidth = 1280
const windowHeight = 720

var backgroundColor = color.RGBA{50, 50, 50, 255}
var timerColor = color.RGBA{75, 75, 75, 255}
var doneColor = color.RGBA{0, 150, 100, 255}

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

var currentTime time.Time = time.Now()
var startTime time.Time = time.Now()
var duration time.Duration = 10 * time.Second
var progressPercent float64

func (g *Game) Update() error {
	currentTime = time.Now()
	progressPercent = float64(currentTime.Sub(startTime)/time.Millisecond) / float64(duration/time.Millisecond) * 100
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if currentTime.Sub(startTime) > duration {
		screen.Fill(doneColor)
	} else {
		screen.Fill(backgroundColor)
		vector.DrawFilledRect(
			screen,
			0,
			0,
			float32(windowWidth*progressPercent/100)*windowMultiplier,
			windowHeight*windowMultiplier,
			timerColor,
			false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth * windowMultiplier, windowHeight * windowMultiplier
}

func main() {
	ebiten.SetWindowSize(windowWidth*windowMultiplier, windowHeight*windowMultiplier)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
