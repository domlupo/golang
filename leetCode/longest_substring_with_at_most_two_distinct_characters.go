// https://leetcode.com/problems/longest-substring-with-at-most-two-distinct-characters/
// N is the length of string
// time complexity is O(N)
// space complexity is O(1)

package main

func lengthOfLongestSubstringTwoDistinct(s string) int {
	var longest int
	l1, r1, l2, r2 := -1, -1, -1, -1
	for i := range len(s) {
		if l1 == -1 {
			l1, r1 = i, i
		} else if l2 == -1 && s[i] != s[l1] {
			l2, r2 = i, i
		}

		if s[i] == s[l1] {
			r1 = i
		} else if s[i] == s[l2] {
			r2 = i
		} else if r1 > r2 {
			l1 = r2 + 1
			l2, r2 = i, i
		} else if r2 > r1 {
			l1 = r1 + 1
			r1 = r2
			l2, r2 = i, i
		}
		longest = max(longest, i-l1+1)
	}

	return longest
}
