// https://leetcode.com/problems/coin-change/
// N is the length of coins, M is the amount
// time complexity is O(N*M)
// space complexity is O(M)

package main

func coinChange(coins []int, amount int) int {
	max := 10001
	dp := make([]int, amount+1)
	for i := 1; i < len(dp); i++ {
		dp[i] = max
	}

	for _, c := range coins {
		for a := c; a <= amount; a++ {
			dp[a] = min(dp[a], 1+dp[a-c])
		}
	}

	if dp[amount] == max {
		return -1
	} else {
		return dp[amount]
	}
}
