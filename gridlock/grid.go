// Represents a low level hexagonal grid.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

type gridConfig struct {
	pixelWidth  int
	pixelHeight int
	basePixels  int
	// TODO: move pixel multiplier to a better location
	// pixelMultiplier must only be used for drawing output like image or processing screen input.
	// Never use for game logic. float32 rather than idiomatic float64 because the Ebiten draw APIs I
	// use are often float32 over float64.
	pixelMultiplier                         float32
	hexPixelWidth                           int
	hexPixelHeight                          int
	hexXPixelOffset                         int
	hexYPixelOffset                         int
	xIndexMax                               int
	yIndexMax                               int
	chebyshevPixelLimit                     int
	minimumPixelsBetweenNonAdjacentHexagons int
	// Returns the radius of a circle which fully surrounds a hexagon. Useful for simple
	// // and conservative hexagon collision detection by treating the hexagon as a circle.
	conservativeHexPixelRadius int
}

func newGridConfig(pixelMultiplier float32) gridConfig {
	gc := gridConfig{
		// For 1280p * 720p resolution, a hexagon height of 40p is ideally 80/sqrt(3) wide.
		// 80/sqrt(3) is 46.188 which is approximately 46p.
		pixelWidth:      1280,
		pixelHeight:     720,
		pixelMultiplier: pixelMultiplier,
		hexPixelWidth:   46,
		hexPixelHeight:  40,
		hexXPixelOffset: 19,
		hexYPixelOffset: 10,
	}
	gc.basePixels = gc.pixelWidth * gc.pixelHeight
	gc.xIndexMax = gc.pixelWidth/gc.hexPixelWidth + 2
	gc.yIndexMax = gc.pixelHeight/(gc.hexPixelHeight*3/4) + 1
	gc.chebyshevPixelLimit = max(gc.pixelWidth, gc.pixelHeight)
	gc.minimumPixelsBetweenNonAdjacentHexagons = min(gc.hexPixelWidth, gc.hexPixelHeight) / 2
	gc.conservativeHexPixelRadius = max(gc.hexPixelWidth, gc.hexPixelHeight) / 2
	return gc
}

func defaultGridConfig() gridConfig {
	return newGridConfig(1.0)
}

// A hexagonal grid.
type grid struct {
	gridConfig
	hexTriangleLUT []hexagonTriangle
}

// The grid struct has a field called hexTriangleLUT which is a look up table (LUT)
//
// The hexagon triangle LUT is a hard coded mapping to aid with converting
// Pixel to hexagon grid position. It avoids the need of CPU intensive math.
// It helps CPU performance by removing the need for floating point or sqrt math logic.
//
// Specific CPU profiling did not prove this optimization was needed but prior CPU profiling
// for this game showed that integer math is much more efficient. Other Pixel to
// hexagon algorithms appeared as complex or even more complex than this method so this
// more optimized solution was implemented first.
//
// It is small portion of the hexagonal grid that is repeated multiple times.
// It covers 460 (46 cols, 10 rows) of Pixels. The top left point in the most
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
// that is repeated multiple times. It is usually referenced as "inverted"
//
//  \     /
//   \   /
//    \ /
//
// This is possible because this section becomes the prior section by following these steps:
//   - horizontally flip the entire section.
//   - split entire section vertically in half.
//   - swap the two vertical halves.
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
// split the entire section vertically in half
//
//  \** **/
//   \* */
//    \ /
//
// swap the two vertical halves.
//
//  **/ \**
//  */   \*
//  /     \
//
// As shown above, the prior section's middle triangle is now converted into
// this section's left and right triangle.

// Position for low resolution game components such as nodes, obstacles, hexagons, etc
type GridPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type clusivity int

const (
	cNone clusivity = iota
	cInclusive
	cExclusive
)

// A value stored in a lookup table to help determine grid position.
type hexagonTriangle int

// A common interior or exterior point of a hexagon.
type hexPoint int

// One of the six sides of a hexagon.
type hexSide int

// Represents a hexagon Pixel shared amongst all hexagons. The left side of
// the hexagon is aligned with x=0. The top point of the hexagon is aligned with y=0.
// By motif definition, x must be less than hex width and y must be less than hex height.
type demotif struct {
	Pixel
}

type direction int

const (
	noDirection direction = iota
	north
	northeast
	east
	southeast
	south
	southwest
	west
	northwest
)

func directions() []direction {
	return []direction{north, northeast, east, southeast, south, southwest, west, northwest}
}

