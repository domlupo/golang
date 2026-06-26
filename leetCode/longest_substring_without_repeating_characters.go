// https://leetcode.com/problems/longest-substring-without-repeating-characters/

package main

func lengthOfLongestSubstring(s string) int {
	lastSeen := make(map[rune]int)
	var left int
	var ans int

	for right, c := range s {
		if i, ok := lastSeen[c]; ok {
			left = max(left, i)
		}
		ans = max(ans, right-left+1)
		lastSeen[c] = right + 1
	}

	return ans
}
