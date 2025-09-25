// main entry for high level game logic
package main

import (
	"errors"
	"image/color"
	"log"
	"log/slog"
	"os"
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
	*drawContext
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
	drawFluxes               bool
	drawBlueTargets          bool
	drawRedTargets           bool
	drawDebugGrid            bool
	slowTick                 bool
	cpuProfilerOn            bool
	memoryProfilerOn         bool
	produceUnits             bool
	cpuProfilerFile          *os.File
	memoryProfilerFile       *os.File
	initialSelectedX         float32
	initialSelectedY         float32
	currentSelectedX         float32
	currentSelectedY         float32
	compOne                  computer
	compTwo                  computer
	gameState                state
	developerToggles         map[string]*toggle
	developerToggleCooldown  time.Duration
	initialNodeHealth        int
	loggerConfig             loggerConfig
}

func newGame() Game {
	g := Game{
		drawContext:             newDrawContext(),
		backgroundColor:         color.RGBA{200, 200, 200, 255},
		produceWait:             1_000_000 * time.Second,
		produceInterval:         3 * time.Second,
		slowTickInterval:        5 * time.Second,
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
		hotRestart:              true,
		cpuProfilerOn:           false,
		memoryProfilerOn:        false,
		developerToggles:        make(map[string]*toggle),
		developerToggleCooldown: 300 * time.Millisecond,
		initialNodeHealth:       100,
		loggerConfig: loggerConfig{
			initialLogLevel:            slog.LevelDebug,
			filterTimeAndLevelFromLogs: true,
			logLoggerLogs:              true,
			logMainLogs:                false,
			logStateLogs:               false,
			logSelectedUnitLogs:        true,
		},
	}
	return g
}

func (g *Game) restart() {
	// TODO: Need to restart timers, etc
	g.gameState = *newState(&g.drawContext.hexConfig)
	g.compOne = *newComputer(bluePlayer)
	g.compTwo = *newComputer(redPlayer)
	setupState(g.gameState, topVsBottomHorizontalOffset, &g.drawContext.hexConfig, g.initialNodeHealth)
	g.updateComputerTargetNodes()
}

func (g *Game) Update() error {
	g.currentTime = time.Now()

	for _, toggle := range g.developerToggles {
		toggle.checkToggle(g.currentTime)
	}

	// refresh
	g.leftMouseReleased = false
	g.leftMouseHeld = false

	// update game state prior to player actions
	if !g.paused && !(g.slowTick && (g.currentTime.Sub(g.slowTickTime) < g.slowTickInterval)) {
		if g.slowTick {
			logger.log(infoLevel, alwaysLogType, "slowTick")
		}
		g.slowTickTime = time.Now()

		g.gameState.refresh()
		g.gameState.analyzeTargets()
		g.gameState.processTerminalPositions()
		g.gameState.attackUnits()
		g.gameState.attackNodes()
		g.gameState.moveUnits()

		if g.produceWait < g.produceInterval {
			g.produceWait = g.produceWait + g.currentTime.Sub(g.priorTime)
		} else if g.produceUnits {
			g.produceWait = 0
			g.gameState.produceUnits()
		}
	}

	// player actions
	graphicalX, graphicalY := ebiten.CursorPosition()
	x := float32(graphicalX) / g.windowMultiplier
	y := float32(graphicalY) / g.windowMultiplier
	pixelPos := pixelPosition{int(x), int(y)}

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
	if g.middleMouseScrollClicked {
		logger.log(infoLevel, mainLogType, "middle click")
	}
	if g.rightMouseClicked {
		// TODO: Determine behavior. Should player choose which layer they target or handle it via game logic?
		// TODO: It seems a bit hacky that I need to pass computer active here. Determine if this is right behavior.
		g.gameState.updateTargets(position{ground, pixelPos}, g.blueComputerActive, g.redComputerActive)
	}

	// update game state after player actions
	if g.leftMouseReleased {
		g.gameState.updateSelected(
			pixelPosition{int(g.initialSelectedX), int(g.initialSelectedY)},
			pixelPosition{int(g.currentSelectedX), int(g.currentSelectedY)})
	}
	nodeOwnershipChanged := g.gameState.updateNodeOwnership()

	// play computer moves. Computers currently affect game state directly rather
	// than the same game state being sent to each one before a final update occurs.
	// This is a hack, but will probably be fine for awhile.
	if nodeOwnershipChanged {
		logger.log(infoLevel, mainLogType, "nodeOwnershipChanged")
		g.updateComputerTargetNodes()
	}

	var infantryCount int
	var sniperCount int
	var unguidedGroundMissileCount int
	var guidedAirMissileCount int
	var unknownCount int
	for _, u := range g.gameState.units {
		switch u.name {
		case infantry:
			infantryCount++
		case sniper:
			sniperCount++
		case unguidedGroundMissile:
			unguidedGroundMissileCount++
		case guidedAirMissile:
			guidedAirMissileCount++
		default:
			unknownCount++
		}

	}
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("total_units", len(g.gameState.units)))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("infantryCount", infantryCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("sniperCount", sniperCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("unguidedGroundMissile", unguidedGroundMissileCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("guidedAirMissile", guidedAirMissileCount))
	logger.log(infoLevel, mainLogType, "game_state", slog.Int("unknownCount", unknownCount))

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
			logger.log(infoLevel, mainLogType, "restartGame")

			g.restart()
		} else {
			return errors.New("regular termination")
		}
	}

	g.priorTime = g.currentTime
	return nil
}

