// https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/

package main

func searchRange(nums []int, target int) []int {
	leftBound, rightBound := -1, -1

	var l int
	r := len(nums) - 1
	for l <= r {
		m := (l + r) / 2
		if nums[m] >= target {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	if l < len(nums) && nums[l] == target {
		leftBound = l
	}

	l = 0
	r = len(nums) - 1
	for l <= r {
		m := (l + r) / 2
		if nums[m] <= target {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	if r >= 0 && nums[r] == target {
		rightBound = r
	}

	return []int{leftBound, rightBound}
}
