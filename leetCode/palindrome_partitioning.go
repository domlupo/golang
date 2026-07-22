// https://leetcode.com/problems/palindrome-partitioning/
// N is the length of s
// time complexity is O(N * 2^N)
// space complexity is O(N)

package main

func partition131(s string) [][]string {
	var ans [][]string
	var path []string
	backtrack131(s, path, &ans)
	return ans
}

func backtrack131(s string, path []string, ans *[][]string) {
	if len(s) == 0 {
		deepCopy := make([]string, len(path))
		copy(deepCopy, path)
		*ans = append(*ans, deepCopy)
		return
	}

	for i := 1; i <= len(s); i++ {
		if isPalindrome131(s[:i]) {
			backtrack131(s[i:], append(path, s[:i]), ans)
		}
	}
}

func isPalindrome131(s string) bool {
	var l int
	r := len(s) - 1
	for l < r {
		if s[l] == s[r] {
			l++
			r--
		} else {
			return false
		}
	}
	return true
}
