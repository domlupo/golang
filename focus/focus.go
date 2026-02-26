package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"slices"
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
	windowMultiplier            float32
	windowWidth                 int
	windowHeight                int
	backgroundColor             color.RGBA
	timerColor                  color.RGBA
	doneColor                   color.RGBA
	logFilepath                 string
	logFile                     *os.File
	lastSecondFileLine          int
	catLogCmds                  []string
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
	logInfoPadding              int
	minuteTicker                time.Ticker
	secondTicker                time.Ticker
	progressPercent             float64
	progressPercentStart        float64
	intProgressPercent          int
	progressLoadedFromLogFile   bool
	err                         error
	mplusFaceSource             *text.GoTextFaceSource
}

func newGame() Game {
	// configs
	directoryName := "/tmp/focus_timer"
	logFileName := "log.txt"
	g := Game{
		windowMultiplier:            1.0,
		windowWidth:                 1280,
		windowHeight:                720,
		intDurationStart:            60,
		logMinsElapsedThresholdMins: 100,
		logPercentPadding:           5, // Largest percent logged (100.0) char length is 5.
		logMinutePadding:            4, // Largest minute log (1440) char length is 4.
		backgroundColor:             color.RGBA{50, 50, 50, 255},
		timerColor:                  color.RGBA{75, 75, 75, 255},
		doneColor:                   color.RGBA{0, 150, 100, 255},
		catLogCmds:                  []string{"l", "lo", "log"},
	}

	// set game fields which are calculated from configs
	g.logInfoPadding = max(g.logPercentPadding, g.logMinutePadding)
	g.logFilepath = fmt.Sprintf("%s/%s", directoryName, logFileName)
	dirErr := os.MkdirAll(directoryName, 0755)
	if dirErr != nil {
		panic(fmt.Sprintf("paths: cannot make directory: %s", dirErr.Error()))
	}
	g.formatStartInfo(false)

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
		g.progressPercent = float64(g.currentTime.Sub(g.startTime)/time.Millisecond) /
			float64(g.duration/time.Millisecond) * maxPercentage

		if g.intProgressPercent != int(g.progressPercent) {
			g.intProgressPercent = int(g.progressPercent)
			loggedProgressPercent := int((float64(g.intProgressPercent) * (float64(maxPercentage) -
				g.progressPercentStart) / maxPercentage) + g.progressPercentStart)
			_, g.err = fmt.Fprintf(g.logFile,
				"%*d%% %s %*d minutes remaining\n",
				g.logPercentPadding,
				loggedProgressPercent,
				g.currentTime.Format(time.RFC3339),
				g.logMinutePadding,
				g.intDurationCount)
			if g.err != nil {
				log.Fatal(g.err)
			}
		}

		go func() {
			for {
				select {
				case <-g.secondTicker.C:
					g.overwriteLastSecond()
				case <-g.minuteTicker.C:
					g.intDurationCount = g.intDurationCount - 1
					loggedProgressPercent := (g.progressPercent * (float64(maxPercentage) -
						g.progressPercentStart) / maxPercentage) + g.progressPercentStart
					if g.intDurationStart > g.logMinsElapsedThresholdMins {
						_, g.err = fmt.Fprintf(g.logFile,
							"%*.1f%% %s %*d minutes remaining\n",
							g.logPercentPadding,
							loggedProgressPercent,
							g.currentTime.Format(time.RFC3339),
							g.logMinutePadding,
							g.intDurationCount)
					}
				default:
					time.Sleep(1 * time.Microsecond)
				}
			}
		}()
	}

	if g.intProgressPercent == maxPercentage {
		g.logFile.Close()
	}

	return nil
}

