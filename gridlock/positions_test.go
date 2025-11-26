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
