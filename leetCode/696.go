// https://leetcode.com/problems/count-binary-substrings/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func countBinarySubstrings(s string) int {
	var ans, prev int
	curr := 1
	for i := 1; i <= len(s); i++ {
		if i == len(s) || s[i-1] != s[i] {
			ans += min(prev, curr)
			prev, curr = curr, 1
		} else {
			curr++
		}
	}
	return ans
}
