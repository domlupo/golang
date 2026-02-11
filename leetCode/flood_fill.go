// https://leetcode.com/problems/flood-fill/

package main

type pos struct {
	r int
	c int
}

func floodFill(image [][]int, sr int, sc int, color int) [][]int {
	startColor := image[sr][sc]
	q := []pos{{r: sr, c: sc}}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		if image[p.r][p.c] != color {
			image[p.r][p.c] = color

			positions := []pos{
				{r: p.r - 1, c: p.c},
				{r: p.r + 1, c: p.c},
				{r: p.r, c: p.c - 1},
				{r: p.r, c: p.c + 1},
			}

			for _, position := range positions {
				if exists1(image, position) && image[position.r][position.c] == startColor {
					q = append(q, position)
				}
			}
		}
	}

	return image
}

func exists1(image [][]int, p pos) bool {
	return (p.r >= 0 && p.r < len(image)) && (p.c >= 0 && p.c < len(image[0]))
}
