// https://leetcode.com/problems/two-sum/

package main

import (
	"slices"
	"testing"
)

func TestTwoSum(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		target int
		nums   []int
		output []int
	}{
		{9, []int{2, 7, 11, 15}, []int{0, 1}},
		{6, []int{3, 2, 4}, []int{1, 2}},
		{6, []int{3, 3}, []int{0, 1}},
	}

	for _, tt := range tests {
		t.Run("twoSum", func(t *testing.T) {
			if !slices.Equal(twoSum(tt.nums, tt.target), tt.output) {
				t.Errorf("incorrect")
			}
		})
	}
}
