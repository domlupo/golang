// https://leetcode.com/problems/3sum
// N is the length of nums
// time complexity is O(N * N)
// space complexity is O(log(N))

package main

import "slices"

func threeSum(nums []int) [][]int {
	var ans [][]int
	slices.Sort(nums)

	for i := 0; i < len(nums)-2; i++ {
		if i > 0 && nums[i-1] == nums[i] {
			// avoid duplicates
			continue
		}

		l := i + 1
		r := len(nums) - 1
		for l < r {
			sum := nums[i] + nums[l] + nums[r]
			if sum == 0 {
				ans = append(ans, []int{nums[i], nums[l], nums[r]})
				l++
				r--
				// avoid duplicates
				for l < r && nums[l] == nums[l-1] {
					l++
				}
			} else if sum < 0 {
				l++
			} else {
				r--
			}
		}
	}

	return ans
}
