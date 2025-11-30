// Represents a low level hexagonal grid.

package main

import (
	"fmt"
	"log/slog"
)

// The grid struct has a field called hexTriangleLUT which is a look up table (LUT)
//
// The hexagon triangle LUT is a hardcoded mapping to aid with converting
// pixel position to hexagon grid position. It avoids the the need of CPU intensive math.
// It helps CPU performance by removing the need for floating point or sqrt math logic.
//
// Specific CPU profiling did not prove this optimizaton was needed but prior CPU profiling
// for this game showed that integer math is much more efficient. Other pixel position to
// hexagon algorithms appeared as complex or even more complex than this method so this
// more optimized solution was implemented first.
//
// It is small portion of the hexagonal grid that is repeated multiple times.
// It covers 460 (46 cols, 10 rows) of pixels. The top left point in the most
// top left pattern is located at {xOffset, yOffset}.
//
// It has a:
//   - right triangle on the left side.
//   - non right triangle in the middle.
//   - right triangle on the right side.
//
// right triangle on left side
//
// ***/ \
// **/   \
// */     \
//
// non right triangle in middle
//
//    /*\
//   /***\
//  /*****\
//
// right triangle on right side
//
//    / \***
//   /   \**
//  /     \*
//
// The hexagon look up table is also used for another section of the hexagonal grid
// that is repeated multiple times. It is usally referenced as "inverted"
//
//  \     /
//   \   /
//    \ /
//
// This is possible because this section becomes the prior section by following these steps:
//   - horizontally flip the entire section.
//   - split entire section vertically in half.
//   - swap the two vertical halfs.
//
// An example highlighting the middle section is shown below:
//
// non right triangle in middle
//
//    /*\
//   /***\
//  /*****\
//
// horizontally flip the entire section.
//
//  \*****/
//   \***/
//    \*/
//
// split the entir section vertically in half
//
//  \** **/
//   \* */
//    \ /
//
// swap the two vertical halfs.
//
//  **/ \**
//  */   \*
//  /     \
//
// As shown above, the prior section's middle triangle is now converted into
// this section's left and right triangle.

// A hexagonal grid.
type grid struct {
	hexConfig
	hexTriangleLUT []hexagonTriangle
}

// Position for multi pixel game components such as nodes, obstacles, hexagons, etc
type GridPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// A value stored in a lookup table to help determine grid position.
type hexagonTriangle int

// A commmon interior or exterior point of a hexagon.
type hexPoint int

// One of the six sides of a hexagon.
type hexSide int

// Represents a hexagon pixel position shared amongst all hexagons. The left side of
// the hexagon is aligned with x=0. The top point of the hexagon is aligned with y=0.
// By motif definition, x must be less than hex width and y must be less than hex height.
type demotif struct {
	Pixel
}

func newGrid() grid {
	g := grid{hexConfig: newHexConfig()}
	g.hexTriangleLUT = g.setupHexTriangleLUT()
	return g
}

// Setup a look up table to help determine grid position
func (g grid) setupHexTriangleLUT() []hexagonTriangle {
	switchTriangle := make(map[Pixel]hexagonTriangle)
	for _, p := range g.demotifLinePixels(hsNorthwest, sExclusive) {
		switchTriangle[p] = leftHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNorthwest, sInclusive) {
		switchTriangle[p] = middleHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNortheast, sInclusive) {
		switchTriangle[p] = middleHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNortheast, sExclusive) {
		switchTriangle[p] = rightHexTriangle
	}

	var lut []hexagonTriangle
	for y := 0; y < (g.hexHeight / 4); y++ {
		triangle := leftHexTriangle
		for x := 0; x < g.hexWidth; x++ {
			p := Pixel{x, y}
			if t, ok := switchTriangle[p]; ok {
				triangle = t
			}
			lut = append(lut, triangle)
		}
	}

	leftTriangleCount := 0
	middleTriangleCount := 0
	rightTriangleCount := 0
	for _, t := range lut {
		switch t {
		case unknownHexTriangle:
			panic("setupHexTriangleLUT: The look up table should not have any unknown hex triangles")
		case leftHexTriangle:
			leftTriangleCount++
		case middleHexTriangle:
			middleTriangleCount++
		case rightHexTriangle:
			rightTriangleCount++
		}
	}
	logger.log(infoLevel, setupLogic, "setupHexTriangleLUT",
		slog.Int("leftHexTriangleCount", leftTriangleCount),
		slog.Int("MiddleHexTriangleCount", middleTriangleCount),
		slog.Int("RightHexTriangleCount", rightTriangleCount),
	)
	if leftTriangleCount+rightTriangleCount != middleTriangleCount {
		panic("setupHexTriangleLUT: invalid hex triangle counts. " +
			"This panic might be a false positive if hex width isn't 46 and hex height isn't 40")
	}
	return lut
}

