// https://leetcode.com/problems/range-sum-query-immutable/

package main

type NumArray struct {
	prefixSums []int
}

func Constructor(nums []int) NumArray {
	prefixSums := []int{0}
	var s int
	for _, n := range nums {
		s += n
		prefixSums = append(prefixSums, s)
	}
	return NumArray{
		prefixSums: prefixSums,
	}
}

func (this *NumArray) SumRange(left int, right int) int {
	return this.prefixSums[right+1] - this.prefixSums[left]
}
