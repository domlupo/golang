// A node is a game component which spawns units and is captured by opposing units.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// All nodes on map
type nodes map[GridPosition]*node

// Lookup table for processed node information.
type nodeLUT map[Pixel]nodeInfo

// Processed node information.
type nodeInfo struct {
	ChebyshevDistance int            `json:"chebyshevDistance"`
	NearestNodes      []GridPosition `json:"nearestNodes"`
}

// Game component which spawns units and is captured by opposing units
type node struct {
	gridPos            GridPosition
	player             player
	control            []*hexagon
	singleSelect       bool
	multiSelect        bool
	maxHealth          int
	currentHealth      int
	capturedBy         player
	targets            []Pixel
	actions            []action
	unitConfig         unitConfig
	pOneUnitAliveCount int
	pTwoUnitAliveCount int
	selected           bool
}

func newNode(p player, gridPos GridPosition, initialNodeHealth int, uc unitConfig) node {
	n := node{
		player:        p,
		gridPos:       gridPos,
		control:       []*hexagon{},
		maxHealth:     initialNodeHealth,
		currentHealth: initialNodeHealth,
		capturedBy:    pNone,
		unitConfig:    uc,
	}
	return n
}

func (n node) log() []slog.Attr {
	var a []slog.Attr
	a = append(a, nest("gridPos", n.gridPos.log()))
	a = append(a, slog.String("player", n.capturedBy.string()))
	if len(n.control) == 0 {
		a = append(a, slog.String("control", "none"))
	} else {
		control := []slog.Attr{}
		for _, c := range n.control {
			control = append(control, c.gridPos.log()...)
		}
		a = append(a, nest("control", control))
	}
	a = append(a, slog.Bool("singleSelect", n.singleSelect))
	a = append(a, slog.Bool("multiSelect", n.multiSelect))
	a = append(a, slog.Int("maxHealth", n.maxHealth))
	a = append(a, slog.Int("currentHealth", n.currentHealth))
	a = append(a, slog.String("capturedBy", n.capturedBy.string()))
	if len(n.targets) == 0 {
		a = append(a, slog.String("targets", "none"))
	} else {
		targets := []slog.Attr{}
		for _, t := range n.targets {
			targets = append(targets, t.log()...)
		}
		a = append(a, nest("targets", targets))
	}
	if len(n.actions) == 0 {
		a = append(a, slog.String("actions", "none"))
	} else {
		actions := []slog.Attr{}
		for _, a := range n.actions {
			actions = append(actions, a.log()...)
		}
		a = append(a, nest("actions", actions))
	}
	a = append(a, slog.String("unitConfig", n.unitConfig.name.string()))
	a = append(a, slog.Int("blueUnitAliveCount", n.pOneUnitAliveCount))
	a = append(a, slog.Int("redUnitAliveCount", n.pTwoUnitAliveCount))
	return a
}

// Returns a unique and deterministic slice of bytes from the passed nodes.
func (nodes nodes) key() []byte {
	var keys keys
	for k := range nodes {
		keys = append(keys, k)
	}
	keys.sort()
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
	if len(info.NearestNodes) == 0 {
		a = append(a, slog.String("NearestNodes", "none"))
	} else {
		nearestNodes := []slog.Attr{}
		for _, p := range info.NearestNodes {
			nearestNodes = append(nearestNodes, p.log()...)
		}
		a = append(a, nest("NearestNodes", nearestNodes))
	}
	return a
}

