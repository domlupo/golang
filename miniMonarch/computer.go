// computer aka non human player engine logic
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

		selfNodePos := nodePos
		selfNode := node
		selfNode.target = nil
		nearestTargetDistanceSquared := math.MaxInt

		for nodePos, node := range s.nodes {
			if node.player == c.player {
				continue
			}

			otherNodePos := nodePos
			otherNode := node
			xDelta := (selfNodePos.x - otherNode.gridPos.x)
			yDelta := (selfNodePos.y - otherNode.gridPos.y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos = getPixelPosition(otherNodePos, innerGridPosition{0, 0})
				selfNode.target = &position{ground, targetPixelPos}
			}

		}
	}

	for unitPos, unit := range s.units {
		if unit.staticTarget {
			continue
		}

		if unit.player != c.player {
			continue
		}

		selfUnitPos := unitPos
		selfUnit := unit
		nearestTargetDistanceSquared := math.MaxInt

		for nodePos, node := range s.nodes {
			if node.player == c.player {
				continue
			}

			otherNodePos := nodePos
			otherNodePixelPos := getPixelPosition(otherNodePos, innerGridPosition{0, 0})
			xDelta := (selfUnitPos.x - otherNodePixelPos.x)
			yDelta := (selfUnitPos.y - otherNodePixelPos.y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos = getPixelPosition(otherNodePos, innerGridPosition{0, 0})
				selfUnit.updateTarget(globalTarget, &position{ground, targetPixelPos})
			}
		}
	}
}
