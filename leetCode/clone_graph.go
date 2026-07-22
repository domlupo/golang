// https://leetcode.com/problems/clone-graph/
// N is the number of nodes, M is the number of edges
// time complexity is O(N + M)
// space complexity is O(N)

package main

type Node133 struct {
	Val       int
	Neighbors []*Node133
}

func cloneGraph(node *Node133) *Node133 {
	if node == nil {
		return nil
	}

	clones := map[*Node133]*Node133{
		node: &Node133{Val: node.Val},
	}

	q := []*Node133{node}
	for len(q) != 0 {
		n := q[0]
		q = q[1:]
		clone := clones[n]
		for _, neighbor := range n.Neighbors {
			neighborClone, neighborCloneExists := clones[neighbor]
			if !neighborCloneExists {
				neighborClone = &Node133{Val: neighbor.Val}
				clones[neighbor] = neighborClone
				q = append(q, neighbor)
			}
			clone.Neighbors = append(clone.Neighbors, neighborClone)
		}
	}

	return clones[node]
}
