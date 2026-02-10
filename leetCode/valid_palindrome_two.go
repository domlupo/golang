// https://leetcode.com/problems/valid-palindrome-ii/

package main

func validPalindrome(s string) bool {
	var l int
	r := len(s) - 1

	for l <= r {
		if s[l] == s[r] {
			l++
			r--
		} else {
			return isPalindrome1(s[(l+1):r+1]) || isPalindrome1(s[l:r])
		}
	}

	return true
}

func isPalindrome1(s string) bool {
	var l int
	r := len(s) - 1
	for l <= r {
		if s[l] != s[r] {
			return false
		}
		l++
		r--
	}

	return true
}
