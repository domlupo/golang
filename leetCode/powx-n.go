// https://leetcode.com/problems/powx-n/
// n is power
// time complexity is O(log(n))
// space complexity is O(1)

package main

func myPow(x float64, n int) float64 {
	return binaryExp(x, n)
}

func binaryExp(x float64, n int) float64 {
	if n == 0 {
		return 1
	}

	if n < 0 {
		n = -1 * n
		x = 1.0 / x
	}

	ans := 1.0
	for n != 0 {
		if n%2 == 1 {
			ans = ans * x
			n--
		}
		x = x * x
		n = n / 2
	}

	return ans
}
