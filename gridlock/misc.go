// Misc code which is used in several other files but not large enough to need its own file.

package main

import "fmt"

const maxPercent = 100
const minPercent = 0

type player int

const (
	pNone player = iota
	pBlue
	pRed
	pNeutral
)

func (p player) String() string {
	switch p {
	case pNone:
		return "pNone"
	case pBlue:
		return "pBlue"
	case pRed:
		return "pRed"
	case pNeutral:
		return "pNeutral"
	default:
		panic(fmt.Sprintf("Cannot convert player to string: %d", p))
	}
}

type windowConfig struct {
	width  int
	height int
}

func newWindowConfig() windowConfig {
	wc := windowConfig{
		width:  1280,
		height: 720,
	}
	return wc
}

func (wc windowConfig) maxPixels() int {
	return wc.height * wc.width
}

func (wc windowConfig) inbounds(p PixelPosition) bool {
	return !(p.X < 0 || p.X > wc.width-1 || p.Y < 0 || p.Y > wc.height-1)
}
