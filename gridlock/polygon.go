// Defines different types of polygons (a convex shape).

package main

import (
	"fmt"
	"slices"
)

// A convex shape.
type polygon struct {
	edges    []line
	vertices []Pixel
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

	return polygon{
		edges:    lines,
		vertices: vertices,
	}
}
