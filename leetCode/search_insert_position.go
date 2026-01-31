// https://leetcode.com/problems/search-insert-position/

package main

func searchInsert(nums []int, target int) int {
	start := 0
	end := len(nums) - 1

	for start <= end {
		mid := (start + end) / 2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return start
}
