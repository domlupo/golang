// https://leetcode.com/problems/find-the-duplicate-number/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func findDuplicate(nums []int) int {
	slow, fast := nums[0], nums[0]
	for {
		slow = nums[slow]
		fast = nums[nums[fast]]
		if slow == fast {
			break
		}
	}

	slow = nums[0]
	for fast != slow {
		slow = nums[slow]
		fast = nums[fast]
	}

	return fast
}
