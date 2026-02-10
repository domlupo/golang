// https://leetcode.com/problems/can-place-flowers/

package main

func canPlaceFlowers(flowerbed []int, n int) bool {
	for i := 0; i < len(flowerbed); i++ {
		l := i == 0 || flowerbed[i-1] == 0
		r := i == len(flowerbed)-1 || flowerbed[i+1] == 0

		if l && r && flowerbed[i] == 0 {
			n--
			flowerbed[i] = 1
		}
	}

	return n <= 0
}
