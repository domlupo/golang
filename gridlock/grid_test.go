package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestGetPixel(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		tile         Tile
		hexPoint     hexPoint
		wantPixelPos Pixel
		wantOk       bool
	}{

		// check hexagon {1, 1}, {1, 2}, {2,1} to ensure hex points are correct
		{Tile{1, 1}, hpCenter, Pixel{42, 30}, true},
		{Tile{1, 1}, hpNorth, Pixel{42, 10}, true},
		{Tile{1, 1}, hpNortheast, Pixel{64, 20}, true},
		{Tile{1, 1}, hpSoutheast, Pixel{64, 39}, true},
		{Tile{1, 1}, hpSouth, Pixel{42, 49}, true},
		{Tile{1, 1}, hpSouthwest, Pixel{19, 39}, true},
		{Tile{1, 1}, hpNorthwest, Pixel{19, 20}, true},
		{Tile{2, 1}, hpCenter, Pixel{88, 30}, true},
		{Tile{2, 1}, hpNorth, Pixel{88, 10}, true},
		{Tile{2, 1}, hpNortheast, Pixel{110, 20}, true},
		{Tile{2, 1}, hpSoutheast, Pixel{110, 39}, true},
		{Tile{2, 1}, hpSouth, Pixel{88, 49}, true},
		{Tile{2, 1}, hpSouthwest, Pixel{65, 39}, true},
		{Tile{2, 1}, hpNorthwest, Pixel{65, 20}, true},
		{Tile{1, 2}, hpCenter, Pixel{65, 60}, true},
		{Tile{1, 2}, hpNorth, Pixel{65, 40}, true},
		{Tile{1, 2}, hpNortheast, Pixel{87, 50}, true},
		{Tile{1, 2}, hpSoutheast, Pixel{87, 69}, true},
		{Tile{1, 2}, hpSouth, Pixel{65, 79}, true},
		{Tile{1, 2}, hpSouthwest, Pixel{42, 69}, true},
		{Tile{1, 2}, hpNorthwest, Pixel{42, 50}, true},

		// check hexagon {26,23} to ensure large x and large y Pixels are correct
		{Tile{27, 23}, hpCenter, Pixel{1238, 690}, true},
		{Tile{27, 23}, hpNorth, Pixel{1238, 670}, true},
		{Tile{27, 23}, hpNortheast, Pixel{1260, 680}, true},
		{Tile{27, 23}, hpSoutheast, Pixel{1260, 699}, true},
		{Tile{27, 23}, hpSouth, Pixel{1238, 709}, true},
		{Tile{27, 23}, hpSouthwest, Pixel{1215, 699}, true},
		{Tile{27, 23}, hpNorthwest, Pixel{1215, 680}, true},

		// check top left hexagon {0,0}
		{Tile{0, 0}, hpCenter, Pixel{19, 0}, true},
		{Tile{0, 0}, hpNorth, Pixel{0, 0}, false},
		{Tile{0, 0}, hpNortheast, Pixel{0, 0}, false},
		{Tile{0, 0}, hpSoutheast, Pixel{41, 9}, true},
		{Tile{0, 0}, hpSouth, Pixel{19, 19}, true},
		{Tile{0, 0}, hpSouthwest, Pixel{0, 0}, false},
		{Tile{0, 0}, hpNorthwest, Pixel{0, 0}, false},

		// check bottom left hexagon {0,24}
		{Tile{0, 24}, hpCenter, Pixel{0, 0}, false},
		{Tile{0, 24}, hpNorth, Pixel{19, 700}, true},
		{Tile{0, 24}, hpNortheast, Pixel{41, 710}, true},
		{Tile{0, 24}, hpSoutheast, Pixel{0, 0}, false},
		{Tile{0, 24}, hpSouth, Pixel{0, 0}, false},
		{Tile{0, 24}, hpSouthwest, Pixel{0, 0}, false},
		{Tile{0, 24}, hpNorthwest, Pixel{0, 0}, false},

		// check top right hexagon {27,0}
		{Tile{27, 0}, hpCenter, Pixel{1261, 0}, true},
		{Tile{27, 0}, hpNorth, Pixel{0, 0}, false},
		{Tile{27, 0}, hpNortheast, Pixel{0, 0}, false},
		{Tile{27, 0}, hpSoutheast, Pixel{0, 0}, false},
		{Tile{27, 0}, hpSouth, Pixel{1261, 19}, true},
		{Tile{27, 0}, hpSouthwest, Pixel{1238, 9}, true},
		{Tile{27, 0}, hpNorthwest, Pixel{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{Tile{27, 24}, hpCenter, Pixel{0, 0}, false},
		{Tile{27, 24}, hpNorth, Pixel{1261, 700}, true},
		{Tile{27, 24}, hpNortheast, Pixel{0, 0}, false},
		{Tile{27, 24}, hpSoutheast, Pixel{0, 0}, false},
		{Tile{27, 24}, hpSouth, Pixel{0, 0}, false},
		{Tile{27, 24}, hpSouthwest, Pixel{0, 0}, false},
		{Tile{27, 24}, hpNorthwest, Pixel{1238, 710}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"tile: %s and hex point %s converts to hex tile %s and want it to %s",
			tt.tile.log(),
			tt.hexPoint.string(),
			tt.wantPixelPos.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.getPixel(tt.tile, tt.hexPoint)
			if ok != tt.wantOk {
				t.Errorf("Want Pixel to %s", sWantOk)
			}
			if !pos.equals(tt.wantPixelPos) {
				t.Errorf("For tile %s and hex point %s, want tile %s but it was %s",
					tt.tile.log(),
					tt.hexPoint.string(),
					tt.wantPixelPos.log(),
					pos.log())
			}
		})
	}
}

