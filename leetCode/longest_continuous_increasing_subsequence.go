// https://leetcode.com/problems/longest-continuous-increasing-subsequence/

package main

func findLengthOfLCIS(nums []int) int {
	var longest, count int

	for i := 0; i < len(nums)-1; i++ {
		if nums[i] < nums[i+1] {
			count++
			longest = max(longest, count)
		} else {
			count = 0
		}
	}

	return longest + 1
}
