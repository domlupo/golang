package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestGetPixelPos(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		gridPos      GridPosition
		hexPoint     hexPoint
		wantPixelPos Pixel
		wantOk       bool
	}{

		// check hexagon {1, 1}, {1, 2}, {2,1} to ensure hex points are correct
		{GridPosition{1, 1}, hpCenter, Pixel{42, 30}, true},
		{GridPosition{1, 1}, hpNorth, Pixel{42, 10}, true},
		{GridPosition{1, 1}, hpNortheast, Pixel{64, 20}, true},
		{GridPosition{1, 1}, hpSoutheast, Pixel{64, 39}, true},
		{GridPosition{1, 1}, hpSouth, Pixel{42, 49}, true},
		{GridPosition{1, 1}, hpSouthwest, Pixel{19, 39}, true},
		{GridPosition{1, 1}, hpNorthwest, Pixel{19, 20}, true},
		{GridPosition{2, 1}, hpCenter, Pixel{88, 30}, true},
		{GridPosition{2, 1}, hpNorth, Pixel{88, 10}, true},
		{GridPosition{2, 1}, hpNortheast, Pixel{110, 20}, true},
		{GridPosition{2, 1}, hpSoutheast, Pixel{110, 39}, true},
		{GridPosition{2, 1}, hpSouth, Pixel{88, 49}, true},
		{GridPosition{2, 1}, hpSouthwest, Pixel{65, 39}, true},
		{GridPosition{2, 1}, hpNorthwest, Pixel{65, 20}, true},
		{GridPosition{1, 2}, hpCenter, Pixel{65, 60}, true},
		{GridPosition{1, 2}, hpNorth, Pixel{65, 40}, true},
		{GridPosition{1, 2}, hpNortheast, Pixel{87, 50}, true},
		{GridPosition{1, 2}, hpSoutheast, Pixel{87, 69}, true},
		{GridPosition{1, 2}, hpSouth, Pixel{65, 79}, true},
		{GridPosition{1, 2}, hpSouthwest, Pixel{42, 69}, true},
		{GridPosition{1, 2}, hpNorthwest, Pixel{42, 50}, true},

		// check hexagon {26,23} to ensure large x and large y positions are correct
		{GridPosition{27, 23}, hpCenter, Pixel{1238, 690}, true},
		{GridPosition{27, 23}, hpNorth, Pixel{1238, 670}, true},
		{GridPosition{27, 23}, hpNortheast, Pixel{1260, 680}, true},
		{GridPosition{27, 23}, hpSoutheast, Pixel{1260, 699}, true},
		{GridPosition{27, 23}, hpSouth, Pixel{1238, 709}, true},
		{GridPosition{27, 23}, hpSouthwest, Pixel{1215, 699}, true},
		{GridPosition{27, 23}, hpNorthwest, Pixel{1215, 680}, true},

		// check top left hexagon {0,0}
		{GridPosition{0, 0}, hpCenter, Pixel{19, 0}, true},
		{GridPosition{0, 0}, hpNorth, Pixel{0, 0}, false},
		{GridPosition{0, 0}, hpNortheast, Pixel{0, 0}, false},
		{GridPosition{0, 0}, hpSoutheast, Pixel{41, 9}, true},
		{GridPosition{0, 0}, hpSouth, Pixel{19, 19}, true},
		{GridPosition{0, 0}, hpSouthwest, Pixel{0, 0}, false},
		{GridPosition{0, 0}, hpNorthwest, Pixel{0, 0}, false},

		// check bottom left hexagon {0,24}
		{GridPosition{0, 24}, hpCenter, Pixel{0, 0}, false},
		{GridPosition{0, 24}, hpNorth, Pixel{19, 700}, true},
		{GridPosition{0, 24}, hpNortheast, Pixel{41, 710}, true},
		{GridPosition{0, 24}, hpSoutheast, Pixel{0, 0}, false},
		{GridPosition{0, 24}, hpSouth, Pixel{0, 0}, false},
		{GridPosition{0, 24}, hpSouthwest, Pixel{0, 0}, false},
		{GridPosition{0, 24}, hpNorthwest, Pixel{0, 0}, false},

		// check top right hexagon {27,0}
		{GridPosition{27, 0}, hpCenter, Pixel{1261, 0}, true},
		{GridPosition{27, 0}, hpNorth, Pixel{0, 0}, false},
		{GridPosition{27, 0}, hpNortheast, Pixel{0, 0}, false},
		{GridPosition{27, 0}, hpSoutheast, Pixel{0, 0}, false},
		{GridPosition{27, 0}, hpSouth, Pixel{1261, 19}, true},
		{GridPosition{27, 0}, hpSouthwest, Pixel{1238, 9}, true},
		{GridPosition{27, 0}, hpNorthwest, Pixel{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{GridPosition{27, 24}, hpCenter, Pixel{0, 0}, false},
		{GridPosition{27, 24}, hpNorth, Pixel{1261, 700}, true},
		{GridPosition{27, 24}, hpNortheast, Pixel{0, 0}, false},
		{GridPosition{27, 24}, hpSoutheast, Pixel{0, 0}, false},
		{GridPosition{27, 24}, hpSouth, Pixel{0, 0}, false},
		{GridPosition{27, 24}, hpSouthwest, Pixel{0, 0}, false},
		{GridPosition{27, 24}, hpNorthwest, Pixel{1238, 710}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"gridPos: %s and hex point %s converts to hex grid position %s and want it to %s",
			tt.gridPos.log(),
			tt.hexPoint.string(),
			tt.wantPixelPos.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.getPixel(tt.gridPos, tt.hexPoint)
			if ok != tt.wantOk {
				t.Errorf("Want pixel position to %s", sWantOk)
			}
			if !pos.equals(tt.wantPixelPos) {
				t.Errorf("For grid pos %s and hex point %s  Want grid position %s but it was %s",
					tt.gridPos.log(),
					tt.hexPoint.string(),
					tt.wantPixelPos.log(),
					pos.log())
			}
		})
	}
}

