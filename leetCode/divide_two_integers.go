// https://leetcode.com/problems/divide-two-integers/

package main

import "math"

func divide(dividend int, divisor int) int {
	// special case for overflow. Value must be truncated.
	if dividend == math.MinInt32 && divisor == -1 {
		return math.MaxInt32
	}

	negative := (dividend < 0) != (divisor < 0)

	// convert to negative instead of positive because negative bound can
	// handle values 1 digit larger 2,147,483,647 vs -2,147,483,648
	if dividend > 0 {
		dividend = -dividend
	}
	if divisor > 0 {
		divisor = -divisor
	}

	// exponential search instead of linear search for better time complexity
	// time complexity is o(log(n) * log(n)) which is much better than o(n).
	var count int
	for dividend <= divisor {
		powerOfTwo := 1
		value := divisor

		for value+value >= dividend {
			value += value
			powerOfTwo += powerOfTwo
		}

		count += powerOfTwo
		dividend -= value
	}

	if negative {
		return -count
	} else {
		return count
	}
}
