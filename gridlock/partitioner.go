// Provides spatial partitions for different unit and analysis combinations.

package main

import (
	"fmt"
	"log/slog"
	"slices"
)

type partitions []*hexagon

// The type of analysis which will process the partitions.
type partitionAnalysis int

// A look up table (LUT) to find an optimized number of partitions. It avoids runtime collision in two ways.
// First, collision checks for all 1380 Pixels in the demotif hexagon are done at startup for the
// 132 adjacent Pixels. If a collision occurs the hexSide which caused it is stored. Second, optimized
// number of partitions means no extra partitions are considered for unit <-> unit collisions since spatial
// partioning checks all units within the partition. The LUT is utilized via: Unit Pixel -> demotif ->
// hex sides from LUT -> check unit grid position for collision -> check hex sides from LUT for collision.
type partitionLUT map[demotif][]hexSide

const (
	paNone partitionAnalysis = iota
	paCollision
	paVision
	paSteer
)

func (pa partitionAnalysis) string() string {
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
	return slog.String("partitionAnalysis", pa.string())
}

type partitioner struct {
	grid                 grid
	hexagons             hexagons
	infantryCollisionLUT partitionLUT
}

func newPartitioner(g grid, h hexagons) partitioner {
	p := partitioner{
		grid:     g,
		hexagons: h,
	}
	p.infantryCollisionLUT = p.setupPartitionLUT(infantryConfig, infantryConfig, paCollision)
	return p
}

func (part partitioner) partitions(u unit, pa partitionAnalysis, buffer int) partitions {
	var positions []GridPosition
	radius := buffer
	switch {
	case u.name == infantry && pa == paCollision:
		demotif, demotifOk := part.grid.demotif(u.Pixel)
		if !demotifOk {
			panic(fmt.Sprintf("partitions: could not demotif pixel %s", u.Pixel.log()))
		}
		hexSides := part.infantryCollisionLUT[demotif]
		gp, gpOk := part.grid.getGridPos(u.Pixel)
		positions = append(positions, gp)
		if !gpOk {
			panic(fmt.Sprintf("partitions: could not find grid position for pixel %s", u.Pixel.log()))
		}
		for _, hs := range hexSides {
			if gp, ok := part.grid.adjacent(gp, hs); ok {
				positions = append(positions, gp)
			}
		}
	case pa == paCollision:
		radius += u.radius
		positions = part.defaultPartitions(u.Pixel, radius)
	case pa == paVision:
		radius += u.vision
		positions = part.defaultPartitions(u.Pixel, radius)
	case pa == paSteer:
		radius += 0
		positions = part.defaultPartitions(u.Pixel, radius)
	default:
		panic(fmt.Sprintf("partitions: unknown partitions analysis %s", pa.string()))
	}
	if pa == paSteer && u.name != infantry {
		panic(fmt.Sprintf("partitions: cannot get partitions for unit %s, analysis %s and buffer %d", u.name, pa.string(), buffer))
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
func (part partitioner) defaultPartitions(p Pixel, d int) []GridPosition {
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

func (part partitioner) setupPartitionLUT(internalUC unitConfig, externalUC unitConfig, pa partitionAnalysis) partitionLUT {
	var internalRadius int
	var externalRadius int
	switch pa {
	case paCollision:
		internalRadius = internalUC.radius
		externalRadius = externalUC.radius
	default:
		panic(fmt.Sprintf("setupPartitionsLUT: unsupported partition analysis %s", pa.string()))
	}

	lut := make(partitionLUT)
	collisionSources := part.grid.demotifHexPixelsAdjacent()
	pixels := part.grid.demotifHexPixels()
	for externalP, hs := range collisionSources {
		for _, internalP := range pixels {
			if _, ok := collides(internalP, externalP, internalRadius, externalRadius); ok {
				hexSides := lut[demotif{internalP}]
				if !slices.Contains(hexSides, hs) {
					hexSides = append(hexSides, hs)
					lut[demotif{internalP}] = hexSides
				}
			}
		}
	}
	return lut
}
