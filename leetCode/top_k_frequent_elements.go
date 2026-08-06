// https://leetcode.com/problems/top-k-frequent-elements/
// N is the length of nums, K is passed k
// time complexity is O(N)
// space complexity is O(N)

package main

func topKFrequent(nums []int, k int) []int {
	counts := make(map[int]int)
	freq := make([][]int, len(nums)+1)

	for _, n := range nums {
		counts[n]++
	}

	for v, c := range counts {
		freq[c] = append(freq[c], v)
	}

	var ans []int
	i := len(freq) - 1
	for i >= 0 && k > 0 {
		ans = append(ans, freq[i]...)
		k -= len(freq[i])
		i--
	}

	return ans
}

/*
Min Heap solution
// time complexity is O(N * log(K))
// space complexity is O(N + K)
func topKFrequent(nums []int, k int) []int {
	counts := make(map[int]int)

	for _, n := range nums {
		counts[n]++
	}

	h := &Heap347{}
	for v, c := range counts {
		heap.Push(h, frequency{val: v, count: c})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	var ans []int
	for k > 0 {
		f := heap.Pop(h).(frequency)
		ans = append(ans, f.val)
		k--
	}

	return ans
}

type frequency struct {
	val, count int
}

type Heap347 []frequency

func (h Heap347) Len() int           { return len(h) }
func (h Heap347) Less(i, j int) bool { return h[i].count < h[j].count }
func (h Heap347) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *Heap347) Push(x any) {
	*h = append(*h, x.(frequency))
}
func (h *Heap347) Pop() any {
	old := *h
	n := len(old) - 1
	x := old[n]
	*h = old[:n]
	return x
}
*/
