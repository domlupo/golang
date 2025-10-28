// A node is a game component which spawns units and is captured by opposing units.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"slices"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// All nodes on map
type nodes map[GridPosition]*node

// Lookup table for processed node information.
type nodeLUT map[PixelPosition]nodeInfo

// Processed node information.
type nodeInfo struct {
	ChebyshevDistance int            `json:"chebyshevDistance"`
	NearestNodes      []GridPosition `json:"nearestNodes"`
}

// Game component which spawns units and is captured by opposing units
type node struct {
	player        player
	gridPos       GridPosition
	control       []*hexagon
	singleSelect  bool
	multiSelect   bool
	maxHealth     int
	currentHealth int
	capturedBy    player
	target        position
	unitConfig    unitConfig
}

func newNode(p player, gridPos GridPosition, initialNodeHealth int, uc unitConfig) node {
	n := node{
		player:        p,
		gridPos:       gridPos,
		control:       []*hexagon{},
		maxHealth:     initialNodeHealth,
		currentHealth: initialNodeHealth,
		capturedBy:    noPlayer,
		unitConfig:    uc,
	}
	return n
}

// Returns a unique and determinite slice of bytes from the passed nodes.
func (nodes nodes) key() []byte {
	var keys []GridPosition
	for k := range nodes {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].X == keys[j].X {
			return keys[i].Y < keys[j].Y
		}
		return keys[i].X < keys[j].X
	})

	var b []byte
	for i := 0; i < len(keys); i++ {
		n := nodes[keys[i]]
		b = append(b, truncate(n.gridPos.X))
		b = append(b, truncate(n.gridPos.Y))
		b = append(b, truncate(int(n.player)))
	}
	return b
}

func (info nodeInfo) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("ChebyshevDistance", info.ChebyshevDistance))
	nearestNodes := []slog.Attr{}
	for _, p := range info.NearestNodes {
		nearestNodes = append(nearestNodes, p.log()...)
	}
	if len(nearestNodes) == 0 {
		a = append(a, slog.String("NearestNodes", "none"))
	} else {
		a = append(a, nest("nearestNodes", nearestNodes))
	}
	return a
}

// CPU intensive because it calculates the optimal distance from any pixel to the nearest node. It
// probably can be optimized further but it is low priority since it it cached. The main optimization
// is dynamic programming. The node pixels are processed first. Next, the pixels 1 pixel away are
// processed. This continues until the all pixels are procesed.
func (nodes nodes) setupNodeLUT(oLUT obstacleLUT, g grid) nodeLUT {
	start := time.Now()
	wc := newWindowConfig()
	distanceLimit := max(wc.width, wc.height)
	nodeLUT := make(map[PixelPosition]nodeInfo, wc.maxPixels())

	for x := range wc.width {
		for y := range wc.height {
			p := PixelPosition{X: x, Y: y}
			gridPos, ok := g.getGridPos(p)
			if !ok {
				panic("setupNodeLUT: expected grid position to exist")
			}
			n := nodes[gridPos]
			if n != nil && n.player != noPlayer {
				nodeLUT[p] = nodeInfo{ChebyshevDistance: 0, NearestNodes: []GridPosition{n.gridPos}}
			} else {
				nodeLUT[p] = nodeInfo{ChebyshevDistance: distanceLimit}
			}
		}
	}

	logger.log(infoLevel, setupLogic, "setupNodeLUT", slog.Int("distance limit", distanceLimit))

	distance := 0
	var priorCount int
	var count int
	for count != wc.maxPixels() {
		priorCount = count
		count = 0

		for x := range wc.width {
			for y := range wc.height {
				p := PixelPosition{X: x, Y: y}
				oInfo := oLUT[p]
				nInfo := nodeLUT[p]

				if !oInfo.Walkable || (oInfo.Walkable && nInfo.ChebyshevDistance != distanceLimit) {
					count += 1
				}

				if !oLUT[p].Walkable || nodeLUT[p].ChebyshevDistance != distance {
					continue
				}

				adjacent := p.allAdjacent()

				for _, a := range adjacent {
					if !oLUT[a].Walkable {
						continue
					}

					naInfo := nodeLUT[a]
					naInfo.ChebyshevDistance = min(distance+1, naInfo.ChebyshevDistance)
					nodeLUT[a] = naInfo
				}

				for _, a := range adjacent {
					if !oLUT[a].Walkable {
						continue
					}

					naInfo := nodeLUT[a]
					if naInfo.ChebyshevDistance < distance {
						nodes := nInfo.NearestNodes
						for _, n := range naInfo.NearestNodes {
							if !slices.Contains(nodes, n) {
								nodes = append(nodes, n)
							}
						}
						nInfo.NearestNodes = nodes
						nodeLUT[p] = nInfo
					}
				}
			}
		}

		percent := float64(count) / float64(wc.maxPixels()) * maxPercent
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("setupNodeLUT: %.1f%% total pixels calculated. %d pixels of distance %d", percent, count-priorCount, distance))

		distance++
	}

	percent := float64(count) / float64(wc.maxPixels()) * maxPercent
	logger.log(infoLevel, setupLogic, fmt.Sprintf("setupNodeLUT: %.0f%% total pixels not skipped with a distance limit of %d", percent, distanceLimit))

	for x := range wc.width {
		for y := range wc.height {
			p := PixelPosition{X: x, Y: y}
			oInfo := oLUT[p]
			if !oInfo.Walkable {
				info := nodeLUT[p]
				info.ChebyshevDistance = distance
				nodeLUT[p] = info
			}
		}
	}

	logger.log(infoLevel, setupLogic, fmt.Sprintf("setupNodeLUT: required %d seconds", int(time.Since(start).Seconds())))
	return nodeLUT
}

func (n node) isZero() bool {
	// TODO: fix implementation
	return n.player == noPlayer && n.gridPos.X == 0 && n.gridPos.Y == 0 && n.target.layer == noLayer
}

func (n node) selected() bool {
	return n.singleSelect || n.multiSelect
}

func (nodes nodes) setupControl(g grid, hexagons hexagons, oLUT obstacleLUT) {

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

			p, ok := g.getAnyPixelPos(current)
			if !ok {
				panic("setup control, could not find any pixel position")
			} else if !oLUT[p].Walkable {
				continue
			}

			if ne, ok := g.adjacent(current, hsNortheast); ok {
				q = append(q, ne)
			}
			if e, ok := g.adjacent(current, hsEast); ok {
				q = append(q, e)
			}
			if se, ok := g.adjacent(current, hsSoutheast); ok {
				q = append(q, se)
			}
			if sw, ok := g.adjacent(current, hsSouthwest); ok {
				q = append(q, sw)
			}
			if w, ok := g.adjacent(current, hsWest); ok {
				q = append(q, w)
			}
			if nw, ok := g.adjacent(current, hsNorthwest); ok {
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

	if n.currentHealth != n.maxHealth {
		c := dc.grid.getPixelPosOutbounds(n.gridPos, hpCenter)
		ne := dc.grid.getPixelPosOutbounds(n.gridPos, hpNortheast)
		nw := dc.grid.getPixelPosOutbounds(n.gridPos, hpNorthwest)
		width := ne.X - nw.X
		healthLineLength := (float32(n.currentHealth) / float32(n.maxHealth)) * float32(width)

		vector.StrokeLine(
			dc.screen,
			float32(c.X-(width/2))*dc.windowMultiplier,
			float32(c.Y)*dc.windowMultiplier,
			(float32(c.X-(width/2))+healthLineLength)*dc.windowMultiplier,
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
		var center, ok = dc.grid.getPixelPos(n.gridPos, hpCenter)
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
