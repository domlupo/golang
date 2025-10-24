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
	backgroundColor          color.RGBA
	produceWait              time.Duration
	produceInterval          time.Duration
	slowTickInterval         time.Duration
	slowTickTime             time.Time
	priorTime                time.Time
	currentTime              time.Time
	leftMouseClicked         bool
	leftMouseHeld            bool
	priorLeftMouseClicked    bool
	leftMouseReleased        bool
	middleMouseScrollClicked bool
	rightMouseClicked        bool
	paused                   bool
	blueComputerActive       bool
	redComputerActive        bool
	hotRestart               bool
	slowTick                 bool
	cpuProfilerOn            bool
	memoryProfilerOn         bool
	produceUnits             bool
	moveUnits                bool
	initialSelectedX         float32
	initialSelectedY         float32
	currentSelectedX         float32
	currentSelectedY         float32
	compOne                  computer
	compTwo                  computer
	gameState                *state
	debugger                 debugger
	developerMode            bool
	developerToggles         map[string]*toggle
	developerToggleCooldown  time.Duration
	initialNodeHealth        int
	pixelPathDistanceMax     int
	loggerConfig             loggerConfig
	outputDirectoryName      string
	cpuProfilerFileName      string
	memProfilerFileName      string
	cpuProfilerFile          *os.File
	memoryProfilerFile       *os.File
}

func newGame() Game {
	drawConfig := newDrawConfig()
	g := Game{
		drawConfig:              drawConfig,
		backgroundColor:         color.RGBA{200, 200, 200, 255},
		produceWait:             1_000_000 * time.Second,
		produceInterval:         3 * time.Second,
		slowTickInterval:        1 * time.Second,
		slowTickTime:            time.Now(),
		priorTime:               time.Now(),
		currentTime:             time.Now(),
		initialSelectedX:        float32(-1.0),
		initialSelectedY:        float32(-1.0),
		currentSelectedX:        float32(-1.0),
		currentSelectedY:        float32(-1.0),
		blueComputerActive:      true,
		redComputerActive:       true,
		produceUnits:            false,
		moveUnits:               true,
		hotRestart:              true,
		cpuProfilerOn:           false,
		memoryProfilerOn:        false,
		developerMode:           true,
		developerToggles:        make(map[string]*toggle),
		developerToggleCooldown: 300 * time.Millisecond,
		initialNodeHealth:       100,
		pixelPathDistanceMax:    drawConfig.hexWidth / 2,
		outputDirectoryName:     "output",
		cpuProfilerFileName:     "cpu.prof",
		memProfilerFileName:     "mem.prof",
		loggerConfig: loggerConfig{
			initialLogLevel:            slog.LevelInfo,
			filterTimeAndLevelFromLogs: true,
			logUndecidedBehavior:       true,
			logSetupLogic:              true,
			logDevTool:                 true,
		},
	}
	return g
}

func (g *Game) restart() {
	// TODO: Need to restart timers, etc
	g.compOne = *newComputer(bluePlayer)
	g.compTwo = *newComputer(redPlayer)

	hexagons := setupHexagons(&g.drawConfig.hexConfig)
	obstacles, pixelPath := setupObstacles(multiObstacle, &g.drawConfig.hexConfig, g.pixelPathDistanceMax)
	nodes, nodeLookup := setupNodes(leftVsRightVerticalOffset, &g.drawConfig.hexConfig, g.initialNodeHealth, hexagons, pixelPath)

	for nPos := range nodes {
		for oPos := range obstacles {
			if nPos.equals(oPos) {
				panic("grid position conflict between node and obstacle")
			}
		}
	}

	g.gameState = newState(nodes, obstacles, hexagons, pixelPath, nodeLookup, &g.drawConfig.hexConfig, g.developerMode)
	g.debugger = newDebugger(g.gameState, g.developerToggleCooldown)
	g.updateComputerTargetNodes()
}

