// https://leetcode.com/problems/valid-parentheses/description/

package main

import "fmt"

func isValid(s string) bool {
	var stack []rune
	for _, r := range s {
		switch r {
		case '(', '{', '[':
			stack = append(stack, r)
		case ')', '}', ']':
			if len(stack) > 0 && leftBracket(r) == stack[len(stack)-1] {
				stack = stack[0 : len(stack)-1]
			} else {
				return false
			}
		default:
			panic(fmt.Sprintf("isValid: unknown rune %c", r))
		}
	}

	return len(stack) == 0
}

func leftBracket(r rune) rune {
	switch r {
	case ')':
		return '('
	case '}':
		return '{'
	case ']':
		return '['
	default:
		panic(fmt.Sprintf("openBracket: unknown rune %c", r))
	}
}
