// https://leetcode.com/problems/smallest-range-i/

package main

import "math"

func smallestRangeI(nums []int, k int) int {
	minimum := math.MaxInt
	maximum := math.MinInt

	for _, n := range nums {
		minimum = min(minimum, n)
		maximum = max(maximum, n)
	}

	diff := maximum - minimum - 2*k
	if diff < 0 {
		return 0
	} else {
		return diff
	}
}
