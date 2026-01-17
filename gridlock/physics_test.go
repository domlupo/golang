package main

import (
	"testing"
)

func TestNewOriginPolygon(t *testing.T) {
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
			newOriginPolygon(tt.lines)
		})
	}
}

func TestCollidesTwoPolygons(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name         string
		p1           polygon
		p2           polygon
		wantCollides bool
	}{

		{
			"two squares do not collide",
			newPolygon([]line{
				{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}},
				{Pos{X: 1, Y: 0}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: 0, Y: 1}},
				{Pos{X: 0, Y: 1}, Pos{X: 0, Y: 0}}}),
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			false,
		},
		{
			"two squares collide top right and bottom left",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			true,
		},
		{
			"two squares collide top left and bottom right",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 3}},
				{Pos{X: 1, Y: 3}, Pos{X: 0, Y: 3}},
				{Pos{X: 0, Y: 3}, Pos{X: 0, Y: 2}},
				{Pos{X: 0, Y: 2}, Pos{X: 1, Y: 2}}}),
			true,
		},
		{
			"two squares collide bottom right and top left",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 2, Y: 1}, Pos{X: 3, Y: 1}},
				{Pos{X: 3, Y: 1}, Pos{X: 3, Y: 0}},
				{Pos{X: 3, Y: 0}, Pos{X: 2, Y: 0}},
				{Pos{X: 2, Y: 0}, Pos{X: 2, Y: 1}}}),
			true,
		},
		{
			"two squares collide bottom left and top right",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 0, Y: 0}, Pos{X: 1, Y: 0}},
				{Pos{X: 1, Y: 0}, Pos{X: 1, Y: 1}},
				{Pos{X: 1, Y: 1}, Pos{X: 0, Y: 1}},
				{Pos{X: 0, Y: 1}, Pos{X: 0, Y: 0}}}),
			true,
		},
		{
			"two equal squares collide ",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			true,
		},
		{
			"two squares collide by subset",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 2, Y: 1}},
				{Pos{X: 2, Y: 1}, Pos{X: 2, Y: 2}},
				{Pos{X: 2, Y: 2}, Pos{X: 1, Y: 2}},
				{Pos{X: 1, Y: 2}, Pos{X: 1, Y: 1}}}),
			newPolygon([]line{
				{Pos{X: 0, Y: 0}, Pos{X: 0, Y: 3}},
				{Pos{X: 0, Y: 3}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 3, Y: 0}},
				{Pos{X: 3, Y: 0}, Pos{X: 0, Y: 0}}}),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1 := dynamic{
				// TODO: Fix this assumes a zero value Pos which is wrong.
				polygon: &tt.p1,
			}
			d2 := dynamic{
				// TODO: Fix this assumes a zero value Pos which is wrong.
				polygon: &tt.p2,
			}
			collides := d1.collides(d2)
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
		p            polygon
		d            dynamic
		wantCollides bool
	}{

		{
			"polygon within circle",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 2, Y: 4}},
				{Pos{X: 2, Y: 4}, Pos{X: 4, Y: 4}},
				{Pos{X: 4, Y: 4}, Pos{X: 4, Y: 2}},
				{Pos{X: 4, Y: 2}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 3},
				circleRadius: intPtr(5),
			},
			true,
		},
		{
			"circle within polygon",
			newPolygon([]line{
				{Pos{X: 1, Y: 1}, Pos{X: 1, Y: 5}},
				{Pos{X: 1, Y: 5}, Pos{X: 5, Y: 5}},
				{Pos{X: 5, Y: 5}, Pos{X: 5, Y: 1}},
				{Pos{X: 5, Y: 1}, Pos{X: 1, Y: 1}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 3},
				circleRadius: intPtr(1),
			},
			true,
		},
		// north east corner, collision from east
		{
			"north east corner, no collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 5, Y: 3},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"north east corner, exact collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 4, Y: 3},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"north east corner, overlap collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 4, Y: 3},
				circleRadius: intPtr(2),
			},
			true,
		},
		// north east corner, collision from north
		{
			"north east corner, no collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 5},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"north east corner, exact collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 4},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"north east corner, overlap collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 4},
				circleRadius: intPtr(2),
			},
			true,
		},
		// north west corner, collision from west
		{
			"north west corner, no collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 0, Y: 3},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"north west corner, exact collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 1, Y: 3},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"north west corner, overlap collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 1, Y: 3},
				circleRadius: intPtr(2),
			},
			true,
		},
		// north west corner, collision from north
		{
			"north west corner, no collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 5},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"north west corner, exact collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 4},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"north west corner, overlap collision from north",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 4},
				circleRadius: intPtr(2),
			},
			true,
		},
		// south west corner, collision from west
		{
			"south west corner, no collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 0, Y: 2},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"south west corner, exact collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 1, Y: 2},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"south west corner, overlap collision from west",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 1, Y: 2},
				circleRadius: intPtr(2),
			},
			true,
		},
		// south west corner, collision from south
		{
			"south west corner, no collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 0},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"south west corner, exact collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 1},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"south west corner, over lap collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 2, Y: 1},
				circleRadius: intPtr(2),
			},
			true,
		},
		// south east corner, collision from east
		{
			"south east corner, no collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 5, Y: 2},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"south east corner, exact collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 4, Y: 2},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"south east corner, overlap collision from east",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 4, Y: 2},
				circleRadius: intPtr(2),
			},
			true,
		},
		// south east corner, collision from south
		{
			"south east corner, no collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 0},
				circleRadius: intPtr(1),
			},
			false,
		},
		{
			"south east corner, exact collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 1},
				circleRadius: intPtr(1),
			},
			true,
		},
		{
			"south east corner, overlap collision from south",
			newPolygon([]line{
				{Pos{X: 2, Y: 2}, Pos{X: 3, Y: 2}},
				{Pos{X: 3, Y: 2}, Pos{X: 3, Y: 3}},
				{Pos{X: 3, Y: 3}, Pos{X: 2, Y: 3}},
				{Pos{X: 2, Y: 3}, Pos{X: 2, Y: 2}}}),
			dynamic{
				Pos:          Pos{X: 3, Y: 1},
				circleRadius: intPtr(2),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := dynamic{
				// TODO: Fix this assumes a zero value Pos which is wrong.
				polygon: &tt.p,
			}
			collides := p.collides(tt.d)
			if tt.wantCollides && !collides {
				t.Errorf("want collides but did not collide")
			}
			if !tt.wantCollides && collides {
				t.Errorf("do not want collides but did collide")
			}
		})
	}
}
