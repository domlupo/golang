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
	*state
	computersConfig
	player player
}

func newComputersConfig() computersConfig {
	c := computersConfig{
		count:               2,
		minUnitDifferential: 2,
	}
	return c
}

func newComputers(s *state) *computers {
	var computers = make(computers, 0)
	c := newComputersConfig()
	for i := range c.count {
		computers = append(computers, newComputer(s, i))
	}
	computers.updateNodetargets(true, Pixel{})
	return &computers
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
		state:           s,
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

func (computers *computers) updateNodetargets(nodeCapture bool, avgUnitPos Pixel) {
	for _, c := range *computers {
		if nodeCapture {
			c.updateForNodeCapture()
		}
		c.updateNodeTargets(avgUnitPos)
	}
}

func (c *computer) updateNodeTargets(avgUnitPos Pixel) {
	var opponent player
	switch c.player {
	case pOne:
		opponent = pTwo
	case pTwo:
		opponent = pOne
	default:
		panic(fmt.Sprintf("updateNodeTargets: unknown player: %d", c.player))
	}

	diffs := c.state.hexagons.maxUnitDifferential(opponent, lLand)
	validDiffsIndex := -1
	var positions []Pixel
	for _, d := range diffs {
		if d.diff > c.minUnitDifferential {
			validDiffsIndex++
		}

		p, ok := c.state.grid.getCenterWithFallback(d.GridPosition)
		if !ok {
			panic("updateNodeTargets: cannot get center with fallback")
		}
		positions = append(positions, p)
	}

	for _, n := range c.state.nodes {
		if n.player != c.player {
			continue
		}

		if n.unitConfig.name == infantry {
			continue
		}

		nP, ok := c.grid.getCenterWithFallback(n.gridPos)
		if !ok {
			panic(fmt.Sprintf("updateNodeTargets: cannot get node pixel position %s", nP.log()))
		}
		selfNode := n
		selfNode.targets = []position{{lLand, nP}}
		switch n.unitConfig.name {
		case missileCruise, missileBallistic:
			if diffs[0].diff > 1 {
				selfNode.targets = []position{{lLand, positions[0]}}
			}
		case missileHeatSeeking, missileAntiBallistic, missileRadarActive:
			if !avgUnitPos.equals(Pixel{}) {
				selfNode.targets = []position{{lLand, avgUnitPos}}
			}
		case fighter:
			if fighterConfig.bombGravityCountMax <= validDiffsIndex+1 {
				c.planBombs(selfNode, positions[:fighterConfig.bombGravityCountMax])
			}
		case bomber:
			if bomberConfig.bombGravityCountMax <= validDiffsIndex+1 {
				c.planBombs(selfNode, positions[:bomberConfig.bombGravityCountMax])
			}
		}
	}
}

func (c *computer) updateForNodeCapture() {
	for _, n := range c.nodes {
		if n.player != c.player {
			continue
		}

		selfNodePos := n.gridPos
		selfNode := n
		selfNode.targets = []position{}
		nearestTargetDistanceSquared := math.MaxInt

		if n.unitConfig.name == missileCruise {
			continue
		}

		for _, n := range c.nodes {
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
				var targetPixelPos, ok = c.grid.getPixel(otherNodePos, hpCenter)
				if !ok {
					panic("could not get pixel pos for update computer node targets")
				}

				selfNode.targets = []position{{lLand, targetPixelPos}}
			}

		}
	}

	for u := range c.units() {
		if u.staticTarget ||
			u.name == bomber ||
			u.name == fighter {
			continue
		}

		if u.player != c.player {
			continue
		}

		nearestTargetDistanceSquared := math.MaxInt

		for _, n := range c.nodes {
			if n.player == c.player {
				continue
			}

			nPixelPos, ok := c.grid.getPixel(n.gridPos, hpCenter)
			if !ok {
				panic("could not get pixel pos for update computer unit targets")
			}
			xDelta := (u.X - nPixelPos.X)
			yDelta := (u.Y - nPixelPos.Y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta

			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				u.updateTarget(globalTarget, position{lLand, nPixelPos})
			}
		}
	}
}

func (c *computer) planBombs(n *node, targets []Pixel) {
	n.targets = []position{}
	n.actions = []action{}

	for _, t := range targets {
		p := position{lAircraft, t}
		n.targets = append(n.targets, p)
		n.actions = append(n.actions, action{position: p, do: func() {
			u := newUnit(
				bombGravityConfig,
				position{layer: bombGravityConfig.produceLayer, Pixel: p.Pixel},
				c.player,
				[]position{{layer: lLand, Pixel: p.Pixel}},
				[]action{},
				GridPosition{},
				false)
			c.spawnRequest(request{spawn: &u})
		}})
	}

	nP, ok := c.grid.getCenterWithFallback(n.gridPos)
	if !ok {
		panic(fmt.Sprintf("planBombs: cannot get node pixel position %s", nP.log()))
	}
	n.targets = append(n.targets, position{layer: bomberConfig.produceLayer, Pixel: nP})
}
