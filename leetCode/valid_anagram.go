// https://leetcode.com/problems/valid-anagram/
// N is the max length of s and t
// time complexity is O(N)
// space complexity is O(1)

package main

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	var counts [26]int
	for i := 0; i < len(s); i++ {
		counts[s[i]-'a']++
		counts[t[i]-'a']--
	}

	for _, count := range counts {
		if count != 0 {
			return false
		}
	}

	return true
}
