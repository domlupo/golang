// https://leetcode.com/problems/valid-palindrome/

package main

func isPalindrome2(s string) bool {
	var p1 int
	p2 := len(s) - 1

	for p1 < p2 {
		if !isAlphaNumeric(s[p1]) {
			p1++
			continue
		}
		if !isAlphaNumeric(s[p2]) {
			p2--
			continue
		}
		if lower(s[p1]) != lower(s[p2]) {
			return false
		}
		p1++
		p2--
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
	if (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') {
		return true
	}
	return false
}
