// https://leetcode.com/problems/next-permutation/

package main

func nextPermutation(nums []int) {
	i := -1
	for j := len(nums) - 2; j >= 0; j-- {
		if nums[j] < nums[j+1] {
			i = j
			break
		}
	}

	// handle special case where array is in descending order
	// it must be reversed to ascending order
	if i == -1 {
		reverse1(nums, 0, len(nums)-1)
		return
	}

	// find the number to the right that is higher than the swap
	// value by the least amount
	swap := -1
	for j := len(nums) - 1; j >= i+1; j-- {
		if nums[j] > nums[i] {
			swap = j
			break
		}
	}

	nums[i], nums[swap] = nums[swap], nums[i]
	reverse1(nums, i+1, len(nums)-1)
}

// parameters: ([]int{1 2 3 4 5 6}, 0, 5
// modifies array: []int{6 5 4 3 2 1}
func reverse1(nums []int, l, r int) {
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}

// high level steps

// find the right most index where (x - 1) < x
// if this does not exist, reverse entire array since
// it is in descending order.

// if index does exist, take the value at (x-1) and find a
// value to the right that is larger by the smallest amount.
// the rest of the array will be sorted in descending.

// swap the value at (x-1) with the other val
// sort the array by reversing values from l++ and r--

// 1 2 3  | 1
// 1 3 2  | 0
// 2 1 3  | 1
// 2 3 1  | 0
// 3 1 2  | 1
// 3 2 1

// 1 2 3 4  | 2
// 1 2 4 3  | 1
// 1 3 2 4  | 2
// 1 3 4 2  | 1
// 1 4 2 3  | 2
// 1 4 3 2  | 0

// 2 1 3 4
// 2 1 4 3
// 2 3 1 4
// 2 3 4 1
// 2 4 1 3
// 2 4 3 1

// 3 1 2 4
// 3 1 4 2
// 3 2 1 4
// 3 2 4 1
// 3 4 1 2
// 3 4 2 1

// 4 1 2 3
// 4 1 3 2
// 4 2 1 3
// 4 2 3 1
// 4 3 1 2
// 4 3 2
