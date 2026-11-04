// https://leetcode.com/problems/sort-colors/
// N is len(nums)
// time complexity is O(N)
// space complexity is O(1)

package main

// single pass solution
func sortColors(nums []int) {
	i, l, r := 0, 0, len(nums)-1

	for i <= r {
		switch nums[i] {
		case 0:
			nums[l], nums[i] = nums[i], nums[l]
			l++
			i++
		case 2:
			nums[r], nums[i] = nums[i], nums[r]
			r--
		default:
			i++
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

	for i := range nums {
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
