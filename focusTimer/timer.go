package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const maxPercentage = 100
const minutesInADay = 1440

type Game struct {
	windowMultiplier            int
	windowWidth                 int
	windowHeight                int
	backgroundColor             color.RGBA
	timerColor                  color.RGBA
	doneColor                   color.RGBA
	logFileName                 string
	timerStarted                bool
	intDurationStart            int
	intDurationStartFromLog     int
	intDurationCount            int
	duration                    time.Duration
	currentTime                 time.Time
	startTime                   time.Time
	logMinsElapsedThresholdMins int
	logPercentPadding           int
	logMinutePadding            int
	minuteTicker                time.Ticker
	progressPercent             float64
	progressPercentStart        float64
	intProgressPercent          int
	file                        *os.File
	progressLoadedFromLogFile   bool
	err                         error
	mplusFaceSource             *text.GoTextFaceSource
}

func newGame() Game {
	g := Game{
		windowMultiplier:            2,
		windowWidth:                 1280,
		windowHeight:                720,
		intDurationStart:            60,
		logMinsElapsedThresholdMins: 100,
		// Largest percent logged (100.0) char length is 5.
		logPercentPadding: 5,
		// Largest minute log (1440) char length is 4.
		logMinutePadding: 4,
		backgroundColor:  color.RGBA{50, 50, 50, 255},
		timerColor:       color.RGBA{75, 75, 75, 255},
		doneColor:        color.RGBA{0, 150, 100, 255},
		logFileName:      "logFocusTimer.txt",
	}
	defer g.minuteTicker.Stop()
	return g
}

func (g *Game) Update() error {
	g.currentTime = time.Now()

	if !g.timerStarted {
		_, dy := ebiten.Wheel()
		g.intDurationStart = max(g.intDurationStart+int(dy), 1)
		g.intDurationStart = min(g.intDurationStart, minutesInADay)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton1) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			g.startTimer()
		}
	}

	if g.timerStarted && g.intProgressPercent < maxPercentage {
		g.progressPercent = float64(g.currentTime.Sub(g.startTime)/time.Millisecond) / float64(g.duration/time.Millisecond) * 100

		if g.intProgressPercent != int(g.progressPercent) {
			g.intProgressPercent = int(g.progressPercent)
			loggedProgressPercent := int((float64(g.intProgressPercent) * (float64(maxPercentage) - g.progressPercentStart) / maxPercentage) + g.progressPercentStart)
			_, g.err = g.file.WriteString(fmt.Sprintf(
				"%*d%% %s %*d minutes remaining\n",
				g.logPercentPadding,
				loggedProgressPercent,
				g.currentTime.Format(time.RFC3339),
				g.logMinutePadding,
				g.intDurationCount))
			if g.err != nil {
				log.Fatal(g.err)
			}
		}

		go func() {
			for {
				select {
				case <-g.minuteTicker.C:
					g.intDurationCount = g.intDurationCount - 1
					loggedProgressPercent := (g.progressPercent * (float64(maxPercentage) - g.progressPercentStart) / maxPercentage) + g.progressPercentStart
					if g.intDurationStart > g.logMinsElapsedThresholdMins {
						_, g.err = g.file.WriteString(fmt.Sprintf(
							"%*.1f%% %s %*d minutes remaining\n",
							g.logPercentPadding,
							loggedProgressPercent,
							g.currentTime.Format(time.RFC3339),
							g.logMinutePadding,
							g.intDurationCount))
					}
				default:
					time.Sleep(1 * time.Microsecond)
				}
			}
		}()
	}

	if g.intProgressPercent == maxPercentage {
		g.file.Close()
	}

	return nil
}

func (g *Game) startTimer() {
	if !g.progressLoadedFromLogFile || g.intDurationStart != g.intDurationStartFromLog {
		g.progressPercentStart = 0
		g.file, g.err = os.Create(g.logFileName)
		if g.err != nil {
			log.Fatal(g.err)
		}
	}

	g.timerStarted = true
	g.startTime = time.Now()
	g.duration = time.Duration(g.intDurationStart) * time.Minute
	g.intDurationCount = g.intDurationStart
	g.minuteTicker = *time.NewTicker(1 * time.Minute)
	_, g.err = g.file.WriteString(
		fmt.Sprintf("Timer set for %d minutes. \nstartTime: %s\n", g.intDurationStart, g.startTime.Format(time.RFC3339)))
	if g.err != nil {
		log.Fatal(g.err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.timerStarted || g.currentTime.Sub(g.startTime) > g.duration {
		screen.Fill(g.doneColor)
	} else {
		screen.Fill(g.backgroundColor)
		// Use int not percent read from log for prior timer progress since timer has 1 minute granularity.
		priorTimerProgress := g.progressPercentStart / 100
		currentTimerProgess := (g.progressPercent * (maxPercentage - g.progressPercentStart) / maxPercentage) / 100
		vector.DrawFilledRect(
			screen,
			0,
			0,
			float32(float64(g.windowWidth*g.windowMultiplier)*(priorTimerProgress+currentTimerProgess)),
			float32(g.windowHeight*g.windowMultiplier),
			g.timerColor,
			false)
	}

	if !g.timerStarted {
		op := &text.DrawOptions{}
		text.Draw(screen, fmt.Sprintf("%d minutes. Change minutes with mouse wheel. Press any mouse button to start timer", g.intDurationStart), &text.GoTextFace{
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

	g.mplusFaceSource, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	ebiten.SetWindowSize(g.windowWidth*g.windowMultiplier, g.windowHeight*g.windowMultiplier)

	intDurationStartFromCLI := readMinsFromCLI()
	file, progressPercentFromLog, intDurationStartFromLog := readProgressFromLog(g.logFileName)

	if intDurationStartFromCLI > 0 {
		g.intDurationStart = intDurationStartFromCLI
		g.startTimer()
	} else if intDurationStartFromLog > 0 && progressPercentFromLog < 100 {
		g.file = file
		g.progressLoadedFromLogFile = true
		g.progressPercentStart = progressPercentFromLog
		g.intDurationStartFromLog = intDurationStartFromLog
		g.intDurationStart = intDurationStartFromLog
	}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

	g.file.Close()
}

func readMinsFromCLI() int {
	args := os.Args
	if len(args) > 2 {
		panic("Invalid args passed, only a single arg is expected.")
	}

	if len(args) == 2 {
		mins, convertMinsErr := strconv.Atoi(os.Args[1])

		if convertMinsErr != nil {
			panic("Invalid arg passed, an integer is expected.")
		}
		if mins < 1 || mins > minutesInADay {
			panic("Invalid arg passed, expected a number between 1 and 1440.")
		}

		return mins
	}

	return 0
}

func readProgressFromLog(filename string) (*os.File, float64, int) {
	file, fileErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	var progressPercentFromLog float64
	var intDurationStartFromLog int
	if fileErr == nil {
		scanner := bufio.NewReader(file)
		for {
			textLine, readLineErr := scanner.ReadString('\n')
			if readLineErr != nil {
				break
			}

			percent := strings.Split(textLine, "%")[0]
			percent = strings.TrimSpace(percent)
			if loggedProgressPercent, convertErr := strconv.ParseFloat(percent, 64); convertErr == nil {
				progressPercentFromLog = loggedProgressPercent
			}

			text := strings.Split(textLine, " ")
			for i, v := range text {
				if v == "minutes" {
					loggedMins, err := strconv.Atoi(text[i-1])
					if err != nil {
						break
					}
					intDurationStartFromLog = loggedMins
				}
			}
		}
	}
	return file, progressPercentFromLog, intDurationStartFromLog
}
