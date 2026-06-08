// https://leetcode.com/problems/valid-parentheses/description/

package main

func isValid(s string) bool {
	leftBrackets := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	var stack []rune
	for _, r := range s {
		leftBracket, ok := leftBrackets[r]
		if ok && len(stack) > 0 && stack[len(stack)-1] == leftBracket {
			stack = stack[:len(stack)-1]
		} else if ok {
			return false
		} else {
			stack = append(stack, r)
		}
	}

	return len(stack) == 0
}