func (d direction) string() string {
	switch d {
	case noDirection:
		return "noDirection"
	case north:
		return "north"
	case northeast:
		return "northeast"
	case east:
		return "east"
	case southeast:
		return "southeast"
	case south:
		return "south"
	case southwest:
		return "southwest"
	case west:
		return "west"
	case northwest:
		return "northwest"
	default:
		panic(fmt.Sprintf("direction.string: unknown direction %d", d))
	}
}

type pixelLine struct {
	p1 Pixel
	p2 Pixel
}

func (l pixelLine) isVertical() bool {
	return l.p1.X == l.p2.X
}

// Typical log method.
func (l pixelLine) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("p1", l.p1.log()))
	a = append(a, nest("p2", l.p2.log()))
	return a
}

func newGrid(gc gridConfig) grid {
	g := grid{
		gridConfig: gc,
	}
	g.hexTriangleLUT = g.setupHexTriangleLUT()
	return g
}

// Setup a look up table to help determine grid position
func (g grid) setupHexTriangleLUT() []hexagonTriangle {
	switchTriangle := make(map[Pixel]hexagonTriangle)
	for _, p := range g.demotifLinePixels(hsNorthwest, cExclusive) {
		switchTriangle[p] = leftHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNorthwest, cInclusive) {
		switchTriangle[p] = middleHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNortheast, cInclusive) {
		switchTriangle[p] = middleHexTriangle
	}
	for _, p := range g.demotifLinePixels(hsNortheast, cExclusive) {
		switchTriangle[p] = rightHexTriangle
	}

	var lut []hexagonTriangle
	for y := 0; y < (g.hexPixelHeight / 4); y++ {
		triangle := leftHexTriangle
		for x := 0; x < g.hexPixelWidth; x++ {
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
			panic("grid.setupHexTriangleLUT: The LUT should not have any unknown hex triangles")
		case leftHexTriangle:
			leftTriangleCount++
		case middleHexTriangle:
			middleTriangleCount++
		case rightHexTriangle:
			rightTriangleCount++
		}
	}
	logger.log(
		infoLevel,
		setupLogic,
		"grid.setupHexTriangleLUT",
		slog.Int("leftHexTriangleCount", leftTriangleCount),
		slog.Int("MiddleHexTriangleCount", middleTriangleCount),
		slog.Int("RightHexTriangleCount", rightTriangleCount),
	)
	if leftTriangleCount+rightTriangleCount != middleTriangleCount {
		panic("grid.setupHexTriangleLUT: invalid hex triangle counts. " +
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
		panic(fmt.Sprintf("hexPoint.string: cannot convert hexPoint %d", hp))
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
		panic(fmt.Sprintf("hexSide.string: cannot convert hexSide %d", hs))
	}
}

func hexPoints() []hexPoint {
	return []hexPoint{hpCenter, hpNorth, hpNortheast, hpSoutheast, hpSouth, hpSouthwest, hpNorthwest}
}

func hexSides() []hexSide {
	return []hexSide{hsNortheast, hsEast, hsSoutheast, hsSouthwest, hsWest, hsNorthwest}
}

// Typical log method.
func (gp GridPosition) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("X", gp.X))
	a = append(a, slog.Int("Y", gp.Y))
	return a
}

// Typical equals method.
func (gp1 GridPosition) equals(gp2 GridPosition) bool {
	return gp1.X == gp2.X && gp1.Y == gp2.Y
}

