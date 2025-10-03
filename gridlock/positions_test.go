package main

import (
	"fmt"
	"testing"
)

func TestGetPixelPos(t *testing.T) {
	lookUpTable := newHexLookUpTable()
	hc := hexConfig{
		lookUpTable:  lookUpTable,
		hexWidth:     46,
		hexHeight:    40,
		xOffset:      19,
		yOffset:      10,
		windowWidth:  1280,
		windowHeight: 720,
	}

	var tests = []struct {
		gridPos        gridPosition
		hexPoint       hexPoint
		wantedPixelPos pixelPosition
		wantedOk       bool
	}{

		// check hexagon {1, 1}, {1, 2}, {2,1} to ensure hex points are correct
		{gridPosition{1, 1}, center, pixelPosition{42, 30}, true},
		{gridPosition{1, 1}, northPoint, pixelPosition{42, 10}, true},
		{gridPosition{1, 1}, northeastPoint, pixelPosition{64, 20}, true},
		{gridPosition{1, 1}, southeastPoint, pixelPosition{64, 39}, true},
		{gridPosition{1, 1}, southPoint, pixelPosition{42, 49}, true},
		{gridPosition{1, 1}, southwestPoint, pixelPosition{19, 39}, true},
		{gridPosition{1, 1}, northwestPoint, pixelPosition{19, 20}, true},
		{gridPosition{2, 1}, center, pixelPosition{88, 30}, true},
		{gridPosition{2, 1}, northPoint, pixelPosition{88, 10}, true},
		{gridPosition{2, 1}, northeastPoint, pixelPosition{110, 20}, true},
		{gridPosition{2, 1}, southeastPoint, pixelPosition{110, 39}, true},
		{gridPosition{2, 1}, southPoint, pixelPosition{88, 49}, true},
		{gridPosition{2, 1}, southwestPoint, pixelPosition{65, 39}, true},
		{gridPosition{2, 1}, northwestPoint, pixelPosition{65, 20}, true},
		{gridPosition{1, 2}, center, pixelPosition{65, 60}, true},
		{gridPosition{1, 2}, northPoint, pixelPosition{65, 40}, true},
		{gridPosition{1, 2}, northeastPoint, pixelPosition{87, 50}, true},
		{gridPosition{1, 2}, southeastPoint, pixelPosition{87, 69}, true},
		{gridPosition{1, 2}, southPoint, pixelPosition{65, 79}, true},
		{gridPosition{1, 2}, southwestPoint, pixelPosition{42, 69}, true},
		{gridPosition{1, 2}, northwestPoint, pixelPosition{42, 50}, true},

		// check hexagon {26,23} to ensure large x and large y positions are correct
		{gridPosition{27, 23}, center, pixelPosition{1238, 690}, true},
		{gridPosition{27, 23}, northPoint, pixelPosition{1238, 670}, true},
		{gridPosition{27, 23}, northeastPoint, pixelPosition{1260, 680}, true},
		{gridPosition{27, 23}, southeastPoint, pixelPosition{1260, 699}, true},
		{gridPosition{27, 23}, southPoint, pixelPosition{1238, 709}, true},
		{gridPosition{27, 23}, southwestPoint, pixelPosition{1215, 699}, true},
		{gridPosition{27, 23}, northwestPoint, pixelPosition{1215, 680}, true},

		// check top left hexagon {0,0}
		{gridPosition{0, 0}, center, pixelPosition{19, 0}, true},
		{gridPosition{0, 0}, northPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 0}, northeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 0}, southeastPoint, pixelPosition{41, 9}, true},
		{gridPosition{0, 0}, southPoint, pixelPosition{19, 19}, true},
		{gridPosition{0, 0}, southwestPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 0}, northwestPoint, pixelPosition{0, 0}, false},

		// check bottom left hexagon {0,24}
		{gridPosition{0, 24}, center, pixelPosition{0, 0}, false},
		{gridPosition{0, 24}, northPoint, pixelPosition{19, 700}, true},
		{gridPosition{0, 24}, northeastPoint, pixelPosition{41, 710}, true},
		{gridPosition{0, 24}, southeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 24}, southPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 24}, southwestPoint, pixelPosition{0, 0}, false},
		{gridPosition{0, 24}, northwestPoint, pixelPosition{0, 0}, false},

		// check top right hexagon {27,0}
		{gridPosition{27, 0}, center, pixelPosition{1261, 0}, true},
		{gridPosition{27, 0}, northPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 0}, northeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 0}, southeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 0}, southPoint, pixelPosition{1261, 19}, true},
		{gridPosition{27, 0}, southwestPoint, pixelPosition{1238, 9}, true},
		{gridPosition{27, 0}, northwestPoint, pixelPosition{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{gridPosition{27, 24}, center, pixelPosition{0, 0}, false},
		{gridPosition{27, 24}, northPoint, pixelPosition{1261, 700}, true},
		{gridPosition{27, 24}, northeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southeastPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southwestPoint, pixelPosition{0, 0}, false},
		{gridPosition{27, 24}, northwestPoint, pixelPosition{1238, 710}, true},
	}

	for _, tt := range tests {
		var sWantedOk string
		if tt.wantedOk {
			sWantedOk = "exist"
		} else {
			sWantedOk = "not exist"
		}

		testname := fmt.Sprintf(
			"gridPos:{%d,%d} and hex point %d converts to hex grid position {%d,%d} and want it to %s",
			tt.gridPos.x,
			tt.gridPos.y,
			tt.hexPoint,
			tt.wantedPixelPos.x,
			tt.wantedPixelPos.y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := getPixelPos(tt.gridPos, tt.hexPoint, &hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if pos.x != tt.wantedPixelPos.x || pos.y != tt.wantedPixelPos.y {
				t.Errorf("For grid pos {%d,%d} and hex point %d  Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.gridPos.x,
					tt.gridPos.y,
					tt.hexPoint,
					tt.wantedPixelPos.x,
					tt.wantedPixelPos.y,
					pos.x,
					pos.y)
			}
		})
	}
}

