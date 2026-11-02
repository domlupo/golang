// https://leetcode.com/problems/longest-palindrome/
// N is the length of s
// time complexity is O(N*N)
// space complexity is O(1)

package main

func longestPalindrome(s string) string {
	var ans string
	for i := 0; i < len(s); i++ {
		odd := expand5(s, i, i)
		if len(odd) > len(ans) {
			ans = odd
		}
		even := expand5(s, i-1, i)
		if len(even) > len(ans) {
			ans = even
		}
	}
	return ans
}

func expand5(s string, l, r int) string {
	for l >= 0 && r < len(s) {
		if s[l] == s[r] {
			l--
			r++
		} else {
			break
		}
	}
	return s[l+1 : r]
}
