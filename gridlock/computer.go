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

func (c *computer) updateNodeTargets(s state, nodeCapture bool, avgUnitPos PixelPosition) {
	if nodeCapture {
		c.updateForNodeCapture(s)
	}

	var opponent player
	switch c.player {
	case pOne:
		opponent = pTwo
	case pTwo:
		opponent = pOne
	default:
		panic(fmt.Sprintf("updateNodeTargets: unknown player: %d", c.player))
	}

	maxDiffs := s.hexagons.maxUnitDifferential(opponent, land)
	var p1 PixelPosition
	var p2 PixelPosition
	for _, pos := range maxDiffs {
		p, ok := s.g.getPixelPos(pos, hpCenter)
		if !ok {
			continue
		}
		if p1.equals(PixelPosition{}) {
			p1 = p
		} else if p2.equals(PixelPosition{}) {
			p2 = p
		} else {
			break
		}
	}
	for _, n := range s.nodes {
		if n.player != c.player {
			continue
		}

		selfNode := n
		switch n.unitConfig.name {
		case infantry:
			continue
		case unguidedGroundMissile:
			selfNode.targets = []position{{land, p1}}
		case guidedAirMissile:
			selfNode.targets = []position{{land, avgUnitPos}}
		case bomber:
			nodeP, ok := s.g.getPixelPos(n.gridPos, hpCenter)
			if !ok {
				panic(fmt.Sprintf("updateNodeTargets: cannot getPixelPos for bomber: %s", n.gridPos.log()))
			}
			selfNode.targets = []position{}
			selfNode.targets = append(selfNode.targets, position{land, p1})
			selfNode.targets = append(selfNode.targets, position{land, p2})
			selfNode.targets = append(selfNode.targets, position{layer: bomberConfig.produceLayer, PixelPosition: nodeP})
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
		selfNode.targets = []position{}
		nearestTargetDistanceSquared := math.MaxInt

		if n.unitConfig.name == unguidedGroundMissile || n.unitConfig.name == guidedAirMissile {
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

				selfNode.targets = []position{position{land, targetPixelPos}}
			}

		}
	}

	for u := range s.units() {
		if u.staticTarget || u.unitConfig.name == guidedAirMissile || u.unitConfig.name == bomber {
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
				u.updateTarget(globalTarget, position{land, nPixelPos})
			}
		}
	}
}
