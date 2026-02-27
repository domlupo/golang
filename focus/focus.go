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
	timeLayout                  string
	lockReleaseDelay            time.Duration
	lastWrittenSecond           time.Time
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
		timeLayout:                  time.RFC3339,
		lockReleaseDelay:            time.Second * 2,
	}

	if g.lockReleaseDelay < time.Second*2 {
		panic("cannot set lock release delay to 1 second or less since write frequency is ~1 second")
	}

	// set game fields which are calculated from configs
	g.logInfoPadding = max(g.logPercentPadding, g.logMinutePadding)
	g.logFilepath = fmt.Sprintf("%s/%s", directoryName, logFileName)
	dirErr := os.MkdirAll(directoryName, 0755)
	if dirErr != nil {
		panic(fmt.Sprintf("paths: cannot make directory: %s", dirErr.Error()))
	}

	return g
}

func (g *Game) Update() error {
	// TODO: verify this is safe to set here. I am worried because I use this same value across two
	// goroutines. I think this is likely to cause race conditions and the focus timer should be
	// restructured.
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
	lastWrittenSecond, ok := readTime(g.logFilepath, g.lastSecondFileLine, g.timeLayout)
	if !ok {
		panic("overwriteLastSecond: cannot read last second from log file")
	}
	epoch := time.Time{}
	if !g.lastWrittenSecond.Equal(epoch) &&
		lastWrittenSecond.Format(g.timeLayout) != g.lastWrittenSecond.Format(g.timeLayout) {
		// invalid log file. Most likely more than one focus timer is running or the log file was
		// manually updated. I believe panic is the best user experience rather than letting any timers
		// continue since the log file may be incorrect. Or a very rare bug happened since there are
		// several protections to avoid multiple timers running. Losing timer progress and failing
		// loudly for rare bugs is better than suppressing rare bugs.

		// TODO: decide if I should add file invalidation
		panic("overwriteLastSecond: invalid log file. detected inconsistent seconds in log file.")
	}

	input, err := os.ReadFile(g.logFilepath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(input), "\n")
	// TODO: verify this is safe. This is related to other TODO for using time across two goroutines.
	g.lastWrittenSecond = g.currentTime
	lines[g.lastSecondFileLine] = lastSecondFormat(g.lastWrittenSecond, g.logInfoPadding)
	output := strings.Join(lines, "\n")

	err = os.WriteFile(g.logFilepath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func readTime(filepath string, lineNumber int, layout string) (time.Time, bool) {
	input, readFileErr := os.ReadFile(filepath)
	if readFileErr != nil {
		// one reason is when file does not exist
	}

	if len(input) == 0 {
		return time.Time{}, false
	}

	lines := strings.Split(string(input), "\n")
	for str := range strings.SplitSeq(lines[lineNumber], " ") {
		t, err := time.Parse(layout, str)
		if err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

func lastSecondFormat(t time.Time, padding int) string {
	return pad(fmt.Sprintf("Last second: %s", t.Format(time.RFC3339)), padding)
}

func pad(info string, padding int) string {
	p := strings.Repeat(" ", padding)
	i := p + info + p
	return i
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
	}

	g.checkFileLock(true)
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
	if readCatLogCmdFromCLI(g.logFilepath, g.catLogCmds) {
		return
	}
	g.checkFileLock(false)

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

// Prevents two focus timers from ticking at once. The last second logged is used as a lock.
func (g *Game) checkFileLock(writeInitialFileInfo bool) {
	// format initial file info. This is needed to determine which line number the lock is located on.
	minStr := "minutes"
	if g.intDurationStart == 1 {
		minStr = "minute"
	}
	timerSet := pad(fmt.Sprintf("Timer set for %d %s.", g.intDurationStart, minStr), g.logInfoPadding)
	timerStarted := pad(fmt.Sprintf(
		"Timer created: %s", g.startTime.Format(time.RFC3339)), g.logInfoPadding)
	lastSecond := lastSecondFormat(g.startTime, g.logInfoPadding)
	dashes := strings.Repeat("-", max(len(timerSet), len(timerStarted), len(lastSecond)))

	// first line is empty to prevent cat format issues when running in background via &
	lines := []string{"", dashes, timerSet, timerStarted, lastSecond, dashes}
	var formattedLines string
	for i, line := range lines {
		if line == lastSecond {
			g.lastSecondFileLine = i
		}
		formattedLines = formattedLines + line + "\n"
	}

	lock, ok := readTime(g.logFilepath, g.lastSecondFileLine, g.timeLayout)

	if ok && lock.Add(g.lockReleaseDelay).After(time.Now()) {
		panic(fmt.Sprintf(
			"Cannot start timer. A timer is currently running or was stopped within the last %d seconds",
			int(g.lockReleaseDelay.Seconds())))
	}

	// TODO: DRY this logic. It is repeated in startTimer fn
	if writeInitialFileInfo && (!g.progressLoadedFromLogFile || g.intDurationStart != g.intDurationStartFromLog) {
		_, g.err = fmt.Fprint(g.logFile, formattedLines)
	} else if writeInitialFileInfo {
		g.overwriteLastSecond()
	}

	if g.err != nil {
		log.Fatal(g.err)
	}
}
