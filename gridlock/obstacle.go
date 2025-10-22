// Obstacles which block movement on a specific grid position.

package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// All obstacles on map
type obstacles map[gridPosition]*obstacle

// All points on map from every obstacle visibility grid
type obstaclePoints map[PixelPosition]*obstaclePoint

// A grid level component which blocks unit movement
type obstacle struct {
	gridPos        gridPosition
	kind           obstacleKind
	visibilityGrid *[]obstaclePoint
}

// A point on an obstacle visibility grid
type obstaclePoint struct {
	PixelPosition
	source           gridPosition
	visibilityGridId uuid.UUID
	cw               *obstaclePoint
	cwDistance       int
	ccw              *obstaclePoint
	ccwDistance      int
}

// look up table for pixel path data
type pixelPath map[PixelPosition]Path

// path data
type Path struct {
	Walkable                    bool `json:"walkable"`
	UnwalkableChebyshevDistance int  `json:"unwalkableChebyshevDistance"`
	UnwalkableSmoothSum         int  `json:"unwalkableSmoothSum"`
}

type obstacleKind int

const (
	none obstacleKind = iota
	lowMountain
	highMountain
)

func (op obstaclePoint) isZero() bool {
	return op.PixelPosition.equals(PixelPosition{0, 0})
}

func (os obstacles) insert(k gridPosition, v *obstacle) {
	_, ok := os[k]
	if ok {
		panic(fmt.Sprintf("grid position {%d, %d} already contains obstacle: ", k.x, k.y))
	}
	os[k] = v
}

func (os obstacles) setup(hc *hexConfig, pixelPathDistanceMax int) pixelPath {
	pixelPath := os.setupPixelPath(hc)
	groups := os.getGroups(hc)
	for _, group := range groups {
		group.setupVisibilityGrid(pixelPath, hc)
	}
	return pixelPath
}

// CPU intensize because it calculates the optimal distance from every walkable pixel to the nearest
// unwalkable pixel. It can be optimized further. The main optimization is dynamic programming. The
// unwalkable pixels are found first. Then the pixels 1 pixel away are calculated next. This
// continues until a pixel path distance limit is reached. The pixel path distance limit is
// not needed for player experience. It helps developer productivity by keeping live reloads fast.
func (os obstacles) setupPixelPath(hc *hexConfig) pixelPath {
	distanceLimit := max(hc.windowWidth, hc.windowHeight)
	pixelPath := make(map[PixelPosition]Path, hc.windowWidth*hc.windowHeight)

	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			p := PixelPosition{X: x, Y: y}
			gridPos, ok := p.getGridPos(hc)
			if !ok {
				panic("Could not setup pixel path")
			}
			o := os[gridPos]
			if o != nil && (o.kind == lowMountain || o.kind == highMountain) {
				pixelPath[p] = Path{Walkable: false}
			} else {
				pixelPath[p] = Path{Walkable: true}
			}
		}
	}

	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			p := PixelPosition{X: x, Y: y}
			path := pixelPath[p]
			if !path.Walkable {
				path.UnwalkableChebyshevDistance = 0
			} else {
				path.UnwalkableChebyshevDistance = distanceLimit
			}
			pixelPath[p] = path
		}
	}

	logger.log(infoLevel, setupLogic, "setup pixel path", slog.Int("distance limit", distanceLimit))

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
				path := pixelPath[p]

				if !path.Walkable || (path.Walkable && path.UnwalkableChebyshevDistance != distanceLimit) {
					count += 1
				}

				if pixelPath[p].UnwalkableChebyshevDistance != distance {
					continue
				}

				adjacent := p.allAdjacent()

				for _, a := range adjacent {
					path := pixelPath[a]
					path.UnwalkableChebyshevDistance = min(distance+1, path.UnwalkableChebyshevDistance)
					pixelPath[a] = path
				}
			}
		}

		percent := float64(count) / float64(hc.windowWidth*hc.windowHeight) * 100
		logger.log(
			infoLevel,
			setupLogic,
			fmt.Sprintf("setup pixel path. %.0f%% total pixels calculated. %d pixels of distance %d", percent, count-priorCount, distance))

		distance++
	}

	percent := float64(count) / float64(hc.windowWidth*hc.windowHeight) * 100
	logger.log(infoLevel, setupLogic, fmt.Sprintf("setup pixel path. %.0f%% total pixels not skipped with a distance limit of %d", percent, distanceLimit))

	logger.log(infoLevel, setupLogic, "start calculating weighted smooth for pixel path distance")
	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			p := PixelPosition{X: x, Y: y}
			adjacent := p.allAdjacent()

			sum := pixelPath[p].UnwalkableChebyshevDistance
			for _, a := range adjacent {
				sum += pixelPath[a].UnwalkableChebyshevDistance

				otherAdjacent := a.allAdjacent()
				for _, oa := range otherAdjacent {
					sum += pixelPath[oa].UnwalkableChebyshevDistance
				}
			}
			path := pixelPath[p]
			path.UnwalkableSmoothSum = sum
			pixelPath[p] = path
		}
	}
	logger.log(infoLevel, setupLogic, "finished calculating weighted smooth for pixel path distance")

	return pixelPath
}

