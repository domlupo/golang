// https://neetcode.io/problems/counting-bits
// N is n value.
// time complexity is O(N)
// space complexity is O(1)

package main

func countBits(n int) []int {
	dp := make([]int, n+1)
	offset, resetOffset := 0, 2
	for i := 1; i < n+1; i++ {
		if i == resetOffset {
			offset = i
			resetOffset *= 2
		}
		dp[i] = dp[i-offset] + 1
	}
	return dp
}

// 0000
// 0001
// 0010
// 0011
// 0100
// 0101
// 0110
// 0111
// 1000
// 1001
// 1010
// 1011
// 1100
// 1101
// 1110
// 1111
