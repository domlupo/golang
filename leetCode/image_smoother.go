// https://leetcode.com/problems/image-smoother/

package main

func imageSmoother(img [][]int) [][]int {
	if len(img) == 0 || len(img[0]) == 0 {
		return [][]int{}
	}

	o := make([][]int, len(img))
	for i := range len(img) {
		o[i] = make([]int, len(img[i]))
	}

	for i := range len(img) {
		for j := range len(img[i]) {
			o[i][j] = average(img, i, j)
		}
	}

	return o
}

func average(img [][]int, x, y int) int {
	var count, sum int
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if exists(img, i, j) {
				sum += img[i][j]
				count++
			}
		}
	}
	return sum / count
}

func exists(img [][]int, x, y int) bool {
	return (x >= 0 && x < len(img)) && (y >= 0 && y < len(img[0]))
}
