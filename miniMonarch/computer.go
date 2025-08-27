package main

import (
	"slices"
)

type computer struct {
	player player
}

func newComputer(p player) *computer {
	c := computer{
		player: p,
	}
	return &c
}

func (c *computer) updateActions(s state) {
	// TODO: Determine behavior here. This does not handle the edge case where multiple computers are targeting the same node
	// and it changs ownership. So the computer will "incorrectly" target an owned node before capturing all neutrals.
	nonCapturedTargets := make(map[gridPosition]struct{}, 0)
	for _, node := range s.nodes {
		if node.player == c.player && node.target != nil && s.nodes[gridPosition{node.target.x, node.target.y}].player != c.player {
			nonCapturedTargets[gridPosition{node.target.x, node.target.y}] = struct{}{}
		}
	}

	if len(nonCapturedTargets) > 0 {
		return
	}

	neutralNodes := make(map[gridPosition]struct{})
	opponentNodes := make(map[gridPosition]struct{})
	for nodeGridPos, node := range s.nodes {
		if node.player == neutralPlayer {
			neutralNodes[nodeGridPos] = struct{}{}
		} else if node.player != c.player {
			opponentNodes[nodeGridPos] = struct{}{}
		}
	}

	// TODO: don't hard code target order
	var targetOrder []gridPosition
	switch c.player {
	case bluePlayer:
		targetOrder = []gridPosition{
			// nearest 3
			{5, 1},
			{5, 5},
			{1, 5},
			// next nearest 5
			{10, 1},
			{10, 5},
			{10, 10},
			{5, 10},
			{1, 10},
			// next nearest 7
			{15, 1},
			{15, 5},
			{15, 10},
			{15, 15},
			{10, 15},
			{5, 15},
			{1, 15},
			// next right 4
			{20, 1},
			{20, 5},
			{20, 10},
			{20, 15},
			// next right 4
			{25, 1},
			{25, 5},
			{25, 10},
			{25, 15},
			// next right 4
			{30, 1},
			{30, 5},
			{30, 10},
			{30, 15},
		}
	case redPlayer:
		targetOrder = []gridPosition{
			// nearest 3
			{25, 15},
			{30, 10},
			{25, 10},
			// next nearest 5
			{20, 15},
			{20, 10},
			{20, 5},
			{25, 5},
			{30, 5},
			// next nearest 7
			{15, 15},
			{15, 10},
			{15, 5},
			{15, 1},
			{20, 1},
			{25, 1},
			{30, 1},
			// next 4 left
			{10, 15},
			{10, 10},
			{10, 5},
			{10, 1},
			// next 4 left
			{5, 15},
			{5, 10},
			{5, 5},
			{5, 1},
			// next 4 left
			{1, 15},
			{1, 10},
			{1, 5},
			{1, 1},
		}
	default:
		panic("computer must be a non neutral player")
	}
	// TODO: optimize, very inefficient
	var target gridPosition
	for _, targetGridPos := range slices.Backward(targetOrder) {
		if _, hasNode := opponentNodes[targetGridPos]; hasNode {
			target = targetGridPos
		}
	}
	for _, targetGridPos := range slices.Backward(targetOrder) {
		if _, hasNode := neutralNodes[targetGridPos]; hasNode {
			target = targetGridPos
		}
	}

	for nodeGridPos, node := range s.nodes {
		if node.player == c.player {
			node.target = &target
			s.nodes[nodeGridPos] = node
		}
	}

	for unitGridPos, unit := range s.units {
		if unit.player == c.player {
			unit.target = &target
			s.units[unitGridPos] = unit
		}
	}
}
