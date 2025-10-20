// Represents a hexagon from the game grid. Primarily used for tracking game state data.

package main

type hexagons map[gridPosition]*hexagon

type hexagon struct {
	gridPos    gridPosition
	linkedList *unit
	unitCount  int
}

func setupHexagons(hc *hexConfig) map[gridPosition]*hexagon {
	hexagons := make(map[gridPosition]*hexagon)
	for i := range hc.xIndexMax {
		for j := range hc.yIndexMax {
			if j%2 == 0 && i == hc.xIndexMax-1 {
				continue
			}

			gridPos := gridPosition{i, j}
			hexagons[gridPos] = &hexagon{gridPos: gridPos}
		}
	}
	return hexagons
}
