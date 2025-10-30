// A low level grid componenet. Spatially partitions the map for efficient computation and local data analytics.

package main

import (
	"fmt"
	"log/slog"
)

// All hexagons on map
type hexagons map[GridPosition]*hexagon

// A hexagon which spatially partitions the map and stores dynamic state.
type hexagon struct {
	gridPos       GridPosition
	linkedList    *unit
	unitCount     int
	unitCountRed  int
	unitCountBlue int
	selected      bool
}

func (h hexagon) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("GridPosition", h.gridPos.log()))
	if h.linkedList == nil {
		a = append(a, slog.String("linkedlist", "none"))
	} else {
		a = append(a, nest("linkedList", h.linkedList.PixelPosition.log()))
	}
	a = append(a, slog.Int("unitCount", h.unitCount))
	a = append(a, slog.Int("unitCountRed", h.unitCountRed))
	a = append(a, slog.Int("unitCountBlue", h.unitCountBlue))
	return a
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

// Returns the hexagon GridPosition which has largest unit count differential for the
// passed player. If multiple grid positions share a max unit count, a random one is selected.
func (hs hexagons) maxUnitDifferential(player player) GridPosition {
	var p GridPosition
	m := 0
	for _, h := range hs {
		var d int
		switch player {
		case redPlayer:
			d = h.unitCountRed - h.unitCountBlue
		case bluePlayer:
			d = h.unitCountBlue - h.unitCountRed
		}
		m = max(m, d)
		if d == m {
			p = h.gridPos
		}
	}
	return p
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
