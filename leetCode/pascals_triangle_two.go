// https://leetcode.com/problems/pascals-triangle-ii/

package main

// optimal due to computation trick next := prev * (rowIndex - col + 1) / col
func getRow(rowIndex int) []int {
	t := []int{1}
	prev := 1
	for col := 1; col <= rowIndex; col++ {
		next := prev * (rowIndex - col + 1) / col
		t = append(t, next)
		prev = next
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
