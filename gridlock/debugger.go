// A debugger for developers which outputs text logs and data visualizations to inspect game state.

package main

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debugger Configurations
type debuggerConfig struct {
	debugNodes                    bool
	debugNodeLUT                  bool
	debugObstacleLUT              bool
	debugVisibilityGrids          bool
	debugPOneTargets              bool
	debugPTwoTargets              bool
	debugHexagons                 bool
	visibilityGridLineWidth       int
	visibilityGridPointRadius     int
	visibilityGridHighlightOffset int
	pixelSpotlightOffset          int
	heatmapR                      uint8
	heatmapG                      uint8
	heatmapB                      uint8
	heatmapAlphaMaxLimit          int
	heatmapPercent                int
	partitionAnalysis             partitionAnalysis
	roygbivColors                 []color.RGBA
}

func newDebuggerConfig() *debuggerConfig {
	dc := debuggerConfig{
		visibilityGridLineWidth:       1,
		visibilityGridPointRadius:     3,
		visibilityGridHighlightOffset: 3,
		pixelSpotlightOffset:          5,
		heatmapR:                      0,
		heatmapG:                      0,
		heatmapB:                      0,
		heatmapAlphaMaxLimit:          200,
		partitionAnalysis:             paVision,
		heatmapPercent:                maxPercent,
		roygbivColors: []color.RGBA{
			{200, 0, 0, 255},
			{250, 150, 0, 255},
			{250, 200, 0, 255},
			{0, 200, 150, 255},
			{0, 0, 250, 255},
			{0, 0, 150, 255},
			{150, 150, 250, 255},
		},
	}
	return &dc
}

// A debugger which outputs text logs and data visualizations to inspect game state.
type debugger struct {
	*state
	*debuggerConfig
	pixel                       Pixel
	obstacleLUTUnwalkable       heatmap
	obstacleLUTDistance         heatmap
	nodeLUTDistance             heatmap
	heatmapAlphaMax             int
	unitPathfind                heatmap
	unitPathfindHeatmapComplete bool
	toggles                     map[string]*toggle
}

// A typical heatmap. Used for visual debugging.
type heatmap map[Pixel]int

func newDebugger(state *state, toggleCooldown time.Duration) *debugger {
	c := newDebuggerConfig()
	debugger := debugger{
		state:           state,
		debuggerConfig:  c,
		heatmapAlphaMax: c.heatmapAlphaMaxLimit,
		unitPathfind:    make(heatmap),
	}
	debugger.setupHeatmaps()
	debugger.setupToggles(toggleCooldown)
	return &debugger
}

func (d *debugger) setupHeatmaps() {
	d.obstacleLUTUnwalkable = make(heatmap)
	d.obstacleLUTDistance = make(heatmap)
	d.nodeLUTDistance = make(heatmap)

	for pixel, info := range d.oLUT {
		switch info.Walkable {
		case true:
			d.obstacleLUTUnwalkable[pixel] = 0
		case false:
			d.obstacleLUTUnwalkable[pixel] = 1
		}
	}

	for pixel, info := range d.oLUT {
		d.obstacleLUTDistance[pixel] = info.UnwalkableChebyshevDistance
	}

	for pixel, info := range d.nLUT {
		d.nodeLUTDistance[pixel] = info.ChebyshevDistance
	}
}

