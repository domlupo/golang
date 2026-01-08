// Different types of polygons (a convex shape).

package main

import (
	"fmt"
	"math"
	"slices"
)

// A convex shape.
type polygon struct {
	edges []line
	// not every edge normal is a separating axis. Parallel edges only require testing one axis.
	separatingAxes []line
	vertices       []Pixel
}

// An origin centered convex shape.
type originPolygon polygon

// Creates a new origin centered polygon. If possible, the passed lines will be shifted to create an
// origin centered polygon. Refer to newPolygon() for further details.
func newOriginPolygon(lines []line) originPolygon {
	poly := newPolygon(lines)

	var midPoint Pixel
	for _, p := range poly.vertices {
		midPoint.X += p.X
		midPoint.Y += p.Y
	}

	if !midPoint.equals(Pixel{}) && !(midPoint.X%len(poly.vertices) == 0 &&
		midPoint.Y%len(poly.vertices) == 0) {
		panic(fmt.Sprintf(
			"newOriginPolygon: polygon is not origin centered and cannot be shifted: %s", midPoint.log()))
	}

	xShift := midPoint.X / len(poly.vertices) / 2
	yShift := midPoint.Y / len(poly.vertices) / 2

	var shiftedVertices []Pixel
	for _, p := range poly.vertices {
		s := Pixel{p.X + xShift, p.Y + yShift}
		shiftedVertices = append(shiftedVertices, s)
	}

	var shiftedEdges []line
	for _, e := range poly.edges {
		s := e
		s.p1.X += xShift
		s.p1.Y += yShift
		s.p2.X += xShift
		s.p2.Y += yShift
		shiftedEdges = append(shiftedEdges, s)
	}

	return originPolygon{
		edges:    shiftedEdges,
		vertices: shiftedVertices,
	}
}

// Creates a new polygon with basic polygon verifications. Polygon verifications will not
// prevent creation for all invalid polygons, but will catch typical invalid polygons. I suspect it
// is very hard to verify all invalid polygon creation so it is not an efficient use of time.
// Specifically collinear lines are probably tricky.
func newPolygon(lines []line) polygon {
	if len(lines) != 4 {
		panic("newPolygon: only 4 sided polygons are implemented")
	}

	sorted := []line{lines[0]}
	current := lines[0].p1
	end := lines[0].p2
	vertices := []Pixel{}
	for {
		vertices = append(vertices, current)
		priorSortCount := len(sorted)
		for _, l := range lines {
			if slices.Contains(sorted, l) {
				continue
			}
			if l.p1 == current {
				current = l.p2
				sorted = append(sorted, l)
				break
			}
			if l.p2 == current {
				current = l.p1
				sorted = append(sorted, l)
				break
			}
		}

		if len(sorted) == priorSortCount || (len(sorted) == len(lines) && current != end) {
			panic("newPolygon: valid polygon connections cannot be created with provided lines")
		} else if len(sorted) == len(lines) {
			vertices = append(vertices, current)
			break
		}
	}

	axes := []line{}
	for _, e := range lines {
		axis := e.perpendicular()
		noParallel := !slices.ContainsFunc(axes, func(a line) bool {
			return axis.parallel(a)
		})
		if noParallel {
			axes = append(axes, axis)
		}
	}

	return polygon{
		edges:          lines,
		separatingAxes: axes,
		vertices:       vertices,
	}
}

// check if two polygons collide via separating axis theorem.
func (p1 polygon) collides(p2 polygon) bool {
	axes := []line{}
	axes = append(axes, p1.separatingAxes...)
	axes = append(axes, p2.separatingAxes...)

	for _, a := range axes {
		x := a.p2.X - a.p1.X
		y := a.p2.Y - a.p1.Y

		p1Min := math.MaxInt
		p1Max := math.MinInt
		for _, v := range p1.vertices {
			dot := x*v.X + y*v.Y
			p1Min = min(p1Min, dot)
			p1Max = max(p1Max, dot)
		}

		p2Min := math.MaxInt
		p2Max := math.MinInt
		for _, v := range p2.vertices {
			dot := x*v.X + y*v.Y
			p2Min = min(p2Min, dot)
			p2Max = max(p2Max, dot)
		}

		overlaps := max(p1Min, p2Min) <= min(p1Max, p2Max)
		if !overlaps {
			// early exit since a separating line can be drawn between the two polygons
			return false
		}
	}

	return true
}
