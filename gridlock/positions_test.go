package main

import (
	"fmt"
	"testing"
)

func TestGetPixelPos(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(46, 40, 19, 10)

	var tests = []struct {
		gridPos        GridPosition
		hexPoint       hexPoint
		wantedPixelPos PixelPosition
		wantedOk       bool
	}{

		// check hexagon {1, 1}, {1, 2}, {2,1} to ensure hex points are correct
		{GridPosition{1, 1}, hpCenter, PixelPosition{42, 30}, true},
		{GridPosition{1, 1}, hpNorth, PixelPosition{42, 10}, true},
		{GridPosition{1, 1}, hpNortheast, PixelPosition{64, 20}, true},
		{GridPosition{1, 1}, hpSoutheast, PixelPosition{64, 39}, true},
		{GridPosition{1, 1}, hpSouth, PixelPosition{42, 49}, true},
		{GridPosition{1, 1}, hpSouthwest, PixelPosition{19, 39}, true},
		{GridPosition{1, 1}, hpNorthwest, PixelPosition{19, 20}, true},
		{GridPosition{2, 1}, hpCenter, PixelPosition{88, 30}, true},
		{GridPosition{2, 1}, hpNorth, PixelPosition{88, 10}, true},
		{GridPosition{2, 1}, hpNortheast, PixelPosition{110, 20}, true},
		{GridPosition{2, 1}, hpSoutheast, PixelPosition{110, 39}, true},
		{GridPosition{2, 1}, hpSouth, PixelPosition{88, 49}, true},
		{GridPosition{2, 1}, hpSouthwest, PixelPosition{65, 39}, true},
		{GridPosition{2, 1}, hpNorthwest, PixelPosition{65, 20}, true},
		{GridPosition{1, 2}, hpCenter, PixelPosition{65, 60}, true},
		{GridPosition{1, 2}, hpNorth, PixelPosition{65, 40}, true},
		{GridPosition{1, 2}, hpNortheast, PixelPosition{87, 50}, true},
		{GridPosition{1, 2}, hpSoutheast, PixelPosition{87, 69}, true},
		{GridPosition{1, 2}, hpSouth, PixelPosition{65, 79}, true},
		{GridPosition{1, 2}, hpSouthwest, PixelPosition{42, 69}, true},
		{GridPosition{1, 2}, hpNorthwest, PixelPosition{42, 50}, true},

		// check hexagon {26,23} to ensure large x and large y positions are correct
		{GridPosition{27, 23}, hpCenter, PixelPosition{1238, 690}, true},
		{GridPosition{27, 23}, hpNorth, PixelPosition{1238, 670}, true},
		{GridPosition{27, 23}, hpNortheast, PixelPosition{1260, 680}, true},
		{GridPosition{27, 23}, hpSoutheast, PixelPosition{1260, 699}, true},
		{GridPosition{27, 23}, hpSouth, PixelPosition{1238, 709}, true},
		{GridPosition{27, 23}, hpSouthwest, PixelPosition{1215, 699}, true},
		{GridPosition{27, 23}, hpNorthwest, PixelPosition{1215, 680}, true},

		// check top left hexagon {0,0}
		{GridPosition{0, 0}, hpCenter, PixelPosition{19, 0}, true},
		{GridPosition{0, 0}, hpNorth, PixelPosition{0, 0}, false},
		{GridPosition{0, 0}, hpNortheast, PixelPosition{0, 0}, false},
		{GridPosition{0, 0}, hpSoutheast, PixelPosition{41, 9}, true},
		{GridPosition{0, 0}, hpSouth, PixelPosition{19, 19}, true},
		{GridPosition{0, 0}, hpSouthwest, PixelPosition{0, 0}, false},
		{GridPosition{0, 0}, hpNorthwest, PixelPosition{0, 0}, false},

		// check bottom left hexagon {0,24}
		{GridPosition{0, 24}, hpCenter, PixelPosition{0, 0}, false},
		{GridPosition{0, 24}, hpNorth, PixelPosition{19, 700}, true},
		{GridPosition{0, 24}, hpNortheast, PixelPosition{41, 710}, true},
		{GridPosition{0, 24}, hpSoutheast, PixelPosition{0, 0}, false},
		{GridPosition{0, 24}, hpSouth, PixelPosition{0, 0}, false},
		{GridPosition{0, 24}, hpSouthwest, PixelPosition{0, 0}, false},
		{GridPosition{0, 24}, hpNorthwest, PixelPosition{0, 0}, false},

		// check top right hexagon {27,0}
		{GridPosition{27, 0}, hpCenter, PixelPosition{1261, 0}, true},
		{GridPosition{27, 0}, hpNorth, PixelPosition{0, 0}, false},
		{GridPosition{27, 0}, hpNortheast, PixelPosition{0, 0}, false},
		{GridPosition{27, 0}, hpSoutheast, PixelPosition{0, 0}, false},
		{GridPosition{27, 0}, hpSouth, PixelPosition{1261, 19}, true},
		{GridPosition{27, 0}, hpSouthwest, PixelPosition{1238, 9}, true},
		{GridPosition{27, 0}, hpNorthwest, PixelPosition{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{GridPosition{27, 24}, hpCenter, PixelPosition{0, 0}, false},
		{GridPosition{27, 24}, hpNorth, PixelPosition{1261, 700}, true},
		{GridPosition{27, 24}, hpNortheast, PixelPosition{0, 0}, false},
		{GridPosition{27, 24}, hpSoutheast, PixelPosition{0, 0}, false},
		{GridPosition{27, 24}, hpSouth, PixelPosition{0, 0}, false},
		{GridPosition{27, 24}, hpSouthwest, PixelPosition{0, 0}, false},
		{GridPosition{27, 24}, hpNorthwest, PixelPosition{1238, 710}, true},
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
			tt.gridPos.X,
			tt.gridPos.Y,
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
					tt.gridPos.X,
					tt.gridPos.Y,
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
	hc := newHexConfig(46, 40, 19, 10)

	var tests = []struct {
		pixelPos  PixelPosition
		wantedPos GridPosition
		wantedOk  bool
	}{
		// top left of entire hex grid
		{PixelPosition{0, 0}, GridPosition{0, 0}, true},
		{PixelPosition{-1, 0}, GridPosition{0, 0}, false},
		{PixelPosition{0, -1}, GridPosition{0, 0}, false},
		{PixelPosition{-1, -1}, GridPosition{0, 0}, false},

		// first left border point hex change around 0,12 pixel position
		{PixelPosition{-1, 11}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 12}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 13}, GridPosition{0, 0}, false},
		{PixelPosition{0, 11}, GridPosition{0, 0}, true},
		{PixelPosition{0, 12}, GridPosition{0, 1}, true},
		{PixelPosition{0, 13}, GridPosition{0, 1}, true},
		{PixelPosition{1, 11}, GridPosition{0, 0}, true},
		{PixelPosition{1, 12}, GridPosition{0, 1}, true},
		{PixelPosition{1, 13}, GridPosition{0, 1}, true},

		// second left border point hex change around 0,48 pixel position
		{PixelPosition{-1, 47}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 48}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 49}, GridPosition{0, 0}, false},
		{PixelPosition{0, 47}, GridPosition{0, 1}, true},
		{PixelPosition{0, 48}, GridPosition{0, 2}, true},
		{PixelPosition{0, 49}, GridPosition{0, 2}, true},
		{PixelPosition{1, 47}, GridPosition{0, 1}, true},
		{PixelPosition{1, 48}, GridPosition{0, 2}, true},
		{PixelPosition{1, 49}, GridPosition{0, 2}, true},

		// third left border point hex change around 0,72 pixel position
		{PixelPosition{-1, 71}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 72}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 73}, GridPosition{0, 0}, false},
		{PixelPosition{0, 71}, GridPosition{0, 2}, true},
		{PixelPosition{0, 72}, GridPosition{0, 3}, true},
		{PixelPosition{0, 73}, GridPosition{0, 3}, true},
		{PixelPosition{1, 71}, GridPosition{0, 2}, true},
		{PixelPosition{1, 72}, GridPosition{0, 3}, true},
		{PixelPosition{1, 73}, GridPosition{0, 3}, true},

		// first top border point hex change around 42,0 pixel position
		{PixelPosition{41, -1}, GridPosition{0, 0}, false},
		{PixelPosition{41, 0}, GridPosition{0, 0}, true},
		{PixelPosition{41, 1}, GridPosition{0, 0}, true},
		{PixelPosition{42, -1}, GridPosition{0, 0}, false},
		{PixelPosition{42, 0}, GridPosition{1, 0}, true},
		{PixelPosition{42, 1}, GridPosition{1, 0}, true},
		{PixelPosition{43, -1}, GridPosition{0, 0}, false},
		{PixelPosition{43, 0}, GridPosition{1, 0}, true},
		{PixelPosition{43, 1}, GridPosition{1, 0}, true},

		// second top border point hex change around 88, 0 pixel position
		{PixelPosition{87, -1}, GridPosition{0, 0}, false},
		{PixelPosition{87, 0}, GridPosition{1, 0}, true},
		{PixelPosition{87, 1}, GridPosition{1, 0}, true},
		{PixelPosition{88, -1}, GridPosition{0, 0}, false},
		{PixelPosition{88, 0}, GridPosition{2, 0}, true},
		{PixelPosition{88, 1}, GridPosition{2, 0}, true},
		{PixelPosition{89, -1}, GridPosition{0, 0}, false},
		{PixelPosition{89, 0}, GridPosition{2, 0}, true},
		{PixelPosition{89, 1}, GridPosition{2, 0}, true},

		// second top border point hex change around 134,0 pixel position
		{PixelPosition{133, -1}, GridPosition{0, 0}, false},
		{PixelPosition{133, 0}, GridPosition{2, 0}, true},
		{PixelPosition{133, 1}, GridPosition{2, 0}, true},
		{PixelPosition{134, -1}, GridPosition{0, 0}, false},
		{PixelPosition{134, 0}, GridPosition{3, 0}, true},
		{PixelPosition{134, 1}, GridPosition{3, 0}, true},
		{PixelPosition{135, -1}, GridPosition{0, 0}, false},
		{PixelPosition{135, 0}, GridPosition{3, 0}, true},
		{PixelPosition{135, 1}, GridPosition{3, 0}, true},

		// hexagon points near 19, 80 pixel position
		{PixelPosition{18, 79}, GridPosition{0, 2}, true},
		{PixelPosition{18, 80}, GridPosition{0, 3}, true},
		{PixelPosition{18, 81}, GridPosition{0, 3}, true},
		{PixelPosition{19, 79}, GridPosition{0, 2}, true},
		{PixelPosition{19, 80}, GridPosition{1, 3}, true},
		{PixelPosition{19, 81}, GridPosition{1, 3}, true},
		{PixelPosition{20, 79}, GridPosition{1, 3}, true},
		{PixelPosition{20, 80}, GridPosition{1, 3}, true},
		{PixelPosition{20, 81}, GridPosition{1, 3}, true},

		// hexagon points near 42, 70 pixel position
		{PixelPosition{41, 69}, GridPosition{0, 2}, true},
		{PixelPosition{41, 70}, GridPosition{1, 3}, true},
		{PixelPosition{41, 71}, GridPosition{1, 3}, true},
		{PixelPosition{42, 69}, GridPosition{1, 2}, true},
		{PixelPosition{42, 70}, GridPosition{1, 3}, true},
		{PixelPosition{42, 71}, GridPosition{1, 3}, true},
		{PixelPosition{43, 69}, GridPosition{1, 2}, true},
		{PixelPosition{43, 70}, GridPosition{1, 2}, true},
		{PixelPosition{43, 71}, GridPosition{1, 3}, true},

		// hexagon points near 65, 80 pixel position
		{PixelPosition{64, 79}, GridPosition{1, 2}, true},
		{PixelPosition{64, 80}, GridPosition{1, 3}, true},
		{PixelPosition{64, 81}, GridPosition{1, 3}, true},
		{PixelPosition{65, 79}, GridPosition{1, 2}, true},
		{PixelPosition{65, 80}, GridPosition{2, 3}, true},
		{PixelPosition{65, 81}, GridPosition{2, 3}, true},
		{PixelPosition{66, 79}, GridPosition{2, 3}, true},
		{PixelPosition{66, 80}, GridPosition{2, 3}, true},
		{PixelPosition{66, 81}, GridPosition{2, 3}, true},

		// hexagon points near 88, 10 pixel position
		{PixelPosition{87, 9}, GridPosition{1, 0}, true},
		{PixelPosition{87, 10}, GridPosition{2, 1}, true},
		{PixelPosition{87, 11}, GridPosition{2, 1}, true},
		{PixelPosition{88, 9}, GridPosition{2, 0}, true},
		{PixelPosition{88, 10}, GridPosition{2, 1}, true},
		{PixelPosition{88, 11}, GridPosition{2, 1}, true},
		{PixelPosition{89, 9}, GridPosition{2, 0}, true},
		{PixelPosition{89, 10}, GridPosition{2, 0}, true},
		{PixelPosition{89, 11}, GridPosition{2, 1}, true},

		// hexagon points near 88, 50 pixel position
		{PixelPosition{87, 49}, GridPosition{2, 1}, true},
		{PixelPosition{87, 50}, GridPosition{1, 2}, true},
		{PixelPosition{87, 51}, GridPosition{1, 2}, true},
		{PixelPosition{88, 49}, GridPosition{2, 1}, true},
		{PixelPosition{88, 50}, GridPosition{2, 2}, true},
		{PixelPosition{88, 51}, GridPosition{2, 2}, true},
		{PixelPosition{89, 49}, GridPosition{2, 2}, true},
		{PixelPosition{89, 50}, GridPosition{2, 2}, true},
		{PixelPosition{89, 51}, GridPosition{2, 2}, true},

		// hexagon points near 88, 70 pixel position
		{PixelPosition{87, 69}, GridPosition{1, 2}, true},
		{PixelPosition{87, 70}, GridPosition{2, 3}, true},
		{PixelPosition{87, 71}, GridPosition{2, 3}, true},
		{PixelPosition{88, 69}, GridPosition{2, 2}, true},
		{PixelPosition{88, 70}, GridPosition{2, 3}, true},
		{PixelPosition{88, 71}, GridPosition{2, 3}, true},
		{PixelPosition{89, 69}, GridPosition{2, 2}, true},
		{PixelPosition{89, 70}, GridPosition{2, 2}, true},
		{PixelPosition{89, 71}, GridPosition{2, 3}, true},

		// bottom left of entire hex grid
		{PixelPosition{-1, 719}, GridPosition{0, 0}, false},
		{PixelPosition{-1, 720}, GridPosition{0, 0}, false},
		{PixelPosition{0, 719}, GridPosition{0, 24}, true},
		{PixelPosition{0, 720}, GridPosition{0, 0}, false},

		// top right of entire hex grid
		{PixelPosition{1279, -1}, GridPosition{0, 0}, false},
		{PixelPosition{1280, -1}, GridPosition{0, 0}, false},
		{PixelPosition{1279, 0}, GridPosition{27, 0}, true},
		{PixelPosition{1280, 0}, GridPosition{0, 0}, false},

		// bottom right of entire hex grid
		{PixelPosition{1279, 719}, GridPosition{27, 24}, true},
		{PixelPosition{1280, 719}, GridPosition{0, 0}, false},
		{PixelPosition{1279, 720}, GridPosition{0, 0}, false},
		{PixelPosition{1280, 720}, GridPosition{0, 0}, false},

		// first hexagon y section (hexagon look up table)
		// between hexagons 0,0 1,1 and 1,0

		// sanity checks
		{PixelPosition{19, 10}, GridPosition{0, 0}, true},
		{PixelPosition{65, 10}, GridPosition{1, 0}, true},
		{PixelPosition{41, 20}, GridPosition{1, 1}, true},
		{PixelPosition{42, 20}, GridPosition{1, 1}, true},
		{PixelPosition{43, 20}, GridPosition{1, 1}, true},

		// 0,0 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{19, 19}, GridPosition{0, 0}, true},
		{PixelPosition{20, 18}, GridPosition{0, 0}, true},
		{PixelPosition{21, 18}, GridPosition{0, 0}, true},
		{PixelPosition{22, 17}, GridPosition{0, 0}, true},
		{PixelPosition{23, 17}, GridPosition{0, 0}, true},
		{PixelPosition{24, 17}, GridPosition{0, 0}, true},
		{PixelPosition{25, 16}, GridPosition{0, 0}, true},
		{PixelPosition{26, 16}, GridPosition{0, 0}, true},
		{PixelPosition{27, 15}, GridPosition{0, 0}, true},
		{PixelPosition{28, 15}, GridPosition{0, 0}, true},
		{PixelPosition{29, 14}, GridPosition{0, 0}, true},
		{PixelPosition{30, 14}, GridPosition{0, 0}, true},
		{PixelPosition{31, 14}, GridPosition{0, 0}, true},
		{PixelPosition{32, 13}, GridPosition{0, 0}, true},
		{PixelPosition{33, 13}, GridPosition{0, 0}, true},
		{PixelPosition{34, 12}, GridPosition{0, 0}, true},
		{PixelPosition{35, 12}, GridPosition{0, 0}, true},
		{PixelPosition{36, 11}, GridPosition{0, 0}, true},
		{PixelPosition{37, 11}, GridPosition{0, 0}, true},
		{PixelPosition{38, 11}, GridPosition{0, 0}, true},
		{PixelPosition{39, 10}, GridPosition{0, 0}, true},
		{PixelPosition{40, 10}, GridPosition{0, 0}, true},
		{PixelPosition{41, 9}, GridPosition{0, 0}, true},

		// 1, 1 hexagon slope check between itself and 0, 0 hexagon
		{PixelPosition{19, 20}, GridPosition{1, 1}, true},
		{PixelPosition{20, 19}, GridPosition{1, 1}, true},
		{PixelPosition{21, 19}, GridPosition{1, 1}, true},
		{PixelPosition{22, 18}, GridPosition{1, 1}, true},
		{PixelPosition{23, 18}, GridPosition{1, 1}, true},
		{PixelPosition{24, 18}, GridPosition{1, 1}, true},
		{PixelPosition{25, 17}, GridPosition{1, 1}, true},
		{PixelPosition{26, 17}, GridPosition{1, 1}, true},
		{PixelPosition{27, 16}, GridPosition{1, 1}, true},
		{PixelPosition{28, 16}, GridPosition{1, 1}, true},
		{PixelPosition{29, 15}, GridPosition{1, 1}, true},
		{PixelPosition{30, 15}, GridPosition{1, 1}, true},
		{PixelPosition{31, 15}, GridPosition{1, 1}, true},
		{PixelPosition{32, 14}, GridPosition{1, 1}, true},
		{PixelPosition{33, 14}, GridPosition{1, 1}, true},
		{PixelPosition{34, 13}, GridPosition{1, 1}, true},
		{PixelPosition{35, 13}, GridPosition{1, 1}, true},
		{PixelPosition{36, 12}, GridPosition{1, 1}, true},
		{PixelPosition{37, 12}, GridPosition{1, 1}, true},
		{PixelPosition{38, 12}, GridPosition{1, 1}, true},
		{PixelPosition{39, 11}, GridPosition{1, 1}, true},
		{PixelPosition{40, 11}, GridPosition{1, 1}, true},
		{PixelPosition{41, 10}, GridPosition{1, 1}, true},

		// 1,0 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{42, 9}, GridPosition{1, 0}, true},
		{PixelPosition{43, 10}, GridPosition{1, 0}, true},
		{PixelPosition{44, 10}, GridPosition{1, 0}, true},
		{PixelPosition{45, 11}, GridPosition{1, 0}, true},
		{PixelPosition{46, 11}, GridPosition{1, 0}, true},
		{PixelPosition{47, 11}, GridPosition{1, 0}, true},
		{PixelPosition{48, 12}, GridPosition{1, 0}, true},
		{PixelPosition{49, 12}, GridPosition{1, 0}, true},
		{PixelPosition{50, 13}, GridPosition{1, 0}, true},
		{PixelPosition{51, 13}, GridPosition{1, 0}, true},

		{PixelPosition{52, 14}, GridPosition{1, 0}, true},
		{PixelPosition{53, 14}, GridPosition{1, 0}, true},
		{PixelPosition{54, 14}, GridPosition{1, 0}, true},
		{PixelPosition{55, 15}, GridPosition{1, 0}, true},
		{PixelPosition{56, 15}, GridPosition{1, 0}, true},
		{PixelPosition{57, 16}, GridPosition{1, 0}, true},
		{PixelPosition{58, 16}, GridPosition{1, 0}, true},
		{PixelPosition{59, 17}, GridPosition{1, 0}, true},
		{PixelPosition{60, 17}, GridPosition{1, 0}, true},
		{PixelPosition{61, 17}, GridPosition{1, 0}, true},
		{PixelPosition{62, 18}, GridPosition{1, 0}, true},
		{PixelPosition{63, 18}, GridPosition{1, 0}, true},
		{PixelPosition{64, 19}, GridPosition{1, 0}, true},

		// 1,1 hexagon slope check between itself and 1, 0 hexagon
		{PixelPosition{42, 10}, GridPosition{1, 1}, true},
		{PixelPosition{43, 11}, GridPosition{1, 1}, true},
		{PixelPosition{44, 11}, GridPosition{1, 1}, true},
		{PixelPosition{45, 12}, GridPosition{1, 1}, true},
		{PixelPosition{46, 12}, GridPosition{1, 1}, true},
		{PixelPosition{47, 12}, GridPosition{1, 1}, true},
		{PixelPosition{48, 13}, GridPosition{1, 1}, true},
		{PixelPosition{49, 13}, GridPosition{1, 1}, true},
		{PixelPosition{50, 14}, GridPosition{1, 1}, true},
		{PixelPosition{51, 14}, GridPosition{1, 1}, true},
		{PixelPosition{52, 15}, GridPosition{1, 1}, true},
		{PixelPosition{53, 15}, GridPosition{1, 1}, true},
		{PixelPosition{54, 15}, GridPosition{1, 1}, true},
		{PixelPosition{55, 16}, GridPosition{1, 1}, true},
		{PixelPosition{56, 16}, GridPosition{1, 1}, true},
		{PixelPosition{57, 17}, GridPosition{1, 1}, true},
		{PixelPosition{58, 17}, GridPosition{1, 1}, true},
		{PixelPosition{59, 18}, GridPosition{1, 1}, true},
		{PixelPosition{60, 18}, GridPosition{1, 1}, true},
		{PixelPosition{61, 18}, GridPosition{1, 1}, true},
		{PixelPosition{62, 19}, GridPosition{1, 1}, true},
		{PixelPosition{63, 19}, GridPosition{1, 1}, true},
		{PixelPosition{64, 20}, GridPosition{1, 1}, true},

		// second hexagon y section (full rectangle of hexagon 1,1)

		// right side of hexagon 0,1
		{PixelPosition{18, 20}, GridPosition{0, 1}, true},
		{PixelPosition{18, 21}, GridPosition{0, 1}, true},
		{PixelPosition{18, 22}, GridPosition{0, 1}, true},
		{PixelPosition{18, 23}, GridPosition{0, 1}, true},
		{PixelPosition{18, 24}, GridPosition{0, 1}, true},
		{PixelPosition{18, 25}, GridPosition{0, 1}, true},
		{PixelPosition{18, 26}, GridPosition{0, 1}, true},
		{PixelPosition{18, 27}, GridPosition{0, 1}, true},
		{PixelPosition{18, 28}, GridPosition{0, 1}, true},
		{PixelPosition{18, 29}, GridPosition{0, 1}, true},
		{PixelPosition{18, 30}, GridPosition{0, 1}, true},
		{PixelPosition{18, 31}, GridPosition{0, 1}, true},
		{PixelPosition{18, 32}, GridPosition{0, 1}, true},
		{PixelPosition{18, 33}, GridPosition{0, 1}, true},
		{PixelPosition{18, 34}, GridPosition{0, 1}, true},
		{PixelPosition{18, 35}, GridPosition{0, 1}, true},
		{PixelPosition{18, 36}, GridPosition{0, 1}, true},
		{PixelPosition{18, 37}, GridPosition{0, 1}, true},
		{PixelPosition{18, 38}, GridPosition{0, 1}, true},
		{PixelPosition{18, 39}, GridPosition{0, 1}, true},

		// left side of hexagon 1,1
		// hexagons are right side inclusive and left side exclusive
		{PixelPosition{19, 20}, GridPosition{1, 1}, true},
		{PixelPosition{19, 21}, GridPosition{1, 1}, true},
		{PixelPosition{19, 22}, GridPosition{1, 1}, true},
		{PixelPosition{19, 23}, GridPosition{1, 1}, true},
		{PixelPosition{19, 24}, GridPosition{1, 1}, true},
		{PixelPosition{19, 25}, GridPosition{1, 1}, true},
		{PixelPosition{19, 26}, GridPosition{1, 1}, true},
		{PixelPosition{19, 27}, GridPosition{1, 1}, true},
		{PixelPosition{19, 28}, GridPosition{1, 1}, true},
		{PixelPosition{19, 29}, GridPosition{1, 1}, true},
		{PixelPosition{19, 30}, GridPosition{1, 1}, true},
		{PixelPosition{19, 31}, GridPosition{1, 1}, true},
		{PixelPosition{19, 32}, GridPosition{1, 1}, true},
		{PixelPosition{19, 33}, GridPosition{1, 1}, true},
		{PixelPosition{19, 34}, GridPosition{1, 1}, true},
		{PixelPosition{19, 35}, GridPosition{1, 1}, true},
		{PixelPosition{19, 36}, GridPosition{1, 1}, true},
		{PixelPosition{19, 37}, GridPosition{1, 1}, true},
		{PixelPosition{19, 38}, GridPosition{1, 1}, true},
		{PixelPosition{19, 39}, GridPosition{1, 1}, true},
		// continue left side of hexagon 1,1
		{PixelPosition{20, 20}, GridPosition{1, 1}, true},
		{PixelPosition{20, 21}, GridPosition{1, 1}, true},
		{PixelPosition{20, 22}, GridPosition{1, 1}, true},
		{PixelPosition{20, 23}, GridPosition{1, 1}, true},
		{PixelPosition{20, 24}, GridPosition{1, 1}, true},
		{PixelPosition{20, 25}, GridPosition{1, 1}, true},
		{PixelPosition{20, 26}, GridPosition{1, 1}, true},
		{PixelPosition{20, 27}, GridPosition{1, 1}, true},
		{PixelPosition{20, 28}, GridPosition{1, 1}, true},
		{PixelPosition{20, 29}, GridPosition{1, 1}, true},
		{PixelPosition{20, 30}, GridPosition{1, 1}, true},
		{PixelPosition{20, 31}, GridPosition{1, 1}, true},
		{PixelPosition{20, 32}, GridPosition{1, 1}, true},
		{PixelPosition{20, 33}, GridPosition{1, 1}, true},
		{PixelPosition{20, 34}, GridPosition{1, 1}, true},
		{PixelPosition{20, 35}, GridPosition{1, 1}, true},
		{PixelPosition{20, 36}, GridPosition{1, 1}, true},
		{PixelPosition{20, 37}, GridPosition{1, 1}, true},
		{PixelPosition{20, 38}, GridPosition{1, 1}, true},
		{PixelPosition{20, 39}, GridPosition{1, 1}, true},

		// right side of hexagon 1,1
		{PixelPosition{64, 20}, GridPosition{1, 1}, true},
		{PixelPosition{64, 21}, GridPosition{1, 1}, true},
		{PixelPosition{64, 22}, GridPosition{1, 1}, true},
		{PixelPosition{64, 23}, GridPosition{1, 1}, true},
		{PixelPosition{64, 24}, GridPosition{1, 1}, true},
		{PixelPosition{64, 25}, GridPosition{1, 1}, true},
		{PixelPosition{64, 26}, GridPosition{1, 1}, true},
		{PixelPosition{64, 27}, GridPosition{1, 1}, true},
		{PixelPosition{64, 28}, GridPosition{1, 1}, true},
		{PixelPosition{64, 29}, GridPosition{1, 1}, true},
		{PixelPosition{64, 30}, GridPosition{1, 1}, true},
		{PixelPosition{64, 31}, GridPosition{1, 1}, true},
		{PixelPosition{64, 32}, GridPosition{1, 1}, true},
		{PixelPosition{64, 33}, GridPosition{1, 1}, true},
		{PixelPosition{64, 34}, GridPosition{1, 1}, true},
		{PixelPosition{64, 35}, GridPosition{1, 1}, true},
		{PixelPosition{64, 36}, GridPosition{1, 1}, true},
		{PixelPosition{64, 37}, GridPosition{1, 1}, true},
		{PixelPosition{64, 38}, GridPosition{1, 1}, true},
		{PixelPosition{64, 39}, GridPosition{1, 1}, true},

		// left side of hexagon 2,1
		// hexagons are right side inclusive and left side exclusive
		{PixelPosition{65, 20}, GridPosition{2, 1}, true},
		{PixelPosition{65, 21}, GridPosition{2, 1}, true},
		{PixelPosition{65, 22}, GridPosition{2, 1}, true},
		{PixelPosition{65, 23}, GridPosition{2, 1}, true},
		{PixelPosition{65, 24}, GridPosition{2, 1}, true},
		{PixelPosition{65, 25}, GridPosition{2, 1}, true},
		{PixelPosition{65, 26}, GridPosition{2, 1}, true},
		{PixelPosition{65, 27}, GridPosition{2, 1}, true},
		{PixelPosition{65, 28}, GridPosition{2, 1}, true},
		{PixelPosition{65, 29}, GridPosition{2, 1}, true},
		{PixelPosition{65, 30}, GridPosition{2, 1}, true},
		{PixelPosition{65, 31}, GridPosition{2, 1}, true},
		{PixelPosition{65, 32}, GridPosition{2, 1}, true},
		{PixelPosition{65, 33}, GridPosition{2, 1}, true},
		{PixelPosition{65, 34}, GridPosition{2, 1}, true},
		{PixelPosition{65, 35}, GridPosition{2, 1}, true},
		{PixelPosition{65, 36}, GridPosition{2, 1}, true},
		{PixelPosition{65, 37}, GridPosition{2, 1}, true},
		{PixelPosition{65, 38}, GridPosition{2, 1}, true},
		{PixelPosition{65, 39}, GridPosition{2, 1}, true},
		// continue left side of hexagon 2,1
		{PixelPosition{66, 20}, GridPosition{2, 1}, true},
		{PixelPosition{66, 21}, GridPosition{2, 1}, true},
		{PixelPosition{66, 22}, GridPosition{2, 1}, true},
		{PixelPosition{66, 23}, GridPosition{2, 1}, true},
		{PixelPosition{66, 24}, GridPosition{2, 1}, true},
		{PixelPosition{66, 25}, GridPosition{2, 1}, true},
		{PixelPosition{66, 26}, GridPosition{2, 1}, true},
		{PixelPosition{66, 27}, GridPosition{2, 1}, true},
		{PixelPosition{66, 28}, GridPosition{2, 1}, true},
		{PixelPosition{66, 29}, GridPosition{2, 1}, true},
		{PixelPosition{66, 30}, GridPosition{2, 1}, true},
		{PixelPosition{66, 31}, GridPosition{2, 1}, true},
		{PixelPosition{66, 32}, GridPosition{2, 1}, true},
		{PixelPosition{66, 33}, GridPosition{2, 1}, true},
		{PixelPosition{66, 34}, GridPosition{2, 1}, true},
		{PixelPosition{66, 35}, GridPosition{2, 1}, true},
		{PixelPosition{66, 36}, GridPosition{2, 1}, true},
		{PixelPosition{66, 37}, GridPosition{2, 1}, true},
		{PixelPosition{66, 38}, GridPosition{2, 1}, true},
		{PixelPosition{66, 39}, GridPosition{2, 1}, true},

		// third hexagon y section (inverted hexagon look up table)
		// between hexagons 0,2 1,1 and 1,2

		// sanity checks
		{PixelPosition{19, 49}, GridPosition{0, 2}, true},
		{PixelPosition{65, 49}, GridPosition{1, 2}, true},
		{PixelPosition{41, 30}, GridPosition{1, 1}, true},
		{PixelPosition{42, 30}, GridPosition{1, 1}, true},
		{PixelPosition{43, 30}, GridPosition{1, 1}, true},

		// 0,2 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{19, 40}, GridPosition{0, 2}, true},
		{PixelPosition{20, 41}, GridPosition{0, 2}, true},
		{PixelPosition{21, 41}, GridPosition{0, 2}, true},
		{PixelPosition{22, 42}, GridPosition{0, 2}, true},
		{PixelPosition{23, 42}, GridPosition{0, 2}, true},
		{PixelPosition{24, 42}, GridPosition{0, 2}, true},
		{PixelPosition{25, 43}, GridPosition{0, 2}, true},
		{PixelPosition{26, 43}, GridPosition{0, 2}, true},
		{PixelPosition{27, 44}, GridPosition{0, 2}, true},
		{PixelPosition{28, 44}, GridPosition{0, 2}, true},
		{PixelPosition{29, 45}, GridPosition{0, 2}, true},
		{PixelPosition{30, 45}, GridPosition{0, 2}, true},
		{PixelPosition{31, 45}, GridPosition{0, 2}, true},
		{PixelPosition{32, 46}, GridPosition{0, 2}, true},
		{PixelPosition{33, 46}, GridPosition{0, 2}, true},
		{PixelPosition{34, 47}, GridPosition{0, 2}, true},
		{PixelPosition{35, 47}, GridPosition{0, 2}, true},
		{PixelPosition{36, 48}, GridPosition{0, 2}, true},
		{PixelPosition{37, 48}, GridPosition{0, 2}, true},
		{PixelPosition{38, 48}, GridPosition{0, 2}, true},
		{PixelPosition{39, 49}, GridPosition{0, 2}, true},
		{PixelPosition{40, 49}, GridPosition{0, 2}, true},
		{PixelPosition{41, 50}, GridPosition{0, 2}, true},

		// 1, 1 hexagon slope check between itself and 0, 2 hexagon
		{PixelPosition{19, 39}, GridPosition{1, 1}, true},
		{PixelPosition{20, 40}, GridPosition{1, 1}, true},
		{PixelPosition{21, 40}, GridPosition{1, 1}, true},
		{PixelPosition{22, 41}, GridPosition{1, 1}, true},
		{PixelPosition{23, 41}, GridPosition{1, 1}, true},
		{PixelPosition{24, 41}, GridPosition{1, 1}, true},
		{PixelPosition{25, 42}, GridPosition{1, 1}, true},
		{PixelPosition{26, 42}, GridPosition{1, 1}, true},
		{PixelPosition{27, 43}, GridPosition{1, 1}, true},
		{PixelPosition{28, 43}, GridPosition{1, 1}, true},
		{PixelPosition{29, 44}, GridPosition{1, 1}, true},
		{PixelPosition{30, 44}, GridPosition{1, 1}, true},
		{PixelPosition{31, 44}, GridPosition{1, 1}, true},
		{PixelPosition{32, 45}, GridPosition{1, 1}, true},
		{PixelPosition{33, 45}, GridPosition{1, 1}, true},
		{PixelPosition{34, 46}, GridPosition{1, 1}, true},
		{PixelPosition{35, 46}, GridPosition{1, 1}, true},
		{PixelPosition{36, 47}, GridPosition{1, 1}, true},
		{PixelPosition{37, 47}, GridPosition{1, 1}, true},
		{PixelPosition{38, 47}, GridPosition{1, 1}, true},
		{PixelPosition{39, 48}, GridPosition{1, 1}, true},
		{PixelPosition{40, 48}, GridPosition{1, 1}, true},
		{PixelPosition{41, 49}, GridPosition{1, 1}, true},

		// 1,1 hexagon slope check between itself and 1, 2 hexagon
		{PixelPosition{42, 49}, GridPosition{1, 1}, true},
		{PixelPosition{43, 48}, GridPosition{1, 1}, true},
		{PixelPosition{44, 48}, GridPosition{1, 1}, true},
		{PixelPosition{45, 47}, GridPosition{1, 1}, true},
		{PixelPosition{46, 47}, GridPosition{1, 1}, true},
		{PixelPosition{47, 47}, GridPosition{1, 1}, true},
		{PixelPosition{48, 46}, GridPosition{1, 1}, true},
		{PixelPosition{49, 46}, GridPosition{1, 1}, true},
		{PixelPosition{50, 45}, GridPosition{1, 1}, true},
		{PixelPosition{51, 45}, GridPosition{1, 1}, true},
		{PixelPosition{52, 44}, GridPosition{1, 1}, true},
		{PixelPosition{53, 44}, GridPosition{1, 1}, true},
		{PixelPosition{54, 44}, GridPosition{1, 1}, true},
		{PixelPosition{55, 43}, GridPosition{1, 1}, true},
		{PixelPosition{56, 43}, GridPosition{1, 1}, true},
		{PixelPosition{57, 42}, GridPosition{1, 1}, true},
		{PixelPosition{58, 42}, GridPosition{1, 1}, true},
		{PixelPosition{59, 41}, GridPosition{1, 1}, true},
		{PixelPosition{60, 41}, GridPosition{1, 1}, true},
		{PixelPosition{61, 41}, GridPosition{1, 1}, true},
		{PixelPosition{62, 40}, GridPosition{1, 1}, true},
		{PixelPosition{63, 40}, GridPosition{1, 1}, true},
		{PixelPosition{64, 39}, GridPosition{1, 1}, true},

		// 1,2 hexagon slope check between itself and 1, 1 hexagon
		{PixelPosition{42, 50}, GridPosition{1, 2}, true},
		{PixelPosition{43, 49}, GridPosition{1, 2}, true},
		{PixelPosition{44, 49}, GridPosition{1, 2}, true},
		{PixelPosition{45, 48}, GridPosition{1, 2}, true},
		{PixelPosition{46, 48}, GridPosition{1, 2}, true},
		{PixelPosition{47, 48}, GridPosition{1, 2}, true},
		{PixelPosition{48, 47}, GridPosition{1, 2}, true},
		{PixelPosition{49, 47}, GridPosition{1, 2}, true},
		{PixelPosition{50, 46}, GridPosition{1, 2}, true},
		{PixelPosition{51, 46}, GridPosition{1, 2}, true},
		{PixelPosition{52, 45}, GridPosition{1, 2}, true},
		{PixelPosition{53, 45}, GridPosition{1, 2}, true},
		{PixelPosition{54, 45}, GridPosition{1, 2}, true},
		{PixelPosition{55, 44}, GridPosition{1, 2}, true},
		{PixelPosition{56, 44}, GridPosition{1, 2}, true},
		{PixelPosition{57, 43}, GridPosition{1, 2}, true},
		{PixelPosition{58, 43}, GridPosition{1, 2}, true},
		{PixelPosition{59, 42}, GridPosition{1, 2}, true},
		{PixelPosition{60, 42}, GridPosition{1, 2}, true},
		{PixelPosition{61, 42}, GridPosition{1, 2}, true},
		{PixelPosition{62, 41}, GridPosition{1, 2}, true},
		{PixelPosition{63, 41}, GridPosition{1, 2}, true},
		{PixelPosition{64, 40}, GridPosition{1, 2}, true},

		// fourth hexagon y section (split rectangle between hexagon 0,2 and 1,2)

		// middle of hexagon 0,2
		{PixelPosition{18, 50}, GridPosition{0, 2}, true},
		{PixelPosition{18, 51}, GridPosition{0, 2}, true},
		{PixelPosition{18, 52}, GridPosition{0, 2}, true},
		{PixelPosition{18, 53}, GridPosition{0, 2}, true},
		{PixelPosition{18, 54}, GridPosition{0, 2}, true},
		{PixelPosition{18, 55}, GridPosition{0, 2}, true},
		{PixelPosition{18, 56}, GridPosition{0, 2}, true},
		{PixelPosition{18, 57}, GridPosition{0, 2}, true},
		{PixelPosition{18, 58}, GridPosition{0, 2}, true},
		{PixelPosition{18, 59}, GridPosition{0, 2}, true},
		{PixelPosition{18, 60}, GridPosition{0, 2}, true},
		{PixelPosition{18, 61}, GridPosition{0, 2}, true},
		{PixelPosition{18, 62}, GridPosition{0, 2}, true},
		{PixelPosition{18, 63}, GridPosition{0, 2}, true},
		{PixelPosition{18, 64}, GridPosition{0, 2}, true},
		{PixelPosition{18, 65}, GridPosition{0, 2}, true},
		{PixelPosition{18, 66}, GridPosition{0, 2}, true},
		{PixelPosition{18, 67}, GridPosition{0, 2}, true},
		{PixelPosition{18, 68}, GridPosition{0, 2}, true},
		{PixelPosition{18, 69}, GridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{PixelPosition{19, 50}, GridPosition{0, 2}, true},
		{PixelPosition{19, 51}, GridPosition{0, 2}, true},
		{PixelPosition{19, 52}, GridPosition{0, 2}, true},
		{PixelPosition{19, 53}, GridPosition{0, 2}, true},
		{PixelPosition{19, 54}, GridPosition{0, 2}, true},
		{PixelPosition{19, 55}, GridPosition{0, 2}, true},
		{PixelPosition{19, 56}, GridPosition{0, 2}, true},
		{PixelPosition{19, 57}, GridPosition{0, 2}, true},
		{PixelPosition{19, 58}, GridPosition{0, 2}, true},
		{PixelPosition{19, 59}, GridPosition{0, 2}, true},
		{PixelPosition{19, 60}, GridPosition{0, 2}, true},
		{PixelPosition{19, 61}, GridPosition{0, 2}, true},
		{PixelPosition{19, 62}, GridPosition{0, 2}, true},
		{PixelPosition{19, 63}, GridPosition{0, 2}, true},
		{PixelPosition{19, 64}, GridPosition{0, 2}, true},
		{PixelPosition{19, 65}, GridPosition{0, 2}, true},
		{PixelPosition{19, 66}, GridPosition{0, 2}, true},
		{PixelPosition{19, 67}, GridPosition{0, 2}, true},
		{PixelPosition{19, 68}, GridPosition{0, 2}, true},
		{PixelPosition{19, 69}, GridPosition{0, 2}, true},

		// continue middle of hexagon 0,2
		{PixelPosition{20, 50}, GridPosition{0, 2}, true},
		{PixelPosition{20, 51}, GridPosition{0, 2}, true},
		{PixelPosition{20, 52}, GridPosition{0, 2}, true},
		{PixelPosition{20, 53}, GridPosition{0, 2}, true},
		{PixelPosition{20, 54}, GridPosition{0, 2}, true},
		{PixelPosition{20, 55}, GridPosition{0, 2}, true},
		{PixelPosition{20, 56}, GridPosition{0, 2}, true},
		{PixelPosition{20, 57}, GridPosition{0, 2}, true},
		{PixelPosition{20, 58}, GridPosition{0, 2}, true},
		{PixelPosition{20, 59}, GridPosition{0, 2}, true},
		{PixelPosition{20, 60}, GridPosition{0, 2}, true},
		{PixelPosition{20, 61}, GridPosition{0, 2}, true},
		{PixelPosition{20, 62}, GridPosition{0, 2}, true},
		{PixelPosition{20, 63}, GridPosition{0, 2}, true},
		{PixelPosition{20, 64}, GridPosition{0, 2}, true},
		{PixelPosition{20, 65}, GridPosition{0, 2}, true},
		{PixelPosition{20, 66}, GridPosition{0, 2}, true},
		{PixelPosition{20, 67}, GridPosition{0, 2}, true},
		{PixelPosition{20, 68}, GridPosition{0, 2}, true},
		{PixelPosition{20, 69}, GridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		{PixelPosition{41, 50}, GridPosition{0, 2}, true},
		{PixelPosition{41, 51}, GridPosition{0, 2}, true},
		{PixelPosition{41, 52}, GridPosition{0, 2}, true},
		{PixelPosition{41, 53}, GridPosition{0, 2}, true},
		{PixelPosition{41, 54}, GridPosition{0, 2}, true},
		{PixelPosition{41, 55}, GridPosition{0, 2}, true},
		{PixelPosition{41, 56}, GridPosition{0, 2}, true},
		{PixelPosition{41, 57}, GridPosition{0, 2}, true},
		{PixelPosition{41, 58}, GridPosition{0, 2}, true},
		{PixelPosition{41, 59}, GridPosition{0, 2}, true},
		{PixelPosition{41, 60}, GridPosition{0, 2}, true},
		{PixelPosition{41, 61}, GridPosition{0, 2}, true},
		{PixelPosition{41, 62}, GridPosition{0, 2}, true},
		{PixelPosition{41, 63}, GridPosition{0, 2}, true},
		{PixelPosition{41, 64}, GridPosition{0, 2}, true},
		{PixelPosition{41, 65}, GridPosition{0, 2}, true},
		{PixelPosition{41, 66}, GridPosition{0, 2}, true},
		{PixelPosition{41, 67}, GridPosition{0, 2}, true},
		{PixelPosition{41, 68}, GridPosition{0, 2}, true},
		{PixelPosition{41, 69}, GridPosition{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		// hexagons are right left side inclusive, right side exclusive
		{PixelPosition{42, 50}, GridPosition{1, 2}, true},
		{PixelPosition{42, 51}, GridPosition{1, 2}, true},
		{PixelPosition{42, 52}, GridPosition{1, 2}, true},
		{PixelPosition{42, 53}, GridPosition{1, 2}, true},
		{PixelPosition{42, 54}, GridPosition{1, 2}, true},
		{PixelPosition{42, 55}, GridPosition{1, 2}, true},
		{PixelPosition{42, 56}, GridPosition{1, 2}, true},
		{PixelPosition{42, 57}, GridPosition{1, 2}, true},
		{PixelPosition{42, 58}, GridPosition{1, 2}, true},
		{PixelPosition{42, 59}, GridPosition{1, 2}, true},
		{PixelPosition{42, 60}, GridPosition{1, 2}, true},
		{PixelPosition{42, 61}, GridPosition{1, 2}, true},
		{PixelPosition{42, 62}, GridPosition{1, 2}, true},
		{PixelPosition{42, 63}, GridPosition{1, 2}, true},
		{PixelPosition{42, 64}, GridPosition{1, 2}, true},
		{PixelPosition{42, 65}, GridPosition{1, 2}, true},
		{PixelPosition{42, 66}, GridPosition{1, 2}, true},
		{PixelPosition{42, 67}, GridPosition{1, 2}, true},
		{PixelPosition{42, 68}, GridPosition{1, 2}, true},
		{PixelPosition{42, 69}, GridPosition{1, 2}, true},

		// border of hexagon 0,2 and 1,2
		{PixelPosition{43, 50}, GridPosition{1, 2}, true},
		{PixelPosition{43, 51}, GridPosition{1, 2}, true},
		{PixelPosition{43, 52}, GridPosition{1, 2}, true},
		{PixelPosition{43, 53}, GridPosition{1, 2}, true},
		{PixelPosition{43, 54}, GridPosition{1, 2}, true},
		{PixelPosition{43, 55}, GridPosition{1, 2}, true},
		{PixelPosition{43, 56}, GridPosition{1, 2}, true},
		{PixelPosition{43, 57}, GridPosition{1, 2}, true},
		{PixelPosition{43, 58}, GridPosition{1, 2}, true},
		{PixelPosition{43, 59}, GridPosition{1, 2}, true},
		{PixelPosition{43, 60}, GridPosition{1, 2}, true},
		{PixelPosition{43, 61}, GridPosition{1, 2}, true},
		{PixelPosition{43, 62}, GridPosition{1, 2}, true},
		{PixelPosition{43, 63}, GridPosition{1, 2}, true},
		{PixelPosition{43, 64}, GridPosition{1, 2}, true},
		{PixelPosition{43, 65}, GridPosition{1, 2}, true},
		{PixelPosition{43, 66}, GridPosition{1, 2}, true},
		{PixelPosition{43, 67}, GridPosition{1, 2}, true},
		{PixelPosition{43, 68}, GridPosition{1, 2}, true},
		{PixelPosition{43, 69}, GridPosition{1, 2}, true},

		// right side of hexagon 1,2
		{PixelPosition{64, 50}, GridPosition{1, 2}, true},
		{PixelPosition{64, 51}, GridPosition{1, 2}, true},
		{PixelPosition{64, 52}, GridPosition{1, 2}, true},
		{PixelPosition{64, 53}, GridPosition{1, 2}, true},
		{PixelPosition{64, 54}, GridPosition{1, 2}, true},
		{PixelPosition{64, 55}, GridPosition{1, 2}, true},
		{PixelPosition{64, 56}, GridPosition{1, 2}, true},
		{PixelPosition{64, 57}, GridPosition{1, 2}, true},
		{PixelPosition{64, 58}, GridPosition{1, 2}, true},
		{PixelPosition{64, 59}, GridPosition{1, 2}, true},
		{PixelPosition{64, 60}, GridPosition{1, 2}, true},
		{PixelPosition{64, 61}, GridPosition{1, 2}, true},
		{PixelPosition{64, 62}, GridPosition{1, 2}, true},
		{PixelPosition{64, 63}, GridPosition{1, 2}, true},
		{PixelPosition{64, 64}, GridPosition{1, 2}, true},
		{PixelPosition{64, 65}, GridPosition{1, 2}, true},
		{PixelPosition{64, 66}, GridPosition{1, 2}, true},
		{PixelPosition{64, 67}, GridPosition{1, 2}, true},
		{PixelPosition{64, 68}, GridPosition{1, 2}, true},
		{PixelPosition{64, 69}, GridPosition{1, 2}, true},

		// continue right side of hexagon 1,2
		{PixelPosition{65, 50}, GridPosition{1, 2}, true},
		{PixelPosition{65, 51}, GridPosition{1, 2}, true},
		{PixelPosition{65, 52}, GridPosition{1, 2}, true},
		{PixelPosition{65, 53}, GridPosition{1, 2}, true},
		{PixelPosition{65, 54}, GridPosition{1, 2}, true},
		{PixelPosition{65, 55}, GridPosition{1, 2}, true},
		{PixelPosition{65, 56}, GridPosition{1, 2}, true},
		{PixelPosition{65, 57}, GridPosition{1, 2}, true},
		{PixelPosition{65, 58}, GridPosition{1, 2}, true},
		{PixelPosition{65, 59}, GridPosition{1, 2}, true},
		{PixelPosition{65, 60}, GridPosition{1, 2}, true},
		{PixelPosition{65, 61}, GridPosition{1, 2}, true},
		{PixelPosition{65, 62}, GridPosition{1, 2}, true},
		{PixelPosition{65, 63}, GridPosition{1, 2}, true},
		{PixelPosition{65, 64}, GridPosition{1, 2}, true},
		{PixelPosition{65, 65}, GridPosition{1, 2}, true},
		{PixelPosition{65, 66}, GridPosition{1, 2}, true},
		{PixelPosition{65, 67}, GridPosition{1, 2}, true},
		{PixelPosition{65, 68}, GridPosition{1, 2}, true},
		{PixelPosition{65, 69}, GridPosition{1, 2}, true},

		// continue right side of hexagon hexagon 1,2
		{PixelPosition{66, 50}, GridPosition{1, 2}, true},
		{PixelPosition{66, 51}, GridPosition{1, 2}, true},
		{PixelPosition{66, 52}, GridPosition{1, 2}, true},
		{PixelPosition{66, 53}, GridPosition{1, 2}, true},
		{PixelPosition{66, 54}, GridPosition{1, 2}, true},
		{PixelPosition{66, 55}, GridPosition{1, 2}, true},
		{PixelPosition{66, 56}, GridPosition{1, 2}, true},
		{PixelPosition{66, 57}, GridPosition{1, 2}, true},
		{PixelPosition{66, 58}, GridPosition{1, 2}, true},
		{PixelPosition{66, 59}, GridPosition{1, 2}, true},
		{PixelPosition{66, 60}, GridPosition{1, 2}, true},
		{PixelPosition{66, 61}, GridPosition{1, 2}, true},
		{PixelPosition{66, 62}, GridPosition{1, 2}, true},
		{PixelPosition{66, 63}, GridPosition{1, 2}, true},
		{PixelPosition{66, 64}, GridPosition{1, 2}, true},
		{PixelPosition{66, 65}, GridPosition{1, 2}, true},
		{PixelPosition{66, 66}, GridPosition{1, 2}, true},
		{PixelPosition{66, 67}, GridPosition{1, 2}, true},
		{PixelPosition{66, 68}, GridPosition{1, 2}, true},
		{PixelPosition{66, 69}, GridPosition{1, 2}, true},
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
			tt.wantedPos.X,
			tt.wantedPos.Y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := tt.pixelPos.getGridPos(&hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if !pos.equals(tt.wantedPos) {
				t.Errorf("For pixel position %s, Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.pixelPos.log(),
					tt.wantedPos.X,
					tt.wantedPos.Y,
					pos.X,
					pos.Y)
			}
		})
	}
}

func TestGridPosAdjacent(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(46, 40, 19, 10)

	var tests = []struct {
		gridPos       GridPosition
		hexSide       hexSide
		wantedGridPos GridPosition
		wantedOk      bool
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
		var sWantedOk string
		if tt.wantedOk {
			sWantedOk = "exist"
		} else {
			sWantedOk = "not exist"
		}

		testname := fmt.Sprintf(
			"gridPos:{%d,%d} and hex side %d has adjacent position {%d,%d} and want it to %s",
			tt.gridPos.X,
			tt.gridPos.Y,
			tt.hexSide,
			tt.wantedGridPos.X,
			tt.wantedGridPos.Y,
			sWantedOk)

		t.Run(testname, func(t *testing.T) {
			pos, ok := tt.gridPos.adjacent(tt.hexSide, &hc)
			if ok != tt.wantedOk {
				t.Errorf("Wanted pixel position to %s", sWantedOk)
			}
			if !pos.equals(tt.wantedGridPos) {
				t.Errorf("For grid pos {%d,%d} and hex side %d Wanted grid position {%d,%d} but it was {%d,%d}",
					tt.gridPos.X,
					tt.gridPos.Y,
					tt.hexSide,
					tt.wantedGridPos.X,
					tt.wantedGridPos.Y,
					pos.X,
					pos.Y)
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
