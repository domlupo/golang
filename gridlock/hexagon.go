// A low level grid componenet. Spatially partitions the map for efficient computation and local data analytics.

package main

import (
	"fmt"
	"log/slog"
	"sort"
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
	unitCountPOne int
	unitCountPTwo int
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
	a = append(a, slog.Int("unitCountPOne", hl.unitCountPOne))
	a = append(a, slog.Int("unitCountPTwo", hl.unitCountPTwo))
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

type unitDifferential struct {
	GridPosition
	diff int
}

// Returns all hexagon GridPositions sorted by largest unit count differential for the passed player.
// No guaranteed sort for multiple grid positions with equal differential.
func (hs hexagons) maxUnitDifferential(opponent player, l layer) []unitDifferential {
	var diffs = []unitDifferential{}
	for p := range hs {
		hl := hs[p].get(l)
		var d int
		switch opponent {
		case pOne:
			d = hl.unitCountPOne - hl.unitCountPTwo
		case pTwo:
			d = hl.unitCountPTwo - hl.unitCountPOne
		default:
			panic(fmt.Sprintf("maxUnitDifferential: unknown player %d", opponent))
		}
		diffs = append(diffs, unitDifferential{p, d})
	}

	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].diff > diffs[j].diff
	})

	return diffs
}

// Add the passed unit into the hexagon partition.
func (h *hexagon) add(u *unit) {
	hl := h.get(u.layer)
	if hl.linkedList == nil && hl.unitCount != 0 {
		panic(fmt.Sprintf("nil linked list unit count should be 0, not %d", hl.unitCount))
	}

	if hl.linkedList == nil {
		hl.linkedList = u
		hl.linkedList.prior = u
		hl.linkedList.next = u
	}
	u.next = hl.linkedList
	u.prior = hl.linkedList.prior
	hl.linkedList.prior.next = u
	hl.linkedList.prior = u

	hl.unitCount++
	switch u.player {
	case pOne:
		hl.unitCountPOne++
	case pTwo:
		hl.unitCountPTwo++
	default:
		panic(fmt.Sprintf("cannot add player %d to hexagon", u.player))
	}

}

// Remove the passed unit from the hexagon partition.
func (h *hexagon) subtract(u unit) {
	hl := h.get(u.layer)
	if u.next.id == u.id && u.prior.id == u.id {
		hl.linkedList = nil
	} else if u.next.id == u.id || u.prior.id == u.id {
		panic("unit is in an invalid state")
	} else {
		hl.linkedList = u.next
	}
	u.next.prior = u.prior
	u.prior.next = u.next

	hl.unitCount--
	switch u.player {
	case pOne:
		hl.unitCountPOne--
	case pTwo:
		hl.unitCountPTwo--
	default:
		panic(fmt.Sprintf("cannot subtract player %d to hexagon", u.player))
	}
}
