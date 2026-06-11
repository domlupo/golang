// https://leetcode.com/problems/valid-palindrome/

package main

func isPalindrome2(s string) bool {
	var l int
	r := len(s) - 1

	for l < r {
		for l < r && !isAlphaNumeric(s[l]) {
			l++
		}
		for l < r && !isAlphaNumeric(s[r]) {
			r--
		}
		if lower(s[l]) != lower(s[r]) {
			return false
		}
		l++
		r--
	}

	return true
}

func lower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

func isAlphaNumeric(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}
