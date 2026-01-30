// https://leetcode.com/problems/roman-to-integer/description/

package main

import (
	"testing"
)

func TestRomanToInteger(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		roman  string
		output int
	}{
		{"III", 3},
		{"LVIII", 58},
		{"MCMXCIV", 1994},
	}

	for _, tt := range tests {
		t.Run("romanToInteger", func(t *testing.T) {
			if romanToInt(tt.roman) != tt.output {
				t.Errorf("incorrect")
			}
		})
	}
}
