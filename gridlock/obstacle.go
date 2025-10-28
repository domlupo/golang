// Obstacles which block movement on a specific grid position.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"sort"
	"time"

	"github.com/gofrs/uuid/v5"
)

// All obstacles on map
type obstacles map[GridPosition]*obstacle

// All points on map from every obstacle visibility grid
type obstaclePoints map[PixelPosition]*obstaclePoint

// A grid level component which blocks unit movement
type obstacle struct {
	gridPos        GridPosition
	kind           obstacleKind
	visibilityGrid *[]obstaclePoint
}

// A point on an obstacle visibility grid
type obstaclePoint struct {
	PixelPosition
	source           GridPosition
	visibilityGridId uuid.UUID
	cw               *obstaclePoint
	cwDistance       int
	ccw              *obstaclePoint
	ccwDistance      int
	selected         bool
}

// Lookup table for processed obstacle information.
type obstacleLUT map[PixelPosition]obstacleInfo

// Processed obstacle information.
type obstacleInfo struct {
	Walkable                    bool          `json:"walkable"`
	NearestUnwalkablesAvg       PixelPosition `json:"nearestUnwalkablesAvg"`
	UnwalkableChebyshevDistance int           `json:"unwalkableChebyshevDistance"`
	UnwalkableSmoothSum         int           `json:"unwalkableSmoothSum"`
}

type obstacleKind int

const (
	none obstacleKind = iota
	lowMountain
	highMountain
)

func (op obstaclePoint) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("PixelPosition", op.PixelPosition.log()))
	a = append(a, nest("source", op.source.log()))
	a = append(a, slog.String("visibilityGridId", op.visibilityGridId.String()))
	a = append(a, nest("cw", op.cw.PixelPosition.log()))
	a = append(a, slog.Int("cwDistance", op.cwDistance))
	a = append(a, nest("ccw", op.ccw.PixelPosition.log()))
	a = append(a, slog.Int("ccwDistance", op.ccwDistance))
	return a
}

func (info obstacleInfo) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Bool("Walkable", info.Walkable))
	a = append(a, nest("nearestUnwalkablesAvg", info.NearestUnwalkablesAvg.log()))
	a = append(a, slog.Int("UnwalkableChebyshevDistance", info.UnwalkableChebyshevDistance))
	a = append(a, slog.Int("UnwalkableSmoothSum", info.UnwalkableSmoothSum))
	return a
}

func (op obstaclePoint) isZero() bool {
	return op.PixelPosition.equals(PixelPosition{0, 0})
}

func (os obstacles) insert(k GridPosition, v *obstacle) {
	_, ok := os[k]
	if ok {
		panic(fmt.Sprintf("grid position {%d, %d} already contains obstacle: ", k.X, k.Y))
	}
	os[k] = v
}

// Returns a unique and determinite slice of bytes from the passed obstacles.
func (os obstacles) key() []byte {
	var keys []GridPosition
	for k := range os {
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
		o := os[keys[i]]
		b = append(b, truncate(o.gridPos.X))
		b = append(b, truncate(o.gridPos.Y))
		b = append(b, truncate(int(o.kind)))
	}
	return b
}

func (os obstacles) setup(g grid) obstacleLUT {
	fileHandler := newFileHandler()
	oLUT := fileHandler.readObstacleLUT(os, g)
	groups := os.getGroups(g)
	for _, group := range groups {
		group.setupVisibilityGrid(oLUT, g)
	}
	return oLUT
}

