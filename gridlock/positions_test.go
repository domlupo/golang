package main

import (
	"fmt"
	"testing"
)

func TestGetPixelPos(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(1280, 720, 46, 40, 19, 10)

	var tests = []struct {
		gridPos        gridPosition
		hexPoint       hexPoint
		wantedPixelPos PixelPosition
		wantedOk       bool
	}{

		// check hexagon {1, 1}, {1, 2}, {2,1} to ensure hex points are correct
		{gridPosition{1, 1}, center, PixelPosition{42, 30}, true},
		{gridPosition{1, 1}, northPoint, PixelPosition{42, 10}, true},
		{gridPosition{1, 1}, northeastPoint, PixelPosition{64, 20}, true},
		{gridPosition{1, 1}, southeastPoint, PixelPosition{64, 39}, true},
		{gridPosition{1, 1}, southPoint, PixelPosition{42, 49}, true},
		{gridPosition{1, 1}, southwestPoint, PixelPosition{19, 39}, true},
		{gridPosition{1, 1}, northwestPoint, PixelPosition{19, 20}, true},
		{gridPosition{2, 1}, center, PixelPosition{88, 30}, true},
		{gridPosition{2, 1}, northPoint, PixelPosition{88, 10}, true},
		{gridPosition{2, 1}, northeastPoint, PixelPosition{110, 20}, true},
		{gridPosition{2, 1}, southeastPoint, PixelPosition{110, 39}, true},
		{gridPosition{2, 1}, southPoint, PixelPosition{88, 49}, true},
		{gridPosition{2, 1}, southwestPoint, PixelPosition{65, 39}, true},
		{gridPosition{2, 1}, northwestPoint, PixelPosition{65, 20}, true},
		{gridPosition{1, 2}, center, PixelPosition{65, 60}, true},
		{gridPosition{1, 2}, northPoint, PixelPosition{65, 40}, true},
		{gridPosition{1, 2}, northeastPoint, PixelPosition{87, 50}, true},
		{gridPosition{1, 2}, southeastPoint, PixelPosition{87, 69}, true},
		{gridPosition{1, 2}, southPoint, PixelPosition{65, 79}, true},
		{gridPosition{1, 2}, southwestPoint, PixelPosition{42, 69}, true},
		{gridPosition{1, 2}, northwestPoint, PixelPosition{42, 50}, true},

		// check hexagon {26,23} to ensure large x and large y positions are correct
		{gridPosition{27, 23}, center, PixelPosition{1238, 690}, true},
		{gridPosition{27, 23}, northPoint, PixelPosition{1238, 670}, true},
		{gridPosition{27, 23}, northeastPoint, PixelPosition{1260, 680}, true},
		{gridPosition{27, 23}, southeastPoint, PixelPosition{1260, 699}, true},
		{gridPosition{27, 23}, southPoint, PixelPosition{1238, 709}, true},
		{gridPosition{27, 23}, southwestPoint, PixelPosition{1215, 699}, true},
		{gridPosition{27, 23}, northwestPoint, PixelPosition{1215, 680}, true},

		// check top left hexagon {0,0}
		{gridPosition{0, 0}, center, PixelPosition{19, 0}, true},
		{gridPosition{0, 0}, northPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 0}, northeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 0}, southeastPoint, PixelPosition{41, 9}, true},
		{gridPosition{0, 0}, southPoint, PixelPosition{19, 19}, true},
		{gridPosition{0, 0}, southwestPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 0}, northwestPoint, PixelPosition{0, 0}, false},

		// check bottom left hexagon {0,24}
		{gridPosition{0, 24}, center, PixelPosition{0, 0}, false},
		{gridPosition{0, 24}, northPoint, PixelPosition{19, 700}, true},
		{gridPosition{0, 24}, northeastPoint, PixelPosition{41, 710}, true},
		{gridPosition{0, 24}, southeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 24}, southPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 24}, southwestPoint, PixelPosition{0, 0}, false},
		{gridPosition{0, 24}, northwestPoint, PixelPosition{0, 0}, false},

		// check top right hexagon {27,0}
		{gridPosition{27, 0}, center, PixelPosition{1261, 0}, true},
		{gridPosition{27, 0}, northPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 0}, northeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 0}, southeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 0}, southPoint, PixelPosition{1261, 19}, true},
		{gridPosition{27, 0}, southwestPoint, PixelPosition{1238, 9}, true},
		{gridPosition{27, 0}, northwestPoint, PixelPosition{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{gridPosition{27, 24}, center, PixelPosition{0, 0}, false},
		{gridPosition{27, 24}, northPoint, PixelPosition{1261, 700}, true},
		{gridPosition{27, 24}, northeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southeastPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 24}, southwestPoint, PixelPosition{0, 0}, false},
		{gridPosition{27, 24}, northwestPoint, PixelPosition{1238, 710}, true},
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
			tt.wantedPixelPos.X,
			tt.wantedPixelPos.Y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := tt.gridPos.getPixelPos(tt.hexPoint, &hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if !pos.equals(tt.wantedPixelPos) {
				t.Errorf("For grid pos {%d,%d} and hex point %d  Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.gridPos.x,
					tt.gridPos.y,
					tt.hexPoint,
					tt.wantedPixelPos.X,
					tt.wantedPixelPos.Y,
					pos.X,
					pos.Y)
			}
		})
	}
}

