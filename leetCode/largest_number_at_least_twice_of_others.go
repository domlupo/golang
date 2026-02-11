// https://leetcode.com/problems/largest-number-at-least-twice-of-others/

package main

import "math"

func dominantIndex(nums []int) int {
	var firstI int
	first := -math.MaxInt
	second := -math.MaxInt

	for i := 0; i < len(nums); i++ {
		if nums[i] > first {
			second = first
			first = nums[i]
			firstI = i
		} else if nums[i] > second {
			second = nums[i]
		}
	}

	if first >= second*2 {
		return firstI
	} else {
		return -1
	}
}
