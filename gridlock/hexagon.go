// A low level grid componenet. Spatially partitions the map for efficient computation and local data analytics.

package main

import "fmt"

// All hexagons on map
type hexagons map[GridPosition]*hexagon

// A hexagon which spatially partitions the map and stores dynamic state.
type hexagon struct {
	gridPos       GridPosition
	linkedList    *unit
	unitCount     int
	unitCountRed  int
	unitCountBlue int
}

func setupHexagons(g grid) map[GridPosition]*hexagon {
	hexagons := make(map[GridPosition]*hexagon)
	for i := range g.xIndexMax {
		for j := range g.yIndexMax {
			if j%2 == 0 && i == g.xIndexMax-1 {
				continue
			}

			gridPos := GridPosition{i, j}
			hexagons[gridPos] = &hexagon{gridPos: gridPos}
		}
	}
	return hexagons
}

func (h *hexagon) updatePointers(u *unit) {
	if h.linkedList == nil {
		h.linkedList = u
		h.linkedList.prior = u
		h.linkedList.next = u
	}
	u.next = h.linkedList
	u.prior = h.linkedList.prior
	h.linkedList.prior.next = u
	h.linkedList.prior = u
}

func (h *hexagon) add(p player) {
	h.unitCount++
	switch p {
	case redPlayer:
		h.unitCountRed++
	case bluePlayer:
		h.unitCountBlue++
	default:
		panic(fmt.Sprintf("cannot add player %d to hexagon", p))
	}
}

func (h *hexagon) subtract(p player) {
	h.unitCount--
	switch p {
	case redPlayer:
		h.unitCountRed--
	case bluePlayer:
		h.unitCountBlue--
	default:
		panic(fmt.Sprintf("cannot subtract player %d to hexagon", p))
	}
}
