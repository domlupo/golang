// Computer aka non human player engine logic.

package main

import (
	"fmt"
	"math"
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

func (c *computer) updateNodeTargets(s state, nodeCapture bool) {
	if nodeCapture {
		c.updateForNodeCapture(s)
	}

	var opponent player
	switch c.player {
	case bluePlayer:
		opponent = redPlayer
	case redPlayer:
		opponent = bluePlayer
	default:
		panic(fmt.Sprintf("updateNodeTargets: unknown player: %d", c.player))
	}

	gridPos := s.hexagons.maxUnitDifferential(opponent)
	if gridPos.equals(GridPosition{}) {
		return
	}
	p, ok := s.g.getAnyPixelPos(gridPos)
	if !ok {
		panic(fmt.Sprintf("updateNodeTargets: cannot getPixelPos: %s", gridPos.log()))
	}

	for _, n := range s.nodes {
		if n.player != c.player {
			continue
		}

		if n.unitConfig == unguidedGroundMissileConfig {
			selfNode := n
			selfNode.target = position{ground, p}
			continue
		}
	}
}

func (c *computer) updateForNodeCapture(s state) {
	for _, n := range s.nodes {
		if n.player != c.player {
			continue
		}

		selfNodePos := n.gridPos
		selfNode := n
		selfNode.target = position{}
		nearestTargetDistanceSquared := math.MaxInt

		if n.unitConfig == unguidedGroundMissileConfig {
			continue
		}

		for _, n := range s.nodes {
			if n.player == c.player {
				continue
			}

			otherNodePos := n.gridPos
			otherNode := n
			xDelta := (selfNodePos.X - otherNode.gridPos.X)
			yDelta := (selfNodePos.Y - otherNode.gridPos.Y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				var targetPixelPos, ok = s.g.getPixelPos(otherNodePos, hpCenter)
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

			nPixelPos, ok := s.g.getPixelPos(n.gridPos, hpCenter)
			if !ok {
				panic("could not get pixel pos for update computer unit targets")
			}
			xDelta := (u.X - nPixelPos.X)
			yDelta := (u.Y - nPixelPos.Y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				u.updateTarget(globalTarget, position{ground, nPixelPos})
			}
		}
	}
}
