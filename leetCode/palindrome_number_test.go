// https://leetcode.com/problems/two-sum/

package main

import (
	"testing"
)

func TestPalindrome(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		num    int
		output bool
	}{
		{121, true},
		{-121, false},
		{10, false},
	}

	for _, tt := range tests {
		t.Run("palindrome", func(t *testing.T) {
			if isPalindrome(tt.num) != tt.output {
				t.Errorf("incorrect")
			}
		})
	}
}