func TestGetGridPos(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(1280, 720, 46, 40, 19, 10)

	var tests = []struct {
		pixelPos  PixelPosition
		wantedPos gridPosition
		wantedOk  bool
	}{
		// top left of entire hex grid
		{PixelPosition{0, 0}, gridPosition{0, 0}, true},
		{PixelPosition{-1, 0}, gridPosition{0, 0}, false},
		{PixelPosition{0, -1}, gridPosition{0, 0}, false},
		{PixelPosition{-1, -1}, gridPosition{0, 0}, false},

		// first left border point hex change around 0,12 pixel position
		{PixelPosition{-1, 11}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 12}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 13}, gridPosition{0, 0}, false},
		{PixelPosition{0, 11}, gridPosition{0, 0}, true},
		{PixelPosition{0, 12}, gridPosition{0, 1}, true},
		{PixelPosition{0, 13}, gridPosition{0, 1}, true},
		{PixelPosition{1, 11}, gridPosition{0, 0}, true},
		{PixelPosition{1, 12}, gridPosition{0, 1}, true},
		{PixelPosition{1, 13}, gridPosition{0, 1}, true},

		// second left border point hex change around 0,48 pixel position
		{PixelPosition{-1, 47}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 48}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 49}, gridPosition{0, 0}, false},
		{PixelPosition{0, 47}, gridPosition{0, 1}, true},
		{PixelPosition{0, 48}, gridPosition{0, 2}, true},
		{PixelPosition{0, 49}, gridPosition{0, 2}, true},
		{PixelPosition{1, 47}, gridPosition{0, 1}, true},
		{PixelPosition{1, 48}, gridPosition{0, 2}, true},
		{PixelPosition{1, 49}, gridPosition{0, 2}, true},

		// third left border point hex change around 0,72 pixel position
		{PixelPosition{-1, 71}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 72}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 73}, gridPosition{0, 0}, false},
		{PixelPosition{0, 71}, gridPosition{0, 2}, true},
		{PixelPosition{0, 72}, gridPosition{0, 3}, true},
		{PixelPosition{0, 73}, gridPosition{0, 3}, true},
		{PixelPosition{1, 71}, gridPosition{0, 2}, true},
		{PixelPosition{1, 72}, gridPosition{0, 3}, true},
		{PixelPosition{1, 73}, gridPosition{0, 3}, true},

		// first top border point hex change around 42,0 pixel position
		{PixelPosition{41, -1}, gridPosition{0, 0}, false},
		{PixelPosition{41, 0}, gridPosition{0, 0}, true},
		{PixelPosition{41, 1}, gridPosition{0, 0}, true},
		{PixelPosition{42, -1}, gridPosition{0, 0}, false},
		{PixelPosition{42, 0}, gridPosition{1, 0}, true},
		{PixelPosition{42, 1}, gridPosition{1, 0}, true},
		{PixelPosition{43, -1}, gridPosition{0, 0}, false},
		{PixelPosition{43, 0}, gridPosition{1, 0}, true},
		{PixelPosition{43, 1}, gridPosition{1, 0}, true},

		// second top border point hex change around 88, 0 pixel position
		{PixelPosition{87, -1}, gridPosition{0, 0}, false},
		{PixelPosition{87, 0}, gridPosition{1, 0}, true},
		{PixelPosition{87, 1}, gridPosition{1, 0}, true},
		{PixelPosition{88, -1}, gridPosition{0, 0}, false},
		{PixelPosition{88, 0}, gridPosition{2, 0}, true},
		{PixelPosition{88, 1}, gridPosition{2, 0}, true},
		{PixelPosition{89, -1}, gridPosition{0, 0}, false},
		{PixelPosition{89, 0}, gridPosition{2, 0}, true},
		{PixelPosition{89, 1}, gridPosition{2, 0}, true},

		// second top border point hex change around 134,0 pixel position
		{PixelPosition{133, -1}, gridPosition{0, 0}, false},
		{PixelPosition{133, 0}, gridPosition{2, 0}, true},
		{PixelPosition{133, 1}, gridPosition{2, 0}, true},
		{PixelPosition{134, -1}, gridPosition{0, 0}, false},
		{PixelPosition{134, 0}, gridPosition{3, 0}, true},
		{PixelPosition{134, 1}, gridPosition{3, 0}, true},
		{PixelPosition{135, -1}, gridPosition{0, 0}, false},
		{PixelPosition{135, 0}, gridPosition{3, 0}, true},
		{PixelPosition{135, 1}, gridPosition{3, 0}, true},

		// hexagon points near 19, 80 pixel position
		{PixelPosition{18, 79}, gridPosition{0, 2}, true},
		{PixelPosition{18, 80}, gridPosition{0, 3}, true},
		{PixelPosition{18, 81}, gridPosition{0, 3}, true},
		{PixelPosition{19, 79}, gridPosition{0, 2}, true},
		{PixelPosition{19, 80}, gridPosition{1, 3}, true},
		{PixelPosition{19, 81}, gridPosition{1, 3}, true},
		{PixelPosition{20, 79}, gridPosition{1, 3}, true},
		{PixelPosition{20, 80}, gridPosition{1, 3}, true},
		{PixelPosition{20, 81}, gridPosition{1, 3}, true},

		// hexagon points near 42, 70 pixel position
		{PixelPosition{41, 69}, gridPosition{0, 2}, true},
		{PixelPosition{41, 70}, gridPosition{1, 3}, true},
		{PixelPosition{41, 71}, gridPosition{1, 3}, true},
		{PixelPosition{42, 69}, gridPosition{1, 2}, true},
		{PixelPosition{42, 70}, gridPosition{1, 3}, true},
		{PixelPosition{42, 71}, gridPosition{1, 3}, true},
		{PixelPosition{43, 69}, gridPosition{1, 2}, true},
		{PixelPosition{43, 70}, gridPosition{1, 2}, true},
		{PixelPosition{43, 71}, gridPosition{1, 3}, true},

		// hexagon points near 65, 80 pixel position
		{PixelPosition{64, 79}, gridPosition{1, 2}, true},
		{PixelPosition{64, 80}, gridPosition{1, 3}, true},
		{PixelPosition{64, 81}, gridPosition{1, 3}, true},
		{PixelPosition{65, 79}, gridPosition{1, 2}, true},
		{PixelPosition{65, 80}, gridPosition{2, 3}, true},
		{PixelPosition{65, 81}, gridPosition{2, 3}, true},
		{PixelPosition{66, 79}, gridPosition{2, 3}, true},
		{PixelPosition{66, 80}, gridPosition{2, 3}, true},
		{PixelPosition{66, 81}, gridPosition{2, 3}, true},

		// hexagon points near 88, 10 pixel position
		{PixelPosition{87, 9}, gridPosition{1, 0}, true},
		{PixelPosition{87, 10}, gridPosition{2, 1}, true},
		{PixelPosition{87, 11}, gridPosition{2, 1}, true},
		{PixelPosition{88, 9}, gridPosition{2, 0}, true},
		{PixelPosition{88, 10}, gridPosition{2, 1}, true},
		{PixelPosition{88, 11}, gridPosition{2, 1}, true},
		{PixelPosition{89, 9}, gridPosition{2, 0}, true},
		{PixelPosition{89, 10}, gridPosition{2, 0}, true},
		{PixelPosition{89, 11}, gridPosition{2, 1}, true},

		// hexagon points near 88, 50 pixel position
		{PixelPosition{87, 49}, gridPosition{2, 1}, true},
		{PixelPosition{87, 50}, gridPosition{1, 2}, true},
		{PixelPosition{87, 51}, gridPosition{1, 2}, true},
		{PixelPosition{88, 49}, gridPosition{2, 1}, true},
		{PixelPosition{88, 50}, gridPosition{2, 2}, true},
		{PixelPosition{88, 51}, gridPosition{2, 2}, true},
		{PixelPosition{89, 49}, gridPosition{2, 2}, true},
		{PixelPosition{89, 50}, gridPosition{2, 2}, true},
		{PixelPosition{89, 51}, gridPosition{2, 2}, true},

		// hexagon points near 88, 70 pixel position
		{PixelPosition{87, 69}, gridPosition{1, 2}, true},
		{PixelPosition{87, 70}, gridPosition{2, 3}, true},
		{PixelPosition{87, 71}, gridPosition{2, 3}, true},
		{PixelPosition{88, 69}, gridPosition{2, 2}, true},
		{PixelPosition{88, 70}, gridPosition{2, 3}, true},
		{PixelPosition{88, 71}, gridPosition{2, 3}, true},
		{PixelPosition{89, 69}, gridPosition{2, 2}, true},
		{PixelPosition{89, 70}, gridPosition{2, 2}, true},
		{PixelPosition{89, 71}, gridPosition{2, 3}, true},

		// bottom left of entire hex grid
		{PixelPosition{-1, 719}, gridPosition{0, 0}, false},
		{PixelPosition{-1, 720}, gridPosition{0, 0}, false},
		{PixelPosition{0, 719}, gridPosition{0, 24}, true},
		{PixelPosition{0, 720}, gridPosition{0, 0}, false},

		// top right of entire hex grid
		{PixelPosition{1279, -1}, gridPosition{0, 0}, false},
		{PixelPosition{1280, -1}, gridPosition{0, 0}, false},
		{PixelPosition{1279, 0}, gridPosition{27, 0}, true},
		{PixelPosition{1280, 0}, gridPosition{0, 0}, false},

		// bottom right of entire hex grid
		{PixelPosition{1279, 719}, gridPosition{27, 24}, true},
		{PixelPosition{1280, 719}, gridPosition{0, 0}, false},
		{PixelPosition{1279, 720}, gridPosition{0, 0}, false},
		{PixelPosition{1280, 720}, gridPosition{0, 0}, false},

		// first hexagon y section (hexagon look up table)
		// between hexagons 0,0 1,1 and 1,0

		// sanity checks
		{PixelPosition{19, 10}, gridPosition{0, 0}, true},
		{PixelPosition{65, 10}, gridPosition{1, 0}, true},
		{PixelPosition{41, 20}, gridPosition{1, 1}, true},
		{PixelPosition{42, 20}, gridPosition{1, 1}, true},
		{PixelPosition{43, 20}, gridPosition{1, 1}, true},

		// 0,0 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{19, 19}, gridPosition{0, 0}, true},
		{PixelPosition{20, 18}, gridPosition{0, 0}, true},
		{PixelPosition{21, 18}, gridPosition{0, 0}, true},
		{PixelPosition{22, 17}, gridPosition{0, 0}, true},
		{PixelPosition{23, 17}, gridPosition{0, 0}, true},
		{PixelPosition{24, 17}, gridPosition{0, 0}, true},
		{PixelPosition{25, 16}, gridPosition{0, 0}, true},
		{PixelPosition{26, 16}, gridPosition{0, 0}, true},
		{PixelPosition{27, 15}, gridPosition{0, 0}, true},
		{PixelPosition{28, 15}, gridPosition{0, 0}, true},
		{PixelPosition{29, 14}, gridPosition{0, 0}, true},
		{PixelPosition{30, 14}, gridPosition{0, 0}, true},
		{PixelPosition{31, 14}, gridPosition{0, 0}, true},
		{PixelPosition{32, 13}, gridPosition{0, 0}, true},
		{PixelPosition{33, 13}, gridPosition{0, 0}, true},
		{PixelPosition{34, 12}, gridPosition{0, 0}, true},
		{PixelPosition{35, 12}, gridPosition{0, 0}, true},
		{PixelPosition{36, 11}, gridPosition{0, 0}, true},
		{PixelPosition{37, 11}, gridPosition{0, 0}, true},
		{PixelPosition{38, 11}, gridPosition{0, 0}, true},
		{PixelPosition{39, 10}, gridPosition{0, 0}, true},
		{PixelPosition{40, 10}, gridPosition{0, 0}, true},
		{PixelPosition{41, 9}, gridPosition{0, 0}, true},

		// 1, 1 hexagon slope check between itself and 0, 0 hexagon
		{PixelPosition{19, 20}, gridPosition{1, 1}, true},
		{PixelPosition{20, 19}, gridPosition{1, 1}, true},
		{PixelPosition{21, 19}, gridPosition{1, 1}, true},
		{PixelPosition{22, 18}, gridPosition{1, 1}, true},
		{PixelPosition{23, 18}, gridPosition{1, 1}, true},
		{PixelPosition{24, 18}, gridPosition{1, 1}, true},
		{PixelPosition{25, 17}, gridPosition{1, 1}, true},
		{PixelPosition{26, 17}, gridPosition{1, 1}, true},
		{PixelPosition{27, 16}, gridPosition{1, 1}, true},
		{PixelPosition{28, 16}, gridPosition{1, 1}, true},
		{PixelPosition{29, 15}, gridPosition{1, 1}, true},
		{PixelPosition{30, 15}, gridPosition{1, 1}, true},
		{PixelPosition{31, 15}, gridPosition{1, 1}, true},
		{PixelPosition{32, 14}, gridPosition{1, 1}, true},
		{PixelPosition{33, 14}, gridPosition{1, 1}, true},
		{PixelPosition{34, 13}, gridPosition{1, 1}, true},
		{PixelPosition{35, 13}, gridPosition{1, 1}, true},
		{PixelPosition{36, 12}, gridPosition{1, 1}, true},
		{PixelPosition{37, 12}, gridPosition{1, 1}, true},
		{PixelPosition{38, 12}, gridPosition{1, 1}, true},
		{PixelPosition{39, 11}, gridPosition{1, 1}, true},
		{PixelPosition{40, 11}, gridPosition{1, 1}, true},
		{PixelPosition{41, 10}, gridPosition{1, 1}, true},

		// 1,0 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{42, 9}, gridPosition{1, 0}, true},
		{PixelPosition{43, 10}, gridPosition{1, 0}, true},
		{PixelPosition{44, 10}, gridPosition{1, 0}, true},
		{PixelPosition{45, 11}, gridPosition{1, 0}, true},
		{PixelPosition{46, 11}, gridPosition{1, 0}, true},
		{PixelPosition{47, 11}, gridPosition{1, 0}, true},
		{PixelPosition{48, 12}, gridPosition{1, 0}, true},
		{PixelPosition{49, 12}, gridPosition{1, 0}, true},
		{PixelPosition{50, 13}, gridPosition{1, 0}, true},
		{PixelPosition{51, 13}, gridPosition{1, 0}, true},

		{PixelPosition{52, 14}, gridPosition{1, 0}, true},
		{PixelPosition{53, 14}, gridPosition{1, 0}, true},
		{PixelPosition{54, 14}, gridPosition{1, 0}, true},
		{PixelPosition{55, 15}, gridPosition{1, 0}, true},
		{PixelPosition{56, 15}, gridPosition{1, 0}, true},
		{PixelPosition{57, 16}, gridPosition{1, 0}, true},
		{PixelPosition{58, 16}, gridPosition{1, 0}, true},
		{PixelPosition{59, 17}, gridPosition{1, 0}, true},
		{PixelPosition{60, 17}, gridPosition{1, 0}, true},
		{PixelPosition{61, 17}, gridPosition{1, 0}, true},
		{PixelPosition{62, 18}, gridPosition{1, 0}, true},
		{PixelPosition{63, 18}, gridPosition{1, 0}, true},
		{PixelPosition{64, 19}, gridPosition{1, 0}, true},

		// 1,1 hexagon slope check between itself and 1, 0 hexagon
		{PixelPosition{42, 10}, gridPosition{1, 1}, true},
		{PixelPosition{43, 11}, gridPosition{1, 1}, true},
		{PixelPosition{44, 11}, gridPosition{1, 1}, true},
		{PixelPosition{45, 12}, gridPosition{1, 1}, true},
		{PixelPosition{46, 12}, gridPosition{1, 1}, true},
		{PixelPosition{47, 12}, gridPosition{1, 1}, true},
		{PixelPosition{48, 13}, gridPosition{1, 1}, true},
		{PixelPosition{49, 13}, gridPosition{1, 1}, true},
		{PixelPosition{50, 14}, gridPosition{1, 1}, true},
		{PixelPosition{51, 14}, gridPosition{1, 1}, true},
		{PixelPosition{52, 15}, gridPosition{1, 1}, true},
		{PixelPosition{53, 15}, gridPosition{1, 1}, true},
		{PixelPosition{54, 15}, gridPosition{1, 1}, true},
		{PixelPosition{55, 16}, gridPosition{1, 1}, true},
		{PixelPosition{56, 16}, gridPosition{1, 1}, true},
		{PixelPosition{57, 17}, gridPosition{1, 1}, true},
		{PixelPosition{58, 17}, gridPosition{1, 1}, true},
		{PixelPosition{59, 18}, gridPosition{1, 1}, true},
		{PixelPosition{60, 18}, gridPosition{1, 1}, true},
		{PixelPosition{61, 18}, gridPosition{1, 1}, true},
		{PixelPosition{62, 19}, gridPosition{1, 1}, true},
		{PixelPosition{63, 19}, gridPosition{1, 1}, true},
		{PixelPosition{64, 20}, gridPosition{1, 1}, true},

		// second hexagon y section (full rectangle of hexagon 1,1)

		// right side of hexagon 0,1
		{PixelPosition{18, 20}, gridPosition{0, 1}, true},
		{PixelPosition{18, 21}, gridPosition{0, 1}, true},
		{PixelPosition{18, 22}, gridPosition{0, 1}, true},
		{PixelPosition{18, 23}, gridPosition{0, 1}, true},
		{PixelPosition{18, 24}, gridPosition{0, 1}, true},
		{PixelPosition{18, 25}, gridPosition{0, 1}, true},
		{PixelPosition{18, 26}, gridPosition{0, 1}, true},
		{PixelPosition{18, 27}, gridPosition{0, 1}, true},
		{PixelPosition{18, 28}, gridPosition{0, 1}, true},
		{PixelPosition{18, 29}, gridPosition{0, 1}, true},
		{PixelPosition{18, 30}, gridPosition{0, 1}, true},
		{PixelPosition{18, 31}, gridPosition{0, 1}, true},
		{PixelPosition{18, 32}, gridPosition{0, 1}, true},
		{PixelPosition{18, 33}, gridPosition{0, 1}, true},
		{PixelPosition{18, 34}, gridPosition{0, 1}, true},
		{PixelPosition{18, 35}, gridPosition{0, 1}, true},
		{PixelPosition{18, 36}, gridPosition{0, 1}, true},
		{PixelPosition{18, 37}, gridPosition{0, 1}, true},
		{PixelPosition{18, 38}, gridPosition{0, 1}, true},
		{PixelPosition{18, 39}, gridPosition{0, 1}, true},

		// left side of hexagon 1,1
		// hexagons are right side inclusive and left side exclusive
		{PixelPosition{19, 20}, gridPosition{1, 1}, true},
		{PixelPosition{19, 21}, gridPosition{1, 1}, true},
		{PixelPosition{19, 22}, gridPosition{1, 1}, true},
		{PixelPosition{19, 23}, gridPosition{1, 1}, true},
		{PixelPosition{19, 24}, gridPosition{1, 1}, true},
		{PixelPosition{19, 25}, gridPosition{1, 1}, true},
		{PixelPosition{19, 26}, gridPosition{1, 1}, true},
		{PixelPosition{19, 27}, gridPosition{1, 1}, true},
		{PixelPosition{19, 28}, gridPosition{1, 1}, true},
		{PixelPosition{19, 29}, gridPosition{1, 1}, true},
		{PixelPosition{19, 30}, gridPosition{1, 1}, true},
		{PixelPosition{19, 31}, gridPosition{1, 1}, true},
		{PixelPosition{19, 32}, gridPosition{1, 1}, true},
		{PixelPosition{19, 33}, gridPosition{1, 1}, true},
		{PixelPosition{19, 34}, gridPosition{1, 1}, true},
		{PixelPosition{19, 35}, gridPosition{1, 1}, true},
		{PixelPosition{19, 36}, gridPosition{1, 1}, true},
		{PixelPosition{19, 37}, gridPosition{1, 1}, true},
		{PixelPosition{19, 38}, gridPosition{1, 1}, true},
		{PixelPosition{19, 39}, gridPosition{1, 1}, true},
		// continue left side of hexagon 1,1
		{PixelPosition{20, 20}, gridPosition{1, 1}, true},
		{PixelPosition{20, 21}, gridPosition{1, 1}, true},
		{PixelPosition{20, 22}, gridPosition{1, 1}, true},
		{PixelPosition{20, 23}, gridPosition{1, 1}, true},
		{PixelPosition{20, 24}, gridPosition{1, 1}, true},
		{PixelPosition{20, 25}, gridPosition{1, 1}, true},
		{PixelPosition{20, 26}, gridPosition{1, 1}, true},
		{PixelPosition{20, 27}, gridPosition{1, 1}, true},
		{PixelPosition{20, 28}, gridPosition{1, 1}, true},
		{PixelPosition{20, 29}, gridPosition{1, 1}, true},
		{PixelPosition{20, 30}, gridPosition{1, 1}, true},
		{PixelPosition{20, 31}, gridPosition{1, 1}, true},
		{PixelPosition{20, 32}, gridPosition{1, 1}, true},
		{PixelPosition{20, 33}, gridPosition{1, 1}, true},
		{PixelPosition{20, 34}, gridPosition{1, 1}, true},
		{PixelPosition{20, 35}, gridPosition{1, 1}, true},
		{PixelPosition{20, 36}, gridPosition{1, 1}, true},
		{PixelPosition{20, 37}, gridPosition{1, 1}, true},
		{PixelPosition{20, 38}, gridPosition{1, 1}, true},
		{PixelPosition{20, 39}, gridPosition{1, 1}, true},

		// right side of hexagon 1,1
		{PixelPosition{64, 20}, gridPosition{1, 1}, true},
		{PixelPosition{64, 21}, gridPosition{1, 1}, true},
		{PixelPosition{64, 22}, gridPosition{1, 1}, true},
		{PixelPosition{64, 23}, gridPosition{1, 1}, true},
		{PixelPosition{64, 24}, gridPosition{1, 1}, true},
		{PixelPosition{64, 25}, gridPosition{1, 1}, true},
		{PixelPosition{64, 26}, gridPosition{1, 1}, true},
		{PixelPosition{64, 27}, gridPosition{1, 1}, true},
		{PixelPosition{64, 28}, gridPosition{1, 1}, true},
		{PixelPosition{64, 29}, gridPosition{1, 1}, true},
		{PixelPosition{64, 30}, gridPosition{1, 1}, true},
		{PixelPosition{64, 31}, gridPosition{1, 1}, true},
		{PixelPosition{64, 32}, gridPosition{1, 1}, true},
		{PixelPosition{64, 33}, gridPosition{1, 1}, true},
		{PixelPosition{64, 34}, gridPosition{1, 1}, true},
		{PixelPosition{64, 35}, gridPosition{1, 1}, true},
		{PixelPosition{64, 36}, gridPosition{1, 1}, true},
		{PixelPosition{64, 37}, gridPosition{1, 1}, true},
		{PixelPosition{64, 38}, gridPosition{1, 1}, true},
		{PixelPosition{64, 39}, gridPosition{1, 1}, true},

		// left side of hexagon 2,1
		// hexagons are right side inclusive and left side exclusive
		{PixelPosition{65, 20}, gridPosition{2, 1}, true},
		{PixelPosition{65, 21}, gridPosition{2, 1}, true},
		{PixelPosition{65, 22}, gridPosition{2, 1}, true},
		{PixelPosition{65, 23}, gridPosition{2, 1}, true},
		{PixelPosition{65, 24}, gridPosition{2, 1}, true},
		{PixelPosition{65, 25}, gridPosition{2, 1}, true},
		{PixelPosition{65, 26}, gridPosition{2, 1}, true},
		{PixelPosition{65, 27}, gridPosition{2, 1}, true},
		{PixelPosition{65, 28}, gridPosition{2, 1}, true},
		{PixelPosition{65, 29}, gridPosition{2, 1}, true},
		{PixelPosition{65, 30}, gridPosition{2, 1}, true},
		{PixelPosition{65, 31}, gridPosition{2, 1}, true},
		{PixelPosition{65, 32}, gridPosition{2, 1}, true},
		{PixelPosition{65, 33}, gridPosition{2, 1}, true},
		{PixelPosition{65, 34}, gridPosition{2, 1}, true},
		{PixelPosition{65, 35}, gridPosition{2, 1}, true},
		{PixelPosition{65, 36}, gridPosition{2, 1}, true},
		{PixelPosition{65, 37}, gridPosition{2, 1}, true},
		{PixelPosition{65, 38}, gridPosition{2, 1}, true},
		{PixelPosition{65, 39}, gridPosition{2, 1}, true},
		// continue left side of hexagon 2,1
		{PixelPosition{66, 20}, gridPosition{2, 1}, true},
		{PixelPosition{66, 21}, gridPosition{2, 1}, true},
		{PixelPosition{66, 22}, gridPosition{2, 1}, true},
		{PixelPosition{66, 23}, gridPosition{2, 1}, true},
		{PixelPosition{66, 24}, gridPosition{2, 1}, true},
		{PixelPosition{66, 25}, gridPosition{2, 1}, true},
		{PixelPosition{66, 26}, gridPosition{2, 1}, true},
		{PixelPosition{66, 27}, gridPosition{2, 1}, true},
		{PixelPosition{66, 28}, gridPosition{2, 1}, true},
		{PixelPosition{66, 29}, gridPosition{2, 1}, true},
		{PixelPosition{66, 30}, gridPosition{2, 1}, true},
		{PixelPosition{66, 31}, gridPosition{2, 1}, true},
		{PixelPosition{66, 32}, gridPosition{2, 1}, true},
		{PixelPosition{66, 33}, gridPosition{2, 1}, true},
		{PixelPosition{66, 34}, gridPosition{2, 1}, true},
		{PixelPosition{66, 35}, gridPosition{2, 1}, true},
		{PixelPosition{66, 36}, gridPosition{2, 1}, true},
		{PixelPosition{66, 37}, gridPosition{2, 1}, true},
		{PixelPosition{66, 38}, gridPosition{2, 1}, true},
		{PixelPosition{66, 39}, gridPosition{2, 1}, true},

		// third hexagon y section (inverted hexagon look up table)
		// between hexagons 0,2 1,1 and 1,2

		// sanity checks
		{PixelPosition{19, 49}, gridPosition{0, 2}, true},
		{PixelPosition{65, 49}, gridPosition{1, 2}, true},
		{PixelPosition{41, 30}, gridPosition{1, 1}, true},
		{PixelPosition{42, 30}, gridPosition{1, 1}, true},
		{PixelPosition{43, 30}, gridPosition{1, 1}, true},

		// 0,2 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{19, 40}, gridPosition{0, 2}, true},
		{PixelPosition{20, 41}, gridPosition{0, 2}, true},
		{PixelPosition{21, 41}, gridPosition{0, 2}, true},
		{PixelPosition{22, 42}, gridPosition{0, 2}, true},
		{PixelPosition{23, 42}, gridPosition{0, 2}, true},
		{PixelPosition{24, 42}, gridPosition{0, 2}, true},
		{PixelPosition{25, 43}, gridPosition{0, 2}, true},
		{PixelPosition{26, 43}, gridPosition{0, 2}, true},
		{PixelPosition{27, 44}, gridPosition{0, 2}, true},
		{PixelPosition{28, 44}, gridPosition{0, 2}, true},
		{PixelPosition{29, 45}, gridPosition{0, 2}, true},
		{PixelPosition{30, 45}, gridPosition{0, 2}, true},
		{PixelPosition{31, 45}, gridPosition{0, 2}, true},
		{PixelPosition{32, 46}, gridPosition{0, 2}, true},
		{PixelPosition{33, 46}, gridPosition{0, 2}, true},
		{PixelPosition{34, 47}, gridPosition{0, 2}, true},
		{PixelPosition{35, 47}, gridPosition{0, 2}, true},
		{PixelPosition{36, 48}, gridPosition{0, 2}, true},
		{PixelPosition{37, 48}, gridPosition{0, 2}, true},
		{PixelPosition{38, 48}, gridPosition{0, 2}, true},
		{PixelPosition{39, 49}, gridPosition{0, 2}, true},
		{PixelPosition{40, 49}, gridPosition{0, 2}, true},
		{PixelPosition{41, 50}, gridPosition{0, 2}, true},

		// 1, 1 hexagon slope check between itself and 0, 2 hexagon
		{PixelPosition{19, 39}, gridPosition{1, 1}, true},
		{PixelPosition{20, 40}, gridPosition{1, 1}, true},
		{PixelPosition{21, 40}, gridPosition{1, 1}, true},
		{PixelPosition{22, 41}, gridPosition{1, 1}, true},
		{PixelPosition{23, 41}, gridPosition{1, 1}, true},
		{PixelPosition{24, 41}, gridPosition{1, 1}, true},
		{PixelPosition{25, 42}, gridPosition{1, 1}, true},
		{PixelPosition{26, 42}, gridPosition{1, 1}, true},
		{PixelPosition{27, 43}, gridPosition{1, 1}, true},
		{PixelPosition{28, 43}, gridPosition{1, 1}, true},
		{PixelPosition{29, 44}, gridPosition{1, 1}, true},
		{PixelPosition{30, 44}, gridPosition{1, 1}, true},
		{PixelPosition{31, 44}, gridPosition{1, 1}, true},
		{PixelPosition{32, 45}, gridPosition{1, 1}, true},
		{PixelPosition{33, 45}, gridPosition{1, 1}, true},
		{PixelPosition{34, 46}, gridPosition{1, 1}, true},
		{PixelPosition{35, 46}, gridPosition{1, 1}, true},
		{PixelPosition{36, 47}, gridPosition{1, 1}, true},
		{PixelPosition{37, 47}, gridPosition{1, 1}, true},
		{PixelPosition{38, 47}, gridPosition{1, 1}, true},
		{PixelPosition{39, 48}, gridPosition{1, 1}, true},
		{PixelPosition{40, 48}, gridPosition{1, 1}, true},
		{PixelPosition{41, 49}, gridPosition{1, 1}, true},

		// 1,1 hexagon slope check between itself and 1, 2 hexagon
		{PixelPosition{42, 49}, gridPosition{1, 1}, true},
		{PixelPosition{43, 48}, gridPosition{1, 1}, true},
		{PixelPosition{44, 48}, gridPosition{1, 1}, true},
		{PixelPosition{45, 47}, gridPosition{1, 1}, true},
		{PixelPosition{46, 47}, gridPosition{1, 1}, true},
		{PixelPosition{47, 47}, gridPosition{1, 1}, true},
		{PixelPosition{48, 46}, gridPosition{1, 1}, true},
		{PixelPosition{49, 46}, gridPosition{1, 1}, true},
		{PixelPosition{50, 45}, gridPosition{1, 1}, true},
		{PixelPosition{51, 45}, gridPosition{1, 1}, true},
		{PixelPosition{52, 44}, gridPosition{1, 1}, true},
		{PixelPosition{53, 44}, gridPosition{1, 1}, true},
		{PixelPosition{54, 44}, gridPosition{1, 1}, true},
		{PixelPosition{55, 43}, gridPosition{1, 1}, true},
		{PixelPosition{56, 43}, gridPosition{1, 1}, true},
		{PixelPosition{57, 42}, gridPosition{1, 1}, true},
		{PixelPosition{58, 42}, gridPosition{1, 1}, true},
		{PixelPosition{59, 41}, gridPosition{1, 1}, true},
		{PixelPosition{60, 41}, gridPosition{1, 1}, true},
		{PixelPosition{61, 41}, gridPosition{1, 1}, true},
		{PixelPosition{62, 40}, gridPosition{1, 1}, true},
		{PixelPosition{63, 40}, gridPosition{1, 1}, true},
		{PixelPosition{64, 39}, gridPosition{1, 1}, true},

		// 1,2 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{42, 50}, gridPosition{1, 2}, true},
		{PixelPosition{43, 49}, gridPosition{1, 2}, true},
		{PixelPosition{44, 49}, gridPosition{1, 2}, true},
		{PixelPosition{45, 48}, gridPosition{1, 2}, true},
		{PixelPosition{46, 48}, gridPosition{1, 2}, true},
		{PixelPosition{47, 48}, gridPosition{1, 2}, true},
		{PixelPosition{48, 47}, gridPosition{1, 2}, true},
		{PixelPosition{49, 47}, gridPosition{1, 2}, true},
		{PixelPosition{50, 46}, gridPosition{1, 2}, true},
		{PixelPosition{51, 46}, gridPosition{1, 2}, true},
		{PixelPosition{52, 45}, gridPosition{1, 2}, true},
		{PixelPosition{53, 45}, gridPosition{1, 2}, true},
		{PixelPosition{54, 45}, gridPosition{1, 2}, true},
		{PixelPosition{55, 44}, gridPosition{1, 2}, true},
		{PixelPosition{56, 44}, gridPosition{1, 2}, true},
		{PixelPosition{57, 43}, gridPosition{1, 2}, true},
		{PixelPosition{58, 43}, gridPosition{1, 2}, true},
		{PixelPosition{59, 42}, gridPosition{1, 2}, true},
		{PixelPosition{60, 42}, gridPosition{1, 2}, true},
		{PixelPosition{61, 42}, gridPosition{1, 2}, true},
		{PixelPosition{62, 41}, gridPosition{1, 2}, true},
		{PixelPosition{63, 41}, gridPosition{1, 2}, true},
		{PixelPosition{64, 40}, gridPosition{1, 2}, true},

		// fourth hexagon y section (split rectangle between hexagon 0,2 and 1,2)

		// middle of hexagon 0,2
		{PixelPosition{18, 50}, gridPosition{0, 2}, true},
		{PixelPosition{18, 51}, gridPosition{0, 2}, true},
		{PixelPosition{18, 52}, gridPosition{0, 2}, true},
		{PixelPosition{18, 53}, gridPosition{0, 2}, true},
		{PixelPosition{18, 54}, gridPosition{0, 2}, true},
		{PixelPosition{18, 55}, gridPosition{0, 2}, true},
		{PixelPosition{18, 56}, gridPosition{0, 2}, true},
		{PixelPosition{18, 57}, gridPosition{0, 2}, true},
		{PixelPosition{18, 58}, gridPosition{0, 2}, true},
		{PixelPosition{18, 59}, gridPosition{0, 2}, true},
		{PixelPosition{18, 60}, gridPosition{0, 2}, true},
		{PixelPosition{18, 61}, gridPosition{0, 2}, true},
		{PixelPosition{18, 62}, gridPosition{0, 2}, true},
		{PixelPosition{18, 63}, gridPosition{0, 2}, true},
		{PixelPosition{18, 64}, gridPosition{0, 2}, true},
		{PixelPosition{18, 65}, gridPosition{0, 2}, true},
		{PixelPosition{18, 66}, gridPosition{0, 2}, true},
		{PixelPosition{18, 67}, gridPosition{0, 2}, true},
		{PixelPosition{18, 68}, gridPosition{0, 2}, true},
		{PixelPosition{18, 69}, gridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{PixelPosition{19, 50}, gridPosition{0, 2}, true},
		{PixelPosition{19, 51}, gridPosition{0, 2}, true},
		{PixelPosition{19, 52}, gridPosition{0, 2}, true},
		{PixelPosition{19, 53}, gridPosition{0, 2}, true},
		{PixelPosition{19, 54}, gridPosition{0, 2}, true},
		{PixelPosition{19, 55}, gridPosition{0, 2}, true},
		{PixelPosition{19, 56}, gridPosition{0, 2}, true},
		{PixelPosition{19, 57}, gridPosition{0, 2}, true},
		{PixelPosition{19, 58}, gridPosition{0, 2}, true},
		{PixelPosition{19, 59}, gridPosition{0, 2}, true},
		{PixelPosition{19, 60}, gridPosition{0, 2}, true},
		{PixelPosition{19, 61}, gridPosition{0, 2}, true},
		{PixelPosition{19, 62}, gridPosition{0, 2}, true},
		{PixelPosition{19, 63}, gridPosition{0, 2}, true},
		{PixelPosition{19, 64}, gridPosition{0, 2}, true},
		{PixelPosition{19, 65}, gridPosition{0, 2}, true},
		{PixelPosition{19, 66}, gridPosition{0, 2}, true},
		{PixelPosition{19, 67}, gridPosition{0, 2}, true},
		{PixelPosition{19, 68}, gridPosition{0, 2}, true},
		{PixelPosition{19, 69}, gridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{PixelPosition{20, 50}, gridPosition{0, 2}, true},
		{PixelPosition{20, 51}, gridPosition{0, 2}, true},
		{PixelPosition{20, 52}, gridPosition{0, 2}, true},
		{PixelPosition{20, 53}, gridPosition{0, 2}, true},
		{PixelPosition{20, 54}, gridPosition{0, 2}, true},
		{PixelPosition{20, 55}, gridPosition{0, 2}, true},
		{PixelPosition{20, 56}, gridPosition{0, 2}, true},
		{PixelPosition{20, 57}, gridPosition{0, 2}, true},
		{PixelPosition{20, 58}, gridPosition{0, 2}, true},
		{PixelPosition{20, 59}, gridPosition{0, 2}, true},
		{PixelPosition{20, 60}, gridPosition{0, 2}, true},
		{PixelPosition{20, 61}, gridPosition{0, 2}, true},
		{PixelPosition{20, 62}, gridPosition{0, 2}, true},
		{PixelPosition{20, 63}, gridPosition{0, 2}, true},
		{PixelPosition{20, 64}, gridPosition{0, 2}, true},
		{PixelPosition{20, 65}, gridPosition{0, 2}, true},
		{PixelPosition{20, 66}, gridPosition{0, 2}, true},
		{PixelPosition{20, 67}, gridPosition{0, 2}, true},
		{PixelPosition{20, 68}, gridPosition{0, 2}, true},
		{PixelPosition{20, 69}, gridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		{PixelPosition{41, 50}, gridPosition{0, 2}, true},
		{PixelPosition{41, 51}, gridPosition{0, 2}, true},
		{PixelPosition{41, 52}, gridPosition{0, 2}, true},
		{PixelPosition{41, 53}, gridPosition{0, 2}, true},
		{PixelPosition{41, 54}, gridPosition{0, 2}, true},
		{PixelPosition{41, 55}, gridPosition{0, 2}, true},
		{PixelPosition{41, 56}, gridPosition{0, 2}, true},
		{PixelPosition{41, 57}, gridPosition{0, 2}, true},
		{PixelPosition{41, 58}, gridPosition{0, 2}, true},
		{PixelPosition{41, 59}, gridPosition{0, 2}, true},
		{PixelPosition{41, 60}, gridPosition{0, 2}, true},
		{PixelPosition{41, 61}, gridPosition{0, 2}, true},
		{PixelPosition{41, 62}, gridPosition{0, 2}, true},
		{PixelPosition{41, 63}, gridPosition{0, 2}, true},
		{PixelPosition{41, 64}, gridPosition{0, 2}, true},
		{PixelPosition{41, 65}, gridPosition{0, 2}, true},
		{PixelPosition{41, 66}, gridPosition{0, 2}, true},
		{PixelPosition{41, 67}, gridPosition{0, 2}, true},
		{PixelPosition{41, 68}, gridPosition{0, 2}, true},
		{PixelPosition{41, 69}, gridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		// hexagons are right left side inclusive, right side exclusive
		{PixelPosition{42, 50}, gridPosition{1, 2}, true},
		{PixelPosition{42, 51}, gridPosition{1, 2}, true},
		{PixelPosition{42, 52}, gridPosition{1, 2}, true},
		{PixelPosition{42, 53}, gridPosition{1, 2}, true},
		{PixelPosition{42, 54}, gridPosition{1, 2}, true},
		{PixelPosition{42, 55}, gridPosition{1, 2}, true},
		{PixelPosition{42, 56}, gridPosition{1, 2}, true},
		{PixelPosition{42, 57}, gridPosition{1, 2}, true},
		{PixelPosition{42, 58}, gridPosition{1, 2}, true},
		{PixelPosition{42, 59}, gridPosition{1, 2}, true},
		{PixelPosition{42, 60}, gridPosition{1, 2}, true},
		{PixelPosition{42, 61}, gridPosition{1, 2}, true},
		{PixelPosition{42, 62}, gridPosition{1, 2}, true},
		{PixelPosition{42, 63}, gridPosition{1, 2}, true},
		{PixelPosition{42, 64}, gridPosition{1, 2}, true},
		{PixelPosition{42, 65}, gridPosition{1, 2}, true},
		{PixelPosition{42, 66}, gridPosition{1, 2}, true},
		{PixelPosition{42, 67}, gridPosition{1, 2}, true},
		{PixelPosition{42, 68}, gridPosition{1, 2}, true},
		{PixelPosition{42, 69}, gridPosition{1, 2}, true},

		// border of hexagon 0,2 and 1,2
		{PixelPosition{43, 50}, gridPosition{1, 2}, true},
		{PixelPosition{43, 51}, gridPosition{1, 2}, true},
		{PixelPosition{43, 52}, gridPosition{1, 2}, true},
		{PixelPosition{43, 53}, gridPosition{1, 2}, true},
		{PixelPosition{43, 54}, gridPosition{1, 2}, true},
		{PixelPosition{43, 55}, gridPosition{1, 2}, true},
		{PixelPosition{43, 56}, gridPosition{1, 2}, true},
		{PixelPosition{43, 57}, gridPosition{1, 2}, true},
		{PixelPosition{43, 58}, gridPosition{1, 2}, true},
		{PixelPosition{43, 59}, gridPosition{1, 2}, true},
		{PixelPosition{43, 60}, gridPosition{1, 2}, true},
		{PixelPosition{43, 61}, gridPosition{1, 2}, true},
		{PixelPosition{43, 62}, gridPosition{1, 2}, true},
		{PixelPosition{43, 63}, gridPosition{1, 2}, true},
		{PixelPosition{43, 64}, gridPosition{1, 2}, true},
		{PixelPosition{43, 65}, gridPosition{1, 2}, true},
		{PixelPosition{43, 66}, gridPosition{1, 2}, true},
		{PixelPosition{43, 67}, gridPosition{1, 2}, true},
		{PixelPosition{43, 68}, gridPosition{1, 2}, true},
		{PixelPosition{43, 69}, gridPosition{1, 2}, true},

		// right side of hexagon 1,2
		{PixelPosition{64, 50}, gridPosition{1, 2}, true},
		{PixelPosition{64, 51}, gridPosition{1, 2}, true},
		{PixelPosition{64, 52}, gridPosition{1, 2}, true},
		{PixelPosition{64, 53}, gridPosition{1, 2}, true},
		{PixelPosition{64, 54}, gridPosition{1, 2}, true},
		{PixelPosition{64, 55}, gridPosition{1, 2}, true},
		{PixelPosition{64, 56}, gridPosition{1, 2}, true},
		{PixelPosition{64, 57}, gridPosition{1, 2}, true},
		{PixelPosition{64, 58}, gridPosition{1, 2}, true},
		{PixelPosition{64, 59}, gridPosition{1, 2}, true},
		{PixelPosition{64, 60}, gridPosition{1, 2}, true},
		{PixelPosition{64, 61}, gridPosition{1, 2}, true},
		{PixelPosition{64, 62}, gridPosition{1, 2}, true},
		{PixelPosition{64, 63}, gridPosition{1, 2}, true},
		{PixelPosition{64, 64}, gridPosition{1, 2}, true},
		{PixelPosition{64, 65}, gridPosition{1, 2}, true},
		{PixelPosition{64, 66}, gridPosition{1, 2}, true},
		{PixelPosition{64, 67}, gridPosition{1, 2}, true},
		{PixelPosition{64, 68}, gridPosition{1, 2}, true},
		{PixelPosition{64, 69}, gridPosition{1, 2}, true},

		// continue right side of hexagon 1,2
		{PixelPosition{65, 50}, gridPosition{1, 2}, true},
		{PixelPosition{65, 51}, gridPosition{1, 2}, true},
		{PixelPosition{65, 52}, gridPosition{1, 2}, true},
		{PixelPosition{65, 53}, gridPosition{1, 2}, true},
		{PixelPosition{65, 54}, gridPosition{1, 2}, true},
		{PixelPosition{65, 55}, gridPosition{1, 2}, true},
		{PixelPosition{65, 56}, gridPosition{1, 2}, true},
		{PixelPosition{65, 57}, gridPosition{1, 2}, true},
		{PixelPosition{65, 58}, gridPosition{1, 2}, true},
		{PixelPosition{65, 59}, gridPosition{1, 2}, true},
		{PixelPosition{65, 60}, gridPosition{1, 2}, true},
		{PixelPosition{65, 61}, gridPosition{1, 2}, true},
		{PixelPosition{65, 62}, gridPosition{1, 2}, true},
		{PixelPosition{65, 63}, gridPosition{1, 2}, true},
		{PixelPosition{65, 64}, gridPosition{1, 2}, true},
		{PixelPosition{65, 65}, gridPosition{1, 2}, true},
		{PixelPosition{65, 66}, gridPosition{1, 2}, true},
		{PixelPosition{65, 67}, gridPosition{1, 2}, true},
		{PixelPosition{65, 68}, gridPosition{1, 2}, true},
		{PixelPosition{65, 69}, gridPosition{1, 2}, true},

		// continue right side of hexagon hexagon 1,2
		{PixelPosition{66, 50}, gridPosition{1, 2}, true},
		{PixelPosition{66, 51}, gridPosition{1, 2}, true},
		{PixelPosition{66, 52}, gridPosition{1, 2}, true},
		{PixelPosition{66, 53}, gridPosition{1, 2}, true},
		{PixelPosition{66, 54}, gridPosition{1, 2}, true},
		{PixelPosition{66, 55}, gridPosition{1, 2}, true},
		{PixelPosition{66, 56}, gridPosition{1, 2}, true},
		{PixelPosition{66, 57}, gridPosition{1, 2}, true},
		{PixelPosition{66, 58}, gridPosition{1, 2}, true},
		{PixelPosition{66, 59}, gridPosition{1, 2}, true},
		{PixelPosition{66, 60}, gridPosition{1, 2}, true},
		{PixelPosition{66, 61}, gridPosition{1, 2}, true},
		{PixelPosition{66, 62}, gridPosition{1, 2}, true},
		{PixelPosition{66, 63}, gridPosition{1, 2}, true},
		{PixelPosition{66, 64}, gridPosition{1, 2}, true},
		{PixelPosition{66, 65}, gridPosition{1, 2}, true},
		{PixelPosition{66, 66}, gridPosition{1, 2}, true},
		{PixelPosition{66, 67}, gridPosition{1, 2}, true},
		{PixelPosition{66, 68}, gridPosition{1, 2}, true},
		{PixelPosition{66, 69}, gridPosition{1, 2}, true},
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
			pos, ok := tt.pixelPos.getGridPos(&hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if !pos.equals(tt.wantedPos) {
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

func TestGridPosAdjacent(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(1280, 720, 46, 40, 19, 10)

	var tests = []struct {
		gridPos       gridPosition
		hexSide       hexSide
		wantedGridPos gridPosition
		wantedOk      bool
	}{
		// check top left hexagon {0,0}
		{gridPosition{0, 0}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{0, 0}, eastSide, gridPosition{1, 0}, true},
		{gridPosition{0, 0}, southeastSide, gridPosition{1, 1}, true},
		{gridPosition{0, 0}, southwestSide, gridPosition{0, 1}, true},
		{gridPosition{0, 0}, westSide, gridPosition{0, 0}, false},
		{gridPosition{0, 0}, northwestSide, gridPosition{0, 0}, false},

		// check top left hexagon {1, 0}
		{gridPosition{1, 0}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{1, 0}, eastSide, gridPosition{2, 0}, true},
		{gridPosition{1, 0}, southeastSide, gridPosition{2, 1}, true},
		{gridPosition{1, 0}, southwestSide, gridPosition{1, 1}, true},
		{gridPosition{1, 0}, westSide, gridPosition{0, 0}, true},
		{gridPosition{1, 0}, northwestSide, gridPosition{0, 0}, false},

		// check top left hexagon {0, 1}
		{gridPosition{0, 1}, northeastSide, gridPosition{0, 0}, true},
		{gridPosition{0, 1}, eastSide, gridPosition{1, 1}, true},
		{gridPosition{0, 1}, southeastSide, gridPosition{0, 2}, true},
		{gridPosition{0, 1}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{0, 1}, westSide, gridPosition{0, 0}, false},
		{gridPosition{0, 1}, northwestSide, gridPosition{0, 0}, false},

		// check bottom left hexagon {0,24}
		{gridPosition{0, 24}, northeastSide, gridPosition{1, 23}, true},
		{gridPosition{0, 24}, eastSide, gridPosition{1, 24}, true},
		{gridPosition{0, 24}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{0, 24}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{0, 24}, westSide, gridPosition{0, 0}, false},
		{gridPosition{0, 24}, northwestSide, gridPosition{0, 23}, true},

		// check bottom left hexagon {1,24}
		{gridPosition{1, 24}, northeastSide, gridPosition{2, 23}, true},
		{gridPosition{1, 24}, eastSide, gridPosition{2, 24}, true},
		{gridPosition{1, 24}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{1, 24}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{1, 24}, westSide, gridPosition{0, 24}, true},
		{gridPosition{1, 24}, northwestSide, gridPosition{1, 23}, true},

		// check bottom left hexagon {0,23}
		{gridPosition{0, 23}, northeastSide, gridPosition{0, 22}, true},
		{gridPosition{0, 23}, eastSide, gridPosition{1, 23}, true},
		{gridPosition{0, 23}, southeastSide, gridPosition{0, 24}, true},
		{gridPosition{0, 23}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{0, 23}, westSide, gridPosition{0, 0}, false},
		{gridPosition{0, 23}, northwestSide, gridPosition{0, 0}, false},

		// check top right hexagon {27,0}
		{gridPosition{27, 0}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{27, 0}, eastSide, gridPosition{0, 0}, false},
		{gridPosition{27, 0}, southeastSide, gridPosition{28, 1}, true},
		{gridPosition{27, 0}, southwestSide, gridPosition{27, 1}, true},
		{gridPosition{27, 0}, westSide, gridPosition{26, 0}, true},
		{gridPosition{27, 0}, northwestSide, gridPosition{0, 0}, false},

		// check top right hexagon {28,1}
		{gridPosition{28, 1}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 1}, eastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 1}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 1}, southwestSide, gridPosition{27, 2}, true},
		{gridPosition{28, 1}, westSide, gridPosition{27, 1}, true},
		{gridPosition{28, 1}, northwestSide, gridPosition{27, 0}, true},

		// check top right hexagon {26,0}
		{gridPosition{26, 0}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{26, 0}, eastSide, gridPosition{27, 0}, true},
		{gridPosition{26, 0}, southeastSide, gridPosition{27, 1}, true},
		{gridPosition{26, 0}, southwestSide, gridPosition{26, 1}, true},
		{gridPosition{26, 0}, westSide, gridPosition{25, 0}, true},
		{gridPosition{26, 0}, northwestSide, gridPosition{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{gridPosition{27, 24}, northeastSide, gridPosition{28, 23}, true},
		{gridPosition{27, 24}, eastSide, gridPosition{0, 0}, false},
		{gridPosition{27, 24}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{27, 24}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{27, 24}, westSide, gridPosition{26, 24}, true},
		{gridPosition{27, 24}, northwestSide, gridPosition{27, 23}, true},

		// check bottom right hexagon {26, 24}
		{gridPosition{26, 24}, northeastSide, gridPosition{27, 23}, true},
		{gridPosition{26, 24}, eastSide, gridPosition{27, 24}, true},
		{gridPosition{26, 24}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{26, 24}, southwestSide, gridPosition{0, 0}, false},
		{gridPosition{26, 24}, westSide, gridPosition{25, 24}, true},
		{gridPosition{26, 24}, northwestSide, gridPosition{26, 23}, true},

		// check bottom right hexagon {28, 23}
		{gridPosition{28, 23}, northeastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 23}, eastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 23}, southeastSide, gridPosition{0, 0}, false},
		{gridPosition{28, 23}, southwestSide, gridPosition{27, 24}, true},
		{gridPosition{28, 23}, westSide, gridPosition{27, 23}, true},
		{gridPosition{28, 23}, northwestSide, gridPosition{27, 22}, true},
	}

	for _, tt := range tests {
		var sWantedOk string
		if tt.wantedOk {
			sWantedOk = "exist"
		} else {
			sWantedOk = "not exist"
		}

		testname := fmt.Sprintf(
			"gridPos:{%d,%d} and hex side %d has adjacent position {%d,%d} and want it to %s",
			tt.gridPos.x,
			tt.gridPos.y,
			tt.hexSide,
			tt.wantedGridPos.x,
			tt.wantedGridPos.y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := tt.gridPos.adjacent(tt.hexSide, &hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if !pos.equals(tt.wantedGridPos) {
				t.Errorf("For grid pos {%d,%d} and hex side %d Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.gridPos.x,
					tt.gridPos.y,
					tt.hexSide,
					tt.wantedGridPos.x,
					tt.wantedGridPos.y,
					pos.x,
					pos.y)
			}
		})
	}

}

// TODO: write higher level tests with more cases once more source code is implemented, these are too low level and sparse.
func TestIntersects(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		l1               line
		l2               line
		wantedIntersects bool
	}{
		// does not intersect
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{3, -1}, PixelPosition{3, 1}}, false},
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{3, 1}, PixelPosition{3, -1}}, false},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{3, -1}, PixelPosition{3, 1}}, false},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{3, 1}, PixelPosition{3, -1}}, false},

		// does intersects
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{1, -1}, PixelPosition{1, 1}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{1, -1}, PixelPosition{1, 1}}, true},
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{1, 1}, PixelPosition{1, -1}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{1, 1}, PixelPosition{1, -1}}, true},

		// 3 collinear points
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{2, 0}, PixelPosition{2, 1}}, true},
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{2, 1}, PixelPosition{2, 0}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{2, 0}, PixelPosition{2, 1}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{2, 1}, PixelPosition{2, 0}}, true},

		// 4 collinear points
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{2, 0}, PixelPosition{4, 0}}, true},
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, line{PixelPosition{4, 0}, PixelPosition{2, 0}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{2, 0}, PixelPosition{4, 0}}, true},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, line{PixelPosition{4, 0}, PixelPosition{2, 0}}, true},

		// lines are equal
		{line{PixelPosition{0, 0}, PixelPosition{1, 0}}, line{PixelPosition{0, 0}, PixelPosition{1, 0}}, true},
		{line{PixelPosition{0, 0}, PixelPosition{1, 0}}, line{PixelPosition{1, 0}, PixelPosition{0, 0}}, true},
		{line{PixelPosition{1, 0}, PixelPosition{0, 0}}, line{PixelPosition{0, 0}, PixelPosition{1, 0}}, true},
		{line{PixelPosition{1, 0}, PixelPosition{0, 0}}, line{PixelPosition{1, 0}, PixelPosition{0, 0}}, true},
	}

	for _, tt := range tests {
		testname := "low level test to be replaced"
		t.Run(testname, func(t *testing.T) {
			intersects := tt.l1.intersects(tt.l2)
			if intersects != tt.wantedIntersects {
				t.Errorf("low level test failed")
			}
		})
	}

}
