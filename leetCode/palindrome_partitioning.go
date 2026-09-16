// https://leetcode.com/problems/palindrome-partitioning/
// N is the length of s
// time complexity is O(N * 2^N)
// space complexity is O(N)

package main

func partition131(s string) [][]string {
	var ans [][]string
	backtrack131(&ans, []string{}, s)
	return ans
}

func backtrack131(ans *[][]string, strs []string, str string) {
	if len(str) == 0 {
		deepCopy := make([]string, len(strs))
		copy(deepCopy, strs)
		*ans = append(*ans, deepCopy)
	}

	for i := 1; i <= len(str); i++ {
		if isPalindrome131(str[:i]) {
			backtrack131(ans, append(strs, str[:i]), str[i:])
		}
	}
}

func isPalindrome131(str string) bool {
	var l int
	r := len(str) - 1
	for l <= r {
		if str[l] != str[r] {
			return false
		}
		l++
		r--
	}
	return true
}
