// https://leetcode.com/problems/can-place-flowers/
// N is the length of flowerbed
// time complexity is O(N)
// space complexity is O(1)

package main

func canPlaceFlowers(flowerbed []int, n int) bool {
	var plants int
	for i := 0; i < len(flowerbed); i++ {
		l := i == 0 || flowerbed[i-1] == 0
		r := i == len(flowerbed)-1 || flowerbed[i+1] == 0
		if l && r && flowerbed[i] == 0 {
			plants++
			i++
		}
	}
	return plants >= n
}
