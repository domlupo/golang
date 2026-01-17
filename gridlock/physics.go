// Different physical shapes used for movement, collision, etc.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

// A physical shape which may be mutated. The single shape is indicated by set pointer.
// Do not create directly, use a constructor.
type dynamic struct {
	Pos
	static       static
	polygon      *polygon
	circleRadius *int
}

// A physical shape which may not be be mutated. The single shape is indicated by set pointer.
// This struct is typically used for either calculating values to create a dynamic or an analysis
// which only considers the physical shape not the physical position.
// Do not create directly, use a constructor.
type static struct {
	offset        Pos
	originPolygon *originPolygon
	circleRadius  *int
}

// Create a dynamic circle. Typical constructor.
func newDynamicCircle(p Pos, r int) *dynamic {
	s := newStaticCircle(r, Pos{})
	d := dynamic{
		static:       s,
		Pos:          p,
		circleRadius: &r,
	}
	return &d
}

// Create a dynamic polygon. Sets the Pos from the passed lines center.
func newDynamicPolygon(lines []line) *dynamic {
	s := newStaticPolygon(lines, Pos{})
	polygon := newPolygon(lines)
	p := Pos{}
	for _, v := range polygon.vertices {
		p.X += v.X
		p.Y += v.Y
	}
	p.X = p.X / float64(len(polygon.vertices))
	p.Y = p.Y / float64(len(polygon.vertices))
	d := dynamic{
		static:  s,
		Pos:     p,
		polygon: &polygon,
	}
	return &d
}

// Create a mutable dynamic from a static which should not be mutated.
func (s static) newDynamic(p Pos) *dynamic {
	d := dynamic{
		static: s,
		Pos:    p,
	}
	switch {
	case s.circleRadius != nil:
		d.circleRadius = s.circleRadius
	case s.originPolygon != nil:
		p := polygon{
			edges:          s.originPolygon.polygon.edges,
			separatingAxes: s.originPolygon.polygon.separatingAxes,
			vertices:       s.originPolygon.vertices,
		}
		d.polygon = &p
	default:
		panic("static.dynamic")
	}
	return &d
}

// Create a static polygon. Shifts the polygon lines if not origin centered. The passed offset
// should be used to indicate position relative to origin.
func newStaticPolygon(lines []line, offset Pos) static {
	op := newOriginPolygon(lines)
	return static{offset: offset, originPolygon: &op}
}

// Create a static circle. Typical constructor.
func newStaticCircle(r int, offset Pos) static {
	return static{offset: offset, circleRadius: &r}
}

// TODO: low priorty but I should add smoke screen unit tests to catch regressions before release.
func (d *dynamic) move(p Pos, a float64) {
	xOffset := d.static.offset.X*math.Cos(a) - d.static.offset.Y*math.Sin(a)
	yOffset := d.static.offset.X*math.Sin(a) + d.static.offset.Y*math.Cos(a)
	x := p.X + xOffset
	y := p.Y + yOffset

	d.Pos = Pos{X: x, Y: y}

	if d.polygon != nil {

		vertices := []Pos{}
		for _, v := range d.static.originPolygon.vertices {
			v := Pos{
				X: x + (v.X*math.Cos(a) - v.Y*math.Sin(a)),
				Y: y + (v.X*math.Sin(a) + v.Y*math.Cos(a)),
			}
			vertices = append(vertices, v)
		}
		d.polygon.vertices = vertices

		edges := []line{}
		for _, e := range d.static.originPolygon.edges {
			p1 := Pos{
				X: x + (e.p1.X*math.Cos(a) - e.p1.Y*math.Sin(a)),
				Y: y + (e.p1.X*math.Sin(a) + e.p1.Y*math.Cos(a)),
			}
			p2 := Pos{
				X: x + (e.p2.X*math.Cos(a) - e.p2.Y*math.Sin(a)),
				Y: y + (e.p2.X*math.Sin(a) + e.p2.Y*math.Cos(a)),
			}
			edges = append(edges, line{p1, p2})
		}
		d.polygon.edges = edges

		separatingAxes := []line{}
		for _, e := range d.static.originPolygon.separatingAxes {
			p1 := Pos{
				X: x + (e.p1.X*math.Cos(a) - e.p1.Y*math.Sin(a)),
				Y: y + (e.p1.X*math.Sin(a) + e.p1.Y*math.Cos(a)),
			}
			p2 := Pos{
				X: x + (e.p2.X*math.Cos(a) - e.p2.Y*math.Sin(a)),
				Y: y + (e.p2.X*math.Sin(a) + e.p2.Y*math.Cos(a)),
			}
			separatingAxes = append(separatingAxes, line{p1, p2})
		}
		d.polygon.separatingAxes = separatingAxes
	}
}

