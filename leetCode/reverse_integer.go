// https://leetcode.com/problems/reverse-integer/

package main

import "math"

func reverse(x int) int {
	negative := x < 0
	var rem int32
	if negative {
		rem = int32(x * -1)
	} else {
		rem = int32(x)
	}

	maximum := math.MaxInt32 / int32(10)
	overflow := math.MaxInt32 % int32(10)
	minimum := math.MinInt32 / int32(10) * int32(-1)
	underflow := math.MinInt32 % int32(10) * int32(-1)
	// maximum = 214748364
	// overflow = 7
	// minimum = 214748364
	// underflow = 8

	var ans int32
	for rem > 0 {
		digit := int32(rem % 10)
		if ans > minimum || ans == minimum && digit > underflow {
			return 0
		} else if ans > maximum || ans == maximum && digit > overflow {
			return 0
		}
		ans = ans*10 + digit
		rem /= 10
	}

	if negative {
		return int(ans) * -1
	} else {
		return int(ans)
	}
}
