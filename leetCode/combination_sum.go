// https://leetcode.com/problems/combination-sum/

package main

func combinationSum(candidates []int, target int) [][]int {
	var ans [][]int
	var comb []int
	backtrack1(target, comb, 0, candidates, &ans)
	return ans
}

func backtrack1(
	remain int,
	comb []int,
	start int,
	candidates []int,
	ans *[][]int,
) {
	if remain == 0 {
		newComb := make([]int, len(comb))
		copy(newComb, comb)
		*ans = append(*ans, newComb)
		return
	} else if remain < 0 {
		return
	}

	for i := start; i < len(candidates); i++ {
		comb = append(comb, candidates[i])
		backtrack1(remain-candidates[i], comb, i, candidates, ans)
		comb = comb[:len(comb)-1]
	}
}
