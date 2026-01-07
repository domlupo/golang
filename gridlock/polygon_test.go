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
				{Pixel{0, 0}, Pixel{1, 0}},
				{Pixel{1, 0}, Pixel{0, 0}},
			},
			true,
		},
		{
			"invalid connections and origin centered",
			[]line{
				{Pixel{-1, -1}, Pixel{1, 1}},
				{Pixel{1, 1}, Pixel{-1, -1}},
				{Pixel{1, 0}, Pixel{0, 1}},
				{Pixel{0, 1}, Pixel{1, 0}},
			},
			true,
		},
		{
			"#1 invalid connections and not origin centered",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 2}},
				{Pixel{-1, 1}, Pixel{1, 1}},
				{Pixel{1, 1}, Pixel{-1, 1}},
				{Pixel{-1, 1}, Pixel{-1, -1}},
			},
			true,
		},
		{
			"#2 invalid connections and not origin centered",
			[]line{
				{Pixel{-1, -1}, Pixel{-2, 1}},
				{Pixel{-1, 1}, Pixel{1, 1}},
				{Pixel{1, 1}, Pixel{-1, 1}},
				{Pixel{-1, 1}, Pixel{-1, -1}},
			},
			true,
		},
		{
			"not origin centered and cannot be x shifted to center origin.",
			[]line{
				{Pixel{-1, -1}, Pixel{1, 1}},
				{Pixel{1, 1}, Pixel{-1, -1}},
				{Pixel{1, 1}, Pixel{0, 1}},
				{Pixel{0, 1}, Pixel{1, 1}},
			},
			true,
		},
		{
			"not origin centered and cannot be y shifted to center origin",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 2}},
				{Pixel{-1, 2}, Pixel{1, 2}},
				{Pixel{1, 2}, Pixel{1, -1}},
				{Pixel{1, -1}, Pixel{-1, -1}},
			},
			true,
		},
		{
			"valid polygon when x shifted to center origin",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 1}},
				{Pixel{-1, 1}, Pixel{3, 1}},
				{Pixel{3, 1}, Pixel{3, -1}},
				{Pixel{3, -1}, Pixel{-1, -1}},
			},
			false,
		},
		{
			"valid polygon when y shifted to center origin",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 3}},
				{Pixel{-1, 3}, Pixel{1, 3}},
				{Pixel{1, 3}, Pixel{1, -1}},
				{Pixel{1, -1}, Pixel{-1, -1}},
			},
			false,
		},
		{
			"valid polygon when x and y shifted to center origin",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 3}},
				{Pixel{-1, 3}, Pixel{3, 3}},
				{Pixel{3, 3}, Pixel{3, -1}},
				{Pixel{3, -1}, Pixel{-1, -1}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pixel{-1, -1}, Pixel{-1, 1}},
				{Pixel{-1, 1}, Pixel{1, 1}},
				{Pixel{1, 1}, Pixel{1, -1}},
				{Pixel{1, -1}, Pixel{-1, -1}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pixel{-2, -2}, Pixel{-2, 2}},
				{Pixel{-2, 2}, Pixel{2, 2}},
				{Pixel{2, 2}, Pixel{2, -2}},
				{Pixel{2, -2}, Pixel{-2, -2}},
			},
			false,
		},
		{
			"valid polygons without x or y shift to center origin",
			[]line{
				{Pixel{-2, -2}, Pixel{-2, 2}},
				{Pixel{-2, 2}, Pixel{2, 2}},
				{Pixel{2, 2}, Pixel{2, -2}},
				{Pixel{2, -2}, Pixel{-2, -2}},
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
