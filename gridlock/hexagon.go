// A low level grid componenet. Spatially partitions the map for efficient computation and local data analytics.

package main

import (
	"fmt"
	"log/slog"
)

// All hexagons on map
type hexagons map[GridPosition]*hexagon

// A hexagon which spatially partitions the map. Stores dynamic data bucketed by layer.
type hexagon struct {
	gridPos  GridPosition
	land     hexLayer
	lowAir   hexLayer
	highAir  hexLayer
	selected bool
}

// A hexagonal spatial partition for a specific layer.
type hexLayer struct {
	linkedList    *unit
	unitCount     int
	unitCountRed  int
	unitCountBlue int
}

func (h hexagon) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("GridPosition", h.gridPos.log()))
	a = append(a, nest("land", h.land.log()))
	a = append(a, nest("lowAir", h.lowAir.log()))
	a = append(a, nest("highAir", h.highAir.log()))
	return a
}

func (hl hexLayer) log() []slog.Attr {
	var a = []slog.Attr{}
	if hl.linkedList == nil {
		a = append(a, slog.String("linkedlist", "none"))
	} else {
		a = append(a, nest("linkedList", hl.linkedList.PixelPosition.log()))
	}
	a = append(a, slog.Int("unitCount", hl.unitCount))
	a = append(a, slog.Int("unitCountRed", hl.unitCountRed))
	a = append(a, slog.Int("unitCountBlue", hl.unitCountBlue))
	return a
}

func (h *hexagon) get(l layer) *hexLayer {
	switch l {
	case land:
		return &h.land
	case lowAir:
		return &h.lowAir
	case highAir:
		return &h.highAir
	default:
		panic(fmt.Sprintf("get: unknown layer %d", l))
	}
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
func (hs hexagons) maxUnitDifferential(player player, l layer) GridPosition {
	var p GridPosition
	m := 0
	for _, h := range hs {
		hl := h.get(l)
		var d int
		switch player {
		case redPlayer:
			d = hl.unitCountRed - hl.unitCountBlue
		case bluePlayer:
			d = hl.unitCountBlue - hl.unitCountRed
		}
		m = max(m, d)
		if d == m {
			p = h.gridPos
		}
	}
	return p
}

func (h *hexagon) updatePointers(u *unit) {
	hl := h.get(u.layer)

	if hl.linkedList == nil {
		hl.linkedList = u
		hl.linkedList.prior = u
		hl.linkedList.next = u
	}
	u.next = hl.linkedList
	u.prior = hl.linkedList.prior
	hl.linkedList.prior.next = u
	hl.linkedList.prior = u
}

func (h *hexagon) add(u unit) {
	hl := h.get(u.layer)
	hl.unitCount++
	switch u.player {
	case redPlayer:
		hl.unitCountRed++
	case bluePlayer:
		hl.unitCountBlue++
	default:
		panic(fmt.Sprintf("cannot add player %d to hexagon", u.player))
	}
}

func (h *hexagon) subtract(u unit) {
	hl := h.get(u.layer)
	hl.unitCount--
	switch u.player {
	case redPlayer:
		hl.unitCountRed--
	case bluePlayer:
		hl.unitCountBlue--
	default:
		panic(fmt.Sprintf("cannot subtract player %d to hexagon", u.player))
	}

	if u.next.id == u.id && u.prior.id == u.id {
		hl.linkedList = nil
	} else if u.next.id == u.id || u.prior.id == u.id {
		panic("unit is in an invalid state")
	} else {
		hl.linkedList = u.next
	}

	u.next.prior = u.prior
	u.prior.next = u.next
}
