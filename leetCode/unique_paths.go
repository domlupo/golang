// https://leetcode.com/problems/unique-paths/
// m is number of rows and n is numbers of columns
// time complexity is O(m*n)
// space complexity is O(n)

package main

func uniquePaths(m int, n int) int {
	dp := make([]int, n+1)

	for i := 1; i < n+1; i++ {
		dp[i] = 1
	}

	for range m - 1 {
		curr := make([]int, n+1)
		for c := 1; c < n+1; c++ {
			curr[c] = dp[c] + curr[c-1]
		}
		dp = curr
	}

	return dp[n]
}
