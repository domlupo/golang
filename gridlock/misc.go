// Misc code which is used in several other files but not large enough to need its own file.

package main

type player int

const (
	noPlayer player = iota
	bluePlayer
	redPlayer
	neutralPlayer
)

type hexagons map[gridPosition]*hexagon

// A hexagon which spatially partitions the map and stores dynamic state.
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
