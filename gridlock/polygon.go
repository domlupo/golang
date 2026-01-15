// Different types of polygons (a convex shape).

package main

import (
	"log/slog"
	"math"
	"slices"
)

// A convex shape.
type polygon struct {
	edges []line
	// not every edge normal is a separating axis. Parallel edges only require testing one axis.
	separatingAxes []line
	vertices       []Pos
}

// An origin centered convex shape. Useful for rotations.
type originPolygon struct {
	polygon
	conservativeRadius float64
}

// Creates an origin centered polygon by shifting passed lines. Refer to newPolygon() for details.
func newOriginPolygon(lines []line) originPolygon {
	poly := newPolygon(lines)

	var center Pos
	for _, p := range poly.vertices {
		center.X += p.X
		center.Y += p.Y
	}

	xShift := center.X / float64(len(poly.vertices)) / 2.0
	yShift := center.Y / float64(len(poly.vertices)) / 2.0

	var shiftedVertices []Pos
	for _, p := range poly.vertices {
		s := Pos{X: p.X + xShift, Y: p.Y + yShift}
		shiftedVertices = append(shiftedVertices, s)
	}

	var conservativeRadius float64
	for _, p := range poly.vertices {
		radius := math.Sqrt(p.X*p.X + p.Y*p.Y)
		conservativeRadius = max(conservativeRadius, radius)
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
		polygon: polygon{
			edges:    shiftedEdges,
			vertices: shiftedVertices,
		},
		conservativeRadius: conservativeRadius,
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
	vertices := []Pos{}
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

	center := Pos{}
	for _, v := range vertices {
		center.X += v.X
		center.Y += v.Y
	}
	center.X = center.X / float64(len(vertices))
	center.Y = center.Y / float64(len(vertices))

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

// Check if two polygons collide via separating axis theorem.
func (p1 polygon) collides(p2 polygon) bool {
	axes := []line{}
	axes = append(axes, p1.separatingAxes...)
	axes = append(axes, p2.separatingAxes...)

	for _, a := range axes {
		x := a.p2.X - a.p1.X
		y := a.p2.Y - a.p1.Y

		p1Min := math.MaxFloat64
		p1Max := -math.MaxFloat64
		for _, v := range p1.vertices {
			dot := x*v.X + y*v.Y
			p1Min = min(p1Min, dot)
			p1Max = max(p1Max, dot)
		}

		p2Min := math.MaxFloat64
		p2Max := -math.MaxFloat64
		for _, v := range p2.vertices {
			dot := x*v.X + y*v.Y
			p2Min = min(p2Min, dot)
			p2Max = max(p2Max, dot)
		}

		overlaps := max(p1Min, p2Min) <= min(p1Max, p2Max)
		if !overlaps {
			// Early exit since a separating line can be drawn between the two polygons.
			return false
		}
	}

	return true
}

// Returns true if the polygon collides with the circle. The algorithm first finds the nearest point
// on every polygon edge to the circle exterior. Then the algorithm returns true any edge's nearest
// point is touching or within the circle. A ray cast is used to handle circle within polygon edge
// case.
func (p polygon) collidesCircle(d dynamic) bool {
	rayCast := line{p1: d.Pos, p2: Pos{X: math.MaxFloat64, Y: d.Pos.Y}}
	intersectionCount := 0

	for _, e := range p.edges {
		if rayCast.intersects(e) && !rayCast.parallel(e) {
			intersectionCount++
		}

		x := math.Abs(e.p2.X - e.p1.X)
		y := math.Abs(e.p2.Y - e.p1.Y)
		cosineSimularity := (x*d.Pos.X + y*d.Pos.Y) / ((x * x) + (y * y))

		nearestPos := Pos{}
		switch {
		case cosineSimularity <= 0:
			nearestPos = e.p1
		case cosineSimularity >= 0:
			nearestPos = e.p2
		default:
			nearestPos = Pos{
				X: cosineSimularity * x,
				Y: cosineSimularity,
			}
		}

		xDelta := nearestPos.X - d.Pos.X
		yDelta := nearestPos.Y - d.Pos.Y
		if (xDelta*xDelta)+(yDelta*yDelta) <= float64(*d.circleRadius)*float64(*d.circleRadius) {
			return true
		}
	}

	return intersectionCount%2 == 1
}

// physics which may not be be mutated.
type static struct {
	offset        Pos
	originPolygon *originPolygon
	circleRadius  *int
}

// Typical log method.
func (s static) log() []slog.Attr {
	a := []slog.Attr{}
	if s.originPolygon == nil {
		a = append(a, slog.String("originPolygon", "none"))
	} else {
		// TODO: add polygon.log() method
		a = append(a, slog.String("originPolygon", "some"))
	}
	if s.circleRadius == nil {
		a = append(a, slog.String("circleRadius", "none"))
	} else {
		a = append(a, slog.Int("circleRadius", *s.circleRadius))
	}
	return a
}

// allows partitioner.go to treat hexagons as circles. This is simple and correct but conservative
// any may need optimized in the future once units use hexagon bodies.
func (s static) conservativeRadius() int {
	switch {
	case s.circleRadius != nil:
		return *s.circleRadius
	case s.originPolygon != nil:
		return int(math.Ceil(s.originPolygon.conservativeRadius))
	default:
		panic("static.conservativeRadius: no static pointer set")
	}
}

// physics which may be mutated.
// TODO: decide if I should add tests? e.g. move a polygon.
// TODO: decide if rotate and move methods be on dynamic struct.
type dynamic struct {
	Pos
	// TODO: add Pos here and remove center within circle and polygon.
	polygon      *polygon
	circleRadius *int
}

func (s static) dynamic(p Pos) dynamic {
	d := dynamic{Pos: p}
	switch {
	case s.circleRadius != nil:
		d.circleRadius = s.circleRadius
	case s.originPolygon != nil:
		d.polygon = &s.originPolygon.polygon
	default:
		panic("static.dynamic")
	}
	return d
}

// Typical log method.
func (d dynamic) log() []slog.Attr {
	a := []slog.Attr{}
	if d.polygon == nil {
		a = append(a, slog.String("polygon", "none"))
	} else {
		// TODO: add polygon log method once debugging requires it or small time slot to implement.
		a = append(a, slog.String("polygon", "some"))
	}
	if d.circleRadius == nil {
		a = append(a, slog.String("circle", "none"))
	} else {
		// TODO: add circle log method once debugging requires it or small time slot to implement.
		a = append(a, slog.String("circle", "some"))
	}
	return a
}
