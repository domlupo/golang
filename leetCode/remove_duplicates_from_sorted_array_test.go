// https://leetcode.com/problems/remove-duplicates-from-sorted-array

package main

import (
	"testing"
)

func TestRemoveDuplicatesFromSortedArray(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		nums         []int
		numsExpected []int
		output       int
	}{
		{[]int{1, 1, 2}, []int{1, 2}, 2},
		{[]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, []int{0, 1, 2, 3, 4}, 5},
	}

	for _, tt := range tests {
		t.Run("removeDuplicatesFromSortedArray", func(t *testing.T) {
			output := removeDuplicates(tt.nums)
			if /*!slices.Equal(tt.nums, tt.numsExpected) || */ tt.output != output {
				t.Errorf("incorrect")
			}
		})
	}
}
