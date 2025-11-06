// Computer aka non human player engine logic.

package main

import (
	"fmt"
	"math"
)

type computer struct {
	player                player
	validUnitDifferential int
}

func newComputer(p player) *computer {
	c := computer{
		player:                p,
		validUnitDifferential: 2,
	}
	return &c
}

func (c *computer) updateNodeTargets(s *state, nodeCapture bool, avgUnitPos PixelPosition) {
	if nodeCapture {
		c.updateForNodeCapture(*s)
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

	diffs := s.hexagons.maxUnitDifferential(opponent, land)
	validDiffsIndex := -1
	var positions []PixelPosition
	for _, d := range diffs {
		if d.diff > c.validUnitDifferential {
			validDiffsIndex++
		}

		p, ok := s.g.getCenterWithFallback(d.GridPosition)
		if !ok {
			panic("updateNodeTargets: cannot get center with fallback")
		}
		positions = append(positions, p)
	}

	for _, n := range s.nodes {
		if n.player != c.player {
			continue
		}

		if n.unitConfig.name == infantry {
			continue
		}

		nP, ok := s.g.getCenterWithFallback(n.gridPos)
		if !ok {
			panic(fmt.Sprintf("updateNodeTargets: cannot get node pixel position %s", nP.log()))
		}
		selfNode := n
		selfNode.targets = []position{{land, nP}}
		switch n.unitConfig.name {
		case unguidedGroundMissile:
			if diffs[0].diff > 1 {
				selfNode.targets = []position{{land, positions[0]}}
			}
		case guidedAirMissile:
			if !avgUnitPos.equals(PixelPosition{}) {
				selfNode.targets = []position{{land, avgUnitPos}}
			}
		case fighter:
			if fighterConfig.bombs <= validDiffsIndex+1 {
				c.planBombs(s, selfNode, positions[:fighterConfig.bombs])
			}
		case bomber:
			if bomberConfig.bombs <= validDiffsIndex+1 {
				c.planBombs(s, selfNode, positions[:bomberConfig.bombs])
			}
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

				selfNode.targets = []position{{land, targetPixelPos}}
			}

		}
	}

	for u := range s.units() {
		if u.staticTarget ||
			u.unitConfig.name == guidedAirMissile ||
			u.unitConfig.name == bomber ||
			u.unitConfig.name == fighter {
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

func (c *computer) planBombs(s *state, n *node, targets []PixelPosition) {
	n.targets = []position{}
	n.actions = []action{}

	for _, t := range targets {
		p := position{highAir, t}
		n.targets = append(n.targets, p)
		n.actions = append(n.actions, action{pos: p, do: func() {
			u := newUnit(
				bombConfig,
				position{layer: lowAir, PixelPosition: p.PixelPosition},
				c.player,
				[]position{{layer: land, PixelPosition: p.PixelPosition}},
				[]action{},
				GridPosition{},
				s.unitSpawnTotal,
				false)
			s.spawn(&u)
			s.events = append(s.events, event{
				step: s.step + bombConfig.explodeStepDelay,
				do: func() {
					s.explode(&u)
				},
			})
		}})
	}

	nP, ok := s.g.getCenterWithFallback(n.gridPos)
	if !ok {
		panic(fmt.Sprintf("planBombs: cannot get node pixel position %s", nP.log()))
	}
	n.targets = append(n.targets, position{layer: bomberConfig.produceLayer, PixelPosition: nP})
}
