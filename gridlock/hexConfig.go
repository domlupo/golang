// The configuration of the hexagon repeated in the grid.

package main

import "fmt"

// Hexagon configuration
type hexConfig struct {
	hexWidth  int
	hexHeight int
	// These offsets define the top left point for the most top left hex look up
	// table pattern. It was chosen since both a positive x and y simplifies code.
	xOffset   int
	yOffset   int
	xIndexMax int
	yIndexMax int
}

func newHexConfig(wc windowConfig, hexWidth int, hexHeight int, xOffset int, yOffset int) hexConfig {
	hc := hexConfig{
		hexWidth:  hexWidth,
		hexHeight: hexHeight,
		xOffset:   xOffset,
		yOffset:   yOffset,
		xIndexMax: wc.width()/hexWidth + 2,
		yIndexMax: wc.height()/(hexHeight/4*3) + 1,
	}
	return hc
}

func (hc hexConfig) minimumDistanceBetweenNonAdjacentHexagons() int {
	return min(hc.hexWidth, hc.hexHeight) / 2
}

// Returns the radius of a circle which fully surrounds a hexagon. Useful for simple
// and conservative hexagon collision detection by treating the hexagon as a circle.
func (hc hexConfig) conservativeHexRadius() int {
	return max(hc.hexWidth, hc.hexHeight) / 2
}

// Base hexagon structure shared amongst all hexagons. The left side of
// the hexagon is aligned with x=0. The top point of the hexagon is aligned with y=0.
func (hc hexConfig) demotifHexLine(hs hexSide) line {
	nw := Pixel{0, hc.hexHeight / 4}
	n := Pixel{hc.hexWidth / 2, 0}
	ne := Pixel{hc.hexWidth, hc.hexHeight / 4}

	sw := Pixel{0, (hc.hexHeight * 3) / 4}
	s := Pixel{hc.hexWidth / 2, hc.hexHeight}
	se := Pixel{hc.hexWidth, (hc.hexHeight * 3) / 4}

	switch hs {
	case hsNone:
		panic("demotifHexLine: hex side none does not have a demotif hex line")
	case hsWest:
		return line{nw, sw}
	case hsNorthwest:
		return line{n, nw}
	case hsNortheast:
		return line{n, ne}
	case hsEast:
		return line{ne, se}
	case hsSouthwest:
		return line{s, sw}
	case hsSoutheast:
		return line{s, se}
	default:
		panic(fmt.Sprintf("demotifHexLine: unknown hex side %s", hs.string()))
	}
}

// Returns all pixels within the demotif hexagon.
func (hc hexConfig) demotifHexPixels() []Pixel {
	demotifPixels := []Pixel{}

	nwPixels := hc.demotifHexLine(hsNorthwest).adjacentPoints(south)
	nePixles := hc.demotifHexLine(hsNortheast).adjacentPoints(south)
	swPixels := hc.demotifHexLine(hsSouthwest).adjacentPoints(north)
	sePixels := hc.demotifHexLine(hsSoutheast).adjacentPoints(north)

	northPixels := append(nwPixels, nePixles...)
	southPixels := append(swPixels, sePixels...)
	for i, np := range northPixels {
		sp := southPixels[i]
		l := line{np, sp}
		if !l.isVertical() {
			panic("demotifHexPixels: expect line to be vertical")
		}
		demotifPixels = append(demotifPixels, l.adjacentPoints(east)...)
	}
	return demotifPixels
}
