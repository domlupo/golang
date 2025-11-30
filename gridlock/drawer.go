// Responsible for rendering pixels and images to the screen.

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type drawer struct {
	*state
	debugger *debugger
}

func newDrawer(s *state, d *debugger) drawer {
	drawer := drawer{
		state:    s,
		debugger: d,
	}
	return drawer
}

func (d drawer) draw(screen *ebiten.Image, dc drawConfig) {
	screen.Fill(dc.backgroundColor)
	d.drawHexGrid(dc)
	for _, n := range d.nodes {
		n.drawNode(dc, d.grid)
	}
	for _, o := range d.obstacles {
		o.drawObstacle(dc, d.grid)
	}
	for _, u := range d.landUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range d.chaffUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range d.bombGravityUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range d.ordanceUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range d.aircraftUnits {
		u.draw(dc, u.selected)
	}
	for _, u := range d.ballisticUnits {
		u.draw(dc, u.selected)
	}
	d.debugger.draw(dc)
}

func (d drawer) drawHexGrid(c drawConfig) {
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

	for i := 0; i < d.grid.xIndexMax; i++ {
		for j := 0; j < d.grid.yIndexMax; j++ {
			if j%2 == 0 && i == d.grid.xIndexMax-1 {
				continue
			}

			p := GridPosition{i, j}

			n := d.grid.getPixelOutbounds(p, hpNorth)
			nX := float32(n.X) * c.windowMultiplier
			nY := float32(n.Y) * c.windowMultiplier

			nw := d.grid.getPixelOutbounds(p, hpNorthwest)
			nwX := float32(nw.X) * c.windowMultiplier
			nwY := float32(nw.Y) * c.windowMultiplier

			ne := d.grid.getPixelOutbounds(p, hpNortheast)
			neX := float32(ne.X+1) * c.windowMultiplier
			neY := float32(ne.Y) * c.windowMultiplier

			sw := d.grid.getPixelOutbounds(p, hpSouthwest)
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
				float32(c.hexGridLineWidth)*c.windowMultiplier,
				c.hexGridColor,
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
				float32(c.hexGridLineWidth)*c.windowMultiplier,
				c.hexGridColor,
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
				float32(c.hexGridLineWidth)*c.windowMultiplier,
				c.hexGridColor,
				false)
		}
	}
}
