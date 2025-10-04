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
	for _, n := range s.nodes {
		if n.player != c.player {
			continue
		}

		selfNodePos := n.gridPos
		selfNode := n
		selfNode.target = position{}
		nearestTargetDistanceSquared := math.MaxInt

		for _, n := range s.nodes {
			if n.player == c.player {
				continue
			}

			otherNodePos := n.gridPos
			otherNode := n
			xDelta := (selfNodePos.x - otherNode.gridPos.x)
			yDelta := (selfNodePos.y - otherNode.gridPos.y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos, ok = otherNodePos.getPixelPos(center, s.hexConfig)
				if !ok {
					panic("could not get pixel pos for update computer node targets")
				}

				selfNode.target = position{ground, targetPixelPos}
			}

		}
	}

	for _, u := range s.units {
		if u.staticTarget {
			continue
		}

		if u.player != c.player {
			continue
		}

		nearestTargetDistanceSquared := math.MaxInt

		for _, n := range s.nodes {
			if n.player == c.player {
				continue
			}

			nPixelPos, ok := n.gridPos.getPixelPos(center, s.hexConfig)
			if !ok {
				panic("could not get pixel pos for update computer unit targets")
			}
			xDelta := (u.x - nPixelPos.x)
			yDelta := (u.y - nPixelPos.y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				u.updateTarget(globalTarget, position{ground, nPixelPos})
			}
		}
	}
}