func (g *Game) updateComputerTargetNodes() {
	if g.blueComputerActive {
		g.compOne.updateTargetNodes(g.gameState)
	}
	if g.redComputerActive {
		g.compTwo.updateTargetNodes(g.gameState)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawContext.screen == nil {
		g.drawContext.screen = screen
	}

	screen.Fill(g.drawContext.backgroundColor)
	drawDebug(*g.drawContext, g.drawDebugGrid)
	g.gameState.draw(*g.drawContext, g.drawBlueTargets, g.drawRedTargets, g.drawFluxes)

	if g.leftMouseHeld {
		drawUnfilledBoxFloat32(
			*g.drawContext,
			g.drawContext.highlightColor,
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

	g.developerToggles["cpuProfile"] = newToggle(
		ebiten.KeyC,
		func() {
			if g.cpuProfilerOn {
				logger.log(infoLevel, alwaysLogType, "Stopping CPU profiler.")
				pprof.StopCPUProfile()
				g.cpuProfilerFile.Close()
			} else {
				logger.log(infoLevel, alwaysLogType, "Starting CPU profiler.")
				f, err := os.Create("cpu.prof")
				if err != nil {
					logger.log(errorLevel, alwaysLogType, "Could not create CPU profiler file.")
				}
				g.cpuProfilerFile = f
				if err := pprof.StartCPUProfile(f); err != nil {
					logger.log(errorLevel, alwaysLogType, "Could not start CPU profiler.")
				}
			}
			g.cpuProfilerOn = !g.cpuProfilerOn
		},
		g.developerToggleCooldown)
	g.developerToggles["memoryProfile"] = newToggle(
		ebiten.KeyV,
		func() {
			if g.memoryProfilerOn {
				logger.log(infoLevel, alwaysLogType, "Stopping memory profiler.")

				if err := pprof.WriteHeapProfile(g.memoryProfilerFile); err != nil {
					logger.log(errorLevel, alwaysLogType, "Could not write heap profile.")
				}
				g.memoryProfilerFile.Close()
			} else {
				logger.log(infoLevel, alwaysLogType, "Starting memory profiler.")
				f, err := os.Create("mem.prof")
				if err != nil {
					logger.log(errorLevel, alwaysLogType, "Could not create memory profiler file.")
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
		func() { logger.toggleType(mainLogType) },
		g.developerToggleCooldown)
	g.developerToggles["stateLogType"] = newToggle(
		ebiten.KeyS,
		func() { logger.toggleType(stateLogType) },
		g.developerToggleCooldown)
	g.developerToggles["selectedUnitLogType"] = newToggle(
		ebiten.KeyU,
		func() { logger.toggleType(selectedUnitLogType) },
		g.developerToggleCooldown)
	g.developerToggles["drawAll"] = newToggle(
		ebiten.KeyA,
		func() {
			var toggle = !(g.drawFluxes && g.drawBlueTargets && g.drawRedTargets && g.drawDebugGrid)
			g.drawFluxes = toggle
			g.drawBlueTargets = toggle
			g.drawRedTargets = toggle
			g.drawDebugGrid = toggle
		},
		g.developerToggleCooldown)
	g.developerToggles["drawFlux"] = newToggle(
		ebiten.KeyF,
		func() { g.drawFluxes = !g.drawFluxes },
		g.developerToggleCooldown)
	g.developerToggles["drawBlueTargets"] = newToggle(
		ebiten.KeyB,
		func() { g.drawBlueTargets = !g.drawBlueTargets },
		g.developerToggleCooldown)
	g.developerToggles["drawRedTargets"] = newToggle(
		ebiten.KeyR,
		func() { g.drawRedTargets = !g.drawRedTargets },
		g.developerToggleCooldown)
	g.developerToggles["drawDebugGrid"] = newToggle(
		ebiten.KeyG,
		func() { g.drawDebugGrid = !g.drawDebugGrid },
		g.developerToggleCooldown)

	g.restart()

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

	if g.cpuProfilerOn {
		pprof.StopCPUProfile()
		g.cpuProfilerFile.Close()
	}
	if g.memoryProfilerOn {
		if err := pprof.WriteHeapProfile(g.memoryProfilerFile); err != nil {
			logger.log(errorLevel, alwaysLogType, "Could not write heap profile.")
		}
		g.memoryProfilerFile.Close()
	}
}
