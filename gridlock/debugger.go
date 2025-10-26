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
	drawHexGrid             bool
	drawObstacleLUT         bool
	drawNodeLUT             bool
	drawVisibilityGrids     bool
	drawBlueTargets         bool
	drawRedTargets          bool
	debugGridWidth          int
	debugGridIndexCount     int
	debugGridColor          color.Color
	visibilityGridLineWidth int
	pixelSpotlightOffset    int
	heatmapR                uint8
	heatmapG                uint8
	heatmapB                uint8
	heatmapAlphaMaxLimit    int
	heatmapPercent          int
	roygbivColors           []color.RGBA
}

func newDebuggerConfig() *debuggerConfig {
	dc := debuggerConfig{
		debugGridWidth:          1,
		debugGridIndexCount:     32,
		debugGridColor:          color.Black,
		visibilityGridLineWidth: 1,
		pixelSpotlightOffset:    5,
		heatmapR:                0,
		heatmapG:                0,
		heatmapB:                0,
		heatmapAlphaMaxLimit:    200,
		heatmapPercent:          MaxPercent,
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
		// TODO: Update this to unwalkable smooth sum once debugger supports percentages
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
			var toggle = !(d.drawBlueTargets && d.drawRedTargets && d.drawHexGrid && d.drawVisibilityGrids)
			d.drawBlueTargets = toggle
			d.drawRedTargets = toggle
			d.drawHexGrid = toggle
			d.drawVisibilityGrids = toggle
		},
		cooldown)
	toggles["drawBlueTargets"] = newToggle(
		ebiten.KeyB,
		func() { d.drawBlueTargets = !d.drawBlueTargets },
		cooldown)
	toggles["drawRedTargets"] = newToggle(
		ebiten.KeyR,
		func() { d.drawRedTargets = !d.drawRedTargets },
		cooldown)
	toggles["drawDebugGrid"] = newToggle(
		ebiten.KeyG,
		func() { d.drawHexGrid = !d.drawHexGrid },
		cooldown)
	toggles["drawVisibilityGrids"] = newToggle(
		ebiten.KeyV,
		func() { d.drawVisibilityGrids = !d.drawVisibilityGrids },
		cooldown)
	toggles["drawObstacleLUT"] = newToggle(
		ebiten.KeyQ,
		func() { d.drawObstacleLUT = !d.drawObstacleLUT },
		cooldown)
	toggles["drawNodeLUT"] = newToggle(
		ebiten.KeyW,
		func() { d.drawNodeLUT = !d.drawNodeLUT },
		cooldown)
	d.toggles = toggles
}

func (d *debugger) debug(p PixelPosition) {
	d.logPixel(p)
	d.createUnitPathfindHeatmap()
}

func (d *debugger) logPixel(p PixelPosition) {
	attrs := []slog.Attr{nest("PixelPosition", p.log())}
	if gridPos, ok := p.getGridPos(d.hc); ok {
		attrs = append(attrs, nest("GridPosition", gridPos.log()))
	} else {
		attrs = append(attrs, slog.String("GridPosition", "notFound"))
	}
	attrs = append(attrs, nest("ObsacleLUT", d.oLUT[p].log()))
	attrs = append(attrs, nest("NodeLUT", d.nLUT[p].log()))
	logger.log(infoLevel, devTool, "debugPixel", attrs...)
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
			if pixelCount%(wc.maxPixels()/MaxPercent) == 0 {
				percentProgess := pixelCount / (wc.maxPixels() / MaxPercent)
				logger.log(infoLevel, devTool, "pathfind heatmap", slog.Int("pixel positions percent progress", percentProgess))
			}

			if info.Walkable {
				unit.updateTarget(globalTarget, position{layer: unit.layer, PixelPosition: pixelPos})
				distance := unit.pathfindGlobalTarget(d.obstacles, d.oLUT, d.hc)
				d.unitPathfind[pixelPos] = distance
			}
		}
	}

	d.unitPathfindHeatmapComplete = true
	logger.log(infoLevel, devTool, "done creating pathfinding heatmap")
}

