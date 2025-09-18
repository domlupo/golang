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

type Game struct {
	windowMultiplier      int
	windowWidth           int
	windowHeight          int
	backgroundColor       color.RGBA
	timerColor            color.RGBA
	doneColor             color.RGBA
	logFileName           string
	timerStarted          bool
	iDuration             int
	duration              time.Duration
	currentTime           time.Time
	startTime             time.Time
	progressPercent       float64
	progressPercentLogged int
	file                  *os.File
	err                   error
	mplusFaceSource       *text.GoTextFaceSource
}

func newGame() Game {
	g := Game{
		windowMultiplier: 2,
		windowWidth:      1280,
		windowHeight:     720,
		iDuration:        60,
		backgroundColor:  color.RGBA{50, 50, 50, 255},
		timerColor:       color.RGBA{75, 75, 75, 255},
		doneColor:        color.RGBA{0, 150, 100, 255},
		logFileName:      "logFocusTimer.txt",
	}
	return g
}

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
	g.currentTime = time.Now()

	if g.timerStarted {
		g.progressPercent = float64(g.currentTime.Sub(g.startTime)/time.Millisecond) / float64(g.duration/time.Millisecond) * 100
		if g.progressPercentLogged < 100 && g.progressPercentLogged != int(g.progressPercent) {
			g.progressPercentLogged = int(g.progressPercent)
			_, g.err = g.file.WriteString(fmt.Sprintf("%d%% %s\n", g.progressPercentLogged, g.currentTime.Format(time.RFC3339)))
			if g.err != nil {
				log.Fatal(g.err)
			}
		}
	} else {
		_, dy := ebiten.Wheel()
		g.iDuration = max(g.iDuration+int(dy), 0)

		if (ebiten.IsMouseButtonPressed(ebiten.MouseButton0) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton1) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton2)) && g.iDuration > 0 {
			g.timerStarted = true
			g.startTime = time.Now()
			g.duration = time.Duration(g.iDuration) * time.Minute
			_, g.err = g.file.WriteString(fmt.Sprintf("Timer set for %d minutes. startTime: %s\n", g.iDuration, g.startTime.Format(time.RFC3339)))
			if g.err != nil {
				log.Fatal(g.err)
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.timerStarted || g.currentTime.Sub(g.startTime) > g.duration {
		screen.Fill(g.doneColor)
	} else {
		screen.Fill(g.backgroundColor)
		vector.DrawFilledRect(
			screen,
			0,
			0,
			float32(float64(g.windowWidth*g.windowMultiplier)*g.progressPercent/float64(100)),
			float32(g.windowHeight*g.windowMultiplier),
			g.timerColor,
			false)
	}
	if !g.timerStarted {
		op := &text.DrawOptions{}
		text.Draw(screen, fmt.Sprintf("%d minutes. Change minutes with mouse wheel. Press any mouse button to start timer", g.iDuration), &text.GoTextFace{
			Source: g.mplusFaceSource,
			Size:   float64(24 * g.windowMultiplier),
		}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowWidth * g.windowMultiplier, g.windowHeight * g.windowMultiplier
}

func main() {
	g := newGame()

	g.file, g.err = os.Create(g.logFileName)
	if g.err != nil {
		log.Fatal(g.err)
	}
	defer g.file.Close()

	g.mplusFaceSource, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

	ebiten.SetWindowSize(g.windowWidth*g.windowMultiplier, g.windowHeight*g.windowMultiplier)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
