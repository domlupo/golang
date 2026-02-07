// https://leetcode.com/problems/two-sum/

package main

import "strconv"

func convertToBase7(num int) string {
	if num == 0 {
		return "0"
	}

	var negative bool
	if num < 0 {
		negative = true
		num *= -1
	}

	var s string
	for num > 0 {
		s = strconv.Itoa(num%7) + s
		num /= 7
	}

	if negative {
		return "-" + s
	} else {
		return s
	}
}
