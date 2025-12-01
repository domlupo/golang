// Main entry point for very high level game logic and configuration.

package main

import (
	"errors"
	"fmt"
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
	nodeSetup                  nodeSetup
	obstacleSetup              obstacleSetup
	slowTickInterval           time.Duration
	slowTickTime               time.Time
	priorTime                  time.Time
	currentTime                time.Time
	paused                     bool
	hotRestart                 bool
	slowTick                   bool
	cpuProfilerOn              bool
	memoryProfilerOn           bool
	state                      *state
	drawer                     *drawer
	computers                  *computers
	controller                 *controller
	debugger                   *debugger
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
		slowTickInterval:        1 * time.Second,
		slowTickTime:            time.Now(),
		priorTime:               time.Now(),
		currentTime:             time.Now(),
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
	return g
}

func (g *Game) restart(nodeSetup nodeSetup, obstacleSetup obstacleSetup, verify bool) {
	// TODO: Need to restart timers, etc
	start := time.Now()
	grid := newGrid(newGridConfig())
	hexagons := newHexagons(grid)
	obstacles, oLUT := setupObstacles(obstacleSetup, grid)
	nodes, nLUT := setupNodes(nodeSetup, grid, g.initialNodeHealth, hexagons, oLUT)
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

	g.state = newState(nodes, obstacles, hexagons, oLUT, nLUT, grid, g.developerMode)
	g.debugger = newDebugger(g.state, g.developerToggleCooldown)
	g.computers = newComputers(g.state)
	g.controller = newController(g.state, g.computers, g.debugger)
	// TODO: I should only invoke newDrawConfig once per game not every restart
	g.drawer = newDrawer(newDrawConfig(), g.state, g.controller, g.debugger)
	ebiten.SetWindowSize(int(float32(g.state.grid.width)*g.state.grid.pixelMultiplier), int(float32(g.state.grid.height)*g.state.grid.pixelMultiplier))
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

	// update game state prior to player actions
	if !g.paused && !(g.slowTick && (g.currentTime.Sub(g.slowTickTime) < g.slowTickInterval)) {
		if g.slowTick {
			logger.log(infoLevel, devTool, "slowTick")
		}
		g.state.tick()
	}
	g.controller.processInput()

	var xSum int
	var ySum int
	var totalCount int
	var infantryCount int
	var bombGravityCount int
	var bomberCount int
	var fighterCount int
	var unknownCount int
	for u := range g.state.units() {
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
		nodeCapture := g.state.updateForNodeCapture()
		g.computers.updateNodetargets(nodeCapture, avgUnitPos)
		g.slowTickTime = time.Now()
	}

	var pOneNodeCount int
	var pTwoNodeCount int

	for _, node := range g.state.nodes {
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
	g.drawer.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(g.state.grid.width) * g.state.grid.pixelMultiplier), int(float32(g.state.grid.height) * g.state.grid.pixelMultiplier)
}

func main() {
	g := newGame()

	var monitors []*ebiten.MonitorType
	monitors = ebiten.AppendMonitors(monitors)
	if g.developerMode {
		for _, m := range monitors {
			if m.Name() == "DELL U2720Q" {
				ebiten.SetMonitor(m)
			}
		}
	}

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