const (
	unknownHexTriangle hexagonTriangle = iota
	leftHexTriangle
	middleHexTriangle
	rightHexTriangle
)

const (
	hpCenter hexPoint = iota
	hpNorth
	hpNortheast
	hpSoutheast
	hpSouth
	hpSouthwest
	hpNorthwest
)

const (
	hsNone hexSide = iota
	hsNortheast
	hsEast
	hsSoutheast
	hsSouthwest
	hsWest
	hsNorthwest
)

func (hp hexPoint) string() string {
	switch hp {
	case hpCenter:
		return "hpCenter"
	case hpNorth:
		return "hpNorth"
	case hpNortheast:
		return "hpNortheast"
	case hpSoutheast:
		return "hpSoutheast"
	case hpSouth:
		return "hpSouth"
	case hpSouthwest:
		return "hpSouthwest"
	case hpNorthwest:
		return "hpNorthwest"
	default:
		panic(fmt.Sprintf("Cannot convert hexPoint to string: %d", hp))
	}
}

func (hs hexSide) string() string {
	switch hs {
	case hsNone:
		return "hsNone"
	case hsNortheast:
		return "hsNortheast"
	case hsEast:
		return "hsEast"
	case hsSoutheast:
		return "hsSoutheast"
	case hsSouthwest:
		return "hsSouthwest"
	case hsWest:
		return "hsWest"
	case hsNorthwest:
		return "hsNorthwest"
	default:
		panic(fmt.Sprintf("Cannot convert hexSide to string: %d", hs))
	}
}

func hexPoints() []hexPoint {
	return []hexPoint{hpCenter, hpNorth, hpNortheast, hpSoutheast, hpSouth, hpSouthwest, hpNorthwest}
}

func hexSides() []hexSide {
	return []hexSide{hsNortheast, hsEast, hsSoutheast, hsSouthwest, hsWest, hsNorthwest}
}

func (gp GridPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("x", gp.X))
	a = append(a, slog.Int("y", gp.Y))
	return a
}

func (gp GridPosition) equals(other GridPosition) bool {
	return gp.X == other.X && gp.Y == other.Y
}

func (g grid) allAdjacent(gp GridPosition) []GridPosition {
	a := []GridPosition{}
	for _, hs := range hexSides() {
		if p, ok := g.adjacent(gp, hs); ok {
			a = append(a, p)
		}
	}
	return a
}

func (g grid) adjacent(p GridPosition, side hexSide) (GridPosition, bool) {
	x := p.X
	y := p.Y
	switch side {
	case hsNone:
		panic("need a direction to find adjacent grid position")
	case hsNortheast:
		y--
		if p.Y%2 == 0 {
			x++
		}
	case hsEast:
		x++
	case hsSoutheast:
		y++
		if p.Y%2 == 0 {
			x++
		}
	case hsWest:
		x--
	case hsSouthwest:
		y++
		if p.Y%2 == 1 {
			x--
		}
	case hsNorthwest:
		y--
		if p.Y%2 == 1 {
			x--
		}
	default:
		panic("Cannot find adjacent grid position. Unknown hex side.")
	}

	if x < 0 || y < 0 || y >= g.yIndexMax || (y%2 == 0 && x >= g.xIndexMax-1) || (y%2 == 1 && x >= g.xIndexMax) {
		return GridPosition{}, false
	} else {
		return GridPosition{x, y}, true
	}
}

func (g grid) isAdjacent(gp1 GridPosition, gp2 GridPosition) bool {
	ne, neOk := g.adjacent(gp1, hsNortheast)
	e, eOk := g.adjacent(gp1, hsEast)
	se, seOk := g.adjacent(gp1, hsSoutheast)
	sw, swOk := g.adjacent(gp1, hsSouthwest)
	w, wOk := g.adjacent(gp1, hsWest)
	nw, nwOk := g.adjacent(gp1, hsNorthwest)

	return (neOk && ne.equals(gp2)) ||
		(eOk && e.equals(gp2)) ||
		(seOk && se.equals(gp2)) ||
		(swOk && sw.equals(gp2)) ||
		(wOk && w.equals(gp2)) ||
		(nwOk && nw.equals(gp2))
}

func (g grid) getPixel(gp GridPosition, hexPoint hexPoint) (Pixel, bool) {
	p := g.getPixelOutbounds(gp, hexPoint)
	if !newWindowConfig().contains(p) {
		return Pixel{}, false
	} else {
		return p, true
	}
}

func (g grid) getCenterWithFallback(gp GridPosition) (Pixel, bool) {
	if p, ok := g.getPixel(gp, hpCenter); ok {
		return p, true
	}

	for _, hp := range hexPoints() {
		if p, ok := g.getPixel(gp, hp); ok {
			return p, true
		}
	}
	return Pixel{}, false
}

