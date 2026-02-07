//https://leetcode.com/problems/perfect-number/

package main

import "math"

func checkPerfectNumber(num int) bool {
	if num == 1 {
		return false
	}

	sum := 1
	m := int(math.Sqrt(float64(num))) + 1

	for i := 2; i < m; i++ {
		if num%i == 0 {
			sum += i
			sum += num / i
		}
	}
	return sum == num
}
