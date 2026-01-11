package main

import (
	"fmt"
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
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 3, Y: -1}, Pos{X: 3, Y: 1}}, false},
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 3, Y: 1}, Pos{X: 3, Y: -1}}, false},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 3, Y: -1}, Pos{X: 3, Y: 1}}, false},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 3, Y: 1}, Pos{X: 3, Y: -1}}, false},

		// does intersects
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 1, Y: -1}, Pos{X: 1, Y: 1}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 1, Y: -1}, Pos{X: 1, Y: 1}}, true},
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 1, Y: 1}, Pos{X: 1, Y: -1}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 1, Y: 1}, Pos{X: 1, Y: -1}}, true},

		// 3 collinear points
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 2, Y: 0}, Pos{X: 2, Y: 1}}, true},
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 0}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 2, Y: 0}, Pos{X: 2, Y: 1}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 0}}, true},

		// 4 collinear points
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 2, Y: 0}, Pos{X: 4, Y: 0}}, true},
		{line{Pos{X: 0, Y: 0}, Pos{X: 2, Y: 0}}, line{Pos{X: 4, Y: 0}, Pos{X: 2, Y: 0}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 2, Y: 0}, Pos{X: 4, Y: 0}}, true},
		{line{Pos{X: 2, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 4, Y: 0}, Pos{X: 2, Y: 0}}, true},

		// lines are equal
		{line{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}}, line{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}}, true},
		{line{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}}, line{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 0}}, true},
		{line{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}}, true},
		{line{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 0}}, line{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 0}}, true},
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