func (d *debugger) setupToggles(cooldown time.Duration) {
	toggles := make(map[string]*toggle)
	toggles["debugAll"] = newToggle(
		ebiten.KeyA,
		func() {
			var toggle = !(d.debugPOneTargets &&
				d.debugPTwoTargets &&
				d.debugVisibilityGrids &&
				d.debugHexagons &&
				d.debugNodes)
			d.debugPOneTargets = toggle
			d.debugPTwoTargets = toggle
			d.debugVisibilityGrids = toggle
			d.debugHexagons = toggle
			d.debugNodes = toggle
		},
		cooldown)
	toggles["debugBlueTargets"] = newToggle(
		ebiten.KeyB,
		func() { d.debugPOneTargets = !d.debugPOneTargets },
		cooldown)
	toggles["debugRedTargets"] = newToggle(
		ebiten.KeyR,
		func() { d.debugPTwoTargets = !d.debugPTwoTargets },
		cooldown)
	toggles["debugVisibilityGrids"] = newToggle(
		ebiten.KeyV,
		func() { d.debugVisibilityGrids = !d.debugVisibilityGrids },
		cooldown)
	toggles["debugNodes"] = newToggle(
		ebiten.KeyN,
		func() { d.debugNodes = !d.debugNodes },
		cooldown)
	toggles["debugNodeLUT"] = newToggle(
		ebiten.KeyW,
		func() { d.debugNodeLUT = !d.debugNodeLUT },
		cooldown)
	toggles["debugObstacleLUT"] = newToggle(
		ebiten.KeyQ,
		func() { d.debugObstacleLUT = !d.debugObstacleLUT },
		cooldown)
	toggles["debugHexagons"] = newToggle(
		ebiten.KeyH,
		func() { d.debugHexagons = !d.debugHexagons },
		cooldown)
	d.toggles = toggles
}

func (d *debugger) debug(p Pixel) {
	d.logPixel(p)
	d.logVisibilityGridObstaclePoints()
	d.logHexagons()
	d.logNodes()
	d.createUnitPathfindHeatmap()
}

func (d *debugger) logPixel(p Pixel) {
	attrs := []slog.Attr{nest("Pixel", p.log())}
	if gridPos, ok := d.grid.getGridPos(p); ok {
		attrs = append(attrs, nest("GridPosition", gridPos.log()))
	} else {
		attrs = append(attrs, slog.String("GridPosition", "notFound"))
	}
	attrs = append(attrs, nest("ObsacleLUT", d.oLUT[p].log()))
	attrs = append(attrs, nest("NodeLUT", d.nLUT[p].log()))
	logger.log(infoLevel, devTool, "pixel", attrs...)
}

func (d *debugger) logVisibilityGridObstaclePoints() {
	if d.debugVisibilityGrids {
		visited := make(map[id]struct{})
		for _, o := range d.obstacles {
			g := *o.visibilityGrid
			op := g[0]
			if _, ok := visited[op.visibilityGridId]; ok {
				continue
			}
			visited[op.visibilityGridId] = struct{}{}
			for _, point := range *o.visibilityGrid {
				if point.selected {
					logger.log(infoLevel, devTool, "obstaclePoint", point.log()...)
				}
			}
		}
	}
}

func (d *debugger) logHexagons() {
	if d.debugHexagons {
		for _, h := range d.hexagons {
			if h.selected {
				logger.log(infoLevel, devTool, "hexagon", h.log()...)
			}
		}
	}
}

func (d *debugger) logNodes() {
	if d.debugNodes {
		for _, n := range d.nodes {
			if n.selected {
				logger.log(infoLevel, devTool, "node", n.log()...)
			}
		}
	}
}

func (d *debugger) selectedUnit() *unit {
	var unit *unit
	for u := range d.units() {
		if u.selected && unit != nil {
			logger.log(debugLevel, devTool, "selectedUnit: cannot select multiple units to debug")
			return nil
		} else if u.selected {
			unit = u
		}
	}
	return unit
}

func (d *debugger) createUnitPathfindHeatmap() {
	logger.onlyEnable(devTool)

	if len(d.unitPathfind) != 0 {
		logger.log(infoLevel, devTool, "createUnitPathfindHeatmap: cannot create heatmap since it was already attempted")
		return
	}

	unit := d.selectedUnit()
	if unit == nil {
		logger.log(debugLevel, devTool, "createUnitPathfindHeatmap: cannot select no units")
		return
	}

	logger.log(infoLevel, devTool, "createUnitPathfindHeatmap: start")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	pixelCount := 0
	for pixel, info := range d.oLUT {
		select {
		case <-ctx.Done():
			return
		default:
			pixelCount++
			if pixelCount%(d.grid.basePixels/maxPercent) == 0 {
				percentProgess := pixelCount / (d.grid.basePixels / maxPercent)
				logger.log(infoLevel, devTool, "createUnitPathfindHeatmap", slog.Int("pixel positions percent progress", percentProgess))
			}

			if info.Walkable {
				unit.updateTarget(globalTarget, position{layer: unit.layer, Pixel: pixel})
				distance := unit.pathfindGlobalTarget(d.obstacles, d.oLUT, d.grid)
				d.unitPathfind[pixel] = distance
			}
		}
	}

	d.unitPathfindHeatmapComplete = true
	logger.log(infoLevel, devTool, "createUnitPathfindHeatmap: done")
}

