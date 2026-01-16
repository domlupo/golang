// Equipment provides functionality to specific units. A unit may have 0, 1 or several equipments.

package main

import (
	"fmt"
	"image/color"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Equipment provides functionality to certain units. A unit may have 0, 1 or several equipments.
type equipment struct {
	name equipName
	*effectiveArea
	attached     *unit
	nearest      *unit
	trackLayers  layers
	ignoreTarget func(t *unit) bool // TODO: decide if this should be renamed. "target" seems wrong.
	ignore       func(e *equipment, t *unit) bool
}

// Typical log method.
func (equip equipment) log() []slog.Attr {
	a := []slog.Attr{}

	a = append(a, slog.String("name", equip.name.string()))
	a = append(a, nest("effectiveArea", equip.effectiveArea.log()))

	if equip.nearest == nil {
		a = append(a, slog.String("nearest", "none"))
	} else {
		a = append(a, slog.Any("nearest", equip.nearest.logShort()))
	}

	return a
}

type equipName int

const (
	enNone equipName = iota
	enRotaryCannnon
	enRadar
)

// Typical string method.
func (en equipName) string() string {
	switch en {
	case enNone:
		return "none"
	case enRotaryCannnon:
		return "rotaryCannon"
	case enRadar:
		return "radar"
	default:
		panic(fmt.Sprintf("equipName.string unknown equip name %d", en))
	}
}

// Denotes the zone affected by equipment.
type effectiveArea struct {
	static  static
	dynamic *dynamic
}

// Typical log method.
func (ea effectiveArea) log() []slog.Attr {
	a := []slog.Attr{}
	a = append(a, nest("static", ea.static.log()))
	a = append(a, nest("dynamic", ea.dynamic.log()))
	return a
}

// TODO: decide if this logic should exist within polygon and circle
// move the equipment around its attached unit.
func (equip *equipment) move() {
	a := equip.attached.angle
	xOffset := equip.static.offset.X*math.Cos(a) - equip.static.offset.Y*math.Sin(a)
	yOffset := equip.static.offset.X*math.Sin(a) + equip.static.offset.Y*math.Cos(a)
	x := equip.attached.X + xOffset
	y := equip.attached.Y + yOffset

	equip.dynamic.Pos = Pos{X: x, Y: y}

	if equip.dynamic.polygon != nil {

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
		equip.dynamic.polygon.separatingAxes = separatingAxes
	}
}

func (equip equipment) draw(dc drawConfig, pixelMultiplier float32) {
	p := equip.effectiveArea.dynamic.polygon
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

	d := equip.effectiveArea.dynamic
	if d.circleRadius != nil {
		vector.StrokeCircle(
			dc.screen,
			float32(d.X)*pixelMultiplier,
			float32(d.Y)*pixelMultiplier,
			float32(*d.circleRadius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			color.Black,
			false,
		)
	}
}
