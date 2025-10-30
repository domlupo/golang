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

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debugger Configurations
type debuggerConfig struct {
	debugHexGrid                  bool
	debugObstacleLUT              bool
	debugNodeLUT                  bool
	debugVisibilityGrids          bool
	debugBlueTargets              bool
	debugRedTargets               bool
	debugHexagons                 bool
	hexGridLineWidth              int
	hexGridIndexCount             int
	debugGridColor                color.Color
	visibilityGridLineWidth       int
	visibilityGridPointRadius     int
	visibilityGridHighlightOffset int
	pixelSpotlightOffset          int
	heatmapR                      uint8
	heatmapG                      uint8
	heatmapB                      uint8
	heatmapAlphaMaxLimit          int
	heatmapPercent                int
	roygbivColors                 []color.RGBA
}

func newDebuggerConfig() *debuggerConfig {
	dc := debuggerConfig{
		hexGridLineWidth:              1,
		hexGridIndexCount:             32,
		debugGridColor:                color.Black,
		visibilityGridLineWidth:       1,
		visibilityGridPointRadius:     3,
		visibilityGridHighlightOffset: 3,
		pixelSpotlightOffset:          5,
		heatmapR:                      0,
		heatmapG:                      0,
		heatmapB:                      0,
		heatmapAlphaMaxLimit:          200,
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
	pixel                       PixelPosition
	obstacleLUTUnwalkable       heatmap
	obstacleLUTDistance         heatmap
	nodeLUTDistance             heatmap
	heatmapAlphaMax             int
	unitPathfind                heatmap
	unitPathfindHeatmapComplete bool
	toggles                     map[string]*toggle
}

// A typical heatmap. Used for visual debugging.
type heatmap map[PixelPosition]int

func newDebugger(state *state, toggleCooldown time.Duration) debugger {
	c := newDebuggerConfig()
	debugger := debugger{
		state:           state,
		debuggerConfig:  c,
		heatmapAlphaMax: c.heatmapAlphaMaxLimit,
		unitPathfind:    make(heatmap),
	}
	debugger.setupHeatmaps()
	debugger.setupToggles(toggleCooldown)
	return debugger
}

func (d *debugger) setupHeatmaps() {
	d.obstacleLUTUnwalkable = make(heatmap)
	d.obstacleLUTDistance = make(heatmap)
	d.nodeLUTDistance = make(heatmap)

	for pixelPos, info := range d.oLUT {
		switch info.Walkable {
		case true:
			d.obstacleLUTUnwalkable[pixelPos] = 0
		case false:
			d.obstacleLUTUnwalkable[pixelPos] = 1
		}
	}

	for pixelPos, info := range d.oLUT {
		// TODO: Revist once local path finding is more polished Should this be unwalkable or smooth sum?
		d.obstacleLUTDistance[pixelPos] = info.UnwalkableChebyshevDistance
	}

	for pixelPos, info := range d.nLUT {
		d.nodeLUTDistance[pixelPos] = info.ChebyshevDistance
	}
}

func (d *debugger) setupToggles(cooldown time.Duration) {
	toggles := make(map[string]*toggle)
	toggles["drawAll"] = newToggle(
		ebiten.KeyA,
		func() {
			var toggle = !(d.debugBlueTargets && d.debugRedTargets && d.debugHexGrid && d.debugVisibilityGrids && d.debugHexagons)
			d.debugBlueTargets = toggle
			d.debugRedTargets = toggle
			d.debugHexGrid = toggle
			d.debugVisibilityGrids = toggle
			d.debugHexagons = toggle

		},
		cooldown)
	toggles["debugBlueTargets"] = newToggle(
		ebiten.KeyB,
		func() { d.debugBlueTargets = !d.debugBlueTargets },
		cooldown)
	toggles["debugRedTargets"] = newToggle(
		ebiten.KeyR,
		func() { d.debugRedTargets = !d.debugRedTargets },
		cooldown)
	toggles["debugHexGrid"] = newToggle(
		ebiten.KeyG,
		func() { d.debugHexGrid = !d.debugHexGrid },
		cooldown)
	toggles["debugVisibilityGrids"] = newToggle(
		ebiten.KeyV,
		func() { d.debugVisibilityGrids = !d.debugVisibilityGrids },
		cooldown)
	toggles["debugObstacleLUT"] = newToggle(
		ebiten.KeyQ,
		func() { d.debugObstacleLUT = !d.debugObstacleLUT },
		cooldown)
	toggles["debugNodeLUT"] = newToggle(
		ebiten.KeyW,
		func() { d.debugNodeLUT = !d.debugNodeLUT },
		cooldown)
	toggles["debugHexagons"] = newToggle(
		ebiten.KeyH,
		func() { d.debugHexagons = !d.debugHexagons },
		cooldown)
	d.toggles = toggles
}

func (d *debugger) debug(p PixelPosition) {
	d.logPixel(p)
	d.logVisibilityGridObstaclePoints()
	d.logHexagons()
	d.createUnitPathfindHeatmap()
}

func (d *debugger) logPixel(p PixelPosition) {
	attrs := []slog.Attr{nest("PixelPosition", p.log())}
	if gridPos, ok := d.g.getGridPos(p); ok {
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
		visited := make(map[uuid.UUID]struct{})
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

func (d *debugger) createUnitPathfindHeatmap() {
	logger.onlyEnable(devTool)

	if len(d.unitPathfind) != 0 {
		logger.log(infoLevel, devTool, "cannot create heatmap since it was already attempted")
		return
	}

	var unit *unit
	for _, u := range d.units {
		if u.selected && unit != nil {
			logger.log(debugLevel, devTool, "can only select one unit when creating heatmap")
			return
		} else if u.selected {
			unit = u
		}
	}

	if unit == nil {
		logger.log(debugLevel, devTool, "must select one unit when creating heat map")
		return
	}

	logger.log(infoLevel, devTool, "start creating pathfinding heatmap")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pixelCount := 0
	wc := newWindowConfig()
	for pixelPos, info := range d.oLUT {
		select {
		case <-ctx.Done():
			return
		default:
			pixelCount++
			if pixelCount%(wc.maxPixels()/maxPercent) == 0 {
				percentProgess := pixelCount / (wc.maxPixels() / maxPercent)
				logger.log(infoLevel, devTool, "pathfind heatmap", slog.Int("pixel positions percent progress", percentProgess))
			}

			if info.Walkable {
				unit.updateTarget(globalTarget, position{layer: unit.layer, PixelPosition: pixelPos})
				distance := unit.pathfindGlobalTarget(d.obstacles, d.oLUT, d.g)
				d.unitPathfind[pixelPos] = distance
			}
		}
	}

	d.unitPathfindHeatmapComplete = true
	logger.log(infoLevel, devTool, "done creating pathfinding heatmap")
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
	if d.debugHexGrid {
		drawHexGrid(c, *d.debuggerConfig)
	}

	if d.debugObstacleLUT {
		d.obstacleLUTUnwalkable.draw(c, *d.debuggerConfig, d.heatmapPercent)
		d.obstacleLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent)

		if newWindowConfig().inbounds(d.pixel) {
			p := d.oLUT[d.pixel].NearestUnwalkablesAvg
			p.draw(c, *d.debuggerConfig, c.highlightColor, true)
		}
	}

	if d.debugNodeLUT {
		d.nodeLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent)

		for _, p := range d.nLUT[d.pixel].NearestNodes {
			drawHexHighlight(c, p)
		}
	}

	if d.unitPathfindHeatmapComplete {
		d.unitPathfind.draw(c, *d.debuggerConfig, d.heatmapAlphaMax)
	}

	for _, n := range d.nodes {
		n.drawTargets(c, d.debugBlueTargets, d.debugRedTargets)
	}

	for _, u := range d.units {
		u.drawDebug(c, *d.debuggerConfig)
	}

	if d.debugHexagons {
		for _, h := range d.hexagons {
			center := d.g.getPixelPosOutbounds(h.gridPos, hpCenter)
			center.draw(c, *d.debuggerConfig, color.Black, false)

			if h.selected {
				drawHexHighlight(c, h.gridPos)
			}
		}
	}

	if d.debugVisibilityGrids {
		visited := make(map[uuid.UUID]struct{})
		for _, o := range d.obstacles {
			g := *o.visibilityGrid
			op := g[0]
			if _, ok := visited[op.visibilityGridId]; ok {
				continue
			}
			visited[op.visibilityGridId] = struct{}{}
			o.drawVisibilityGrid(c, *d.debuggerConfig)
		}
	}
}

func (p PixelPosition) draw(c drawConfig, dC debuggerConfig, color color.Color, spotlight bool) {
	vector.DrawFilledRect(
		c.screen,
		float32(p.X)*c.windowMultiplier,
		float32(p.Y)*c.windowMultiplier,
		float32(c.pixelSize)*c.windowMultiplier,
		float32(c.pixelSize)*c.windowMultiplier,
		color,
		false)

	if spotlight {
		drawUnfilledBoxInt(
			c,
			c.highlightColor,
			p.X-dC.pixelSpotlightOffset,
			p.Y-dC.pixelSpotlightOffset,
			p.X+dC.pixelSpotlightOffset,
			p.Y+dC.pixelSpotlightOffset)
	}
}

func (heatmap heatmap) draw(c drawConfig, dC debuggerConfig, heatmapPercent int) {
	var maxV int
	for _, v := range heatmap {
		maxV = max(maxV, v)
	}
	trimmedV := maxV * heatmapPercent / maxPercent
	for pixelPos, v := range heatmap {
		v := min(v, trimmedV)
		color := color.RGBA{
			dC.heatmapR,
			dC.heatmapB,
			dC.heatmapG,
			uint8(float64(v) / float64(trimmedV) * float64(dC.heatmapAlphaMaxLimit))}
		vector.DrawFilledRect(
			c.screen,
			float32(pixelPos.X)*c.windowMultiplier,
			float32(pixelPos.Y)*c.windowMultiplier,
			1*c.windowMultiplier,
			1*c.windowMultiplier,
			color,
			false)
	}
}

func (o obstacle) drawVisibilityGrid(c drawConfig, dC debuggerConfig) {
	for _, op := range *o.visibilityGrid {
		vector.DrawFilledCircle(
			c.screen,
			float32(op.X)*c.windowMultiplier,
			float32(op.Y)*c.windowMultiplier,
			float32(dC.visibilityGridPointRadius)*c.windowMultiplier,
			c.highlightColor,
			false)
		vector.StrokeLine(
			c.screen,
			float32(op.X)*c.windowMultiplier,
			float32(op.Y)*c.windowMultiplier,
			float32(op.cw.X)*c.windowMultiplier,
			float32(op.cw.Y)*c.windowMultiplier,
			float32(dC.visibilityGridLineWidth)*c.windowMultiplier,
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
			)
		}
	}
}

