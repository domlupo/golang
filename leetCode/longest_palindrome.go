// https://leetcode.com/problems/longest-palindrome/

package main

func longestPalindrome(s string) int {
	counts := make(map[rune]int)
	var middleChar int
	var l int

	for _, c := range s {
		counts[c] = counts[c] + 1
		if counts[c] == 2 {
			counts[c] = 0
			l += 2
			middleChar--
		} else {
			middleChar++
		}
	}

	if middleChar > 0 {
		return l + 1
	} else {
		return l
	}
}
