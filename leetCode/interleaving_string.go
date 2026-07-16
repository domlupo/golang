// https://leetcode.com/problems/interleaving-string/
// M is length of s1 and N is the length of s2
// time complexity is O(M * N)
// space complexity is O(N)

package main

func isInterleave(s1 string, s2 string, s3 string) bool {
	// This is needed to simplify logic by a lot.
	// Specifically, to avoid many index out of bounds
	if len(s1)+len(s2) != len(s3) {
		return false
	}

	dp := make([]bool, len(s2)+1)

	for i := 0; i <= len(s1); i++ {
		for j := 0; j <= len(s2); j++ {
			if i == 0 && j == 0 {
				dp[0] = true
			} else if i == 0 {
				dp[j] = dp[j-1] && s2[j-1] == s3[j-1]
			} else if j == 0 {
				dp[j] = dp[j] && s1[i-1] == s3[i-1]
			} else {
				first := dp[j] && s1[i-1] == s3[i+j-1]
				second := dp[j-1] && s2[j-1] == s3[i+j-1]
				dp[j] = first || second
			}
		}
	}

	return dp[len(s2)]
}
