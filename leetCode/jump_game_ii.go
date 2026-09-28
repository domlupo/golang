// https://leetcode.com/problems/jump-game-ii/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func jump(nums []int) int {
	var jumps, maxPos, pos int

	for i, n := range nums {
		maxPos = max(maxPos, i+n)
		if pos >= len(nums)-1 {
			break
		} else if i == pos {
			jumps++
			pos = maxPos
		}
	}

	return jumps
}