func (g *Game) overwriteLastSecond() {
	input, err := os.ReadFile(g.logFilepath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")
	lines[g.lastSecondFileLine] = lastSecondFormat(g.currentTime, g.logInfoPadding)
	output := strings.Join(lines, "\n")

	err = os.WriteFile(g.logFilepath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func lastSecondFormat(t time.Time, padding int) string {
	return pad(fmt.Sprintf("Last second %s", t.Format(time.RFC3339)), padding)
}

func pad(info string, padding int) string {
	p := strings.Repeat(" ", padding)
	i := p + info + p
	return i
}

func (g *Game) formatStartInfo(writeToFile bool) {
	timerSet := pad(fmt.Sprintf("Timer set for %d minutes.", g.intDurationStart), g.logInfoPadding)
	timerStarted := pad(fmt.Sprintf(
		"Timer creation %s", g.startTime.Format(time.RFC3339)), g.logInfoPadding)
	lastSecond := lastSecondFormat(g.startTime, g.logInfoPadding)
	dashes := strings.Repeat("-", max(len(timerSet), len(timerStarted), len(lastSecond)))

	lines := []string{dashes, timerSet, timerStarted, lastSecond, dashes}
	var formattedLines string
	for i, line := range lines {
		if line == lastSecond {
			g.lastSecondFileLine = i
		}
		formattedLines = formattedLines + line + "\n"
	}

	if writeToFile {
		_, g.err = fmt.Fprint(g.logFile, formattedLines)
	}

	if g.err != nil {
		log.Fatal(g.err)
	}
}

func (g *Game) startTimer() {
	g.timerStarted = true
	g.startTime = time.Now()
	g.duration = time.Duration(g.intDurationStart) * time.Minute
	g.intDurationCount = g.intDurationStart
	g.minuteTicker = *time.NewTicker(1 * time.Minute)
	g.secondTicker = *time.NewTicker(1 * time.Second)

	if !g.progressLoadedFromLogFile || g.intDurationStart != g.intDurationStartFromLog {
		g.progressPercentStart = 0
		g.logFile, g.err = os.Create(g.logFilepath)
		if g.err != nil {
			log.Fatal(g.err)
		}
		g.formatStartInfo(true)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.timerStarted || g.currentTime.Sub(g.startTime) > g.duration {
		screen.Fill(g.doneColor)
	} else {
		screen.Fill(g.backgroundColor)
		// Use int not percent read from log for prior timer progress since timer has 1 min granularity.
		priorTimerProgress := g.progressPercentStart / maxPercentage
		currentTimerProgess := (g.progressPercent * (maxPercentage - g.progressPercentStart) /
			maxPercentage) / maxPercentage
		vector.FillRect(
			screen,
			0,
			0,
			float32(g.windowWidth)*g.windowMultiplier*float32(priorTimerProgress+currentTimerProgess),
			float32(g.windowHeight)*g.windowMultiplier,
			g.timerColor,
			false)
	}

	if !g.timerStarted {
		op := &text.DrawOptions{}
		s := "%d minutes. Change minutes with mouse wheel. Press any mouse button to start timer"
		text.Draw(screen, fmt.Sprintf(s, g.intDurationStart), &text.GoTextFace{
			Source: g.mplusFaceSource,
			Size:   float64(24 * g.windowMultiplier),
		}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(g.windowWidth) * g.windowMultiplier),
		int(float32(g.windowHeight) * g.windowMultiplier)
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
			g.windowMultiplier = 1.5
		}
	}

	ebiten.SetWindowSize(
		int(float32(g.windowWidth)*g.windowMultiplier),
		int(float32(g.windowHeight)*g.windowMultiplier))

	if readCatLogCmdFromCLI(g.logFilepath, g.catLogCmds) {
		return
	}

	durationStartFromCLI := readMinsFromCLI()
	file, progressPercentFromLog, durationStartFromLog := readProgressFromLog(g.logFilepath)

	if durationStartFromCLI > 0 {
		g.intDurationStart = durationStartFromCLI
	} else if durationStartFromLog > 0 && progressPercentFromLog < maxPercentage {
		g.logFile = file
		g.progressLoadedFromLogFile = true
		g.progressPercentStart = progressPercentFromLog
		g.intDurationStartFromLog = durationStartFromLog
		g.intDurationStart = durationStartFromLog
	}

	if err := ebiten.RunGameWithOptions(&g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal(err)
	}

	g.logFile.Close()
}

func readCatLogCmdFromCLI(logPath string, catLogCmds []string) bool {
	args := os.Args
	if len(args) == 2 && (slices.Contains(catLogCmds, strings.ToLower(args[1]))) {
		cmd := exec.Command("cat", logPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return true
	}
	return false
}

func readMinsFromCLI() int {
	args := os.Args
	switch {
	case len(args) > 2:
		panic("Invalid number of arguments passed, only one is expected.")
	case len(args) == 2:
		mins, convertMinsErr := strconv.Atoi(os.Args[1])
		if convertMinsErr != nil {
			panic("Invalid type passed, an integer is expected.")
		}
		if mins < 1 || mins > minutesInADay {
			panic(fmt.Sprintf(
				"Invalid number passed, expected a number between 1 and %d. (maximum minutes in a day)",
				minutesInADay))
		}
		return mins
	default:
		return 0
	}
}

func readProgressFromLog(filename string) (*os.File, float64, int) {
	file, fileErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	var progressPercentFromLog float64
	var durationStartFromLog int

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
					durationStartFromLog = loggedMins
				}
			}
		}
	}
	return file, progressPercentFromLog, durationStartFromLog
}
