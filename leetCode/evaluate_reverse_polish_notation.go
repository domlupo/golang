// https://leetcode.com/problems/evaluate-reverse-polish-notation/
// N is the length of tokens
// time complexity is O(N)
// space complexity is O(N)

package main

import "strconv"

func evalRPN(tokens []string) int {
	var stack []int
	for _, token := range tokens {
		number, err := strconv.Atoi(token)
		if err == nil {
			stack = append(stack, number)
			continue
		}

		var result int
		left := stack[len(stack)-2]
		right := stack[len(stack)-1]
		stack = stack[:len(stack)-2]
		switch token {
		case "+":
			result = left + right
		case "-":
			result = left - right
		case "*":
			result = left * right
		case "/":
			result = left / right
		}
		stack = append(stack, result)
	}

	return stack[0]
}
