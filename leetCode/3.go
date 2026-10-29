// https://leetcode.com/problems/longest-substring-without-repeating-characters/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func lengthOfLongestSubstring(s string) int {
	lastSeen := make(map[rune]int)
	var ans, lastRepeat int
	for i, r := range s {
		lastRepeat = max(lastRepeat, lastSeen[r])
		ans = max(ans, i-lastRepeat+1)
		lastSeen[r] = i + 1
	}
	return ans
}
