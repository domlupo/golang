// https://leetcode.com/problems/add-strings/

package main

import "strconv"

func addStrings(num1 string, num2 string) string {
	p1 := len(num1) - 1
	p2 := len(num2) - 1
	var carry bool
	var s string

	for p1 >= 0 || p2 >= 0 || carry {
		var total int
		if p1 >= 0 {
			i, _ := strconv.Atoi(string(num1[p1]))
			total += i
		}
		if p2 >= 0 {
			i, _ := strconv.Atoi(string(num2[p2]))
			total += i
		}
		if carry {
			total += 1
		}

		s = strconv.Itoa(total%10) + s

		if total/10 == 1 {
			carry = true
		} else {
			carry = false
		}

		p1--
		p2--
	}

	return s
}
