// https://leetcode.com/problems/first-unique-character-in-a-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func firstUniqChar(s string) int {
	var chars [26]int
	for _, c := range s {
		chars[c-'a']++
	}
	for ci, c := range s {
		if chars[c-'a'] == 1 {
			return ci
		}
	}
	return -1
}
