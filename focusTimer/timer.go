package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const windowMultiplier = 2

// aspect ratio config (16:9 aspect ratio)
const windowWidth = 1280
const windowHeight = 720

var backgroundColor = color.RGBA{50, 50, 50, 255}
var timerColor = color.RGBA{75, 75, 75, 255}
var doneColor = color.RGBA{0, 150, 100, 255}
var initialTimerMinutes = 60

var mplusFaceSource *text.GoTextFaceSource

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

	s, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	mplusFaceSource = s
}

var durationInt = initialTimerMinutes
var currentTime time.Time
var startTime time.Time
var duration time.Duration
var progressPercent float64
var progressPercentLogged int = 0
var timerStarted = false
var file *os.File
var err error

func (g *Game) Update() error {

	currentTime = time.Now()
	if timerStarted {
		progressPercent = float64(currentTime.Sub(startTime)/time.Millisecond) / float64(duration/time.Millisecond) * 100
		if progressPercentLogged < 100 && progressPercentLogged != int(progressPercent) {
			progressPercentLogged = int(progressPercent)
			_, err = file.WriteString(fmt.Sprintf("%d%% %s\n", progressPercentLogged, currentTime.Format(time.RFC3339)))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		_, dy := ebiten.Wheel()
		durationInt = max(durationInt+int(dy), 0)

		if (ebiten.IsMouseButtonPressed(ebiten.MouseButton0) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton1) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton2)) && durationInt > 0 {
			timerStarted = true
			startTime = time.Now()
			duration = time.Duration(durationInt) * time.Minute
			_, err = file.WriteString(fmt.Sprintf("Timer set for %d minutes. startTime: %s\n", durationInt, startTime.Format(time.RFC3339)))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !timerStarted || currentTime.Sub(startTime) > duration {
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
	if !timerStarted {
		var op = &text.DrawOptions{}
		text.Draw(screen, fmt.Sprintf("%d minutes. Change minutes with mouse wheel. Press any mouse button to start timer", durationInt), &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   24 * windowMultiplier,
		}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth * windowMultiplier, windowHeight * windowMultiplier
}

func main() {
	file, err = os.Create("logFocusTimer.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ebiten.SetWindowSize(windowWidth*windowMultiplier, windowHeight*windowMultiplier)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
