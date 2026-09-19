// https://leetcode.com/problems/climbing-stairs/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(1)

package main

// iterative (best space and time complexity)

func climbStairs(n int) int {
	prev, curr := 1, 1
	for range n {
		prev, curr = curr, prev+curr
	}
	return prev
}

// dynamic programming

var ways map[int]int = map[int]int{1: 1, 2: 2}

func climbStairsDP(n int) int {
	if w, ok := ways[n]; ok {
		return w
	} else {
		ways[n] = climbStairsDP(n-1) + climbStairsDP(n-2)
		return ways[n]
	}
}
