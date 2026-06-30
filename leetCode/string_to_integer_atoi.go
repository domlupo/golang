// https://leetcode.com/problems/string-to-integer-atoi/

package main

import "math"

func myAtoi(s string) int {
	var i int
	var ans int
	digits := map[byte]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
	}

	// skip white space
	for i < len(s) && s[i] == ' ' {
		i++
	}

	// parse sign
	var negative bool
	if i < len(s) && s[i] == '-' {
		negative = true
		i++
	} else if i < len(s) && s[i] == '+' {
		i++
	}

	// skip leading zeros
	for i < len(s) && s[i] == '0' {
		i++
	}

	// parse digits
	for i < len(s) {
		digit, isDigit := digits[s[i]]
		if !isDigit {
			break
		}

		// round 32 bit overflow and underflow
		ans = ans*10 + digit
		if !negative && ans >= math.MaxInt32 {
			return math.MaxInt32
		} else if negative && ans*-1 <= math.MinInt32 {
			return math.MinInt32
		}

		i++
	}

	if negative {
		return ans * -1
	} else {
		return ans
	}
}