func (d *debugger) updateHeatmapPercent(delta int) {
	if delta == 0 || !(d.debugObstacleLUT || d.debugNodeLUT) {
		return
	}
	d.heatmapPercent = d.heatmapPercent + delta
	d.heatmapPercent = max(d.heatmapPercent, minPercent+1)
	d.heatmapPercent = min(d.heatmapPercent, maxPercent)
	logger.log(infoLevel, devTool, fmt.Sprintf("updateHeatmapPercent: Set to %d", d.heatmapPercent))
}

func (d debugger) draw(c drawConfig) {
	if d.debugObstacleLUT {
		d.obstacleLUTUnwalkable.draw(c, *d.debuggerConfig, d.heatmapPercent, d.grid.pixelMultiplier)
		d.obstacleLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent, d.grid.pixelMultiplier)

		if d.grid.contains(d.pixel) {
			p := d.oLUT[d.pixel].NearestUnwalkablesAvg
			p.draw(c, *d.debuggerConfig, c.highlightColor, true, d.grid.pixelMultiplier)
		}
	}

	if d.debugNodeLUT {
		d.nodeLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent, d.grid.pixelMultiplier)

		for _, p := range d.nLUT[d.pixel].NearestNodes {
			drawHexHighlight(c, d.grid, p)
		}
	}

	if d.unitPathfindHeatmapComplete {
		d.unitPathfind.draw(c, *d.debuggerConfig, d.heatmapAlphaMax, d.grid.pixelMultiplier)
	}

	for _, n := range d.nodes {
		n.drawTargets(c, d.grid, d.debugPOneTargets, d.debugPTwoTargets)
	}

	for u := range d.units() {
		u.drawDebug(c, *d.debuggerConfig, d.grid.pixelMultiplier)
	}

	if d.debugHexagons {
		u := d.selectedUnit()
		if u != nil {
			logger.log(infoLevel, devTool, "debugHexagons", d.partitionAnalysis.log())
			partitions := d.partitions(*u, d.partitionAnalysis)
			for _, h := range partitions {
				drawHexHighlight(c, d.grid, h.gridPos)
			}
		} else {
			for _, h := range d.hexagons {
				center := d.grid.getPixelOutbounds(h.gridPos, hpCenter)
				center.draw(c, *d.debuggerConfig, color.White, false, d.grid.pixelMultiplier)
				if h.selected {
					drawHexHighlight(c, d.grid, h.gridPos)
				}
			}
		}
	}

	if d.debugVisibilityGrids {
		visited := make(map[id]struct{})
		for _, o := range d.obstacles {
			g := *o.visibilityGrid
			op := g[0]
			if _, ok := visited[op.visibilityGridId]; ok {
				continue
			}
			visited[op.visibilityGridId] = struct{}{}
			o.drawVisibilityGrid(c, *d.debuggerConfig, d.grid.pixelMultiplier)
		}
	}
}

func (p Pixel) draw(c drawConfig, dC debuggerConfig, color color.Color, spotlight bool, windowpixelMultiplier float32) {
	vector.FillRect(
		c.screen,
		float32(p.X)*windowpixelMultiplier,
		float32(p.Y)*windowpixelMultiplier,
		float32(c.pixelSize)*windowpixelMultiplier,
		float32(c.pixelSize)*windowpixelMultiplier,
		color,
		false)

	if spotlight {
		drawUnfilledBoxInt(
			c,
			c.highlightColor,
			p.X-dC.pixelSpotlightOffset,
			p.Y-dC.pixelSpotlightOffset,
			p.X+dC.pixelSpotlightOffset,
			p.Y+dC.pixelSpotlightOffset,
			windowpixelMultiplier)
	}
}

