package main

import (
	"fmt"
	"testing"
)

// low level and sparse tests but not high priority to fix.
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

func TestNewStaticPolygon(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name      string
		lines     []line
		wantPanic bool
	}{

		{
			"zero lines and origin centered",
			[]line{},
			true,
		},
		{
			"two lines and origin centered",
			[]line{
				{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}},
				{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 0}},
			},
			true,
		},
		{
			"invalid connections and origin centered",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: -1, Y: -1}},
				{Pos{X: 1, Y: 0}, Pos{X: 0, Y: 1}},
				{Pos{X: 0, Y: 1}, Pos{X: 1, Y: 0}},
			},
			true,
		},
		{
			"#1 invalid connections and not origin centered",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -1, Y: 2}},
				{Pos{X: -1, Y: 1}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: -1, Y: 1}},
				{Pos{X: -1, Y: 1}, Pos{X: -1, Y: -1}},
			},
			true,
		},
		{
			"#2 invalid connections and not origin centered",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -2, Y: 1}},
				{Pos{X: -1, Y: 1}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: -1, Y: 1}},
				{Pos{X: -1, Y: 1}, Pos{X: -1, Y: -1}},
			},
			true,
		},
		{
			"valid polygon when x shifted to center origin",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -1, Y: 1}},
				{Pos{X: -1, Y: 1}, Pos{X: 3, Y: 1}},
				{Pos{X: 3, Y: 1}, Pos{X: 3, Y: -1}},
				{Pos{X: 3, Y: -1}, Pos{X: -1, Y: -1}},
			},
			false,
		},
		{
			"valid polygon when y shifted to center origin",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -1, Y: 3}},
				{Pos{X: -1, Y: 3}, Pos{X: 1, Y: 3}},
				{Pos{X: 1, Y: 3}, Pos{X: 1, Y: -1}},
				{Pos{X: 1, Y: -1}, Pos{X: -1, Y: -1}},
			},
			false,
		},
		{
			"valid polygon when x and y shifted to center origin",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -1, Y: 3}},
				{Pos{X: -1, Y: 3}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 3, Y: -1}},
				{Pos{X: 3, Y: -1}, Pos{X: -1, Y: -1}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pos{X: -1, Y: -1}, Pos{X: -1, Y: 1}},
				{Pos{X: -1, Y: 1}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: 1, Y: -1}},
				{Pos{X: 1, Y: -1}, Pos{X: -1, Y: -1}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pos{X: -2, Y: -2}, Pos{X: -2, Y: 2}},
				{Pos{X: -2, Y: 2}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 2, Y: -2}},
				{Pos{X: 2, Y: -2}, Pos{X: -2, Y: -2}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pos{X: -2, Y: -2}, Pos{X: -2, Y: 2}},
				{Pos{X: -2, Y: 2}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 2, Y: -2}},
				{Pos{X: 2, Y: -2}, Pos{X: -2, Y: -2}},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil && tt.wantPanic:
					t.Errorf("want panic but did not panic")
				case r != nil && !tt.wantPanic:
					t.Errorf("do not want panic but did panic: %s", r)
				}
			}()
			newStaticPolygon(tt.lines, Pos{})
		})
	}
}

func TestCollidesTwoPolygons(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name          string
		polygonLines1 []line
		polygonLines2 []line
		wantCollides  bool
	}{

		{
			"two squares do not collide",
			[]line{
				{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}},
				{Pos{X: 1, Y: 0}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: 0, Y: 1}},
				{Pos{X: 0, Y: 1}, Pos{X: 0, Y: 0}}},
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			false,
		},
		{
			"two squares collide top right and bottom left",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			true,
		},
		{
			"two squares collide top left and bottom right",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 3}},
				{Pos{X: 1, Y: 3}, Pos{X: 0, Y: 3}},
				{Pos{X: 0, Y: 3}, Pos{X: 0, Y: 2}},
				{Pos{X: 0, Y: 2}, Pos{X: 1, Y: 2}}},
			true,
		},
		{
			"two squares collide bottom right and top left",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 2, Y: 1}, Pos{X: 3, Y: 1}},
				{Pos{X: 3, Y: 1}, Pos{X: 3, Y: 0}},
				{Pos{X: 3, Y: 0}, Pos{X: 2, Y: 0}},
				{Pos{X: 2, Y: 0}, Pos{X: 2, Y: 1}}},
			true,
		},
		{
			"two squares collide bottom left and top right",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}},
				{Pos{X: 1, Y: 0}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: 0, Y: 1}},
				{Pos{X: 0, Y: 1}, Pos{X: 0, Y: 0}}},
			true,
		},
		{
			"two equal squares collide ",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			true,
		},
		{
			"two squares collide by subset",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}},
			[]line{
				{Pos{X: 0, Y: 0}, Pos{X: 0, Y: 3}},
				{Pos{X: 0, Y: 3}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 3, Y: 0}},
				{Pos{X: 3, Y: 0}, Pos{X: 0, Y: 0}}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := newDynamicPolygon(tt.polygonLines1)
			p2 := newDynamicPolygon(tt.polygonLines2)
			collides := p1.collides(*p2)
			otherCollides := p2.collides(*p1)
			if collides != otherCollides {
				t.Errorf("the struct used as the receiver should not affect collision")
			}
			if tt.wantCollides && !collides {
				t.Errorf("want collides but did not collide")
			}
			if !tt.wantCollides && collides {
				t.Errorf("do not want collides but did collide")
			}
		})
	}
}

