// Obstacles which block movement on a specific grid position.

package main

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type obstacleKind int

const (
	none obstacleKind = iota
	mountain
)

type obstacle struct {
	gridPos        gridPosition
	kind           obstacleKind
	visibilityGrid []pixelPosition
}

func (o obstacle) drawObstacle(dc drawContext) {
	drawHex(dc, o.gridPos, drawTriangleColor{red: 0.2, blue: 0.2, green: 0.2, alpha: 0.8}, false)
}

func (o obstacle) drawVisibilityGrid(dc drawContext) {
	for _, pos := range o.visibilityGrid {
		vector.DrawFilledCircle(
			dc.screen,
			float32(pos.x)*dc.windowMultiplier,
			float32(pos.y)*dc.windowMultiplier,
			3*dc.windowMultiplier,
			dc.highlightColor,
			false)
	}
}

type obstacles []obstacle

func (os obstacles) setupVisibilityGrid(hc *hexConfig) {
	positions := make(map[pixelPosition]struct{})

	// TODO: don't hardcode offset values for infantry unit with radius of 3
	nOffset := pixelPosition{0, -4}
	neOffset := pixelPosition{5, -2}
	seOffset := pixelPosition{5, 2}
	sOffset := pixelPosition{0, 5}
	swOffset := pixelPosition{-4, 2}
	nwOffset := pixelPosition{-4, -2}

	for _, o := range os {
		n := o.gridPos.getPixelPosOutbounds(northPoint, hc)
		ne := o.gridPos.getPixelPosOutbounds(northeastPoint, hc)
		se := o.gridPos.getPixelPosOutbounds(southeastPoint, hc)
		s := o.gridPos.getPixelPosOutbounds(southPoint, hc)
		sw := o.gridPos.getPixelPosOutbounds(southwestPoint, hc)
		nw := o.gridPos.getPixelPosOutbounds(northwestPoint, hc)

		if nOffset, ok := n.add(nOffset, hc); ok {
			positions[nOffset] = struct{}{}
		}
		if neOffset, ok := ne.add(neOffset, hc); ok {
			positions[neOffset] = struct{}{}
		}
		if seOffset, ok := se.add(seOffset, hc); ok {
			positions[seOffset] = struct{}{}
		}
		if sOffset, ok := s.add(sOffset, hc); ok {
			positions[sOffset] = struct{}{}
		}
		if swOffset, ok := sw.add(swOffset, hc); ok {
			positions[swOffset] = struct{}{}
		}
		if nwOffset, ok := nw.add(nwOffset, hc); ok {
			positions[nwOffset] = struct{}{}
		}
	}

	for pos := range positions {
		for otherPos := range positions {
			if pos.x == otherPos.x && pos.y == otherPos.y {
				continue
			}

			// remove unneeded positions. Very hacky approach is to remove
			// any positions that collides with a radius of 5.
			// TODO: Don't hard code 5 and don't hack
			xDelta := pos.x - otherPos.x
			yDelta := pos.y - otherPos.y
			rSum := 5 + 5
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				delete(positions, pos)
				delete(positions, otherPos)
			}
		}
	}

	var grid []pixelPosition
	for pos := range positions {
		grid = append(grid, pos)
	}

	for i := range os {
		os[i].visibilityGrid = grid
	}

}
