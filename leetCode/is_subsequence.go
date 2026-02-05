// https://leetcode.com/problems/is-subsequence/

package main

func isSubsequence(s string, t string) bool {
	var p1, p2 int

	for p1 < len(s) && p2 < len(t) {
		if s[p1] == t[p2] {
			p1++
		}
		p2++
	}

	return p1 == len(s)
}
