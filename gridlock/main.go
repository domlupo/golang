// Main entry point for very high level game logic and configuration.

package main

import (
	"errors"
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// A relatively safe global logger. A few functions affect the state of the logger.
// e.g. level of logs emitted, type of logs emitted. These should generally be called
// only from the main game logic file.
var logger l

type Game struct {
	*drawConfig
	nodeSetup                  nodeSetup
	obstacleSetup              obstacleSetup
	computers                  computers
	backgroundColor            color.RGBA
	produceWait                time.Duration
	slowTickInterval           time.Duration
	slowTickTime               time.Time
	priorTime                  time.Time
	currentTime                time.Time
	leftMouseClicked           bool
	leftMouseHeld              bool
	priorLeftMouseClicked      bool
	leftMouseReleased          bool
	middleMouseScrollClicked   bool
	rightMouseClicked          bool
	paused                     bool
	hotRestart                 bool
	slowTick                   bool
	cpuProfilerOn              bool
	memoryProfilerOn           bool
	initialSelectedX           float32
	initialSelectedY           float32
	currentSelectedX           float32
	currentSelectedY           float32
	gameState                  *state
	debugger                   debugger
	developerMode              bool
	developerToggles           map[string]*toggle
	developerToggleCooldown    time.Duration
	initialNodeHealth          int
	outputDirectoryName        string
	cpuProfilerToggleName      string
	memProfilerToggleName      string
	cpuProfilerFileName        string
	memProfilerFileName        string
	cpuProfilerStartBeforeGame bool
	memProfilerStartBeforeGame bool
	cacheSetups                bool
	cpuProfilerFile            *os.File
	memoryProfilerFile         *os.File
}

func newGame() Game {
	g := Game{
		developerMode:           true,
		obstacleSetup:           multiObstacle,
		nodeSetup:               leftVsRightVerticalOffsetAircraftAndOrdnance,
		backgroundColor:         color.RGBA{200, 200, 200, 255},
		produceWait:             1_000_000 * time.Second,
		slowTickInterval:        1 * time.Second,
		slowTickTime:            time.Now(),
		priorTime:               time.Now(),
		currentTime:             time.Now(),
		initialSelectedX:        float32(-1.0),
		initialSelectedY:        float32(-1.0),
		currentSelectedX:        float32(-1.0),
		currentSelectedY:        float32(-1.0),
		hotRestart:              true,
		cpuProfilerToggleName:   "cpuProfiler",
		memProfilerToggleName:   "memProfiler",
		developerToggles:        make(map[string]*toggle),
		developerToggleCooldown: 300 * time.Millisecond,
		initialNodeHealth:       100,
		outputDirectoryName:     "output",
		cpuProfilerFileName:     "cpu.prof",
		memProfilerFileName:     "mem.prof",
	}
	logger = setupLogger(g.developerToggleCooldown)
	g.drawConfig = newDrawConfig()
	return g
}

func (g *Game) restart(nodeSetup nodeSetup, obstacleSetup obstacleSetup, verify bool) {
	// TODO: Need to restart timers, etc

	// TODO: Determine if grid being passed here is correct. It might be best to have both lower level grid
	// for positions and a higher level struct for the grid components
	start := time.Now()
	hexagons := setupHexagons(g.drawConfig.grid)
	obstacles, oLUT := setupObstacles(obstacleSetup, g.drawConfig.grid)
	nodes, nLUT := setupNodes(nodeSetup, g.drawConfig.grid, g.initialNodeHealth, hexagons, oLUT)
	logger.log(infoLevel, setupLogic, fmt.Sprintf("Setup required %0.2f seconds", float64(time.Since(start).Seconds())))

	if verify {
		for nPos := range nodes {
			for oPos := range obstacles {
				if nPos.equals(oPos) {
					panic("grid position conflict between node and obstacle")
				}
			}
		}
	}

	g.gameState = newState(nodes, obstacles, hexagons, oLUT, nLUT, g.drawConfig.grid, g.developerMode)
	g.debugger = newDebugger(g.gameState, g.developerToggleCooldown)
	g.computers = newComputers(g.gameState)
}

func (g *Game) Update() error {
	g.priorTime = g.currentTime
	g.currentTime = time.Now()
	logger.log(
		infoLevel,
		mainLogic,
		"game performance",
		slog.String("fps", fmt.Sprintf("%0.2f", ebiten.ActualFPS())),
		slog.String("tps", fmt.Sprintf("%0.2f", ebiten.ActualTPS())),
	)

	for _, t := range g.developerToggles {
		t.checkToggle(g.currentTime)
	}
	for _, t := range g.debugger.toggles {
		t.checkToggle(g.currentTime)
	}
	for _, t := range logger.toggles {
		t.checkToggle(g.currentTime)
	}

	// refresh
	g.leftMouseReleased = false
	g.leftMouseHeld = false

	// update game state prior to player actions
	if !g.paused && !(g.slowTick && (g.currentTime.Sub(g.slowTickTime) < g.slowTickInterval)) {
		if g.slowTick {
			logger.log(infoLevel, devTool, "slowTick")
		}
		g.gameState.tick()
	}

	// player actions
	graphicalX, graphicalY := ebiten.CursorPosition()
	_, mouseWheelDelta := ebiten.Wheel()
	g.debugger.updateHeatmapPercent(-int(mouseWheelDelta))

	x := float32(graphicalX) / g.windowMultiplier
	y := float32(graphicalY) / g.windowMultiplier
	pixel := Pixel{int(x), int(y)}

	g.priorLeftMouseClicked = g.leftMouseClicked
	g.leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	g.middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	g.rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	g.currentSelectedX = x
	g.currentSelectedY = y
	g.debugger.pixel = pixel
	if g.priorLeftMouseClicked && !g.leftMouseClicked {
		g.leftMouseReleased = true
	} else if !g.priorLeftMouseClicked && g.leftMouseClicked {
		g.initialSelectedX = x
		g.initialSelectedY = y
	} else if g.leftMouseClicked {
		g.leftMouseHeld = true
	}
	if g.developerMode && g.middleMouseScrollClicked {
		g.debugger.debug(pixel)
	}
	if g.rightMouseClicked {
		filter := []player{pNeutral}
		filter = append(filter, g.computers.active()...)
		g.gameState.updateTargets(position{lLand, pixel}, filter)
	}

	// update game state after player actions
	if g.leftMouseReleased {
		g.gameState.updateSelected(
			Pixel{int(g.initialSelectedX), int(g.initialSelectedY)},
			Pixel{int(g.currentSelectedX), int(g.currentSelectedY)})
	}

	var xSum int
	var ySum int
	var totalCount int
	var infantryCount int
	var bombGravityCount int
	var bomberCount int
	var fighterCount int
	var unknownCount int
	for u := range g.gameState.units() {
		xSum += u.X
		ySum += u.Y
		totalCount++
		switch u.name {
		case infantry:
			infantryCount++
		case bombGravity:
			bombGravityCount++
		case bomber:
			bomberCount++
		case fighter:
			fighterCount++
		default:
			unknownCount++
		}
	}
	var avgUnitPos Pixel
	if totalCount != 0 {
		avgUnitPos = Pixel{X: xSum / totalCount, Y: ySum / totalCount}
	}
	logger.log(infoLevel, mainLogic, "unit counts",
		slog.Int("total_units", totalCount),
		slog.Int("infantryCount", infantryCount),
		slog.Int("bombGravity", bombGravityCount),
		slog.Int("bomber", bomberCount),
		slog.Int("fighter", fighterCount),
		slog.Int("unknownCount", unknownCount),
	)

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if !g.paused && !(g.slowTick && (g.currentTime.Sub(g.slowTickTime) < g.slowTickInterval)) {
		nodeCapture := g.gameState.updateForNodeCapture()
		g.computers.updateNodetargets(nodeCapture, avgUnitPos)
		g.slowTickTime = time.Now()
	}

	var pOneNodeCount int
	var pTwoNodeCount int

	for _, node := range g.gameState.nodes {
		switch node.player {
		case pOne:
			pOneNodeCount++
		case pTwo:
			pTwoNodeCount++
		}
	}

	if pOneNodeCount == 0 || pTwoNodeCount == 0 {
		if g.hotRestart {
			logger.log(infoLevel, mainLogic, "restartGame")

			g.restart(g.nodeSetup, g.obstacleSetup, true)
		} else {
			return errors.New("regular termination")
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawConfig.screen == nil {
		g.drawConfig.screen = screen
	}

	screen.Fill(g.drawConfig.backgroundColor)

	g.debugger.draw(*g.drawConfig)
	g.gameState.draw(*g.drawConfig)

	if g.leftMouseHeld {
		drawUnfilledBoxFloat32(
			*g.drawConfig,
			g.drawConfig.highlightColor,
			g.initialSelectedX,
			g.initialSelectedY,
			g.currentSelectedX,
			g.currentSelectedY,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	wc := newWindowConfig()
	return int(float32(wc.width()) * g.windowMultiplier), int(float32(wc.height()) * g.windowMultiplier)
}

func main() {
	g := newGame()

	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	// TODO: Remove before release in order to avoid bugs for customers which own this monitor
	for _, m := range monitors {
		if m.Name() == "DELL U2720Q" {
			ebiten.SetMonitor(m)
		}
	}

	wc := newWindowConfig()
	ebiten.SetWindowSize(int(float32(wc.width())*g.windowMultiplier), int(float32(wc.height())*g.windowMultiplier))
	os.Mkdir(g.outputDirectoryName, 0755)

	g.developerToggles[g.cpuProfilerToggleName] = newToggle(
		ebiten.KeyC,
		func() {
			if g.cpuProfilerOn {
				logger.log(infoLevel, devTool, "Stopping CPU profiler.")
				pprof.StopCPUProfile()
				g.cpuProfilerFile.Close()
			} else {
				logger.log(infoLevel, devTool, "Starting CPU profiler.")
				f, err := os.Create(fmt.Sprintf("%s/%s", g.outputDirectoryName, g.cpuProfilerFileName))
				if err != nil {
					logger.log(errorLevel, devTool, "Could not create CPU profiler file.")
				}
				g.cpuProfilerFile = f
				if err := pprof.StartCPUProfile(f); err != nil {
					logger.log(errorLevel, devTool, "Could not start CPU profiler.")
				}
			}
			g.cpuProfilerOn = !g.cpuProfilerOn
		},
		g.developerToggleCooldown)
	g.developerToggles[g.memProfilerToggleName] = newToggle(
		ebiten.KeyX,
		func() {
			if g.memoryProfilerOn {
				logger.log(infoLevel, devTool, "Stopping memory profiler.")

				if err := pprof.WriteHeapProfile(g.memoryProfilerFile); err != nil {
					logger.log(errorLevel, devTool, "Could not write heap profile.")
				}
				g.memoryProfilerFile.Close()
			} else {
				logger.log(infoLevel, devTool, "Starting memory profiler.")
				f, err := os.Create(fmt.Sprintf("%s/%s", g.outputDirectoryName, g.memProfilerFileName))
				if err != nil {
					logger.log(errorLevel, devTool, "Could not create memory profiler file.")
				}
				g.memoryProfilerFile = f
				runtime.GC()
			}
			g.memoryProfilerOn = !g.memoryProfilerOn
		},
		g.developerToggleCooldown)
	g.developerToggles["pause"] = newToggle(
		ebiten.KeyP,
		func() { g.paused = !g.paused },
		g.developerToggleCooldown)
	g.developerToggles["slowTick"] = newToggle(
		ebiten.KeyT,
		func() { g.slowTick = !g.slowTick },
		g.developerToggleCooldown)

	if !g.developerMode {
		g.developerToggles = make(map[string]*toggle)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(g Game) {
		for range c {
			logger.log(warnLevel, mainLogic, "received interrupt signal")
			g.forceCleanup()
		}
	}(g)

	if g.cpuProfilerStartBeforeGame {
		g.developerToggles[g.cpuProfilerToggleName].callBack()
	}
	if g.memProfilerStartBeforeGame {
		g.developerToggles[g.memProfilerToggleName].callBack()
	}

	if g.cacheSetups {
		logger.log(infoLevel, setupLogic, "cacheSetups: cache all missing setups")
		start := time.Now()
		n := nodeSetups()
		o := obstacleSetups()
		m := max(len(n), len(o))
		for i := range m {
			oI := i
			nI := i
			if i > len(obstacleSetups())-1 {
				oI = 0
			} else {
				logger.log(infoLevel, setupLogic, fmt.Sprintf("cacheSetups: cache obstacle %s", o[oI].String()))
			}
			if i > len(nodeSetups())-1 {
				nI = 0
			} else {
				logger.log(infoLevel, setupLogic, fmt.Sprintf("cacheSetups: cache node %s", n[nI].String()))
			}
			g.restart(n[nI], o[oI], false)
		}
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("cacheSetups: caching all missing setups required %d seconds", int(time.Since(start).Seconds())))
		return
	}

	g.restart(g.nodeSetup, g.obstacleSetup, true)

	if err := ebiten.RunGame(&g); err != nil {
		logger.log(errorLevel, mainLogic, err.Error())
	}

	g.forceCleanup()
}

func (g Game) forceCleanup() {
	pprof.StopCPUProfile()
	g.cpuProfilerFile.Close()
	logger.log(debugLevel, devTool, "force closed cpu profiler file")
	if err := pprof.WriteHeapProfile(g.memoryProfilerFile); err != nil {
		logger.log(errorLevel, devTool, "Could not write heap profile.")
	}
	g.memoryProfilerFile.Close()
	logger.log(debugLevel, devTool, "force closed memory profiler file")
}
