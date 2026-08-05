// https://leetcode.com/problems/different-ways-to-add-parentheses/
// N is the number of characters
// time complexity is O(2^N)
// space complexity is O(2^N)

package main

import "strconv"

func diffWaysToCompute(expression string) []int {
	var res []int
	for i := 0; i < len(expression); i++ {
		oper := expression[i]
		if oper == '+' || oper == '-' || oper == '*' {
			s1 := diffWaysToCompute(expression[:i])
			s2 := diffWaysToCompute(expression[i+1:])
			for _, a := range s1 {
				for _, b := range s2 {
					if oper == '+' {
						res = append(res, a+b)
					} else if oper == '-' {
						res = append(res, a-b)
					} else if oper == '*' {
						res = append(res, a*b)
					}
				}
			}
		}
	}

	if len(res) == 0 {
		num, _ := strconv.Atoi(expression)
		res = append(res, num)
	}

	return res
}
