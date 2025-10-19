// Obstacles which block movement on a specific grid position.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type obstacles map[gridPosition]*obstacle

type obstaclePoints map[pixelPosition]*obstaclePoint

type obstacle struct {
	gridPos        gridPosition
	kind           obstacleKind
	visibilityGrid *[]obstaclePoint
}

type obstaclePoint struct {
	pixelPosition
	source           gridPosition
	visibilityGridId uuid.UUID
	cw               *obstaclePoint
	cwDistance       int
	ccw              *obstaclePoint
	ccwDistance      int
}

type pixelPathing map[pixelPosition]pathing

type pathing struct {
	walkable bool
}

type obstacleKind int

const (
	none obstacleKind = iota
	lowMountain
	highMountain
)

func (op obstaclePoint) isZero() bool {
	return op.pixelPosition.equals(pixelPosition{0, 0})
}

func (os obstacles) insert(k gridPosition, v *obstacle) {
	_, ok := os[k]
	if ok {
		panic(fmt.Sprintf("grid position {%d, %d} already contains obstacle: ", k.x, k.y))
	}
	os[k] = v
}

func (os obstacles) setup(hc *hexConfig) pixelPathing {
	pixelPathing := os.setupPixelPathing(hc)
	groups := os.getGroups(hc)
	for _, group := range groups {
		group.setupVisibilityGrid(pixelPathing, hc)
	}
	return pixelPathing
}

func (os obstacles) setupPixelPathing(hc *hexConfig) pixelPathing {
	pixelPathing := make(map[pixelPosition]pathing, hc.windowWidth*hc.windowHeight)

	for x := range hc.windowWidth {
		for y := range hc.windowHeight {
			pixelPos := pixelPosition{x: x, y: y}
			gridPos, ok := pixelPos.getGridPos(hc)
			if !ok {
				panic("Could not setup pixel path")
			}
			o := os[gridPos]
			walkable := true
			if o != nil && (o.kind == lowMountain || o.kind == highMountain) {
				walkable = false
			}
			pixelPathing[pixelPos] = pathing{walkable: walkable}
		}
	}

	return pixelPathing
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
		queue := []gridPosition{k}
		for len(queue) != 0 {
			current := queue[0]
			queue = queue[1:]
			group[current] = os[current]

			ne, neOk := current.adjacent(northeastSide, hc)
			if neOk && needsGrouped[ne] {
				needsGrouped[ne] = false
				queue = append(queue, ne)
			}
			e, eOk := current.adjacent(eastSide, hc)
			if eOk && needsGrouped[e] {
				needsGrouped[e] = false
				queue = append(queue, e)
			}
			se, seOk := current.adjacent(southeastSide, hc)
			if seOk && needsGrouped[se] {
				needsGrouped[se] = false
				queue = append(queue, se)
			}
			sw, swOk := current.adjacent(southwestSide, hc)
			if swOk && needsGrouped[sw] {
				needsGrouped[sw] = false
				queue = append(queue, sw)
			}
			w, wOk := current.adjacent(westSide, hc)
			if wOk && needsGrouped[w] {
				needsGrouped[w] = false
				queue = append(queue, w)
			}
			nw, nwOk := current.adjacent(northwestSide, hc)
			if nwOk && needsGrouped[nw] {
				needsGrouped[nw] = false
				queue = append(queue, nw)
			}
		}
		groups = append(groups, group)
	}
	return groups
}

// Returns true if a point to point line of sight exists using Bresenham's algorithm
func (pos pixelPosition) lineOfSight(otherPos pixelPosition, pathing pixelPathing) (pixelPosition, bool) {
	x0 := pos.x
	y0 := pos.y
	x1 := otherPos.x
	y1 := otherPos.y

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
		pathPos := pixelPosition{x: x0, y: y0}
		if !pathing[pathPos].walkable {
			return pathPos, false
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

	return pixelPosition{}, true
}

// returns the obstacle point nearest to the line of sight block from start to end
func (start pixelPosition) obstaclePoint(end pixelPosition, pathing pixelPathing, obstacles obstacles, hc *hexConfig) (obstaclePoint, bool) {
	obstaclePixelPos, los := start.lineOfSight(end, pathing)

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
		distance := obstaclePixelPos.chebyshevDistance(op.pixelPosition)
		if distance < nearest {
			nearest = distance
			nearestIndex = index
		}
	}

	return visibilityGrid[nearestIndex], los
}