func TestGetGridPos(t *testing.T) {
	lookUpTable := newHexLookUpTable()
	hc := hexConfig{
		lookUpTable:  lookUpTable,
		hexWidth:     46,
		hexHeight:    40,
		xOffset:      19,
		yOffset:      10,
		windowWidth:  1280,
		windowHeight: 720,
	}

	var tests = []struct {
		pixelPos  pixelPosition
		wantedPos gridPosition
		wantedOk  bool
	}{
		// top left of entire hex grid
		{pixelPosition{0, 0}, gridPosition{0, 0}, true},
		{pixelPosition{-1, 0}, gridPosition{0, 0}, false},
		{pixelPosition{0, -1}, gridPosition{0, 0}, false},
		{pixelPosition{-1, -1}, gridPosition{0, 0}, false},

		// first left border point hex change around 0,12 pixel position
		{pixelPosition{-1, 11}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 12}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 13}, gridPosition{0, 0}, false},
		{pixelPosition{0, 11}, gridPosition{0, 0}, true},
		{pixelPosition{0, 12}, gridPosition{0, 1}, true},
		{pixelPosition{0, 13}, gridPosition{0, 1}, true},
		{pixelPosition{1, 11}, gridPosition{0, 0}, true},
		{pixelPosition{1, 12}, gridPosition{0, 1}, true},
		{pixelPosition{1, 13}, gridPosition{0, 1}, true},

		// second left border point hex change around 0,48 pixel position
		{pixelPosition{-1, 47}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 48}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 49}, gridPosition{0, 0}, false},
		{pixelPosition{0, 47}, gridPosition{0, 1}, true},
		{pixelPosition{0, 48}, gridPosition{0, 2}, true},
		{pixelPosition{0, 49}, gridPosition{0, 2}, true},
		{pixelPosition{1, 47}, gridPosition{0, 1}, true},
		{pixelPosition{1, 48}, gridPosition{0, 2}, true},
		{pixelPosition{1, 49}, gridPosition{0, 2}, true},

		// third left border point hex change around 0,72 pixel position
		{pixelPosition{-1, 71}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 72}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 73}, gridPosition{0, 0}, false},
		{pixelPosition{0, 71}, gridPosition{0, 2}, true},
		{pixelPosition{0, 72}, gridPosition{0, 3}, true},
		{pixelPosition{0, 73}, gridPosition{0, 3}, true},
		{pixelPosition{1, 71}, gridPosition{0, 2}, true},
		{pixelPosition{1, 72}, gridPosition{0, 3}, true},
		{pixelPosition{1, 73}, gridPosition{0, 3}, true},

		// first top border point hex change around 42,0 pixel position
		{pixelPosition{41, -1}, gridPosition{0, 0}, false},
		{pixelPosition{41, 0}, gridPosition{0, 0}, true},
		{pixelPosition{41, 1}, gridPosition{0, 0}, true},
		{pixelPosition{42, -1}, gridPosition{0, 0}, false},
		{pixelPosition{42, 0}, gridPosition{1, 0}, true},
		{pixelPosition{42, 1}, gridPosition{1, 0}, true},
		{pixelPosition{43, -1}, gridPosition{0, 0}, false},
		{pixelPosition{43, 0}, gridPosition{1, 0}, true},
		{pixelPosition{43, 1}, gridPosition{1, 0}, true},

		// second top border point hex change around 88, 0 pixel position
		{pixelPosition{87, -1}, gridPosition{0, 0}, false},
		{pixelPosition{87, 0}, gridPosition{1, 0}, true},
		{pixelPosition{87, 1}, gridPosition{1, 0}, true},
		{pixelPosition{88, -1}, gridPosition{0, 0}, false},
		{pixelPosition{88, 0}, gridPosition{2, 0}, true},
		{pixelPosition{88, 1}, gridPosition{2, 0}, true},
		{pixelPosition{89, -1}, gridPosition{0, 0}, false},
		{pixelPosition{89, 0}, gridPosition{2, 0}, true},
		{pixelPosition{89, 1}, gridPosition{2, 0}, true},

		// second top border point hex change around 134,0 pixel position
		{pixelPosition{133, -1}, gridPosition{0, 0}, false},
		{pixelPosition{133, 0}, gridPosition{2, 0}, true},
		{pixelPosition{133, 1}, gridPosition{2, 0}, true},
		{pixelPosition{134, -1}, gridPosition{0, 0}, false},
		{pixelPosition{134, 0}, gridPosition{3, 0}, true},
		{pixelPosition{134, 1}, gridPosition{3, 0}, true},
		{pixelPosition{135, -1}, gridPosition{0, 0}, false},
		{pixelPosition{135, 0}, gridPosition{3, 0}, true},
		{pixelPosition{135, 1}, gridPosition{3, 0}, true},

		// hexagon points near 19, 80 pixel position
		{pixelPosition{18, 79}, gridPosition{0, 2}, true},
		{pixelPosition{18, 80}, gridPosition{0, 3}, true},
		{pixelPosition{18, 81}, gridPosition{0, 3}, true},
		{pixelPosition{19, 79}, gridPosition{0, 2}, true},
		{pixelPosition{19, 80}, gridPosition{1, 3}, true},
		{pixelPosition{19, 81}, gridPosition{1, 3}, true},
		{pixelPosition{20, 79}, gridPosition{1, 3}, true},
		{pixelPosition{20, 80}, gridPosition{1, 3}, true},
		{pixelPosition{20, 81}, gridPosition{1, 3}, true},

		// hexagon points near 42, 70 pixel position
		{pixelPosition{41, 69}, gridPosition{0, 2}, true},
		{pixelPosition{41, 70}, gridPosition{1, 3}, true},
		{pixelPosition{41, 71}, gridPosition{1, 3}, true},
		{pixelPosition{42, 69}, gridPosition{1, 2}, true},
		{pixelPosition{42, 70}, gridPosition{1, 3}, true},
		{pixelPosition{42, 71}, gridPosition{1, 3}, true},
		{pixelPosition{43, 69}, gridPosition{1, 2}, true},
		{pixelPosition{43, 70}, gridPosition{1, 2}, true},
		{pixelPosition{43, 71}, gridPosition{1, 3}, true},

		// hexagon points near 65, 80 pixel position
		{pixelPosition{64, 79}, gridPosition{1, 2}, true},
		{pixelPosition{64, 80}, gridPosition{1, 3}, true},
		{pixelPosition{64, 81}, gridPosition{1, 3}, true},
		{pixelPosition{65, 79}, gridPosition{1, 2}, true},
		{pixelPosition{65, 80}, gridPosition{2, 3}, true},
		{pixelPosition{65, 81}, gridPosition{2, 3}, true},
		{pixelPosition{66, 79}, gridPosition{2, 3}, true},
		{pixelPosition{66, 80}, gridPosition{2, 3}, true},
		{pixelPosition{66, 81}, gridPosition{2, 3}, true},

		// hexagon points near 88, 10 pixel position
		{pixelPosition{87, 9}, gridPosition{1, 0}, true},
		{pixelPosition{87, 10}, gridPosition{2, 1}, true},
		{pixelPosition{87, 11}, gridPosition{2, 1}, true},
		{pixelPosition{88, 9}, gridPosition{2, 0}, true},
		{pixelPosition{88, 10}, gridPosition{2, 1}, true},
		{pixelPosition{88, 11}, gridPosition{2, 1}, true},
		{pixelPosition{89, 9}, gridPosition{2, 0}, true},
		{pixelPosition{89, 10}, gridPosition{2, 0}, true},
		{pixelPosition{89, 11}, gridPosition{2, 1}, true},

		// hexagon points near 88, 50 pixel position
		{pixelPosition{87, 49}, gridPosition{2, 1}, true},
		{pixelPosition{87, 50}, gridPosition{1, 2}, true},
		{pixelPosition{87, 51}, gridPosition{1, 2}, true},
		{pixelPosition{88, 49}, gridPosition{2, 1}, true},
		{pixelPosition{88, 50}, gridPosition{2, 2}, true},
		{pixelPosition{88, 51}, gridPosition{2, 2}, true},
		{pixelPosition{89, 49}, gridPosition{2, 2}, true},
		{pixelPosition{89, 50}, gridPosition{2, 2}, true},
		{pixelPosition{89, 51}, gridPosition{2, 2}, true},

		// hexagon points near 88, 70 pixel position
		{pixelPosition{87, 69}, gridPosition{1, 2}, true},
		{pixelPosition{87, 70}, gridPosition{2, 3}, true},
		{pixelPosition{87, 71}, gridPosition{2, 3}, true},
		{pixelPosition{88, 69}, gridPosition{2, 2}, true},
		{pixelPosition{88, 70}, gridPosition{2, 3}, true},
		{pixelPosition{88, 71}, gridPosition{2, 3}, true},
		{pixelPosition{89, 69}, gridPosition{2, 2}, true},
		{pixelPosition{89, 70}, gridPosition{2, 2}, true},
		{pixelPosition{89, 71}, gridPosition{2, 3}, true},

		// bottom left of entire hex grid
		{pixelPosition{-1, 719}, gridPosition{0, 0}, false},
		{pixelPosition{-1, 720}, gridPosition{0, 0}, false},
		{pixelPosition{0, 719}, gridPosition{0, 24}, true},
		{pixelPosition{0, 720}, gridPosition{0, 0}, false},

		// top right of entire hex grid
		{pixelPosition{1279, -1}, gridPosition{0, 0}, false},
		{pixelPosition{1280, -1}, gridPosition{0, 0}, false},
		{pixelPosition{1279, 0}, gridPosition{27, 0}, true},
		{pixelPosition{1280, 0}, gridPosition{0, 0}, false},

		// bottom right of entire hex grid
		{pixelPosition{1279, 719}, gridPosition{27, 24}, true},
		{pixelPosition{1280, 719}, gridPosition{0, 0}, false},
		{pixelPosition{1279, 720}, gridPosition{0, 0}, false},
		{pixelPosition{1280, 720}, gridPosition{0, 0}, false},

		// first hexagon y section (hexagon look up table)
		// between hexagons 0,0 1,1 and 1,0

		// sanity checks
		{pixelPosition{19, 10}, gridPosition{0, 0}, true},
		{pixelPosition{65, 10}, gridPosition{1, 0}, true},
		{pixelPosition{41, 20}, gridPosition{1, 1}, true},
		{pixelPosition{42, 20}, gridPosition{1, 1}, true},
		{pixelPosition{43, 20}, gridPosition{1, 1}, true},

		// 0,0 hexagon slope check between itself and 1, 1 hexagon
		{pixelPosition{19, 19}, gridPosition{0, 0}, true},
		{pixelPosition{20, 18}, gridPosition{0, 0}, true},
		{pixelPosition{21, 18}, gridPosition{0, 0}, true},
		{pixelPosition{22, 17}, gridPosition{0, 0}, true},
		{pixelPosition{23, 17}, gridPosition{0, 0}, true},
		{pixelPosition{24, 17}, gridPosition{0, 0}, true},
		{pixelPosition{25, 16}, gridPosition{0, 0}, true},
		{pixelPosition{26, 16}, gridPosition{0, 0}, true},
		{pixelPosition{27, 15}, gridPosition{0, 0}, true},
		{pixelPosition{28, 15}, gridPosition{0, 0}, true},
		{pixelPosition{29, 14}, gridPosition{0, 0}, true},
		{pixelPosition{30, 14}, gridPosition{0, 0}, true},
		{pixelPosition{31, 14}, gridPosition{0, 0}, true},
		{pixelPosition{32, 13}, gridPosition{0, 0}, true},
		{pixelPosition{33, 13}, gridPosition{0, 0}, true},
		{pixelPosition{34, 12}, gridPosition{0, 0}, true},
		{pixelPosition{35, 12}, gridPosition{0, 0}, true},
		{pixelPosition{36, 11}, gridPosition{0, 0}, true},
		{pixelPosition{37, 11}, gridPosition{0, 0}, true},
		{pixelPosition{38, 11}, gridPosition{0, 0}, true},
		{pixelPosition{39, 10}, gridPosition{0, 0}, true},
		{pixelPosition{40, 10}, gridPosition{0, 0}, true},
		//{pixelPosition{41, 9}, gridPosition{0, 0}, true}, // prob bug

		// 1, 1 hexagon slope check between itself and 0, 0 hexagon
		{pixelPosition{19, 20}, gridPosition{1, 1}, true},
		{pixelPosition{20, 19}, gridPosition{1, 1}, true},
		{pixelPosition{21, 19}, gridPosition{1, 1}, true},
		{pixelPosition{22, 18}, gridPosition{1, 1}, true},
		{pixelPosition{23, 18}, gridPosition{1, 1}, true},
		{pixelPosition{24, 18}, gridPosition{1, 1}, true},
		{pixelPosition{25, 17}, gridPosition{1, 1}, true},
		{pixelPosition{26, 17}, gridPosition{1, 1}, true},
		{pixelPosition{27, 16}, gridPosition{1, 1}, true},
		{pixelPosition{28, 16}, gridPosition{1, 1}, true},
		{pixelPosition{29, 15}, gridPosition{1, 1}, true},
		{pixelPosition{30, 15}, gridPosition{1, 1}, true},
		{pixelPosition{31, 15}, gridPosition{1, 1}, true},
		{pixelPosition{32, 14}, gridPosition{1, 1}, true},
		{pixelPosition{33, 14}, gridPosition{1, 1}, true},
		{pixelPosition{34, 13}, gridPosition{1, 1}, true},
		{pixelPosition{35, 13}, gridPosition{1, 1}, true},
		{pixelPosition{36, 12}, gridPosition{1, 1}, true},
		{pixelPosition{37, 12}, gridPosition{1, 1}, true},
		{pixelPosition{38, 12}, gridPosition{1, 1}, true},
		{pixelPosition{39, 11}, gridPosition{1, 1}, true},
		{pixelPosition{40, 11}, gridPosition{1, 1}, true},
		{pixelPosition{41, 10}, gridPosition{1, 1}, true},

		// 1,0 hexagon slope check between itself and 1, 1 hexagon
		//{pixelPosition{42, 9}, gridPosition{1, 0}, true}, prob bug
		{pixelPosition{43, 10}, gridPosition{1, 0}, true},
		{pixelPosition{44, 10}, gridPosition{1, 0}, true},
		{pixelPosition{45, 11}, gridPosition{1, 0}, true},
		{pixelPosition{46, 11}, gridPosition{1, 0}, true},
		{pixelPosition{47, 11}, gridPosition{1, 0}, true},
		{pixelPosition{48, 12}, gridPosition{1, 0}, true},
		{pixelPosition{49, 12}, gridPosition{1, 0}, true},
		{pixelPosition{50, 13}, gridPosition{1, 0}, true},
		{pixelPosition{51, 13}, gridPosition{1, 0}, true},

		{pixelPosition{52, 14}, gridPosition{1, 0}, true},
		{pixelPosition{53, 14}, gridPosition{1, 0}, true},
		{pixelPosition{54, 14}, gridPosition{1, 0}, true},
		{pixelPosition{55, 15}, gridPosition{1, 0}, true},
		{pixelPosition{56, 15}, gridPosition{1, 0}, true},
		{pixelPosition{57, 16}, gridPosition{1, 0}, true},
		{pixelPosition{58, 16}, gridPosition{1, 0}, true},
		{pixelPosition{59, 17}, gridPosition{1, 0}, true},
		{pixelPosition{60, 17}, gridPosition{1, 0}, true},
		{pixelPosition{61, 17}, gridPosition{1, 0}, true},
		{pixelPosition{62, 18}, gridPosition{1, 0}, true},
		{pixelPosition{63, 18}, gridPosition{1, 0}, true},
		{pixelPosition{64, 19}, gridPosition{1, 0}, true},

		// 1,1 hexagon slope check between itself and 1, 0 hexagon
		{pixelPosition{42, 10}, gridPosition{1, 1}, true},
		{pixelPosition{43, 11}, gridPosition{1, 1}, true},
		{pixelPosition{44, 11}, gridPosition{1, 1}, true},
		{pixelPosition{45, 12}, gridPosition{1, 1}, true},
		{pixelPosition{46, 12}, gridPosition{1, 1}, true},
		{pixelPosition{47, 12}, gridPosition{1, 1}, true},
		{pixelPosition{48, 13}, gridPosition{1, 1}, true},
		{pixelPosition{49, 13}, gridPosition{1, 1}, true},
		{pixelPosition{50, 14}, gridPosition{1, 1}, true},
		{pixelPosition{51, 14}, gridPosition{1, 1}, true},
		{pixelPosition{52, 15}, gridPosition{1, 1}, true},
		{pixelPosition{53, 15}, gridPosition{1, 1}, true},
		{pixelPosition{54, 15}, gridPosition{1, 1}, true},
		{pixelPosition{55, 16}, gridPosition{1, 1}, true},
		{pixelPosition{56, 16}, gridPosition{1, 1}, true},
		{pixelPosition{57, 17}, gridPosition{1, 1}, true},
		{pixelPosition{58, 17}, gridPosition{1, 1}, true},
		{pixelPosition{59, 18}, gridPosition{1, 1}, true},
		{pixelPosition{60, 18}, gridPosition{1, 1}, true},
		{pixelPosition{61, 18}, gridPosition{1, 1}, true},
		{pixelPosition{62, 19}, gridPosition{1, 1}, true},
		{pixelPosition{63, 19}, gridPosition{1, 1}, true},
		{pixelPosition{64, 20}, gridPosition{1, 1}, true},

		// second hexagon y section (full rectangle of hexagon 1,1)

		// right side of hexagon 0,1
		{pixelPosition{18, 20}, gridPosition{0, 1}, true},
		{pixelPosition{18, 21}, gridPosition{0, 1}, true},
		{pixelPosition{18, 22}, gridPosition{0, 1}, true},
		{pixelPosition{18, 23}, gridPosition{0, 1}, true},
		{pixelPosition{18, 24}, gridPosition{0, 1}, true},
		{pixelPosition{18, 25}, gridPosition{0, 1}, true},
		{pixelPosition{18, 26}, gridPosition{0, 1}, true},
		{pixelPosition{18, 27}, gridPosition{0, 1}, true},
		{pixelPosition{18, 28}, gridPosition{0, 1}, true},
		{pixelPosition{18, 29}, gridPosition{0, 1}, true},
		{pixelPosition{18, 30}, gridPosition{0, 1}, true},
		{pixelPosition{18, 31}, gridPosition{0, 1}, true},
		{pixelPosition{18, 32}, gridPosition{0, 1}, true},
		{pixelPosition{18, 33}, gridPosition{0, 1}, true},
		{pixelPosition{18, 34}, gridPosition{0, 1}, true},
		{pixelPosition{18, 35}, gridPosition{0, 1}, true},
		{pixelPosition{18, 36}, gridPosition{0, 1}, true},
		{pixelPosition{18, 37}, gridPosition{0, 1}, true},
		{pixelPosition{18, 38}, gridPosition{0, 1}, true},
		{pixelPosition{18, 39}, gridPosition{0, 1}, true},

		// left side of hexagon 1,1
		// hexagons are right side inclusive and left side exclusive
		{pixelPosition{19, 20}, gridPosition{1, 1}, true},
		{pixelPosition{19, 21}, gridPosition{1, 1}, true},
		{pixelPosition{19, 22}, gridPosition{1, 1}, true},
		{pixelPosition{19, 23}, gridPosition{1, 1}, true},
		{pixelPosition{19, 24}, gridPosition{1, 1}, true},
		{pixelPosition{19, 25}, gridPosition{1, 1}, true},
		{pixelPosition{19, 26}, gridPosition{1, 1}, true},
		{pixelPosition{19, 27}, gridPosition{1, 1}, true},
		{pixelPosition{19, 28}, gridPosition{1, 1}, true},
		{pixelPosition{19, 29}, gridPosition{1, 1}, true},
		{pixelPosition{19, 30}, gridPosition{1, 1}, true},
		{pixelPosition{19, 31}, gridPosition{1, 1}, true},
		{pixelPosition{19, 32}, gridPosition{1, 1}, true},
		{pixelPosition{19, 33}, gridPosition{1, 1}, true},
		{pixelPosition{19, 34}, gridPosition{1, 1}, true},
		{pixelPosition{19, 35}, gridPosition{1, 1}, true},
		{pixelPosition{19, 36}, gridPosition{1, 1}, true},
		{pixelPosition{19, 37}, gridPosition{1, 1}, true},
		{pixelPosition{19, 38}, gridPosition{1, 1}, true},
		{pixelPosition{19, 39}, gridPosition{1, 1}, true},
		// continue left side of hexagon 1,1
		{pixelPosition{20, 20}, gridPosition{1, 1}, true},
		{pixelPosition{20, 21}, gridPosition{1, 1}, true},
		{pixelPosition{20, 22}, gridPosition{1, 1}, true},
		{pixelPosition{20, 23}, gridPosition{1, 1}, true},
		{pixelPosition{20, 24}, gridPosition{1, 1}, true},
		{pixelPosition{20, 25}, gridPosition{1, 1}, true},
		{pixelPosition{20, 26}, gridPosition{1, 1}, true},
		{pixelPosition{20, 27}, gridPosition{1, 1}, true},
		{pixelPosition{20, 28}, gridPosition{1, 1}, true},
		{pixelPosition{20, 29}, gridPosition{1, 1}, true},
		{pixelPosition{20, 30}, gridPosition{1, 1}, true},
		{pixelPosition{20, 31}, gridPosition{1, 1}, true},
		{pixelPosition{20, 32}, gridPosition{1, 1}, true},
		{pixelPosition{20, 33}, gridPosition{1, 1}, true},
		{pixelPosition{20, 34}, gridPosition{1, 1}, true},
		{pixelPosition{20, 35}, gridPosition{1, 1}, true},
		{pixelPosition{20, 36}, gridPosition{1, 1}, true},
		{pixelPosition{20, 37}, gridPosition{1, 1}, true},
		{pixelPosition{20, 38}, gridPosition{1, 1}, true},
		{pixelPosition{20, 39}, gridPosition{1, 1}, true},

		// right side of hexagon 1,1
		{pixelPosition{64, 20}, gridPosition{1, 1}, true},
		{pixelPosition{64, 21}, gridPosition{1, 1}, true},
		{pixelPosition{64, 22}, gridPosition{1, 1}, true},
		{pixelPosition{64, 23}, gridPosition{1, 1}, true},
		{pixelPosition{64, 24}, gridPosition{1, 1}, true},
		{pixelPosition{64, 25}, gridPosition{1, 1}, true},
		{pixelPosition{64, 26}, gridPosition{1, 1}, true},
		{pixelPosition{64, 27}, gridPosition{1, 1}, true},
		{pixelPosition{64, 28}, gridPosition{1, 1}, true},
		{pixelPosition{64, 29}, gridPosition{1, 1}, true},
		{pixelPosition{64, 30}, gridPosition{1, 1}, true},
		{pixelPosition{64, 31}, gridPosition{1, 1}, true},
		{pixelPosition{64, 32}, gridPosition{1, 1}, true},
		{pixelPosition{64, 33}, gridPosition{1, 1}, true},
		{pixelPosition{64, 34}, gridPosition{1, 1}, true},
		{pixelPosition{64, 35}, gridPosition{1, 1}, true},
		{pixelPosition{64, 36}, gridPosition{1, 1}, true},
		{pixelPosition{64, 37}, gridPosition{1, 1}, true},
		{pixelPosition{64, 38}, gridPosition{1, 1}, true},
		{pixelPosition{64, 39}, gridPosition{1, 1}, true},

		// left side of hexagon 2,1
		// hexagons are right side inclusive and left side exclusive
		{pixelPosition{65, 20}, gridPosition{2, 1}, true},
		{pixelPosition{65, 21}, gridPosition{2, 1}, true},
		{pixelPosition{65, 22}, gridPosition{2, 1}, true},
		{pixelPosition{65, 23}, gridPosition{2, 1}, true},
		{pixelPosition{65, 24}, gridPosition{2, 1}, true},
		{pixelPosition{65, 25}, gridPosition{2, 1}, true},
		{pixelPosition{65, 26}, gridPosition{2, 1}, true},
		{pixelPosition{65, 27}, gridPosition{2, 1}, true},
		{pixelPosition{65, 28}, gridPosition{2, 1}, true},
		{pixelPosition{65, 29}, gridPosition{2, 1}, true},
		{pixelPosition{65, 30}, gridPosition{2, 1}, true},
		{pixelPosition{65, 31}, gridPosition{2, 1}, true},
		{pixelPosition{65, 32}, gridPosition{2, 1}, true},
		{pixelPosition{65, 33}, gridPosition{2, 1}, true},
		{pixelPosition{65, 34}, gridPosition{2, 1}, true},
		{pixelPosition{65, 35}, gridPosition{2, 1}, true},
		{pixelPosition{65, 36}, gridPosition{2, 1}, true},
		{pixelPosition{65, 37}, gridPosition{2, 1}, true},
		{pixelPosition{65, 38}, gridPosition{2, 1}, true},
		{pixelPosition{65, 39}, gridPosition{2, 1}, true},
		// continue left side of hexagon 2,1
		{pixelPosition{66, 20}, gridPosition{2, 1}, true},
		{pixelPosition{66, 21}, gridPosition{2, 1}, true},
		{pixelPosition{66, 22}, gridPosition{2, 1}, true},
		{pixelPosition{66, 23}, gridPosition{2, 1}, true},
		{pixelPosition{66, 24}, gridPosition{2, 1}, true},
		{pixelPosition{66, 25}, gridPosition{2, 1}, true},
		{pixelPosition{66, 26}, gridPosition{2, 1}, true},
		{pixelPosition{66, 27}, gridPosition{2, 1}, true},
		{pixelPosition{66, 28}, gridPosition{2, 1}, true},
		{pixelPosition{66, 29}, gridPosition{2, 1}, true},
		{pixelPosition{66, 30}, gridPosition{2, 1}, true},
		{pixelPosition{66, 31}, gridPosition{2, 1}, true},
		{pixelPosition{66, 32}, gridPosition{2, 1}, true},
		{pixelPosition{66, 33}, gridPosition{2, 1}, true},
		{pixelPosition{66, 34}, gridPosition{2, 1}, true},
		{pixelPosition{66, 35}, gridPosition{2, 1}, true},
		{pixelPosition{66, 36}, gridPosition{2, 1}, true},
		{pixelPosition{66, 37}, gridPosition{2, 1}, true},
		{pixelPosition{66, 38}, gridPosition{2, 1}, true},
		{pixelPosition{66, 39}, gridPosition{2, 1}, true},

		// third hexagon y section (inverted hexagon look up table)
		// between hexagons 0,2 1,1 and 1,2

		// sanity checks
		{pixelPosition{19, 49}, gridPosition{0, 2}, true},
		{pixelPosition{65, 49}, gridPosition{1, 2}, true},
		{pixelPosition{41, 30}, gridPosition{1, 1}, true},
		{pixelPosition{42, 30}, gridPosition{1, 1}, true},
		{pixelPosition{43, 30}, gridPosition{1, 1}, true},

		// 0,2 hexagon slope check between itself and 1, 1 hexagon
		{pixelPosition{19, 40}, gridPosition{0, 2}, true},
		{pixelPosition{20, 41}, gridPosition{0, 2}, true},
		{pixelPosition{21, 41}, gridPosition{0, 2}, true},
		{pixelPosition{22, 42}, gridPosition{0, 2}, true},
		{pixelPosition{23, 42}, gridPosition{0, 2}, true},
		{pixelPosition{24, 42}, gridPosition{0, 2}, true},
		{pixelPosition{25, 43}, gridPosition{0, 2}, true},
		{pixelPosition{26, 43}, gridPosition{0, 2}, true},
		{pixelPosition{27, 44}, gridPosition{0, 2}, true},
		{pixelPosition{28, 44}, gridPosition{0, 2}, true},
		{pixelPosition{29, 45}, gridPosition{0, 2}, true},
		{pixelPosition{30, 45}, gridPosition{0, 2}, true},
		{pixelPosition{31, 45}, gridPosition{0, 2}, true},
		{pixelPosition{32, 46}, gridPosition{0, 2}, true},
		{pixelPosition{33, 46}, gridPosition{0, 2}, true},
		{pixelPosition{34, 47}, gridPosition{0, 2}, true},
		{pixelPosition{35, 47}, gridPosition{0, 2}, true},
		{pixelPosition{36, 48}, gridPosition{0, 2}, true},
		{pixelPosition{37, 48}, gridPosition{0, 2}, true},
		{pixelPosition{38, 48}, gridPosition{0, 2}, true},
		{pixelPosition{39, 49}, gridPosition{0, 2}, true},
		{pixelPosition{40, 49}, gridPosition{0, 2}, true},
		{pixelPosition{41, 50}, gridPosition{0, 2}, true},

		// 1, 1 hexagon slope check between itself and 0, 2 hexagon
		{pixelPosition{19, 39}, gridPosition{1, 1}, true},
		{pixelPosition{20, 40}, gridPosition{1, 1}, true},
		{pixelPosition{21, 40}, gridPosition{1, 1}, true},
		{pixelPosition{22, 41}, gridPosition{1, 1}, true},
		{pixelPosition{23, 41}, gridPosition{1, 1}, true},
		{pixelPosition{24, 41}, gridPosition{1, 1}, true},
		{pixelPosition{25, 42}, gridPosition{1, 1}, true},
		{pixelPosition{26, 42}, gridPosition{1, 1}, true},
		{pixelPosition{27, 43}, gridPosition{1, 1}, true},
		{pixelPosition{28, 43}, gridPosition{1, 1}, true},
		{pixelPosition{29, 44}, gridPosition{1, 1}, true},
		{pixelPosition{30, 44}, gridPosition{1, 1}, true},
		{pixelPosition{31, 44}, gridPosition{1, 1}, true},
		{pixelPosition{32, 45}, gridPosition{1, 1}, true},
		{pixelPosition{33, 45}, gridPosition{1, 1}, true},
		{pixelPosition{34, 46}, gridPosition{1, 1}, true},
		{pixelPosition{35, 46}, gridPosition{1, 1}, true},
		{pixelPosition{36, 47}, gridPosition{1, 1}, true},
		{pixelPosition{37, 47}, gridPosition{1, 1}, true},
		{pixelPosition{38, 47}, gridPosition{1, 1}, true},
		{pixelPosition{39, 48}, gridPosition{1, 1}, true},
		{pixelPosition{40, 48}, gridPosition{1, 1}, true},
		{pixelPosition{41, 49}, gridPosition{1, 1}, true},

		// 1,1 hexagon slope check between itself and 1, 2 hexagon
		{pixelPosition{42, 49}, gridPosition{1, 1}, true},
		{pixelPosition{43, 48}, gridPosition{1, 1}, true},
		{pixelPosition{44, 48}, gridPosition{1, 1}, true},
		{pixelPosition{45, 47}, gridPosition{1, 1}, true},
		{pixelPosition{46, 47}, gridPosition{1, 1}, true},
		{pixelPosition{47, 47}, gridPosition{1, 1}, true},
		{pixelPosition{48, 46}, gridPosition{1, 1}, true},
		{pixelPosition{49, 46}, gridPosition{1, 1}, true},
		{pixelPosition{50, 45}, gridPosition{1, 1}, true},
		{pixelPosition{51, 45}, gridPosition{1, 1}, true},
		{pixelPosition{52, 44}, gridPosition{1, 1}, true},
		{pixelPosition{53, 44}, gridPosition{1, 1}, true},
		{pixelPosition{54, 44}, gridPosition{1, 1}, true},
		{pixelPosition{55, 43}, gridPosition{1, 1}, true},
		{pixelPosition{56, 43}, gridPosition{1, 1}, true},
		{pixelPosition{57, 42}, gridPosition{1, 1}, true},
		{pixelPosition{58, 42}, gridPosition{1, 1}, true},
		{pixelPosition{59, 41}, gridPosition{1, 1}, true},
		{pixelPosition{60, 41}, gridPosition{1, 1}, true},
		{pixelPosition{61, 41}, gridPosition{1, 1}, true},
		{pixelPosition{62, 40}, gridPosition{1, 1}, true},
		{pixelPosition{63, 40}, gridPosition{1, 1}, true},
		{pixelPosition{64, 39}, gridPosition{1, 1}, true},

		// 1,2 hexagon slope check between itself and 1, 1 hexagon
		{pixelPosition{42, 50}, gridPosition{1, 2}, true},
		{pixelPosition{43, 49}, gridPosition{1, 2}, true},
		{pixelPosition{44, 49}, gridPosition{1, 2}, true},
		{pixelPosition{45, 48}, gridPosition{1, 2}, true},
		{pixelPosition{46, 48}, gridPosition{1, 2}, true},
		{pixelPosition{47, 48}, gridPosition{1, 2}, true},
		{pixelPosition{48, 47}, gridPosition{1, 2}, true},
		{pixelPosition{49, 47}, gridPosition{1, 2}, true},
		{pixelPosition{50, 46}, gridPosition{1, 2}, true},
		{pixelPosition{51, 46}, gridPosition{1, 2}, true},
		{pixelPosition{52, 45}, gridPosition{1, 2}, true},
		{pixelPosition{53, 45}, gridPosition{1, 2}, true},
		{pixelPosition{54, 45}, gridPosition{1, 2}, true},
		{pixelPosition{55, 44}, gridPosition{1, 2}, true},
		{pixelPosition{56, 44}, gridPosition{1, 2}, true},
		{pixelPosition{57, 43}, gridPosition{1, 2}, true},
		{pixelPosition{58, 43}, gridPosition{1, 2}, true},
		{pixelPosition{59, 42}, gridPosition{1, 2}, true},
		{pixelPosition{60, 42}, gridPosition{1, 2}, true},
		{pixelPosition{61, 42}, gridPosition{1, 2}, true},
		{pixelPosition{62, 41}, gridPosition{1, 2}, true},
		{pixelPosition{63, 41}, gridPosition{1, 2}, true},
		{pixelPosition{64, 40}, gridPosition{1, 2}, true},

		// fourth hexagon y section (split rectangle between hexagon 0,2 and 1,2)

		// middle of hexagon 0,2
		{pixelPosition{18, 50}, gridPosition{0, 2}, true},
		{pixelPosition{18, 51}, gridPosition{0, 2}, true},
		{pixelPosition{18, 52}, gridPosition{0, 2}, true},
		{pixelPosition{18, 53}, gridPosition{0, 2}, true},
		{pixelPosition{18, 54}, gridPosition{0, 2}, true},
		{pixelPosition{18, 55}, gridPosition{0, 2}, true},
		{pixelPosition{18, 56}, gridPosition{0, 2}, true},
		{pixelPosition{18, 57}, gridPosition{0, 2}, true},
		{pixelPosition{18, 58}, gridPosition{0, 2}, true},
		{pixelPosition{18, 59}, gridPosition{0, 2}, true},
		{pixelPosition{18, 60}, gridPosition{0, 2}, true},
		{pixelPosition{18, 61}, gridPosition{0, 2}, true},
		{pixelPosition{18, 62}, gridPosition{0, 2}, true},
		{pixelPosition{18, 63}, gridPosition{0, 2}, true},
		{pixelPosition{18, 64}, gridPosition{0, 2}, true},
		{pixelPosition{18, 65}, gridPosition{0, 2}, true},
		{pixelPosition{18, 66}, gridPosition{0, 2}, true},
		{pixelPosition{18, 67}, gridPosition{0, 2}, true},
		{pixelPosition{18, 68}, gridPosition{0, 2}, true},
		{pixelPosition{18, 69}, gridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{pixelPosition{19, 50}, gridPosition{0, 2}, true},
		{pixelPosition{19, 51}, gridPosition{0, 2}, true},
		{pixelPosition{19, 52}, gridPosition{0, 2}, true},
		{pixelPosition{19, 53}, gridPosition{0, 2}, true},
		{pixelPosition{19, 54}, gridPosition{0, 2}, true},
		{pixelPosition{19, 55}, gridPosition{0, 2}, true},
		{pixelPosition{19, 56}, gridPosition{0, 2}, true},
		{pixelPosition{19, 57}, gridPosition{0, 2}, true},
		{pixelPosition{19, 58}, gridPosition{0, 2}, true},
		{pixelPosition{19, 59}, gridPosition{0, 2}, true},
		{pixelPosition{19, 60}, gridPosition{0, 2}, true},
		{pixelPosition{19, 61}, gridPosition{0, 2}, true},
		{pixelPosition{19, 62}, gridPosition{0, 2}, true},
		{pixelPosition{19, 63}, gridPosition{0, 2}, true},
		{pixelPosition{19, 64}, gridPosition{0, 2}, true},
		{pixelPosition{19, 65}, gridPosition{0, 2}, true},
		{pixelPosition{19, 66}, gridPosition{0, 2}, true},
		{pixelPosition{19, 67}, gridPosition{0, 2}, true},
		{pixelPosition{19, 68}, gridPosition{0, 2}, true},
		{pixelPosition{19, 69}, gridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{pixelPosition{20, 50}, gridPosition{0, 2}, true},
		{pixelPosition{20, 51}, gridPosition{0, 2}, true},
		{pixelPosition{20, 52}, gridPosition{0, 2}, true},
		{pixelPosition{20, 53}, gridPosition{0, 2}, true},
		{pixelPosition{20, 54}, gridPosition{0, 2}, true},
		{pixelPosition{20, 55}, gridPosition{0, 2}, true},
		{pixelPosition{20, 56}, gridPosition{0, 2}, true},
		{pixelPosition{20, 57}, gridPosition{0, 2}, true},
		{pixelPosition{20, 58}, gridPosition{0, 2}, true},
		{pixelPosition{20, 59}, gridPosition{0, 2}, true},
		{pixelPosition{20, 60}, gridPosition{0, 2}, true},
		{pixelPosition{20, 61}, gridPosition{0, 2}, true},
		{pixelPosition{20, 62}, gridPosition{0, 2}, true},
		{pixelPosition{20, 63}, gridPosition{0, 2}, true},
		{pixelPosition{20, 64}, gridPosition{0, 2}, true},
		{pixelPosition{20, 65}, gridPosition{0, 2}, true},
		{pixelPosition{20, 66}, gridPosition{0, 2}, true},
		{pixelPosition{20, 67}, gridPosition{0, 2}, true},
		{pixelPosition{20, 68}, gridPosition{0, 2}, true},
		{pixelPosition{20, 69}, gridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		{pixelPosition{41, 50}, gridPosition{0, 2}, true},
		{pixelPosition{41, 51}, gridPosition{0, 2}, true},
		{pixelPosition{41, 52}, gridPosition{0, 2}, true},
		{pixelPosition{41, 53}, gridPosition{0, 2}, true},
		{pixelPosition{41, 54}, gridPosition{0, 2}, true},
		{pixelPosition{41, 55}, gridPosition{0, 2}, true},
		{pixelPosition{41, 56}, gridPosition{0, 2}, true},
		{pixelPosition{41, 57}, gridPosition{0, 2}, true},
		{pixelPosition{41, 58}, gridPosition{0, 2}, true},
		{pixelPosition{41, 59}, gridPosition{0, 2}, true},
		{pixelPosition{41, 60}, gridPosition{0, 2}, true},
		{pixelPosition{41, 61}, gridPosition{0, 2}, true},
		{pixelPosition{41, 62}, gridPosition{0, 2}, true},
		{pixelPosition{41, 63}, gridPosition{0, 2}, true},
		{pixelPosition{41, 64}, gridPosition{0, 2}, true},
		{pixelPosition{41, 65}, gridPosition{0, 2}, true},
		{pixelPosition{41, 66}, gridPosition{0, 2}, true},
		{pixelPosition{41, 67}, gridPosition{0, 2}, true},
		{pixelPosition{41, 68}, gridPosition{0, 2}, true},
		{pixelPosition{41, 69}, gridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		// hexagons are right left side inclusive, right side exclusive
		{pixelPosition{42, 50}, gridPosition{1, 2}, true},
		{pixelPosition{42, 51}, gridPosition{1, 2}, true},
		{pixelPosition{42, 52}, gridPosition{1, 2}, true},
		{pixelPosition{42, 53}, gridPosition{1, 2}, true},
		{pixelPosition{42, 54}, gridPosition{1, 2}, true},
		{pixelPosition{42, 55}, gridPosition{1, 2}, true},
		{pixelPosition{42, 56}, gridPosition{1, 2}, true},
		{pixelPosition{42, 57}, gridPosition{1, 2}, true},
		{pixelPosition{42, 58}, gridPosition{1, 2}, true},
		{pixelPosition{42, 59}, gridPosition{1, 2}, true},
		{pixelPosition{42, 60}, gridPosition{1, 2}, true},
		{pixelPosition{42, 61}, gridPosition{1, 2}, true},
		{pixelPosition{42, 62}, gridPosition{1, 2}, true},
		{pixelPosition{42, 63}, gridPosition{1, 2}, true},
		{pixelPosition{42, 64}, gridPosition{1, 2}, true},
		{pixelPosition{42, 65}, gridPosition{1, 2}, true},
		{pixelPosition{42, 66}, gridPosition{1, 2}, true},
		{pixelPosition{42, 67}, gridPosition{1, 2}, true},
		{pixelPosition{42, 68}, gridPosition{1, 2}, true},
		{pixelPosition{42, 69}, gridPosition{1, 2}, true},

		// border of hexagon 0,2 and 1,2
		{pixelPosition{43, 50}, gridPosition{1, 2}, true},
		{pixelPosition{43, 51}, gridPosition{1, 2}, true},
		{pixelPosition{43, 52}, gridPosition{1, 2}, true},
		{pixelPosition{43, 53}, gridPosition{1, 2}, true},
		{pixelPosition{43, 54}, gridPosition{1, 2}, true},
		{pixelPosition{43, 55}, gridPosition{1, 2}, true},
		{pixelPosition{43, 56}, gridPosition{1, 2}, true},
		{pixelPosition{43, 57}, gridPosition{1, 2}, true},
		{pixelPosition{43, 58}, gridPosition{1, 2}, true},
		{pixelPosition{43, 59}, gridPosition{1, 2}, true},
		{pixelPosition{43, 60}, gridPosition{1, 2}, true},
		{pixelPosition{43, 61}, gridPosition{1, 2}, true},
		{pixelPosition{43, 62}, gridPosition{1, 2}, true},
		{pixelPosition{43, 63}, gridPosition{1, 2}, true},
		{pixelPosition{43, 64}, gridPosition{1, 2}, true},
		{pixelPosition{43, 65}, gridPosition{1, 2}, true},
		{pixelPosition{43, 66}, gridPosition{1, 2}, true},
		{pixelPosition{43, 67}, gridPosition{1, 2}, true},
		{pixelPosition{43, 68}, gridPosition{1, 2}, true},
		{pixelPosition{43, 69}, gridPosition{1, 2}, true},

		// right side of hexagon 1,2
		{pixelPosition{64, 50}, gridPosition{1, 2}, true},
		{pixelPosition{64, 51}, gridPosition{1, 2}, true},
		{pixelPosition{64, 52}, gridPosition{1, 2}, true},
		{pixelPosition{64, 53}, gridPosition{1, 2}, true},
		{pixelPosition{64, 54}, gridPosition{1, 2}, true},
		{pixelPosition{64, 55}, gridPosition{1, 2}, true},
		{pixelPosition{64, 56}, gridPosition{1, 2}, true},
		{pixelPosition{64, 57}, gridPosition{1, 2}, true},
		{pixelPosition{64, 58}, gridPosition{1, 2}, true},
		{pixelPosition{64, 59}, gridPosition{1, 2}, true},
		{pixelPosition{64, 60}, gridPosition{1, 2}, true},
		{pixelPosition{64, 61}, gridPosition{1, 2}, true},
		{pixelPosition{64, 62}, gridPosition{1, 2}, true},
		{pixelPosition{64, 63}, gridPosition{1, 2}, true},
		{pixelPosition{64, 64}, gridPosition{1, 2}, true},
		{pixelPosition{64, 65}, gridPosition{1, 2}, true},
		{pixelPosition{64, 66}, gridPosition{1, 2}, true},
		{pixelPosition{64, 67}, gridPosition{1, 2}, true},
		{pixelPosition{64, 68}, gridPosition{1, 2}, true},
		{pixelPosition{64, 69}, gridPosition{1, 2}, true},

		// cotinue right side of hexagon 1,2
		{pixelPosition{65, 50}, gridPosition{1, 2}, true},
		{pixelPosition{65, 51}, gridPosition{1, 2}, true},
		{pixelPosition{65, 52}, gridPosition{1, 2}, true},
		{pixelPosition{65, 53}, gridPosition{1, 2}, true},
		{pixelPosition{65, 54}, gridPosition{1, 2}, true},
		{pixelPosition{65, 55}, gridPosition{1, 2}, true},
		{pixelPosition{65, 56}, gridPosition{1, 2}, true},
		{pixelPosition{65, 57}, gridPosition{1, 2}, true},
		{pixelPosition{65, 58}, gridPosition{1, 2}, true},
		{pixelPosition{65, 59}, gridPosition{1, 2}, true},
		{pixelPosition{65, 60}, gridPosition{1, 2}, true},
		{pixelPosition{65, 61}, gridPosition{1, 2}, true},
		{pixelPosition{65, 62}, gridPosition{1, 2}, true},
		{pixelPosition{65, 63}, gridPosition{1, 2}, true},
		{pixelPosition{65, 64}, gridPosition{1, 2}, true},
		{pixelPosition{65, 65}, gridPosition{1, 2}, true},
		{pixelPosition{65, 66}, gridPosition{1, 2}, true},
		{pixelPosition{65, 67}, gridPosition{1, 2}, true},
		{pixelPosition{65, 68}, gridPosition{1, 2}, true},
		{pixelPosition{65, 69}, gridPosition{1, 2}, true},

		// continue right side of hexagon hexagon 1,2
		{pixelPosition{66, 50}, gridPosition{1, 2}, true},
		{pixelPosition{66, 51}, gridPosition{1, 2}, true},
		{pixelPosition{66, 52}, gridPosition{1, 2}, true},
		{pixelPosition{66, 53}, gridPosition{1, 2}, true},
		{pixelPosition{66, 54}, gridPosition{1, 2}, true},
		{pixelPosition{66, 55}, gridPosition{1, 2}, true},
		{pixelPosition{66, 56}, gridPosition{1, 2}, true},
		{pixelPosition{66, 57}, gridPosition{1, 2}, true},
		{pixelPosition{66, 58}, gridPosition{1, 2}, true},
		{pixelPosition{66, 59}, gridPosition{1, 2}, true},
		{pixelPosition{66, 60}, gridPosition{1, 2}, true},
		{pixelPosition{66, 61}, gridPosition{1, 2}, true},
		{pixelPosition{66, 62}, gridPosition{1, 2}, true},
		{pixelPosition{66, 63}, gridPosition{1, 2}, true},
		{pixelPosition{66, 64}, gridPosition{1, 2}, true},
		{pixelPosition{66, 65}, gridPosition{1, 2}, true},
		{pixelPosition{66, 66}, gridPosition{1, 2}, true},
		{pixelPosition{66, 67}, gridPosition{1, 2}, true},
		{pixelPosition{66, 68}, gridPosition{1, 2}, true},
		{pixelPosition{66, 69}, gridPosition{1, 2}, true},
	}

	for _, tt := range tests {
		var sWantedOk string
		if tt.wantedOk {
			sWantedOk = "exist"
		} else {
			sWantedOk = "not exist"
		}

		testname := fmt.Sprintf(
			"PixelPosition:%s converts to hex grid position {%d,%d} and want it to %s",
			tt.pixelPos.log(),
			tt.wantedPos.x,
			tt.wantedPos.y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := getGridPos(tt.pixelPos, &hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if pos.x != tt.wantedPos.x || pos.y != tt.wantedPos.y {
				t.Errorf("For pixel position %s, Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.pixelPos.log(),
					tt.wantedPos.x,
					tt.wantedPos.y,
					pos.x,
					pos.y)
			}
		})
	}
}
