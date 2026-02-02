// https://leetcode.com/problems/single-number/

package main

func singleNumber(nums []int) int {
	var ans int
	for _, n := range nums {
		ans ^= n
	}
	return ans
}
