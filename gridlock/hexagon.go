// Represents a hexagon from the game grid. Primarily used for tracking game state data.

package main

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type hexagons map[gridPosition]*hexagon

type hexagon struct {
	gridPos             gridPosition
	linkedList          *unit
	unitCount           int
	sumUnitX            int
	sumUnitY            int
	avgUnit             pixelPosition
	blueUnitsColliding  int
	blueUnitsUnwalkable int
	blueUnitsAttacking  int
	redUnitsColliding   int
	redUnitsUnwalkable  int
	redUnitsAttacking   int
}

func setupHexagons(hc *hexConfig) map[gridPosition]*hexagon {
	hexagons := make(map[gridPosition]*hexagon)
	for i := range hc.xIndexMax {
		for j := range hc.yIndexMax {
			gridPos := gridPosition{i, j}
			hexagons[gridPos] = &hexagon{gridPos: gridPos}
		}
	}
	return hexagons
}

func (h *hexagon) refresh(hc *hexConfig) {
	h.sumUnitY = 0
	h.sumUnitX = 0
	h.blueUnitsUnwalkable = 0
	h.blueUnitsColliding = 0
	h.blueUnitsAttacking = 0
	h.redUnitsUnwalkable = 0
	h.redUnitsColliding = 0
	h.redUnitsAttacking = 0
	if pos, ok := h.gridPos.getPixelPos(center, hc); ok {
		h.avgUnit = pos
	}
}

func (h *hexagon) drawFlux(dc drawConfig) {
	vector.DrawFilledCircle(
		dc.screen,
		float32(h.avgUnit.x)*dc.windowMultiplier,
		float32(h.avgUnit.y)*dc.windowMultiplier,
		float32(dc.avgUnitPosRadius)*dc.windowMultiplier,
		dc.avgUnitPosColor,
		false,
	)
}
