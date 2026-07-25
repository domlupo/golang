// https://leetcode.com/problems/find-peak-element/
// N is the length of nums
// time complexity is O(log(N))
// space complexity is O(1)

package main

func findPeakElement(nums []int) int {
	if len(nums) >= 2 && nums[0] > nums[1] {
		return 0
	} else if len(nums) >= 2 && nums[len(nums)-1] > nums[len(nums)-2] {
		return len(nums) - 1
	}

	var l int
	r := len(nums) - 1
	for l < r {
		m := (l + r) / 2
		if m > 0 && m < len(nums)-1 && nums[m] > nums[m-1] && nums[m] > nums[m+1] {
			return m
		} else if m > 0 && nums[m] > nums[m+1] {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	return l
}
