// https://leetcode.com/problems/palindrome-number/
// N is the value of x
// time complexity is O(log10(N))
// space complexity is O(1)

package main

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	rem := x
	var rev int
	for rem > 0 {
		rev = rev*10 + rem%10
		rem /= 10
	}

	return x == rev
}
