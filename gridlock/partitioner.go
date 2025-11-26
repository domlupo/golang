// Provides spatial partitions for different unit and analysis combinations.

package main

import (
	"fmt"
	"log/slog"
)

type partitions []*hexagon

// The type of analysis which will process the partitions.
type partitionAnalysis int

const (
	paNone partitionAnalysis = iota
	paCollision
	paVision
	paSteer
)

func (pa partitionAnalysis) String() string {
	switch pa {
	case paNone:
		return "paNone"
	case paCollision:
		return "paCollision"
	case paVision:
		return "paVision"
	case paSteer:
		return "paSteer"
	default:
		panic(fmt.Sprintf("Cannot convert partition analysis to string: %d", pa))
	}
}

func (pa partitionAnalysis) log() slog.Attr {
	return slog.String("partitionAnalysis", pa.String())
}

type partitioner struct {
	grid     grid
	hexagons hexagons
}

func newPartitioner(g grid, h hexagons) partitioner {
	p := partitioner{
		grid:     g,
		hexagons: h,
	}
	return p
}

func (part partitioner) partitions(u unit, pa partitionAnalysis, buffer int) partitions {
	var positions []GridPosition
	switch {
	case u.name == infantry && pa == paCollision:
		positions = part.defaultPartitions(u.PixelPosition, u.radius+buffer)
	case u.name == infantry && pa == paVision:
		positions = part.defaultPartitions(u.PixelPosition, u.vision+buffer)
	case u.name == infantry && pa == paSteer:
		positions = part.defaultPartitions(u.PixelPosition, buffer)
	case pa == paCollision:
		positions = part.defaultPartitions(u.PixelPosition, u.radius+buffer)
	case pa == paVision:
		positions = part.defaultPartitions(u.PixelPosition, u.vision+buffer)
	default:
		panic(fmt.Sprintf("partitions: cannot get partitions for unit %s, analysis %s and buffer %d", u.name, pa.String(), buffer))
	}
	var ps partitions
	for _, p := range positions {
		ps = append(ps, part.hexagons[p])
	}
	return ps
}

// Returns conservative partitions which collide with the passed position and distance.
// The caller is responsible for verifying the analysis performed with the partitions is valid.
// For example, analyzing partitions for collision detection between two circles would be invalid
// if the circle with the position and distance not passed has a radius which exceeds the
// conservative hex radius. Unoptimized because the hexagon collision uses conservative circles.
func (part partitioner) defaultPartitions(p PixelPosition, d int) []GridPosition {
	gridPos, ok := part.grid.getGridPos(p)
	if !ok {
		panic(fmt.Sprintf("defaultPartitions: cannot find grid pos for %s", p.log()))
	}
	visited := make(map[GridPosition]struct{})
	visited[gridPos] = struct{}{}
	q := []GridPosition{gridPos}
	nextQ := []GridPosition{}
	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		for _, a := range part.grid.allAdjacent(current) {
			if _, ok := visited[a]; !ok {
				visited[a] = struct{}{}
				// TODO: Fix this. This isn't correct but will probably work fine for awhile.
				// It doesn't handle border partition collisions 100% correct but this should be rare.
				// Also, I suspect this partitions method will most likey be refactored before issue
				// blocks development.
				pa, okA := part.grid.getCenterWithFallback(a)
				if !okA {
					panic(fmt.Sprintf("defaultPartitions: cannot find center for pos %s", a.log()))
				}
				_, collides := collides(p, pa, d, part.grid.conservativeHexRadius())
				if collides {
					nextQ = append(nextQ, a)
				}
			}
		}

		if len(q) == 0 && len(nextQ) == 0 {
			break
		} else if len(q) == 0 {
			q = nextQ
			nextQ = []GridPosition{}
		}
	}

	partitions := make([]GridPosition, len(visited))
	i := 0
	for p := range visited {
		partitions[i] = p
		i++
	}
	return partitions
}
