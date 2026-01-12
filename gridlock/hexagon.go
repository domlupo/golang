// A low level grid component. Spatially partitions the map for efficient computation and local
// data analytics.

package main

import (
	"fmt"
	"log/slog"
	"slices"
)

// All hexagons on grid.
type hexagons map[GridPosition]*hexagon

// A hexagon which spatially partitions the map. Stores dynamic data bucketed by layer.
type hexagon struct {
	gridPos     GridPosition
	land        hexLayer
	ordnance    hexLayer
	aircraft    hexLayer
	ballistic   hexLayer
	bombGravity hexLayer
	chaff       hexLayer
	selected    bool
}

// A hexagonal spatial partition for a specific layer.
type hexLayer struct {
	linkedList    *unit
	unitCount     int
	unitCountPOne int
	unitCountPTwo int
}

// Typical log method.
func (h hexagon) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("GridPosition", h.gridPos.log()))
	a = append(a, nest("land", h.land.log()))
	a = append(a, nest("ordnance", h.ordnance.log()))
	a = append(a, nest("aircraft", h.aircraft.log()))
	a = append(a, nest("ballistic", h.ballistic.log()))
	a = append(a, nest("bombGravity", h.bombGravity.log()))
	a = append(a, nest("chaff", h.chaff.log()))
	return a
}

// Typical log method.
func (hl hexLayer) log() []slog.Attr {
	var a = []slog.Attr{}
	if hl.linkedList == nil {
		a = append(a, slog.String("linkedList", "none"))
	} else {
		a = append(a, nest("linkedList", hl.linkedList.Pos.log()))
	}
	a = append(a, slog.Int("unitCount", hl.unitCount))
	a = append(a, slog.Int("unitCountPOne", hl.unitCountPOne))
	a = append(a, slog.Int("unitCountPTwo", hl.unitCountPTwo))
	return a
}

func (h *hexagon) get(l layer) *hexLayer {
	switch l {
	case lLand:
		return &h.land
	case lOrdnance:
		return &h.ordnance
	case lAircraft:
		return &h.aircraft
	case lBallistic:
		return &h.ballistic
	case lBombGravity:
		return &h.bombGravity
	case lChaff:
		return &h.chaff
	default:
		panic(fmt.Sprintf("hexagon.get: unknown layer %s", l.log()))
	}
}

func newHexagons(g grid) map[GridPosition]*hexagon {
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

// Returns all hexagon GridPositions sorted by largest unit count differential for the passed
// player. No guaranteed sort for multiple grid positions with equal differential.
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
			panic(fmt.Sprintf("hexagons.maxUnitDifferential: unknown player %d", opponent))
		}
		diffs = append(diffs, unitDifferential{p, d})
	}

	slices.SortFunc(diffs, func(i, j unitDifferential) int {
		switch {
		case i.diff < j.diff:
			return 1
		case i.diff > j.diff:
			return -1
		default:
			return 1
		}
	})

	return diffs
}

// Add the passed unit into the hexagon partition.
func (h *hexagon) add(u *unit) {
	hl := h.get(u.L)
	if hl.linkedList == nil && hl.unitCount != 0 {
		panic(fmt.Sprintf(
			"hexagon.add: nil linked list unit count should be 0, not %d. \n hexagon: %s \n unit: %s",
			hl.unitCount,
			h.log(),
			u.log()))
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
		panic(fmt.Sprintf("hexagon.add: cannot add player %d to hexagon", u.player))
	}

}

// Remove the passed unit from the hexagon partition.
func (h *hexagon) subtract(u unit) {
	hl := h.get(u.L)
	if u.next.id == u.id && u.prior.id == u.id {
		hl.linkedList = nil
	} else if u.next.id == u.id || u.prior.id == u.id {
		panic("hexagon.subtract: unit is in an invalid state")
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
		panic(fmt.Sprintf("hexagon.subtract: cannot subtract player %d from hexagon", u.player))
	}
}
