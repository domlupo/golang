// https://leetcode.com/problems/number-of-lines-to-write-string/

package main

func numberOfLines(widths []int, s string) []int {
	lines := 1
	var pixels int
	for _, c := range s {
		i := int(c - 'a')
		charPixels := widths[i]
		if pixels+charPixels > 100 {
			lines++
			pixels = 0
		}
		pixels += charPixels
	}

	return []int{lines, pixels}
}