func TestGetGridPos(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		pixel   Pixel
		wantPos GridPosition
		wantOk  bool
	}{
		// top left of entire hex grid
		{Pixel{0, 0}, GridPosition{0, 0}, true},
		{Pixel{-1, 0}, GridPosition{0, 0}, false},
		{Pixel{0, -1}, GridPosition{0, 0}, false},
		{Pixel{-1, -1}, GridPosition{0, 0}, false},

		// first left border point hex change around 0,12 pixel position
		{Pixel{-1, 11}, GridPosition{0, 0}, false},
		{Pixel{-1, 12}, GridPosition{0, 0}, false},
		{Pixel{-1, 13}, GridPosition{0, 0}, false},
		{Pixel{0, 11}, GridPosition{0, 0}, true},
		{Pixel{0, 12}, GridPosition{0, 1}, true},
		{Pixel{0, 13}, GridPosition{0, 1}, true},
		{Pixel{1, 11}, GridPosition{0, 0}, true},
		{Pixel{1, 12}, GridPosition{0, 1}, true},
		{Pixel{1, 13}, GridPosition{0, 1}, true},

		// second left border point hex change around 0,48 pixel position
		{Pixel{-1, 47}, GridPosition{0, 0}, false},
		{Pixel{-1, 48}, GridPosition{0, 0}, false},
		{Pixel{-1, 49}, GridPosition{0, 0}, false},
		{Pixel{0, 47}, GridPosition{0, 1}, true},
		{Pixel{0, 48}, GridPosition{0, 2}, true},
		{Pixel{0, 49}, GridPosition{0, 2}, true},
		{Pixel{1, 47}, GridPosition{0, 1}, true},
		{Pixel{1, 48}, GridPosition{0, 2}, true},
		{Pixel{1, 49}, GridPosition{0, 2}, true},

		// third left border point hex change around 0,72 pixel position
		{Pixel{-1, 71}, GridPosition{0, 0}, false},
		{Pixel{-1, 72}, GridPosition{0, 0}, false},
		{Pixel{-1, 73}, GridPosition{0, 0}, false},
		{Pixel{0, 71}, GridPosition{0, 2}, true},
		{Pixel{0, 72}, GridPosition{0, 3}, true},
		{Pixel{0, 73}, GridPosition{0, 3}, true},
		{Pixel{1, 71}, GridPosition{0, 2}, true},
		{Pixel{1, 72}, GridPosition{0, 3}, true},
		{Pixel{1, 73}, GridPosition{0, 3}, true},

		// first top border point hex change around 42,0 pixel position
		{Pixel{41, -1}, GridPosition{0, 0}, false},
		{Pixel{41, 0}, GridPosition{0, 0}, true},
		{Pixel{41, 1}, GridPosition{0, 0}, true},
		{Pixel{42, -1}, GridPosition{0, 0}, false},
		{Pixel{42, 0}, GridPosition{1, 0}, true},
		{Pixel{42, 1}, GridPosition{1, 0}, true},
		{Pixel{43, -1}, GridPosition{0, 0}, false},
		{Pixel{43, 0}, GridPosition{1, 0}, true},
		{Pixel{43, 1}, GridPosition{1, 0}, true},

		// second top border point hex change around 88, 0 pixel position
		{Pixel{87, -1}, GridPosition{0, 0}, false},
		{Pixel{87, 0}, GridPosition{1, 0}, true},
		{Pixel{87, 1}, GridPosition{1, 0}, true},
		{Pixel{88, -1}, GridPosition{0, 0}, false},
		{Pixel{88, 0}, GridPosition{2, 0}, true},
		{Pixel{88, 1}, GridPosition{2, 0}, true},
		{Pixel{89, -1}, GridPosition{0, 0}, false},
		{Pixel{89, 0}, GridPosition{2, 0}, true},
		{Pixel{89, 1}, GridPosition{2, 0}, true},

		// third top border point hex change around 134,0 pixel position
		{Pixel{133, -1}, GridPosition{0, 0}, false},
		{Pixel{133, 0}, GridPosition{2, 0}, true},
		{Pixel{133, 1}, GridPosition{2, 0}, true},
		{Pixel{134, -1}, GridPosition{0, 0}, false},
		{Pixel{134, 0}, GridPosition{3, 0}, true},
		{Pixel{134, 1}, GridPosition{3, 0}, true},
		{Pixel{135, -1}, GridPosition{0, 0}, false},
		{Pixel{135, 0}, GridPosition{3, 0}, true},
		{Pixel{135, 1}, GridPosition{3, 0}, true},

		// hexagon points near 19, 80 pixel position
		{Pixel{18, 79}, GridPosition{0, 2}, true},
		{Pixel{18, 80}, GridPosition{0, 3}, true},
		{Pixel{18, 81}, GridPosition{0, 3}, true},
		{Pixel{19, 79}, GridPosition{0, 2}, true},
		{Pixel{19, 80}, GridPosition{1, 3}, true},
		{Pixel{19, 81}, GridPosition{1, 3}, true},
		{Pixel{20, 79}, GridPosition{1, 3}, true},
		{Pixel{20, 80}, GridPosition{1, 3}, true},
		{Pixel{20, 81}, GridPosition{1, 3}, true},

		// hexagon points near 42, 70 pixel position
		{Pixel{41, 69}, GridPosition{0, 2}, true},
		{Pixel{41, 70}, GridPosition{1, 3}, true},
		{Pixel{41, 71}, GridPosition{1, 3}, true},
		{Pixel{42, 69}, GridPosition{1, 2}, true},
		{Pixel{42, 70}, GridPosition{1, 3}, true},
		{Pixel{42, 71}, GridPosition{1, 3}, true},
		{Pixel{43, 69}, GridPosition{1, 2}, true},
		{Pixel{43, 70}, GridPosition{1, 2}, true},
		{Pixel{43, 71}, GridPosition{1, 3}, true},

		// hexagon points near 65, 80 pixel position
		{Pixel{64, 79}, GridPosition{1, 2}, true},
		{Pixel{64, 80}, GridPosition{1, 3}, true},
		{Pixel{64, 81}, GridPosition{1, 3}, true},
		{Pixel{65, 79}, GridPosition{1, 2}, true},
		{Pixel{65, 80}, GridPosition{2, 3}, true},
		{Pixel{65, 81}, GridPosition{2, 3}, true},
		{Pixel{66, 79}, GridPosition{2, 3}, true},
		{Pixel{66, 80}, GridPosition{2, 3}, true},
		{Pixel{66, 81}, GridPosition{2, 3}, true},

		// hexagon points near 88, 10 pixel position
		{Pixel{87, 9}, GridPosition{1, 0}, true},
		{Pixel{87, 10}, GridPosition{2, 1}, true},
		{Pixel{87, 11}, GridPosition{2, 1}, true},
		{Pixel{88, 9}, GridPosition{2, 0}, true},
		{Pixel{88, 10}, GridPosition{2, 1}, true},
		{Pixel{88, 11}, GridPosition{2, 1}, true},
		{Pixel{89, 9}, GridPosition{2, 0}, true},
		{Pixel{89, 10}, GridPosition{2, 0}, true},
		{Pixel{89, 11}, GridPosition{2, 1}, true},

		// hexagon points near 88, 50 pixel position
		{Pixel{87, 49}, GridPosition{2, 1}, true},
		{Pixel{87, 50}, GridPosition{1, 2}, true},
		{Pixel{87, 51}, GridPosition{1, 2}, true},
		{Pixel{88, 49}, GridPosition{2, 1}, true},
		{Pixel{88, 50}, GridPosition{2, 2}, true},
		{Pixel{88, 51}, GridPosition{2, 2}, true},
		{Pixel{89, 49}, GridPosition{2, 2}, true},
		{Pixel{89, 50}, GridPosition{2, 2}, true},
		{Pixel{89, 51}, GridPosition{2, 2}, true},

		// hexagon points near 88, 70 pixel position
		{Pixel{87, 69}, GridPosition{1, 2}, true},
		{Pixel{87, 70}, GridPosition{2, 3}, true},
		{Pixel{87, 71}, GridPosition{2, 3}, true},
		{Pixel{88, 69}, GridPosition{2, 2}, true},
		{Pixel{88, 70}, GridPosition{2, 3}, true},
		{Pixel{88, 71}, GridPosition{2, 3}, true},
		{Pixel{89, 69}, GridPosition{2, 2}, true},
		{Pixel{89, 70}, GridPosition{2, 2}, true},
		{Pixel{89, 71}, GridPosition{2, 3}, true},

		// bottom left of entire hex grid
		{Pixel{-1, 719}, GridPosition{0, 0}, false},
		{Pixel{-1, 720}, GridPosition{0, 0}, false},
		{Pixel{0, 719}, GridPosition{0, 24}, true},
		{Pixel{0, 720}, GridPosition{0, 0}, false},

		// top right of entire hex grid
		{Pixel{1279, -1}, GridPosition{0, 0}, false},
		{Pixel{1280, -1}, GridPosition{0, 0}, false},
		{Pixel{1279, 0}, GridPosition{27, 0}, true},
		{Pixel{1280, 0}, GridPosition{0, 0}, false},

		// bottom right of entire hex grid
		{Pixel{1279, 719}, GridPosition{27, 24}, true},
		{Pixel{1280, 719}, GridPosition{0, 0}, false},
		{Pixel{1279, 720}, GridPosition{0, 0}, false},
		{Pixel{1280, 720}, GridPosition{0, 0}, false},

		// first hexagon y section (hexagon look up table)
		// between hexagons 0,0 1,1 and 1,0

		// sanity checks
		{Pixel{19, 10}, GridPosition{0, 0}, true},
		{Pixel{65, 10}, GridPosition{1, 0}, true},
		{Pixel{41, 20}, GridPosition{1, 1}, true},
		{Pixel{42, 20}, GridPosition{1, 1}, true},
		{Pixel{43, 20}, GridPosition{1, 1}, true},

		// 0,0 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{19, 19}, GridPosition{0, 0}, true},
		{Pixel{20, 18}, GridPosition{0, 0}, true},
		{Pixel{21, 18}, GridPosition{0, 0}, true},
		{Pixel{22, 17}, GridPosition{0, 0}, true},
		{Pixel{23, 17}, GridPosition{0, 0}, true},
		{Pixel{24, 17}, GridPosition{0, 0}, true},
		{Pixel{25, 16}, GridPosition{0, 0}, true},
		{Pixel{26, 16}, GridPosition{0, 0}, true},
		{Pixel{27, 15}, GridPosition{0, 0}, true},
		{Pixel{28, 15}, GridPosition{0, 0}, true},
		{Pixel{29, 14}, GridPosition{0, 0}, true},
		{Pixel{30, 14}, GridPosition{0, 0}, true},
		{Pixel{31, 14}, GridPosition{0, 0}, true},
		{Pixel{32, 13}, GridPosition{0, 0}, true},
		{Pixel{33, 13}, GridPosition{0, 0}, true},
		{Pixel{34, 12}, GridPosition{0, 0}, true},
		{Pixel{35, 12}, GridPosition{0, 0}, true},
		{Pixel{36, 11}, GridPosition{0, 0}, true},
		{Pixel{37, 11}, GridPosition{0, 0}, true},
		{Pixel{38, 11}, GridPosition{0, 0}, true},
		{Pixel{39, 10}, GridPosition{0, 0}, true},
		{Pixel{40, 10}, GridPosition{0, 0}, true},
		{Pixel{41, 9}, GridPosition{0, 0}, true},

		// 1, 1 hexagon slope check between itself and 0, 0 hexagon
		{Pixel{19, 20}, GridPosition{1, 1}, true},
		{Pixel{20, 19}, GridPosition{1, 1}, true},
		{Pixel{21, 19}, GridPosition{1, 1}, true},
		{Pixel{22, 18}, GridPosition{1, 1}, true},
		{Pixel{23, 18}, GridPosition{1, 1}, true},
		{Pixel{24, 18}, GridPosition{1, 1}, true},
		{Pixel{25, 17}, GridPosition{1, 1}, true},
		{Pixel{26, 17}, GridPosition{1, 1}, true},
		{Pixel{27, 16}, GridPosition{1, 1}, true},
		{Pixel{28, 16}, GridPosition{1, 1}, true},
		{Pixel{29, 15}, GridPosition{1, 1}, true},
		{Pixel{30, 15}, GridPosition{1, 1}, true},
		{Pixel{31, 15}, GridPosition{1, 1}, true},
		{Pixel{32, 14}, GridPosition{1, 1}, true},
		{Pixel{33, 14}, GridPosition{1, 1}, true},
		{Pixel{34, 13}, GridPosition{1, 1}, true},
		{Pixel{35, 13}, GridPosition{1, 1}, true},
		{Pixel{36, 12}, GridPosition{1, 1}, true},
		{Pixel{37, 12}, GridPosition{1, 1}, true},
		{Pixel{38, 12}, GridPosition{1, 1}, true},
		{Pixel{39, 11}, GridPosition{1, 1}, true},
		{Pixel{40, 11}, GridPosition{1, 1}, true},
		{Pixel{41, 10}, GridPosition{1, 1}, true},

		// 1,0 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{42, 9}, GridPosition{1, 0}, true},
		{Pixel{43, 10}, GridPosition{1, 0}, true},
		{Pixel{44, 10}, GridPosition{1, 0}, true},
		{Pixel{45, 11}, GridPosition{1, 0}, true},
		{Pixel{46, 11}, GridPosition{1, 0}, true},
		{Pixel{47, 11}, GridPosition{1, 0}, true},
		{Pixel{48, 12}, GridPosition{1, 0}, true},
		{Pixel{49, 12}, GridPosition{1, 0}, true},
		{Pixel{50, 13}, GridPosition{1, 0}, true},
		{Pixel{51, 13}, GridPosition{1, 0}, true},

		{Pixel{52, 14}, GridPosition{1, 0}, true},
		{Pixel{53, 14}, GridPosition{1, 0}, true},
		{Pixel{54, 14}, GridPosition{1, 0}, true},
		{Pixel{55, 15}, GridPosition{1, 0}, true},
		{Pixel{56, 15}, GridPosition{1, 0}, true},
		{Pixel{57, 16}, GridPosition{1, 0}, true},
		{Pixel{58, 16}, GridPosition{1, 0}, true},
		{Pixel{59, 17}, GridPosition{1, 0}, true},
		{Pixel{60, 17}, GridPosition{1, 0}, true},
		{Pixel{61, 17}, GridPosition{1, 0}, true},
		{Pixel{62, 18}, GridPosition{1, 0}, true},
		{Pixel{63, 18}, GridPosition{1, 0}, true},
		{Pixel{64, 19}, GridPosition{1, 0}, true},

		// 1,1 hexagon slope check between itself and 1, 0 hexagon
		{Pixel{42, 10}, GridPosition{1, 1}, true},
		{Pixel{43, 11}, GridPosition{1, 1}, true},
		{Pixel{44, 11}, GridPosition{1, 1}, true},
		{Pixel{45, 12}, GridPosition{1, 1}, true},
		{Pixel{46, 12}, GridPosition{1, 1}, true},
		{Pixel{47, 12}, GridPosition{1, 1}, true},
		{Pixel{48, 13}, GridPosition{1, 1}, true},
		{Pixel{49, 13}, GridPosition{1, 1}, true},
		{Pixel{50, 14}, GridPosition{1, 1}, true},
		{Pixel{51, 14}, GridPosition{1, 1}, true},
		{Pixel{52, 15}, GridPosition{1, 1}, true},
		{Pixel{53, 15}, GridPosition{1, 1}, true},
		{Pixel{54, 15}, GridPosition{1, 1}, true},
		{Pixel{55, 16}, GridPosition{1, 1}, true},
		{Pixel{56, 16}, GridPosition{1, 1}, true},
		{Pixel{57, 17}, GridPosition{1, 1}, true},
		{Pixel{58, 17}, GridPosition{1, 1}, true},
		{Pixel{59, 18}, GridPosition{1, 1}, true},
		{Pixel{60, 18}, GridPosition{1, 1}, true},
		{Pixel{61, 18}, GridPosition{1, 1}, true},
		{Pixel{62, 19}, GridPosition{1, 1}, true},
		{Pixel{63, 19}, GridPosition{1, 1}, true},
		{Pixel{64, 20}, GridPosition{1, 1}, true},

		// second hexagon y section (full rectangle of hexagon 1,1)

		// right side of hexagon 0,1
		{Pixel{18, 20}, GridPosition{0, 1}, true},
		{Pixel{18, 21}, GridPosition{0, 1}, true},
		{Pixel{18, 22}, GridPosition{0, 1}, true},
		{Pixel{18, 23}, GridPosition{0, 1}, true},
		{Pixel{18, 24}, GridPosition{0, 1}, true},
		{Pixel{18, 25}, GridPosition{0, 1}, true},
		{Pixel{18, 26}, GridPosition{0, 1}, true},
		{Pixel{18, 27}, GridPosition{0, 1}, true},
		{Pixel{18, 28}, GridPosition{0, 1}, true},
		{Pixel{18, 29}, GridPosition{0, 1}, true},
		{Pixel{18, 30}, GridPosition{0, 1}, true},
		{Pixel{18, 31}, GridPosition{0, 1}, true},
		{Pixel{18, 32}, GridPosition{0, 1}, true},
		{Pixel{18, 33}, GridPosition{0, 1}, true},
		{Pixel{18, 34}, GridPosition{0, 1}, true},
		{Pixel{18, 35}, GridPosition{0, 1}, true},
		{Pixel{18, 36}, GridPosition{0, 1}, true},
		{Pixel{18, 37}, GridPosition{0, 1}, true},
		{Pixel{18, 38}, GridPosition{0, 1}, true},
		{Pixel{18, 39}, GridPosition{0, 1}, true},

		// left side of hexagon 1,1
		// hexagons are right side inclusive and left side exclusive
		{Pixel{19, 20}, GridPosition{1, 1}, true},
		{Pixel{19, 21}, GridPosition{1, 1}, true},
		{Pixel{19, 22}, GridPosition{1, 1}, true},
		{Pixel{19, 23}, GridPosition{1, 1}, true},
		{Pixel{19, 24}, GridPosition{1, 1}, true},
		{Pixel{19, 25}, GridPosition{1, 1}, true},
		{Pixel{19, 26}, GridPosition{1, 1}, true},
		{Pixel{19, 27}, GridPosition{1, 1}, true},
		{Pixel{19, 28}, GridPosition{1, 1}, true},
		{Pixel{19, 29}, GridPosition{1, 1}, true},
		{Pixel{19, 30}, GridPosition{1, 1}, true},
		{Pixel{19, 31}, GridPosition{1, 1}, true},
		{Pixel{19, 32}, GridPosition{1, 1}, true},
		{Pixel{19, 33}, GridPosition{1, 1}, true},
		{Pixel{19, 34}, GridPosition{1, 1}, true},
		{Pixel{19, 35}, GridPosition{1, 1}, true},
		{Pixel{19, 36}, GridPosition{1, 1}, true},
		{Pixel{19, 37}, GridPosition{1, 1}, true},
		{Pixel{19, 38}, GridPosition{1, 1}, true},
		{Pixel{19, 39}, GridPosition{1, 1}, true},
		// continue left side of hexagon 1,1
		{Pixel{20, 20}, GridPosition{1, 1}, true},
		{Pixel{20, 21}, GridPosition{1, 1}, true},
		{Pixel{20, 22}, GridPosition{1, 1}, true},
		{Pixel{20, 23}, GridPosition{1, 1}, true},
		{Pixel{20, 24}, GridPosition{1, 1}, true},
		{Pixel{20, 25}, GridPosition{1, 1}, true},
		{Pixel{20, 26}, GridPosition{1, 1}, true},
		{Pixel{20, 27}, GridPosition{1, 1}, true},
		{Pixel{20, 28}, GridPosition{1, 1}, true},
		{Pixel{20, 29}, GridPosition{1, 1}, true},
		{Pixel{20, 30}, GridPosition{1, 1}, true},
		{Pixel{20, 31}, GridPosition{1, 1}, true},
		{Pixel{20, 32}, GridPosition{1, 1}, true},
		{Pixel{20, 33}, GridPosition{1, 1}, true},
		{Pixel{20, 34}, GridPosition{1, 1}, true},
		{Pixel{20, 35}, GridPosition{1, 1}, true},
		{Pixel{20, 36}, GridPosition{1, 1}, true},
		{Pixel{20, 37}, GridPosition{1, 1}, true},
		{Pixel{20, 38}, GridPosition{1, 1}, true},
		{Pixel{20, 39}, GridPosition{1, 1}, true},

		// right side of hexagon 1,1
		{Pixel{64, 20}, GridPosition{1, 1}, true},
		{Pixel{64, 21}, GridPosition{1, 1}, true},
		{Pixel{64, 22}, GridPosition{1, 1}, true},
		{Pixel{64, 23}, GridPosition{1, 1}, true},
		{Pixel{64, 24}, GridPosition{1, 1}, true},
		{Pixel{64, 25}, GridPosition{1, 1}, true},
		{Pixel{64, 26}, GridPosition{1, 1}, true},
		{Pixel{64, 27}, GridPosition{1, 1}, true},
		{Pixel{64, 28}, GridPosition{1, 1}, true},
		{Pixel{64, 29}, GridPosition{1, 1}, true},
		{Pixel{64, 30}, GridPosition{1, 1}, true},
		{Pixel{64, 31}, GridPosition{1, 1}, true},
		{Pixel{64, 32}, GridPosition{1, 1}, true},
		{Pixel{64, 33}, GridPosition{1, 1}, true},
		{Pixel{64, 34}, GridPosition{1, 1}, true},
		{Pixel{64, 35}, GridPosition{1, 1}, true},
		{Pixel{64, 36}, GridPosition{1, 1}, true},
		{Pixel{64, 37}, GridPosition{1, 1}, true},
		{Pixel{64, 38}, GridPosition{1, 1}, true},
		{Pixel{64, 39}, GridPosition{1, 1}, true},

		// left side of hexagon 2,1
		// hexagons are right side inclusive and left side exclusive
		{Pixel{65, 20}, GridPosition{2, 1}, true},
		{Pixel{65, 21}, GridPosition{2, 1}, true},
		{Pixel{65, 22}, GridPosition{2, 1}, true},
		{Pixel{65, 23}, GridPosition{2, 1}, true},
		{Pixel{65, 24}, GridPosition{2, 1}, true},
		{Pixel{65, 25}, GridPosition{2, 1}, true},
		{Pixel{65, 26}, GridPosition{2, 1}, true},
		{Pixel{65, 27}, GridPosition{2, 1}, true},
		{Pixel{65, 28}, GridPosition{2, 1}, true},
		{Pixel{65, 29}, GridPosition{2, 1}, true},
		{Pixel{65, 30}, GridPosition{2, 1}, true},
		{Pixel{65, 31}, GridPosition{2, 1}, true},
		{Pixel{65, 32}, GridPosition{2, 1}, true},
		{Pixel{65, 33}, GridPosition{2, 1}, true},
		{Pixel{65, 34}, GridPosition{2, 1}, true},
		{Pixel{65, 35}, GridPosition{2, 1}, true},
		{Pixel{65, 36}, GridPosition{2, 1}, true},
		{Pixel{65, 37}, GridPosition{2, 1}, true},
		{Pixel{65, 38}, GridPosition{2, 1}, true},
		{Pixel{65, 39}, GridPosition{2, 1}, true},
		// continue left side of hexagon 2,1
		{Pixel{66, 20}, GridPosition{2, 1}, true},
		{Pixel{66, 21}, GridPosition{2, 1}, true},
		{Pixel{66, 22}, GridPosition{2, 1}, true},
		{Pixel{66, 23}, GridPosition{2, 1}, true},
		{Pixel{66, 24}, GridPosition{2, 1}, true},
		{Pixel{66, 25}, GridPosition{2, 1}, true},
		{Pixel{66, 26}, GridPosition{2, 1}, true},
		{Pixel{66, 27}, GridPosition{2, 1}, true},
		{Pixel{66, 28}, GridPosition{2, 1}, true},
		{Pixel{66, 29}, GridPosition{2, 1}, true},
		{Pixel{66, 30}, GridPosition{2, 1}, true},
		{Pixel{66, 31}, GridPosition{2, 1}, true},
		{Pixel{66, 32}, GridPosition{2, 1}, true},
		{Pixel{66, 33}, GridPosition{2, 1}, true},
		{Pixel{66, 34}, GridPosition{2, 1}, true},
		{Pixel{66, 35}, GridPosition{2, 1}, true},
		{Pixel{66, 36}, GridPosition{2, 1}, true},
		{Pixel{66, 37}, GridPosition{2, 1}, true},
		{Pixel{66, 38}, GridPosition{2, 1}, true},
		{Pixel{66, 39}, GridPosition{2, 1}, true},

		// third hexagon y section (inverted hexagon look up table)
		// between hexagons 0,2 1,1 and 1,2

		// sanity checks
		{Pixel{19, 49}, GridPosition{0, 2}, true},
		{Pixel{65, 49}, GridPosition{1, 2}, true},
		{Pixel{41, 30}, GridPosition{1, 1}, true},
		{Pixel{42, 30}, GridPosition{1, 1}, true},
		{Pixel{43, 30}, GridPosition{1, 1}, true},

		// 0,2 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{19, 40}, GridPosition{0, 2}, true},
		{Pixel{20, 41}, GridPosition{0, 2}, true},
		{Pixel{21, 41}, GridPosition{0, 2}, true},
		{Pixel{22, 42}, GridPosition{0, 2}, true},
		{Pixel{23, 42}, GridPosition{0, 2}, true},
		{Pixel{24, 42}, GridPosition{0, 2}, true},
		{Pixel{25, 43}, GridPosition{0, 2}, true},
		{Pixel{26, 43}, GridPosition{0, 2}, true},
		{Pixel{27, 44}, GridPosition{0, 2}, true},
		{Pixel{28, 44}, GridPosition{0, 2}, true},
		{Pixel{29, 45}, GridPosition{0, 2}, true},
		{Pixel{30, 45}, GridPosition{0, 2}, true},
		{Pixel{31, 45}, GridPosition{0, 2}, true},
		{Pixel{32, 46}, GridPosition{0, 2}, true},
		{Pixel{33, 46}, GridPosition{0, 2}, true},
		{Pixel{34, 47}, GridPosition{0, 2}, true},
		{Pixel{35, 47}, GridPosition{0, 2}, true},
		{Pixel{36, 48}, GridPosition{0, 2}, true},
		{Pixel{37, 48}, GridPosition{0, 2}, true},
		{Pixel{38, 48}, GridPosition{0, 2}, true},
		{Pixel{39, 49}, GridPosition{0, 2}, true},
		{Pixel{40, 49}, GridPosition{0, 2}, true},
		{Pixel{41, 50}, GridPosition{0, 2}, true},

		// 1, 1 hexagon slope check between itself and 0, 2 hexagon
		{Pixel{19, 39}, GridPosition{1, 1}, true},
		{Pixel{20, 40}, GridPosition{1, 1}, true},
		{Pixel{21, 40}, GridPosition{1, 1}, true},
		{Pixel{22, 41}, GridPosition{1, 1}, true},
		{Pixel{23, 41}, GridPosition{1, 1}, true},
		{Pixel{24, 41}, GridPosition{1, 1}, true},
		{Pixel{25, 42}, GridPosition{1, 1}, true},
		{Pixel{26, 42}, GridPosition{1, 1}, true},
		{Pixel{27, 43}, GridPosition{1, 1}, true},
		{Pixel{28, 43}, GridPosition{1, 1}, true},
		{Pixel{29, 44}, GridPosition{1, 1}, true},
		{Pixel{30, 44}, GridPosition{1, 1}, true},
		{Pixel{31, 44}, GridPosition{1, 1}, true},
		{Pixel{32, 45}, GridPosition{1, 1}, true},
		{Pixel{33, 45}, GridPosition{1, 1}, true},
		{Pixel{34, 46}, GridPosition{1, 1}, true},
		{Pixel{35, 46}, GridPosition{1, 1}, true},
		{Pixel{36, 47}, GridPosition{1, 1}, true},
		{Pixel{37, 47}, GridPosition{1, 1}, true},
		{Pixel{38, 47}, GridPosition{1, 1}, true},
		{Pixel{39, 48}, GridPosition{1, 1}, true},
		{Pixel{40, 48}, GridPosition{1, 1}, true},
		{Pixel{41, 49}, GridPosition{1, 1}, true},

		// 1,1 hexagon slope check between itself and 1, 2 hexagon
		{Pixel{42, 49}, GridPosition{1, 1}, true},
		{Pixel{43, 48}, GridPosition{1, 1}, true},
		{Pixel{44, 48}, GridPosition{1, 1}, true},
		{Pixel{45, 47}, GridPosition{1, 1}, true},
		{Pixel{46, 47}, GridPosition{1, 1}, true},
		{Pixel{47, 47}, GridPosition{1, 1}, true},
		{Pixel{48, 46}, GridPosition{1, 1}, true},
		{Pixel{49, 46}, GridPosition{1, 1}, true},
		{Pixel{50, 45}, GridPosition{1, 1}, true},
		{Pixel{51, 45}, GridPosition{1, 1}, true},
		{Pixel{52, 44}, GridPosition{1, 1}, true},
		{Pixel{53, 44}, GridPosition{1, 1}, true},
		{Pixel{54, 44}, GridPosition{1, 1}, true},
		{Pixel{55, 43}, GridPosition{1, 1}, true},
		{Pixel{56, 43}, GridPosition{1, 1}, true},
		{Pixel{57, 42}, GridPosition{1, 1}, true},
		{Pixel{58, 42}, GridPosition{1, 1}, true},
		{Pixel{59, 41}, GridPosition{1, 1}, true},
		{Pixel{60, 41}, GridPosition{1, 1}, true},
		{Pixel{61, 41}, GridPosition{1, 1}, true},
		{Pixel{62, 40}, GridPosition{1, 1}, true},
		{Pixel{63, 40}, GridPosition{1, 1}, true},
		{Pixel{64, 39}, GridPosition{1, 1}, true},

		// 1,2 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{42, 50}, GridPosition{1, 2}, true},
		{Pixel{43, 49}, GridPosition{1, 2}, true},
		{Pixel{44, 49}, GridPosition{1, 2}, true},
		{Pixel{45, 48}, GridPosition{1, 2}, true},
		{Pixel{46, 48}, GridPosition{1, 2}, true},
		{Pixel{47, 48}, GridPosition{1, 2}, true},
		{Pixel{48, 47}, GridPosition{1, 2}, true},
		{Pixel{49, 47}, GridPosition{1, 2}, true},
		{Pixel{50, 46}, GridPosition{1, 2}, true},
		{Pixel{51, 46}, GridPosition{1, 2}, true},
		{Pixel{52, 45}, GridPosition{1, 2}, true},
		{Pixel{53, 45}, GridPosition{1, 2}, true},
		{Pixel{54, 45}, GridPosition{1, 2}, true},
		{Pixel{55, 44}, GridPosition{1, 2}, true},
		{Pixel{56, 44}, GridPosition{1, 2}, true},
		{Pixel{57, 43}, GridPosition{1, 2}, true},
		{Pixel{58, 43}, GridPosition{1, 2}, true},
		{Pixel{59, 42}, GridPosition{1, 2}, true},
		{Pixel{60, 42}, GridPosition{1, 2}, true},
		{Pixel{61, 42}, GridPosition{1, 2}, true},
		{Pixel{62, 41}, GridPosition{1, 2}, true},
		{Pixel{63, 41}, GridPosition{1, 2}, true},
		{Pixel{64, 40}, GridPosition{1, 2}, true},

		// fourth hexagon y section (split rectangle between hexagon 0,2 and 1,2)

		// middle of hexagon 0,2
		{Pixel{18, 50}, GridPosition{0, 2}, true},
		{Pixel{18, 51}, GridPosition{0, 2}, true},
		{Pixel{18, 52}, GridPosition{0, 2}, true},
		{Pixel{18, 53}, GridPosition{0, 2}, true},
		{Pixel{18, 54}, GridPosition{0, 2}, true},
		{Pixel{18, 55}, GridPosition{0, 2}, true},
		{Pixel{18, 56}, GridPosition{0, 2}, true},
		{Pixel{18, 57}, GridPosition{0, 2}, true},
		{Pixel{18, 58}, GridPosition{0, 2}, true},
		{Pixel{18, 59}, GridPosition{0, 2}, true},
		{Pixel{18, 60}, GridPosition{0, 2}, true},
		{Pixel{18, 61}, GridPosition{0, 2}, true},
		{Pixel{18, 62}, GridPosition{0, 2}, true},
		{Pixel{18, 63}, GridPosition{0, 2}, true},
		{Pixel{18, 64}, GridPosition{0, 2}, true},
		{Pixel{18, 65}, GridPosition{0, 2}, true},
		{Pixel{18, 66}, GridPosition{0, 2}, true},
		{Pixel{18, 67}, GridPosition{0, 2}, true},
		{Pixel{18, 68}, GridPosition{0, 2}, true},
		{Pixel{18, 69}, GridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{Pixel{19, 50}, GridPosition{0, 2}, true},
		{Pixel{19, 51}, GridPosition{0, 2}, true},
		{Pixel{19, 52}, GridPosition{0, 2}, true},
		{Pixel{19, 53}, GridPosition{0, 2}, true},
		{Pixel{19, 54}, GridPosition{0, 2}, true},
		{Pixel{19, 55}, GridPosition{0, 2}, true},
		{Pixel{19, 56}, GridPosition{0, 2}, true},
		{Pixel{19, 57}, GridPosition{0, 2}, true},
		{Pixel{19, 58}, GridPosition{0, 2}, true},
		{Pixel{19, 59}, GridPosition{0, 2}, true},
		{Pixel{19, 60}, GridPosition{0, 2}, true},
		{Pixel{19, 61}, GridPosition{0, 2}, true},
		{Pixel{19, 62}, GridPosition{0, 2}, true},
		{Pixel{19, 63}, GridPosition{0, 2}, true},
		{Pixel{19, 64}, GridPosition{0, 2}, true},
		{Pixel{19, 65}, GridPosition{0, 2}, true},
		{Pixel{19, 66}, GridPosition{0, 2}, true},
		{Pixel{19, 67}, GridPosition{0, 2}, true},
		{Pixel{19, 68}, GridPosition{0, 2}, true},
		{Pixel{19, 69}, GridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{Pixel{20, 50}, GridPosition{0, 2}, true},
		{Pixel{20, 51}, GridPosition{0, 2}, true},
		{Pixel{20, 52}, GridPosition{0, 2}, true},
		{Pixel{20, 53}, GridPosition{0, 2}, true},
		{Pixel{20, 54}, GridPosition{0, 2}, true},
		{Pixel{20, 55}, GridPosition{0, 2}, true},
		{Pixel{20, 56}, GridPosition{0, 2}, true},
		{Pixel{20, 57}, GridPosition{0, 2}, true},
		{Pixel{20, 58}, GridPosition{0, 2}, true},
		{Pixel{20, 59}, GridPosition{0, 2}, true},
		{Pixel{20, 60}, GridPosition{0, 2}, true},
		{Pixel{20, 61}, GridPosition{0, 2}, true},
		{Pixel{20, 62}, GridPosition{0, 2}, true},
		{Pixel{20, 63}, GridPosition{0, 2}, true},
		{Pixel{20, 64}, GridPosition{0, 2}, true},
		{Pixel{20, 65}, GridPosition{0, 2}, true},
		{Pixel{20, 66}, GridPosition{0, 2}, true},
		{Pixel{20, 67}, GridPosition{0, 2}, true},
		{Pixel{20, 68}, GridPosition{0, 2}, true},
		{Pixel{20, 69}, GridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		{Pixel{41, 50}, GridPosition{0, 2}, true},
		{Pixel{41, 51}, GridPosition{0, 2}, true},
		{Pixel{41, 52}, GridPosition{0, 2}, true},
		{Pixel{41, 53}, GridPosition{0, 2}, true},
		{Pixel{41, 54}, GridPosition{0, 2}, true},
		{Pixel{41, 55}, GridPosition{0, 2}, true},
		{Pixel{41, 56}, GridPosition{0, 2}, true},
		{Pixel{41, 57}, GridPosition{0, 2}, true},
		{Pixel{41, 58}, GridPosition{0, 2}, true},
		{Pixel{41, 59}, GridPosition{0, 2}, true},
		{Pixel{41, 60}, GridPosition{0, 2}, true},
		{Pixel{41, 61}, GridPosition{0, 2}, true},
		{Pixel{41, 62}, GridPosition{0, 2}, true},
		{Pixel{41, 63}, GridPosition{0, 2}, true},
		{Pixel{41, 64}, GridPosition{0, 2}, true},
		{Pixel{41, 65}, GridPosition{0, 2}, true},
		{Pixel{41, 66}, GridPosition{0, 2}, true},
		{Pixel{41, 67}, GridPosition{0, 2}, true},
		{Pixel{41, 68}, GridPosition{0, 2}, true},
		{Pixel{41, 69}, GridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		// hexagons are right left side inclusive, right side exclusive
		{Pixel{42, 50}, GridPosition{1, 2}, true},
		{Pixel{42, 51}, GridPosition{1, 2}, true},
		{Pixel{42, 52}, GridPosition{1, 2}, true},
		{Pixel{42, 53}, GridPosition{1, 2}, true},
		{Pixel{42, 54}, GridPosition{1, 2}, true},
		{Pixel{42, 55}, GridPosition{1, 2}, true},
		{Pixel{42, 56}, GridPosition{1, 2}, true},
		{Pixel{42, 57}, GridPosition{1, 2}, true},
		{Pixel{42, 58}, GridPosition{1, 2}, true},
		{Pixel{42, 59}, GridPosition{1, 2}, true},
		{Pixel{42, 60}, GridPosition{1, 2}, true},
		{Pixel{42, 61}, GridPosition{1, 2}, true},
		{Pixel{42, 62}, GridPosition{1, 2}, true},
		{Pixel{42, 63}, GridPosition{1, 2}, true},
		{Pixel{42, 64}, GridPosition{1, 2}, true},
		{Pixel{42, 65}, GridPosition{1, 2}, true},
		{Pixel{42, 66}, GridPosition{1, 2}, true},
		{Pixel{42, 67}, GridPosition{1, 2}, true},
		{Pixel{42, 68}, GridPosition{1, 2}, true},
		{Pixel{42, 69}, GridPosition{1, 2}, true},

		// border of hexagon 0,2 and 1,2
		{Pixel{43, 50}, GridPosition{1, 2}, true},
		{Pixel{43, 51}, GridPosition{1, 2}, true},
		{Pixel{43, 52}, GridPosition{1, 2}, true},
		{Pixel{43, 53}, GridPosition{1, 2}, true},
		{Pixel{43, 54}, GridPosition{1, 2}, true},
		{Pixel{43, 55}, GridPosition{1, 2}, true},
		{Pixel{43, 56}, GridPosition{1, 2}, true},
		{Pixel{43, 57}, GridPosition{1, 2}, true},
		{Pixel{43, 58}, GridPosition{1, 2}, true},
		{Pixel{43, 59}, GridPosition{1, 2}, true},
		{Pixel{43, 60}, GridPosition{1, 2}, true},
		{Pixel{43, 61}, GridPosition{1, 2}, true},
		{Pixel{43, 62}, GridPosition{1, 2}, true},
		{Pixel{43, 63}, GridPosition{1, 2}, true},
		{Pixel{43, 64}, GridPosition{1, 2}, true},
		{Pixel{43, 65}, GridPosition{1, 2}, true},
		{Pixel{43, 66}, GridPosition{1, 2}, true},
		{Pixel{43, 67}, GridPosition{1, 2}, true},
		{Pixel{43, 68}, GridPosition{1, 2}, true},
		{Pixel{43, 69}, GridPosition{1, 2}, true},

		// right side of hexagon 1,2
		{Pixel{64, 50}, GridPosition{1, 2}, true},
		{Pixel{64, 51}, GridPosition{1, 2}, true},
		{Pixel{64, 52}, GridPosition{1, 2}, true},
		{Pixel{64, 53}, GridPosition{1, 2}, true},
		{Pixel{64, 54}, GridPosition{1, 2}, true},
		{Pixel{64, 55}, GridPosition{1, 2}, true},
		{Pixel{64, 56}, GridPosition{1, 2}, true},
		{Pixel{64, 57}, GridPosition{1, 2}, true},
		{Pixel{64, 58}, GridPosition{1, 2}, true},
		{Pixel{64, 59}, GridPosition{1, 2}, true},
		{Pixel{64, 60}, GridPosition{1, 2}, true},
		{Pixel{64, 61}, GridPosition{1, 2}, true},
		{Pixel{64, 62}, GridPosition{1, 2}, true},
		{Pixel{64, 63}, GridPosition{1, 2}, true},
		{Pixel{64, 64}, GridPosition{1, 2}, true},
		{Pixel{64, 65}, GridPosition{1, 2}, true},
		{Pixel{64, 66}, GridPosition{1, 2}, true},
		{Pixel{64, 67}, GridPosition{1, 2}, true},
		{Pixel{64, 68}, GridPosition{1, 2}, true},
		{Pixel{64, 69}, GridPosition{1, 2}, true},

		// continue right side of hexagon 1,2
		{Pixel{65, 50}, GridPosition{1, 2}, true},
		{Pixel{65, 51}, GridPosition{1, 2}, true},
		{Pixel{65, 52}, GridPosition{1, 2}, true},
		{Pixel{65, 53}, GridPosition{1, 2}, true},
		{Pixel{65, 54}, GridPosition{1, 2}, true},
		{Pixel{65, 55}, GridPosition{1, 2}, true},
		{Pixel{65, 56}, GridPosition{1, 2}, true},
		{Pixel{65, 57}, GridPosition{1, 2}, true},
		{Pixel{65, 58}, GridPosition{1, 2}, true},
		{Pixel{65, 59}, GridPosition{1, 2}, true},
		{Pixel{65, 60}, GridPosition{1, 2}, true},
		{Pixel{65, 61}, GridPosition{1, 2}, true},
		{Pixel{65, 62}, GridPosition{1, 2}, true},
		{Pixel{65, 63}, GridPosition{1, 2}, true},
		{Pixel{65, 64}, GridPosition{1, 2}, true},
		{Pixel{65, 65}, GridPosition{1, 2}, true},
		{Pixel{65, 66}, GridPosition{1, 2}, true},
		{Pixel{65, 67}, GridPosition{1, 2}, true},
		{Pixel{65, 68}, GridPosition{1, 2}, true},
		{Pixel{65, 69}, GridPosition{1, 2}, true},

		// continue right side of hexagon hexagon 1,2
		{Pixel{66, 50}, GridPosition{1, 2}, true},
		{Pixel{66, 51}, GridPosition{1, 2}, true},
		{Pixel{66, 52}, GridPosition{1, 2}, true},
		{Pixel{66, 53}, GridPosition{1, 2}, true},
		{Pixel{66, 54}, GridPosition{1, 2}, true},
		{Pixel{66, 55}, GridPosition{1, 2}, true},
		{Pixel{66, 56}, GridPosition{1, 2}, true},
		{Pixel{66, 57}, GridPosition{1, 2}, true},
		{Pixel{66, 58}, GridPosition{1, 2}, true},
		{Pixel{66, 59}, GridPosition{1, 2}, true},
		{Pixel{66, 60}, GridPosition{1, 2}, true},
		{Pixel{66, 61}, GridPosition{1, 2}, true},
		{Pixel{66, 62}, GridPosition{1, 2}, true},
		{Pixel{66, 63}, GridPosition{1, 2}, true},
		{Pixel{66, 64}, GridPosition{1, 2}, true},
		{Pixel{66, 65}, GridPosition{1, 2}, true},
		{Pixel{66, 66}, GridPosition{1, 2}, true},
		{Pixel{66, 67}, GridPosition{1, 2}, true},
		{Pixel{66, 68}, GridPosition{1, 2}, true},
		{Pixel{66, 69}, GridPosition{1, 2}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"Pixel: %s converts to hex grid position %s and want it to %s",
			tt.pixel.log(),
			tt.wantPos.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.getGridPos(tt.pixel)
			if ok != tt.wantOk {
				t.Errorf("Want pixel position to %s", sWantOk)
			}
			if !pos.equals(tt.wantPos) {
				t.Errorf("For pixel position %s, Want grid position %s but it was %s",
					tt.pixel.log(),
					tt.wantPos.log(),
					pos.log())
			}
		})
	}
}

func TestAdjacent(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		gridPos     GridPosition
		hexSide     hexSide
		wantGridPos GridPosition
		wantOk      bool
	}{
		// check top left hexagon {0,0}
		{GridPosition{0, 0}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{0, 0}, hsEast, GridPosition{1, 0}, true},
		{GridPosition{0, 0}, hsSoutheast, GridPosition{1, 1}, true},
		{GridPosition{0, 0}, hsSouthwest, GridPosition{0, 1}, true},
		{GridPosition{0, 0}, hsWest, GridPosition{0, 0}, false},
		{GridPosition{0, 0}, hsNorthwest, GridPosition{0, 0}, false},

		// check top left hexagon {1, 0}
		{GridPosition{1, 0}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{1, 0}, hsEast, GridPosition{2, 0}, true},
		{GridPosition{1, 0}, hsSoutheast, GridPosition{2, 1}, true},
		{GridPosition{1, 0}, hsSouthwest, GridPosition{1, 1}, true},
		{GridPosition{1, 0}, hsWest, GridPosition{0, 0}, true},
		{GridPosition{1, 0}, hsNorthwest, GridPosition{0, 0}, false},

		// check top left hexagon {0, 1}
		{GridPosition{0, 1}, hsNortheast, GridPosition{0, 0}, true},
		{GridPosition{0, 1}, hsEast, GridPosition{1, 1}, true},
		{GridPosition{0, 1}, hsSoutheast, GridPosition{0, 2}, true},
		{GridPosition{0, 1}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{0, 1}, hsWest, GridPosition{0, 0}, false},
		{GridPosition{0, 1}, hsNorthwest, GridPosition{0, 0}, false},

		// check bottom left hexagon {0,24}
		{GridPosition{0, 24}, hsNortheast, GridPosition{1, 23}, true},
		{GridPosition{0, 24}, hsEast, GridPosition{1, 24}, true},
		{GridPosition{0, 24}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{0, 24}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{0, 24}, hsWest, GridPosition{0, 0}, false},
		{GridPosition{0, 24}, hsNorthwest, GridPosition{0, 23}, true},

		// check bottom left hexagon {1,24}
		{GridPosition{1, 24}, hsNortheast, GridPosition{2, 23}, true},
		{GridPosition{1, 24}, hsEast, GridPosition{2, 24}, true},
		{GridPosition{1, 24}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{1, 24}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{1, 24}, hsWest, GridPosition{0, 24}, true},
		{GridPosition{1, 24}, hsNorthwest, GridPosition{1, 23}, true},

		// check bottom left hexagon {0,23}
		{GridPosition{0, 23}, hsNortheast, GridPosition{0, 22}, true},
		{GridPosition{0, 23}, hsEast, GridPosition{1, 23}, true},
		{GridPosition{0, 23}, hsSoutheast, GridPosition{0, 24}, true},
		{GridPosition{0, 23}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{0, 23}, hsWest, GridPosition{0, 0}, false},
		{GridPosition{0, 23}, hsNorthwest, GridPosition{0, 0}, false},

		// check top right hexagon {27,0}
		{GridPosition{27, 0}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{27, 0}, hsEast, GridPosition{0, 0}, false},
		{GridPosition{27, 0}, hsSoutheast, GridPosition{28, 1}, true},
		{GridPosition{27, 0}, hsSouthwest, GridPosition{27, 1}, true},
		{GridPosition{27, 0}, hsWest, GridPosition{26, 0}, true},
		{GridPosition{27, 0}, hsNorthwest, GridPosition{0, 0}, false},

		// check top right hexagon {28,1}
		{GridPosition{28, 1}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{28, 1}, hsEast, GridPosition{0, 0}, false},
		{GridPosition{28, 1}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{28, 1}, hsSouthwest, GridPosition{27, 2}, true},
		{GridPosition{28, 1}, hsWest, GridPosition{27, 1}, true},
		{GridPosition{28, 1}, hsNorthwest, GridPosition{27, 0}, true},

		// check top right hexagon {26,0}
		{GridPosition{26, 0}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{26, 0}, hsEast, GridPosition{27, 0}, true},
		{GridPosition{26, 0}, hsSoutheast, GridPosition{27, 1}, true},
		{GridPosition{26, 0}, hsSouthwest, GridPosition{26, 1}, true},
		{GridPosition{26, 0}, hsWest, GridPosition{25, 0}, true},
		{GridPosition{26, 0}, hsNorthwest, GridPosition{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{GridPosition{27, 24}, hsNortheast, GridPosition{28, 23}, true},
		{GridPosition{27, 24}, hsEast, GridPosition{0, 0}, false},
		{GridPosition{27, 24}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{27, 24}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{27, 24}, hsWest, GridPosition{26, 24}, true},
		{GridPosition{27, 24}, hsNorthwest, GridPosition{27, 23}, true},

		// check bottom right hexagon {26, 24}
		{GridPosition{26, 24}, hsNortheast, GridPosition{27, 23}, true},
		{GridPosition{26, 24}, hsEast, GridPosition{27, 24}, true},
		{GridPosition{26, 24}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{26, 24}, hsSouthwest, GridPosition{0, 0}, false},
		{GridPosition{26, 24}, hsWest, GridPosition{25, 24}, true},
		{GridPosition{26, 24}, hsNorthwest, GridPosition{26, 23}, true},

		// check bottom right hexagon {28, 23}
		{GridPosition{28, 23}, hsNortheast, GridPosition{0, 0}, false},
		{GridPosition{28, 23}, hsEast, GridPosition{0, 0}, false},
		{GridPosition{28, 23}, hsSoutheast, GridPosition{0, 0}, false},
		{GridPosition{28, 23}, hsSouthwest, GridPosition{27, 24}, true},
		{GridPosition{28, 23}, hsWest, GridPosition{27, 23}, true},
		{GridPosition{28, 23}, hsNorthwest, GridPosition{27, 22}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"gridPos: %s and hex side %s has adjacent position %s and want it to %s",
			tt.gridPos.log(),
			tt.hexSide.string(),
			tt.wantGridPos.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.adjacent(tt.gridPos, tt.hexSide)
			if ok != tt.wantOk {
				t.Errorf("Want pixel position to %s", sWantOk)
			}
			if !pos.equals(tt.wantGridPos) {
				t.Errorf("For grid pos %s and hex side %s Want grid position %s but it was %s",
					tt.gridPos.log(),
					tt.hexSide.string(),
					tt.wantGridPos.log(),
					pos.log())
			}
		})
	}

}

func TestDemotif(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		pixel       Pixel
		wantDemotif demotif
		wantOk      bool
	}{
		// hexagon points near 19, 20 pixel position
		{Pixel{18, 19}, demotif{Pixel{22, 39}}, true},
		{Pixel{18, 20}, demotif{Pixel{45, 10}}, true},
		{Pixel{18, 21}, demotif{Pixel{45, 11}}, true},
		{Pixel{19, 19}, demotif{Pixel{23, 39}}, true},
		{Pixel{19, 20}, demotif{Pixel{0, 10}}, true},
		{Pixel{19, 21}, demotif{Pixel{0, 11}}, true},
		{Pixel{20, 19}, demotif{Pixel{1, 9}}, true},
		{Pixel{20, 20}, demotif{Pixel{1, 10}}, true},
		{Pixel{20, 21}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 19, 40 pixel position
		{Pixel{18, 39}, demotif{Pixel{45, 29}}, true},
		{Pixel{18, 40}, demotif{Pixel{22, 0}}, true},
		{Pixel{18, 41}, demotif{Pixel{22, 1}}, true},
		{Pixel{19, 39}, demotif{Pixel{0, 29}}, true},
		{Pixel{19, 40}, demotif{Pixel{23, 0}}, true},
		{Pixel{19, 41}, demotif{Pixel{23, 1}}, true},
		{Pixel{20, 39}, demotif{Pixel{1, 29}}, true},
		{Pixel{20, 40}, demotif{Pixel{1, 30}}, true},
		{Pixel{20, 41}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 42, 10 pixel position
		{Pixel{41, 9}, demotif{Pixel{45, 29}}, true},
		{Pixel{41, 10}, demotif{Pixel{22, 0}}, true},
		{Pixel{41, 11}, demotif{Pixel{22, 1}}, true},
		{Pixel{42, 9}, demotif{Pixel{0, 29}}, true},
		{Pixel{42, 10}, demotif{Pixel{23, 0}}, true},
		{Pixel{42, 11}, demotif{Pixel{23, 1}}, true},
		{Pixel{43, 9}, demotif{Pixel{1, 29}}, true},
		{Pixel{43, 10}, demotif{Pixel{1, 30}}, true},
		{Pixel{43, 11}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 42, 50 pixel position
		{Pixel{41, 49}, demotif{Pixel{22, 39}}, true},
		{Pixel{41, 50}, demotif{Pixel{45, 10}}, true},
		{Pixel{41, 51}, demotif{Pixel{45, 11}}, true},
		{Pixel{42, 49}, demotif{Pixel{23, 39}}, true},
		{Pixel{42, 50}, demotif{Pixel{0, 10}}, true},
		{Pixel{42, 51}, demotif{Pixel{0, 11}}, true},
		{Pixel{43, 49}, demotif{Pixel{1, 9}}, true},
		{Pixel{43, 50}, demotif{Pixel{1, 10}}, true},
		{Pixel{43, 51}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 65, 20 pixel position
		{Pixel{64, 19}, demotif{Pixel{22, 39}}, true},
		{Pixel{64, 20}, demotif{Pixel{45, 10}}, true},
		{Pixel{64, 21}, demotif{Pixel{45, 11}}, true},
		{Pixel{65, 19}, demotif{Pixel{23, 39}}, true},
		{Pixel{65, 20}, demotif{Pixel{0, 10}}, true},
		{Pixel{65, 21}, demotif{Pixel{0, 11}}, true},
		{Pixel{66, 19}, demotif{Pixel{1, 9}}, true},
		{Pixel{66, 20}, demotif{Pixel{1, 10}}, true},
		{Pixel{66, 21}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 65, 40 pixel position
		{Pixel{64, 39}, demotif{Pixel{45, 29}}, true},
		{Pixel{64, 40}, demotif{Pixel{22, 0}}, true},
		{Pixel{64, 41}, demotif{Pixel{22, 1}}, true},
		{Pixel{65, 39}, demotif{Pixel{0, 29}}, true},
		{Pixel{65, 40}, demotif{Pixel{23, 0}}, true},
		{Pixel{65, 41}, demotif{Pixel{23, 1}}, true},
		{Pixel{66, 39}, demotif{Pixel{1, 29}}, true},
		{Pixel{66, 40}, demotif{Pixel{1, 30}}, true},
		{Pixel{66, 41}, demotif{Pixel{24, 1}}, true},

		// top left of entire hex grid
		{Pixel{0, 0}, demotif{Pixel{4, 20}}, true},
		{Pixel{-1, 0}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, -1}, demotif{Pixel{0, 0}}, false},

		// first left border point hex change around 0,12 pixel position
		{Pixel{-1, 11}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 12}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 13}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 11}, demotif{Pixel{4, 31}}, true},
		{Pixel{0, 12}, demotif{Pixel{27, 2}}, true},
		{Pixel{0, 13}, demotif{Pixel{27, 3}}, true},
		{Pixel{1, 11}, demotif{Pixel{5, 31}}, true},
		{Pixel{1, 12}, demotif{Pixel{28, 2}}, true},
		{Pixel{1, 13}, demotif{Pixel{28, 3}}, true},

		// second left border point hex change around 0,48 pixel position
		{Pixel{-1, 47}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 48}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 49}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 47}, demotif{Pixel{27, 37}}, true},
		{Pixel{0, 48}, demotif{Pixel{4, 8}}, true},
		{Pixel{0, 49}, demotif{Pixel{4, 9}}, true},
		{Pixel{1, 47}, demotif{Pixel{28, 37}}, true},
		{Pixel{1, 48}, demotif{Pixel{5, 8}}, true},
		{Pixel{1, 49}, demotif{Pixel{5, 9}}, true},

		// third left border point hex change around 0,72 pixel position
		{Pixel{-1, 71}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 72}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 73}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 71}, demotif{Pixel{4, 31}}, true},
		{Pixel{0, 72}, demotif{Pixel{27, 2}}, true},
		{Pixel{0, 73}, demotif{Pixel{27, 3}}, true},
		{Pixel{1, 71}, demotif{Pixel{5, 31}}, true},
		{Pixel{1, 72}, demotif{Pixel{28, 2}}, true},
		{Pixel{1, 73}, demotif{Pixel{28, 3}}, true},

		// first top border point hex change around 42,0 pixel position
		{Pixel{41, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{41, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{41, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{42, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{42, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{42, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{43, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{43, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{43, 1}, demotif{Pixel{1, 21}}, true},

		// second top border point hex change around 88, 0 pixel position
		{Pixel{87, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{87, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{87, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{88, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{88, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{88, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{89, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{89, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{89, 1}, demotif{Pixel{1, 21}}, true},

		// third top border point hex change around 134,0 pixel position
		{Pixel{133, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{133, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{133, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{134, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{134, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{134, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{135, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{135, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{135, 1}, demotif{Pixel{1, 21}}, true},

		// hexagon points near 19, 80 pixel position
		{Pixel{18, 79}, demotif{Pixel{22, 39}}, true},
		{Pixel{18, 80}, demotif{Pixel{45, 10}}, true},
		{Pixel{18, 81}, demotif{Pixel{45, 11}}, true},
		{Pixel{19, 79}, demotif{Pixel{23, 39}}, true},
		{Pixel{19, 80}, demotif{Pixel{0, 10}}, true},
		{Pixel{19, 81}, demotif{Pixel{0, 11}}, true},
		{Pixel{20, 79}, demotif{Pixel{1, 9}}, true},
		{Pixel{20, 80}, demotif{Pixel{1, 10}}, true},
		{Pixel{20, 81}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 42, 70 pixel position
		{Pixel{41, 69}, demotif{Pixel{45, 29}}, true},
		{Pixel{41, 70}, demotif{Pixel{22, 0}}, true},
		{Pixel{41, 71}, demotif{Pixel{22, 1}}, true},
		{Pixel{42, 69}, demotif{Pixel{0, 29}}, true},
		{Pixel{42, 70}, demotif{Pixel{23, 0}}, true},
		{Pixel{42, 71}, demotif{Pixel{23, 1}}, true},
		{Pixel{43, 69}, demotif{Pixel{1, 29}}, true},
		{Pixel{43, 70}, demotif{Pixel{1, 30}}, true},
		{Pixel{43, 71}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 65, 80 pixel position
		{Pixel{64, 79}, demotif{Pixel{22, 39}}, true},
		{Pixel{64, 80}, demotif{Pixel{45, 10}}, true},
		{Pixel{64, 81}, demotif{Pixel{45, 11}}, true},
		{Pixel{65, 79}, demotif{Pixel{23, 39}}, true},
		{Pixel{65, 80}, demotif{Pixel{0, 10}}, true},
		{Pixel{65, 81}, demotif{Pixel{0, 11}}, true},
		{Pixel{66, 79}, demotif{Pixel{1, 9}}, true},
		{Pixel{66, 80}, demotif{Pixel{1, 10}}, true},
		{Pixel{66, 81}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 88, 10 pixel position
		{Pixel{87, 9}, demotif{Pixel{45, 29}}, true},
		{Pixel{87, 10}, demotif{Pixel{22, 0}}, true},
		{Pixel{87, 11}, demotif{Pixel{22, 1}}, true},
		{Pixel{88, 9}, demotif{Pixel{0, 29}}, true},
		{Pixel{88, 10}, demotif{Pixel{23, 0}}, true},
		{Pixel{88, 11}, demotif{Pixel{23, 1}}, true},
		{Pixel{89, 9}, demotif{Pixel{1, 29}}, true},
		{Pixel{89, 10}, demotif{Pixel{1, 30}}, true},
		{Pixel{89, 11}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 88, 50 pixel position
		{Pixel{87, 49}, demotif{Pixel{22, 39}}, true},
		{Pixel{87, 50}, demotif{Pixel{45, 10}}, true},
		{Pixel{87, 51}, demotif{Pixel{45, 11}}, true},
		{Pixel{88, 49}, demotif{Pixel{23, 39}}, true},
		{Pixel{88, 50}, demotif{Pixel{0, 10}}, true},
		{Pixel{88, 51}, demotif{Pixel{0, 11}}, true},
		{Pixel{89, 49}, demotif{Pixel{1, 9}}, true},
		{Pixel{89, 50}, demotif{Pixel{1, 10}}, true},
		{Pixel{89, 51}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 88, 70 pixel position
		{Pixel{87, 69}, demotif{Pixel{45, 29}}, true},
		{Pixel{87, 70}, demotif{Pixel{22, 0}}, true},
		{Pixel{87, 71}, demotif{Pixel{22, 1}}, true},
		{Pixel{88, 69}, demotif{Pixel{0, 29}}, true},
		{Pixel{88, 70}, demotif{Pixel{23, 0}}, true},
		{Pixel{88, 71}, demotif{Pixel{23, 1}}, true},
		{Pixel{89, 69}, demotif{Pixel{1, 29}}, true},
		{Pixel{89, 70}, demotif{Pixel{1, 30}}, true},
		{Pixel{89, 71}, demotif{Pixel{24, 1}}, true},

		// bottom left of entire hex grid
		{Pixel{-1, 719}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 720}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 719}, demotif{Pixel{4, 19}}, true},
		{Pixel{0, 720}, demotif{Pixel{0, 0}}, false},

		// top right of entire hex grid
		{Pixel{1279, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{1280, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{1279, 0}, demotif{Pixel{41, 20}}, true},
		{Pixel{1280, 0}, demotif{Pixel{0, 0}}, false},

		// bottom right of entire hex grid
		{Pixel{1279, 719}, demotif{Pixel{41, 19}}, true},
		{Pixel{1280, 719}, demotif{Pixel{0, 0}}, false},
		{Pixel{1279, 720}, demotif{Pixel{0, 0}}, false},
		{Pixel{1280, 720}, demotif{Pixel{0, 0}}, false},
	}

	g := newGrid(defaultGridConfig())
	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"Pixel %s demotifs to %s and want it to %s",
			tt.pixel.log(),
			tt.wantDemotif.log(),
			sWantOk)

		t.Run(testname, func(t *testing.T) {
			demotif, ok := g.demotif(tt.pixel)
			if ok != tt.wantOk {
				t.Errorf("Want pixel position to %s", sWantOk)
			}
			if !demotif.equals(tt.wantDemotif.Pixel) {
				t.Errorf("For Pixel %s want demotif %s but it was %s",
					tt.pixel.log(),
					tt.wantDemotif.log(),
					demotif.log())
			}
		})
	}
}

func TestDemotifLinePixels(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		hs        hexSide
		c         clusivity
		wantLine  []Pixel
		wantPanic bool
	}{
		{hsNone, cInclusive, []Pixel{}, true},
		{hsWest, cInclusive, []Pixel{
			{0, 10}, {0, 11}, {0, 12}, {0, 13}, {0, 14},
			{0, 15}, {0, 16}, {0, 17}, {0, 18}, {0, 19},
			{0, 20}, {0, 21}, {0, 22}, {0, 23}, {0, 24},
			{0, 25}, {0, 26}, {0, 27}, {0, 28}, {0, 29},
		}, false},
		{hsWest, cExclusive, []Pixel{
			{-1, 10}, {-1, 11}, {-1, 12}, {-1, 13}, {-1, 14},
			{-1, 15}, {-1, 16}, {-1, 17}, {-1, 18}, {-1, 19},
			{-1, 20}, {-1, 21}, {-1, 22}, {-1, 23}, {-1, 24},
			{-1, 25}, {-1, 26}, {-1, 27}, {-1, 28}, {-1, 29},
		}, false},
		{hsNorthwest, cInclusive, []Pixel{
			{0, 10}, {1, 9}, {2, 9}, {3, 8}, {4, 8},
			{5, 8}, {6, 7}, {7, 7}, {8, 6}, {9, 6},
			{10, 5}, {11, 5}, {12, 5}, {13, 4}, {14, 4},
			{15, 3}, {16, 3}, {17, 2}, {18, 2}, {19, 2},
			{20, 1}, {21, 1}, {22, 0},
		}, false},
		{hsNorthwest, cExclusive, []Pixel{
			{0, 9}, {1, 8}, {2, 8}, {3, 7}, {4, 7},
			{5, 7}, {6, 6}, {7, 6}, {8, 5}, {9, 5},
			{10, 4}, {11, 4}, {12, 4}, {13, 3}, {14, 3},
			{15, 2}, {16, 2}, {17, 1}, {18, 1}, {19, 1},
			{20, 0}, {21, 0}, {22, -1},
		}, false},
		{hsNortheast, cInclusive, []Pixel{
			{23, 0}, {24, 1}, {25, 1}, {26, 2}, {27, 2},
			{28, 2}, {29, 3}, {30, 3}, {31, 4}, {32, 4},
			{33, 5}, {34, 5}, {35, 5}, {36, 6}, {37, 6},
			{38, 7}, {39, 7}, {40, 8}, {41, 8}, {42, 8},
			{43, 9}, {44, 9}, {45, 10},
		}, false},
		{hsNortheast, cExclusive, []Pixel{
			{23, -1}, {24, 0}, {25, 0}, {26, 1}, {27, 1},
			{28, 1}, {29, 2}, {30, 2}, {31, 3}, {32, 3},
			{33, 4}, {34, 4}, {35, 4}, {36, 5}, {37, 5},
			{38, 6}, {39, 6}, {40, 7}, {41, 7}, {42, 7},
			{43, 8}, {44, 8}, {45, 9},
		}, false},
		{hsEast, cInclusive, []Pixel{
			{45, 10}, {45, 11}, {45, 12}, {45, 13}, {45, 14},
			{45, 15}, {45, 16}, {45, 17}, {45, 18}, {45, 19},
			{45, 20}, {45, 21}, {45, 22}, {45, 23}, {45, 24},
			{45, 25}, {45, 26}, {45, 27}, {45, 28}, {45, 29},
		}, false},
		{hsEast, cExclusive, []Pixel{
			{46, 10}, {46, 11}, {46, 12}, {46, 13}, {46, 14},
			{46, 15}, {46, 16}, {46, 17}, {46, 18}, {46, 19},
			{46, 20}, {46, 21}, {46, 22}, {46, 23}, {46, 24},
			{46, 25}, {46, 26}, {46, 27}, {46, 28}, {46, 29},
		}, false},
		{hsSouthwest, cInclusive, []Pixel{
			{0, 29}, {1, 30}, {2, 30}, {3, 31}, {4, 31},
			{5, 31}, {6, 32}, {7, 32}, {8, 33}, {9, 33},
			{10, 34}, {11, 34}, {12, 34}, {13, 35}, {14, 35},
			{15, 36}, {16, 36}, {17, 37}, {18, 37}, {19, 37},
			{20, 38}, {21, 38}, {22, 39},
		}, false},
		{hsSouthwest, cExclusive, []Pixel{
			{0, 30}, {1, 31}, {2, 31}, {3, 32}, {4, 32},
			{5, 32}, {6, 33}, {7, 33}, {8, 34}, {9, 34},
			{10, 35}, {11, 35}, {12, 35}, {13, 36}, {14, 36},
			{15, 37}, {16, 37}, {17, 38}, {18, 38}, {19, 38},
			{20, 39}, {21, 39}, {22, 40},
		}, false},
		{hsSoutheast, cInclusive, []Pixel{
			{23, 39}, {24, 38}, {25, 38}, {26, 37}, {27, 37},
			{28, 37}, {29, 36}, {30, 36}, {31, 35}, {32, 35},
			{33, 34}, {34, 34}, {35, 34}, {36, 33}, {37, 33},
			{38, 32}, {39, 32}, {40, 31}, {41, 31}, {42, 31},
			{43, 30}, {44, 30}, {45, 29},
		}, false},
		{hsSoutheast, cExclusive, []Pixel{
			{23, 40}, {24, 39}, {25, 39}, {26, 38}, {27, 38},
			{28, 38}, {29, 37}, {30, 37}, {31, 36}, {32, 36},
			{33, 35}, {34, 35}, {35, 35}, {36, 34}, {37, 34},
			{38, 33}, {39, 33}, {40, 32}, {41, 32}, {42, 32},
			{43, 31}, {44, 31}, {45, 30},
		}, false},
	}

	g := newGrid(defaultGridConfig())
	for _, tt := range tests {
		testname := fmt.Sprintf("test demotif hex line for hex side %s", tt.hs.string())
		t.Run(testname, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil && tt.wantPanic:
					t.Errorf("want panic but did not panic")
				case r != nil && !tt.wantPanic:
					t.Errorf("do not want panic but did panic")
				}
			}()
			line := g.demotifLinePixels(tt.hs, tt.c)
			if !slices.Equal(line, tt.wantLine) {
				t.Errorf("demotifLinePixel has pixels %#v rather than expected pixels \n %#v",
					line, tt.wantLine)
			}
		})
	}
}

func TestDemotifPixels(t *testing.T) {
	g := newGrid(defaultGridConfig())
	demotifPixels := g.demotifHexPixels()
	length := len(demotifPixels)
	if length != 1380 {
		t.Error("demotifPixels sanity check failed because lengths is", length)
	}
}

func TestDemotifHexPixelsAdjacent(t *testing.T) {
	g := newGrid(defaultGridConfig())
	demotifPixelsAdjacent := g.demotifHexPixelsAdjacent()

	length := len(demotifPixelsAdjacent)
	ePixelCount, wPixelCount := g.hexHeight/2, g.hexHeight/2
	nwPixelCount := g.hexWidth / 2
	nePixelCount := g.hexWidth / 2
	swPixelCount := g.hexWidth / 2
	sePixelCount := g.hexWidth / 2
	expectedLength := ePixelCount +
		wPixelCount + nwPixelCount + nePixelCount + swPixelCount + sePixelCount
	if length != expectedLength {
		t.Errorf("demotifPixels sanity check failed because length is %d rather than expected %d",
			length, expectedLength)
		fmt.Println(demotifPixelsAdjacent)
	}

	wPixels := []Pixel{
		{-1, 10}, {-1, 11}, {-1, 12}, {-1, 13}, {-1, 14},
		{-1, 15}, {-1, 16}, {-1, 17}, {-1, 18}, {-1, 19},
		{-1, 20}, {-1, 21}, {-1, 22}, {-1, 23}, {-1, 24},
		{-1, 25}, {-1, 26}, {-1, 27}, {-1, 28}, {-1, 29},
	}
	nwPixels := []Pixel{
		{0, 9}, {1, 8}, {2, 8}, {3, 7}, {4, 7},
		{5, 7}, {6, 6}, {7, 6}, {8, 5}, {9, 5},
		{10, 4}, {11, 4}, {12, 4}, {13, 3}, {14, 3},
		{15, 2}, {16, 2}, {17, 1}, {18, 1}, {19, 1},
		{20, 0}, {21, 0}, {22, -1},
	}
	nePixels := []Pixel{
		{23, -1}, {24, 0}, {25, 0}, {26, 1}, {27, 1},
		{28, 1}, {29, 2}, {30, 2}, {31, 3}, {32, 3},
		{33, 4}, {34, 4}, {35, 4}, {36, 5}, {37, 5},
		{38, 6}, {39, 6}, {40, 7}, {41, 7}, {42, 7},
		{43, 8}, {44, 8}, {45, 9},
	}
	ePixels := []Pixel{
		{46, 10}, {46, 11}, {46, 12}, {46, 13}, {46, 14},
		{46, 15}, {46, 16}, {46, 17}, {46, 18}, {46, 19},
		{46, 20}, {46, 21}, {46, 22}, {46, 23}, {46, 24},
		{46, 25}, {46, 26}, {46, 27}, {46, 28}, {46, 29},
	}
	swPixels := []Pixel{
		{0, 30}, {1, 31}, {2, 31}, {3, 32}, {4, 32},
		{5, 32}, {6, 33}, {7, 33}, {8, 34}, {9, 34},
		{10, 35}, {11, 35}, {12, 35}, {13, 36}, {14, 36},
		{15, 37}, {16, 37}, {17, 38}, {18, 38}, {19, 38},
		{20, 39}, {21, 39}, {22, 40},
	}
	sePixels := []Pixel{
		{23, 40}, {24, 39}, {25, 39}, {26, 38}, {27, 38},
		{28, 38}, {29, 37}, {30, 37}, {31, 36}, {32, 36},
		{33, 35}, {34, 35}, {35, 35}, {36, 34}, {37, 34},
		{38, 33}, {39, 33}, {40, 32}, {41, 32}, {42, 32},
		{43, 31}, {44, 31}, {45, 30},
	}
	format := "expected pixel %s to be %s"
	for _, pixel := range wPixels {
		if demotifPixelsAdjacent[pixel] != hsWest {
			t.Error(format, pixel.log(), hsWest.string())
		}
	}
	for _, pixel := range nwPixels {
		if demotifPixelsAdjacent[pixel] != hsNorthwest {
			t.Error(format, pixel.log(), hsNorthwest.string())
		}
	}
	for _, pixel := range nePixels {
		if demotifPixelsAdjacent[pixel] != hsNortheast {
			t.Error(format, pixel.log(), hsNortheast.string())
		}
	}
	for _, pixel := range ePixels {
		if demotifPixelsAdjacent[pixel] != hsEast {
			t.Error(format, pixel.log(), hsEast.string())
		}
	}
	for _, pixel := range swPixels {
		if demotifPixelsAdjacent[pixel] != hsSouthwest {
			t.Error(format, pixel.log(), hsSouthwest.string())
		}
	}
	for _, pixel := range sePixels {
		if demotifPixelsAdjacent[pixel] != hsSoutheast {
			t.Error(format, pixel.log(), hsSoutheast.string())
		}
	}
}
