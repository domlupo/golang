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
	player             player
	unitConfigProvider unitConfigProvider
}

func newComputersConfig() computersConfig {
	c := computersConfig{
		count:               2,
		minUnitDifferential: 2,
	}
	return c
}

func newComputers(s *state, ucProvider unitConfigProvider) *computers {
	var computers = make(computers, 0)
	c := newComputersConfig()
	for i := range c.count {
		computers = append(computers, newComputer(s, ucProvider, i))
	}
	computers.updateNodetargets(true, Pos{})
	return &computers
}

func newComputer(s *state, ucProvider unitConfigProvider, i int) computer {
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
		state:              s,
		unitConfigProvider: ucProvider,
		player:             p,
		computersConfig:    newComputersConfig(),
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

func (computers *computers) updateNodetargets(nodeCapture bool, avgUnitPos Pos) {
	for _, c := range *computers {
		if nodeCapture {
			c.updateForNodeCapture()
		}
		c.updateNodeTargets(avgUnitPos)
	}
}

func (c *computer) updateNodeTargets(avgUnitPos Pos) {
	var opponent player
	switch c.player {
	case pOne:
		opponent = pTwo
	case pTwo:
		opponent = pOne
	default:
		panic(fmt.Sprintf("computer.updateNodeTargets: unknown player: %d", c.player))
	}

	diffs := c.state.hexagons.maxUnitDifferential(opponent, lLand)
	validDiffsIndex := -1
	var pixels []Pixel
	for _, d := range diffs {
		if d.diff > c.minUnitDifferential {
			validDiffsIndex++
		}

		p, ok := c.state.grid.getCenterWithFallback(d.GridPosition)
		if !ok {
			panic(fmt.Sprintf(
				"computer.updateNodeTargets: cannot get center with fallback for grid pos %s",
				d.GridPosition.log()))
		}
		pixels = append(pixels, p)
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
			panic(fmt.Sprintf(
				"computer.updateNodeTargets: cannot get pixel for grid pos %s",
				n.gridPos.log()))
		}
		selfNode := n
		selfNode.targets = []Pixel{nP}
		switch n.unitConfig.name {
		case missileCruise, missileBallistic:
			if diffs[0].diff > 1 {
				selfNode.targets = []Pixel{pixels[0]}
			}
		case missileHeatSeeking, missileAntiBallistic, missileRadarActive:
			if !avgUnitPos.equals(Pos{}) {
				selfNode.targets = []Pixel{avgUnitPos.pixel()}
			}
		case fighter, bomber:
			bombCount := n.unitConfig.consumables[bombGravity].count
			if bombCount <= validDiffsIndex+1 {
				c.planBombs(selfNode, pixels[:bombCount])
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
		selfNode.targets = []Pixel{}
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
				var targetPixel, ok = c.grid.getPixel(otherNodePos, hpCenter)
				if !ok {
					panic(fmt.Sprintf(
						"computer.updateForNodeCapture: could not get pixel for grid pos %s",
						otherNodePos.log()))
				}

				selfNode.targets = []Pixel{targetPixel}
			}

		}
	}

	for u := range c.units() {
		if u.isStaticTarget ||
			u.name == bomber ||
			u.name == fighter {
			continue
		}

		if u.player != c.player {
			continue
		}

		nearestTargetDistanceSquared := math.MaxFloat64

		for _, n := range c.nodes {
			if n.player == c.player {
				continue
			}

			nPixel, ok := c.grid.getPixel(n.gridPos, hpCenter)
			if !ok {
				panic(fmt.Sprintf(
					"computer.updateForNodeCapture: could not get pixel for grid pos %s",
					n.gridPos.log()))
			}
			n := nPixel.pos()
			xDelta := (u.X - n.X)
			yDelta := (u.Y - n.Y)
			otherTargetDistanceSquared := xDelta*xDelta + yDelta*yDelta
			if otherTargetDistanceSquared < nearestTargetDistanceSquared {
				nearestTargetDistanceSquared = otherTargetDistanceSquared
				u.updateTarget(globalTarget, nPixel)
			}
		}
	}
}

func (c *computer) planBombs(n *node, targets []Pixel) {
	n.targets = []Pixel{}
	n.actions = []action{}

	bombGravityUnitConfig := c.unitConfigProvider.provide(bombGravity)
	for _, t := range targets {
		n.targets = append(n.targets, t)
		n.actions = append(n.actions, action{Pixel: t, do: func() {
			u := newUnit(
				bombGravityUnitConfig,
				t.posLayer(bombGravityUnitConfig.produceLayer),
				c.player,
				[]Pixel{t},
				[]action{},
				GridPosition{},
				false)
			c.spawnRequest(request{spawn: u})
		}})
	}

	nP, ok := c.grid.getCenterWithFallback(n.gridPos)
	if !ok {
		panic(fmt.Sprintf("computer.planBombs: cannot get pixel for grid pos %s", n.gridPos.log()))
	}
	n.targets = append(n.targets, nP)
}
