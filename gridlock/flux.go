// A project specific concept for forces which contriubute to local pathfinding
package main

import "github.com/hajimehoshi/ebiten/v2/vector"

type fluxVector struct {
	magnitude int
	hexSide   int
}

type flux struct {
	gridPos     gridPosition
	staticUnits int
	fluxVector
}

func newFlux(gridPos gridPosition) flux {
	f := flux{
		gridPos:     gridPos,
		staticUnits: 0,
		fluxVector:  fluxVector{0, noSide},
	}
	return f
}

func (f *flux) getFluxDelta(magnitude int) gridPosition {
	switch f.hexSide {
	case noSide:
		return gridPosition{0, 0}
	case northeastSide:
		return gridPosition{magnitude, -magnitude}
	case eastSide:
		return gridPosition{magnitude, 0}
	case southeastSide:
		return gridPosition{magnitude, magnitude}
	case southwestSide:
		return gridPosition{-magnitude, magnitude}
	case westSide:
		return gridPosition{-magnitude, 0}
	case northwestSide:
		return gridPosition{-magnitude, -magnitude}
	default:
		panic("unknown direction for getFluxDelta")
	}
}

func (f *flux) setFluxVector(s *state) {
	northEastPos, northEastPosOk := f.gridPos.adjacent(northeastSide, s.hexConfig)
	eastPos, eastPosOk := f.gridPos.adjacent(eastSide, s.hexConfig)
	southEastPos, southEastPosOk := f.gridPos.adjacent(southeastSide, s.hexConfig)
	southWestPos, southWestPosOk := f.gridPos.adjacent(southwestSide, s.hexConfig)
	westPos, westPosOk := f.gridPos.adjacent(westSide, s.hexConfig)
	northWestPos, northWestPosOk := f.gridPos.adjacent(northwestSide, s.hexConfig)

	xFlux := 0
	yFlux := 0

	if northEastPosOk {
		xFlux = xFlux - s.fluxes[northEastPos].staticUnits
		yFlux = yFlux + s.fluxes[northEastPos].staticUnits
	}
	if eastPosOk {
		xFlux = xFlux - s.fluxes[eastPos].staticUnits
	}
	if southEastPosOk {
		xFlux = xFlux - s.fluxes[southEastPos].staticUnits
		yFlux = yFlux - s.fluxes[southEastPos].staticUnits
	}
	if southWestPosOk {
		xFlux = xFlux + s.fluxes[southWestPos].staticUnits
		yFlux = yFlux - s.fluxes[southWestPos].staticUnits
	}
	if westPosOk {
		xFlux = xFlux + s.fluxes[westPos].staticUnits
	}
	if northWestPosOk {
		xFlux = xFlux + s.fluxes[northWestPos].staticUnits
		yFlux = yFlux + s.fluxes[northWestPos].staticUnits
	}

	switch {
	case xFlux == 0 && yFlux == 0:
		f.fluxVector = fluxVector{0, noSide}
	case xFlux > 0 && yFlux < 0:
		f.fluxVector = fluxVector{1, northeastSide}
	case xFlux > 0 && yFlux == 0:
		f.fluxVector = fluxVector{1, eastSide}
	case xFlux > 0 && yFlux > 0:
		f.fluxVector = fluxVector{1, southeastSide}
	case xFlux < 0 && yFlux > 0:
		f.fluxVector = fluxVector{1, southwestSide}
	case xFlux < 0 && yFlux == 0:
		f.fluxVector = fluxVector{1, westSide}
	case xFlux < 0 && yFlux < 0:
		f.fluxVector = fluxVector{1, northwestSide}
	}
}

func (f *flux) refresh() {
	f.staticUnits = 0
	f.magnitude = 0
	f.hexSide = noSide
}

func (f *flux) draw(dc drawContext, fluxPos gridPosition) {
	if f.magnitude == 0 {
		return
	}

	center, _ := getPixelPos(fluxPos, center, &dc.hexConfig)
	xEnd := center.x
	yEnd := center.y

	switch f.hexSide {
	case noSide:
		panic("need a direction to draw flux")
	case northeastSide:
		xEnd += (dc.hexWidth / 2)
		yEnd -= (dc.hexWidth / 2)
	case eastSide:
		xEnd += (dc.hexWidth / 2)
	case southeastSide:
		xEnd += (dc.hexWidth / 2)
		yEnd += (dc.hexWidth / 2)
	case southwestSide:
		xEnd -= (dc.hexWidth / 2)
		yEnd += (dc.hexWidth / 2)
	case westSide:
		xEnd -= (dc.hexWidth / 2)
	case northwestSide:
		xEnd -= (dc.hexWidth / 2)
		yEnd -= (dc.hexWidth / 2)
	}

	vector.StrokeLine(
		dc.screen,
		float32(center.x)*dc.windowMultiplier,
		float32(center.y)*dc.windowMultiplier,
		float32(xEnd)*dc.windowMultiplier,
		float32(yEnd)*dc.windowMultiplier,
		float32(dc.fluxLineWidth)*dc.windowMultiplier,
		dc.fluxColor,
		false,
	)

	vector.DrawFilledCircle(
		dc.screen,
		float32(xEnd)*dc.windowMultiplier,
		float32(yEnd)*dc.windowMultiplier,
		float32(dc.fluxDirectionRadius)*dc.windowMultiplier,
		dc.fluxColor,
		false,
	)
}