func (g grid) allAdjacentGridPos(gp GridPosition) []GridPosition {
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
		panic("grid.adjacent: need a direction to find adjacent grid position")
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
		panic("grid.adjacent: cannot find adjacent grid position due to unknown hex side")
	}

	if x < 0 || y < 0 || y >= g.yIndexMax || (y%2 == 0 && x >= g.xIndexMax-1) ||
		(y%2 == 1 && x >= g.xIndexMax) {
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
	if !g.contains(p) {
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
	x := gp.X * g.hexPixelWidth
	y := gp.Y * (g.hexPixelHeight / 4 * 3)
	if gp.Y%2 == 1 {
		x -= (g.hexPixelWidth) / 2
	}

	// The values are found from the 1. hex look up table defined in hexConfig.go
	// and 2. arbitrary rule of hexagons sides are left inclusive and right exclusive.
	// These offsets are found from the hex look up table start point which is
	// located at (xOffset, yOffset).
	switch hexPoint {
	case hpCenter:
		y -= g.hexPixelHeight / 4
	case hpNorth:
		y -= g.hexPixelHeight / 4 * 3
	case hpNorthwest:
		x -= g.hexPixelWidth / 2
		y -= g.hexPixelHeight / 2
	case hpNortheast:
		x += g.hexPixelWidth/2 - 1
		y -= g.hexPixelHeight / 2
	case hpSouth:
		y += g.hexPixelHeight/4 - 1
	case hpSouthwest:
		x -= g.hexPixelWidth / 2
		y -= 1
	case hpSoutheast:
		x += g.hexPixelWidth/2 - 1
		y -= 1
	}

	return Pixel{x + g.hexXPixelOffset, y + g.hexYPixelOffset}
}

func (g grid) getGridPos(p Pixel) (GridPosition, bool) {
	if !g.contains(p) {
		return GridPosition{}, false
	}
	return g.getGridPosOutbounds(p), true
}

func (g grid) getGridPosOutbounds(p Pixel) GridPosition {
	xIndex := (p.X+g.hexPixelWidth-g.hexXPixelOffset)/g.hexPixelWidth - 1
	yIndex := (p.Y+(g.hexPixelHeight/4*3)-g.hexYPixelOffset)/(g.hexPixelHeight/4*3) - 2
	innerX := (p.X + g.hexPixelWidth - g.hexXPixelOffset) % g.hexPixelWidth
	innerY := (p.Y + (g.hexPixelHeight * 3 / 2) - g.hexYPixelOffset) % (g.hexPixelHeight * 3 / 2)

	firstYSectionStart := 0
	secondYSectionStart := firstYSectionStart + g.hexPixelHeight/4
	thirdYSectionStart := secondYSectionStart + g.hexPixelHeight/2
	fourthYSectionStart := thirdYSectionStart + g.hexPixelHeight/4

	var xIndexOffset int
	var yIndexOffset int
	switch {
	case innerY >= firstYSectionStart && innerY < secondYSectionStart:
		triangle := g.hexTriangleLUT[innerX+innerY*g.hexPixelWidth]
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
		hexLookUpYOffset := (g.hexPixelWidth * g.hexPixelHeight / 4) * 3
		var hexLookUpXOffset int
		if innerX >= 0 && innerX < (g.hexPixelWidth/2) {
			hexLookUpXOffset = (g.hexPixelWidth / 2)
		} else {
			hexLookUpXOffset = -(g.hexPixelWidth / 2)
		}

		triangle := g.hexTriangleLUT[innerX+(innerY*g.hexPixelWidth)+hexLookUpXOffset-hexLookUpYOffset]
		switch {
		case triangle == leftHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		case triangle == middleHexTriangle && innerX < (g.hexPixelWidth/2):
			xIndexOffset = 0
			yIndexOffset = 2
		case triangle == middleHexTriangle && innerX >= (g.hexPixelWidth/2):
			xIndexOffset = 1
			yIndexOffset = 2
		case triangle == rightHexTriangle:
			xIndexOffset = 1
			yIndexOffset = 1
		}
	case innerY >= fourthYSectionStart && innerY < (g.hexPixelHeight*3/2):
		if innerX >= 0 && innerX < (g.hexPixelWidth/2) {
			xIndexOffset = 0
			yIndexOffset = 2
		} else {
			xIndexOffset = 1
			yIndexOffset = 2
		}
	default:
		panic(fmt.Sprintf(
			"grid.getGridPosOutbounds: could not find a y section for the inner y: %d and Pixel %s",
			innerY,
			p.log()))
	}

	return GridPosition{xIndex + xIndexOffset, yIndex + yIndexOffset}
}

// Remove the repeated hexagon grid pattern from the Pixel to find the base
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
	p.X += (g.hexPixelWidth / 2) - g.hexXPixelOffset
	p.Y += (g.hexPixelHeight / 2)
	// Remove the repeated hexagon grid pattern.
	p.Y -= gp.Y * (g.hexPixelHeight * 3 / 4)
	switch gp.Y % 2 {
	case 0:
		p.X -= gp.X * g.hexPixelWidth
	case 1:
		p.X -= gp.X * g.hexPixelWidth
		p.X += g.hexPixelWidth / 2
	default:
		panic("grid.demotif: invalid modulo")
	}
	if p.X >= g.hexPixelWidth || p.Y >= g.hexPixelHeight {
		panic(fmt.Sprintf(
			"grid.demotif: by motif definition, x must be less than hex width and y must be less than "+
				"hex height %s",
			p.log()))
	}
	return demotif{p}, true
}

