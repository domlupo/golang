// https://leetcode.com/problems/max-consecutive-ones/

package main

func findMaxConsecutiveOnes(nums []int) int {
	var m, c int

	for _, n := range nums {
		c = c*n + n
		if c > m {
			m = c
		}
	}

	return m
}
