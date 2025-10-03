// The hexagon configuration is relatively complex. A look up table is implemented to optimize CPU usage.
package main

// This global should only be used for function signatures, structs, etc. It can
// be calculated during runtime via hexConfig fields. It is calculated from hexagon width
// * (hexagon height/4). The 4 is from the hexagon look up table being 1/4 of the total
// hexagon height.
const hexLookUpTablePixelCount = 460

type hexConfig struct {
	lookUpTable [hexLookUpTablePixelCount]hexagonTriangle
	hexWidth    int
	hexHeight   int
	// These offsets define the top left point for the most top left hex look up
	// table pattern. It was chosen since both a positive x and y simplifies code.
	xOffset      int
	yOffset      int
	windowWidth  int
	windowHeight int
	xIndexMax    int
	yIndexMax    int
}

func newHexConfig(windowWidth int, windowHeight int) hexConfig {
	lookUpTable := newHexLookUpTable()
	hexWidth := 46
	hexHeight := 40
	xIndexMax := windowWidth/hexWidth + 1
	yIndexMax := windowHeight/(hexHeight/4*3) + 1
	hc := hexConfig{
		lookUpTable: lookUpTable,
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

// The hexagon look up table is a hardcoded mapping to aid with converting
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

type hexagonTriangle int

const (
	unknownHexTriangle hexagonTriangle = iota
	leftHexTriangle
	middleHexTriangle
	rightHexTriangle
)

func newHexLookUpTable() [hexLookUpTablePixelCount]hexagonTriangle {
	var pc [hexLookUpTablePixelCount]hexagonTriangle

	rowPixelCount := 46

	// row 1
	rowStart := 0
	rowMiddleOffset := 1
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 2
	rowStart = rowPixelCount
	rowMiddleOffset = 3
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 3
	rowStart = rowPixelCount * 2
	rowMiddleOffset = 6
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 4
	rowStart = rowPixelCount * 3
	rowMiddleOffset = 8
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 5
	rowStart = rowPixelCount * 4
	rowMiddleOffset = 10
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 6
	rowStart = rowPixelCount * 5
	rowMiddleOffset = 13
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 7
	rowStart = rowPixelCount * 6
	rowMiddleOffset = 15
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 8
	rowStart = rowPixelCount * 7
	rowMiddleOffset = 17
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 9
	rowStart = rowPixelCount * 8
	rowMiddleOffset = 20
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)
	// row 10
	rowStart = rowPixelCount * 9
	rowMiddleOffset = 22
	setHexLookUpTable(&pc, rowPixelCount, rowStart, rowMiddleOffset)

	return pc
}

func setHexLookUpTable(pc *[hexLookUpTablePixelCount]hexagonTriangle, rowPixelCount int, rowStart int, rowMiddleOffset int) {
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
