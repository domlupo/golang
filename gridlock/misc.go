// Misc code which is used in several other files but not large enough to need its own file.

package main

type player int

const (
	noPlayer player = iota
	bluePlayer
	redPlayer
	neutralPlayer
)

// All hexagons on map
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
