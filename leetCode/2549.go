// https://leetcode.com/problems/count-distinct-numbers-on-board/
// time complexity is O(1)
// space complexity is O(1)

package main

func distinctIntegers(n int) int {
	if n == 1 {
		return 1
	} else {
		return n - 1
	}
}
