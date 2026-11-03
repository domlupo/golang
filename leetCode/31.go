// https://leetcode.com/problems/next-permutation/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

func nextPermutation(nums []int) {
	l := -1

	// 1. find right most index where nums[x-1] < nums[x].
	for i := 1; i < len(nums); i++ {
		if nums[i-1] < nums[i] {
			l = i - 1
		}
	}

	// if this index does not exist, nums is sorted in descending and
	// reverse entire array to get the next permutation
	if l == -1 {
		reverseToEnd(nums, 0)
		return
	}

	// 2. find index most right and greater than the value of nums[l]
	r := l + 1
	for i := r; i < len(nums); i++ {
		// find the right most to handle cases like below
		// 1 3 3 3 3 -> 3 1 3 3 3
		// l       r
		// 3 3 3 3 1
		// reverse l+1 to end
		// 3 1 3 3 3
		// instead of incorrect
		// 1 3 3 3 3 -> 3 3 3 3 1
		// l r
		// 3 1 3 3 3
		// reverse l+1 to end
		// 3 3 3 3 1
		//
		// find greater than value to handle cases like below
		// 2 5 4 3 2 -> 3 2 2 4 5
		// l     r
		// 3 5 4 2 2
		// reverse l+1 to end
		// 3 2 2 4 5
		// instead of incorrect
		// 2 5 4 3 2 -> 2 2 3 4 5
		// l       r
		// 2 5 4 3 2
		// reverse l+1 to end
		// 2 2 3 4 5
		if nums[i] > nums[l] {
			r = i
		}
	}

	nums[l], nums[r] = nums[r], nums[l]
	reverseToEnd(nums, l+1)
}

func reverseToEnd(nums []int, l int) {
	r := len(nums) - 1
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}

// remember how to get next for 1 6 5 4 3 2 -> 2 1 3 4 5 6
// 1 6 5 4 3 2
// 1. find right most index where nums[x-1] < nums[x].
// l := 0
// 2. find index most right and greater than the value of nums[l]
// r := 5
// 3. swap nums[l] and nums[r]
// 2 6 5 4 3 1
// 4. reverse the numbers between l+1 and len(nums)-1
// 2 1 3 4 5

// all next permutations for 1 2 3 4
// 1 2 3 4
// 1 2 4 3
// 1 3 2 4
// 1 3 4 2
// 1 4 2 3
// 1 4 3 2

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
// 4 3 2 1