func (os obstacles) getGroups(hc *hexConfig) []obstacles {
	needsGrouped := make(map[gridPosition]bool)
	for k := range os {
		needsGrouped[k] = true
	}

	var groups []obstacles

	for k := range os {
		if !needsGrouped[k] {
			continue
		}

		group := make(obstacles)
		q := []gridPosition{k}
		for len(q) != 0 {
			current := q[0]
			q = q[1:]
			group[current] = os[current]

			ne, neOk := current.adjacent(northeastSide, hc)
			if neOk && needsGrouped[ne] {
				needsGrouped[ne] = false
				q = append(q, ne)
			}
			e, eOk := current.adjacent(eastSide, hc)
			if eOk && needsGrouped[e] {
				needsGrouped[e] = false
				q = append(q, e)
			}
			se, seOk := current.adjacent(southeastSide, hc)
			if seOk && needsGrouped[se] {
				needsGrouped[se] = false
				q = append(q, se)
			}
			sw, swOk := current.adjacent(southwestSide, hc)
			if swOk && needsGrouped[sw] {
				needsGrouped[sw] = false
				q = append(q, sw)
			}
			w, wOk := current.adjacent(westSide, hc)
			if wOk && needsGrouped[w] {
				needsGrouped[w] = false
				q = append(q, w)
			}
			nw, nwOk := current.adjacent(northwestSide, hc)
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
func (pos PixelPosition) lineOfSight(otherPos PixelPosition, path pixelPath, radius int) (PixelPosition, bool) {
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
		if !path[p].Walkable || path[p].UnwalkableChebyshevDistance < radius-1 {
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
func (start PixelPosition) obstaclePoint(end PixelPosition, path pixelPath, obstacles obstacles, hc *hexConfig) (obstaclePoint, bool) {
	obstaclePixelPos, los := start.lineOfSight(end, path, 0)

	if los {
		return obstaclePoint{}, los
	}

	nearest := math.MaxInt64
	var nearestIndex int
	gridPos, ok := obstaclePixelPos.getGridPos(hc)
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

func (os obstacles) setupVisibilityGrid(path pixelPath, hc *hexConfig) {
	points := make(obstaclePoints)
	points.createAllPoints(os, hc)
	points.removeDuplicatePoints()
	points.connectPoints(path, hc)

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

func (points obstaclePoints) createAllPoints(os obstacles, hc *hexConfig) {
	// TODO: don't hardcode offset values for infantry unit with radius of 3
	nOffset := PixelPosition{0, -4}
	neOffset := PixelPosition{5, -2}
	seOffset := PixelPosition{5, 2}
	sOffset := PixelPosition{0, 5}
	swOffset := PixelPosition{-4, 2}
	nwOffset := PixelPosition{-4, -2}

	for _, o := range os {
		nObstaclePoint := &obstaclePoint{source: o.gridPos}
		neObstaclePoint := &obstaclePoint{source: o.gridPos}
		seObstaclePoint := &obstaclePoint{source: o.gridPos}
		sObstaclePoint := &obstaclePoint{source: o.gridPos}
		swObstaclePoint := &obstaclePoint{source: o.gridPos}
		nwObstaclePoint := &obstaclePoint{source: o.gridPos}

		n := o.gridPos.getPixelPosOutbounds(northPoint, hc)
		ne := o.gridPos.getPixelPosOutbounds(northeastPoint, hc)
		se := o.gridPos.getPixelPosOutbounds(southeastPoint, hc)
		s := o.gridPos.getPixelPosOutbounds(southPoint, hc)
		sw := o.gridPos.getPixelPosOutbounds(southwestPoint, hc)
		nw := o.gridPos.getPixelPosOutbounds(northwestPoint, hc)

		if nPoint, ok := n.add(nOffset, hc); ok {
			nObstaclePoint.PixelPosition = nPoint
			points[nPoint] = nObstaclePoint
		}
		if neTarget, ok := ne.add(neOffset, hc); ok {
			neObstaclePoint.PixelPosition = neTarget
			points[neTarget] = neObstaclePoint
		}
		if seTarget, ok := se.add(seOffset, hc); ok {
			seObstaclePoint.PixelPosition = seTarget
			points[seTarget] = seObstaclePoint
		}
		if sTarget, ok := s.add(sOffset, hc); ok {
			sObstaclePoint.PixelPosition = sTarget
			points[sTarget] = sObstaclePoint
		}
		if swTarget, ok := sw.add(swOffset, hc); ok {
			swObstaclePoint.PixelPosition = swTarget
			points[swTarget] = swObstaclePoint
		}
		if nwTarget, ok := nw.add(nwOffset, hc); ok {
			nwObstaclePoint.PixelPosition = nwTarget
			points[nwTarget] = nwObstaclePoint
		}
	}
}

func (points obstaclePoints) removeDuplicatePoints() {
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
				if op.cw != nil {
					op.cw.ccw = nil
				}
				if op.ccw != nil {
					op.ccw.cw = nil
				}
				if otherOp.cw != nil {
					otherOp.cw.ccw = nil
				}
				if otherOp.ccw != nil {
					otherOp.ccw.cw = nil
				}
				delete(points, PixelPosition{X: op.X, Y: op.Y})
				delete(points, PixelPosition{X: otherOp.X, Y: otherOp.Y})
			}
		}
	}
}

func (points obstaclePoints) connectPoints(path pixelPath, hc *hexConfig) {
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

			if _, lineOfSight := current.lineOfSight(op.PixelPosition, path, 0); !lineOfSight {
				continue
			}

			if !foundFirstCwConnection &&
				current.source.equals(op.source) &&
				orientation(current.PixelPosition, current.source.getPixelPosOutbounds(center, hc), op.PixelPosition) == clockwise {
				foundFirstCwConnection = true
			}

			if !foundFirstCwConnection {
				continue
			}

			isAdjacent := current.source.isAdjacent(op.source, hc)
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

func (o obstacle) drawVisibilityGrid(dc drawConfig) {
	for _, op := range *o.visibilityGrid {
		vector.DrawFilledCircle(
			dc.screen,
			float32(op.X)*dc.windowMultiplier,
			float32(op.Y)*dc.windowMultiplier,
			3*dc.windowMultiplier,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			float32(op.X)*dc.windowMultiplier,
			float32(op.Y)*dc.windowMultiplier,
			float32(op.cw.X)*dc.windowMultiplier,
			float32(op.cw.Y)*dc.windowMultiplier,
			float32(dc.visibilityGridLineWidth)*dc.windowMultiplier,
			dc.highlightColor,
			false)
	}
}

func (pp pixelPath) draw(dc drawConfig) {
	unwalkable := make(heatmap)
	for pixelPos, path := range pp {
		switch path.Walkable {
		case true:
			unwalkable[pixelPos] = 0
		case false:
			unwalkable[pixelPos] = 1
		}
	}

	distance := make(heatmap)
	for pixelPos, path := range pp {
		distance[pixelPos] = path.UnwalkableSmoothSum
	}

	unwalkable.draw(dc)
	distance.draw(dc)
}
