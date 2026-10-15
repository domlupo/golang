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
	developerMode              bool
	initialNodeHealth          int
	outputDirectoryName        string
	cpuProfilerFileName        string
	memProfilerFileName        string
	cpuProfilerStartBeforeGame bool
	memProfilerStartBeforeGame bool
	cacheSetups                bool
	cpuProfilerFile            *os.File
	memoryProfilerFile         *os.File
}

func newGame() Game {
	developerMode := true
	g := Game{
		developerMode:       developerMode,
		obstacleSetup:       multiObstacle,
		nodeSetup:           leftVsRightVerticalOffsetFighters,
		slowTickInterval:    1 * time.Second,
		slowTickTime:        time.Now(),
		priorTime:           time.Now(),
		currentTime:         time.Now(),
		hotRestart:          true,
		initialNodeHealth:   100,
		outputDirectoryName: "output",
		cpuProfilerFileName: "cpu.prof",
		memProfilerFileName: "mem.prof",
	}
	logger = setupLogger()
	return g
}

func (g *Game) restart(nodeSetup nodeSetup, obstacleSetup obstacleSetup, verify bool) {
	// TODO: Need to restart timers, etc
	start := time.Now()

	pixelMultiplier := float32(1.25)
	if _, ok := g.isDeveloperMacStudio(); ok {
		pixelMultiplier = 2.5
	}

	logger.log(infoLevel, setupLogic, "Game.restart: newGrid")
	grid := newGrid(newGridConfig(pixelMultiplier))
	logger.log(infoLevel, setupLogic, "Game.restart: newHexagons")
	hexagons := newHexagons(grid)
	logger.log(infoLevel, setupLogic, "Game.restart: setupObstacles")
	obstacles, oLUT := setupObstacles(obstacleSetup, grid)
	logger.log(infoLevel, setupLogic, "Game.restart: newUnitConfigProvider")
	unitConfigProvider := newUnitConfigProvider()
	logger.log(infoLevel, setupLogic, "Game.restart: setupNodes")
	nodes, nLUT := setupNodes(
		nodeSetup, grid, g.initialNodeHealth, hexagons, oLUT, unitConfigProvider)

	logger.log(
		infoLevel,
		setupLogic,
		fmt.Sprintf(
			"Game.restart: Setup required %0.2f seconds",
			float64(time.Since(start).Seconds())))

	if verify {
		for nPos := range nodes {
			for oPos := range obstacles {
				if nPos.equals(oPos) {
					panic("Game.restart: tile conflict between node and obstacle")
				}
			}
		}
	}

	logger.log(infoLevel, setupLogic, "Game.restart: newPartitioner")
	partitioner := newPartitioner(grid, hexagons, unitConfigProvider)
	logger.log(infoLevel, setupLogic, "Game.restart: newState")
	g.state = newState(nodes, obstacles, hexagons, oLUT, nLUT, grid, partitioner, g.developerMode)
	logger.log(infoLevel, setupLogic, "Game.restart: newDebugger")
	debugger := newDebugger(g.state)
	logger.log(infoLevel, setupLogic, "Game.restart: newToggleRegister")
	tr := newToggleRegister(g, debugger)
	logger.log(infoLevel, setupLogic, "Game.restart: newComputers")
	g.computers = newComputers(g.state, unitConfigProvider)
	logger.log(infoLevel, setupLogic, "Game.restart: newControllers")
	g.controller = newController(g.state, g.computers, debugger, tr)

	// TODO: I should only invoke newDrawConfig once per game not every restart
	logger.log(infoLevel, setupLogic, "Game.restart: newDrawer")
	g.drawer = newDrawer(newDrawConfig(), g.state, g.controller, debugger)
	ebiten.SetWindowSize(
		int(float32(g.state.grid.pixelWidth)*g.state.grid.pixelMultiplier),
		int(float32(g.state.grid.pixelHeight)*g.state.grid.pixelMultiplier))
}

func (g *Game) Update() error {
	g.priorTime = g.currentTime
	g.currentTime = time.Now()
	logger.log(
		infoLevel,
		mainLogic,
		"Game,Update: game performance",
		slog.String("FPS", fmt.Sprintf("%0.2f", ebiten.ActualFPS())),
		slog.String("TPS", fmt.Sprintf("%0.2f", ebiten.ActualTPS())),
	)

	// update game state prior to player actions
	if !g.paused && !(g.slowTick && (g.currentTime.Sub(g.slowTickTime) < g.slowTickInterval)) {
		if g.slowTick {
			logger.log(infoLevel, devTool, "Game.Update: slowTick")
		}
		g.state.tick()
	}
	g.controller.processInput()

	var xSum float64
	var ySum float64
	var totalCount float64
	unitCounts := make(map[name]int)
	for u := range g.state.units() {
		xSum += u.X
		ySum += u.Y
		totalCount++
		unitCounts[u.name]++
	}
	var avgUnitPos Pos
	if totalCount != 0 {
		avgUnitPos = Pos{X: xSum / totalCount, Y: ySum / totalCount}
	}
	counts := []slog.Attr{slog.Float64("total", totalCount)}
	for name, count := range unitCounts {
		counts = append(counts, slog.Int(name.string(), count))
	}
	logger.log(infoLevel, mainLogic, "Game.Update: unitCounts", counts...)

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
	return int(float32(g.state.grid.pixelWidth) * g.state.grid.pixelMultiplier),
		int(float32(g.state.grid.pixelHeight) * g.state.grid.pixelMultiplier)
}

func main() {
	g := newGame()
	if m, ok := g.isDeveloperMacStudio(); ok {
		ebiten.SetMonitor(m)
	}

	os.Mkdir(g.outputDirectoryName, 0755)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(g Game) {
		for range c {
			logger.log(warnLevel, mainLogic, "received interrupt signal")
			g.forceCleanup()
		}
	}(g)

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
				logger.log(infoLevel, setupLogic,
					fmt.Sprintf("cacheSetups: cache obstacle %s", o[oI].string()))
			}
			if i > len(nodeSetups())-1 {
				nI = 0
			} else {
				logger.log(infoLevel, setupLogic,
					fmt.Sprintf("cacheSetups: cache node %s", n[nI].string()))
			}
			g.restart(n[nI], o[oI], false)
		}
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("cacheSetups: caching all missing setups required %d seconds",
				int(time.Since(start).Seconds())))
		return
	}

	g.restart(g.nodeSetup, g.obstacleSetup, true)

	if err := ebiten.RunGameWithOptions(&g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
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

func (g Game) isDeveloperMacStudio() (*ebiten.MonitorType, bool) {
	if g.developerMode {
		var monitors []*ebiten.MonitorType
		monitors = ebiten.AppendMonitors(monitors)
		for _, m := range monitors {
			if m.Name() == "DELL U2720Q" {
				return m, true
			}
		}
	}
	return nil, false
}

func (g *Game) registerToggles(tr *toggleRegister) {
	// setup cpu profiler
	toggleCPUProfiler := func() {
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
	}

	// setup memory profiler
	toggleMEMProfiler := func() {
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
	}

	if g.cpuProfilerStartBeforeGame && !g.cpuProfilerOn {
		toggleCPUProfiler()
	}
	if g.memProfilerStartBeforeGame && !g.memoryProfilerOn {
		toggleMEMProfiler()
	}

	tr.register(
		keyC,
		toggleCPUProfiler,
	)
	tr.register(
		keyX,
		toggleMEMProfiler,
	)
	tr.register(
		keyP,
		func() { g.paused = !g.paused },
	)
	tr.register(
		keyT,
		func() { g.slowTick = !g.slowTick },
	)
}
