// https://leetcode.com/problems/reverse-only-letters/

package main

func reverseOnlyLetters(s string) string {
	c := []byte(s)
	var l int
	r := len(c) - 1

	for l < r {
		if !isChar(c[l]) {
			l++
		} else if !isChar(c[r]) {
			r--
		} else {
			c[l], c[r] = c[r], c[l]
			l++
			r--
		}
	}

	return string(c)
}

func isChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
