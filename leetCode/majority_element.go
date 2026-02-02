// https://leetcode.com/problems/majority-element/

package main

func majorityElement(nums []int) int {
	var e, c int

	for _, n := range nums {
		if c == 0 {
			e = n
		}
		if e == n {
			c++
		} else {
			c--
		}
	}

	return e
}
