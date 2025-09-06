package main

import "math"

type computer struct {
	player player
}

func newComputer(p player) *computer {
	c := computer{
		player: p,
	}
	return &c
}

func (c *computer) updateTargetNodes(s state) {
	for nodePos, node := range s.nodes {
		if node.player != c.player {
			continue
		}
		var selfNodePos = nodePos
		var selfNode = node
		selfNode.target = nil
		var nearestTargetDistanceSquared = math.MaxInt

		for nodePos, node := range s.nodes {
			if node.player == c.player {
				continue
			}
			var otherNodePos = nodePos
			var otherNode = node

			var xDelta = (selfNodePos.x - otherNode.gridPos.x)
			var yDelta = (selfNodePos.y - otherNode.gridPos.y)
			var otherTargetDistanceSquared = xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos = getPixelPosition(otherNodePos, innerGridPosition{0, 0})
				selfNode.target = &targetPixelPos
			}

		}
	}

	for unitPos, unit := range s.units {
		if unit.player != c.player {
			continue
		}
		var selfUnitPos = unitPos
		var selfUnit = unit
		selfUnit.globalTarget = nil
		var nearestTargetDistanceSquared = math.MaxInt

		for nodePos, node := range s.nodes {
			if node.player == c.player {
				continue
			}
			var otherNodePos = nodePos
			var otherNodePixelPos = getPixelPosition(otherNodePos, innerGridPosition{0, 0})

			var xDelta = (selfUnitPos.x - otherNodePixelPos.x)
			var yDelta = (selfUnitPos.y - otherNodePixelPos.y)
			var otherTargetDistanceSquared = xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos = getPixelPosition(otherNodePos, innerGridPosition{0, 0})
				selfUnit.globalTarget = &targetPixelPos
			}
		}
	}
}
