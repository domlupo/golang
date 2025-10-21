// A game component which spawns units and captured by opposing units
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

type nodes map[gridPosition]*node

type node struct {
	player        player
	gridPos       gridPosition
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

func newNode(p player, gridPos gridPosition, initialNodeHealth int, hexConfig *hexConfig, uc unitConfig) node {
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
	return n.player == noPlayer && n.gridPos.x == 0 && n.gridPos.y == 0 && n.target.layer == noLayer
}

func (n node) selected() bool {
	return n.singleSelect || n.multiSelect
}

func (nodes nodes) setupControl(hc *hexConfig, hexagons hexagons, obstacles obstacles) {
	visited := make(map[gridPosition]struct{})

	for _, h := range hexagons {
		queue := []gridPosition{h.gridPos}
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if _, ok := visited[current]; ok {
				continue
			}

			visited[current] = struct{}{}

			if _, obstacle := obstacles[current]; obstacle {
				continue
			}

			if ne, ok := current.adjacent(northeastSide, hc); ok {
				queue = append(queue, ne)
			}
			if e, ok := current.adjacent(eastSide, hc); ok {
				queue = append(queue, e)
			}
			if se, ok := current.adjacent(southeastSide, hc); ok {
				queue = append(queue, se)
			}
			if sw, ok := current.adjacent(southwestSide, hc); ok {
				queue = append(queue, sw)
			}
			if w, ok := current.adjacent(westSide, hc); ok {
				queue = append(queue, w)
			}
			if nw, ok := current.adjacent(northwestSide, hc); ok {
				queue = append(queue, nw)
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
				queue = []gridPosition{}
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
			float32(c.x-n.radius)*dc.windowMultiplier,
			float32(c.y)*dc.windowMultiplier,
			(float32(c.x-n.radius)+healthLineLength)*dc.windowMultiplier,
			float32(c.y)*dc.windowMultiplier,
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
			float32(center.x)*dc.windowMultiplier,
			float32(center.y)*dc.windowMultiplier,
			float32(n.target.x)*dc.windowMultiplier,
			float32(n.target.y)*dc.windowMultiplier,
			float32(dc.targetNodeLineWidth)*dc.windowMultiplier,
			nodeColor,
			false,
		)
	}
}
