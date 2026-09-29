// https://leetcode.com/problems/valid-parenthesis-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(N)

package main

func checkValidString(s string) bool {
	var open, close int

	for i := 0; i < len(s); i++ {
		l := s[i]
		if l == '(' || l == '*' {
			open++
		} else {
			open--
		}

		r := s[len(s)-1-i]
		if r == ')' || r == '*' {
			close++
		} else {
			close--
		}

		if open < 0 || close < 0 {
			return false
		}
	}

	return true
}
