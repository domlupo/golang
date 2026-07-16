// https://leetcode.com/problems/3sum-closest/
// N is the length of nums
// time complexity is O(N * N)
// space complexity is O(log(N)) from slices.Sort

package main

import (
	"math"
	"slices"
)

func threeSumClosest(nums []int, target int) int {
	slices.Sort(nums)
	minDiff := math.MaxInt
	ans := nums[0] + nums[1] + nums[2]

	for i := 0; i < len(nums)-2; i++ {
		l := i + 1
		r := len(nums) - 1

		for l < r {
			sum := nums[i] + nums[l] + nums[r]
			diff := sum - target
			if diff < 0 {
				diff *= -1
			}

			minDiff = min(minDiff, diff)
			if diff == minDiff {
				ans = sum
			}

			if sum == target {
				return ans
			} else if sum < target {
				l++
			} else {
				r--
			}
		}
	}

	return ans
}
