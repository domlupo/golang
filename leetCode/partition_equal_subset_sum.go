// https://leetcode.com/problems/partition-equal-subset-sum/
// N is the length of nums, and M is the subset sum (sum/2)
// time complexity is O(N * M)
// space complexity is O(M)

package main

func canPartition(nums []int) bool {
	var sum int
	for _, n := range nums {
		sum += n
	}

	if sum%2 == 1 {
		return false
	}

	target := (sum / 2)
	dp := make([]bool, target+1)
	dp[0] = true

	for _, n := range nums {
		for j := target; j >= n; j-- {
			dp[j] = dp[j] || dp[j-n]
		}
	}

	return dp[len(dp)-1]
}