// CPU intensive because it calculates the optimal distance from any pixel to the nearest node. It
// probably can be optimized further but it is low priority since it it cached. The main
// optimization is dynamic programming. The node pixels are processed first. Next, the pixels 1
// pixel away are processed. This continues until the all pixels are processed.
func (nodes nodes) setupNodeLUT(oLUT obstacleLUT, g grid) nodeLUT {
	start := time.Now()
	distanceLimit := g.chebyshevPixelLimit
	nodeLUT := make(map[Pixel]nodeInfo, g.basePixels)

	for x := range g.pixelWidth {
		for y := range g.pixelHeight {
			p := Pixel{X: x, Y: y}
			gridPos, ok := g.getGridPos(p)
			if !ok {
				panic("setupNodeLUT: expected grid position to exist")
			}
			n := nodes[gridPos]
			if n != nil && n.player != pNone {
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
	for count != g.basePixels {
		priorCount = count
		count = 0

		for x := range g.pixelWidth {
			for y := range g.pixelHeight {
				p := Pixel{X: x, Y: y}
				oInfo := oLUT[p]
				nInfo := nodeLUT[p]

				if !oInfo.Walkable || (oInfo.Walkable && nInfo.ChebyshevDistance != distanceLimit) {
					count += 1
				}

				if !oLUT[p].Walkable || nodeLUT[p].ChebyshevDistance != distance {
					continue
				}

				adjacent := g.allAdjacentPixels(p)

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

		percent := float64(count) / float64(g.basePixels) * maxPercent
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("setupNodeLUT: %.1f%% total pixels calculated. %d pixels of distance %d",
				percent, count-priorCount, distance))

		distance++
		if distance > max(g.pixelHeight, g.pixelWidth) {
			panic("setupNodeLUT: invalid setup, distance is greater than max allowed by window.")
		}
	}

	percent := float64(count) / float64(g.basePixels) * maxPercent
	logger.log(infoLevel, setupLogic, fmt.Sprintf(
		"setupNodeLUT: %.0f%% total pixels not skipped with a distance limit of %d",
		percent,
		distanceLimit))

	for x := range g.pixelWidth {
		for y := range g.pixelHeight {
			p := Pixel{X: x, Y: y}
			oInfo := oLUT[p]
			if !oInfo.Walkable {
				info := nodeLUT[p]
				info.ChebyshevDistance = distance
				nodeLUT[p] = info
			}
		}
	}

	logger.log(infoLevel, setupLogic,
		fmt.Sprintf("setupNodeLUT: required %d seconds", int(time.Since(start).Seconds())))

	if len(nodeLUT) != g.basePixels {
		logger.log(errorLevel, setupLogic,
			"setupNodeLUT: cannot setup node look up table since created LUT is an invalid size")
		panic("setupNodeLUT: game requires a node LUT to start")
	}

	return nodeLUT
}

func (n node) isZero() bool {
	// TODO: fix implementation
	return n.player == pNone && n.gridPos.X == 0 && n.gridPos.Y == 0
}

func (n node) selectActive() bool {
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

			p, ok := g.getCenterWithFallback(current)
			if !ok {
				panic(fmt.Sprintf("setup control, could not find any pixel for grid pos %s", current.log()))
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
			case pOne, pTwo, pNeutral:
				n.control = append(n.control, h)
				// TODO: Need to determine behavior for same depth conflicts.
				// Currently, code selects first encountered node which is not correct.
				q = []GridPosition{}
			case pNone:
				continue
			default:
				panic("cannot setup hex control. Unknown player")
			}
		}
	}
}

func (n *node) drawNode(dc drawConfig, g grid) {
	var drawTriangleColor drawTriangleColor
	switch n.player {
	case pOne:
		drawTriangleColor = dc.pOneNodeColor
	case pTwo:
		drawTriangleColor = dc.pTwoNodeColor
	case pNeutral:
		drawTriangleColor = dc.pNeutralNodeColor
	default:
		panic("unknown player cannot drawNode")
	}
	drawHex(dc, g, n.gridPos, drawTriangleColor, n.multiSelect)

	if n.currentHealth != n.maxHealth {
		c := g.getPixelOutbounds(n.gridPos, hpCenter)
		ne := g.getPixelOutbounds(n.gridPos, hpNortheast)
		nw := g.getPixelOutbounds(n.gridPos, hpNorthwest)
		width := ne.X - nw.X
		healthLineLength := (float32(n.currentHealth) / float32(n.maxHealth)) * float32(width)

		vector.StrokeLine(
			dc.screen,
			float32(c.X-(width/2))*g.pixelMultiplier,
			float32(c.Y)*g.pixelMultiplier,
			(float32(c.X-(width/2))+healthLineLength)*g.pixelMultiplier,
			float32(c.Y)*g.pixelMultiplier,
			float32(dc.healthLineWidth)*g.pixelMultiplier,
			dc.nodeHealthColor,
			false)
	}

	if n.singleSelect {
		for _, h := range n.control {
			if !n.gridPos.equals(h.gridPos) {
				drawHex(dc, g, h.gridPos, dc.nodeControlColor, false)
			}
		}

	}
}

func (n *node) drawTargets(dc drawConfig, g grid, debugNodeTargets bool) {
	var color color.Color
	switch n.player {
	case pOne:
		color = dc.pOneNodeTargetColor
	case pTwo:
		color = dc.pTwoNodeTargetColor
	case pNeutral:
		color = dc.pNeutralNodeTargetColor
	default:
		panic("unknown player for drawTargets")
	}

	if !n.isZero() && debugNodeTargets {
		var centerPixel, ok = g.getPixel(n.gridPos, hpCenter)
		if !ok {
			panic("could not get pixel pos for node draw targets")
		}
		for _, t := range n.targets {
			vector.StrokeLine(
				dc.screen,
				float32(centerPixel.X)*g.pixelMultiplier,
				float32(centerPixel.Y)*g.pixelMultiplier,
				float32(t.X)*g.pixelMultiplier,
				float32(t.Y)*g.pixelMultiplier,
				float32(dc.targetNodeLineWidth)*g.pixelMultiplier,
				color,
				false,
			)
		}
	}
}
