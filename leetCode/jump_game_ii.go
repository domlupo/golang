// https://leetcode.com/problems/jump-game-ii/

package main

func jump(nums []int) int {
	var maxReachable, jumps, jumpTo int

	for i := 0; i < len(nums)-1; i++ {
		if i+nums[i] > maxReachable {
			maxReachable = i + nums[i]
		}

		if i == jumpTo {
			jumps++
			jumpTo = maxReachable
		}
	}

	return jumps
}
