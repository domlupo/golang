// https://leetcode.com/problems/longest-common-prefix/description/

package main

import (
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		strings []string
		output  string
	}{
		{[]string{"flower", "flow", "flight"}, "fl"},
		{[]string{"dog", "racecar", "car"}, ""},
	}

	for _, tt := range tests {
		t.Run("longestCommonPrefix", func(t *testing.T) {
			if longestCommonPrefix(tt.strings) != tt.output {
				t.Errorf("incorrect")
			}
		})
	}
}
