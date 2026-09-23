// https://leetcode.com/problems/longest-palindromic-substring/
// N is the length of string
// time complexity is O(N^2)
// space complexity is O(1)

package main

func longestPalindrome1(s string) string {
	var ans string

	for i := 0; i < len(s); i++ {
		odd := expand(s, i, i)
		if len(odd) > len(ans) {
			ans = odd
		}
		even := expand(s, i, i+1)
		if len(even) > len(ans) {
			ans = even
		}
	}
	return ans
}

func expand(s string, l, r int) string {
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}
	return s[l+1 : r]
}