func (g grid) demotifLinePixels(hs hexSide, clusivity clusivity) []Pixel {
	nw := Pixel{0, g.hexPixelHeight / 4}
	n := Pixel{g.hexPixelWidth / 2, 0}
	ne := Pixel{g.hexPixelWidth, g.hexPixelHeight / 4}
	sw := Pixel{0, g.hexPixelHeight * 3 / 4}
	s := Pixel{g.hexPixelWidth / 2, g.hexPixelHeight}
	se := Pixel{g.hexPixelWidth, g.hexPixelHeight * 3 / 4}

	var d direction
	trim := []Pixel{}
	switch {
	case hs == hsWest && clusivity == cInclusive:
		d = east
		trim = append(trim, sw)
	case hs == hsNorthwest && clusivity == cInclusive:
		d = south
		trim = append(trim, n)
	case hs == hsNortheast && clusivity == cInclusive:
		d = south
		trim = append(trim, ne)
	case hs == hsEast && clusivity == cInclusive:
		d = west
		trim = append(trim, Pixel{g.hexPixelWidth - 1, g.hexPixelHeight * 3 / 4})
	case hs == hsSouthwest && clusivity == cInclusive:
		d = north
		trim = append(trim, Pixel{g.hexPixelWidth / 2, g.hexPixelHeight - 1})
	case hs == hsSoutheast && clusivity == cInclusive:
		d = north
		trim = append(trim, Pixel{g.hexPixelWidth, g.hexPixelHeight*3/4 - 1})
	case hs == hsWest && clusivity == cExclusive:
		d = west
		trim = append(trim, Pixel{sw.X - 1, sw.Y})
	case hs == hsNorthwest && clusivity == cExclusive:
		d = north
		trim = append(trim, Pixel{g.hexPixelWidth / 2, -1})
	case hs == hsNortheast && clusivity == cExclusive:
		d = north
		trim = append(trim, Pixel{g.hexPixelWidth, g.hexPixelHeight/4 - 1})
	case hs == hsEast && clusivity == cExclusive:
		d = east
		trim = append(trim, Pixel{g.hexPixelWidth, g.hexPixelHeight * 3 / 4})
	case hs == hsSouthwest && clusivity == cExclusive:
		d = south
		trim = append(trim, Pixel{g.hexPixelWidth / 2, g.hexPixelHeight})
	case hs == hsSoutheast && clusivity == cExclusive:
		d = south
		trim = append(trim, Pixel{g.hexPixelWidth, g.hexPixelHeight * 3 / 4})
	default:
		panic("grid.demotifLinePixels unknown hex side and clusivity combination")
	}

	var l pixelLine
	switch hs {
	case hsNone:
		panic("grid.demotifLinePixels: hex side none does not have Pixels")
	case hsWest:
		l = pixelLine{nw, sw}
	case hsNorthwest:
		l = pixelLine{nw, n}
	case hsNortheast:
		l = pixelLine{n, ne}
	case hsEast:
		l = pixelLine{ne, se}
	case hsSouthwest:
		l = pixelLine{s, sw}
	case hsSoutheast:
		l = pixelLine{s, se}
	default:
		panic(fmt.Sprintf("grid.demotifLinePixels: unknown hex side %s", hs.string()))
	}
	return l.adjacentPixels(d, trim...)
}

// Returns all pixels within the demotif hexagon.
func (g grid) demotifHexPixels() []Pixel {
	demotifPixels := []Pixel{}

	nwPixels := g.demotifLinePixels(hsNorthwest, cInclusive)
	nePixels := g.demotifLinePixels(hsNortheast, cInclusive)
	swPixels := g.demotifLinePixels(hsSouthwest, cInclusive)
	sePixels := g.demotifLinePixels(hsSoutheast, cInclusive)

	northPixels := append(nwPixels, nePixels...)
	southPixels := append(swPixels, sePixels...)
	for i, np := range northPixels {
		sp := southPixels[i]
		l := pixelLine{np, sp}
		if !l.isVertical() {
			panic("grid.demotifHexPixels: expect line to be vertical")
		}
		demotifPixels = append(demotifPixels, l.adjacentPixels(east)...)
	}
	return demotifPixels
}

// Returns all pixels adjacent to the demotif hexagon.
func (g grid) demotifHexPixelsAdjacent() map[Pixel]hexSide {
	demotifPixels := make(map[Pixel]hexSide)
	for _, hs := range hexSides() {
		for _, p := range g.demotifLinePixels(hs, cExclusive) {
			demotifPixels[p] = hs
		}
	}
	return demotifPixels
}

func (g grid) contains(p Pixel) bool {
	xContained := p.X >= 0 && p.X < g.pixelWidth
	yContained := p.Y >= 0 && p.Y < g.pixelHeight
	return xContained && yContained
}

func (g grid) collidesWithBorder(p Pixel, offset int) bool {
	if p.X+offset > g.pixelWidth ||
		p.X-offset < 0 ||
		p.Y+offset > g.pixelHeight ||
		p.Y-offset < 0 {
		return true
	} else {
		return false
	}
}

