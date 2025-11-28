package main

import (
	"fmt"
	"slices"
	"testing"
)

// TODO: Sparse tests but not a high priority.
func TestIntersects(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		l1             line
		l2             line
		wantIntersects bool
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
		testname := fmt.Sprintf("test intersects between line %s and line %s", tt.l1.log(), tt.l2.log())
		t.Run(testname, func(t *testing.T) {
			intersects := tt.l1.intersects(tt.l2)
			if intersects != tt.wantIntersects {
				t.Errorf("interects does not equal wantIntersects")
			}
		})
	}

}

func TestInclusivePoints(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		l          line
		wantPoints []PixelPosition
		wantPanic  bool
	}{
		// horizontal lines
		{line{PixelPosition{0, 0}, PixelPosition{1, 0}}, []PixelPosition{{0, 0}, {1, 0}}, false},
		{line{PixelPosition{0, 0}, PixelPosition{2, 0}}, []PixelPosition{{0, 0}, {1, 0}, {2, 0}}, false},
		{line{PixelPosition{1, 0}, PixelPosition{0, 0}}, []PixelPosition{{0, 0}, {1, 0}}, false},
		{line{PixelPosition{2, 0}, PixelPosition{0, 0}}, []PixelPosition{{0, 0}, {1, 0}, {2, 0}}, false},

		// vertical lines
		{line{PixelPosition{0, 0}, PixelPosition{0, 1}}, []PixelPosition{{0, 0}, {0, 1}}, false},
		{line{PixelPosition{0, 0}, PixelPosition{0, 2}}, []PixelPosition{{0, 0}, {0, 1}, {0, 2}}, false},
		{line{PixelPosition{0, 1}, PixelPosition{0, 0}}, []PixelPosition{{0, 0}, {0, 1}}, false},
		{line{PixelPosition{0, 2}, PixelPosition{0, 0}}, []PixelPosition{{0, 0}, {0, 1}, {0, 2}}, false},

		// lines not either horizontal or vertical
		{line{PixelPosition{0, 0}, PixelPosition{0, 0}}, []PixelPosition{}, true},
		{line{PixelPosition{0, 0}, PixelPosition{1, 1}}, []PixelPosition{}, true},
		{line{PixelPosition{1, 1}, PixelPosition{1, 1}}, []PixelPosition{}, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("test inclusive points for line %s", tt.l.log())
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

			inclusivePoints := tt.l.inclusivePoints()
			if !slices.Equal(inclusivePoints, tt.wantPoints) {
				t.Errorf("inclusive points %#v not equal to wanted points %#v", inclusivePoints, tt.wantPoints)
			}
		})
	}
}

func TestAdjacentPointsNoBounds(t *testing.T) {
	l := line{PixelPosition{-1, -1}, PixelPosition{-1, -1}}
	point := l.adjacentPoints(south, nil)[0]
	if !point.equals(PixelPosition{-1, -1}) {
		t.Error("A point should be returned when bounds is nil")
	}
}

