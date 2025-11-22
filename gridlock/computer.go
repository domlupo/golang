// Computer aka non human player engine logic.

package main

import (
	"fmt"
	"math"
)

// Configuration shared by all computers
type computersConfig struct {
	count               int
	minUnitDifferential int
}

type computers []computer

type computer struct {
	s      *state
	player player
	computersConfig
}

func newComputersConfig() computersConfig {
	c := computersConfig{
		count:               2,
		minUnitDifferential: 2,
	}
	return c
}

func newComputers(s *state) computers {
	var computers = make(computers, 0)
	c := newComputersConfig()
	for i := range c.count {
		computers = append(computers, newComputer(s, i))
	}
	computers.updateNodetargets(true, PixelPosition{})
	return computers
}

func newComputer(s *state, i int) computer {
	var p player
	switch i {
	case 0:
		p = pOne
	case 1:
		p = pTwo
	default:
		panic("newComputer: only two computers are supported")
	}

	c := computer{
		s:               s,
		player:          p,
		computersConfig: newComputersConfig(),
	}
	return c
}

func (computers computers) active() []player {
	a := []player{}
	for _, c := range computers {
		a = append(a, c.player)
	}
	return a
}

func (computers *computers) updateNodetargets(nodeCapture bool, avgUnitPos PixelPosition) {
	for _, c := range *computers {
		if nodeCapture {
			c.updateForNodeCapture()
		}
		c.updateNodeTargets(avgUnitPos)
	}
}

func (c *computer) updateNodeTargets(avgUnitPos PixelPosition) {
	var opponent player
	switch c.player {
	case pOne:
		opponent = pTwo
	case pTwo:
		opponent = pOne
	default:
		panic(fmt.Sprintf("updateNodeTargets: unknown player: %d", c.player))
	}

	diffs := c.s.hexagons.maxUnitDifferential(opponent, land)
	validDiffsIndex := -1
	var positions []PixelPosition
	for _, d := range diffs {
		if d.diff > c.minUnitDifferential {
			validDiffsIndex++
		}

		p, ok := c.s.g.getCenterWithFallback(d.GridPosition)
		if !ok {
			panic("updateNodeTargets: cannot get center with fallback")
		}
		positions = append(positions, p)
	}

	for _, n := range c.s.nodes {
		if n.player != c.player {
			continue
		}

		if n.unitConfig.name == infantry {
			continue
		}

		nP, ok := c.s.g.getCenterWithFallback(n.gridPos)
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
		case missile:
			if !avgUnitPos.equals(PixelPosition{}) {
				selfNode.targets = []position{{land, avgUnitPos}}
			}
		case fighter:
			if fighterConfig.bombCountMax <= validDiffsIndex+1 {
				c.planBombs(selfNode, positions[:fighterConfig.bombCountMax])
			}
		case bomber:
			if bomberConfig.bombCountMax <= validDiffsIndex+1 {
				c.planBombs(selfNode, positions[:bomberConfig.bombCountMax])
			}
		}
	}
}

func (c *computer) updateForNodeCapture() {
	for _, n := range c.s.nodes {
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

		for _, n := range c.s.nodes {
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
				var targetPixelPos, ok = c.s.g.getPixelPos(otherNodePos, hpCenter)
				if !ok {
					panic("could not get pixel pos for update computer node targets")
				}

				selfNode.targets = []position{{land, targetPixelPos}}
			}

		}
	}

	for u := range c.s.units() {
		if u.staticTarget ||
			u.name == guidedAirMissile ||
			u.name == bomber ||
			u.name == fighter {
			continue
		}

		if u.player != c.player {
			continue
		}

		nearestTargetDistanceSquared := math.MaxInt

		for _, n := range c.s.nodes {
			if n.player == c.player {
				continue
			}

			nPixelPos, ok := c.s.g.getPixelPos(n.gridPos, hpCenter)
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

func (c *computer) planBombs(n *node, targets []PixelPosition) {
	n.targets = []position{}
	n.actions = []action{}

	for _, t := range targets {
		p := position{aircraft, t}
		n.targets = append(n.targets, p)
		n.actions = append(n.actions, action{pos: p, do: func() {
			u := newUnit(
				bombConfig,
				position{layer: ordnance, PixelPosition: p.PixelPosition},
				c.player,
				[]position{{layer: land, PixelPosition: p.PixelPosition}},
				[]action{},
				GridPosition{},
				false)
			// TODO: Update this to a spawn request since it can panic
			c.s.spawn(&u)
			c.s.events = append(c.s.events, event{
				step: c.s.step + bombConfig.explodeStepDelay,
				do: func() {
					c.s.explode(&u)
				},
			})
		}})
	}

	nP, ok := c.s.g.getCenterWithFallback(n.gridPos)
	if !ok {
		panic(fmt.Sprintf("planBombs: cannot get node pixel position %s", nP.log()))
	}
	n.targets = append(n.targets, position{layer: bomberConfig.produceLayer, PixelPosition: nP})
}
