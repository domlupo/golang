package main

import "testing"

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