func TestAdjacentPoints(t *testing.T) {
	t.Parallel()
	wc := newWindowConfig()
	var tests = []struct {
		l          line
		d          direction
		wantPoints []PixelPosition
		wantPanic  bool
	}{
		// trims points out of bounds 1
		{line{PixelPosition{0, 0}, PixelPosition{0, 0}}, north, []PixelPosition{}, false},
		{line{PixelPosition{0, 0}, PixelPosition{0, 0}}, south, []PixelPosition{{0, 0}}, false},
		{line{PixelPosition{0, -1}, PixelPosition{0, -1}}, north, []PixelPosition{}, false},
		{line{PixelPosition{0, -1}, PixelPosition{0, -1}}, south, []PixelPosition{}, false},
		{line{PixelPosition{0, 0}, PixelPosition{0, 0}}, west, []PixelPosition{}, false},
		{line{PixelPosition{0, 0}, PixelPosition{0, 0}}, east, []PixelPosition{{0, 0}}, false},
		{line{PixelPosition{-1, 0}, PixelPosition{-1, 0}}, west, []PixelPosition{}, false},
		{line{PixelPosition{-1, 0}, PixelPosition{-1, 0}}, east, []PixelPosition{}, false},
		//trims points out of bounds 2
		{line{PixelPosition{-1, 1}, PixelPosition{-1, 1}}, north, []PixelPosition{}, false},
		{line{PixelPosition{-1, 1}, PixelPosition{-1, 1}}, south, []PixelPosition{}, false},
		{line{PixelPosition{-1, 1}, PixelPosition{-1, 1}}, east, []PixelPosition{}, false},
		{line{PixelPosition{-1, 1}, PixelPosition{-1, 1}}, west, []PixelPosition{}, false},
		// trims point out of bounds 3
		{line{PixelPosition{wc.width, 1}, PixelPosition{wc.width, 1}}, north, []PixelPosition{}, false},
		{line{PixelPosition{wc.width, 1}, PixelPosition{wc.width, 1}}, south, []PixelPosition{}, false},
		{line{PixelPosition{wc.width, 1}, PixelPosition{wc.width, 1}}, west, []PixelPosition{{wc.width - 1, 1}}, false},
		{line{PixelPosition{wc.width, 1}, PixelPosition{wc.width, 1}}, east, []PixelPosition{}, false},
		// trims points out of bounds 4
		{line{PixelPosition{1, wc.height}, PixelPosition{1, wc.height}}, north, []PixelPosition{{1, wc.height - 1}}, false},
		{line{PixelPosition{1, wc.height}, PixelPosition{1, wc.height}}, south, []PixelPosition{}, false},
		{line{PixelPosition{1, wc.height}, PixelPosition{1, wc.height}}, west, []PixelPosition{}, false},
		{line{PixelPosition{1, wc.height}, PixelPosition{1, wc.height}}, east, []PixelPosition{}, false},
		{line{PixelPosition{1, wc.height - 1}, PixelPosition{1, wc.height - 1}}, north, []PixelPosition{{1, wc.height - 2}}, false},
		{line{PixelPosition{1, wc.height - 1}, PixelPosition{1, wc.height - 1}}, south, []PixelPosition{{1, wc.height - 1}}, false},
		{line{PixelPosition{1, wc.height - 1}, PixelPosition{1, wc.height - 1}}, west, []PixelPosition{{0, wc.height - 1}}, false},
		{line{PixelPosition{1, wc.height - 1}, PixelPosition{1, wc.height - 1}}, east, []PixelPosition{{1, wc.height - 1}}, false},
		// trims points out of bounds 5
		{line{PixelPosition{0, -3}, PixelPosition{0, 0}}, north, []PixelPosition{}, false},
		{line{PixelPosition{0, -3}, PixelPosition{0, 0}}, south, []PixelPosition{{0, 0}}, false},
		{line{PixelPosition{0, -3}, PixelPosition{0, -1}}, north, []PixelPosition{}, false},
		{line{PixelPosition{0, -3}, PixelPosition{0, -1}}, south, []PixelPosition{}, false},

		// point 1
		{line{PixelPosition{1, 1}, PixelPosition{1, 1}}, north, []PixelPosition{{1, 0}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 1}}, south, []PixelPosition{{1, 1}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 1}}, west, []PixelPosition{{0, 1}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 1}}, east, []PixelPosition{{1, 1}}, false},
		// point 2
		{line{PixelPosition{2, 3}, PixelPosition{2, 3}}, north, []PixelPosition{{2, 2}}, false},
		{line{PixelPosition{2, 3}, PixelPosition{2, 3}}, south, []PixelPosition{{2, 3}}, false},
		{line{PixelPosition{2, 3}, PixelPosition{2, 3}}, west, []PixelPosition{{1, 3}}, false},
		{line{PixelPosition{2, 3}, PixelPosition{2, 3}}, east, []PixelPosition{{2, 3}}, false},
		// point 3
		{line{PixelPosition{3, 2}, PixelPosition{3, 2}}, north, []PixelPosition{{3, 1}}, false},
		{line{PixelPosition{3, 2}, PixelPosition{3, 2}}, south, []PixelPosition{{3, 2}}, false},
		{line{PixelPosition{3, 2}, PixelPosition{3, 2}}, east, []PixelPosition{{3, 2}}, false},
		{line{PixelPosition{3, 2}, PixelPosition{3, 2}}, west, []PixelPosition{{2, 2}}, false},

		// vertical line 1
		{line{PixelPosition{1, 1}, PixelPosition{1, 2}}, north, []PixelPosition{{1, 0}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 2}}, south, []PixelPosition{{1, 2}}, false},
		{line{PixelPosition{1, 2}, PixelPosition{1, 1}}, north, []PixelPosition{{1, 0}}, false},
		{line{PixelPosition{1, 2}, PixelPosition{1, 1}}, south, []PixelPosition{{1, 2}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 2}}, west, []PixelPosition{{0, 1}, {0, 2}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 2}}, east, []PixelPosition{{1, 1}, {1, 2}}, false},
		{line{PixelPosition{1, 2}, PixelPosition{1, 1}}, west, []PixelPosition{{0, 1}, {0, 2}}, false},
		{line{PixelPosition{1, 2}, PixelPosition{1, 1}}, east, []PixelPosition{{1, 1}, {1, 2}}, false},
		// vertical line 2
		{line{PixelPosition{1, 1}, PixelPosition{1, 3}}, north, []PixelPosition{{1, 0}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 3}}, south, []PixelPosition{{1, 3}}, false},
		{line{PixelPosition{1, 3}, PixelPosition{1, 1}}, north, []PixelPosition{{1, 0}}, false},
		{line{PixelPosition{1, 3}, PixelPosition{1, 1}}, south, []PixelPosition{{1, 3}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 3}}, west, []PixelPosition{{0, 1}, {0, 2}, {0, 3}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{1, 3}}, east, []PixelPosition{{1, 1}, {1, 2}, {1, 3}}, false},
		{line{PixelPosition{1, 3}, PixelPosition{1, 1}}, west, []PixelPosition{{0, 1}, {0, 2}, {0, 3}}, false},
		{line{PixelPosition{1, 3}, PixelPosition{1, 1}}, east, []PixelPosition{{1, 1}, {1, 2}, {1, 3}}, false},
		// vertical line 3
		{line{PixelPosition{2, 2}, PixelPosition{2, 4}}, north, []PixelPosition{{2, 1}}, false},
		{line{PixelPosition{2, 2}, PixelPosition{2, 4}}, south, []PixelPosition{{2, 4}}, false},
		{line{PixelPosition{2, 4}, PixelPosition{2, 2}}, north, []PixelPosition{{2, 1}}, false},
		{line{PixelPosition{2, 4}, PixelPosition{2, 2}}, south, []PixelPosition{{2, 4}}, false},
		{line{PixelPosition{2, 2}, PixelPosition{2, 4}}, west, []PixelPosition{{1, 2}, {1, 3}, {1, 4}}, false},
		{line{PixelPosition{2, 2}, PixelPosition{2, 4}}, east, []PixelPosition{{2, 2}, {2, 3}, {2, 4}}, false},
		{line{PixelPosition{2, 4}, PixelPosition{2, 2}}, west, []PixelPosition{{1, 2}, {1, 3}, {1, 4}}, false},
		{line{PixelPosition{2, 4}, PixelPosition{2, 2}}, east, []PixelPosition{{2, 2}, {2, 3}, {2, 4}}, false},

		// unimplemented horizontal line
		{line{PixelPosition{1, 1}, PixelPosition{2, 1}}, west, []PixelPosition{{}}, true},
		{line{PixelPosition{2, 1}, PixelPosition{1, 1}}, east, []PixelPosition{{}}, true},
		{line{PixelPosition{1, 1}, PixelPosition{2, 1}}, west, []PixelPosition{{}}, true},
		{line{PixelPosition{2, 1}, PixelPosition{1, 1}}, east, []PixelPosition{{}}, true},
		// horizontal line 1
		{line{PixelPosition{1, 1}, PixelPosition{2, 1}}, north, []PixelPosition{{1, 0}, {2, 0}}, false},
		{line{PixelPosition{2, 1}, PixelPosition{1, 1}}, north, []PixelPosition{{1, 0}, {2, 0}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{2, 1}}, south, []PixelPosition{{1, 1}, {2, 1}}, false},
		{line{PixelPosition{2, 1}, PixelPosition{1, 1}}, south, []PixelPosition{{1, 1}, {2, 1}}, false},
		// horizontal line 2
		{line{PixelPosition{1, 1}, PixelPosition{3, 1}}, north, []PixelPosition{{1, 0}, {2, 0}, {3, 0}}, false},
		{line{PixelPosition{3, 1}, PixelPosition{1, 1}}, north, []PixelPosition{{1, 0}, {2, 0}, {3, 0}}, false},
		{line{PixelPosition{1, 1}, PixelPosition{3, 1}}, south, []PixelPosition{{1, 1}, {2, 1}, {3, 1}}, false},
		{line{PixelPosition{3, 1}, PixelPosition{1, 1}}, south, []PixelPosition{{1, 1}, {2, 1}, {3, 1}}, false},
		// horizontal line 3
		{line{PixelPosition{2, 3}, PixelPosition{4, 3}}, north, []PixelPosition{{2, 2}, {3, 2}, {4, 2}}, false},
		{line{PixelPosition{4, 3}, PixelPosition{2, 3}}, north, []PixelPosition{{2, 2}, {3, 2}, {4, 2}}, false},
		{line{PixelPosition{2, 3}, PixelPosition{4, 3}}, south, []PixelPosition{{2, 3}, {3, 3}, {4, 3}}, false},
		{line{PixelPosition{4, 3}, PixelPosition{2, 3}}, south, []PixelPosition{{2, 3}, {3, 3}, {4, 3}}, false},

		// unimplemented postitive slop
		{line{PixelPosition{19, 20}, PixelPosition{42, 10}}, west, []PixelPosition{}, true},
		{line{PixelPosition{42, 10}, PixelPosition{19, 20}}, east, []PixelPosition{}, true},
		{line{PixelPosition{19, 20}, PixelPosition{42, 10}}, west, []PixelPosition{}, true},
		{line{PixelPosition{42, 10}, PixelPosition{19, 20}}, east, []PixelPosition{}, true},
		// unimplemented negative slop
		{line{PixelPosition{42, 10}, PixelPosition{65, 20}}, west, []PixelPosition{}, true},
		{line{PixelPosition{65, 20}, PixelPosition{42, 10}}, east, []PixelPosition{}, true},
		{line{PixelPosition{42, 10}, PixelPosition{65, 20}}, west, []PixelPosition{}, true},
		{line{PixelPosition{65, 20}, PixelPosition{42, 10}}, east, []PixelPosition{}, true},

		// positive slope 1
		{line{PixelPosition{19, 20}, PixelPosition{42, 10}}, north,
			[]PixelPosition{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false},
		{line{PixelPosition{42, 10}, PixelPosition{19, 20}}, north,
			[]PixelPosition{
				{19, 19}, {20, 18}, {21, 18}, {22, 17}, {23, 17},
				{24, 17}, {25, 16}, {26, 16}, {27, 15}, {28, 15},
				{29, 14}, {30, 14}, {31, 14}, {32, 13}, {33, 13},
				{34, 12}, {35, 12}, {36, 11}, {37, 11}, {38, 11},
				{39, 10}, {40, 10}, {41, 9}, {42, 9}},
			false},
		{line{PixelPosition{19, 20}, PixelPosition{42, 10}}, south,
			[]PixelPosition{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false},
		{line{PixelPosition{42, 10}, PixelPosition{19, 20}}, south,
			[]PixelPosition{
				{19, 20}, {20, 19}, {21, 19}, {22, 18}, {23, 18},
				{24, 18}, {25, 17}, {26, 17}, {27, 16}, {28, 16},
				{29, 15}, {30, 15}, {31, 15}, {32, 14}, {33, 14},
				{34, 13}, {35, 13}, {36, 12}, {37, 12}, {38, 12},
				{39, 11}, {40, 11}, {41, 10}, {42, 10}},
			false},
		// positive slope 2
		{line{PixelPosition{42, 50}, PixelPosition{65, 40}}, north,
			[]PixelPosition{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false},
		{line{PixelPosition{65, 40}, PixelPosition{42, 50}}, north,
			[]PixelPosition{
				{42, 49}, {43, 48}, {44, 48}, {45, 47}, {46, 47},
				{47, 47}, {48, 46}, {49, 46}, {50, 45}, {51, 45},
				{52, 44}, {53, 44}, {54, 44}, {55, 43}, {56, 43},
				{57, 42}, {58, 42}, {59, 41}, {60, 41}, {61, 41},
				{62, 40}, {63, 40}, {64, 39}, {65, 39}},
			false},
		{line{PixelPosition{42, 50}, PixelPosition{65, 40}}, south,
			[]PixelPosition{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false},
		{line{PixelPosition{65, 40}, PixelPosition{42, 50}}, south,
			[]PixelPosition{
				{42, 50}, {43, 49}, {44, 49}, {45, 48}, {46, 48},
				{47, 48}, {48, 47}, {49, 47}, {50, 46}, {51, 46},
				{52, 45}, {53, 45}, {54, 45}, {55, 44}, {56, 44},
				{57, 43}, {58, 43}, {59, 42}, {60, 42}, {61, 42},
				{62, 41}, {63, 41}, {64, 40}, {65, 40}},
			false},

		// negative slope 1
		{line{PixelPosition{42, 10}, PixelPosition{65, 20}}, north,
			[]PixelPosition{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false},
		{line{PixelPosition{65, 20}, PixelPosition{42, 10}}, north,
			[]PixelPosition{
				{42, 9}, {43, 10}, {44, 10}, {45, 11}, {46, 11},
				{47, 11}, {48, 12}, {49, 12}, {50, 13}, {51, 13},
				{52, 14}, {53, 14}, {54, 14}, {55, 15}, {56, 15},
				{57, 16}, {58, 16}, {59, 17}, {60, 17}, {61, 17},
				{62, 18}, {63, 18}, {64, 19}, {65, 19}},
			false},
		{line{PixelPosition{42, 10}, PixelPosition{65, 20}}, south,
			[]PixelPosition{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false},
		{line{PixelPosition{65, 20}, PixelPosition{42, 10}}, south,
			[]PixelPosition{
				{42, 10}, {43, 11}, {44, 11}, {45, 12}, {46, 12},
				{47, 12}, {48, 13}, {49, 13}, {50, 14}, {51, 14},
				{52, 15}, {53, 15}, {54, 15}, {55, 16}, {56, 16},
				{57, 17}, {58, 17}, {59, 18}, {60, 18}, {61, 18},
				{62, 19}, {63, 19}, {64, 20}, {65, 20}},
			false},
		// negative slope 2
		{line{PixelPosition{19, 40}, PixelPosition{42, 50}}, north,
			[]PixelPosition{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false},
		{line{PixelPosition{42, 50}, PixelPosition{19, 40}}, north,
			[]PixelPosition{
				{19, 39}, {20, 40}, {21, 40}, {22, 41}, {23, 41},
				{24, 41}, {25, 42}, {26, 42}, {27, 43}, {28, 43},
				{29, 44}, {30, 44}, {31, 44}, {32, 45}, {33, 45},
				{34, 46}, {35, 46}, {36, 47}, {37, 47}, {38, 47},
				{39, 48}, {40, 48}, {41, 49}, {42, 49}},
			false},
		{line{PixelPosition{19, 40}, PixelPosition{42, 50}}, south,
			[]PixelPosition{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false},
		{line{PixelPosition{42, 50}, PixelPosition{19, 40}}, south,
			[]PixelPosition{
				{19, 40}, {20, 41}, {21, 41}, {22, 42}, {23, 42},
				{24, 42}, {25, 43}, {26, 43}, {27, 44}, {28, 44},
				{29, 45}, {30, 45}, {31, 45}, {32, 46}, {33, 46},
				{34, 47}, {35, 47}, {36, 48}, {37, 48}, {38, 48},
				{39, 49}, {40, 49}, {41, 50}, {42, 50}},
			false},

		{line{PixelPosition{}, PixelPosition{}}, noDirecton, []PixelPosition{}, true},
	}
	for _, d := range directions() {
		switch d {
		case north, south, west, east:
			continue
		}
		test := struct {
			l          line
			d          direction
			wantPoints []PixelPosition
			wantPanic  bool
		}{
			line{PixelPosition{}, PixelPosition{}}, d, []PixelPosition{}, true}
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
			wc := newWindowConfig()
			points := tt.l.adjacentPoints(tt.d, &wc.bounds)
			if !slices.Equal(points, tt.wantPoints) {
				t.Errorf("points %s of line %s returned %#v. not equal to wanted points\n %#v", tt.d.string(), tt.l.log(), points, tt.wantPoints)
			}
		})
	}
}
