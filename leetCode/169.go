// https://leetcode.com/problems/majority-element/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

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
