// https://leetcode.com/problems/combination-sum-ii/

package main

func combinationSum2(candidates []int, target int) [][]int {
	var ans [][]int
	var comb []int
	sort40(candidates)
	backtrack2(target, 0, comb, candidates, &ans)
	return ans
}

func backtrack2(
	remain, start int,
	comb, candidates []int,
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
		// skip duplicates. For [1 1 2 5 6 7 10] this line ensures
		// only the first occurred 1 is considered to get [1, 2, 5] and [1, 7]
		// The second occurred 1 is only considered if the first one was also
		// considered. This gets [1, 1, 6].
		if i > start && candidates[i] == candidates[i-1] {
			continue
		}

		comb = append(comb, candidates[i])
		backtrack2(remain-candidates[i], i+1, comb, candidates, ans)
		comb = comb[:len(comb)-1]
	}
}

func sort40(a []int) {
	for i := 0; i < len(a)-1; i++ {
		if i >= 0 && a[i] > a[i+1] {
			a[i], a[i+1] = a[i+1], a[i]
			i -= 2
		}
	}
}
