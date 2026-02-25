// https://leetcode.com/problems/sum-of-root-to-leaf-binary-numbers/

package main

func sumRootToLeafRecursive(root *TreeNode) int {
	var sum int
	dfs(root, 0, &sum)
	return sum
}

func dfs(n *TreeNode, v int, s *int) {
	if n == nil {
		return
	}
	v = v*2 + n.Val
	if n.Left == nil && n.Right == nil {
		*s += v
	}
	dfs(n.Left, v, s)
	dfs(n.Right, v, s)
}

func sumRootToLeaf(root *TreeNode) int {
	s := []*TreeNode{root}
	v := []int{0}
	var sum int

	for len(s) > 0 {
		curr := s[0]
		val := v[0]
		val *= 2

		s = s[1:]
		v = v[1:]

		if curr.Left != nil {
			s = append(s, curr.Left)
			v = append(v, val+curr.Val)
		}
		if curr.Right != nil {
			s = append(s, curr.Right)
			v = append(v, val+curr.Val)
		}
		if curr.Right == nil && curr.Left == nil {
			sum += val + curr.Val
		}
	}

	return sum
}
