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
	var p1 PixelPosition
	var p2 PixelPosition
	for _, diff := range diffs {
		p, ok := s.g.getCenterWithFallback(diff.GridPosition)
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
				selfNode.targets = []position{{land, p1}}
			}
		case guidedAirMissile:
			if !avgUnitPos.equals(PixelPosition{}) {
				selfNode.targets = []position{{land, avgUnitPos}}
			}
		case fighter:
			if diffs[0].diff > 2 && diffs[1].diff > 2 {
				selfNode.targets = []position{}
				first := position{highAir, p1}
				selfNode.targets = append(selfNode.targets, first)
				second := position{highAir, p2}
				selfNode.targets = append(selfNode.targets, second)
				selfNode.targets = append(selfNode.targets, position{layer: bomberConfig.produceLayer, PixelPosition: nP})
			}
		case bomber:
			if diffs[0].diff > 2 && diffs[1].diff > 2 {
				selfNode.targets = []position{}
				selfNode.actions = []action{}
				first := position{highAir, p1}
				selfNode.targets = append(selfNode.targets, first)
				selfNode.actions = append(selfNode.actions, action{pos: first, do: func() {
					u := newUnit(
						bombConfig,
						position{layer: lowAir, PixelPosition: first.PixelPosition},
						c.player,
						[]position{{layer: land, PixelPosition: first.PixelPosition}},
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
				second := position{highAir, p2}
				selfNode.targets = append(selfNode.targets, second)
				selfNode.actions = append(selfNode.actions, action{pos: second, do: func() {
					u := newUnit(
						bombConfig,
						position{layer: lowAir, PixelPosition: second.PixelPosition},
						c.player,
						[]position{{layer: land, PixelPosition: second.PixelPosition}},
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
				selfNode.targets = append(selfNode.targets, position{layer: bomberConfig.produceLayer, PixelPosition: nP})
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
