// The configuration of the hexagon repeated in the grid.

package main

type side int

const (
	sNone side = iota
	sInclusive
	sExclusive
)

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

// For 1280p * 720p resolution, a hexagon height of 40p is ideally 80/sqrt(3) wide.
// 80/sqrt(3) is 46.188 which is approximately 46p.

func newHexConfig() hexConfig {
	wc := newWindowConfig()
	// For 1280p * 720p resolution, a hexagon height of 40p is ideally 80/sqrt(3) wide.
	// 80/sqrt(3) is 46.188 which is approximately 46p.
	hexWidth := 46
	hexHeight := 40
	hc := hexConfig{
		hexWidth:  hexWidth,
		hexHeight: hexHeight,
		xOffset:   19,
		yOffset:   10,
		xIndexMax: wc.width/hexWidth + 2,
		yIndexMax: wc.height/(hexHeight/4*3) + 1,
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
