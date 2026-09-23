// https://leetcode.com/problems/palindromic-substrings/
// N is the length of nums
// time complexity is O(N^2)
// space complexity is O(1)

package main

func countSubstrings(s string) int {
	var ans int
	for i := 0; i < len(s); i++ {
		ans += count(i-1, i, s)
		ans += count(i, i, s)
	}
	return ans
}

func count(l, r int, s string) int {
	var count int
	for l >= 0 && r < len(s) && s[l] == s[r] {
		count++
		l--
		r++
	}
	return count
}
