// A node is a game component which spawns units and is captured by opposing units.

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// All nodes on map
type nodes map[GridPosition]*node

// Game component which spawns units and is captured by opposing units
type node struct {
	player        player
	gridPos       GridPosition
	control       []*hexagon
	radius        int
	singleSelect  bool
	multiSelect   bool
	maxHealth     int
	currentHealth int
	capturedBy    player
	target        position
	unitConfig    unitConfig
}

func newNode(p player, gridPos GridPosition, initialNodeHealth int, hexConfig *hexConfig, uc unitConfig) node {
	n := node{
		player:  p,
		gridPos: gridPos,
		control: []*hexagon{},
		// HACK: "/ 9 * 2 is temporarily set until hexagon collision is implemented
		// This value was chosen since it is large but still fully within hexagon.
		radius:        (hexConfig.hexHeight + hexConfig.hexWidth) / 9 * 2,
		maxHealth:     initialNodeHealth,
		currentHealth: initialNodeHealth,
		capturedBy:    noPlayer,
		unitConfig:    uc,
	}
	return n
}

func (n node) isZero() bool {
	// TODO: fix implementation
	return n.player == noPlayer && n.gridPos.X == 0 && n.gridPos.Y == 0 && n.target.layer == noLayer
}

func (n node) selected() bool {
	return n.singleSelect || n.multiSelect
}

func (nodes nodes) setupControl(hc *hexConfig, hexagons hexagons, obstacles obstacles) {

	for _, h := range hexagons {
		q := []GridPosition{h.gridPos}
		visited := make(map[GridPosition]struct{})

		for len(q) > 0 {
			current := q[0]
			q = q[1:]

			if _, ok := visited[current]; ok {
				continue
			}

			visited[current] = struct{}{}

			if _, obstacle := obstacles[current]; obstacle {
				continue
			}

			if ne, ok := current.adjacent(northeastSide, hc); ok {
				q = append(q, ne)
			}
			if e, ok := current.adjacent(eastSide, hc); ok {
				q = append(q, e)
			}
			if se, ok := current.adjacent(southeastSide, hc); ok {
				q = append(q, se)
			}
			if sw, ok := current.adjacent(southwestSide, hc); ok {
				q = append(q, sw)
			}
			if w, ok := current.adjacent(westSide, hc); ok {
				q = append(q, w)
			}
			if nw, ok := current.adjacent(northwestSide, hc); ok {
				q = append(q, nw)
			}

			n, nOk := nodes[current]
			if !nOk {
				continue
			}

			switch n.player {
			case redPlayer, bluePlayer, neutralPlayer:
				n.control = append(n.control, h)
				// TODO: Need to determine behavior for same depth conflicts.
				// Currently, code selects first enountered node which is not correct.
				q = []GridPosition{}
			case noPlayer:
				continue
			default:
				panic("cannot setup hex control. Unknown player")
			}
		}
	}
}

func (n *node) drawNode(dc drawConfig) {
	var red float32
	var blue float32
	var green float32
	var alpha float32
	switch n.player {
	case bluePlayer:
		red = 0.5
		blue = 1
		green = 0.5
		alpha = 1
	case redPlayer:
		red = 1
		blue = 0.5
		green = 0.5
		alpha = 1
	case neutralPlayer:
		red = 0.5
		blue = 0.5
		green = 0.5
		alpha = 1
	default:
		panic("unknown player cannot drawNode")
	}

	drawHex(dc, n.gridPos, drawTriangleColor{red: red, blue: blue, green: green, alpha: alpha}, n.multiSelect)

	c := n.gridPos.getPixelPosOutbounds(center, &dc.hexConfig)

	if n.currentHealth != n.maxHealth {
		var healthLineLength = (float32(n.currentHealth) / float32(n.maxHealth)) * float32(n.radius*2)
		vector.StrokeLine(
			dc.screen,
			float32(c.X-n.radius)*dc.windowMultiplier,
			float32(c.Y)*dc.windowMultiplier,
			(float32(c.X-n.radius)+healthLineLength)*dc.windowMultiplier,
			float32(c.Y)*dc.windowMultiplier,
			float32(dc.healthLineWidth)*dc.windowMultiplier,
			dc.nodeHealthColor,
			false)
	}

	if n.singleSelect {
		for _, h := range n.control {
			if !n.gridPos.equals(h.gridPos) {
				drawHex(dc, h.gridPos, drawTriangleColor{red: 0.7, blue: 0.7, green: 0.7, alpha: alpha}, false)
			}
		}

	}
}

func (n *node) drawTargets(dc drawConfig, drawBlueTargets bool, drawRedTargets bool) {
	var nodeColor color.Color
	switch n.player {
	case bluePlayer:
		nodeColor = dc.bluePlayerNodeColor
	case redPlayer:
		nodeColor = dc.redPlayerNodeColor
	case neutralPlayer:
		nodeColor = dc.neutralNodeColor
	default:
		panic("unknown player for drawTargets")
	}

	if !n.target.isZero() && ((n.player == bluePlayer && drawBlueTargets) || (n.player == redPlayer && drawRedTargets)) {
		var center, ok = n.gridPos.getPixelPos(center, &dc.hexConfig)
		if !ok {
			panic("could not get pixel pos for node draw targets")
		}
		vector.StrokeLine(
			dc.screen,
			float32(center.X)*dc.windowMultiplier,
			float32(center.Y)*dc.windowMultiplier,
			float32(n.target.X)*dc.windowMultiplier,
			float32(n.target.Y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			nodeColor,
			false,
		)
	}
}
