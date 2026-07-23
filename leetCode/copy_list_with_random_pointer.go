// https://leetcode.com/problems/copy-list-with-random-pointer/
// N is the length of list
// time complexity is O(N)
// space complexity is O(N)
// There is a O(1) space solution but it requires modifying input so I skipped it.

package main

type Node138 struct {
	Val    int
	Next   *Node138
	Random *Node138
}

func copyRandomList(head *Node138) *Node138 {
	clones := make(map[*Node138]*Node138)

	curr := head
	for curr != nil {
		clone := &Node138{Val: curr.Val}
		clones[curr] = clone
		curr = curr.Next
	}

	curr = head
	for curr != nil {
		clone := clones[curr]
		clone.Next = clones[curr.Next]
		clone.Random = clones[curr.Random]
		curr = curr.Next
	}

	return clones[head]
}