func (u *unit) drawDebug(c drawConfig, dC debuggerConfig) {
	if u.selected && !u.activeTarget.isZero() {
		vector.StrokeLine(
			c.screen,
			float32(u.X)*c.windowMultiplier,
			float32(u.Y)*c.windowMultiplier,
			float32(u.activeTarget.X)*c.windowMultiplier,
			float32(u.activeTarget.Y)*c.windowMultiplier,
			float32(c.targetNodeLineWidth)*c.windowMultiplier,
			color.Black,
			false,
		)

		vector.StrokeLine(
			c.screen,
			float32(u.X)*c.windowMultiplier,
			float32(u.Y)*c.windowMultiplier,
			float32(u.activeTarget.X)*c.windowMultiplier,
			float32(u.activeTarget.Y)*c.windowMultiplier,
			float32(c.targetNodeLineWidth)*c.windowMultiplier,
			color.Black,
			false,
		)

		j := 0
		for i := len(u.localTargets) - 1; i >= 0; i-- {
			lt := u.localTargets[i]
			color := dC.roygbivColors[j]

			vector.StrokeLine(
				c.screen,
				float32(u.X)*c.windowMultiplier,
				float32(u.Y)*c.windowMultiplier,
				float32(lt.X)*c.windowMultiplier,
				float32(lt.Y)*c.windowMultiplier,
				float32(c.targetNodeLineWidth)*c.windowMultiplier,
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

// This is likely to be become apart of game and then it should be removed from debugger.go
func drawHexGrid(c drawConfig, dC debuggerConfig) {
	// This is the pattern drawn.
	//    / \
	//   |
	//   |
	//
	// when repeated, it forms a hexagonal grid.
	//    / \ / \
	//   |   |
	//   |   |
	//  / \ / \ / \
	// |   |   |
	// |   |   |
	//    / \ / \
	//   |   |
	//   |   |

	for i := 0; i < c.grid.xIndexMax; i++ {
		for j := 0; j < c.grid.yIndexMax; j++ {
			if j%2 == 0 && i == c.grid.xIndexMax-1 {
				continue
			}

			p := GridPosition{i, j}

			n := c.grid.getPixelPosOutbounds(p, hpNorth)
			nX := float32(n.X) * c.windowMultiplier
			nY := float32(n.Y) * c.windowMultiplier

			nw := c.grid.getPixelPosOutbounds(p, hpNorthwest)
			nwX := float32(nw.X) * c.windowMultiplier
			nwY := float32(nw.Y) * c.windowMultiplier

			ne := c.grid.getPixelPosOutbounds(p, hpNortheast)
			neX := float32(ne.X+1) * c.windowMultiplier
			neY := float32(ne.Y) * c.windowMultiplier

			sw := c.grid.getPixelPosOutbounds(p, hpSouthwest)
			swX := float32(sw.X) * c.windowMultiplier
			swY := float32(sw.Y+1) * c.windowMultiplier

			// draw west side
			//    / \
			//  *|*
			//  *|*
			vector.StrokeLine(
				c.screen,
				nwX,
				nwY,
				swX,
				swY,
				float32(dC.hexGridLineWidth)*c.windowMultiplier,
				dC.debugGridColor,
				false)

			// draw northwest
			//   */*\
			//   |
			//   |
			vector.StrokeLine(
				c.screen,
				nX,
				nY,
				nwX,
				nwY,
				float32(dC.hexGridLineWidth)*c.windowMultiplier,
				dC.debugGridColor,
				false)

			// draw northeast
			//    /*\*
			//   |
			//   |
			vector.StrokeLine(
				c.screen,
				nX,
				nY,
				neX,
				neY,
				float32(dC.hexGridLineWidth)*c.windowMultiplier,
				dC.debugGridColor,
				false)

		}
	}
}
