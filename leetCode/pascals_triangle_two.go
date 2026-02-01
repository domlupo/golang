// https://leetcode.com/problems/pascals-triangle-ii/

package main

// optimal due to computation trick n := p * (rowIndex - k + 1) / k
func getRow(rowIndex int) []int {
	t := []int{1}
	p := 1
	for k := 1; k <= rowIndex; k++ {
		n := p * (rowIndex - k + 1) / k
		t = append(t, n)
		p = n
	}
	return t
}

// not optimal
func getRowNotOptimal(rowIndex int) []int {
	t := []int{1}
	for rowIndex != 0 {
		rowIndex--
		p := t

		t = []int{1}
		for i := 1; i < len(p); i++ {
			t = append(t, p[i-1]+p[i])
		}
		t = append(t, 1)
	}

	return t
}
