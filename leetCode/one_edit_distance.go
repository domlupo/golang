// https://leetcode.com/problems/one-edit-distance/
// N is the length of string
// time complexity is O(N)
// space complexity is O(1)

package main

func isOneEditDistance(s string, t string) bool {
	if s == t {
		return false
	}

	var i int
	for i < len(s) && i < len(t) {
		if s[i] != t[i] {
			break
		}
		i++
	}

	if len(s) == len(t) {
		return s[i+1:] == t[i+1:]
	} else if len(s) < len(t) {
		return s[i:] == t[i+1:]
	} else {
		return s[i+1:] == t[i:]
	}
}
