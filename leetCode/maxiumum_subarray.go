// https://leetcode.com/problems/maximum-subarray/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func maxSubArray(nums []int) int {
	var currSum int
	maxSum := math.MinInt

	for _, n := range nums {
		currSum += n
		if n > currSum {
			currSum = n
		}
		maxSum = max(maxSum, currSum)
	}

	return maxSum
}