func TestCollidesOnePolygonOneCircle(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name         string
		polygonLines []line
		circlePos    Pos
		radius       int
		wantCollides bool
	}{

		{
			"polygon within circle",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 2, Y: 4}},
				{Pos{X: 2, Y: 4}, Pos{X: 4, Y: 4}},
				{Pos{X: 4, Y: 4}, Pos{X: 4, Y: 2}},
				{Pos{X: 4, Y: 2}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 3},
			5,
			true,
		},
		{
			"circle within polygon",
			[]line{
				{Pos{X: 1, Y: 1}, Pos{X: 1, Y: 5}},
				{Pos{X: 1, Y: 5}, Pos{X: 5, Y: 5}},
				{Pos{X: 5, Y: 5}, Pos{X: 5, Y: 1}},
				{Pos{X: 5, Y: 1}, Pos{X: 1, Y: 1}}},
			Pos{X: 3, Y: 3},
			1,
			true,
		},
		// north east corner, collision from east
		{
			"north east corner, no collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 5, Y: 3},
			1,
			false,
		},
		{
			"north east corner, exact collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 4, Y: 3},
			1,
			true,
		},
		{
			"north east corner, overlap collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 4, Y: 3},
			2,
			true,
		},
		// north east corner, collision from north
		{
			"north east corner, no collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 5},
			1,
			false,
		},
		{
			"north east corner, exact collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 4},
			1,
			true,
		},
		{
			"north east corner, overlap collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 4},
			2,
			true,
		},
		// north west corner, collision from west
		{
			"north west corner, no collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 0, Y: 3},
			1,
			false,
		},
		{"north west corner, exact collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 1, Y: 3},
			1,
			true,
		},
		{
			"north west corner, overlap collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 1, Y: 3},
			2,
			true,
		},
		// north west corner, collision from north
		{
			"north west corner, no collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 5},
			1,
			false,
		},
		{
			"north west corner, exact collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 4},
			1,
			true,
		},
		{
			"north west corner, overlap collision from north",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 4},
			2,
			true,
		},
		// south west corner, collision from west
		{
			"south west corner, no collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 0, Y: 2},
			1,
			false,
		},
		{
			"south west corner, exact collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 1, Y: 2},
			1,
			true,
		},
		{
			"south west corner, overlap collision from west",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 1, Y: 2},
			2,
			true,
		},
		// south west corner, collision from south
		{
			"south west corner, no collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 0},
			1,
			false,
		},
		{
			"south west corner, exact collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 1},
			1,
			true,
		},
		{
			"south west corner, over lap collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 2, Y: 1},
			2,
			true,
		},
		// south east corner, collision from east
		{
			"south east corner, no collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 5, Y: 2},
			1,
			false,
		},
		{
			"south east corner, exact collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 4, Y: 2},
			1,
			true,
		},
		{
			"south east corner, overlap collision from east",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 4, Y: 2},
			2,
			true,
		},
		// south east corner, collision from south
		{
			"south east corner, no collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 0},
			1,
			false,
		},
		{
			"south east corner, exact collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 1},
			1,
			true,
		},
		{
			"south east corner, overlap collision from south",
			[]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}},
			Pos{X: 3, Y: 1},
			2,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newDynamicPolygon(tt.polygonLines)
			c := newDynamicCircle(tt.circlePos, tt.radius)
			collides := p.collides(*c)
			otherCollides := c.collides(*p)
			if collides != otherCollides {
				t.Errorf("the struct used as the receiver should not affect collision")
			}
			if tt.wantCollides && !collides {
				t.Errorf("want collides but did not collide")
			}
			if !tt.wantCollides && collides {
				t.Errorf("do not want collides but did collide")
			}
		})
	}
}
