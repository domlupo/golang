// Responsible for rendering pixels and images to the screen.

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type drawer struct {
	state    *state
	debugger *debugger
}

func newDrawer(s *state, d *debugger) drawer {
	drawer := drawer{
		state:    s,
		debugger: d,
	}
	return drawer
}

func (d drawer) draw(screen *ebiten.Image, dc drawConfig) {
	screen.Fill(dc.backgroundColor)
	d.debugger.draw(dc)
	d.state.draw(dc)
}
