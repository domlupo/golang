// https://leetcode.com/problems/valid-anagram/

package main

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	counts := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		counts[s[i]] = counts[s[i]] + 1
		counts[t[i]] = counts[t[i]] - 1
	}

	for _, v := range counts {
		if v != 0 {
			return false
		}
	}
	return true
}
