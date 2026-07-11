// https://leetcode.com/problems/remove-duplicates-from-sorted-array-ii/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func removeDuplicates80(nums []int) int {
	var count, swap int
	n := nums[0]

	for i := 0; i < len(nums); i++ {
		if count < 2 && nums[i] == n {
			nums[i], nums[swap] = nums[swap], nums[i]
			swap++
			count++
		} else if nums[i] > n {
			n = nums[i]
			nums[i], nums[swap] = nums[swap], nums[i]
			swap++
			count = 1
		}
	}

	return swap
}
