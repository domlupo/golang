// https://leetcode.com/problems/maximum-average-subarray-i/

package main

import "math"

func findMaxAverage(nums []int, k int) float64 {
	var sum int
	for i := 0; i < k-1; i++ {
		sum += nums[i]
	}

	m := float64(-math.MaxInt)
	for i := k - 1; i < len(nums); i++ {
		sum += nums[i]
		m = max(m, float64(sum)/float64(k))
		sum -= nums[i+1-k]
	}

	return m
}