// CPU intensive for several reasons. One, it calculates the optimal distance from every walkable pixel
// to the nearest unwalkable pixel (dynamic programming is needed since it easily takes >30 mins without).
// Two, it calculates the nearestUnwalkablesAvg which is (o(n)^3) complexity. This has a worse case of
// (1280 * 720 * ~1440) aka (window width pixel * window height pixels * a pixel where the nearest obstacles
// are a the largest possible square around it. 1440 is (720 * 4 / 2). 4 sides but the hexagons sides would
// alternate cutting the value in half.
//
// Can be optimized further but it is low priority since LUT is cached and optimizatons are error prone.
// The main optimization currently used is dynamic programming. The unwalkable pixels are processed first. Next,
// the pixels 1 pixel away are processed. This continues until all the pixels are processed.
func (os obstacles) setupObstacleLUT(g grid) obstacleLUT {
	start := time.Now()
	wc := newWindowConfig()
	distanceLimit := max(wc.width, wc.height)
	nearestUnwalkables := make(map[PixelPosition][]PixelPosition)
	oLUT := make(map[PixelPosition]obstacleInfo, wc.width*wc.height)

	for x := range wc.width {
		for y := range wc.height {
			p := PixelPosition{X: x, Y: y}
			gridPos, ok := g.getGridPos(p)
			if !ok {
				panic("setupObstacleLUT: Could not setup obstacle lookup table")
			}
			o := os[gridPos]
			if o != nil && (o.kind == lowMountain || o.kind == highMountain) {
				oLUT[p] = obstacleInfo{Walkable: false}
				nearestUnwalkables[p] = []PixelPosition{p}
			} else {
				oLUT[p] = obstacleInfo{Walkable: true}
			}
		}
	}

	for x := range wc.width {
		for y := range wc.height {
			p := PixelPosition{X: x, Y: y}
			info := oLUT[p]
			if !info.Walkable {
				info.UnwalkableChebyshevDistance = 0
			} else {
				info.UnwalkableChebyshevDistance = distanceLimit
			}
			oLUT[p] = info
		}
	}

	logger.log(infoLevel, setupLogic, "setupObstacleLUT", slog.Int("distance limit", distanceLimit))

	distance := 0
	var priorCount int
	var count int
	for count != wc.maxPixels() {
		priorCount = count
		count = 0

		for x := range wc.width {
			for y := range wc.height {
				p := PixelPosition{X: x, Y: y}
				info := oLUT[p]

				if !info.Walkable || (info.Walkable && info.UnwalkableChebyshevDistance != distanceLimit) {
					count += 1
				}

				if info.UnwalkableChebyshevDistance != distance {
					continue
				}

				adjacent := p.allAdjacent()

				for _, a := range adjacent {
					aInfo := oLUT[a]
					aInfo.UnwalkableChebyshevDistance = min(distance+1, aInfo.UnwalkableChebyshevDistance)
					oLUT[a] = aInfo
				}

				for _, a := range adjacent {
					aInfo := oLUT[a]
					if aInfo.UnwalkableChebyshevDistance < distance {
						n := nearestUnwalkables[p]
						for _, unwalkable := range nearestUnwalkables[a] {
							if !slices.Contains(n, unwalkable) {
								n = append(n, unwalkable)
							}
						}
						nearestUnwalkables[p] = n
						oLUT[p] = info
					}
				}
			}
		}

		percent := float64(count) / float64(wc.maxPixels()) * maxPercent
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("setupObstacleLUT: %.1f%% total pixels calculated. %d pixels of distance %d", percent, count-priorCount, distance))

		// Currently, this calculation takes the longest e.g. (30s of 45s total).
		// I looked into removing floating point math. A few CPU profiles showed this is only
		// slightly less performant (< 5 seconds) than integer math. I suspect using integer math
		// may lead to very hard to deubg issues which I want to avoid. Another potential error prone
		// optmization that I am avoiding is summing the pixel position as they are added to the array.
		// Optimizing may be needed in the future to speed up setup, but LUT caching is sufficient for now.
		for x := range wc.width {
			for y := range wc.height {
				p := PixelPosition{X: x, Y: y}
				n := nearestUnwalkables[p]

				xSum := 0
				ySum := 0
				for _, pos := range n {
					xSum += pos.X
					ySum += pos.Y
				}
				xAvg := int(math.Round(float64(xSum) / float64(len(n))))
				yAvg := int(math.Round(float64(ySum) / float64(len(n))))

				info := oLUT[p]
				info.NearestUnwalkablesAvg = PixelPosition{X: xAvg, Y: yAvg}
				oLUT[p] = info
			}
		}

		distance++
	}

	percent := float64(count) / float64(wc.maxPixels()) * maxPercent
	logger.log(
		infoLevel,
		setupLogic,
		fmt.Sprintf(
			"setupObstacleLUT: %.0f%% total pixels not skipped with a distance limit of %d",
			percent,
			distanceLimit))

	logger.log(infoLevel, setupLogic, "setupObstacleLUT: start calculating weighted smooth distance")
	for x := range wc.width {
		for y := range wc.height {
			p := PixelPosition{X: x, Y: y}
			adjacent := p.allAdjacent()

			sum := oLUT[p].UnwalkableChebyshevDistance
			for _, a := range adjacent {
				sum += oLUT[a].UnwalkableChebyshevDistance

				otherAdjacent := a.allAdjacent()
				for _, oa := range otherAdjacent {
					sum += oLUT[oa].UnwalkableChebyshevDistance
				}
			}
			info := oLUT[p]
			info.UnwalkableSmoothSum = sum
			oLUT[p] = info
		}
	}

	logger.log(infoLevel, setupLogic, fmt.Sprintf("setupObstacleLUT: required %d seconds", int(time.Since(start).Seconds())))
	return oLUT
}

