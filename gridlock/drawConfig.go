// The draw configuration for game components shown to player.

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type drawTriangleColor struct {
	red   float32
	green float32
	blue  float32
	alpha float32
}

// Drawing Configurations for game components which are shown to player
type drawConfig struct {
	screen *ebiten.Image
	// windowMultiplier must only be used for drawing images and windows.
	// Never use for game logic. float32 rather than idiomatic float64
	// because the Ebiten draw APIs I use are often float32 over float64.
	windowMultiplier              float32
	defaultLineWidth              int
	selectectionBoxLineWidth      int
	missileExplodeRadiusLineWidth int
	healthLineWidth               int
	targetNodeLineWidth           int
	hexLineWidth                  int
	pixelSize                     int
	hexGridLineWidth              int
	hexGridIndexCount             int
	backgroundColor               color.RGBA
	highlightColor                color.RGBA
	unitHealthColor               color.RGBA
	nodeHealthColor               color.RGBA
	pOneNodeTargetColor           color.RGBA
	pOneUnitColor                 color.RGBA
	pOneInfantryColor             color.RGBA
	pTwoNodeTargetColor           color.RGBA
	pTwoUnitColor                 color.RGBA
	pTwoInfantryColor             color.RGBA
	pNeutralNodeTargetColor       color.RGBA
	hexGridColor                  color.Color
	drawTriangleDefaultImage      *ebiten.Image
	pOneNodeColor                 drawTriangleColor
	pTwoNodeColor                 drawTriangleColor
	pNeutralNodeColor             drawTriangleColor
	nodeControlColor              drawTriangleColor
	pngBomber                     *ebiten.Image
	pngFighter                    *ebiten.Image
	pngMissile                    *ebiten.Image
	pngFlare                      *ebiten.Image
}

// Invoke only once per game session. NewImage calls are expensive and are not
// idempotent for specific edge cases. https://github.com/hajimehoshi/ebiten/issues/3335
func newDrawConfig() drawConfig {
	// TODO: check if this should use a subimage if final game uses draw triangle.
	// Not investigating now since draw triangle may not be used in final game.
	drawTriangleDefaultImage := ebiten.NewImage(1, 1)
	drawTriangleDefaultImage.Fill(color.White)
	highlightColor := color.RGBA{70, 255, 0, 200}
	pngBomer, _, bomberErr := ebitenutil.NewImageFromFile("assets/bomber.png")
	if bomberErr != nil {
		panic("newDrawConfig: cannot load bomber png")
	}
	pngFighter, _, fighterErr := ebitenutil.NewImageFromFile("assets/fighter.png")
	if fighterErr != nil {
		panic("newDrawConfig: cannot load fighter png")
	}
	pngMissile, _, missileErr := ebitenutil.NewImageFromFile("assets/missile.png")
	if missileErr != nil {
		panic("newDrawConfig: cannot load missile png")
	}
	wc := newWindowConfig()
	dc := drawConfig{
		screen:                        nil,
		windowMultiplier:              wc.multiplier,
		defaultLineWidth:              1,
		selectectionBoxLineWidth:      1,
		missileExplodeRadiusLineWidth: 1,
		healthLineWidth:               2,
		targetNodeLineWidth:           2,
		hexLineWidth:                  2,
		pixelSize:                     1,
		hexGridLineWidth:              1,
		hexGridIndexCount:             32,
		backgroundColor:               color.RGBA{200, 200, 200, 255},
		highlightColor:                highlightColor,
		unitHealthColor:               highlightColor,
		nodeHealthColor:               highlightColor,
		pOneNodeTargetColor:           color.RGBA{0, 0, 255, 100},
		pOneUnitColor:                 color.RGBA{0, 0, 255, 100},
		pOneInfantryColor:             color.RGBA{0, 0, 255, 50},
		pTwoNodeTargetColor:           color.RGBA{255, 0, 0, 100},
		pTwoUnitColor:                 color.RGBA{255, 0, 0, 100},
		pTwoInfantryColor:             color.RGBA{255, 0, 0, 50},
		pNeutralNodeTargetColor:       color.RGBA{100, 60, 100, 200},
		hexGridColor:                  color.White,
		drawTriangleDefaultImage:      drawTriangleDefaultImage,
		pOneNodeColor:                 drawTriangleColor{red: 0.5, blue: 1, green: 0.5, alpha: 0.5},
		pTwoNodeColor:                 drawTriangleColor{red: 1, blue: 0.5, green: 0.5, alpha: 0.5},
		pNeutralNodeColor:             drawTriangleColor{red: 0.5, blue: 0.5, green: 0.5, alpha: 0.5},
		nodeControlColor:              drawTriangleColor{red: 0.7, blue: 0.7, green: 0.7, alpha: 1},
		pngBomber:                     pngBomer,
		pngFighter:                    pngFighter,
		pngMissile:                    pngMissile,
		// TODO: Add flare png
		pngFlare: pngMissile,
	}
	return dc
}