func (g grid) sum(p1 Pixel, p2 Pixel) (Pixel, bool) {
	p := Pixel{X: p1.X + p2.X, Y: p1.Y + p2.Y}
	if !g.contains(p) {
		return Pixel{}, false
	} else {
		return p, true
	}
}

func (g grid) allAdjacentPixels(pixel Pixel) []Pixel {
	a := []Pixel{}
	for _, d := range directions() {
		if p, ok := g.adjacentPixel(pixel, d); ok {
			a = append(a, p)
		}
	}
	return a
}

func (g grid) adjacentPixel(p Pixel, d direction) (Pixel, bool) {
	var x int
	var y int

	switch d {
	case north:
		y = -1
	case northeast:
		x = 1
		y = -1
	case east:
		x = 1
	case southeast:
		x = 1
		y = 1
	case south:
		y = 1
	case southwest:
		x = -1
		y = 1
	case west:
		x = -1
		x = -1
	case northwest:
		x = -1
		y = -1
	default:
		panic(fmt.Sprintf(
			"grid.adjacentPixel cannot get adjacent pixel for %s unknown direction: %d",
			p.log(),
			d))
	}

	a := Pixel{X: p.X + x, Y: p.Y + y}
	if !g.contains(a) {
		return Pixel{}, false
	}
	return a, true
}

// Returns the adjacent Pixels in the passed direction. The function has multiple assumptions.
// Assumes a square grid with typical compass directions. Assumes a square grid where east is
// positive X and south is positive Y. X is typical and Y is opposite of typical for Cartesian
// graph. Assumes the line should be drawn in the top left corner of the grid square. Sloped lines
// determine Pixel direction of split grid squares via which Pixel has more area. Trims any Pixels
// passed as trim but panics if they are not an adjacent Pixel.
func (l pixelLine) adjacentPixels(d direction, trim ...Pixel) []Pixel {
	switch d {
	case north, south, west, east:
	default:
		panic(fmt.Sprintf(
			"grid.adjacentPixels: cannot get Pixels in direction %s since it is unimplemented",
			d.string()))
	}

	var w Pixel
	var e Pixel
	switch {
	case l.p1.X < l.p2.X:
		w = l.p1
		e = l.p2
	default:
		w = l.p2
		e = l.p1
	}
	nY := max(l.p1.Y, l.p2.Y)
	sY := min(l.p1.Y, l.p2.Y)

	var pixels []Pixel
	switch {
	case l.isVertical() && d == west:
		y := sY
		for y <= nY {
			pixels = append(pixels, Pixel{w.X - 1, y})
			y += 1
		}
	case l.isVertical() && d == east:
		y := sY
		for y <= nY {
			pixels = append(pixels, Pixel{w.X, y})
			y += 1
		}
	case !l.isVertical() && d == west, !l.isVertical() && d == east:
		panic(fmt.Sprintf(
			"grid.adjacentPixels: cannot get pixels in direction %s for non vertical lines since it is "+
				"unimplemented",
			d.string()))
	case l.isVertical() && d == north:
		pixels = append(pixels, Pixel{w.X, sY - 1})
	case l.isVertical() && d == south:
		pixels = append(pixels, Pixel{w.X, nY})
	default:
		// Ebitengine y is opposite of typical Cartesian grid, slope equation is not typical.
		// e.Y - w.Y instead of typical w.Y - e.Y
		slope := float64(e.Y-w.Y) / float64(w.X-e.X)
		x := float64(w.X) + 0.5
		for x <= float64(e.X+1) {
			yIntercept := float64(w.Y) + (slope * float64(w.X))
			// Ebitengine y is opposite of typical Cartesian grid, y intercept equation is not typical.
			// y := yIntercept - slope*x instead of y := slope*x + yIntercept.
			y := yIntercept - slope*x
			xInteger, _ := math.Modf(x)
			yInteger, yDecimal := math.Modf(y)
			xInt := int(xInteger)
			yInt := int(yInteger)
			switch {
			case yDecimal == 0.5:
				panic("grid.adjacentPixels: undefined behavior when y midpoint occurs")
			case yDecimal < 0.5:
				yInt -= 1
			}
			var yOffset int
			if d == south {
				yOffset++
			}
			pixels = append(pixels, Pixel{xInt, yInt + yOffset})
			x += 1
		}
	}

	for _, p := range trim {
		i := slices.Index(pixels, p)
		if i == -1 {
			panic(fmt.Sprintf(
				"grid.adjacentPixels: cannot trim %s since it is not an adjacent pixel",
				p.log()))
		}
		pixels = slices.Delete(pixels, i, i+1)
	}
	return pixels
}
