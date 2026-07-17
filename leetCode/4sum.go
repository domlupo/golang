// https://leetcode.com/problems/4sum/
// N is the length of nums
// time complexity is O(N^3) or O(N^(k-1)) where K is for K sum
// space complexity is O(N)

package main

import "slices"

func fourSum(nums []int, target int) [][]int {
	slices.Sort(nums)
	return kSum(nums, target, 0, 4)
}

func kSum(nums []int, target, start, k int) [][]int {
	var ans [][]int

	if k == 2 {
		return twoSum(nums, target, start)
	}

	for i := start; i < len(nums); i++ {
		// avoid duplicates
		if i == start || nums[i] != nums[i-1] {
			for _, comb := range kSum(nums, target-nums[i], i+1, k-1) {
				ans = append(ans, append([]int{nums[i]}, comb...))
			}
		}
	}

	return ans
}

func twoSum(nums []int, target, l int) [][]int {
	var ans [][]int
	r := len(nums) - 1
	for l < r {
		sum := nums[l] + nums[r]
		if sum == target {
			ans = append(ans, []int{nums[l], nums[r]})
			l++
			r--
			for l < r && nums[l] == nums[l-1] {
				// avoid duplicates
				l++
			}
		} else if sum < target {
			l++
		} else {
			r--
		}
	}

	return ans
}

/*

four sum instead of k sum

func fourSum(nums []int, target int) [][]int {
	var ans [][]int
	slices.Sort(nums)

	for i := 0; i < len(nums)-3; i++ {
		if i > 0 && nums[i-1] == nums[i] {
			// avoid duplicates
			continue
		}
		for j := i + 1; j < len(nums)-2; j++ {
			if j > i+1 && nums[j-1] == nums[j] {
				// avoid duplicates
				continue
			}

			l := j + 1
			r := len(nums) - 1
			for l < r {
				sum := nums[i] + nums[j] + nums[l] + nums[r]
				if sum == target {
					ans = append(ans, []int{nums[i], nums[j], nums[l], nums[r]})
					l++
					r--
					for l < r && nums[l] == nums[l-1] {
						// avoid duplicates
						l++
					}
				} else if sum < target {
					l++
				} else {
					r--
				}
			}
		}
	}

	return ans
}
*/
