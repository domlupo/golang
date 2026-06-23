// https://leetcode.com/problems/search-insert-position/

package main

func searchInsert(nums []int, target int) int {
	l := 0
	r := len(nums) - 1

	// when target is not found, l is always equal to r + 1
	for l <= r {
		m := (l + r) / 2
		if nums[m] < target {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return l
}
