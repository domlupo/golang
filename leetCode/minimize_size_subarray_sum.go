// https://leetcode.com/problems/minimum-size-subarray-sum/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func minSubArrayLen(target int, nums []int) int {
	var l, r int
	var sum int
	ans := math.MaxInt
	for r < len(nums) {
		sum += nums[r]
		r++
		for sum >= target {
			ans = min(ans, r-l)
			sum -= nums[l]
			l++
		}
	}

	if ans == math.MaxInt {
		return 0
	} else {
		return ans
	}
}
