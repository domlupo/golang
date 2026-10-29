// https://leetcode.com/problems/maximum-subarray/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func maxSubArray(nums []int) int {
	var sum int
	maxSum := nums[0]
	for _, n := range nums {
		sum += n
		if n > sum {
			sum = n
		}
		maxSum = max(maxSum, sum)
	}
	return maxSum
}
