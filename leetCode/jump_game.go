// https://leetcode.com/problems/jump-game/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func canJump(nums []int) bool {
	var maxReach int

	for ni, n := range nums {
		if ni <= maxReach {
			maxReach = max(maxReach, ni+n)
		}
	}

	return maxReach >= len(nums)-1
}
