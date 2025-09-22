package main

import (
	"testing"
)

// A very low level test. Ideally this would be higher level test for behavior rather than implementation.
// The higher level test would assert pixel positions -> grid hexagons rather than
// inner grid position -> hex triangle but that is not possible since the end to end
// behavior is not fully implemented.
func TestHexPositionContext(t *testing.T) {
	pc := newHexPositionContext()

	if len(pc) != 460 {
		t.Errorf("position context must represent 460 pixels (46 cols * 10 rows)")
	}

	// sanity check
	if pc[hexPosContextIndex(innerGridPosition{0, 0})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{45, 0})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{22, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{23, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 10 left side
	if pc[hexPosContextIndex(innerGridPosition{0, 9})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{1, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{2, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 9 left side
	if pc[hexPosContextIndex(innerGridPosition{2, 8})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{3, 8})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{4, 8})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 8 left side
	if pc[hexPosContextIndex(innerGridPosition{4, 7})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{5, 7})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{6, 7})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 7 left side
	if pc[hexPosContextIndex(innerGridPosition{6, 6})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{7, 6})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{8, 6})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{9, 6})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 6 left side
	if pc[hexPosContextIndex(innerGridPosition{9, 5})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{10, 5})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{11, 5})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 5 left side
	if pc[hexPosContextIndex(innerGridPosition{11, 4})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{12, 4})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{13, 4})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 4 left side
	if pc[hexPosContextIndex(innerGridPosition{13, 3})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{14, 3})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{15, 3})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{16, 3})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 3 left side
	if pc[hexPosContextIndex(innerGridPosition{16, 2})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{17, 2})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{18, 2})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 2 left side
	if pc[hexPosContextIndex(innerGridPosition{18, 1})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{19, 1})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{20, 1})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 1 left side
	if pc[hexPosContextIndex(innerGridPosition{20, 0})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{21, 0})] != leftHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{22, 0})] != middleHexTriangle {
		t.Errorf("error")
	}
	// row 1 right side
	if pc[hexPosContextIndex(innerGridPosition{23, 0})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{24, 0})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{25, 0})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 2 right side
	if pc[hexPosContextIndex(innerGridPosition{25, 1})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{26, 1})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{27, 1})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 3 right side
	if pc[hexPosContextIndex(innerGridPosition{27, 2})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{28, 2})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{29, 2})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 4 right side
	if pc[hexPosContextIndex(innerGridPosition{29, 3})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{30, 3})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{31, 3})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{32, 3})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 5 right side
	if pc[hexPosContextIndex(innerGridPosition{32, 4})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{33, 4})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{34, 4})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 6 right side
	if pc[hexPosContextIndex(innerGridPosition{34, 5})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{35, 5})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{36, 5})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 7 right side
	if pc[hexPosContextIndex(innerGridPosition{36, 6})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{37, 6})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{38, 6})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{39, 6})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 8 right side
	if pc[hexPosContextIndex(innerGridPosition{39, 7})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{40, 7})] != rightHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{41, 7})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 9 right side
	if pc[hexPosContextIndex(innerGridPosition{41, 8})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{42, 8})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{43, 8})] != rightHexTriangle {
		t.Errorf("error")
	}
	// row 10 right side
	if pc[hexPosContextIndex(innerGridPosition{43, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{44, 9})] != middleHexTriangle {
		t.Errorf("error")
	}
	if pc[hexPosContextIndex(innerGridPosition{45, 9})] != rightHexTriangle {
		t.Errorf("error")
	}
}
