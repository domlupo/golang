// https://leetcode.com/problems/longest-common-subsequence/
// N is the length of text1, M is the length of text2
// time complexity is O(M * N)
// space complexity is O(min(M, N))

package main

func longestCommonSubsequence(text1 string, text2 string) int {
	if len(text2) > len(text1) {
		text1, text2 = text2, text1
	}

	dp := make([]int, len(text2)+1)

	for r := range len(text1) {
		curr := make([]int, len(text2)+1)
		for c := 1; c < len(text2)+1; c++ {
			curr[c] = max(curr[c-1], dp[c])
			if text1[r] == text2[c-1] {
				curr[c] = max(curr[c], dp[c-1]+1)
			}
		}
		dp = curr
	}

	return dp[len(text2)]
}

// 0 a a b b
// b 0 0 1 1
// b 0 0 1 2
// a 1 1 1 2
// a 1 2 2 2

// 0 c r a b t
// c 1 1 1 1 1
// a 1 1 2 2 2
// t 1 1 2 2 3

// 0 b s b i n i n m
// j 0 0 0 0 0 0 0 0
// m 0 0 0 0 0 0 0 1
// j 0 0 0 0 0 0 0 1
// k 0 0 0 0 0 0 0 1
// b 1 1 1 1 1 1 1 1
// k 1 1 1 1 1 1 1 1
// j 1 1 1 1 1 1 1 1
// k 1 1 1 1 1 1 1 1
// v 1 1 1 1 1 1 1 1
