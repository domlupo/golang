// Equipment provides functionality to specific units. A unit may have 0, 1 or several equipments.

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Equipment provides functionality to certain units. A unit may have 0, 1 or several equipments.
type equipment struct {
	attached *unit
	effectiveArea
}

// Denotes the zone affected by equipment. Different types of equipment have different zone shapes
// which is indicated via set or unset pointer. Every equipment has a single zone shape.
type effectiveArea struct {
	static
	dynamic
}

// The equipment configuration which must not be mutated.
type static struct {
	offset  Pos
	polygon *originPolygon
	radius  *int
}

// The equipment physics which may be mutated.
type dynamic struct {
	polygon polygon
	circle  Pos
}

// Rotate the equipment around its attached unit.
func (equip *equipment) rotate() {
	a := equip.attached.angle
	xOffset := equip.static.offset.X*math.Cos(a) - equip.static.offset.Y*math.Sin(a)
	yOffset := equip.static.offset.X*math.Sin(a) + equip.static.offset.Y*math.Cos(a)
	x := equip.attached.X + xOffset
	y := equip.attached.Y + yOffset
	if equip.static.polygon != nil {
		edges := []line{}
		for _, e := range equip.static.polygon.edges {
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
		equip.dynamic.polygon.edges = edges
	}
	if equip.static.radius != nil {
		equip.circle = Pos{X: x, Y: y}
	}
}

func (equip equipment) draw(dc drawConfig, pixelMultiplier float32) {
	if equip.static.polygon != nil {
		for _, e := range equip.dynamic.polygon.edges {
			vector.StrokeLine(
				dc.screen,
				float32(e.p1.X)*pixelMultiplier,
				float32(e.p1.Y)*pixelMultiplier,
				float32(e.p2.X)*pixelMultiplier,
				float32(e.p2.Y)*pixelMultiplier,
				float32(dc.defaultLineWidth)*pixelMultiplier,
				color.Black,
				true)
		}
	}

	if equip.radius != nil {
		vector.StrokeCircle(
			dc.screen,
			float32(equip.dynamic.circle.X)*pixelMultiplier,
			float32(equip.dynamic.circle.Y)*pixelMultiplier,
			float32(*equip.radius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			color.Black,
			false,
		)
	}
}
