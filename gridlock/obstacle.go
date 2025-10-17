// Obstacles which block movement on a specific grid position.

package main

import (
	"image/color"
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type obstacles map[gridPosition]*obstacle

type obstacle struct {
	gridPos gridPosition
	groupId uuid.UUID
	kind    obstacleKind
	// TODO: determine if I should remove this. It seems a bit hacky. groupId probably works better.
	visibilityGrid *[]obstacleTarget
}

type obstacleTarget struct {
	pixelPosition
	groupId     uuid.UUID
	cw          *obstacleTarget
	cwDistance  int
	ccw         *obstacleTarget
	ccwDistance int
}

type pixelPathing map[pixelPosition]pathing

type pathing struct {
	walkable bool
}

type obstacleKind int

const (
	none obstacleKind = iota
	mountain
)

func (os obstacles) setup(hc *hexConfig) pixelPathing {
	pixelPathing := os.setupPixelPathing(hc)
	groups := os.getGroups(hc)
	for _, group := range groups {
		group.setupVisibilityGrid(hc, pixelPathing)
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
			if o != nil && o.kind == mountain {
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

		id, err := uuid.NewV4()
		if err != nil {
			logger.log(errorLevel, noLogType, "getGroups", slog.String("failed to generate UUID", err.Error()))
			panic("cannot getGroups since uuid generation failed")
		}

		group := make(obstacles)
		queue := []gridPosition{k}
		for len(queue) != 0 {
			current := queue[0]
			queue = queue[1:]
			os[current].groupId = id
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

// returns the obstacle target nearest to the start position for the obstacle blocking the line of sight between start and end
func (start pixelPosition) obstacleTarget(end pixelPosition, pathing pixelPathing, obstacles obstacles, hc *hexConfig) (obstacleTarget, bool) {
	obstaclePixelPos, los := start.lineOfSight(end, pathing)

	if los {
		return obstacleTarget{}, los
	}

	nearest := math.MaxInt64
	var nearestIndex int
	gridPos, ok := obstaclePixelPos.getGridPos(hc)
	if !ok {
		// TODO: Fix, this occurs when a human player selects a node + unit, click and holds a pixel then drags off screen
		panic("expected obstacle grid position to exist")
	}
	visibilityGrid := *obstacles[gridPos].visibilityGrid
	for index, ot := range visibilityGrid {
		distance := start.chebyshevDistance(ot.pixelPosition)
		if distance < nearest {
			nearest = distance
			nearestIndex = index
		}
	}

	return visibilityGrid[nearestIndex], los
}

func (os obstacles) setupVisibilityGrid(hc *hexConfig, pathing pixelPathing) {
	targets := os.setupObstacleTargets(hc)

	for _, ot := range targets {
		for _, otherOt := range targets {
			if ot.x == otherOt.x && ot.y == otherOt.y {
				continue
			}

			// remove unneeded positions. Very hacky approach is to remove
			// any positions that collides with a radius of 5.
			// TODO: Don't hard code 5 and don't hack
			xDelta := ot.x - otherOt.x
			yDelta := ot.y - otherOt.y
			rSum := 5 + 5
			if xDelta*xDelta+yDelta*yDelta <= rSum*rSum {
				if ot.cw != nil {
					ot.cw.ccw = nil
				}
				if ot.ccw != nil {
					ot.ccw.cw = nil
				}
				if otherOt.cw != nil {
					otherOt.cw.ccw = nil
				}
				if otherOt.ccw != nil {
					otherOt.ccw.cw = nil
				}
				delete(targets, pixelPosition{x: ot.x, y: ot.y})
				delete(targets, pixelPosition{x: otherOt.x, y: otherOt.y})
			}
		}
	}

	var nilCWTargets []*obstacleTarget
	var nilCCWTargets []*obstacleTarget
	for _, ot := range targets {
		if ot.cw == nil {
			nilCWTargets = append(nilCWTargets, ot)
		}
		if ot.ccw == nil {
			nilCCWTargets = append(nilCCWTargets, ot)
		}
	}

	for _, cw := range nilCWTargets {
		shortestDistance := math.MaxInt64
		var cwFix *obstacleTarget
		for _, ccw := range nilCCWTargets {
			if cw.x == ccw.x && cw.y == ccw.y {
				continue
			}

			if cw.cw != nil || ccw.ccw != nil {
				continue
			}

			if cw.ccw == ccw || ccw.cw == cw {
				continue
			}

			distance := cw.chebyshevDistance(ccw.pixelPosition)
			_, lineOfSight := cw.lineOfSight(ccw.pixelPosition, pathing)
			if distance < shortestDistance && lineOfSight {
				shortestDistance = distance
				cwFix = ccw
			} else if distance == shortestDistance && lineOfSight {
				panic("Handle visibility grid setup when points are equal. This will probably be rare or impossible")
			}
		}
		cw.cw = cwFix
		cwFix.ccw = cw
	}

	// TOODO: check if there is a better way to do this
	var id uuid.UUID
	for _, v := range os {
		id = v.groupId
		break
	}

	var grid []obstacleTarget
	for _, ot := range targets {
		ot.cwDistance = ot.chebyshevDistance(ot.cw.pixelPosition)
		ot.ccwDistance = ot.chebyshevDistance(ot.ccw.pixelPosition)
		ot.groupId = id
		grid = append(grid, *ot)
	}
	for i := range os {
		os[i].visibilityGrid = &grid
	}
}

func (os obstacles) setupObstacleTargets(hc *hexConfig) map[pixelPosition]*obstacleTarget {
	targets := make(map[pixelPosition]*obstacleTarget)

	// TODO: don't hardcode offset values for infantry unit with radius of 3
	nOffset := pixelPosition{0, -4}
	neOffset := pixelPosition{5, -2}
	seOffset := pixelPosition{5, 2}
	sOffset := pixelPosition{0, 5}
	swOffset := pixelPosition{-4, 2}
	nwOffset := pixelPosition{-4, -2}

	for _, o := range os {
		nObstacleTarget := &obstacleTarget{}
		neObstacleTarget := &obstacleTarget{}
		seObstacleTarget := &obstacleTarget{}
		sObstacleTarget := &obstacleTarget{}
		swObstacleTarget := &obstacleTarget{}
		nwObstacleTarget := &obstacleTarget{}

		nObstacleTarget.cw = neObstacleTarget
		nObstacleTarget.ccw = nwObstacleTarget
		neObstacleTarget.cw = seObstacleTarget
		neObstacleTarget.ccw = nObstacleTarget
		seObstacleTarget.cw = sObstacleTarget
		seObstacleTarget.ccw = neObstacleTarget
		sObstacleTarget.cw = swObstacleTarget
		sObstacleTarget.ccw = seObstacleTarget
		swObstacleTarget.cw = nwObstacleTarget
		swObstacleTarget.ccw = sObstacleTarget
		nwObstacleTarget.cw = nObstacleTarget
		nwObstacleTarget.ccw = swObstacleTarget

		n := o.gridPos.getPixelPosOutbounds(northPoint, hc)
		ne := o.gridPos.getPixelPosOutbounds(northeastPoint, hc)
		se := o.gridPos.getPixelPosOutbounds(southeastPoint, hc)
		s := o.gridPos.getPixelPosOutbounds(southPoint, hc)
		sw := o.gridPos.getPixelPosOutbounds(southwestPoint, hc)
		nw := o.gridPos.getPixelPosOutbounds(northwestPoint, hc)

		if nTarget, ok := n.add(nOffset, hc); ok {
			nObstacleTarget.pixelPosition = nTarget
			targets[nTarget] = nObstacleTarget
		} else {
			nObstacleTarget.ccw.cw = nObstacleTarget.cw
			nObstacleTarget.cw.ccw = nObstacleTarget.ccw
		}
		if neTarget, ok := ne.add(neOffset, hc); ok {
			neObstacleTarget.pixelPosition = neTarget
			targets[neTarget] = neObstacleTarget
		} else {
			neObstacleTarget.ccw.cw = neObstacleTarget.cw
			neObstacleTarget.cw.ccw = neObstacleTarget.ccw
		}
		if seTarget, ok := se.add(seOffset, hc); ok {
			seObstacleTarget.pixelPosition = seTarget
			targets[seTarget] = seObstacleTarget
		} else {
			seObstacleTarget.ccw.cw = seObstacleTarget.cw
			seObstacleTarget.cw.ccw = seObstacleTarget.ccw
		}
		if sTarget, ok := s.add(sOffset, hc); ok {
			sObstacleTarget.pixelPosition = sTarget
			targets[sTarget] = sObstacleTarget
		} else {
			sObstacleTarget.ccw.cw = sObstacleTarget.cw
			sObstacleTarget.cw.ccw = sObstacleTarget.ccw
		}
		if swTarget, ok := sw.add(swOffset, hc); ok {
			swObstacleTarget.pixelPosition = swTarget
			targets[swTarget] = swObstacleTarget
		} else {
			swObstacleTarget.ccw.cw = swObstacleTarget.cw
			swObstacleTarget.cw.ccw = swObstacleTarget.ccw
		}
		if nwTarget, ok := nw.add(nwOffset, hc); ok {
			nwObstacleTarget.pixelPosition = nwTarget
			targets[nwTarget] = nwObstacleTarget
		} else {
			nwObstacleTarget.ccw.cw = nwObstacleTarget.cw
			nwObstacleTarget.cw.ccw = nwObstacleTarget.ccw
		}
	}
	return targets
}

func (o obstacle) drawObstacle(dc drawConfig) {
	drawHex(dc, o.gridPos, drawTriangleColor{red: 0.2, blue: 0.2, green: 0.2, alpha: 0.8}, false)
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