// Returns true if the two dynamics collide.
// TODO: long term this will probably need to return distance
func (d1 dynamic) collides(d2 dynamic) bool {
	c1 := d1.circleRadius
	p1 := d1.polygon
	c2 := d2.circleRadius
	p2 := d2.polygon

	switch {
	// Returns true if two circles collide. Uses distance formula.
	case c1 != nil && c2 != nil:
		xDelta := d1.X - d2.X
		yDelta := d1.Y - d2.Y
		rSum := float64(*c1 + *c2)
		squaredDistance := xDelta*xDelta + yDelta*yDelta
		if squaredDistance <= rSum*rSum {
			return true
		} else {
			return false
		}

	// Returns true if two polygons collide. Uses separating axis theorem.
	case p1 != nil && p2 != nil:
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

	// Returns true if the polygon collides with the circle. The algorithm first finds the nearest
	// point on every polygon edge to the circle exterior. Then the algorithm returns true any edge's
	// nearest point is touching or within the circle. A ray cast is used to handle circle within
	// polygon edge case.
	case c1 != nil && p2 != nil, p1 != nil && c2 != nil:
		c := d1
		if c.circleRadius == nil {
			c = d2
		}
		p := p1
		if p == nil {
			p = p2
		}

		rayCast := line{p1: c.Pos, p2: Pos{X: math.MaxFloat64, Y: c.Pos.Y}}
		intersectionCount := 0

		for _, e := range p.edges {
			if rayCast.intersects(e) && !rayCast.parallel(e) {
				intersectionCount++
			}

			x := math.Abs(e.p2.X - e.p1.X)
			y := math.Abs(e.p2.Y - e.p1.Y)
			cosineSimularity := (x*c.Pos.X + y*c.Pos.Y) / ((x * x) + (y * y))

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

			xDelta := nearestPos.X - c.Pos.X
			yDelta := nearestPos.Y - c.Pos.Y
			if (xDelta*xDelta)+(yDelta*yDelta) <= float64(*c.circleRadius)*float64(*c.circleRadius) {
				return true
			}
		}
		return intersectionCount%2 == 1

	default:
		panic(fmt.Sprintf(
			"dynamic.collides: unknown shape combination, d1: %s, d2: %s",
			d1.log(),
			d2.log()))
	}

}

// Allows partitioner layer to treat hexagons as circles. This is simple and correct but
// conservative and may need optimized in the future once units use hexagon bodies.
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

// Two unordered connected Positions.
type line struct {
	p1 Pos
	p2 Pos
}

// Typical equals method.
func (l1 line) equals(l2 line) bool {
	return (l1.p1.equals(l2.p1) && l1.p2.equals(l2.p2)) ||
		(l1.p2.equals(l2.p1) && l1.p1.equals(l2.p2))
}

func (l line) isVertical() bool {
	return l.p1.X == l.p2.X
}

func (l line) contains(p Pos) bool {
	return p.X <= max(l.p1.X, l.p2.X) &&
		p.X >= min(l.p1.X, l.p2.X) && p.Y <= max(l.p1.Y, l.p2.Y) && p.Y >= min(l.p1.Y, l.p2.Y)
}

func (l line) perpendicular() line {
	return line{Pos{X: -l.p1.Y, Y: l.p1.X}, Pos{X: -l.p2.Y, Y: l.p2.X}}
}

func (l1 line) parallel(l2 line) bool {
	switch {
	case l1.isVertical() && l2.isVertical():
		return true
	case l1.isVertical() && !l2.isVertical():
		return false
	case !l1.isVertical() && l2.isVertical():
		return false
	default:
		slope1 := (l1.p2.Y - l1.p1.Y) / (l1.p2.X - l1.p1.X)
		slope2 := (l2.p2.Y - l2.p1.Y) / (l2.p2.X - l2.p1.X)
		return slope1 == slope2
	}
}

func (l1 line) intersects(l2 line) bool {
	o1 := orientation(l1.p1, l1.p2, l2.p1)
	o2 := orientation(l1.p1, l1.p2, l2.p2)
	o3 := orientation(l2.p1, l2.p2, l1.p1)
	o4 := orientation(l2.p1, l2.p2, l1.p2)

	switch {
	case o1 != o2 && o3 != o4:
		return true
	case o1 == collinear && l1.contains(l2.p1):
		return true
	case o2 == collinear && l1.contains(l2.p2):
		return true
	case o3 == collinear && l2.contains(l1.p1):
		return true
	case o4 == collinear && l2.contains(l1.p2):
		return true
	default:
		return false
	}
}