func TestGetTile(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		pixel   Pixel
		wantPos Tile
		wantOk  bool
	}{
		// top left of entire hex grid
		{Pixel{0, 0}, Tile{0, 0}, true},
		{Pixel{-1, 0}, Tile{0, 0}, false},
		{Pixel{0, -1}, Tile{0, 0}, false},
		{Pixel{-1, -1}, Tile{0, 0}, false},

		// first left border point hex change around 0,12 Pixel
		{Pixel{-1, 11}, Tile{0, 0}, false},
		{Pixel{-1, 12}, Tile{0, 0}, false},
		{Pixel{-1, 13}, Tile{0, 0}, false},
		{Pixel{0, 11}, Tile{0, 0}, true},
		{Pixel{0, 12}, Tile{0, 1}, true},
		{Pixel{0, 13}, Tile{0, 1}, true},
		{Pixel{1, 11}, Tile{0, 0}, true},
		{Pixel{1, 12}, Tile{0, 1}, true},
		{Pixel{1, 13}, Tile{0, 1}, true},

		// second left border point hex change around 0,48 Pixel
		{Pixel{-1, 47}, Tile{0, 0}, false},
		{Pixel{-1, 48}, Tile{0, 0}, false},
		{Pixel{-1, 49}, Tile{0, 0}, false},
		{Pixel{0, 47}, Tile{0, 1}, true},
		{Pixel{0, 48}, Tile{0, 2}, true},
		{Pixel{0, 49}, Tile{0, 2}, true},
		{Pixel{1, 47}, Tile{0, 1}, true},
		{Pixel{1, 48}, Tile{0, 2}, true},
		{Pixel{1, 49}, Tile{0, 2}, true},

		// third left border point hex change around 0,72 Pixel
		{Pixel{-1, 71}, Tile{0, 0}, false},
		{Pixel{-1, 72}, Tile{0, 0}, false},
		{Pixel{-1, 73}, Tile{0, 0}, false},
		{Pixel{0, 71}, Tile{0, 2}, true},
		{Pixel{0, 72}, Tile{0, 3}, true},
		{Pixel{0, 73}, Tile{0, 3}, true},
		{Pixel{1, 71}, Tile{0, 2}, true},
		{Pixel{1, 72}, Tile{0, 3}, true},
		{Pixel{1, 73}, Tile{0, 3}, true},

		// first top border point hex change around 42,0 Pixel
		{Pixel{41, -1}, Tile{0, 0}, false},
		{Pixel{41, 0}, Tile{0, 0}, true},
		{Pixel{41, 1}, Tile{0, 0}, true},
		{Pixel{42, -1}, Tile{0, 0}, false},
		{Pixel{42, 0}, Tile{1, 0}, true},
		{Pixel{42, 1}, Tile{1, 0}, true},
		{Pixel{43, -1}, Tile{0, 0}, false},
		{Pixel{43, 0}, Tile{1, 0}, true},
		{Pixel{43, 1}, Tile{1, 0}, true},

		// second top border point hex change around 88, 0 Pixel
		{Pixel{87, -1}, Tile{0, 0}, false},
		{Pixel{87, 0}, Tile{1, 0}, true},
		{Pixel{87, 1}, Tile{1, 0}, true},
		{Pixel{88, -1}, Tile{0, 0}, false},
		{Pixel{88, 0}, Tile{2, 0}, true},
		{Pixel{88, 1}, Tile{2, 0}, true},
		{Pixel{89, -1}, Tile{0, 0}, false},
		{Pixel{89, 0}, Tile{2, 0}, true},
		{Pixel{89, 1}, Tile{2, 0}, true},

		// third top border point hex change around 134,0 Pixel
		{Pixel{133, -1}, Tile{0, 0}, false},
		{Pixel{133, 0}, Tile{2, 0}, true},
		{Pixel{133, 1}, Tile{2, 0}, true},
		{Pixel{134, -1}, Tile{0, 0}, false},
		{Pixel{134, 0}, Tile{3, 0}, true},
		{Pixel{134, 1}, Tile{3, 0}, true},
		{Pixel{135, -1}, Tile{0, 0}, false},
		{Pixel{135, 0}, Tile{3, 0}, true},
		{Pixel{135, 1}, Tile{3, 0}, true},

		// hexagon points near 19, 80 Pixel
		{Pixel{18, 79}, Tile{0, 2}, true},
		{Pixel{18, 80}, Tile{0, 3}, true},
		{Pixel{18, 81}, Tile{0, 3}, true},
		{Pixel{19, 79}, Tile{0, 2}, true},
		{Pixel{19, 80}, Tile{1, 3}, true},
		{Pixel{19, 81}, Tile{1, 3}, true},
		{Pixel{20, 79}, Tile{1, 3}, true},
		{Pixel{20, 80}, Tile{1, 3}, true},
		{Pixel{20, 81}, Tile{1, 3}, true},

		// hexagon points near 42, 70 Pixel
		{Pixel{41, 69}, Tile{0, 2}, true},
		{Pixel{41, 70}, Tile{1, 3}, true},
		{Pixel{41, 71}, Tile{1, 3}, true},
		{Pixel{42, 69}, Tile{1, 2}, true},
		{Pixel{42, 70}, Tile{1, 3}, true},
		{Pixel{42, 71}, Tile{1, 3}, true},
		{Pixel{43, 69}, Tile{1, 2}, true},
		{Pixel{43, 70}, Tile{1, 2}, true},
		{Pixel{43, 71}, Tile{1, 3}, true},

		// hexagon points near 65, 80 Pixel
		{Pixel{64, 79}, Tile{1, 2}, true},
		{Pixel{64, 80}, Tile{1, 3}, true},
		{Pixel{64, 81}, Tile{1, 3}, true},
		{Pixel{65, 79}, Tile{1, 2}, true},
		{Pixel{65, 80}, Tile{2, 3}, true},
		{Pixel{65, 81}, Tile{2, 3}, true},
		{Pixel{66, 79}, Tile{2, 3}, true},
		{Pixel{66, 80}, Tile{2, 3}, true},
		{Pixel{66, 81}, Tile{2, 3}, true},

		// hexagon points near 88, 10 Pixel
		{Pixel{87, 9}, Tile{1, 0}, true},
		{Pixel{87, 10}, Tile{2, 1}, true},
		{Pixel{87, 11}, Tile{2, 1}, true},
		{Pixel{88, 9}, Tile{2, 0}, true},
		{Pixel{88, 10}, Tile{2, 1}, true},
		{Pixel{88, 11}, Tile{2, 1}, true},
		{Pixel{89, 9}, Tile{2, 0}, true},
		{Pixel{89, 10}, Tile{2, 0}, true},
		{Pixel{89, 11}, Tile{2, 1}, true},

		// hexagon points near 88, 50 Pixel
		{Pixel{87, 49}, Tile{2, 1}, true},
		{Pixel{87, 50}, Tile{1, 2}, true},
		{Pixel{87, 51}, Tile{1, 2}, true},
		{Pixel{88, 49}, Tile{2, 1}, true},
		{Pixel{88, 50}, Tile{2, 2}, true},
		{Pixel{88, 51}, Tile{2, 2}, true},
		{Pixel{89, 49}, Tile{2, 2}, true},
		{Pixel{89, 50}, Tile{2, 2}, true},
		{Pixel{89, 51}, Tile{2, 2}, true},

		// hexagon points near 88, 70 Pixel
		{Pixel{87, 69}, Tile{1, 2}, true},
		{Pixel{87, 70}, Tile{2, 3}, true},
		{Pixel{87, 71}, Tile{2, 3}, true},
		{Pixel{88, 69}, Tile{2, 2}, true},
		{Pixel{88, 70}, Tile{2, 3}, true},
		{Pixel{88, 71}, Tile{2, 3}, true},
		{Pixel{89, 69}, Tile{2, 2}, true},
		{Pixel{89, 70}, Tile{2, 2}, true},
		{Pixel{89, 71}, Tile{2, 3}, true},

		// bottom left of entire hex grid
		{Pixel{-1, 719}, Tile{0, 0}, false},
		{Pixel{-1, 720}, Tile{0, 0}, false},
		{Pixel{0, 719}, Tile{0, 24}, true},
		{Pixel{0, 720}, Tile{0, 0}, false},

		// top right of entire hex grid
		{Pixel{1279, -1}, Tile{0, 0}, false},
		{Pixel{1280, -1}, Tile{0, 0}, false},
		{Pixel{1279, 0}, Tile{27, 0}, true},
		{Pixel{1280, 0}, Tile{0, 0}, false},

		// bottom right of entire hex grid
		{Pixel{1279, 719}, Tile{27, 24}, true},
		{Pixel{1280, 719}, Tile{0, 0}, false},
		{Pixel{1279, 720}, Tile{0, 0}, false},
		{Pixel{1280, 720}, Tile{0, 0}, false},

		// first hexagon y section (hexagon look up table)
		// between hexagons 0,0 1,1 and 1,0

		// sanity checks
		{Pixel{19, 10}, Tile{0, 0}, true},
		{Pixel{65, 10}, Tile{1, 0}, true},
		{Pixel{41, 20}, Tile{1, 1}, true},
		{Pixel{42, 20}, Tile{1, 1}, true},
		{Pixel{43, 20}, Tile{1, 1}, true},

		// 0,0 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{19, 19}, Tile{0, 0}, true},
		{Pixel{20, 18}, Tile{0, 0}, true},
		{Pixel{21, 18}, Tile{0, 0}, true},
		{Pixel{22, 17}, Tile{0, 0}, true},
		{Pixel{23, 17}, Tile{0, 0}, true},
		{Pixel{24, 17}, Tile{0, 0}, true},
		{Pixel{25, 16}, Tile{0, 0}, true},
		{Pixel{26, 16}, Tile{0, 0}, true},
		{Pixel{27, 15}, Tile{0, 0}, true},
		{Pixel{28, 15}, Tile{0, 0}, true},
		{Pixel{29, 14}, Tile{0, 0}, true},
		{Pixel{30, 14}, Tile{0, 0}, true},
		{Pixel{31, 14}, Tile{0, 0}, true},
		{Pixel{32, 13}, Tile{0, 0}, true},
		{Pixel{33, 13}, Tile{0, 0}, true},
		{Pixel{34, 12}, Tile{0, 0}, true},
		{Pixel{35, 12}, Tile{0, 0}, true},
		{Pixel{36, 11}, Tile{0, 0}, true},
		{Pixel{37, 11}, Tile{0, 0}, true},
		{Pixel{38, 11}, Tile{0, 0}, true},
		{Pixel{39, 10}, Tile{0, 0}, true},
		{Pixel{40, 10}, Tile{0, 0}, true},
		{Pixel{41, 9}, Tile{0, 0}, true},

		// 1, 1 hexagon slope check between itself and 0, 0 hexagon
		{Pixel{19, 20}, Tile{1, 1}, true},
		{Pixel{20, 19}, Tile{1, 1}, true},
		{Pixel{21, 19}, Tile{1, 1}, true},
		{Pixel{22, 18}, Tile{1, 1}, true},
		{Pixel{23, 18}, Tile{1, 1}, true},
		{Pixel{24, 18}, Tile{1, 1}, true},
		{Pixel{25, 17}, Tile{1, 1}, true},
		{Pixel{26, 17}, Tile{1, 1}, true},
		{Pixel{27, 16}, Tile{1, 1}, true},
		{Pixel{28, 16}, Tile{1, 1}, true},
		{Pixel{29, 15}, Tile{1, 1}, true},
		{Pixel{30, 15}, Tile{1, 1}, true},
		{Pixel{31, 15}, Tile{1, 1}, true},
		{Pixel{32, 14}, Tile{1, 1}, true},
		{Pixel{33, 14}, Tile{1, 1}, true},
		{Pixel{34, 13}, Tile{1, 1}, true},
		{Pixel{35, 13}, Tile{1, 1}, true},
		{Pixel{36, 12}, Tile{1, 1}, true},
		{Pixel{37, 12}, Tile{1, 1}, true},
		{Pixel{38, 12}, Tile{1, 1}, true},
		{Pixel{39, 11}, Tile{1, 1}, true},
		{Pixel{40, 11}, Tile{1, 1}, true},
		{Pixel{41, 10}, Tile{1, 1}, true},

		// 1,0 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{42, 9}, Tile{1, 0}, true},
		{Pixel{43, 10}, Tile{1, 0}, true},
		{Pixel{44, 10}, Tile{1, 0}, true},
		{Pixel{45, 11}, Tile{1, 0}, true},
		{Pixel{46, 11}, Tile{1, 0}, true},
		{Pixel{47, 11}, Tile{1, 0}, true},
		{Pixel{48, 12}, Tile{1, 0}, true},
		{Pixel{49, 12}, Tile{1, 0}, true},
		{Pixel{50, 13}, Tile{1, 0}, true},
		{Pixel{51, 13}, Tile{1, 0}, true},

		{Pixel{52, 14}, Tile{1, 0}, true},
		{Pixel{53, 14}, Tile{1, 0}, true},
		{Pixel{54, 14}, Tile{1, 0}, true},
		{Pixel{55, 15}, Tile{1, 0}, true},
		{Pixel{56, 15}, Tile{1, 0}, true},
		{Pixel{57, 16}, Tile{1, 0}, true},
		{Pixel{58, 16}, Tile{1, 0}, true},
		{Pixel{59, 17}, Tile{1, 0}, true},
		{Pixel{60, 17}, Tile{1, 0}, true},
		{Pixel{61, 17}, Tile{1, 0}, true},
		{Pixel{62, 18}, Tile{1, 0}, true},
		{Pixel{63, 18}, Tile{1, 0}, true},
		{Pixel{64, 19}, Tile{1, 0}, true},

		// 1,1 hexagon slope check between itself and 1, 0 hexagon
		{Pixel{42, 10}, Tile{1, 1}, true},
		{Pixel{43, 11}, Tile{1, 1}, true},
		{Pixel{44, 11}, Tile{1, 1}, true},
		{Pixel{45, 12}, Tile{1, 1}, true},
		{Pixel{46, 12}, Tile{1, 1}, true},
		{Pixel{47, 12}, Tile{1, 1}, true},
		{Pixel{48, 13}, Tile{1, 1}, true},
		{Pixel{49, 13}, Tile{1, 1}, true},
		{Pixel{50, 14}, Tile{1, 1}, true},
		{Pixel{51, 14}, Tile{1, 1}, true},
		{Pixel{52, 15}, Tile{1, 1}, true},
		{Pixel{53, 15}, Tile{1, 1}, true},
		{Pixel{54, 15}, Tile{1, 1}, true},
		{Pixel{55, 16}, Tile{1, 1}, true},
		{Pixel{56, 16}, Tile{1, 1}, true},
		{Pixel{57, 17}, Tile{1, 1}, true},
		{Pixel{58, 17}, Tile{1, 1}, true},
		{Pixel{59, 18}, Tile{1, 1}, true},
		{Pixel{60, 18}, Tile{1, 1}, true},
		{Pixel{61, 18}, Tile{1, 1}, true},
		{Pixel{62, 19}, Tile{1, 1}, true},
		{Pixel{63, 19}, Tile{1, 1}, true},
		{Pixel{64, 20}, Tile{1, 1}, true},

		// second hexagon y section (full rectangle of hexagon 1,1)

		// right side of hexagon 0,1
		{Pixel{18, 20}, Tile{0, 1}, true},
		{Pixel{18, 21}, Tile{0, 1}, true},
		{Pixel{18, 22}, Tile{0, 1}, true},
		{Pixel{18, 23}, Tile{0, 1}, true},
		{Pixel{18, 24}, Tile{0, 1}, true},
		{Pixel{18, 25}, Tile{0, 1}, true},
		{Pixel{18, 26}, Tile{0, 1}, true},
		{Pixel{18, 27}, Tile{0, 1}, true},
		{Pixel{18, 28}, Tile{0, 1}, true},
		{Pixel{18, 29}, Tile{0, 1}, true},
		{Pixel{18, 30}, Tile{0, 1}, true},
		{Pixel{18, 31}, Tile{0, 1}, true},
		{Pixel{18, 32}, Tile{0, 1}, true},
		{Pixel{18, 33}, Tile{0, 1}, true},
		{Pixel{18, 34}, Tile{0, 1}, true},
		{Pixel{18, 35}, Tile{0, 1}, true},
		{Pixel{18, 36}, Tile{0, 1}, true},
		{Pixel{18, 37}, Tile{0, 1}, true},
		{Pixel{18, 38}, Tile{0, 1}, true},
		{Pixel{18, 39}, Tile{0, 1}, true},

		// left side of hexagon 1,1
		// hexagons are right side inclusive and left side exclusive
		{Pixel{19, 20}, Tile{1, 1}, true},
		{Pixel{19, 21}, Tile{1, 1}, true},
		{Pixel{19, 22}, Tile{1, 1}, true},
		{Pixel{19, 23}, Tile{1, 1}, true},
		{Pixel{19, 24}, Tile{1, 1}, true},
		{Pixel{19, 25}, Tile{1, 1}, true},
		{Pixel{19, 26}, Tile{1, 1}, true},
		{Pixel{19, 27}, Tile{1, 1}, true},
		{Pixel{19, 28}, Tile{1, 1}, true},
		{Pixel{19, 29}, Tile{1, 1}, true},
		{Pixel{19, 30}, Tile{1, 1}, true},
		{Pixel{19, 31}, Tile{1, 1}, true},
		{Pixel{19, 32}, Tile{1, 1}, true},
		{Pixel{19, 33}, Tile{1, 1}, true},
		{Pixel{19, 34}, Tile{1, 1}, true},
		{Pixel{19, 35}, Tile{1, 1}, true},
		{Pixel{19, 36}, Tile{1, 1}, true},
		{Pixel{19, 37}, Tile{1, 1}, true},
		{Pixel{19, 38}, Tile{1, 1}, true},
		{Pixel{19, 39}, Tile{1, 1}, true},
		// continue left side of hexagon 1,1
		{Pixel{20, 20}, Tile{1, 1}, true},
		{Pixel{20, 21}, Tile{1, 1}, true},
		{Pixel{20, 22}, Tile{1, 1}, true},
		{Pixel{20, 23}, Tile{1, 1}, true},
		{Pixel{20, 24}, Tile{1, 1}, true},
		{Pixel{20, 25}, Tile{1, 1}, true},
		{Pixel{20, 26}, Tile{1, 1}, true},
		{Pixel{20, 27}, Tile{1, 1}, true},
		{Pixel{20, 28}, Tile{1, 1}, true},
		{Pixel{20, 29}, Tile{1, 1}, true},
		{Pixel{20, 30}, Tile{1, 1}, true},
		{Pixel{20, 31}, Tile{1, 1}, true},
		{Pixel{20, 32}, Tile{1, 1}, true},
		{Pixel{20, 33}, Tile{1, 1}, true},
		{Pixel{20, 34}, Tile{1, 1}, true},
		{Pixel{20, 35}, Tile{1, 1}, true},
		{Pixel{20, 36}, Tile{1, 1}, true},
		{Pixel{20, 37}, Tile{1, 1}, true},
		{Pixel{20, 38}, Tile{1, 1}, true},
		{Pixel{20, 39}, Tile{1, 1}, true},

		// right side of hexagon 1,1
		{Pixel{64, 20}, Tile{1, 1}, true},
		{Pixel{64, 21}, Tile{1, 1}, true},
		{Pixel{64, 22}, Tile{1, 1}, true},
		{Pixel{64, 23}, Tile{1, 1}, true},
		{Pixel{64, 24}, Tile{1, 1}, true},
		{Pixel{64, 25}, Tile{1, 1}, true},
		{Pixel{64, 26}, Tile{1, 1}, true},
		{Pixel{64, 27}, Tile{1, 1}, true},
		{Pixel{64, 28}, Tile{1, 1}, true},
		{Pixel{64, 29}, Tile{1, 1}, true},
		{Pixel{64, 30}, Tile{1, 1}, true},
		{Pixel{64, 31}, Tile{1, 1}, true},
		{Pixel{64, 32}, Tile{1, 1}, true},
		{Pixel{64, 33}, Tile{1, 1}, true},
		{Pixel{64, 34}, Tile{1, 1}, true},
		{Pixel{64, 35}, Tile{1, 1}, true},
		{Pixel{64, 36}, Tile{1, 1}, true},
		{Pixel{64, 37}, Tile{1, 1}, true},
		{Pixel{64, 38}, Tile{1, 1}, true},
		{Pixel{64, 39}, Tile{1, 1}, true},

		// left side of hexagon 2,1
		// hexagons are right side inclusive and left side exclusive
		{Pixel{65, 20}, Tile{2, 1}, true},
		{Pixel{65, 21}, Tile{2, 1}, true},
		{Pixel{65, 22}, Tile{2, 1}, true},
		{Pixel{65, 23}, Tile{2, 1}, true},
		{Pixel{65, 24}, Tile{2, 1}, true},
		{Pixel{65, 25}, Tile{2, 1}, true},
		{Pixel{65, 26}, Tile{2, 1}, true},
		{Pixel{65, 27}, Tile{2, 1}, true},
		{Pixel{65, 28}, Tile{2, 1}, true},
		{Pixel{65, 29}, Tile{2, 1}, true},
		{Pixel{65, 30}, Tile{2, 1}, true},
		{Pixel{65, 31}, Tile{2, 1}, true},
		{Pixel{65, 32}, Tile{2, 1}, true},
		{Pixel{65, 33}, Tile{2, 1}, true},
		{Pixel{65, 34}, Tile{2, 1}, true},
		{Pixel{65, 35}, Tile{2, 1}, true},
		{Pixel{65, 36}, Tile{2, 1}, true},
		{Pixel{65, 37}, Tile{2, 1}, true},
		{Pixel{65, 38}, Tile{2, 1}, true},
		{Pixel{65, 39}, Tile{2, 1}, true},
		// continue left side of hexagon 2,1
		{Pixel{66, 20}, Tile{2, 1}, true},
		{Pixel{66, 21}, Tile{2, 1}, true},
		{Pixel{66, 22}, Tile{2, 1}, true},
		{Pixel{66, 23}, Tile{2, 1}, true},
		{Pixel{66, 24}, Tile{2, 1}, true},
		{Pixel{66, 25}, Tile{2, 1}, true},
		{Pixel{66, 26}, Tile{2, 1}, true},
		{Pixel{66, 27}, Tile{2, 1}, true},
		{Pixel{66, 28}, Tile{2, 1}, true},
		{Pixel{66, 29}, Tile{2, 1}, true},
		{Pixel{66, 30}, Tile{2, 1}, true},
		{Pixel{66, 31}, Tile{2, 1}, true},
		{Pixel{66, 32}, Tile{2, 1}, true},
		{Pixel{66, 33}, Tile{2, 1}, true},
		{Pixel{66, 34}, Tile{2, 1}, true},
		{Pixel{66, 35}, Tile{2, 1}, true},
		{Pixel{66, 36}, Tile{2, 1}, true},
		{Pixel{66, 37}, Tile{2, 1}, true},
		{Pixel{66, 38}, Tile{2, 1}, true},
		{Pixel{66, 39}, Tile{2, 1}, true},

		// third hexagon y section (inverted hexagon look up table)
		// between hexagons 0,2 1,1 and 1,2

		// sanity checks
		{Pixel{19, 49}, Tile{0, 2}, true},
		{Pixel{65, 49}, Tile{1, 2}, true},
		{Pixel{41, 30}, Tile{1, 1}, true},
		{Pixel{42, 30}, Tile{1, 1}, true},
		{Pixel{43, 30}, Tile{1, 1}, true},

		// 0,2 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{19, 40}, Tile{0, 2}, true},
		{Pixel{20, 41}, Tile{0, 2}, true},
		{Pixel{21, 41}, Tile{0, 2}, true},
		{Pixel{22, 42}, Tile{0, 2}, true},
		{Pixel{23, 42}, Tile{0, 2}, true},
		{Pixel{24, 42}, Tile{0, 2}, true},
		{Pixel{25, 43}, Tile{0, 2}, true},
		{Pixel{26, 43}, Tile{0, 2}, true},
		{Pixel{27, 44}, Tile{0, 2}, true},
		{Pixel{28, 44}, Tile{0, 2}, true},
		{Pixel{29, 45}, Tile{0, 2}, true},
		{Pixel{30, 45}, Tile{0, 2}, true},
		{Pixel{31, 45}, Tile{0, 2}, true},
		{Pixel{32, 46}, Tile{0, 2}, true},
		{Pixel{33, 46}, Tile{0, 2}, true},
		{Pixel{34, 47}, Tile{0, 2}, true},
		{Pixel{35, 47}, Tile{0, 2}, true},
		{Pixel{36, 48}, Tile{0, 2}, true},
		{Pixel{37, 48}, Tile{0, 2}, true},
		{Pixel{38, 48}, Tile{0, 2}, true},
		{Pixel{39, 49}, Tile{0, 2}, true},
		{Pixel{40, 49}, Tile{0, 2}, true},
		{Pixel{41, 50}, Tile{0, 2}, true},

		// 1, 1 hexagon slope check between itself and 0, 2 hexagon
		{Pixel{19, 39}, Tile{1, 1}, true},
		{Pixel{20, 40}, Tile{1, 1}, true},
		{Pixel{21, 40}, Tile{1, 1}, true},
		{Pixel{22, 41}, Tile{1, 1}, true},
		{Pixel{23, 41}, Tile{1, 1}, true},
		{Pixel{24, 41}, Tile{1, 1}, true},
		{Pixel{25, 42}, Tile{1, 1}, true},
		{Pixel{26, 42}, Tile{1, 1}, true},
		{Pixel{27, 43}, Tile{1, 1}, true},
		{Pixel{28, 43}, Tile{1, 1}, true},
		{Pixel{29, 44}, Tile{1, 1}, true},
		{Pixel{30, 44}, Tile{1, 1}, true},
		{Pixel{31, 44}, Tile{1, 1}, true},
		{Pixel{32, 45}, Tile{1, 1}, true},
		{Pixel{33, 45}, Tile{1, 1}, true},
		{Pixel{34, 46}, Tile{1, 1}, true},
		{Pixel{35, 46}, Tile{1, 1}, true},
		{Pixel{36, 47}, Tile{1, 1}, true},
		{Pixel{37, 47}, Tile{1, 1}, true},
		{Pixel{38, 47}, Tile{1, 1}, true},
		{Pixel{39, 48}, Tile{1, 1}, true},
		{Pixel{40, 48}, Tile{1, 1}, true},
		{Pixel{41, 49}, Tile{1, 1}, true},

		// 1,1 hexagon slope check between itself and 1, 2 hexagon
		{Pixel{42, 49}, Tile{1, 1}, true},
		{Pixel{43, 48}, Tile{1, 1}, true},
		{Pixel{44, 48}, Tile{1, 1}, true},
		{Pixel{45, 47}, Tile{1, 1}, true},
		{Pixel{46, 47}, Tile{1, 1}, true},
		{Pixel{47, 47}, Tile{1, 1}, true},
		{Pixel{48, 46}, Tile{1, 1}, true},
		{Pixel{49, 46}, Tile{1, 1}, true},
		{Pixel{50, 45}, Tile{1, 1}, true},
		{Pixel{51, 45}, Tile{1, 1}, true},
		{Pixel{52, 44}, Tile{1, 1}, true},
		{Pixel{53, 44}, Tile{1, 1}, true},
		{Pixel{54, 44}, Tile{1, 1}, true},
		{Pixel{55, 43}, Tile{1, 1}, true},
		{Pixel{56, 43}, Tile{1, 1}, true},
		{Pixel{57, 42}, Tile{1, 1}, true},
		{Pixel{58, 42}, Tile{1, 1}, true},
		{Pixel{59, 41}, Tile{1, 1}, true},
		{Pixel{60, 41}, Tile{1, 1}, true},
		{Pixel{61, 41}, Tile{1, 1}, true},
		{Pixel{62, 40}, Tile{1, 1}, true},
		{Pixel{63, 40}, Tile{1, 1}, true},
		{Pixel{64, 39}, Tile{1, 1}, true},

		// 1,2 hexagon slope check between itself and 1, 1 hexagon
		{Pixel{42, 50}, Tile{1, 2}, true},
		{Pixel{43, 49}, Tile{1, 2}, true},
		{Pixel{44, 49}, Tile{1, 2}, true},
		{Pixel{45, 48}, Tile{1, 2}, true},
		{Pixel{46, 48}, Tile{1, 2}, true},
		{Pixel{47, 48}, Tile{1, 2}, true},
		{Pixel{48, 47}, Tile{1, 2}, true},
		{Pixel{49, 47}, Tile{1, 2}, true},
		{Pixel{50, 46}, Tile{1, 2}, true},
		{Pixel{51, 46}, Tile{1, 2}, true},
		{Pixel{52, 45}, Tile{1, 2}, true},
		{Pixel{53, 45}, Tile{1, 2}, true},
		{Pixel{54, 45}, Tile{1, 2}, true},
		{Pixel{55, 44}, Tile{1, 2}, true},
		{Pixel{56, 44}, Tile{1, 2}, true},
		{Pixel{57, 43}, Tile{1, 2}, true},
		{Pixel{58, 43}, Tile{1, 2}, true},
		{Pixel{59, 42}, Tile{1, 2}, true},
		{Pixel{60, 42}, Tile{1, 2}, true},
		{Pixel{61, 42}, Tile{1, 2}, true},
		{Pixel{62, 41}, Tile{1, 2}, true},
		{Pixel{63, 41}, Tile{1, 2}, true},
		{Pixel{64, 40}, Tile{1, 2}, true},

		// fourth hexagon y section (split rectangle between hexagon 0,2 and 1,2)

		// middle of hexagon 0,2
		{Pixel{18, 50}, Tile{0, 2}, true},
		{Pixel{18, 51}, Tile{0, 2}, true},
		{Pixel{18, 52}, Tile{0, 2}, true},
		{Pixel{18, 53}, Tile{0, 2}, true},
		{Pixel{18, 54}, Tile{0, 2}, true},
		{Pixel{18, 55}, Tile{0, 2}, true},
		{Pixel{18, 56}, Tile{0, 2}, true},
		{Pixel{18, 57}, Tile{0, 2}, true},
		{Pixel{18, 58}, Tile{0, 2}, true},
		{Pixel{18, 59}, Tile{0, 2}, true},
		{Pixel{18, 60}, Tile{0, 2}, true},
		{Pixel{18, 61}, Tile{0, 2}, true},
		{Pixel{18, 62}, Tile{0, 2}, true},
		{Pixel{18, 63}, Tile{0, 2}, true},
		{Pixel{18, 64}, Tile{0, 2}, true},
		{Pixel{18, 65}, Tile{0, 2}, true},
		{Pixel{18, 66}, Tile{0, 2}, true},
		{Pixel{18, 67}, Tile{0, 2}, true},
		{Pixel{18, 68}, Tile{0, 2}, true},
		{Pixel{18, 69}, Tile{0, 2}, true},

		// continue middle of hexagon 0,2
		{Pixel{19, 50}, Tile{0, 2}, true},
		{Pixel{19, 51}, Tile{0, 2}, true},
		{Pixel{19, 52}, Tile{0, 2}, true},
		{Pixel{19, 53}, Tile{0, 2}, true},
		{Pixel{19, 54}, Tile{0, 2}, true},
		{Pixel{19, 55}, Tile{0, 2}, true},
		{Pixel{19, 56}, Tile{0, 2}, true},
		{Pixel{19, 57}, Tile{0, 2}, true},
		{Pixel{19, 58}, Tile{0, 2}, true},
		{Pixel{19, 59}, Tile{0, 2}, true},
		{Pixel{19, 60}, Tile{0, 2}, true},
		{Pixel{19, 61}, Tile{0, 2}, true},
		{Pixel{19, 62}, Tile{0, 2}, true},
		{Pixel{19, 63}, Tile{0, 2}, true},
		{Pixel{19, 64}, Tile{0, 2}, true},
		{Pixel{19, 65}, Tile{0, 2}, true},
		{Pixel{19, 66}, Tile{0, 2}, true},
		{Pixel{19, 67}, Tile{0, 2}, true},
		{Pixel{19, 68}, Tile{0, 2}, true},
		{Pixel{19, 69}, Tile{0, 2}, true},

		// continue middle of hexagon 0,2
		{Pixel{20, 50}, Tile{0, 2}, true},
		{Pixel{20, 51}, Tile{0, 2}, true},
		{Pixel{20, 52}, Tile{0, 2}, true},
		{Pixel{20, 53}, Tile{0, 2}, true},
		{Pixel{20, 54}, Tile{0, 2}, true},
		{Pixel{20, 55}, Tile{0, 2}, true},
		{Pixel{20, 56}, Tile{0, 2}, true},
		{Pixel{20, 57}, Tile{0, 2}, true},
		{Pixel{20, 58}, Tile{0, 2}, true},
		{Pixel{20, 59}, Tile{0, 2}, true},
		{Pixel{20, 60}, Tile{0, 2}, true},
		{Pixel{20, 61}, Tile{0, 2}, true},
		{Pixel{20, 62}, Tile{0, 2}, true},
		{Pixel{20, 63}, Tile{0, 2}, true},
		{Pixel{20, 64}, Tile{0, 2}, true},
		{Pixel{20, 65}, Tile{0, 2}, true},
		{Pixel{20, 66}, Tile{0, 2}, true},
		{Pixel{20, 67}, Tile{0, 2}, true},
		{Pixel{20, 68}, Tile{0, 2}, true},
		{Pixel{20, 69}, Tile{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		{Pixel{41, 50}, Tile{0, 2}, true},
		{Pixel{41, 51}, Tile{0, 2}, true},
		{Pixel{41, 52}, Tile{0, 2}, true},
		{Pixel{41, 53}, Tile{0, 2}, true},
		{Pixel{41, 54}, Tile{0, 2}, true},
		{Pixel{41, 55}, Tile{0, 2}, true},
		{Pixel{41, 56}, Tile{0, 2}, true},
		{Pixel{41, 57}, Tile{0, 2}, true},
		{Pixel{41, 58}, Tile{0, 2}, true},
		{Pixel{41, 59}, Tile{0, 2}, true},
		{Pixel{41, 60}, Tile{0, 2}, true},
		{Pixel{41, 61}, Tile{0, 2}, true},
		{Pixel{41, 62}, Tile{0, 2}, true},
		{Pixel{41, 63}, Tile{0, 2}, true},
		{Pixel{41, 64}, Tile{0, 2}, true},
		{Pixel{41, 65}, Tile{0, 2}, true},
		{Pixel{41, 66}, Tile{0, 2}, true},
		{Pixel{41, 67}, Tile{0, 2}, true},
		{Pixel{41, 68}, Tile{0, 2}, true},
		{Pixel{41, 69}, Tile{0, 2}, true},

		// border of hexagon 0,2 and 1,2
		// hexagons are right left side inclusive, right side exclusive
		{Pixel{42, 50}, Tile{1, 2}, true},
		{Pixel{42, 51}, Tile{1, 2}, true},
		{Pixel{42, 52}, Tile{1, 2}, true},
		{Pixel{42, 53}, Tile{1, 2}, true},
		{Pixel{42, 54}, Tile{1, 2}, true},
		{Pixel{42, 55}, Tile{1, 2}, true},
		{Pixel{42, 56}, Tile{1, 2}, true},
		{Pixel{42, 57}, Tile{1, 2}, true},
		{Pixel{42, 58}, Tile{1, 2}, true},
		{Pixel{42, 59}, Tile{1, 2}, true},
		{Pixel{42, 60}, Tile{1, 2}, true},
		{Pixel{42, 61}, Tile{1, 2}, true},
		{Pixel{42, 62}, Tile{1, 2}, true},
		{Pixel{42, 63}, Tile{1, 2}, true},
		{Pixel{42, 64}, Tile{1, 2}, true},
		{Pixel{42, 65}, Tile{1, 2}, true},
		{Pixel{42, 66}, Tile{1, 2}, true},
		{Pixel{42, 67}, Tile{1, 2}, true},
		{Pixel{42, 68}, Tile{1, 2}, true},
		{Pixel{42, 69}, Tile{1, 2}, true},

		// border of hexagon 0,2 and 1,2
		{Pixel{43, 50}, Tile{1, 2}, true},
		{Pixel{43, 51}, Tile{1, 2}, true},
		{Pixel{43, 52}, Tile{1, 2}, true},
		{Pixel{43, 53}, Tile{1, 2}, true},
		{Pixel{43, 54}, Tile{1, 2}, true},
		{Pixel{43, 55}, Tile{1, 2}, true},
		{Pixel{43, 56}, Tile{1, 2}, true},
		{Pixel{43, 57}, Tile{1, 2}, true},
		{Pixel{43, 58}, Tile{1, 2}, true},
		{Pixel{43, 59}, Tile{1, 2}, true},
		{Pixel{43, 60}, Tile{1, 2}, true},
		{Pixel{43, 61}, Tile{1, 2}, true},
		{Pixel{43, 62}, Tile{1, 2}, true},
		{Pixel{43, 63}, Tile{1, 2}, true},
		{Pixel{43, 64}, Tile{1, 2}, true},
		{Pixel{43, 65}, Tile{1, 2}, true},
		{Pixel{43, 66}, Tile{1, 2}, true},
		{Pixel{43, 67}, Tile{1, 2}, true},
		{Pixel{43, 68}, Tile{1, 2}, true},
		{Pixel{43, 69}, Tile{1, 2}, true},

		// right side of hexagon 1,2
		{Pixel{64, 50}, Tile{1, 2}, true},
		{Pixel{64, 51}, Tile{1, 2}, true},
		{Pixel{64, 52}, Tile{1, 2}, true},
		{Pixel{64, 53}, Tile{1, 2}, true},
		{Pixel{64, 54}, Tile{1, 2}, true},
		{Pixel{64, 55}, Tile{1, 2}, true},
		{Pixel{64, 56}, Tile{1, 2}, true},
		{Pixel{64, 57}, Tile{1, 2}, true},
		{Pixel{64, 58}, Tile{1, 2}, true},
		{Pixel{64, 59}, Tile{1, 2}, true},
		{Pixel{64, 60}, Tile{1, 2}, true},
		{Pixel{64, 61}, Tile{1, 2}, true},
		{Pixel{64, 62}, Tile{1, 2}, true},
		{Pixel{64, 63}, Tile{1, 2}, true},
		{Pixel{64, 64}, Tile{1, 2}, true},
		{Pixel{64, 65}, Tile{1, 2}, true},
		{Pixel{64, 66}, Tile{1, 2}, true},
		{Pixel{64, 67}, Tile{1, 2}, true},
		{Pixel{64, 68}, Tile{1, 2}, true},
		{Pixel{64, 69}, Tile{1, 2}, true},

		// continue right side of hexagon 1,2
		{Pixel{65, 50}, Tile{1, 2}, true},
		{Pixel{65, 51}, Tile{1, 2}, true},
		{Pixel{65, 52}, Tile{1, 2}, true},
		{Pixel{65, 53}, Tile{1, 2}, true},
		{Pixel{65, 54}, Tile{1, 2}, true},
		{Pixel{65, 55}, Tile{1, 2}, true},
		{Pixel{65, 56}, Tile{1, 2}, true},
		{Pixel{65, 57}, Tile{1, 2}, true},
		{Pixel{65, 58}, Tile{1, 2}, true},
		{Pixel{65, 59}, Tile{1, 2}, true},
		{Pixel{65, 60}, Tile{1, 2}, true},
		{Pixel{65, 61}, Tile{1, 2}, true},
		{Pixel{65, 62}, Tile{1, 2}, true},
		{Pixel{65, 63}, Tile{1, 2}, true},
		{Pixel{65, 64}, Tile{1, 2}, true},
		{Pixel{65, 65}, Tile{1, 2}, true},
		{Pixel{65, 66}, Tile{1, 2}, true},
		{Pixel{65, 67}, Tile{1, 2}, true},
		{Pixel{65, 68}, Tile{1, 2}, true},
		{Pixel{65, 69}, Tile{1, 2}, true},

		// continue right side of hexagon hexagon 1,2
		{Pixel{66, 50}, Tile{1, 2}, true},
		{Pixel{66, 51}, Tile{1, 2}, true},
		{Pixel{66, 52}, Tile{1, 2}, true},
		{Pixel{66, 53}, Tile{1, 2}, true},
		{Pixel{66, 54}, Tile{1, 2}, true},
		{Pixel{66, 55}, Tile{1, 2}, true},
		{Pixel{66, 56}, Tile{1, 2}, true},
		{Pixel{66, 57}, Tile{1, 2}, true},
		{Pixel{66, 58}, Tile{1, 2}, true},
		{Pixel{66, 59}, Tile{1, 2}, true},
		{Pixel{66, 60}, Tile{1, 2}, true},
		{Pixel{66, 61}, Tile{1, 2}, true},
		{Pixel{66, 62}, Tile{1, 2}, true},
		{Pixel{66, 63}, Tile{1, 2}, true},
		{Pixel{66, 64}, Tile{1, 2}, true},
		{Pixel{66, 65}, Tile{1, 2}, true},
		{Pixel{66, 66}, Tile{1, 2}, true},
		{Pixel{66, 67}, Tile{1, 2}, true},
		{Pixel{66, 68}, Tile{1, 2}, true},
		{Pixel{66, 69}, Tile{1, 2}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"Pixel: %s converts to hex Tile %s and want it to %s",
			tt.pixel.log(),
			tt.wantPos.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.getTile(tt.pixel)
			if ok != tt.wantOk {
				t.Errorf("Want Pixel to %s", sWantOk)
			}
			if !pos.equals(tt.wantPos) {
				t.Errorf("For Pixel %s, want Tile %s but it was %s",
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
		tile     Tile
		hexSide  hexSide
		wantTile Tile
		wantOk   bool
	}{
		// check top left hexagon {0,0}
		{Tile{0, 0}, hsNortheast, Tile{0, 0}, false},
		{Tile{0, 0}, hsEast, Tile{1, 0}, true},
		{Tile{0, 0}, hsSoutheast, Tile{1, 1}, true},
		{Tile{0, 0}, hsSouthwest, Tile{0, 1}, true},
		{Tile{0, 0}, hsWest, Tile{0, 0}, false},
		{Tile{0, 0}, hsNorthwest, Tile{0, 0}, false},

		// check top left hexagon {1, 0}
		{Tile{1, 0}, hsNortheast, Tile{0, 0}, false},
		{Tile{1, 0}, hsEast, Tile{2, 0}, true},
		{Tile{1, 0}, hsSoutheast, Tile{2, 1}, true},
		{Tile{1, 0}, hsSouthwest, Tile{1, 1}, true},
		{Tile{1, 0}, hsWest, Tile{0, 0}, true},
		{Tile{1, 0}, hsNorthwest, Tile{0, 0}, false},

		// check top left hexagon {0, 1}
		{Tile{0, 1}, hsNortheast, Tile{0, 0}, true},
		{Tile{0, 1}, hsEast, Tile{1, 1}, true},
		{Tile{0, 1}, hsSoutheast, Tile{0, 2}, true},
		{Tile{0, 1}, hsSouthwest, Tile{0, 0}, false},
		{Tile{0, 1}, hsWest, Tile{0, 0}, false},
		{Tile{0, 1}, hsNorthwest, Tile{0, 0}, false},

		// check bottom left hexagon {0,24}
		{Tile{0, 24}, hsNortheast, Tile{1, 23}, true},
		{Tile{0, 24}, hsEast, Tile{1, 24}, true},
		{Tile{0, 24}, hsSoutheast, Tile{0, 0}, false},
		{Tile{0, 24}, hsSouthwest, Tile{0, 0}, false},
		{Tile{0, 24}, hsWest, Tile{0, 0}, false},
		{Tile{0, 24}, hsNorthwest, Tile{0, 23}, true},

		// check bottom left hexagon {1,24}
		{Tile{1, 24}, hsNortheast, Tile{2, 23}, true},
		{Tile{1, 24}, hsEast, Tile{2, 24}, true},
		{Tile{1, 24}, hsSoutheast, Tile{0, 0}, false},
		{Tile{1, 24}, hsSouthwest, Tile{0, 0}, false},
		{Tile{1, 24}, hsWest, Tile{0, 24}, true},
		{Tile{1, 24}, hsNorthwest, Tile{1, 23}, true},

		// check bottom left hexagon {0,23}
		{Tile{0, 23}, hsNortheast, Tile{0, 22}, true},
		{Tile{0, 23}, hsEast, Tile{1, 23}, true},
		{Tile{0, 23}, hsSoutheast, Tile{0, 24}, true},
		{Tile{0, 23}, hsSouthwest, Tile{0, 0}, false},
		{Tile{0, 23}, hsWest, Tile{0, 0}, false},
		{Tile{0, 23}, hsNorthwest, Tile{0, 0}, false},

		// check top right hexagon {27,0}
		{Tile{27, 0}, hsNortheast, Tile{0, 0}, false},
		{Tile{27, 0}, hsEast, Tile{0, 0}, false},
		{Tile{27, 0}, hsSoutheast, Tile{28, 1}, true},
		{Tile{27, 0}, hsSouthwest, Tile{27, 1}, true},
		{Tile{27, 0}, hsWest, Tile{26, 0}, true},
		{Tile{27, 0}, hsNorthwest, Tile{0, 0}, false},

		// check top right hexagon {28,1}
		{Tile{28, 1}, hsNortheast, Tile{0, 0}, false},
		{Tile{28, 1}, hsEast, Tile{0, 0}, false},
		{Tile{28, 1}, hsSoutheast, Tile{0, 0}, false},
		{Tile{28, 1}, hsSouthwest, Tile{27, 2}, true},
		{Tile{28, 1}, hsWest, Tile{27, 1}, true},
		{Tile{28, 1}, hsNorthwest, Tile{27, 0}, true},

		// check top right hexagon {26,0}
		{Tile{26, 0}, hsNortheast, Tile{0, 0}, false},
		{Tile{26, 0}, hsEast, Tile{27, 0}, true},
		{Tile{26, 0}, hsSoutheast, Tile{27, 1}, true},
		{Tile{26, 0}, hsSouthwest, Tile{26, 1}, true},
		{Tile{26, 0}, hsWest, Tile{25, 0}, true},
		{Tile{26, 0}, hsNorthwest, Tile{0, 0}, false},

		// check bottom right hexagon {27, 24}
		{Tile{27, 24}, hsNortheast, Tile{28, 23}, true},
		{Tile{27, 24}, hsEast, Tile{0, 0}, false},
		{Tile{27, 24}, hsSoutheast, Tile{0, 0}, false},
		{Tile{27, 24}, hsSouthwest, Tile{0, 0}, false},
		{Tile{27, 24}, hsWest, Tile{26, 24}, true},
		{Tile{27, 24}, hsNorthwest, Tile{27, 23}, true},

		// check bottom right hexagon {26, 24}
		{Tile{26, 24}, hsNortheast, Tile{27, 23}, true},
		{Tile{26, 24}, hsEast, Tile{27, 24}, true},
		{Tile{26, 24}, hsSoutheast, Tile{0, 0}, false},
		{Tile{26, 24}, hsSouthwest, Tile{0, 0}, false},
		{Tile{26, 24}, hsWest, Tile{25, 24}, true},
		{Tile{26, 24}, hsNorthwest, Tile{26, 23}, true},

		// check bottom right hexagon {28, 23}
		{Tile{28, 23}, hsNortheast, Tile{0, 0}, false},
		{Tile{28, 23}, hsEast, Tile{0, 0}, false},
		{Tile{28, 23}, hsSoutheast, Tile{0, 0}, false},
		{Tile{28, 23}, hsSouthwest, Tile{27, 24}, true},
		{Tile{28, 23}, hsWest, Tile{27, 23}, true},
		{Tile{28, 23}, hsNorthwest, Tile{27, 22}, true},
	}

	for _, tt := range tests {
		var sWantOk string
		if tt.wantOk {
			sWantOk = "exist"
		} else {
			sWantOk = "not exist"
		}

		testname := fmt.Sprintf(
			"tile: %s and hex side %s has adjacent itle %s and want it to %s",
			tt.tile.log(),
			tt.hexSide.string(),
			tt.wantTile.log(),
			sWantOk)

		g := newGrid(defaultGridConfig())
		t.Run(testname, func(t *testing.T) {
			pos, ok := g.adjacent(tt.tile, tt.hexSide)
			if ok != tt.wantOk {
				t.Errorf("Want Pixel to %s", sWantOk)
			}
			if !pos.equals(tt.wantTile) {
				t.Errorf("For tile %s and hex side %s Want tile %s but it was %s",
					tt.tile.log(),
					tt.hexSide.string(),
					tt.wantTile.log(),
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
		// hexagon points near 19, 20 Pixel
		{Pixel{18, 19}, demotif{Pixel{22, 39}}, true},
		{Pixel{18, 20}, demotif{Pixel{45, 10}}, true},
		{Pixel{18, 21}, demotif{Pixel{45, 11}}, true},
		{Pixel{19, 19}, demotif{Pixel{23, 39}}, true},
		{Pixel{19, 20}, demotif{Pixel{0, 10}}, true},
		{Pixel{19, 21}, demotif{Pixel{0, 11}}, true},
		{Pixel{20, 19}, demotif{Pixel{1, 9}}, true},
		{Pixel{20, 20}, demotif{Pixel{1, 10}}, true},
		{Pixel{20, 21}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 19, 40 Pixel
		{Pixel{18, 39}, demotif{Pixel{45, 29}}, true},
		{Pixel{18, 40}, demotif{Pixel{22, 0}}, true},
		{Pixel{18, 41}, demotif{Pixel{22, 1}}, true},
		{Pixel{19, 39}, demotif{Pixel{0, 29}}, true},
		{Pixel{19, 40}, demotif{Pixel{23, 0}}, true},
		{Pixel{19, 41}, demotif{Pixel{23, 1}}, true},
		{Pixel{20, 39}, demotif{Pixel{1, 29}}, true},
		{Pixel{20, 40}, demotif{Pixel{1, 30}}, true},
		{Pixel{20, 41}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 42, 10 Pixel
		{Pixel{41, 9}, demotif{Pixel{45, 29}}, true},
		{Pixel{41, 10}, demotif{Pixel{22, 0}}, true},
		{Pixel{41, 11}, demotif{Pixel{22, 1}}, true},
		{Pixel{42, 9}, demotif{Pixel{0, 29}}, true},
		{Pixel{42, 10}, demotif{Pixel{23, 0}}, true},
		{Pixel{42, 11}, demotif{Pixel{23, 1}}, true},
		{Pixel{43, 9}, demotif{Pixel{1, 29}}, true},
		{Pixel{43, 10}, demotif{Pixel{1, 30}}, true},
		{Pixel{43, 11}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 42, 50 Pixel
		{Pixel{41, 49}, demotif{Pixel{22, 39}}, true},
		{Pixel{41, 50}, demotif{Pixel{45, 10}}, true},
		{Pixel{41, 51}, demotif{Pixel{45, 11}}, true},
		{Pixel{42, 49}, demotif{Pixel{23, 39}}, true},
		{Pixel{42, 50}, demotif{Pixel{0, 10}}, true},
		{Pixel{42, 51}, demotif{Pixel{0, 11}}, true},
		{Pixel{43, 49}, demotif{Pixel{1, 9}}, true},
		{Pixel{43, 50}, demotif{Pixel{1, 10}}, true},
		{Pixel{43, 51}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 65, 20 Pixel
		{Pixel{64, 19}, demotif{Pixel{22, 39}}, true},
		{Pixel{64, 20}, demotif{Pixel{45, 10}}, true},
		{Pixel{64, 21}, demotif{Pixel{45, 11}}, true},
		{Pixel{65, 19}, demotif{Pixel{23, 39}}, true},
		{Pixel{65, 20}, demotif{Pixel{0, 10}}, true},
		{Pixel{65, 21}, demotif{Pixel{0, 11}}, true},
		{Pixel{66, 19}, demotif{Pixel{1, 9}}, true},
		{Pixel{66, 20}, demotif{Pixel{1, 10}}, true},
		{Pixel{66, 21}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 65, 40 Pixel
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

		// first left border point hex change around 0,12 Pixel
		{Pixel{-1, 11}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 12}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 13}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 11}, demotif{Pixel{4, 31}}, true},
		{Pixel{0, 12}, demotif{Pixel{27, 2}}, true},
		{Pixel{0, 13}, demotif{Pixel{27, 3}}, true},
		{Pixel{1, 11}, demotif{Pixel{5, 31}}, true},
		{Pixel{1, 12}, demotif{Pixel{28, 2}}, true},
		{Pixel{1, 13}, demotif{Pixel{28, 3}}, true},

		// second left border point hex change around 0,48 Pixel
		{Pixel{-1, 47}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 48}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 49}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 47}, demotif{Pixel{27, 37}}, true},
		{Pixel{0, 48}, demotif{Pixel{4, 8}}, true},
		{Pixel{0, 49}, demotif{Pixel{4, 9}}, true},
		{Pixel{1, 47}, demotif{Pixel{28, 37}}, true},
		{Pixel{1, 48}, demotif{Pixel{5, 8}}, true},
		{Pixel{1, 49}, demotif{Pixel{5, 9}}, true},

		// third left border point hex change around 0,72 Pixel
		{Pixel{-1, 71}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 72}, demotif{Pixel{0, 0}}, false},
		{Pixel{-1, 73}, demotif{Pixel{0, 0}}, false},
		{Pixel{0, 71}, demotif{Pixel{4, 31}}, true},
		{Pixel{0, 72}, demotif{Pixel{27, 2}}, true},
		{Pixel{0, 73}, demotif{Pixel{27, 3}}, true},
		{Pixel{1, 71}, demotif{Pixel{5, 31}}, true},
		{Pixel{1, 72}, demotif{Pixel{28, 2}}, true},
		{Pixel{1, 73}, demotif{Pixel{28, 3}}, true},

		// first top border point hex change around 42,0 Pixel
		{Pixel{41, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{41, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{41, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{42, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{42, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{42, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{43, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{43, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{43, 1}, demotif{Pixel{1, 21}}, true},

		// second top border point hex change around 88, 0 Pixel
		{Pixel{87, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{87, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{87, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{88, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{88, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{88, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{89, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{89, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{89, 1}, demotif{Pixel{1, 21}}, true},

		// third top border point hex change around 134,0 Pixel
		{Pixel{133, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{133, 0}, demotif{Pixel{45, 20}}, true},
		{Pixel{133, 1}, demotif{Pixel{45, 21}}, true},
		{Pixel{134, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{134, 0}, demotif{Pixel{0, 20}}, true},
		{Pixel{134, 1}, demotif{Pixel{0, 21}}, true},
		{Pixel{135, -1}, demotif{Pixel{0, 0}}, false},
		{Pixel{135, 0}, demotif{Pixel{1, 20}}, true},
		{Pixel{135, 1}, demotif{Pixel{1, 21}}, true},

		// hexagon points near 19, 80 Pixel
		{Pixel{18, 79}, demotif{Pixel{22, 39}}, true},
		{Pixel{18, 80}, demotif{Pixel{45, 10}}, true},
		{Pixel{18, 81}, demotif{Pixel{45, 11}}, true},
		{Pixel{19, 79}, demotif{Pixel{23, 39}}, true},
		{Pixel{19, 80}, demotif{Pixel{0, 10}}, true},
		{Pixel{19, 81}, demotif{Pixel{0, 11}}, true},
		{Pixel{20, 79}, demotif{Pixel{1, 9}}, true},
		{Pixel{20, 80}, demotif{Pixel{1, 10}}, true},
		{Pixel{20, 81}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 42, 70 Pixel
		{Pixel{41, 69}, demotif{Pixel{45, 29}}, true},
		{Pixel{41, 70}, demotif{Pixel{22, 0}}, true},
		{Pixel{41, 71}, demotif{Pixel{22, 1}}, true},
		{Pixel{42, 69}, demotif{Pixel{0, 29}}, true},
		{Pixel{42, 70}, demotif{Pixel{23, 0}}, true},
		{Pixel{42, 71}, demotif{Pixel{23, 1}}, true},
		{Pixel{43, 69}, demotif{Pixel{1, 29}}, true},
		{Pixel{43, 70}, demotif{Pixel{1, 30}}, true},
		{Pixel{43, 71}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 65, 80 Pixel
		{Pixel{64, 79}, demotif{Pixel{22, 39}}, true},
		{Pixel{64, 80}, demotif{Pixel{45, 10}}, true},
		{Pixel{64, 81}, demotif{Pixel{45, 11}}, true},
		{Pixel{65, 79}, demotif{Pixel{23, 39}}, true},
		{Pixel{65, 80}, demotif{Pixel{0, 10}}, true},
		{Pixel{65, 81}, demotif{Pixel{0, 11}}, true},
		{Pixel{66, 79}, demotif{Pixel{1, 9}}, true},
		{Pixel{66, 80}, demotif{Pixel{1, 10}}, true},
		{Pixel{66, 81}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 88, 10 Pixel
		{Pixel{87, 9}, demotif{Pixel{45, 29}}, true},
		{Pixel{87, 10}, demotif{Pixel{22, 0}}, true},
		{Pixel{87, 11}, demotif{Pixel{22, 1}}, true},
		{Pixel{88, 9}, demotif{Pixel{0, 29}}, true},
		{Pixel{88, 10}, demotif{Pixel{23, 0}}, true},
		{Pixel{88, 11}, demotif{Pixel{23, 1}}, true},
		{Pixel{89, 9}, demotif{Pixel{1, 29}}, true},
		{Pixel{89, 10}, demotif{Pixel{1, 30}}, true},
		{Pixel{89, 11}, demotif{Pixel{24, 1}}, true},

		// hexagon points near 88, 50 Pixel
		{Pixel{87, 49}, demotif{Pixel{22, 39}}, true},
		{Pixel{87, 50}, demotif{Pixel{45, 10}}, true},
		{Pixel{87, 51}, demotif{Pixel{45, 11}}, true},
		{Pixel{88, 49}, demotif{Pixel{23, 39}}, true},
		{Pixel{88, 50}, demotif{Pixel{0, 10}}, true},
		{Pixel{88, 51}, demotif{Pixel{0, 11}}, true},
		{Pixel{89, 49}, demotif{Pixel{1, 9}}, true},
		{Pixel{89, 50}, demotif{Pixel{1, 10}}, true},
		{Pixel{89, 51}, demotif{Pixel{1, 11}}, true},

		// hexagon points near 88, 70 Pixel
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
				t.Errorf("Want Pixel to %s", sWantOk)
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
	ePixelCount, wPixelCount := g.hexPixelHeight/2, g.hexPixelHeight/2
	nwPixelCount := g.hexPixelWidth / 2
	nePixelCount := g.hexPixelWidth / 2
	swPixelCount := g.hexPixelWidth / 2
	sePixelCount := g.hexPixelWidth / 2
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

type adjacentPixelsTestCase struct {
	l          pixelLine
	d          direction
	wantPixels []Pixel
	wantPanic  bool
	trim       []Pixel
}

func TestAdjacentPixels(t *testing.T) {
	t.Parallel()
	var tests = []adjacentPixelsTestCase{
		// trim panic
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{0, 0}}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{1, 1}, {0, 0}}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{1, 1}, {1, 1}}},

		// trim success
		{pixelLine{Pixel{0, 0}, Pixel{0, 0}}, south, []Pixel{}, false, []Pixel{{0, 0}}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, false, []Pixel{{1, 1}}},
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{2, 1}}, false, []Pixel{{1, 1}}},
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{1, 1}}, false, []Pixel{{2, 1}}},
		{pixelLine{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{1, 1}, {3, 1}}, false, []Pixel{{2, 1}}},
		{pixelLine{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{2, 1}}, false, []Pixel{{1, 1}, {3, 1}}},
		{pixelLine{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{}, false, []Pixel{{1, 1}, {2, 1}, {3, 1}}},

		// Pixel 1
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, west, []Pixel{{0, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 1}}, east, []Pixel{{1, 1}}, false, []Pixel{}},
		// Pixel 2
		{pixelLine{Pixel{2, 3}, Pixel{2, 3}}, north, []Pixel{{2, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 3}, Pixel{2, 3}}, south, []Pixel{{2, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 3}, Pixel{2, 3}}, west, []Pixel{{1, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 3}, Pixel{2, 3}}, east, []Pixel{{2, 3}}, false, []Pixel{}},
		// Pixel 3
		{pixelLine{Pixel{3, 2}, Pixel{3, 2}}, north, []Pixel{{3, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{3, 2}, Pixel{3, 2}}, south, []Pixel{{3, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{3, 2}, Pixel{3, 2}}, east, []Pixel{{3, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{3, 2}, Pixel{3, 2}}, west, []Pixel{{2, 2}}, false, []Pixel{}},

		// vertical pixelLine 1
		{pixelLine{Pixel{1, 1}, Pixel{1, 2}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 2}}, south, []Pixel{{1, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 2}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 2}, Pixel{1, 1}}, south, []Pixel{{1, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 2}}, west, []Pixel{{0, 1}, {0, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 2}}, east, []Pixel{{1, 1}, {1, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 2}, Pixel{1, 1}}, west, []Pixel{{0, 1}, {0, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 2}, Pixel{1, 1}}, east, []Pixel{{1, 1}, {1, 2}}, false, []Pixel{}},
		// vertical pixelLine 2
		{pixelLine{Pixel{1, 1}, Pixel{1, 3}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 3}}, south, []Pixel{{1, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 3}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 3}, Pixel{1, 1}}, south, []Pixel{{1, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 3}}, west, []Pixel{{0, 1}, {0, 2}, {0, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{1, 3}}, east, []Pixel{{1, 1}, {1, 2}, {1, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 3}, Pixel{1, 1}}, west, []Pixel{{0, 1}, {0, 2}, {0, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 3}, Pixel{1, 1}}, east, []Pixel{{1, 1}, {1, 2}, {1, 3}}, false, []Pixel{}},
		// vertical pixelLine 3
		{pixelLine{Pixel{2, 2}, Pixel{2, 4}}, north, []Pixel{{2, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 2}, Pixel{2, 4}}, south, []Pixel{{2, 4}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 4}, Pixel{2, 2}}, north, []Pixel{{2, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 4}, Pixel{2, 2}}, south, []Pixel{{2, 4}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 2}, Pixel{2, 4}}, west, []Pixel{{1, 2}, {1, 3}, {1, 4}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 2}, Pixel{2, 4}}, east, []Pixel{{2, 2}, {2, 3}, {2, 4}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 4}, Pixel{2, 2}}, west, []Pixel{{1, 2}, {1, 3}, {1, 4}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 4}, Pixel{2, 2}}, east, []Pixel{{2, 2}, {2, 3}, {2, 4}}, false, []Pixel{}},

		// unimplemented horizontal pixelLine
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, west, []Pixel{{}}, true, []Pixel{}},
		{pixelLine{Pixel{2, 1}, Pixel{1, 1}}, east, []Pixel{{}}, true, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, west, []Pixel{{}}, true, []Pixel{}},
		{pixelLine{Pixel{2, 1}, Pixel{1, 1}}, east, []Pixel{{}}, true, []Pixel{}},
		// horizontal pixelLine 1
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, north, []Pixel{{1, 0}, {2, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}, {2, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{1, 1}, {2, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}, {2, 1}}, false, []Pixel{}},
		// horizontal pixelLine 2
		{pixelLine{Pixel{1, 1}, Pixel{3, 1}}, north, []Pixel{{1, 0}, {2, 0}, {3, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{3, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}, {2, 0}, {3, 0}}, false, []Pixel{}},
		{pixelLine{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{1, 1}, {2, 1}, {3, 1}}, false, []Pixel{}},
		{pixelLine{Pixel{3, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}, {2, 1}, {3, 1}}, false, []Pixel{}},
		// horizontal pixelLine 3
		{pixelLine{Pixel{2, 3}, Pixel{4, 3}}, north, []Pixel{{2, 2}, {3, 2}, {4, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{4, 3}, Pixel{2, 3}}, north, []Pixel{{2, 2}, {3, 2}, {4, 2}}, false, []Pixel{}},
		{pixelLine{Pixel{2, 3}, Pixel{4, 3}}, south, []Pixel{{2, 3}, {3, 3}, {4, 3}}, false, []Pixel{}},
		{pixelLine{Pixel{4, 3}, Pixel{2, 3}}, south, []Pixel{{2, 3}, {3, 3}, {4, 3}}, false, []Pixel{}},

		// unimplemented positive slope
		{pixelLine{Pixel{19, 20}, Pixel{42, 10}}, west, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{19, 20}}, east, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{19, 20}, Pixel{42, 10}}, west, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{19, 20}}, east, []Pixel{}, true, []Pixel{}},
		// unimplemented negative slope
		{pixelLine{Pixel{42, 10}, Pixel{65, 20}}, west, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{65, 20}, Pixel{42, 10}}, east, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{65, 20}}, west, []Pixel{}, true, []Pixel{}},
		{pixelLine{Pixel{65, 20}, Pixel{42, 10}}, east, []Pixel{}, true, []Pixel{}},

		// positive slope 1
		{pixelLine{Pixel{19, 20}, Pixel{42, 10}}, north,
			[]Pixel{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{19, 20}}, north,
			[]Pixel{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false, []Pixel{}},
		{pixelLine{Pixel{19, 20}, Pixel{42, 10}}, south,
			[]Pixel{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{19, 20}}, south,
			[]Pixel{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false, []Pixel{}},
		// positive slope 2
		{pixelLine{Pixel{42, 50}, Pixel{65, 40}}, north,
			[]Pixel{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false, []Pixel{}},
		{pixelLine{Pixel{65, 40}, Pixel{42, 50}}, north,
			[]Pixel{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 50}, Pixel{65, 40}}, south,
			[]Pixel{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false, []Pixel{}},
		{pixelLine{Pixel{65, 40}, Pixel{42, 50}}, south,
			[]Pixel{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false, []Pixel{}},

		// negative slope 1
		{pixelLine{Pixel{42, 10}, Pixel{65, 20}}, north,
			[]Pixel{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false, []Pixel{}},
		{pixelLine{Pixel{65, 20}, Pixel{42, 10}}, north,
			[]Pixel{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 10}, Pixel{65, 20}}, south,
			[]Pixel{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false, []Pixel{}},
		{pixelLine{Pixel{65, 20}, Pixel{42, 10}}, south,
			[]Pixel{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false, []Pixel{}},
		// negative slope 2
		{pixelLine{Pixel{19, 40}, Pixel{42, 50}}, north,
			[]Pixel{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 50}, Pixel{19, 40}}, north,
			[]Pixel{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false, []Pixel{}},
		{pixelLine{Pixel{19, 40}, Pixel{42, 50}}, south,
			[]Pixel{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false, []Pixel{}},
		{pixelLine{Pixel{42, 50}, Pixel{19, 40}}, south,
			[]Pixel{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false, []Pixel{}},

		{pixelLine{Pixel{}, Pixel{}}, noDirection, []Pixel{}, true, []Pixel{}},
	}

	for _, d := range directions() {
		switch d {
		case north, south, west, east:
			continue
		}
		test := adjacentPixelsTestCase{pixelLine{Pixel{}, Pixel{}}, d, []Pixel{}, true, []Pixel{}}
		tests = append(tests, test)
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("test adjacent Pixels %s of pixelLine %s", tt.d.string(), tt.l.log())
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
			pixels := tt.l.adjacentPixels(tt.d, tt.trim...)
			if !slices.Equal(pixels, tt.wantPixels) {
				t.Errorf("Pixels %s of pixelLine %s returned %#v. not equal to wanted Pixels \n %#v",
					tt.d.string(), tt.l.log(), pixels, tt.wantPixels)
			}
		})
	}
}
