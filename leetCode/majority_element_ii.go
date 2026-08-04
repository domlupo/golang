// https://leetcode.com/problems/majority-element-ii/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func majorityElement229(nums []int) []int {
	first, second := math.MaxInt, math.MaxInt
	var firstCount, secondCount int
	// first pass to find potential majorities
	for _, n := range nums {
		if n == first {
			firstCount++
		} else if n == second {
			secondCount++
		} else if firstCount == 0 {
			first = n
			firstCount++
		} else if secondCount == 0 {
			second = n
			secondCount++
		} else {
			firstCount--
			secondCount--
		}
	}

	var ans []int
	firstCount, secondCount = 0, 0
	// second pass to verify actual majorities
	for _, n := range nums {
		if n == first {
			firstCount++
		} else if n == second {
			secondCount++
		}
	}
	if firstCount > len(nums)/3 {
		ans = append(ans, first)
	}
	if secondCount > len(nums)/3 {
		ans = append(ans, second)
	}
	return ans
}
