// A project specific concept for forces which contriubute to local pathfinding
package main

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type fluxVector struct {
	magnitude int
	direction int
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
		fluxVector:  fluxVector{0, noDirection},
	}
	return f
}

func (f *flux) getFluxDelta() *gridPosition {
	var magnitude = fluxMagnitude
	switch f.direction {
	case noDirection:
		return &gridPosition{0, 0}
	case north:
		return &gridPosition{0, -magnitude}
	case northEast:
		return &gridPosition{magnitude, -magnitude}
	case east:
		return &gridPosition{magnitude, 0}
	case southEast:
		return &gridPosition{magnitude, magnitude}
	case south:
		return &gridPosition{0, magnitude}
	case southWest:
		return &gridPosition{-magnitude, magnitude}
	case west:
		return &gridPosition{-magnitude, 0}
	case northWest:
		return &gridPosition{-magnitude, -magnitude}
	default:
		panic("unknown direction for getFluxDelta")
	}
}

func (f *flux) setFluxVector(s *state, xMax int, yMax int) {
	var northPos = f.gridPos.adjacent(north, xMax, yMax)
	var northEastPos = f.gridPos.adjacent(northEast, xMax, yMax)
	var eastPos = f.gridPos.adjacent(east, xMax, yMax)
	var southEastPos = f.gridPos.adjacent(southEast, xMax, yMax)
	var southPos = f.gridPos.adjacent(south, xMax, yMax)
	var southWestPos = f.gridPos.adjacent(southWest, xMax, yMax)
	var westPos = f.gridPos.adjacent(west, xMax, yMax)
	var northWestPos = f.gridPos.adjacent(northWest, xMax, yMax)

	var xFlux = 0
	var yFlux = 0

	if northPos != nil {
		yFlux = yFlux + s.fluxes[*northPos].staticUnits
	}
	if northEastPos != nil {
		xFlux = xFlux - s.fluxes[*northEastPos].staticUnits
		yFlux = yFlux + s.fluxes[*northEastPos].staticUnits
	}
	if eastPos != nil {
		xFlux = xFlux - s.fluxes[*eastPos].staticUnits
	}
	if southEastPos != nil {
		xFlux = xFlux - s.fluxes[*southEastPos].staticUnits
		yFlux = yFlux - s.fluxes[*southEastPos].staticUnits
	}
	if southPos != nil {
		yFlux = yFlux - s.fluxes[*southPos].staticUnits
	}
	if southWestPos != nil {
		xFlux = xFlux + s.fluxes[*southWestPos].staticUnits
		yFlux = yFlux - s.fluxes[*southWestPos].staticUnits
	}
	if westPos != nil {
		xFlux = xFlux + s.fluxes[*westPos].staticUnits
	}
	if northWestPos != nil {
		xFlux = xFlux + s.fluxes[*northWestPos].staticUnits
		yFlux = yFlux + s.fluxes[*northWestPos].staticUnits
	}

	if xFlux == 0 && yFlux == 0 {
		f.fluxVector = fluxVector{0, noDirection}
	} else if xFlux == 0 && yFlux < 0 {
		f.fluxVector = fluxVector{1, north}
	} else if xFlux > 0 && yFlux < 0 {
		f.fluxVector = fluxVector{1, northEast}
	} else if xFlux > 0 && yFlux == 0 {
		f.fluxVector = fluxVector{1, east}
	} else if xFlux > 0 && yFlux > 0 {
		f.fluxVector = fluxVector{1, southEast}
	} else if xFlux == 0 && yFlux > 0 {
		f.fluxVector = fluxVector{1, south}
	} else if xFlux < 0 && yFlux > 0 {
		f.fluxVector = fluxVector{1, southWest}
	} else if xFlux < 0 && yFlux == 0 {
		f.fluxVector = fluxVector{1, west}
	} else if xFlux < 0 && yFlux < 0 {
		f.fluxVector = fluxVector{1, northWest}
	}
}

func (f *flux) refresh() {
	f.staticUnits = 0
	f.magnitude = 0
	f.direction = noDirection
}

func (f *flux) draw(dc drawContext, fluxPos gridPosition) {
	if f.magnitude == 0 {
		return
	}

	var center = getPixelPosition(fluxPos, innerGridPosition{20, 20})
	var xEnd = center.x
	var yEnd = center.y

	switch f.direction {
	case noDirection:
		panic("need a direction to draw flux")
	case north:
		yEnd -= (dc.gridSize / 2)
	case northEast:
		xEnd += (dc.gridSize / 2)
		yEnd -= (dc.gridSize / 2)
	case east:
		xEnd += (dc.gridSize / 2)
	case southEast:
		xEnd += (dc.gridSize / 2)
		yEnd += (dc.gridSize / 2)
	case south:
		yEnd += (dc.gridSize / 2)
	case southWest:
		xEnd -= (dc.gridSize / 2)
		yEnd += (dc.gridSize / 2)
	case west:
		xEnd -= (dc.gridSize / 2)
	case northWest:
		xEnd -= (dc.gridSize / 2)
		yEnd -= (dc.gridSize / 2)
	}

	vector.StrokeLine(
		dc.screen,
		float32(center.x)*dc.windowMultiplier,
		float32(center.y)*dc.windowMultiplier,
		float32(xEnd)*dc.windowMultiplier,
		float32(yEnd)*dc.windowMultiplier,
		float32(fluxLineWidth)*dc.windowMultiplier,
		fluxColor,
		false,
	)

	vector.DrawFilledCircle(
		dc.screen,
		float32(xEnd)*dc.windowMultiplier,
		float32(yEnd)*dc.windowMultiplier,
		float32(fluxDirectionRadius)*dc.windowMultiplier,
		fluxColor,
		false,
	)
}