func (heatmap heatmap) draw(c drawConfig, dC debuggerConfig, heatmapPercent int, windowpixelMultiplier float32) {
	var maxV int
	for _, v := range heatmap {
		maxV = max(maxV, v)
	}
	trimmedV := maxV * heatmapPercent / maxPercent
	for pixel, v := range heatmap {
		v := min(v, trimmedV)
		color := color.RGBA{
			dC.heatmapR,
			dC.heatmapB,
			dC.heatmapG,
			uint8(float64(v) / float64(trimmedV) * float64(dC.heatmapAlphaMaxLimit))}
		vector.FillRect(
			c.screen,
			float32(pixel.X)*windowpixelMultiplier,
			float32(pixel.Y)*windowpixelMultiplier,
			1*windowpixelMultiplier,
			1*windowpixelMultiplier,
			color,
			false)
	}
}

func (o obstacle) drawVisibilityGrid(c drawConfig, dC debuggerConfig, windowpixelMultiplier float32) {
	for _, op := range *o.visibilityGrid {
		vector.FillCircle(
			c.screen,
			float32(op.X)*windowpixelMultiplier,
			float32(op.Y)*windowpixelMultiplier,
			float32(dC.visibilityGridPointRadius)*windowpixelMultiplier,
			c.highlightColor,
			false)
		vector.StrokeLine(
			c.screen,
			float32(op.X)*windowpixelMultiplier,
			float32(op.Y)*windowpixelMultiplier,
			float32(op.cw.X)*windowpixelMultiplier,
			float32(op.cw.Y)*windowpixelMultiplier,
			float32(dC.visibilityGridLineWidth)*windowpixelMultiplier,
			c.highlightColor,
			false)

		if op.selected {
			offset := dC.visibilityGridPointRadius + dC.visibilityGridHighlightOffset
			drawUnfilledBoxInt(
				c,
				c.highlightColor,
				op.X-offset,
				op.Y-offset,
				op.X+offset,
				op.Y+offset,
				windowpixelMultiplier,
			)
		}
	}
}

func (u *unit) drawDebug(c drawConfig, dC debuggerConfig, windowpixelMultiplier float32) {
	if u.selected && !u.activeTarget.equals(position{}) {
		vector.StrokeLine(
			c.screen,
			float32(u.X)*windowpixelMultiplier,
			float32(u.Y)*windowpixelMultiplier,
			float32(u.activeTarget.X)*windowpixelMultiplier,
			float32(u.activeTarget.Y)*windowpixelMultiplier,
			float32(c.targetNodeLineWidth)*windowpixelMultiplier,
			color.Black,
			false,
		)

		vector.StrokeLine(
			c.screen,
			float32(u.X)*windowpixelMultiplier,
			float32(u.Y)*windowpixelMultiplier,
			float32(u.activeTarget.X)*windowpixelMultiplier,
			float32(u.activeTarget.Y)*windowpixelMultiplier,
			float32(c.targetNodeLineWidth)*windowpixelMultiplier,
			color.Black,
			false,
		)

		j := 0
		for i := len(u.localTargets) - 1; i >= 0; i-- {
			lt := u.localTargets[i]
			color := dC.roygbivColors[j]

			vector.StrokeLine(
				c.screen,
				float32(u.X)*windowpixelMultiplier,
				float32(u.Y)*windowpixelMultiplier,
				float32(lt.X)*windowpixelMultiplier,
				float32(lt.Y)*windowpixelMultiplier,
				float32(c.targetNodeLineWidth)*windowpixelMultiplier,
				color,
				false,
			)

			j += 1
			if j == len(dC.roygbivColors) {
				j = 0
			}
		}

		j = 0
		for i := len(u.globalTargets) - 1; i >= 0; i-- {
			gt := u.globalTargets[i]
			color := dC.roygbivColors[j]

			vector.StrokeLine(
				c.screen,
				float32(u.X)*windowpixelMultiplier,
				float32(u.Y)*windowpixelMultiplier,
				float32(gt.X)*windowpixelMultiplier,
				float32(gt.Y)*windowpixelMultiplier,
				float32(c.targetNodeLineWidth)*windowpixelMultiplier,
				color,
				false,
			)

			j += 1
			if j == len(dC.roygbivColors) {
				j = 0
			}
		}
	}
}
