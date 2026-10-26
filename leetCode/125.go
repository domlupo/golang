// https://leetcode.com/problems/valid-palindrome/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func isPalindromeTwo(s string) bool {
	l, r := 0, len(s)-1
	for l < r {
		if !isAlphaNumeric(s[l]) {
			l++
		} else if !isAlphaNumeric(s[r]) {
			r--
		} else if toLower(s[l]) != toLower(s[r]) {
			return false
		} else {
			l++
			r--
		}

	}
	return true
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

func isAlphaNumeric(b byte) bool {
	lower := b >= 'a' && b <= 'z'
	upper := b >= 'A' && b <= 'Z'
	numeric := b >= '0' && b <= '9'
	return lower || upper || numeric
}
