package main

import (
	"fmt"
	"testing"
)

func TestNewShape(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		boundary  []line
		wantPanic bool
	}{
		// Basic sanity checks for obvious invalid shapes
		{[]line{
			{}},
			true},
		{[]line{
			{},
			{}},
			true},
		{[]line{
			{},
			{},
			{}},
			true},
		{[]line{
			{PixelPosition{0, 0}, PixelPosition{1, 0}},
			{PixelPosition{1, 0}, PixelPosition{0, 1}},
			{PixelPosition{0, 1}, PixelPosition{0, 1}}},
			true},
		{[]line{
			{PixelPosition{0, 0}, PixelPosition{1, 0}},
			{PixelPosition{1, 0}, PixelPosition{0, 1}},
			{PixelPosition{0, 1}, PixelPosition{0, 2}}},
			true},
		{[]line{
			{PixelPosition{0, 0}, PixelPosition{1, 0}},
			{PixelPosition{1, 0}, PixelPosition{2, 0}},
			{PixelPosition{2, 0}, PixelPosition{0, 0}}},
			false},
		{[]line{
			{PixelPosition{0, 0}, PixelPosition{1, 0}},
			{PixelPosition{1, 0}, PixelPosition{1, 1}},
			{PixelPosition{1, 1}, PixelPosition{0, 1}},
			{PixelPosition{0, 1}, PixelPosition{0, 0}}},
			false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("newShape %v", tt.boundary)
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
			newShape(tt.boundary)
		})
	}

}