func (os obstacles) getGroups(g grid) []obstacles {
	needsGrouped := make(map[GridPosition]bool)
	for k := range os {
		needsGrouped[k] = true
	}

	var groups []obstacles

	for k := range os {
		if !needsGrouped[k] {
			continue
		}

		group := make(obstacles)
		q := []GridPosition{k}
		for len(q) != 0 {
			current := q[0]
			q = q[1:]
			group[current] = os[current]

			ne, neOk := g.adjacent(current, hsNortheast)
			if neOk && needsGrouped[ne] {
				needsGrouped[ne] = false
				q = append(q, ne)
			}
			e, eOk := g.adjacent(current, hsEast)
			if eOk && needsGrouped[e] {
				needsGrouped[e] = false
				q = append(q, e)
			}
			se, seOk := g.adjacent(current, hsSoutheast)
			if seOk && needsGrouped[se] {
				needsGrouped[se] = false
				q = append(q, se)
			}
			sw, swOk := g.adjacent(current, hsSouthwest)
			if swOk && needsGrouped[sw] {
				needsGrouped[sw] = false
				q = append(q, sw)
			}
			w, wOk := g.adjacent(current, hsWest)
			if wOk && needsGrouped[w] {
				needsGrouped[w] = false
				q = append(q, w)
			}
			nw, nwOk := g.adjacent(current, hsNorthwest)
			if nwOk && needsGrouped[nw] {
				needsGrouped[nw] = false
				q = append(q, nw)
			}
		}
		groups = append(groups, group)
	}
	return groups
}

// Returns true if a point to point line of sight exists using Bresenham's algorithm
func (pos PixelPosition) lineOfSight(otherPos PixelPosition, oLUT obstacleLUT, radius int) (PixelPosition, bool) {
	x0 := pos.X
	y0 := pos.Y
	x1 := otherPos.X
	y1 := otherPos.Y

	var sx int
	var sy int
	dx := x1 - x0
	dy := y1 - y0
	if dx < 0 {
		dx = dx * -1
	}
	if dy > 0 {
		dy = dy * -1
	}
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	bresenhamError := dx + dy

	for x0 != x1 || y0 != y1 {
		p := PixelPosition{X: x0, Y: y0}
		if !oLUT[p].Walkable || oLUT[p].UnwalkableChebyshevDistance < radius-1 {
			return p, false
		}

		doubleError := 2 * bresenhamError
		if doubleError >= dy && x0 != x1 {
			bresenhamError += dy
			x0 += sx
		}
		if doubleError <= dx && y0 != y1 {
			bresenhamError += dx
			y0 += sy
		}

	}

	return PixelPosition{}, true
}

// returns the obstacle point nearest to the line of sight block from start to end
func (start PixelPosition) obstaclePoint(end PixelPosition, oLUT obstacleLUT, obstacles obstacles, g grid) (obstaclePoint, bool) {
	obstaclePixelPos, los := start.lineOfSight(end, oLUT, 0)

	if los {
		return obstaclePoint{}, los
	}

	nearest := math.MaxInt64
	var nearestIndex int
	gridPos, ok := g.getGridPos(obstaclePixelPos)
	if !ok {
		// occurs when cursor is off screen
		return obstaclePoint{}, los
	}
	visibilityGrid := *obstacles[gridPos].visibilityGrid
	for index, op := range visibilityGrid {
		distance := obstaclePixelPos.chebyshevDistance(op.PixelPosition)
		if distance < nearest {
			nearest = distance
			nearestIndex = index
		}
	}

	return visibilityGrid[nearestIndex], los
}

func (os obstacles) setupVisibilityGrid(oLUT obstacleLUT, g grid) {
	points := make(obstaclePoints)
	points.createAllPoints(os, g)
	points.removeDuplicatePoints()
	points.connectPoints(oLUT, g)

	id, err := uuid.NewV4()
	if err != nil {
		logger.log(errorLevel, noLogType, "getGroups", slog.String("failed to generate UUID", err.Error()))
		panic("cannot getGroups since uuid generation failed")
	}

	var grid []obstaclePoint
	for _, op := range points {
		op.visibilityGridId = id
		grid = append(grid, *op)
	}
	for i := range os {
		os[i].visibilityGrid = &grid
	}
}

