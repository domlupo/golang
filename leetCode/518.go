// https://leetcode.com/problems/coin-change-ii/
// C is the length of coins, A is the amount
// time complexity is O(C * A)
// space complexity is O(A)

package main

func change(amount int, coins []int) int {
	dp := make([]int, amount+1)
	dp[0] = 1

	for _, c := range coins {
		for a := c; a <= amount; a++ {
			dp[a] += dp[a-c]
		}
	}

	return dp[amount]
}