func (g grid) getPixelOutbounds(gp GridPosition, hexPoint hexPoint) Pixel {
	x := gp.X * g.hexWidth
	y := gp.Y * (g.hexHeight / 4 * 3)
	if gp.Y%2 == 1 {
		x -= (g.hexWidth) / 2
	}

	// The values are found from the 1. hex look up table defined in hexConfig.go
	// and 2. arbitrary rule of hexagons sides are left inclusive and right exclusive.
	// These offsets are found from the hex look up table start point which is
	// located at (xOffset, yOffset).
	switch hexPoint {
	case hpCenter:
		y -= g.hexHeight / 4
	case hpNorth:
		y -= g.hexHeight / 4 * 3
	case hpNorthwest:
		x -= g.hexWidth / 2
		y -= g.hexHeight / 2
	case hpNortheast:
		x += g.hexWidth/2 - 1
		y -= g.hexHeight / 2
	case hpSouth:
		y += g.hexHeight/4 - 1
	case hpSouthwest:
		x -= g.hexWidth / 2
		y -= 1
	case hpSoutheast:
		x += g.hexWidth/2 - 1
		y -= 1
	}

	return Pixel{x + g.xOffset, y + g.yOffset}
}

func (g grid) getGridPos(p Pixel) (GridPosition, bool) {
	wc := newWindowConfig()
	if !wc.contains(p) {
		return GridPosition{}, false
	}
	return g.getGridPosOutbounds(p), true
}

