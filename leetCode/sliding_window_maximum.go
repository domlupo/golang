// https://leetcode.com/problems/sliding-window-maximum/
// N is the length of nums, K is the length of window
// time complexity is O(N)
// space complexity is O(K)

package main

func maxSlidingWindow(nums []int, k int) []int {
	var dq, ans []int

	for i, n := range nums {
		// remove out of window values from front of deque
		for len(dq) > 0 && dq[0] <= i-k {
			dq = dq[1:]
		}

		// remove existing values lower/equal to current value from back of deque
		for len(dq) > 0 && nums[dq[len(dq)-1]] <= n {
			dq = dq[:len(dq)-1]
		}

		// add new value to back of deque
		dq = append(dq, i)

		// create answer once window formed
		if i >= k-1 {
			ans = append(ans, nums[dq[0]])
		}
	}

	return ans
}
