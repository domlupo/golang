// https://leetcode.com/problems/combination-sum-iii/
// K is k
// time complexity is O(K×C(9,K))
// space complexity is O(K)

package main

func combinationSum3(k int, n int) [][]int {
	var ans [][]int
	var comb []int
	backtrack216(&ans, comb, 1, 0, k, n)
	return ans
}

func backtrack216(ans *[][]int, comb []int, start, sum, k, n int) {
	if sum == n && len(comb) == k {
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	} else if sum >= n || len(comb) >= k {
		return
	}

	for i := start; i < 10; i++ {
		comb = append(comb, i)
		backtrack216(ans, comb, i+1, sum+i, k, n)
		comb = comb[:len(comb)-1]
	}
}
