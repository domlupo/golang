// A line is two unordered connected Pixels.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

// Two unordered connected Pixels.
type line struct {
	p1 Pixel
	p2 Pixel
}

func (l line) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("p1", l.p1.log()))
	a = append(a, nest("p2", l.p2.log()))
	return a
}

func (l line) isVertical() bool {
	return l.p1.X == l.p2.X
}

func (l line) contains(p Pixel) bool {
	return p.X <= max(l.p1.X, l.p2.X) &&
		p.X >= min(l.p1.X, l.p2.X) && p.Y <= max(l.p1.Y, l.p2.Y) && p.Y >= min(l.p1.Y, l.p2.Y)
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

// Returns the adjacent points in the passed direction. The function has multiple assumptions.
// Assumes a square grid with typical compass directions. Assumes a square grid where east is
// positive x and south is positive y. x is typical and y is opposite of typical for Cartesian
// graph. Assumes the line should be drawn in the top left corner of the grid square. Sloped lines
// determine point direction of split grid squares via which point has more area. Trims any points
// passed as trim but panics if they are not an adjacent point.
func (l line) adjacentPoints(d direction, trim ...Pixel) []Pixel {
	switch d {
	case north, south, west, east:
	default:
		panic(fmt.Sprintf(
			"adjacentPoints: cannot get points in direction %s since it is unimplemented", d.string()))
	}

	var w Pixel
	var e Pixel
	switch {
	case l.p1.X < l.p2.X:
		w = l.p1
		e = l.p2
	default:
		w = l.p2
		e = l.p1
	}
	nY := max(l.p1.Y, l.p2.Y)
	sY := min(l.p1.Y, l.p2.Y)

	var points []Pixel
	switch {
	case l.isVertical() && d == west:
		y := sY
		for y <= nY {
			points = append(points, Pixel{w.X - 1, y})
			y += 1
		}
	case l.isVertical() && d == east:
		y := sY
		for y <= nY {
			points = append(points, Pixel{w.X, y})
			y += 1
		}
	case !l.isVertical() && d == west, !l.isVertical() && d == east:
		s := "adjacentPoints: cannot get points in direction %s for non vertical lines since it is " +
			"unimplemented"
		panic(fmt.Sprintf(s, d.string()))
	case l.isVertical() && d == north:
		points = append(points, Pixel{w.X, sY - 1})
	case l.isVertical() && d == south:
		points = append(points, Pixel{w.X, nY})
	default:
		// Ebitengine y is opposite of typical Cartesian grid, slope equation is not typical.
		// e.Y - w.Y instead of typical w.Y - e.Y
		slope := float64(e.Y-w.Y) / float64(w.X-e.X)
		x := float64(w.X) + 0.5
		for x <= float64(e.X+1) {
			yIntercept := float64(w.Y) + (slope * float64(w.X))
			// Ebitengine y is opposite of typical Cartesian grid, y intercept equation is not typical.
			// y := yIntercept - slope*x instead of y := slope*x + yIntercept.
			y := yIntercept - slope*x
			xInteger, _ := math.Modf(x)
			yInteger, yDecimal := math.Modf(y)
			xInt := int(xInteger)
			yInt := int(yInteger)
			switch {
			case yDecimal == 0.5:
				panic("adjacentPoints: undefined behavior when y midpoint occurs")
			case yDecimal < 0.5:
				yInt -= 1
			}
			var yOffset int
			if d == south {
				yOffset++
			}
			points = append(points, Pixel{xInt, yInt + yOffset})
			x += 1
		}
	}

	for _, p := range trim {
		i := slices.Index(points, p)
		if i == -1 {
			panic(
				fmt.Sprintf("adjacentPoints: cannot trim %s since it is not an adjacent point", p.log()))
		}
		points = slices.Delete(points, i, i+1)
	}
	return points
}

type originPolygon []line

// Creates a new origin centered polygon with basic polygon verifications. If possible, the passed
// lines will be shifted to create an origin centered polygon. Polygon verifications will not
// prevent creation for all invalid polygons, but will catch typical invalid polygons. I suspect it
// is very hard to verify all invalid polygon creation so it is not an efficient use of time.
// Specifically collinear lines are probably tricky. func newOriginPolygon(lines []line)
func newOriginPolygon(lines []line) originPolygon {
	if len(lines) != 4 {
		panic("newOriginPolygon: only 4 sided polygons are implemented")
	}

	var p Pixel
	for _, l := range lines {
		p.X += l.p1.X
		p.Y += l.p1.Y
	}

	if !p.equals(Pixel{}) && !(p.X%len(lines) == 0 && p.Y%len(lines) == 0) {
		panic(fmt.Sprintf(
			"newOriginPolygon: polygon is not origin centered and cannot be shifted: %s", p.log()))
	}

	xShift := p.X / len(lines) / 2
	yShift := p.Y / len(lines) / 2
	var shifted []line
	for _, l := range lines {
		s := l
		s.p1.X += xShift
		s.p1.Y += yShift
		s.p2.X += xShift
		s.p2.Y += yShift
		shifted = append(shifted, s)
	}

	sorted := []line{shifted[0]}
	current := shifted[0].p1
	for {
		priorSortCount := len(sorted)
		for _, l := range shifted {
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

		if len(sorted) == len(shifted) {
			break
		} else if len(sorted) == priorSortCount {
			panic("newOriginPolygon: valid polygon connections cannot be created with provided lines")
		}
	}

	return shifted
}
