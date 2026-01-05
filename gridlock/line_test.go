package main

import (
	"fmt"
	"slices"
	"testing"
)

// sparse tests but not high priority
func TestIntersects(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		l1             line
		l2             line
		wantIntersects bool
	}{
		// does not intersect
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{3, -1}, Pixel{3, 1}}, false},
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{3, 1}, Pixel{3, -1}}, false},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{3, -1}, Pixel{3, 1}}, false},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{3, 1}, Pixel{3, -1}}, false},

		// does intersects
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{1, -1}, Pixel{1, 1}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{1, -1}, Pixel{1, 1}}, true},
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{1, 1}, Pixel{1, -1}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{1, 1}, Pixel{1, -1}}, true},

		// 3 collinear points
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{2, 0}, Pixel{2, 1}}, true},
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{2, 1}, Pixel{2, 0}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{2, 0}, Pixel{2, 1}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{2, 1}, Pixel{2, 0}}, true},

		// 4 collinear points
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{2, 0}, Pixel{4, 0}}, true},
		{line{Pixel{0, 0}, Pixel{2, 0}}, line{Pixel{4, 0}, Pixel{2, 0}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{2, 0}, Pixel{4, 0}}, true},
		{line{Pixel{2, 0}, Pixel{0, 0}}, line{Pixel{4, 0}, Pixel{2, 0}}, true},

		// lines are equal
		{line{Pixel{0, 0}, Pixel{1, 0}}, line{Pixel{0, 0}, Pixel{1, 0}}, true},
		{line{Pixel{0, 0}, Pixel{1, 0}}, line{Pixel{1, 0}, Pixel{0, 0}}, true},
		{line{Pixel{1, 0}, Pixel{0, 0}}, line{Pixel{0, 0}, Pixel{1, 0}}, true},
		{line{Pixel{1, 0}, Pixel{0, 0}}, line{Pixel{1, 0}, Pixel{0, 0}}, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf(
			"test intersects between line %s and line %s", tt.l1.log(), tt.l2.log())
		t.Run(testname, func(t *testing.T) {
			intersects := tt.l1.intersects(tt.l2)
			if intersects != tt.wantIntersects {
				t.Errorf("intersects does not equal wantIntersects")
			}
		})
	}

}

type adjacentPointsTest struct {
	l          line
	d          direction
	wantPoints []Pixel
	wantPanic  bool
	trim       []Pixel
}

