// Accepts player and developer input to affect game components and debugger.

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type controller struct {
	developerMode         bool
	state                 *state
	computers             *computers
	debugger              *debugger
	toggleRegister        *toggleRegister
	initialSelectedX      float32
	initialSelectedY      float32
	currentSelectedX      float32
	currentSelectedY      float32
	leftMouseHeld         bool
	leftMouseClicked      bool
	priorLeftMouseClicked bool
	leftMouseReleased     bool
	rightMouseClicked     bool
}

func newController(s *state, c *computers, d *debugger, tr *toggleRegister) *controller {
	controller := controller{
		developerMode:    true,
		state:            s,
		computers:        c,
		debugger:         d,
		toggleRegister:   tr,
		initialSelectedX: float32(-1.0),
		initialSelectedY: float32(-1.0),
		currentSelectedX: float32(-1.0),
		currentSelectedY: float32(-1.0),
	}
	return &controller
}

func (c *controller) processInput() {
	c.toggleRegister.checkAllTogglePresses()

	c.leftMouseReleased = false
	c.leftMouseHeld = false

	graphicalX, graphicalY := ebiten.CursorPosition()
	_, mouseWheelDelta := ebiten.Wheel()
	c.debugger.updateHeatmapPercent(-int(mouseWheelDelta))

	x := float32(graphicalX) / c.state.grid.pixelMultiplier
	y := float32(graphicalY) / c.state.grid.pixelMultiplier
	pixel := Pixel{int(x), int(y)}

	c.priorLeftMouseClicked = c.leftMouseClicked
	c.leftMouseClicked = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
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
	c.debugger.debug(c.developerMode, ebiten.IsKeyPressed(ebiten.KeySpace), pixel)
	if c.rightMouseClicked {
		filter := []player{pNeutral}
		filter = append(filter, c.computers.active()...)
		c.state.updateTargets(pixel, filter)
	}
	if c.leftMouseReleased {
		c.updateSelected(
			Pixel{int(c.initialSelectedX), int(c.initialSelectedY)},
			Pixel{int(c.currentSelectedX), int(c.currentSelectedY)})
	}
}

func (c controller) updateSelected(clickPosition Pixel, releasePosition Pixel) {
	var priorSingleSelect GridPosition
	for _, n := range c.state.nodes {
		if n.singleSelect {
			priorSingleSelect = n.gridPos
		}
		n.singleSelect = false
		n.multiSelect = false
	}

	clickedGridPos, _ := c.state.grid.getGridPos(clickPosition)
	clickedNode, nodeClicked := c.state.nodes[clickedGridPos]
	if nodeClicked && clickPosition.equals(releasePosition) &&
		clickedGridPos.equals(priorSingleSelect) {
		for _, n := range c.state.nodes {
			if n.player == clickedNode.player && n.unitConfig.name == clickedNode.unitConfig.name {
				n.multiSelect = true
			}
		}
		return
	} else if nodeClicked && clickPosition.equals(releasePosition) {
		clickedNode.singleSelect = true
		return
	}

	var topLeftPixel Pixel
	var bottomRightPixel Pixel

	switch {
	case clickPosition.X < releasePosition.X && clickPosition.Y < releasePosition.Y:
		topLeftPixel = clickPosition
		bottomRightPixel = releasePosition
	case clickPosition.X > releasePosition.X && clickPosition.Y > releasePosition.Y:
		topLeftPixel = releasePosition
		bottomRightPixel = clickPosition
	case clickPosition.X > releasePosition.X && clickPosition.Y < releasePosition.Y:
		topLeftPixel = Pixel{releasePosition.X, clickPosition.Y}
		bottomRightPixel = Pixel{clickPosition.X, releasePosition.Y}
	case clickPosition.X < releasePosition.X && clickPosition.Y > releasePosition.Y:
		topLeftPixel = Pixel{clickPosition.X, releasePosition.Y}
		bottomRightPixel = Pixel{releasePosition.X, clickPosition.Y}
	}

	// TODO: Fix bad UI behavior when cursor ends off screen. Right now, when a cursor
	// ends off screen the behavior differs based on which side it ends. If the cursor
	// ends left or above window, the position is incorrectly set to {0, 0} which
	// usually causes too much selection to occur. The behavior is correct when
	// the cursor ends right or below window. It is tricky and error prone to fix this
	// due to hexagon LUT repetition, so I am punting it.
	topLeftGridPos, _ := c.state.grid.getGridPos(topLeftPixel)
	bottomRightGridPos := c.state.grid.getGridPosOutbounds(bottomRightPixel)

	var firstSelectedNode *node
	for _, n := range c.state.nodes {
		// TODO: Fix broken code or remove it? Not sure if this is still needed.
		if (topLeftGridPos.X < n.gridPos.X) &&
			(n.gridPos.X < bottomRightGridPos.X) &&
			(topLeftGridPos.Y < n.gridPos.Y) &&
			(n.gridPos.Y < bottomRightGridPos.Y) {

			if firstSelectedNode == nil {
				firstSelectedNode = n
				firstSelectedNode.singleSelect = true
				firstSelectedNode.multiSelect = false
			} else {
				firstSelectedNode.singleSelect = false
				firstSelectedNode.multiSelect = true
				n.multiSelect = true
			}
		}
	}

	if c.developerMode {
		for u := range c.state.units() {
			if (topLeftPixel.pos().X < u.X) &&
				(u.X < bottomRightPixel.pos().X) &&
				(topLeftPixel.pos().Y < u.Y) &&
				(u.Y < bottomRightPixel.pos().Y) {
				u.selected = true
			} else {
				u.selected = false
			}
		}

		for _, h := range c.state.hexagons {
			p := c.state.grid.getPixelOutbounds(h.gridPos, hpCenter)
			if (topLeftPixel.X < p.X) &&
				(p.X < bottomRightPixel.X) &&
				(topLeftPixel.Y < p.Y) &&
				(p.Y < bottomRightPixel.Y) {
				h.selected = true
			} else {
				h.selected = false
			}
		}

		for _, n := range c.state.nodes {
			p := c.state.grid.getPixelOutbounds(n.gridPos, hpCenter)
			if (topLeftPixel.X < p.X) &&
				(p.X < bottomRightPixel.X) &&
				(topLeftPixel.Y < p.Y) &&
				(p.Y < bottomRightPixel.Y) {
				n.selected = true
			} else {
				n.selected = false
			}
		}

		visited := make(map[id]struct{})
		for _, o := range c.state.obstacles {
			g := *o.visibilityGrid
			obstaclePoint := g[0]
			if _, ok := visited[obstaclePoint.visibilityGridId]; ok {
				continue
			}
			visited[obstaclePoint.visibilityGridId] = struct{}{}
			for i, op := range *o.visibilityGrid {
				if (topLeftPixel.X < op.X) &&
					(op.X < bottomRightPixel.X) &&
					(topLeftPixel.Y < op.Y) &&
					(op.Y < bottomRightPixel.Y) {
					op.selected = true
				} else {
					op.selected = false
				}
				ops := *o.visibilityGrid
				ops[i] = op
			}
		}
	}

}