func (d *debugger) updateHeatmapPercent(delta int) {
	if delta == 0 {
		return
	}
	d.heatmapPercent = d.heatmapPercent + delta
	d.heatmapPercent = max(d.heatmapPercent, MinPercent+1)
	d.heatmapPercent = min(d.heatmapPercent, MaxPercent)
	logger.log(infoLevel, devTool, fmt.Sprintf("updateHeatmapPercent: Set to %d", d.heatmapPercent))
}

func (d debugger) draw(c drawConfig) {
	if d.drawHexGrid {
		drawHexGrid(c, *d.debuggerConfig)
	}

	if d.drawObstacleLUT {
		d.obstacleLUTUnwalkable.draw(c, *d.debuggerConfig, d.heatmapPercent)
		d.obstacleLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent)

		// TODO: Handle drawing obstacle LUT pixel when cursor if out of bounds
		p := d.oLUT[d.pixel].NearestUnwalkablesAvg
		p.draw(c, *d.debuggerConfig)
	}

	if d.drawNodeLUT {
		d.nodeLUTDistance.draw(c, *d.debuggerConfig, d.heatmapPercent)
	}

	if d.unitPathfindHeatmapComplete {
		d.unitPathfind.draw(c, *d.debuggerConfig, d.heatmapAlphaMax)
	}

	for _, n := range d.nodes {
		n.drawTargets(c, d.drawBlueTargets, d.drawRedTargets)
	}

	for _, u := range d.units {
		u.drawDebug(c, *d.debuggerConfig)
	}

	if d.drawVisibilityGrids {
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

func (p PixelPosition) draw(c drawConfig, dC debuggerConfig) {
	vector.DrawFilledRect(
		c.screen,
		float32(p.X)*c.windowMultiplier,
		float32(p.Y)*c.windowMultiplier,
		1*c.windowMultiplier,
		1*c.windowMultiplier,
		c.highlightColor,
		false)

	drawUnfilledBoxInt(
		c,
		c.highlightColor,
		p.X-dC.pixelSpotlightOffset,
		p.Y-dC.pixelSpotlightOffset,
		p.X+dC.pixelSpotlightOffset,
		p.Y+dC.pixelSpotlightOffset)
}

func (heatmap heatmap) draw(c drawConfig, dC debuggerConfig, heatmapPercent int) {
	var maxV int
	for _, v := range heatmap {
		maxV = max(maxV, v)
	}
	trimmedV := maxV * heatmapPercent / MaxPercent
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
			3*c.windowMultiplier,
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

	for i := 0; i < c.xIndexMax; i++ {
		for j := 0; j < c.yIndexMax; j++ {
			if j%2 == 0 && i == c.xIndexMax-1 {
				continue
			}

			lineWidth := float32(dC.debugGridWidth) * c.windowMultiplier
			pos := GridPosition{i, j}

			n := pos.getPixelPosOutbounds(hpNorth, &c.hexConfig)
			nX := float32(n.X) * c.windowMultiplier
			nY := float32(n.Y) * c.windowMultiplier

			nw := pos.getPixelPosOutbounds(hpNorthwest, &c.hexConfig)
			nwX := float32(nw.X) * c.windowMultiplier
			nwY := float32(nw.Y) * c.windowMultiplier

			ne := pos.getPixelPosOutbounds(hpNortheast, &c.hexConfig)
			neX := float32(ne.X+1) * c.windowMultiplier
			neY := float32(ne.Y) * c.windowMultiplier

			sw := pos.getPixelPosOutbounds(hpSouthwest, &c.hexConfig)
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
				lineWidth,
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
				lineWidth,
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
				lineWidth,
				dC.debugGridColor,
				false)

		}
	}
}
