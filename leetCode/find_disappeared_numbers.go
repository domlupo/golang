// https://leetcode.com/problems/find-all-numbers-disappeared-in-an-array/

package main

func findDisappearedNumbers(nums []int) []int {
	o := make([]int, len(nums))

	for _, n := range nums {
		o[n-1] = n
	}

	var p int
	for i, n := range o {
		if n != i+1 {
			o[p] = i + 1
			p++
		}
	}

	return o[:p]
}