func (os obstacles) setupVisibilityGrid(pathing pixelPathing, hc *hexConfig) {
	points := make(obstaclePoints)
	points.createAllPoints(os, hc)
	points.removeDuplicatePoints()
	points.connectPoints(pathing, hc)

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
	nOffset := pixelPosition{0, -4}
	neOffset := pixelPosition{5, -2}
	seOffset := pixelPosition{5, 2}
	sOffset := pixelPosition{0, 5}
	swOffset := pixelPosition{-4, 2}
	nwOffset := pixelPosition{-4, -2}

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
			nObstaclePoint.pixelPosition = nPoint
			points[nPoint] = nObstaclePoint
		}
		if neTarget, ok := ne.add(neOffset, hc); ok {
			neObstaclePoint.pixelPosition = neTarget
			points[neTarget] = neObstaclePoint
		}
		if seTarget, ok := se.add(seOffset, hc); ok {
			seObstaclePoint.pixelPosition = seTarget
			points[seTarget] = seObstaclePoint
		}
		if sTarget, ok := s.add(sOffset, hc); ok {
			sObstaclePoint.pixelPosition = sTarget
			points[sTarget] = sObstaclePoint
		}
		if swTarget, ok := sw.add(swOffset, hc); ok {
			swObstaclePoint.pixelPosition = swTarget
			points[swTarget] = swObstaclePoint
		}
		if nwTarget, ok := nw.add(nwOffset, hc); ok {
			nwObstaclePoint.pixelPosition = nwTarget
			points[nwTarget] = nwObstaclePoint
		}
	}
}

func (points obstaclePoints) removeDuplicatePoints() {
	for _, op := range points {
		for _, otherOp := range points {
			if op.equals(otherOp.pixelPosition) {
				continue
			}

			// remove unneeded positions. Very hacky approach is to remove
			// any positions that collides with a radius of 5.
			// TODO: Don't hard code 5 and don't hack
			xDelta := op.x - otherOp.x
			yDelta := op.y - otherOp.y
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
				delete(points, pixelPosition{x: op.x, y: op.y})
				delete(points, pixelPosition{x: otherOp.x, y: otherOp.y})
			}
		}
	}
}

func (points obstaclePoints) connectPoints(pathing pixelPathing, hc *hexConfig) {
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
			if current.pixelPosition.equals(op.pixelPosition) {
				continue
			}

			if op.cw != nil {
				continue
			}

			if _, lineOfSight := current.lineOfSight(op.pixelPosition, pathing); !lineOfSight {
				continue
			}

			if !foundFirstCwConnection &&
				current.source.equals(op.source) &&
				orientation(current.pixelPosition, current.source.getPixelPosOutbounds(center, hc), op.pixelPosition) == clockwise {
				foundFirstCwConnection = true
			}

			if !foundFirstCwConnection {
				continue
			}

			isAdjacent := current.source.isAdjacent(op.source, hc)
			distance := current.chebyshevDistance(op.pixelPosition)
			switch {
			case isAdjacent:
				cwConnection = op
				shortestDistance = 0
			case distance < shortestDistance:
				cwConnection = op
				shortestDistance = distance
			case distance == shortestDistance &&
				orientation(current.pixelPosition, op.pixelPosition, cwConnection.pixelPosition) == clockwise:
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
	distance := op.chebyshevDistance(other.pixelPosition)
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
	for _, ot := range *o.visibilityGrid {
		vector.DrawFilledCircle(
			dc.screen,
			float32(ot.x)*dc.windowMultiplier,
			float32(ot.y)*dc.windowMultiplier,
			3*dc.windowMultiplier,
			dc.highlightColor,
			false)
		vector.StrokeLine(
			dc.screen,
			float32(ot.x)*dc.windowMultiplier,
			float32(ot.y)*dc.windowMultiplier,
			float32(ot.cw.x)*dc.windowMultiplier,
			float32(ot.cw.y)*dc.windowMultiplier,
			float32(dc.visibilityGridLineWidth)*dc.windowMultiplier,
			dc.highlightColor,
			false)
	}
}

func (pp pixelPathing) draw(dc drawConfig) {
	for pixelPos, pathing := range pp {
		if !pathing.walkable {
			vector.DrawFilledRect(
				dc.screen,
				float32(pixelPos.x)*dc.windowMultiplier,
				float32(pixelPos.y)*dc.windowMultiplier,
				1*dc.windowMultiplier,
				1*dc.windowMultiplier,
				color.Black,
				false)
		}
	}
}