// Typical log method.
func (d dynamic) log() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, nest("Pos", d.Pos.log()))
	if d.polygon != nil {
		a = append(a, nest("polygon", d.polygon.log()))
	}
	if d.circleRadius != nil {
		a = append(a, slog.Int("circleRadius", *d.circleRadius))
	}
	a = append(a, nest("static", d.static.log()))
	return a
}

// Typical log method.
func (s static) log() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, nest("offset", s.offset.log()))
	if s.originPolygon != nil {
		a = append(a, nest("originPolygon", s.originPolygon.log()))
	}
	if s.circleRadius != nil {
		a = append(a, slog.Int("circleRadius", *s.circleRadius))
	}
	return a
}

// Typical log method.
func (p polygon) log() []slog.Attr {
	a := []slog.Attr{}

	edges := []slog.Attr{}
	for _, e := range p.edges {
		edges = append(edges, e.log()...)
	}
	a = append(a, nest("edges", edges))

	axes := []slog.Attr{}
	for _, sa := range p.separatingAxes {
		axes = append(axes, sa.log()...)
	}
	a = append(a, nest("separatingAxes", axes))

	vertices := []slog.Attr{}
	for _, v := range p.vertices {
		vertices = append(vertices, v.log()...)
	}
	a = append(a, nest("vertices", vertices))

	return a
}

// Typical log method.
func (l line) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("p1", l.p1.log()))
	a = append(a, nest("p2", l.p2.log()))
	return a
}

// Do not use outside physics.go. A convex shape.
type polygon struct {
	edges []line
	// not every edge normal is a separating axis. Parallel edges only require testing one axis.
	separatingAxes []line
	vertices       []Pos
}

// Do not use outside physics.go. An origin centered convex shape. Useful for rotations.
type originPolygon struct {
	*polygon
	conservativeRadius float64
}

// Do not use outside physics.go. Instead use a dynamic or static constructor. Creates an origin
// centered polygon by shifting passed lines. Refer to newPolygon() for details.
func newOriginPolygon(lines []line) originPolygon {
	poly := newPolygon(lines)

	var center Pos
	for _, p := range poly.vertices {
		center.X += p.X
		center.Y += p.Y
	}
	xShift := center.X / float64(len(poly.vertices))
	yShift := center.Y / float64(len(poly.vertices))

	uniqPoints := make(map[Pos]struct{})
	var shiftedEdges []line
	for _, e := range poly.edges {
		s := e
		s.p1.X -= xShift
		s.p1.Y -= yShift
		s.p2.X -= xShift
		s.p2.Y -= yShift
		shiftedEdges = append(shiftedEdges, s)
		uniqPoints[s.p1] = struct{}{}
		uniqPoints[s.p2] = struct{}{}
	}

	if len(uniqPoints) != len(poly.edges) {
		panic(fmt.Sprintf(
			"newOriginPolygon: expected %d unique points but found %d",
			len(poly.edges),
			len(uniqPoints)))
	}

	var conservativeRadius float64
	for p := range uniqPoints {
		x := (p.X - center.X)
		y := (p.Y - center.Y)
		r := math.Sqrt(x*x + y*y)
		conservativeRadius = max(conservativeRadius, r)
	}

	shiftedPolygon := newPolygon(shiftedEdges)

	vCenter := Pos{}
	for _, v := range shiftedPolygon.vertices {
		vCenter.X += v.X
		vCenter.Y += v.Y
	}
	if !vCenter.equals(Pos{}) {
		panic("newOriginPolygon: shifted vertices not centered")
	}

	eCenter := Pos{}
	for _, e := range shiftedPolygon.edges {
		eCenter.X += e.p1.X
		eCenter.X += e.p2.X
		eCenter.Y += e.p1.Y
		eCenter.Y += e.p2.Y
	}
	if !eCenter.equals(Pos{}) {
		panic("newOriginPolygon: shifted edges not centered")
	}

	return originPolygon{
		polygon:            &shiftedPolygon,
		conservativeRadius: conservativeRadius,
	}
}

// Do not use outside physics.go. Instead use a dynamic or static constructor. Creates a new polygon
// with basic polygon verifications. Polygon verifications will not prevent creation for all invalid
// polygons, but will catch typical invalid polygons. I suspect it is very hard to verify all
// invalid polygon creation so it is not an efficient use of time. Specifically collinear lines are
// probably tricky.
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
