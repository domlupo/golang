// A debugger which outputs text logs and data visualizations to inspect game state.

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// A debugger which outputs text logs and data visualizations to inspect game state.
type debugger struct {
	*state
	unitPathfind                heatmap
	unitPathfindHeatmapComplete bool
}

func newDebugger(state *state) debugger {
	return debugger{
		state:        state,
		unitPathfind: make(heatmap),
	}
}

func (d *debugger) debug(p PixelPosition) {
	d.logPixel(p)
	d.createUnitPathfindHeatmap()
}

func (d *debugger) logPixel(p PixelPosition) {
	attrs := []slog.Attr{nest("PixelPosition", p.log())}
	if gridPos, ok := p.getGridPos(d.hexConfig); ok {
		attrs = append(attrs, nest("GridPosition", gridPos.log()))
	} else {
		attrs = append(attrs, slog.String("GridPosition", "notFound"))
	}
	attrs = append(attrs, nest("PixelPath", d.pixelPath[p].log()))
	attrs = append(attrs, nest("NodeLookup", d.nodeLookup[p].log()))
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
	for pixelPos, path := range d.pixelPath {
		select {
		case <-ctx.Done():
			return
		default:
			pixelCount++
			if pixelCount%(d.windowWidth*d.windowHeight/100) == 0 {
				percentProgess := pixelCount / (d.windowWidth * d.windowHeight / 100)
				logger.log(infoLevel, devTool, "pathfind heatmap", slog.Int("pixel positions percent progress", percentProgess))
			}

			if path.Walkable {
				unit.updateTarget(globalTarget, position{layer: unit.layer, PixelPosition: pixelPos})
				distance := unit.pathfindGlobalTarget(d.obstacles, d.pixelPath, d.hexConfig)
				d.unitPathfind[pixelPos] = distance
			}
		}
	}

	d.unitPathfindHeatmapComplete = true
	logger.log(infoLevel, devTool, "done creating pathfinding heatmap")
}

func (d debugger) draw(
	dc drawConfig,
	drawDebugGrid bool,
	drawPixelPath bool,
	drawNodeLookup bool,
	drawVisibilityGrids bool,
	drawBlueTargets bool,
	drawRedTargets bool) {

	if drawDebugGrid {
		drawHexDebugGrid(dc)
	}

	if drawPixelPath {
		d.pixelPath.draw(dc)
	}

	if drawNodeLookup {
		d.nodeLookup.draw(dc)
	}

	if d.unitPathfindHeatmapComplete {
		d.unitPathfind.draw(dc)
	}

	for _, n := range d.nodes {
		n.drawTargets(dc, drawBlueTargets, drawRedTargets)
	}

	if drawVisibilityGrids {
		visited := make(map[uuid.UUID]struct{})
		for _, o := range d.obstacles {
			g := *o.visibilityGrid
			op := g[0]
			if _, ok := visited[op.visibilityGridId]; ok {
				continue
			}
			visited[op.visibilityGridId] = struct{}{}
			o.drawVisibilityGrid(dc)
		}
	}
}

// This is likely to be become apart of game and then it should be removed from debugger.go
func drawHexDebugGrid(dc drawConfig) {
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

	for i := 0; i < dc.xIndexMax; i++ {
		for j := 0; j < dc.yIndexMax; j++ {
			if j%2 == 0 && i == dc.xIndexMax-1 {
				continue
			}

			lineWidth := float32(dc.debugGridWidth) * dc.windowMultiplier
			pos := GridPosition{i, j}

			n := pos.getPixelPosOutbounds(hpNorth, &dc.hexConfig)
			nX := float32(n.X) * dc.windowMultiplier
			nY := float32(n.Y) * dc.windowMultiplier

			nw := pos.getPixelPosOutbounds(hpNorthwest, &dc.hexConfig)
			nwX := float32(nw.X) * dc.windowMultiplier
			nwY := float32(nw.Y) * dc.windowMultiplier

			ne := pos.getPixelPosOutbounds(hpNortheast, &dc.hexConfig)
			neX := float32(ne.X+1) * dc.windowMultiplier
			neY := float32(ne.Y) * dc.windowMultiplier

			sw := pos.getPixelPosOutbounds(hpSouthwest, &dc.hexConfig)
			swX := float32(sw.X) * dc.windowMultiplier
			swY := float32(sw.Y+1) * dc.windowMultiplier

			// draw west side
			//    / \
			//  *|*
			//  *|*
			vector.StrokeLine(
				dc.screen,
				nwX,
				nwY,
				swX,
				swY,
				lineWidth,
				dc.debugGridColor,
				false)

			// draw northwest
			//   */*\
			//   |
			//   |
			vector.StrokeLine(
				dc.screen,
				nX,
				nY,
				nwX,
				nwY,
				lineWidth,
				dc.debugGridColor,
				false)

			// draw northeast
			//    /*\*
			//   |
			//   |
			vector.StrokeLine(
				dc.screen,
				nX,
				nY,
				neX,
				neY,
				lineWidth,
				dc.debugGridColor,
				false)

		}
	}
}
