// https://leetcode.com/problems/third-maximum-number/solutions/

package main

import "math"

func thirdMax(nums []int) int {
	first := nums[0]
	second := math.MinInt
	third := math.MinInt

	for _, n := range nums {
		if n == first || n == second || n == third {
			continue
		}

		if n > first {
			third = second
			second = first
			first = n
		} else if n > second {
			third = second
			second = n
		} else if n > third {
			third = n
		}
	}

	if third == math.MinInt {
		return first
	} else {
		return third
	}
}
