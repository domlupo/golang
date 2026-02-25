// https://leetcode.com/problems/partition-array-into-three-parts-with-equal-sum/

package main

func canThreePartsEqualSum(arr []int) bool {
	var total int
	for _, n := range arr {
		total += n
	}
	if total%3 != 0 {
		return false
	}
	target := total / 3

	var sum, sums int
	for _, n := range arr {
		sum += n
		if sum == target && sums < 2 {
			sum = 0
			sums += 1
		} else if sum == target && sums == 2 {
			sums += 1
		}
	}

	return sums == 3 && sum == target
}
