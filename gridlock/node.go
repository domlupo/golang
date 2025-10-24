// A node is a game component which spawns units and is captured by opposing units.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// All nodes on map
type nodes map[GridPosition]*node

// Look up table for nearest nodes.
type nodeLookup map[PixelPosition]NearestNodes

// TODO: Don't export NearestNodes it can be private
type NearestNodes struct {
	ChebyshevDistance int            `json:"nearestNodeChebyshevDistance"`
	Nodes             []GridPosition `json:"nodes"`
}

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

func (nn NearestNodes) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("ChebyshevDistance", nn.ChebyshevDistance))
	positions := []slog.Attr{}
	for _, p := range nn.Nodes {
		positions = append(positions, p.log()...)
	}
	if len(positions) == 0 {
		a = append(a, slog.String("nodes", "none"))
	} else {
		a = append(a, nest("nodes", positions))
	}
	return a
}

// CPU intensive because it calculates the optimal distance from any pixel to the nearest node. It
// probably can be optimized further but it is low priority since it it cached. The main optimization
// is dynamic programming. The node pixels are processed first. Next, the pixels 1 pixel away are
// processed. This continues until the all pixels are procesed.
func (nodes nodes) setupLookup(oLUT obstacleLUT, hc *hexConfig) nodeLookup {
	distanceLimit := max(hc.windowWidth, hc.windowHeight)
	lookup := make(map[PixelPosition]NearestNodes, hc.windowWidth*hc.windowHeight)

	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			p := PixelPosition{X: x, Y: y}
			gridPos, ok := p.getGridPos(hc)
			if !ok {
				panic("nodes setup lookup expected grid position to exist")
			}
			n := nodes[gridPos]
			if n != nil && n.player != noPlayer {
				lookup[p] = NearestNodes{ChebyshevDistance: 0, Nodes: []GridPosition{n.gridPos}}
			} else {
				lookup[p] = NearestNodes{ChebyshevDistance: distanceLimit}
			}
		}
	}

	logger.log(infoLevel, setupLogic, "setup node lookup", slog.Int("distance limit", distanceLimit))

	distance := 0
	var priorCount int
	var count int
	for count != hc.windowWidth*hc.windowHeight {
		if distance == distanceLimit {
			break
		}
		priorCount = count
		count = 0

		for x := range hc.windowWidth {
			for y := range hc.windowHeight {
				p := PixelPosition{X: x, Y: y}
				oInfo := oLUT[p]
				nnO := lookup[p]

				if !oInfo.Walkable || (oInfo.Walkable && nnO.ChebyshevDistance != distanceLimit) {
					count += 1
				}

				if !oLUT[p].Walkable || lookup[p].ChebyshevDistance != distance {
					continue
				}

				adjacent := p.allAdjacent(hc)

				for _, a := range adjacent {
					if !oLUT[a].Walkable {
						continue
					}

					nn := lookup[a]
					nn.ChebyshevDistance = min(distance+1, nn.ChebyshevDistance)
					lookup[a] = nn
				}

				for _, a := range adjacent {
					if !oLUT[a].Walkable {
						continue
					}

					nn := lookup[a]
					if nn.ChebyshevDistance < distance {
						nodes := nnO.Nodes
						for _, n := range nn.Nodes {
							if !slices.Contains(nodes, n) {
								nodes = append(nodes, n)
							}
						}
						nnO.Nodes = nodes
						lookup[p] = nnO
					}
				}
			}
		}

		percent := float64(count) / float64(hc.windowWidth*hc.windowHeight) * 100
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("setup node lookup. %.1f%% total pixels calculated. %d pixels of distance %d", percent, count-priorCount, distance))

		distance++
	}

	percent := float64(count) / float64(hc.windowWidth*hc.windowHeight) * 100
	logger.log(infoLevel, setupLogic, fmt.Sprintf("setup node lookup. %.0f%% total pixels not skipped with a distance limit of %d", percent, distanceLimit))

	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			p := PixelPosition{X: x, Y: y}
			oInfo := oLUT[p]
			if !oInfo.Walkable {
				nn := lookup[p]
				nn.ChebyshevDistance = distance
				lookup[p] = nn
			}
		}
	}

	return lookup
}

func (n node) isZero() bool {
	// TODO: fix implementation
	return n.player == noPlayer && n.gridPos.X == 0 && n.gridPos.Y == 0 && n.target.layer == noLayer
}

func (n node) selected() bool {
	return n.singleSelect || n.multiSelect
}

func (nodes nodes) setupControl(hc *hexConfig, hexagons hexagons, oLUT obstacleLUT) {

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

			p, ok := current.getAnyPixelPos(hc)
			if !ok {
				panic("setup control, could not find any pixel position")
			} else if !oLUT[p].Walkable {
				continue
			}

			if ne, ok := current.adjacent(hsNortheast, hc); ok {
				q = append(q, ne)
			}
			if e, ok := current.adjacent(hsEast, hc); ok {
				q = append(q, e)
			}
			if se, ok := current.adjacent(hsSoutheast, hc); ok {
				q = append(q, se)
			}
			if sw, ok := current.adjacent(hsSouthwest, hc); ok {
				q = append(q, sw)
			}
			if w, ok := current.adjacent(hsWest, hc); ok {
				q = append(q, w)
			}
			if nw, ok := current.adjacent(hsNorthwest, hc); ok {
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

	c := n.gridPos.getPixelPosOutbounds(hpCenter, &dc.hexConfig)

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
		var center, ok = n.gridPos.getPixelPos(hpCenter, &dc.hexConfig)
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
