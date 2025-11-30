// Accepts player and developer input to affect game components and debugger.

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type controller struct {
	developerMode            bool
	state                    *state
	computers                *computers
	debugger                 *debugger
	initialSelectedX         float32
	initialSelectedY         float32
	currentSelectedX         float32
	currentSelectedY         float32
	leftMouseHeld            bool
	leftMouseClicked         bool
	priorLeftMouseClicked    bool
	leftMouseReleased        bool
	middleMouseScrollClicked bool
	rightMouseClicked        bool
}

func newController(s *state, c *computers, d *debugger) *controller {
	controller := controller{
		developerMode:    true,
		state:            s,
		computers:        c,
		debugger:         d,
		initialSelectedX: float32(-1.0),
		initialSelectedY: float32(-1.0),
		currentSelectedX: float32(-1.0),
		currentSelectedY: float32(-1.0),
	}
	return &controller
}

func (c *controller) control() {
	c.leftMouseReleased = false
	c.leftMouseHeld = false

	graphicalX, graphicalY := ebiten.CursorPosition()
	_, mouseWheelDelta := ebiten.Wheel()
	c.debugger.updateHeatmapPercent(-int(mouseWheelDelta))

	wc := newWindowConfig()
	x := float32(graphicalX) / wc.multiplier
	y := float32(graphicalY) / wc.multiplier
	pixel := Pixel{int(x), int(y)}

	c.priorLeftMouseClicked = c.leftMouseClicked
	c.leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	c.middleMouseScrollClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	c.rightMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	c.currentSelectedX = x
	c.currentSelectedY = y
	c.debugger.pixel = pixel
	if c.priorLeftMouseClicked && !c.leftMouseClicked {
		c.leftMouseReleased = true
	} else if !c.priorLeftMouseClicked && c.leftMouseClicked {
		c.initialSelectedX = x
		c.initialSelectedY = y
	} else if c.leftMouseClicked {
		c.leftMouseHeld = true
	}
	if c.developerMode && c.middleMouseScrollClicked {
		c.debugger.debug(pixel)
	}
	if c.rightMouseClicked {
		filter := []player{pNeutral}
		filter = append(filter, c.computers.active()...)
		c.state.updateTargets(position{lLand, pixel}, filter)
	}
	if c.leftMouseReleased {
		c.state.updateSelected(
			Pixel{int(c.initialSelectedX), int(c.initialSelectedY)},
			Pixel{int(c.currentSelectedX), int(c.currentSelectedY)})
	}
}
