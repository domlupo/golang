package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestDemotifHexLine(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)

	nw := Pixel{0, 10}
	n := Pixel{23, 0}
	ne := Pixel{46, 10}

	sw := Pixel{0, 30}
	s := Pixel{23, 40}
	se := Pixel{46, 30}

	var tests = []struct {
		hs        hexSide
		p1        Pixel
		p2        Pixel
		wantPanic bool
	}{
		{hsNone, Pixel{}, Pixel{}, true},
		{hsWest, nw, sw, false},
		{hsNorthwest, n, nw, false},
		{hsNortheast, n, ne, false},
		{hsEast, ne, se, false},
		{hsSouthwest, s, sw, false},
		{hsSoutheast, s, se, false},
	}

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
			line := hc.demotifHexLine(tt.hs)
			if !(line.p1.equals(tt.p1) || line.p1.equals(tt.p2)) {
				t.Errorf("demotifHexLine has point %s rather than expected points %s %s", line.p1.log(), tt.p1.log(), tt.p2.log())
			}
			if !(line.p2.equals(tt.p1) || line.p2.equals(tt.p2)) {
				t.Errorf("demotifHexLine has point %s rather than expected points %s %s", line.p2.log(), tt.p1.log(), tt.p2.log())
			}
		})
	}
}

