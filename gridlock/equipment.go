// Equipment provides functionality to specific units. A unit may have 0, 1 or several equipments.

package main

import (
	"image/color"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Equipment provides functionality to certain units. A unit may have 0, 1 or several equipments.
type equipment struct {
	attached     *unit
	nearest      *unit
	trackLayers  layers
	ignoreTarget func(t *unit) bool // TODO: decide if this should be renamed. "target" seems wrong.
	ignore       func(e *equipment, t *unit) bool
	effectiveArea
}

// Denotes the zone affected by equipment. Different types of equipment have different zone shapes
// which is indicated via set or unset pointer. Every equipment has a single zone shape.
type effectiveArea struct {
	static
	dynamic
}

func (equip equipment) log() []slog.Attr {
	a := []slog.Attr{}
	if equip.nearest == nil {
		a = append(a, slog.String("nearest", "none"))
	} else {
		a = append(a, slog.Any("nearest", equip.nearest.logShort()))
	}
	return a
}

// move the equipment around its attached unit.
// TODO: decide if this logic should exist within polygon and circle
// TODO: Need to rotate vertex then line? This is incorrect
// TODO: Need to handle separating axis.
func (equip *equipment) move() {
	a := equip.attached.angle
	xOffset := equip.static.offset.X*math.Cos(a) - equip.static.offset.Y*math.Sin(a)
	yOffset := equip.static.offset.X*math.Sin(a) + equip.static.offset.Y*math.Cos(a)
	x := equip.attached.X + xOffset
	y := equip.attached.Y + yOffset

	if equip.dynamic.polygon != nil {
		// TODO: refactor to remove the duplicate logic?
		vertices := []Pos{}
		for _, v := range equip.static.originPolygon.vertices {
			v := Pos{
				X: x + (v.X*math.Cos(a) - v.Y*math.Sin(a)),
				Y: y + (v.X*math.Sin(a) + v.Y*math.Cos(a)),
			}
			vertices = append(vertices, v)
		}
		equip.dynamic.polygon.vertices = vertices

		edges := []line{}
		for _, e := range equip.static.originPolygon.edges {
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

		separatingAxes := []line{}
		for _, e := range equip.static.originPolygon.separatingAxes {
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
	}

	c := equip.dynamic.circle
	if c != nil {
		c.center = Pos{X: x, Y: y}
	}
}

func (equip equipment) draw(dc drawConfig, pixelMultiplier float32) {
	p := equip.dynamic.polygon
	if p != nil {
		for _, e := range p.edges {
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

	c := equip.dynamic.circle
	if c != nil {
		vector.StrokeCircle(
			dc.screen,
			float32(c.center.X)*pixelMultiplier,
			float32(c.center.Y)*pixelMultiplier,
			float32(c.radius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			color.Black,
			false,
		)
	}
}
