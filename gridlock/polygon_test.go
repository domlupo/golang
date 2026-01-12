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

func TestCollides(t *testing.T) {
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
			collides := tt.p1.collides(tt.p2)
			if tt.wantCollides && !collides {
				t.Errorf("want collides but did not collide")
			}
			if !tt.wantCollides && collides {
				t.Errorf("do not want collides but did collide")
			}
		})
	}
}
