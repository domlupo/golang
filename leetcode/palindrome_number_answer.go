// https://leetcode.com/problems/palindrome-number/

package main

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	xRemainder := x
	var reversed int
	for i := 0; xRemainder > 0; i++ {
		rightDigit := xRemainder % 10
		reversed = reversed*10 + rightDigit
		xRemainder /= 10

		if xRemainder == 0 {
			break
		}
	}

	return x == reversed
}
