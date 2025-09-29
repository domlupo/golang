// The hexagon configuration is relatively complex due to CPU optimizations, so it is given its own file.
package main

// This global should only be used for function signatures, structs, etc. It can
// be calculated during runtime via hexConfig fields. It is calculated from hexagon width
// * (hexagon height/4). The 4 is from the hexagon position context being 1/4 of the total
// hexagon area.
const hexagonPostitionContextPixelCount = 460

type hexConfig struct {
	hexPositionContext [hexagonPostitionContextPixelCount]hexagonTriangle
	hexWidth           int
	hexHeight          int
	// These offsets define the top left point for the most top left hex position
	// context pattern. It was chosen since both a positive x and y simplifies code.
	xOffset      int
	yOffset      int
	windowWidth  int
	windowHeight int
	xIndexMax    int
	yIndexMax    int
}

func newHexConfig(windowWidth int, windowHeight int) hexConfig {
	hexPositionContext := newHexPositionContext()
	hexWidth := 46
	hexHeight := 40
	xIndexMax := windowWidth/hexWidth + 1
	yIndexMax := windowHeight/(hexHeight/4*3) + 1
	hc := hexConfig{
		hexPositionContext: hexPositionContext,
		// For 1280p * 720p resolution, a hexagon height of 40p is ideally 80/sqrt(3) wide.
		// 80/sqrt(3) is 46.188 which is approximately 46p.
		hexWidth:     hexWidth,
		hexHeight:    hexHeight,
		xOffset:      19,
		yOffset:      10,
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
		xIndexMax:    xIndexMax,
		yIndexMax:    yIndexMax,
	}
	return hc
}

// The hexagon position context is a hardcoded mapping to aid with converting
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
// The hexagon position context is also used for another section of the hexagonal grid
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
// As shown above, the prior sections middile triangle is now converted into
// this sections left and right triangle.

type hexagonTriangle int

const (
	unknownHexTriangle hexagonTriangle = iota
	leftHexTriangle
	middleHexTriangle
	rightHexTriangle
)

func newHexPositionContext() [hexagonPostitionContextPixelCount]hexagonTriangle {
	var pc [hexagonPostitionContextPixelCount]hexagonTriangle

	rowPixelCount := 46

	// row 1
	rowStart := 0
	rowMiddleOffset := 1
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 2
	rowStart = rowPixelCount
	rowMiddleOffset = 3
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 3
	rowStart = rowPixelCount * 2
	rowMiddleOffset = 6
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 4
	rowStart = rowPixelCount * 3
	rowMiddleOffset = 8
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 5
	rowStart = rowPixelCount * 4
	rowMiddleOffset = 10
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 6
	rowStart = rowPixelCount * 5
	rowMiddleOffset = 13
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 7
	rowStart = rowPixelCount * 6
	rowMiddleOffset = 15
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 8
	rowStart = rowPixelCount * 7
	rowMiddleOffset = 17
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 9
	rowStart = rowPixelCount * 8
	rowMiddleOffset = 20
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 10
	rowStart = rowPixelCount * 9
	rowMiddleOffset = 22
	setHexPositionContext(&pc, rowPixelCount, rowStart, rowMiddleOffset)

	return pc
}

func setHexPositionContext(pc *[hexagonPostitionContextPixelCount]hexagonTriangle, rowPixelCount int, rowStart int, rowMiddleOffset int) {
	pixelCount := rowPixelCount / 2

	leftHexStart := rowStart
	middleHexStart := rowStart + pixelCount - rowMiddleOffset
	rightHexStart := middleHexStart + 2*rowMiddleOffset
	for i := leftHexStart; i < middleHexStart; i++ {
		pc[i] = leftHexTriangle
	}
	for i := middleHexStart; i < rightHexStart; i++ {
		pc[i] = middleHexTriangle
	}
	for i := rightHexStart; i < rowStart+rowPixelCount; i++ {
		pc[i] = rightHexTriangle
	}
}
