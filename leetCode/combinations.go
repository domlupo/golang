// https://leetcode.com/problems/two-sum/
// time complexity is very hard to calculate since backtracking is used.
// space complexity is O(k) where k is length of combination due to recursion call stack.

package main

func combine(n int, k int) [][]int {
	var ans [][]int
	var comb []int
	backtrack77(&ans, comb, k, 1, n)
	return ans
}

func backtrack77(ans *[][]int, comb []int, k, start, n int) {
	if len(comb) == k {
		deepCopy := make([]int, k)
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	for i := start; i <= n; i++ {
		comb = append(comb, i)
		backtrack77(ans, comb, k, i+1, n)
		comb = comb[:len(comb)-1]
	}
}
