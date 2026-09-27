// https://leetcode.com/problems/word-break/
// N is the length of s, M is len of wordDict, K is the avg length of words in M
// time complexity is O(N * M * K)
// space complexity is O(N)

package main

func wordBreak(s string, wordDict []string) bool {
	dp := make([]bool, len(s)+1)
	dp[0] = true

	for i := range len(s) + 1 {
		for _, word := range wordDict {
			if i >= len(word) && dp[i-len(word)] && s[i-len(word):i] == word {
				dp[i] = true
				break
			}
		}
	}

	return dp[len(s)]
}