func TestDemotifLinePixels(t *testing.T) {
	t.Parallel()
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)
	var tests = []struct {
		hs        hexSide
		d         side
		wantLine  []Pixel
		wantPanic bool
	}{
		{hsNone, sInclusive, []Pixel{}, true},
		{hsWest, sInclusive, []Pixel{
			{0, 10}, {0, 11}, {0, 12}, {0, 13}, {0, 14},
			{0, 15}, {0, 16}, {0, 17}, {0, 18}, {0, 19},
			{0, 20}, {0, 21}, {0, 22}, {0, 23}, {0, 24},
			{0, 25}, {0, 26}, {0, 27}, {0, 28}, {0, 29},
		}, false},
		{hsWest, sExclusive, []Pixel{
			{-1, 10}, {-1, 11}, {-1, 12}, {-1, 13}, {-1, 14},
			{-1, 15}, {-1, 16}, {-1, 17}, {-1, 18}, {-1, 19},
			{-1, 20}, {-1, 21}, {-1, 22}, {-1, 23}, {-1, 24},
			{-1, 25}, {-1, 26}, {-1, 27}, {-1, 28}, {-1, 29},
		}, false},
		{hsNorthwest, sInclusive, []Pixel{
			{0, 10}, {1, 9}, {2, 9}, {3, 8}, {4, 8},
			{5, 8}, {6, 7}, {7, 7}, {8, 6}, {9, 6},
			{10, 5}, {11, 5}, {12, 5}, {13, 4}, {14, 4},
			{15, 3}, {16, 3}, {17, 2}, {18, 2}, {19, 2},
			{20, 1}, {21, 1}, {22, 0},
		}, false},
		{hsNorthwest, sExclusive, []Pixel{
			{0, 9}, {1, 8}, {2, 8}, {3, 7}, {4, 7},
			{5, 7}, {6, 6}, {7, 6}, {8, 5}, {9, 5},
			{10, 4}, {11, 4}, {12, 4}, {13, 3}, {14, 3},
			{15, 2}, {16, 2}, {17, 1}, {18, 1}, {19, 1},
			{20, 0}, {21, 0},
		}, false},
		{hsNortheast, sInclusive, []Pixel{
			{23, 0}, {24, 1}, {25, 1}, {26, 2}, {27, 2},
			{28, 2}, {29, 3}, {30, 3}, {31, 4}, {32, 4},
			{33, 5}, {34, 5}, {35, 5}, {36, 6}, {37, 6},
			{38, 7}, {39, 7}, {40, 8}, {41, 8}, {42, 8},
			{43, 9}, {44, 9}, {45, 10},
		}, false},
		{hsNortheast, sExclusive, []Pixel{
			{24, 0}, {25, 0}, {26, 1}, {27, 1},
			{28, 1}, {29, 2}, {30, 2}, {31, 3}, {32, 3},
			{33, 4}, {34, 4}, {35, 4}, {36, 5}, {37, 5},
			{38, 6}, {39, 6}, {40, 7}, {41, 7}, {42, 7},
			{43, 8}, {44, 8}, {45, 9},
		}, false},
		{hsEast, sInclusive, []Pixel{
			{45, 10}, {45, 11}, {45, 12}, {45, 13}, {45, 14},
			{45, 15}, {45, 16}, {45, 17}, {45, 18}, {45, 19},
			{45, 20}, {45, 21}, {45, 22}, {45, 23}, {45, 24},
			{45, 25}, {45, 26}, {45, 27}, {45, 28}, {45, 29},
		}, false},
		{hsEast, sExclusive, []Pixel{
			{46, 10}, {46, 11}, {46, 12}, {46, 13}, {46, 14},
			{46, 15}, {46, 16}, {46, 17}, {46, 18}, {46, 19},
			{46, 20}, {46, 21}, {46, 22}, {46, 23}, {46, 24},
			{46, 25}, {46, 26}, {46, 27}, {46, 28}, {46, 29},
		}, false},
		{hsSouthwest, sInclusive, []Pixel{
			{0, 29}, {1, 30}, {2, 30}, {3, 31}, {4, 31},
			{5, 31}, {6, 32}, {7, 32}, {8, 33}, {9, 33},
			{10, 34}, {11, 34}, {12, 34}, {13, 35}, {14, 35},
			{15, 36}, {16, 36}, {17, 37}, {18, 37}, {19, 37},
			{20, 38}, {21, 38}, {22, 39},
		}, false},
		{hsSouthwest, sExclusive, []Pixel{
			{0, 30}, {1, 31}, {2, 31}, {3, 32}, {4, 32},
			{5, 32}, {6, 33}, {7, 33}, {8, 34}, {9, 34},
			{10, 35}, {11, 35}, {12, 35}, {13, 36}, {14, 36},
			{15, 37}, {16, 37}, {17, 38}, {18, 38}, {19, 38},
			{20, 39}, {21, 39}, {22, 40},
		}, false},
		{hsSoutheast, sInclusive, []Pixel{
			{23, 39}, {24, 38}, {25, 38}, {26, 37}, {27, 37},
			{28, 37}, {29, 36}, {30, 36}, {31, 35}, {32, 35},
			{33, 34}, {34, 34}, {35, 34}, {36, 33}, {37, 33},
			{38, 32}, {39, 32}, {40, 31}, {41, 31}, {42, 31},
			{43, 30}, {44, 30}, {45, 29},
		}, false},
		{hsSoutheast, sExclusive, []Pixel{
			{23, 40}, {24, 39}, {25, 39}, {26, 38}, {27, 38},
			{28, 38}, {29, 37}, {30, 37}, {31, 36}, {32, 36},
			{33, 35}, {34, 35}, {35, 35}, {36, 34}, {37, 34},
			{38, 33}, {39, 33}, {40, 32}, {41, 32}, {42, 32},
			{43, 31}, {44, 31}, {45, 30},
		}, false},
	}

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
			line := hc.demotifLinePixels(tt.hs, tt.d)
			if !slices.Equal(line, tt.wantLine) {
				t.Errorf("demotifLinePixel has pixels %#v rather than expected pixels \n %#v", line, tt.wantLine)
			}
		})
	}
}

func TestDemotifPixels(t *testing.T) {
	hc := newHexConfig(newWindowConfig(), 46, 40, 19, 10)
	demotifPixels := hc.demotifHexPixels()
	if len(demotifPixels) != 1380 {
		t.Error("demotifPixels sanity check failed")
	}
}
