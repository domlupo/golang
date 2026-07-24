// https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/
// N is the length of nums
// time complexity is O(log(N))
// space complexity is O(1)

package main

func findMin(nums []int) int {
	var l int
	r := len(nums) - 1
	for l < r {
		m := (l + r) / 2
		if nums[m] > nums[len(nums)-1] {
			l = m + 1
		} else {
			r = m
		}
	}

	return nums[l]
}
