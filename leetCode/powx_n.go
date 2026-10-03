// https://leetcode.com/problems/powx-n/
// n is power
// time complexity is O(log(n))
// space complexity is O(1)

package main

func myPow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}

	var negative bool
	if n < 0 {
		negative = true
		n *= -1
	}

	product := 1.0
	for n != 0 {
		if n%2 == 1 {
			product *= x
			n--
		} else {
			x *= x
			n /= 2
		}
	}

	if negative {
		return 1 / product
	} else {
		return product
	}
}

// 2 4 8 16 32 64 128 256 512 1024
// 2^10 = 1024
// 4^5
// 4 * 4^4
// 4 * 16^2
// 4 * 256
