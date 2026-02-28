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
	progressReadFromLogFile     bool
	minuteDurationToggle        int
	secondDurationToggle        int
	timerDurationTotal          time.Duration
	timerDurationLeft           time.Duration
	currentTime                 time.Time
	startTime                   time.Time
	logMinsElapsedThresholdMins int
	logPercentPadding           int
	logMinutePadding            int
	logSecondPadding            int
	logInfoPadding              int
	minuteTicker                time.Ticker
	secondTicker                time.Ticker
	progressPercent             float64
	intProgressPercent          int
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
		minuteDurationToggle:        60,
		logMinsElapsedThresholdMins: 100,
		logPercentPadding:           5, // Largest percent logged (100.0) char length is 5.
		logMinutePadding:            4, // Largest minute logged (1440) char length is 4.
		logSecondPadding:            2, // Largest second logged (60) char length is 2.
		backgroundColor:             color.RGBA{50, 50, 50, 255},
		timerColor:                  color.RGBA{75, 75, 75, 255},
		doneColor:                   color.RGBA{0, 150, 100, 255},
		catLogCmds:                  []string{"l", "lo", "log"},
		timeLayout:                  time.RFC3339,
		lockReleaseDelay:            time.Second * 2,
	}

	if g.lockReleaseDelay < time.Second*2 {
		panic("newGame: cannot set lock release delay less than 2 seconds since writes occur every ~1" +
			"second")
	}

	// set game fields which are calculated from configs
	g.logInfoPadding = max(g.logPercentPadding, g.logMinutePadding) * 2
	g.logFilepath = fmt.Sprintf("%s/%s", directoryName, logFileName)
	dirErr := os.MkdirAll(directoryName, 0755)
	if dirErr != nil {
		panic(fmt.Sprintf("newGame: cannot make directory: %s", dirErr.Error()))
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
		g.minuteDurationToggle = max(g.minuteDurationToggle+int(dy), 0)
		g.minuteDurationToggle = min(g.minuteDurationToggle, minutesInADay)

		if g.isTimerSetToUseLogFileProgress() {
			g.secondDurationToggle = int(g.timerDurationLeft.Seconds()) - g.minuteDurationToggle*60
		} else {
			g.secondDurationToggle = 0
		}

		if (g.minuteDurationToggle > 0 || g.secondDurationToggle > 0) &&
			ebiten.IsMouseButtonPressed(ebiten.MouseButton0) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton1) ||
			ebiten.IsMouseButtonPressed(ebiten.MouseButton2) {
			g.startTimer()
		}
	}

	if g.timerStarted && g.progressPercent != maxPercentage {
		priorSessionTime := float64(g.timerDurationTotal - g.timerDurationLeft)
		currentSessionTime := float64(time.Since(g.startTime))
		g.progressPercent = (priorSessionTime + currentSessionTime) / float64(g.timerDurationTotal) *
			maxPercentage
		g.progressPercent = min(g.progressPercent, maxPercentage)

		if g.intProgressPercent != int(g.progressPercent) {
			g.intProgressPercent = int(g.progressPercent)
			durationLeft := g.timerDurationLeft - (time.Since(g.startTime))
			minsLeft := int(durationLeft.Minutes())
			secsLeft := int(durationLeft.Seconds()) - minsLeft*60

			_, g.err = fmt.Fprintf(g.logFile,
				"%*d%% %s %*d minutes %0*d seconds left\n",
				g.logPercentPadding,
				g.intProgressPercent,
				g.currentTime.Format(g.timeLayout),
				g.logMinutePadding,
				minsLeft,
				g.logSecondPadding,
				secsLeft)

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
					durationLeft := g.timerDurationLeft - (time.Since(g.startTime))
					minsLeft := int(durationLeft.Minutes())
					secsLeft := int(durationLeft.Seconds()) - minsLeft*60

					if g.minuteDurationToggle > g.logMinsElapsedThresholdMins {
						_, g.err = fmt.Fprintf(g.logFile,
							"%*.1f%% %s %*d minutes %0*d seconds left\n",
							g.logPercentPadding,
							g.progressPercent,
							g.currentTime.Format(g.timeLayout),
							g.logMinutePadding,
							minsLeft,
							g.logSecondPadding,
							secsLeft)
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
		panic("overwriteLastSecond: invalid log file. detected inconsistent seconds in log file.")
	}

	input, err := os.ReadFile(g.logFilepath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(input), "\n")
	// TODO: verify this is safe. This is related to other TODO for using time across two goroutines.
	g.lastWrittenSecond = g.currentTime
	lines[g.lastSecondFileLine] = lastSecondFormat(
		g.lastWrittenSecond, g.logInfoPadding, g.timeLayout)
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

func lastSecondFormat(t time.Time, padding int, layout string) string {
	return pad(fmt.Sprintf("Last second:   %s", t.Format(layout)), padding)
}

func pad(info string, padding int) string {
	p := strings.Repeat(" ", padding)
	i := p + info + p
	return i
}

func (g *Game) startTimer() {
	g.timerStarted = true
	g.startTime = time.Now()
	g.minuteTicker = *time.NewTicker(1 * time.Minute)
	g.secondTicker = *time.NewTicker(1 * time.Second)

	if !g.isTimerSetToUseLogFileProgress() {
		g.timerDurationTotal = time.Duration(g.minuteDurationToggle) * time.Minute
		g.timerDurationLeft = g.timerDurationTotal
		g.logFile, g.err = os.Create(g.logFilepath)
		if g.err != nil {
			log.Fatal(g.err)
		}
	}

	g.checkFileLock(true)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.timerStarted || g.progressPercent == 100 {
		screen.Fill(g.doneColor)
	} else {
		screen.Fill(g.backgroundColor)
		vector.FillRect(
			screen,
			0,
			0,
			float32(g.windowWidth)*g.windowMultiplier*float32(g.progressPercent)/maxPercentage,
			float32(g.windowHeight)*g.windowMultiplier,
			g.timerColor,
			false)
	}

	if !g.timerStarted {
		op := &text.DrawOptions{}
		suffix := " Mouse wheel changes minutes. Press any mouse button to start timer."
		if g.secondDurationToggle == 0 {
			s := "Timer set to %d minutes." + suffix
			text.Draw(screen, fmt.Sprintf(s, g.minuteDurationToggle),
				&text.GoTextFace{
					Source: g.mplusFaceSource,
					Size:   float64(24 * g.windowMultiplier),
				}, op)
		} else {
			s := "Timer set to %d minutes %d seconds." + suffix
			text.Draw(screen, fmt.Sprintf(s, g.minuteDurationToggle, g.secondDurationToggle),
				&text.GoTextFace{
					Source: g.mplusFaceSource,
					Size:   float64(24 * g.windowMultiplier),
				}, op)
		}
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
	file, durationLeft, durationTotal := readTimerProgressFromLog(g.logFilepath, g.timeLayout)

	if durationStartFromCLI > 0 {
		g.minuteDurationToggle = durationStartFromCLI
	} else if durationLeft.Seconds() > 0 {
		g.progressReadFromLogFile = true
		g.logFile = file
		g.timerDurationLeft = durationLeft
		g.timerDurationTotal = durationTotal
		g.minuteDurationToggle = int(durationLeft.Minutes())
		g.secondDurationToggle = int(durationLeft.Seconds()) - g.minuteDurationToggle*60
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
		panic("readMinsFromCLI: Invalid number of arguments passed, only one is expected.")
	case len(args) == 2:
		mins, convertMinsErr := strconv.Atoi(os.Args[1])
		if convertMinsErr != nil {
			panic("readMinsFromCLI: Invalid type passed, an integer is expected.")
		}
		if mins < 1 || mins > minutesInADay {
			panic(fmt.Sprintf("readMinsFromCLI: Invalid number passed, expected a number between 1 and "+
				"%d. (maximum minutes in a day)", minutesInADay))
		}
		return mins
	default:
		return 0
	}
}

// Attempt to read timer progress from log. Returns the file, duration left, and total duration.
func readTimerProgressFromLog(filename string, layout string) (*os.File, time.Duration, time.Duration) {
	file, fileErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	var timerDurationTotal time.Duration
	var timerStarts []time.Time
	var timerStops []time.Time

	if fileErr == nil {
		scanner := bufio.NewReader(file)
		for {
			line, readLineErr := scanner.ReadString('\n')
			if readLineErr != nil {
				break
			}
			line = line[:len(line)-1]
			lowercaseLine := strings.ToLower(line)

			if strings.Contains(lowercaseLine, "set") {
				for str := range strings.SplitSeq((line), " ") {
					mins, err := strconv.Atoi(str)
					if err == nil {
						timerDurationTotal = time.Minute * time.Duration(mins)
					}
				}
			}

			// TODO: add keywords to config rather than hard code
			var timerStart, timerStop bool
			if strings.Contains(lowercaseLine, "created") ||
				strings.Contains(lowercaseLine, "unpaused") {
				timerStart = true
			} else if strings.Contains(lowercaseLine, "last") ||
				strings.Contains(lowercaseLine, "paused") {
				timerStop = true
			}

			for str := range strings.SplitSeq(line, " ") {
				t, err := time.Parse(layout, str)
				if err == nil && timerStart {
					timerStarts = append(timerStarts, t)
				} else if err == nil && timerStop {
					timerStops = append(timerStops, t)
				}
			}
		}
	}

	if len(timerStarts) != len(timerStops) {
		panic(fmt.Sprintf("readDurationLeftFromLog: cannot determine duration since start and stops "+
			" are not equal. starts: %d stops: %d", len(timerStarts), len(timerStops)))
	}

	if len(timerStops) > 0 {
		lastTimerStop := timerStops[0]
		timerStops = append(timerStops[1:], lastTimerStop)
	}

	var timerDurationComplete time.Duration
	for i := 0; i < len(timerStarts); i++ {
		timerDurationComplete += timerStops[i].Sub(timerStarts[i])
	}

	return file, timerDurationTotal - timerDurationComplete, timerDurationTotal
}

// Prevents two focus timers from ticking at once. The last second logged is used as a lock.
func (g *Game) checkFileLock(writeInitialFileInfo bool) {
	// format initial file info. This is needed to determine which line number the lock is located on.
	minStr := "minutes"
	if g.minuteDurationToggle == 1 {
		minStr = "minute"
	}
	timerSet := pad(fmt.Sprintf("Timer set for %d %s.", g.minuteDurationToggle, minStr), g.logInfoPadding)
	timerStarted := pad(fmt.Sprintf(
		"Timer created: %s", g.startTime.Format(g.timeLayout)), g.logInfoPadding)
	lastSecond := lastSecondFormat(g.startTime, g.logInfoPadding, g.timeLayout)
	dashes := strings.Repeat("-", max(len(timerSet), len(timerStarted), len(lastSecond)))

	// set last second file line aka the lock line
	// first line is empty to prevent cat format issues when running in background via &
	lines := []string{"", dashes, timerSet, timerStarted, lastSecond, dashes}
	for i, line := range lines {
		if line == lastSecond {
			g.lastSecondFileLine = i
		}
	}

	// check lock
	lock, ok := readTime(g.logFilepath, g.lastSecondFileLine, g.timeLayout)
	if ok && lock.Add(g.lockReleaseDelay).After(time.Now()) {
		panic(fmt.Sprintf("checkFileLock: cannot start timer. A timer is currently running or was "+
			"stopped within the last %d seconds", int(g.lockReleaseDelay.Seconds())))
	}

	var formattedLines string
	if writeInitialFileInfo && !g.isTimerSetToUseLogFileProgress() {
		for _, line := range lines {
			formattedLines = formattedLines + line + "\n"
		}
		_, g.err = fmt.Fprint(g.logFile, formattedLines)
	} else if writeInitialFileInfo && g.isTimerSetToUseLogFileProgress() {
		timerPaused := pad(fmt.Sprintf("Timer paused:   %s", lock.Format(g.timeLayout)), g.logInfoPadding)
		timerUnpaused := pad(fmt.Sprintf("Timer unpaused: %s", g.currentTime.Format(g.timeLayout)), g.logInfoPadding)
		lines = []string{dashes, timerPaused, timerUnpaused, dashes}
		for _, line := range lines {
			formattedLines = formattedLines + line + "\n"
		}
		_, g.err = fmt.Fprint(g.logFile, formattedLines)
		g.overwriteLastSecond()
	}

	if g.err != nil {
		log.Fatal(g.err)
	}
}

func (g *Game) isTimerSetToUseLogFileProgress() bool {
	return g.progressReadFromLogFile &&
		g.minuteDurationToggle == int(g.timerDurationLeft.Minutes()) &&
		int(g.timerDurationLeft.Seconds()) > 0
}
