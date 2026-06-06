// https://leetcode.com/problems/palindrome-number/

package main

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	rem := x
	var reversed int
	for rem > 0 {
		reversed = reversed*10 + rem%10
		rem /= 10
	}

	return x == reversed
}
