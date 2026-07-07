// https://leetcode.com/problems/permutations-ii/

package main

func permuteUnique(nums []int) [][]int {
	var ans [][]int
	var perm []int
	counter := make(map[int]int)
	for _, n := range nums {
		counter[n]++
	}
	backtrack4(perm, &ans, counter, len(nums))
	return ans
}

func backtrack4(perm []int, ans *[][]int, counter map[int]int, size int) {
	if len(perm) == size {
		deepCopy := make([]int, size)
		copy(deepCopy, perm)
		*ans = append(*ans, deepCopy)
		return
	}

	for num, count := range counter {
		if count == 0 {
			continue
		}
		perm = append(perm, num)
		counter[num]--
		backtrack4(perm, ans, counter, size)
		counter[num]++
		perm = perm[:len(perm)-1]
	}
}
