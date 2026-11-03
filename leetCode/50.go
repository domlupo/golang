// https://leetcode.com/problems/powx-n/
// n is passed n parameter
// time complexity is O(log(n))
// space complexity is O(1)

package main

func myPow(x float64, n int) float64 {
	denominator := n < 0
	if denominator {
		n *= -1
	}

	val := 1.0
	for n > 0 {
		mult := x
		exp := 1
		for exp*2 <= n {
			exp *= 2
			mult *= mult
		}
		val *= mult
		n -= exp
	}

	if denominator {
		return 1.0 / val
	} else {
		return val
	}
}