func TestAdjacentPoints(t *testing.T) {
	t.Parallel()
	var tests = []adjacentPointsTest{
		// trim panic
		{line{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{0, 0}}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{1, 1}, {0, 0}}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, true, []Pixel{{1, 1}, {1, 1}}},

		// trim success
		{line{Pixel{0, 0}, Pixel{0, 0}}, south, []Pixel{}, false, []Pixel{{0, 0}}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{}, false, []Pixel{{1, 1}}},
		{line{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{2, 1}}, false, []Pixel{{1, 1}}},
		{line{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{1, 1}}, false, []Pixel{{2, 1}}},
		{line{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{1, 1}, {3, 1}}, false, []Pixel{{2, 1}}},
		{line{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{2, 1}}, false, []Pixel{{1, 1}, {3, 1}}},
		{line{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{}, false, []Pixel{{1, 1}, {2, 1}, {3, 1}}},

		// point 1
		{line{Pixel{1, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, west, []Pixel{{0, 1}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 1}}, east, []Pixel{{1, 1}}, false, []Pixel{}},
		// point 2
		{line{Pixel{2, 3}, Pixel{2, 3}}, north, []Pixel{{2, 2}}, false, []Pixel{}},
		{line{Pixel{2, 3}, Pixel{2, 3}}, south, []Pixel{{2, 3}}, false, []Pixel{}},
		{line{Pixel{2, 3}, Pixel{2, 3}}, west, []Pixel{{1, 3}}, false, []Pixel{}},
		{line{Pixel{2, 3}, Pixel{2, 3}}, east, []Pixel{{2, 3}}, false, []Pixel{}},
		// point 3
		{line{Pixel{3, 2}, Pixel{3, 2}}, north, []Pixel{{3, 1}}, false, []Pixel{}},
		{line{Pixel{3, 2}, Pixel{3, 2}}, south, []Pixel{{3, 2}}, false, []Pixel{}},
		{line{Pixel{3, 2}, Pixel{3, 2}}, east, []Pixel{{3, 2}}, false, []Pixel{}},
		{line{Pixel{3, 2}, Pixel{3, 2}}, west, []Pixel{{2, 2}}, false, []Pixel{}},

		// vertical line 1
		{line{Pixel{1, 1}, Pixel{1, 2}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 2}}, south, []Pixel{{1, 2}}, false, []Pixel{}},
		{line{Pixel{1, 2}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{line{Pixel{1, 2}, Pixel{1, 1}}, south, []Pixel{{1, 2}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 2}}, west, []Pixel{{0, 1}, {0, 2}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 2}}, east, []Pixel{{1, 1}, {1, 2}}, false, []Pixel{}},
		{line{Pixel{1, 2}, Pixel{1, 1}}, west, []Pixel{{0, 1}, {0, 2}}, false, []Pixel{}},
		{line{Pixel{1, 2}, Pixel{1, 1}}, east, []Pixel{{1, 1}, {1, 2}}, false, []Pixel{}},
		// vertical line 2
		{line{Pixel{1, 1}, Pixel{1, 3}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 3}}, south, []Pixel{{1, 3}}, false, []Pixel{}},
		{line{Pixel{1, 3}, Pixel{1, 1}}, north, []Pixel{{1, 0}}, false, []Pixel{}},
		{line{Pixel{1, 3}, Pixel{1, 1}}, south, []Pixel{{1, 3}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 3}}, west, []Pixel{{0, 1}, {0, 2}, {0, 3}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{1, 3}}, east, []Pixel{{1, 1}, {1, 2}, {1, 3}}, false, []Pixel{}},
		{line{Pixel{1, 3}, Pixel{1, 1}}, west, []Pixel{{0, 1}, {0, 2}, {0, 3}}, false, []Pixel{}},
		{line{Pixel{1, 3}, Pixel{1, 1}}, east, []Pixel{{1, 1}, {1, 2}, {1, 3}}, false, []Pixel{}},
		// vertical line 3
		{line{Pixel{2, 2}, Pixel{2, 4}}, north, []Pixel{{2, 1}}, false, []Pixel{}},
		{line{Pixel{2, 2}, Pixel{2, 4}}, south, []Pixel{{2, 4}}, false, []Pixel{}},
		{line{Pixel{2, 4}, Pixel{2, 2}}, north, []Pixel{{2, 1}}, false, []Pixel{}},
		{line{Pixel{2, 4}, Pixel{2, 2}}, south, []Pixel{{2, 4}}, false, []Pixel{}},
		{line{Pixel{2, 2}, Pixel{2, 4}}, west, []Pixel{{1, 2}, {1, 3}, {1, 4}}, false, []Pixel{}},
		{line{Pixel{2, 2}, Pixel{2, 4}}, east, []Pixel{{2, 2}, {2, 3}, {2, 4}}, false, []Pixel{}},
		{line{Pixel{2, 4}, Pixel{2, 2}}, west, []Pixel{{1, 2}, {1, 3}, {1, 4}}, false, []Pixel{}},
		{line{Pixel{2, 4}, Pixel{2, 2}}, east, []Pixel{{2, 2}, {2, 3}, {2, 4}}, false, []Pixel{}},

		// unimplemented horizontal line
		{line{Pixel{1, 1}, Pixel{2, 1}}, west, []Pixel{{}}, true, []Pixel{}},
		{line{Pixel{2, 1}, Pixel{1, 1}}, east, []Pixel{{}}, true, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{2, 1}}, west, []Pixel{{}}, true, []Pixel{}},
		{line{Pixel{2, 1}, Pixel{1, 1}}, east, []Pixel{{}}, true, []Pixel{}},
		// horizontal line 1
		{line{Pixel{1, 1}, Pixel{2, 1}}, north, []Pixel{{1, 0}, {2, 0}}, false, []Pixel{}},
		{line{Pixel{2, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}, {2, 0}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{2, 1}}, south, []Pixel{{1, 1}, {2, 1}}, false, []Pixel{}},
		{line{Pixel{2, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}, {2, 1}}, false, []Pixel{}},
		// horizontal line 2
		{line{Pixel{1, 1}, Pixel{3, 1}}, north, []Pixel{{1, 0}, {2, 0}, {3, 0}}, false, []Pixel{}},
		{line{Pixel{3, 1}, Pixel{1, 1}}, north, []Pixel{{1, 0}, {2, 0}, {3, 0}}, false, []Pixel{}},
		{line{Pixel{1, 1}, Pixel{3, 1}}, south, []Pixel{{1, 1}, {2, 1}, {3, 1}}, false, []Pixel{}},
		{line{Pixel{3, 1}, Pixel{1, 1}}, south, []Pixel{{1, 1}, {2, 1}, {3, 1}}, false, []Pixel{}},
		// horizontal line 3
		{line{Pixel{2, 3}, Pixel{4, 3}}, north, []Pixel{{2, 2}, {3, 2}, {4, 2}}, false, []Pixel{}},
		{line{Pixel{4, 3}, Pixel{2, 3}}, north, []Pixel{{2, 2}, {3, 2}, {4, 2}}, false, []Pixel{}},
		{line{Pixel{2, 3}, Pixel{4, 3}}, south, []Pixel{{2, 3}, {3, 3}, {4, 3}}, false, []Pixel{}},
		{line{Pixel{4, 3}, Pixel{2, 3}}, south, []Pixel{{2, 3}, {3, 3}, {4, 3}}, false, []Pixel{}},

		// unimplemented positive slope
		{line{Pixel{19, 20}, Pixel{42, 10}}, west, []Pixel{}, true, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{19, 20}}, east, []Pixel{}, true, []Pixel{}},
		{line{Pixel{19, 20}, Pixel{42, 10}}, west, []Pixel{}, true, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{19, 20}}, east, []Pixel{}, true, []Pixel{}},
		// unimplemented negative slope
		{line{Pixel{42, 10}, Pixel{65, 20}}, west, []Pixel{}, true, []Pixel{}},
		{line{Pixel{65, 20}, Pixel{42, 10}}, east, []Pixel{}, true, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{65, 20}}, west, []Pixel{}, true, []Pixel{}},
		{line{Pixel{65, 20}, Pixel{42, 10}}, east, []Pixel{}, true, []Pixel{}},

		// positive slope 1
		{line{Pixel{19, 20}, Pixel{42, 10}}, north,
			[]Pixel{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{19, 20}}, north,
			[]Pixel{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false, []Pixel{}},
		{line{Pixel{19, 20}, Pixel{42, 10}}, south,
			[]Pixel{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{19, 20}}, south,
			[]Pixel{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false, []Pixel{}},
		// positive slope 2
		{line{Pixel{42, 50}, Pixel{65, 40}}, north,
			[]Pixel{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false, []Pixel{}},
		{line{Pixel{65, 40}, Pixel{42, 50}}, north,
			[]Pixel{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false, []Pixel{}},
		{line{Pixel{42, 50}, Pixel{65, 40}}, south,
			[]Pixel{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false, []Pixel{}},
		{line{Pixel{65, 40}, Pixel{42, 50}}, south,
			[]Pixel{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false, []Pixel{}},

		// negative slope 1
		{line{Pixel{42, 10}, Pixel{65, 20}}, north,
			[]Pixel{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false, []Pixel{}},
		{line{Pixel{65, 20}, Pixel{42, 10}}, north,
			[]Pixel{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false, []Pixel{}},
		{line{Pixel{42, 10}, Pixel{65, 20}}, south,
			[]Pixel{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false, []Pixel{}},
		{line{Pixel{65, 20}, Pixel{42, 10}}, south,
			[]Pixel{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false, []Pixel{}},
		// negative slope 2
		{line{Pixel{19, 40}, Pixel{42, 50}}, north,
			[]Pixel{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false, []Pixel{}},
		{line{Pixel{42, 50}, Pixel{19, 40}}, north,
			[]Pixel{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false, []Pixel{}},
		{line{Pixel{19, 40}, Pixel{42, 50}}, south,
			[]Pixel{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false, []Pixel{}},
		{line{Pixel{42, 50}, Pixel{19, 40}}, south,
			[]Pixel{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false, []Pixel{}},

		{line{Pixel{}, Pixel{}}, noDirection, []Pixel{}, true, []Pixel{}},
	}

	for _, d := range directions() {
		switch d {
		case north, south, west, east:
			continue
		}
		test := adjacentPointsTest{line{Pixel{}, Pixel{}}, d, []Pixel{}, true, []Pixel{}}
		tests = append(tests, test)
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("test adjacent points %s of line %s", tt.d.string(), tt.l.log())
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
			points := tt.l.adjacentPoints(tt.d, tt.trim...)
			if !slices.Equal(points, tt.wantPoints) {
				t.Errorf("points %s of line %s returned %#v. not equal to wanted points\n %#v",
					tt.d.string(), tt.l.log(), points, tt.wantPoints)
			}
		})
	}
}

func TestNewOriginPolygon(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		lines     []line
		wantPanic bool
	}{
		// too few lines and origin centered
		{[]line{}, true},
		{[]line{
			{Pixel{0, 0}, Pixel{1, 0}},
			{Pixel{1, 0}, Pixel{0, 0}},
		}, true},

		// invalid connections and origin centered
		{[]line{
			{Pixel{-1, -1}, Pixel{1, 1}},
			{Pixel{1, 1}, Pixel{-1, -1}},
			{Pixel{1, 1}, Pixel{0, 1}},
			{Pixel{0, 1}, Pixel{1, 1}},
		}, true},

		// not origin centered and cannot be x shifted to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 1}},
			{Pixel{-1, 1}, Pixel{2, 1}},
			{Pixel{2, 1}, Pixel{2, -1}},
			{Pixel{2, -1}, Pixel{-1, -1}},
		}, true},

		// not origin centered and cannot be y shifted to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 2}},
			{Pixel{-1, 2}, Pixel{1, 2}},
			{Pixel{1, 2}, Pixel{1, -1}},
			{Pixel{1, -1}, Pixel{-1, -1}},
		}, true},

		// valid polygon when x shifted to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 1}},
			{Pixel{-1, 1}, Pixel{3, 1}},
			{Pixel{3, 1}, Pixel{3, -1}},
			{Pixel{3, -1}, Pixel{-1, -1}},
		}, false},

		// valid polygon when y shifted to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 3}},
			{Pixel{-1, 3}, Pixel{1, 3}},
			{Pixel{1, 3}, Pixel{1, -1}},
			{Pixel{1, -1}, Pixel{-1, -1}},
		}, false},

		// valid polygon when x and y shifted to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 3}},
			{Pixel{-1, 3}, Pixel{3, 3}},
			{Pixel{3, 3}, Pixel{3, -1}},
			{Pixel{3, -1}, Pixel{-1, -1}},
		}, false},

		// valid polygons without x or y shift to center origin.
		{[]line{
			{Pixel{-1, -1}, Pixel{-1, 1}},
			{Pixel{-1, 1}, Pixel{1, 1}},
			{Pixel{1, 1}, Pixel{1, -1}},
			{Pixel{1, -1}, Pixel{-1, -1}},
		}, false},
		{[]line{
			{Pixel{-2, -2}, Pixel{-2, 2}},
			{Pixel{-2, 2}, Pixel{2, 2}},
			{Pixel{2, 2}, Pixel{2, -2}},
			{Pixel{2, -2}, Pixel{-2, -2}},
		}, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf(
			"test newOriginPolygon for lines %#v", tt.lines)
		t.Run(testname, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil && tt.wantPanic:
					t.Errorf("want panic but did not panic")
				case r != nil && !tt.wantPanic:
					t.Errorf("do not want panic but did panic: %s", r)
				}
			}()
			newOriginPolygon(tt.lines)
		})
	}

}