func (points obstaclePoints) createAllPoints(os obstacles, g grid) {
	// TODO: don't hardcode offset values for infantry unit with radius of 3
	nOffset := PixelPosition{0, -4}
	neOffset := PixelPosition{4, -2}
	seOffset := PixelPosition{4, 2}
	sOffset := PixelPosition{0, 4}
	swOffset := PixelPosition{-4, 2}
	nwOffset := PixelPosition{-4, -2}

	for _, o := range os {
		nObstaclePoint := &obstaclePoint{source: o.gridPos}
		neObstaclePoint := &obstaclePoint{source: o.gridPos}
		seObstaclePoint := &obstaclePoint{source: o.gridPos}
		sObstaclePoint := &obstaclePoint{source: o.gridPos}
		swObstaclePoint := &obstaclePoint{source: o.gridPos}
		nwObstaclePoint := &obstaclePoint{source: o.gridPos}

		n := g.getPixelPosOutbounds(o.gridPos, hpNorth)
		ne := g.getPixelPosOutbounds(o.gridPos, hpNortheast)
		se := g.getPixelPosOutbounds(o.gridPos, hpSoutheast)
		s := g.getPixelPosOutbounds(o.gridPos, hpSouth)
		sw := g.getPixelPosOutbounds(o.gridPos, hpSouthwest)
		nw := g.getPixelPosOutbounds(o.gridPos, hpNorthwest)

		if nPoint, ok := n.add(nOffset); ok {
			nObstaclePoint.PixelPosition = nPoint
			points[nPoint] = nObstaclePoint
		}
		if neTarget, ok := ne.add(neOffset); ok {
			neObstaclePoint.PixelPosition = neTarget
			points[neTarget] = neObstaclePoint
		}
		if seTarget, ok := se.add(seOffset); ok {
			seObstaclePoint.PixelPosition = seTarget
			points[seTarget] = seObstaclePoint
		}
		if sTarget, ok := s.add(sOffset); ok {
			sObstaclePoint.PixelPosition = sTarget
			points[sTarget] = sObstaclePoint
		}
		if swTarget, ok := sw.add(swOffset); ok {
			swObstaclePoint.PixelPosition = swTarget
			points[swTarget] = swObstaclePoint
		}
		if nwTarget, ok := nw.add(nwOffset); ok {
			nwObstaclePoint.PixelPosition = nwTarget
			points[nwTarget] = nwObstaclePoint
		}
	}
}

func (points obstaclePoints) removeDuplicatePoints() {
	duplicate := make(map[PixelPosition]struct{})
	for _, op := range points {
		for _, otherOp := range points {
			if op.equals(otherOp.PixelPosition) {
				continue
			}

			// remove unneeded positions. Very hacky approach is to remove
			// any positions that collides with a radius of 5.
			// TODO: Don't hard code 5 and don't hack
			xDelta := op.X - otherOp.X
			yDelta := op.Y - otherOp.Y
			rSum := 5 + 5
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				duplicate[PixelPosition{X: op.X, Y: op.Y}] = struct{}{}
				duplicate[PixelPosition{X: otherOp.X, Y: otherOp.Y}] = struct{}{}
			}
		}
	}

	for p := range duplicate {
		delete(points, p)
	}
}

func (points obstaclePoints) connectPoints(oLUT obstacleLUT, g grid) {
	var notConnected []*obstaclePoint
	for _, op := range points {
		notConnected = append(notConnected, op)
	}

	var foundFirstCwConnection bool
	startIndex := 0
	start := notConnected[startIndex]
	current := start
	for start.ccw == nil {
		shortestDistance := math.MaxInt64
		var cwConnection *obstaclePoint

		for _, op := range notConnected {
			if current.PixelPosition.equals(op.PixelPosition) {
				continue
			}

			if op.cw != nil {
				continue
			}

			if _, lineOfSight := current.lineOfSight(op.PixelPosition, oLUT, 0); !lineOfSight {
				continue
			}

			if !foundFirstCwConnection &&
				current.source.equals(op.source) &&
				orientation(current.PixelPosition, g.getPixelPosOutbounds(current.source, hpCenter), op.PixelPosition) == clockwise {
				foundFirstCwConnection = true
			}

			if !foundFirstCwConnection {
				continue
			}

			isAdjacent := g.isAdjacent(current.source, op.source)
			distance := current.chebyshevDistance(op.PixelPosition)
			switch {
			case isAdjacent:
				cwConnection = op
				shortestDistance = 0
			case distance < shortestDistance:
				cwConnection = op
				shortestDistance = distance
			case distance == shortestDistance &&
				orientation(current.PixelPosition, op.PixelPosition, cwConnection.PixelPosition) == clockwise:
				cwConnection = op
			}
		}

		switch {
		case cwConnection != nil:
			current.connectCwTo(cwConnection)
			current = current.cw
		case foundFirstCwConnection:
			current.connectCwTo(start)
		default:
			foundFirstCwConnection = false
			startIndex++
			start = notConnected[startIndex]
			current = start
		}
	}
}

func (op *obstaclePoint) connectCwTo(other *obstaclePoint) {
	distance := op.chebyshevDistance(other.PixelPosition)
	op.cw = other
	other.ccw = op
	op.cwDistance = distance
	other.ccwDistance = distance
}

func (o obstacle) drawObstacle(dc drawConfig) {
	switch o.kind {
	case lowMountain:
		drawHex(dc, o.gridPos, drawTriangleColor{red: 0.2, blue: 0.2, green: 0.2, alpha: 0.8}, false)
	case highMountain:
		drawHex(dc, o.gridPos, drawTriangleColor{red: 0.0, blue: 0.0, green: 0.0, alpha: 0.8}, false)
	}
}
