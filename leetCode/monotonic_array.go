// https://leetcode.com/problems/monotonic-array/

package main

func isMonotonic(nums []int) bool {
	decreasing := true
	increasing := true

	for i := 0; i < len(nums)-1; i++ {
		if nums[i+1] < nums[i] {
			increasing = false
		}
		if nums[i+1] > nums[i] {
			decreasing = false
		}
	}

	return decreasing || increasing
}
