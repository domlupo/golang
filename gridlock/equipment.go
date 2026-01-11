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
	offset  Pos
	polygon polygon
}

// Rotate the equipment around its attached unit.
func (equip *equipment) rotate() {
	a := equip.attached.angle
	equip.dynamic.offset = Pos{
		X: equip.static.offset.X*math.Cos(a) - equip.static.offset.Y*math.Sin(a),
		Y: equip.static.offset.X*math.Sin(a) + equip.static.offset.Y*math.Cos(a),
	}
	if equip.static.polygon != nil {
		edges := []line{}
		for _, e := range equip.static.polygon.edges {
			p1 := Pos{
				X: e.p1.X*math.Cos(a) - e.p1.Y*math.Sin(a),
				Y: e.p1.X*math.Sin(a) + e.p1.Y*math.Cos(a),
			}
			p2 := Pos{
				X: e.p2.X*math.Cos(a) - e.p2.Y*math.Sin(a),
				Y: e.p2.X*math.Sin(a) + e.p2.Y*math.Cos(a),
			}
			edges = append(edges, line{p1, p2})
		}
		equip.dynamic.polygon.edges = edges
	}
}

func (equip equipment) draw(dc drawConfig, pixelMultiplier float32) {
	if equip.static.polygon != nil {
		for _, e := range equip.dynamic.polygon.edges {
			p1X := equip.attached.X + equip.dynamic.offset.X + e.p1.X
			p1Y := equip.attached.Y + equip.dynamic.offset.Y + e.p1.Y
			p2X := equip.attached.X + equip.dynamic.offset.X + e.p2.X
			p2Y := equip.attached.Y + equip.dynamic.offset.Y + e.p2.Y

			vector.StrokeLine(
				dc.screen,
				float32(p1X)*pixelMultiplier,
				float32(p1Y)*pixelMultiplier,
				float32(p2X)*pixelMultiplier,
				float32(p2Y)*pixelMultiplier,
				float32(dc.defaultLineWidth)*pixelMultiplier,
				color.Black,
				true)
		}
	}

	if equip.radius != nil {
		x := equip.attached.X + equip.dynamic.offset.X
		y := equip.attached.Y + equip.dynamic.offset.Y
		vector.StrokeCircle(
			dc.screen,
			float32(x)*pixelMultiplier,
			float32(y)*pixelMultiplier,
			float32(*equip.radius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			color.Black,
			false,
		)
	}
}
