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
	xOffset       int
	yOffset       int
	radius        *int
	originPolygon *originPolygon
}

// TODO: don't pass unit. in the future when I start setting attached field, unit does not needed passed.
func (e equipment) draw(dc drawConfig, pixelMultiplier float32, u *unit) {
	a := u.angle
	xOffset := float64(e.xOffset)*math.Cos(a) - float64(e.yOffset)*math.Sin(a)
	yOffset := float64(e.xOffset)*math.Sin(a) + float64(e.yOffset)*math.Cos(a)

	if e.originPolygon != nil {
		for _, l := range *e.originPolygon {
			r1X := float64(l.p1.X)*math.Cos(a) - float64(l.p1.Y)*math.Sin(a)
			r1Y := float64(l.p1.X)*math.Sin(a) + float64(l.p1.Y)*math.Cos(a)
			r2X := float64(l.p2.X)*math.Cos(a) - float64(l.p2.Y)*math.Sin(a)
			r2Y := float64(l.p2.X)*math.Sin(a) + float64(l.p2.Y)*math.Cos(a)

			vector.StrokeLine(
				dc.screen,
				(float32(u.X)+float32(r1X+xOffset))*pixelMultiplier,
				(float32(u.Y)+float32(r1Y+yOffset))*pixelMultiplier,
				(float32(u.X)+float32(r2X+xOffset))*pixelMultiplier,
				(float32(u.Y)+float32(r2Y+yOffset))*pixelMultiplier,
				float32(dc.defaultLineWidth)*pixelMultiplier,
				color.Black,
				true)
		}
	}

	if e.radius != nil {
		vector.StrokeCircle(
			dc.screen,
			(float32(u.X)+float32(xOffset))*pixelMultiplier,
			(float32(u.Y)+float32(yOffset))*pixelMultiplier,
			float32(*e.radius)*pixelMultiplier,
			float32(dc.defaultLineWidth)*pixelMultiplier,
			color.Black,
			false,
		)
	}
}
