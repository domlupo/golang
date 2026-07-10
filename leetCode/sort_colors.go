// https://leetcode.com/problems/sort-colors/
// N is len(nums)
// time complexity is O(N)
// space complexity is O(1)

package main

// single pass solution
func sortColors(nums []int) {
	var l, curr int
	r := len(nums) - 1

	for curr <= r {
		if nums[curr] == 0 {
			nums[l], nums[curr] = nums[curr], nums[l]
			l++
			curr++
		} else if nums[curr] == 2 {
			nums[r], nums[curr] = nums[curr], nums[r]
			r--
		} else {
			curr++
		}
	}
}

// counting sort, not single pass
func sortColors1(nums []int) {
	var zeros, ones int
	for _, n := range nums {
		switch n {
		case 0:
			zeros++
		case 1:
			ones++
		}
	}

	for i := 0; i < len(nums); i++ {
		if zeros > 0 {
			nums[i] = 0
			zeros--
		} else if ones > 0 {
			nums[i] = 1
			ones--
		} else {
			nums[i] = 2
		}
	}
}