func (g *Game) Update() error {
	g.priorTime = g.currentTime
	g.currentTime = time.Now()
	logger.log(
		infoLevel,
		mainLogic,
		"game performance",
		slog.Float64("fps", ebiten.ActualFPS()), slog.Float64("tps", ebiten.ActualTPS()),
	)

	for _, t := range g.developerToggles {
		t.checkToggle(g.currentTime)
	}
	for _, t := range g.debugger.toggles {
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
		g.slowTickTime = time.Now()

		g.gameState.analyzeTargets()
		g.gameState.processTerminalPositions()
		g.gameState.refresh()
		g.gameState.attackUnits()
		g.gameState.attackNodes()
		if g.moveUnits {
			g.gameState.moveUnits()
		}

		if g.produceWait < g.produceInterval {
			g.produceWait = g.produceWait + g.currentTime.Sub(g.priorTime)
		} else if g.produceUnits {
			g.produceWait = 0
			g.gameState.produceUnits()
		}
	}

	// player actions
	graphicalX, graphicalY := ebiten.CursorPosition()
	_, mouseWheelDelta := ebiten.Wheel()
	g.debugger.updateHeatmapMax(int(mouseWheelDelta))

	x := float32(graphicalX) / g.windowMultiplier
	y := float32(graphicalY) / g.windowMultiplier
	pixelPos := PixelPosition{int(x), int(y)}

	g.priorLeftMouseClicked = g.leftMouseClicked
	g.leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	g.middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	g.rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	g.currentSelectedX = x
	g.currentSelectedY = y
	if g.priorLeftMouseClicked && !g.leftMouseClicked {
		g.leftMouseReleased = true
	} else if !g.priorLeftMouseClicked && g.leftMouseClicked {
		g.initialSelectedX = x
		g.initialSelectedY = y
	} else if g.leftMouseClicked {
		g.leftMouseHeld = true
	}
	if g.developerMode && g.middleMouseScrollClicked {
		g.debugger.debug(pixelPos)
	}
	if g.rightMouseClicked {
		// TODO: Determine behavior. Should player choose which layer they target or handle it via game logic?
		// TODO: It seems a bit hacky that I need to pass computer active here. Determine if this is right behavior.
		g.gameState.updateTargets(position{ground, pixelPos}, g.blueComputerActive, g.redComputerActive)
	}

	// update game state after player actions
	if g.leftMouseReleased {
		g.gameState.updateSelected(
			PixelPosition{int(g.initialSelectedX), int(g.initialSelectedY)},
			PixelPosition{int(g.currentSelectedX), int(g.currentSelectedY)})
	}
	nodeOwnershipChanged := g.gameState.updateNodeOwnership()

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if nodeOwnershipChanged {
		logger.log(infoLevel, mainLogic, "nodeOwnershipChanged")
		g.updateComputerTargetNodes()
	}

	var infantryCount int
	var unguidedGroundMissileCount int
	var guidedAirMissileCount int
	var unknownCount int
	for _, u := range g.gameState.units {
		switch u.name {
		case infantry:
			infantryCount++
		case unguidedGroundMissile:
			unguidedGroundMissileCount++
		case guidedAirMissile:
			guidedAirMissileCount++
		default:
			unknownCount++
		}
	}
	logger.log(infoLevel, mainLogic, "unit counts",
		slog.Int("total_units", len(g.gameState.units)),
		slog.Int("infantryCount", infantryCount),
		slog.Int("unguidedGroundMissile", unguidedGroundMissileCount),
		slog.Int("guidedAirMissile", guidedAirMissileCount),
		slog.Int("unknownCount", unknownCount),
	)

	var blueNodeCount int
	var redNodeCount int

	for _, node := range g.gameState.nodes {
		switch node.player {
		case bluePlayer:
			blueNodeCount++
		case redPlayer:
			redNodeCount++
		}
	}

	if blueNodeCount == 0 || redNodeCount == 0 {
		if g.hotRestart {
			logger.log(infoLevel, mainLogic, "restartGame")

			g.restart()
		} else {
			return errors.New("regular termination")
		}
	}

	return nil
}

func (g *Game) updateComputerTargetNodes() {
	if g.blueComputerActive {
		g.compOne.updateTargetNodes(*g.gameState)
	}
	if g.redComputerActive {
		g.compTwo.updateTargetNodes(*g.gameState)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawConfig.screen == nil {
		g.drawConfig.screen = screen
	}

	screen.Fill(g.drawConfig.backgroundColor)
	g.gameState.draw(*g.drawConfig)
	g.debugger.draw(*g.drawConfig)

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
	return int(float32(g.windowWidth) * g.windowMultiplier), int(float32(g.windowHeight) * g.windowMultiplier)
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

	ebiten.SetWindowSize(int(float32(g.windowWidth)*g.windowMultiplier), int(float32(g.windowHeight)*g.windowMultiplier))

	logger = newLogger(g.loggerConfig)

	os.Mkdir(g.outputDirectoryName, 0755)

	g.developerToggles["cpuProfile"] = newToggle(
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
	g.developerToggles["memoryProfile"] = newToggle(
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
	g.developerToggles["logLevel"] = newToggle(
		ebiten.KeyL,
		func() { logger.toggleLevel() },
		g.developerToggleCooldown)
	g.developerToggles["mainLogType"] = newToggle(
		ebiten.KeyM,
		func() { logger.toggleType(mainLogic) },
		g.developerToggleCooldown)
	g.developerToggles["stateLogType"] = newToggle(
		ebiten.KeyS,
		func() { logger.toggleType(stateLogic) },
		g.developerToggleCooldown)
	g.developerToggles["selectedUnitLogType"] = newToggle(
		ebiten.KeyU,
		func() { logger.toggleType(selectedUnit) },
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

	g.restart()

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
