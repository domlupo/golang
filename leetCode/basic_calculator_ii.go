// https://leetcode.com/problems/basic-calculator-ii/
// N is the length of string
// time complexity is O(N)
// space complexity is O(1)

package main

import "strconv"

func calculate(s string) int {
	var ans, currNum, lastNum int
	sign := byte('+')
	bytes := []byte(s)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == ' ' {
			continue
		}

		for i < len(bytes) && bytes[i] >= '0' && bytes[i] <= '9' {
			digit, _ := strconv.Atoi(string(bytes[i]))
			currNum = currNum*10 + digit
			i++
		}

		switch sign {
		case '+':
			ans += lastNum
			lastNum = currNum
		case '-':
			ans += lastNum
			lastNum = -currNum
		case '*':
			lastNum = lastNum * currNum
		case '/':
			lastNum = lastNum / currNum
		}

		currNum = 0
		if i < len(bytes) {
			sign = bytes[i]
		}
	}

	return ans + lastNum
}
