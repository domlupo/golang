// https://leetcode.com/problems/search-in-rotated-sorted-array/
// N is the length of nums
// time complexity is O(log(N))
// space complexity is O(1)

package main

func search33(nums []int, target int) int {
	l, r := 0, len(nums)-1

	// set both l and r to split point which is the lowest value in nums
	for l < r {
		m := (l + r) / 2
		if nums[m] < nums[len(nums)-1] {
			r = m
		} else {
			l = m + 1
		}
	}

	split := l
	if target >= nums[split] && target <= nums[len(nums)-1] {
		l = split
		r = len(nums) - 1
	} else {
		l = 0
		r = split - 1
	}

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
