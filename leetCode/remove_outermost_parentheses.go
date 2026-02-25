// https://leetcode.com/problems/remove-outermost-parentheses/

package main

func removeOuterParentheses(s string) string {
	var count int
	var ans []byte

	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			if count > 0 {
				ans = append(ans, s[i])
			}
			count++
		} else {
			count--
			if count > 0 {
				ans = append(ans, s[i])
			}
		}
	}

	return string(ans)
}
