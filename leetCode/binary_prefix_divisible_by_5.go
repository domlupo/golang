// https://leetcode.com/problems/binary-prefix-divisible-by-5/

package main

func prefixesDivBy5(nums []int) []bool {
	var ans []bool
	var remainder int

	for _, n := range nums {
		remainder = (remainder*2 + n) % 5
		var divisible bool
		if remainder == 0 {
			divisible = true
		}
		ans = append(ans, divisible)
	}

	return ans
}
