// https://leetcode.com/problems/palindrome-number/

package main

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	xRemainder := x
	var reversed int
	for xRemainder > 0 {
		reversed = reversed*10 + xRemainder%10
		xRemainder /= 10
	}

	return x == reversed
}
