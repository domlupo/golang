// https://leetcode.com/problems/search-in-rotated-sorted-array/

package main

func search1(nums []int, target int) int {
	// binary search to find split point
	var l int
	r := len(nums) - 1
	for l <= r {
		m := (l + r) / 2
		if nums[m] > nums[len(nums)-1] {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	// binary search both sides of split point
	left := s(nums, 0, l-1, target)
	right := s(nums, l, len(nums)-1, target)
	if left == -1 {
		return right
	} else {
		return left
	}
}

func s(nums []int, l, r, target int) int {
	for l <= r {
		m := (l + r) / 2
		if nums[m] == target {
			return m
		} else if nums[m] > target {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}