func (g grid) getGridPosOutbounds(p Pixel) GridPosition {
	xIndex := (p.X+g.hexWidth-g.xOffset)/g.hexWidth - 1
	yIndex := (p.Y+(g.hexHeight/4*3)-g.yOffset)/(g.hexHeight/4*3) - 2
	innerX := (p.X + g.hexWidth - g.xOffset) % g.hexWidth
	innerY := (p.Y + (g.hexHeight * 3 / 2) - g.yOffset) % (g.hexHeight * 3 / 2)

	firstYSectionStart := 0
	secondYSectionStart := firstYSectionStart + g.hexHeight/4
	thirdYSectionStart := secondYSectionStart + g.hexHeight/2
	fourthYSectionStart := thirdYSectionStart + g.hexHeight/4

	var xIndexOffset int
	var yIndexOffset int
	switch {
	case innerY >= firstYSectionStart && innerY < secondYSectionStart:
		triangle := g.hexTriangleLUT[innerX+innerY*g.hexWidth]
		switch triangle {
		case leftHexTriangle:
			xIndexOffset = 0
			yIndexOffset = 1
		case middleHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 2
		case rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= secondYSectionStart && innerY < thirdYSectionStart:
		xIndexOffset = 1
		yIndexOffset = 2
	case innerY >= thirdYSectionStart && innerY < fourthYSectionStart:
		hexLookUpYOffset := (g.hexWidth * g.hexHeight / 4) * 3
		var hexLookUpXOffset int
		if innerX >= 0 && innerX < (g.hexWidth/2) {
			hexLookUpXOffset = (g.hexWidth / 2)
		} else {
			hexLookUpXOffset = -(g.hexWidth / 2)
		}

		triangle := g.hexTriangleLUT[innerX+(innerY*g.hexWidth)+hexLookUpXOffset-hexLookUpYOffset]
		switch {
		case triangle == leftHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		case triangle == middleHexTriangle && innerX < (g.hexWidth/2):
			xIndexOffset = 0
			yIndexOffset = 2
		case triangle == middleHexTriangle && innerX >= (g.hexWidth/2):
			xIndexOffset = 1
			yIndexOffset = 2
		case triangle == rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= fourthYSectionStart && innerY < (g.hexHeight*3/2):
		if innerX >= 0 && innerX < (g.hexWidth/2) {
			xIndexOffset = 0
			yIndexOffset = 2
		} else {
			xIndexOffset = 1
			yIndexOffset = 2
		}
	default:
		panic(fmt.Sprintf("getHexGridPostion could not find a y section for the inner y: %d for pixel position %s", innerY, p.log()))
	}

	return GridPosition{xIndex + xIndexOffset, yIndex + yIndexOffset}
}

// Remove the repeated hexagon grid pattern from the pixel position to find the base
// hexagon structure shared amongst all hexagons. The left side of the hexagon is
// aligned with x=0. The top point of the hexagon is aligned with y=0. By motif
// definition, x must be less than hex width and y must be less than hex height.
func (g grid) demotif(p Pixel) (demotif, bool) {
	gp, ok := g.getGridPos(p)
	if !ok {
		return demotif{}, false
	}
	// Offset hexagon to {0, 0}.
	// The left side of the hexagon is aligned with x=0.
	// The top point of the hexagon is aligned with y=0.
	p.X += (g.hexWidth / 2) - g.xOffset
	p.Y += (g.hexHeight / 2)
	// Remove the repeated hexagon grid pattern.
	p.Y -= gp.Y * (g.hexHeight * 3 / 4)
	switch gp.Y % 2 {
	case 0:
		p.X -= gp.X * g.hexWidth
	case 1:
		p.X -= gp.X * g.hexWidth
		p.X += g.hexWidth / 2
	default:
		panic("demotif: invalid modulo")
	}
	if p.X >= g.hexWidth || p.Y >= g.hexHeight {
		panic(fmt.Sprintf("demotif: by motif definition, x must be less than hex width and y must be less than hex height %s", p.log()))
	}
	return demotif{p}, true
}

func (g grid) demotifLinePixels(hs hexSide, side side) []Pixel {
	hc := g.hexConfig

	nw := Pixel{0, hc.hexHeight / 4}
	n := Pixel{hc.hexWidth / 2, 0}
	ne := Pixel{hc.hexWidth, hc.hexHeight / 4}
	sw := Pixel{0, hc.hexHeight * 3 / 4}
	s := Pixel{hc.hexWidth / 2, hc.hexHeight}
	se := Pixel{hc.hexWidth, hc.hexHeight * 3 / 4}

	var d direction
	trim := []Pixel{}
	switch {
	case hs == hsWest && side == sInclusive:
		d = east
		trim = append(trim, sw)
	case hs == hsNorthwest && side == sInclusive:
		d = south
		trim = append(trim, n)
	case hs == hsNortheast && side == sInclusive:
		d = south
		trim = append(trim, ne)
	case hs == hsEast && side == sInclusive:
		d = west
		trim = append(trim, Pixel{hc.hexWidth - 1, hc.hexHeight * 3 / 4})
	case hs == hsSouthwest && side == sInclusive:
		d = north
		trim = append(trim, Pixel{hc.hexWidth / 2, hc.hexHeight - 1})
	case hs == hsSoutheast && side == sInclusive:
		d = north
		trim = append(trim, Pixel{hc.hexWidth, hc.hexHeight*3/4 - 1})
	case hs == hsWest && side == sExclusive:
		d = west
		trim = append(trim, Pixel{sw.X - 1, sw.Y})
	case hs == hsNorthwest && side == sExclusive:
		d = north
		trim = append(trim, Pixel{hc.hexWidth / 2, -1})
	case hs == hsNortheast && side == sExclusive:
		d = north
		trim = append(trim, Pixel{hc.hexWidth, hc.hexHeight/4 - 1})
	case hs == hsEast && side == sExclusive:
		d = east
		trim = append(trim, Pixel{hc.hexWidth, hc.hexHeight * 3 / 4})
	case hs == hsSouthwest && side == sExclusive:
		d = south
		trim = append(trim, Pixel{hc.hexWidth / 2, hc.hexHeight})
	case hs == hsSoutheast && side == sExclusive:
		d = south
		trim = append(trim, Pixel{hc.hexWidth, hc.hexHeight * 3 / 4})
	default:
		panic("demotifLinePixels unknown hex side and side combination")
	}

	var l line
	switch hs {
	case hsNone:
		panic("demotifLinePixels: hex side none does not have pixels")
	case hsWest:
		l = line{nw, sw}
	case hsNorthwest:
		l = line{nw, n}
	case hsNortheast:
		l = line{n, ne}
	case hsEast:
		l = line{ne, se}
	case hsSouthwest:
		l = line{s, sw}
	case hsSoutheast:
		l = line{s, se}
	default:
		panic(fmt.Sprintf("demotifLinePixels: unknown hex side %s", hs.string()))
	}
	return l.adjacentPoints(d, trim...)
}

// Returns all pixels within the demotif hexagon.
func (g grid) demotifHexPixels() []Pixel {
	demotifPixels := []Pixel{}

	nwPixels := g.demotifLinePixels(hsNorthwest, sInclusive)
	nePixles := g.demotifLinePixels(hsNortheast, sInclusive)
	swPixels := g.demotifLinePixels(hsSouthwest, sInclusive)
	sePixels := g.demotifLinePixels(hsSoutheast, sInclusive)

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

// Returns all pixels adjacent to the demotif hexagon.
func (g grid) demotifHexPixelsAdjacent() map[Pixel]hexSide {
	demotifPixels := make(map[Pixel]hexSide)
	for _, hs := range hexSides() {
		for _, p := range g.demotifLinePixels(hs, sExclusive) {
			demotifPixels[p] = hs
		}
	}
	return demotifPixels
}
