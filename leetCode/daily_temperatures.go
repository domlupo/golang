// https://leetcode.com/problems/daily-temperatures/
// N is the length of temperatures
// time complexity is O(N)
// space complexity is O(1) due to using returned array

package main

func dailyTemperatures(temperatures []int) []int {
	var hottest int
	ans := make([]int, len(temperatures))

	for i := len(temperatures) - 1; i >= 0; i-- {
		t := temperatures[i]
		hottest = max(hottest, t)

		if t == hottest {
			ans[i] = 0
			continue
		}

		next := i + 1
		for t >= temperatures[next] {
			next += ans[next]
		}
		ans[i] = next - i
	}

	return ans
}

/*
// N is the length of temperatures
// time complexity is O(N)
// space complexity is O(N)

func dailyTemperatures(temperatures []int) []int {
	ans := make([]int, len(temperatures))
	var stack []int

	for i, t := range temperatures {
		for len(stack) > 0 && temperatures[stack[len(stack)-1]] < t {
			prev := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ans[prev] = i - prev
		}

		stack = append(stack, i)
	}

	return ans
}

*/
