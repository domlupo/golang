// Misc code which is used in several other files but not large enough to need its own file.

package main

import "fmt"

const MaxPercent = 100
const MinPercent = 0

type player int

const (
	noPlayer player = iota
	bluePlayer
	redPlayer
	neutralPlayer
)

func (p player) String() string {
	switch p {
	case noPlayer:
		return "noPlayer"
	case bluePlayer:
		return "bluePlayer"
	case redPlayer:
		return "redPlayer"
	case neutralPlayer:
		return "neutralPlayer"
	default:
		panic(fmt.Sprintf("Cannot convert player to string: %d", p))
	}
}

type windowConfig struct {
	width  int
	height int
}

func newWindowConfig() windowConfig {
	wc := windowConfig{
		width:  1280,
		height: 720,
	}
	return wc
}

func (wc windowConfig) maxPixels() int {
	return wc.height * wc.width
}

func (wc windowConfig) inbounds(p PixelPosition) bool {
	return !(p.X < 0 || p.X > wc.width-1 || p.Y < 0 || p.Y > wc.height-1)
}

// All hexagons on map
type hexagons map[GridPosition]*hexagon

// A hexagon which spatially partitions the map and stores dynamic state.
type hexagon struct {
	gridPos    GridPosition
	linkedList *unit
	unitCount  int
}

func setupHexagons(hc *hexConfig) map[GridPosition]*hexagon {
	hexagons := make(map[GridPosition]*hexagon)
	for i := range hc.xIndexMax {
		for j := range hc.yIndexMax {
			if j%2 == 0 && i == hc.xIndexMax-1 {
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
