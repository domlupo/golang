// Contains a file handler for reading and writing various data to files.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"os"
	"slices"
)

type fileHandler struct {
	directoryName             string
	obstacleLUTFileNamePrefix string
	nodeLUTFileNamePrefix     string
}

func newFileHandler() fileHandler {
	return fileHandler{
		directoryName:             "output",
		obstacleLUTFileNamePrefix: "obstacleLUT",
		nodeLUTFileNamePrefix:     "nodeLUT",
	}
}

type file int

const (
	fNone file = iota
	fObstacle
	fNode
)

type keys []Tile

// Deterministically sort the keys.
func (keys keys) sort() {
	slices.SortFunc(keys, func(i, j Tile) int {
		switch {
		case i.X == j.X && i.Y < j.Y:
			return -1
		case i.X < j.X:
			return -1
		default:
			return 1
		}
	})
}

// Golang encoding/json does not support direct marshalling/unmarshalling for custom types
// Multiple intermediate entry types are defined.
type obstacleLUTEntry struct {
	Key   Pixel        `json:"key"`
	Value obstacleInfo `json:"value"`
}
type nodeLUTEntry struct {
	Key   Pixel    `json:"key"`
	Value nodeInfo `json:"value"`
}

func (f fileHandler) readObstacleLUT(obstacles obstacles, g grid) obstacleLUT {
	fileName := f.fileName(fObstacle, obstacles.key())
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, fileName))

	if errors.Is(err, os.ErrNotExist) {
		logger.log(
			infoLevel,
			setupLogic,
			"fileHandler.readObstacleLUT: create look up table file since it does not exist")
		return f.writeObstacleLUT(obstacles, g)
	} else if err != nil {
		logger.log(
			errorLevel,
			setupLogic,
			fmt.Sprintf("fileHandler.readObstacleLUT: file handler could not open json file: %v", err))
		return f.writeObstacleLUT(obstacles, g)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(
			errorLevel,
			setupLogic,
			fmt.Sprintf("fileHandler.readObstacleLUT: file handler could not read json file: %e", err))
		return f.writeObstacleLUT(obstacles, g)
	}

	entries := make([]obstacleLUTEntry, 0, g.basePixels)
	json.Unmarshal(bytes, &entries)

	oLUT := make(obstacleLUT, g.basePixels)
	for _, e := range entries {
		oLUT[e.Key] = e.Value
	}

	if len(oLUT) != g.basePixels {
		logger.log(
			errorLevel,
			setupLogic,
			fmt.Sprintf(
				"filerHandler.readObstacleLUT: cannot read obstacle look up table since loaded LUT is %d "+
					"pixels rather than expected %d",
				len(oLUT),
				g.basePixels))
		oLUT = f.writeObstacleLUT(obstacles, g)
	} else {
		logger.log(
			infoLevel,
			setupLogic,
			"fileHandler.readObstacleLUT: loaded obstacle look up table, skipped creating it.")
	}

	return oLUT
}

func (f fileHandler) writeObstacleLUT(obstacles obstacles, g grid) obstacleLUT {
	oLUT := obstacles.setupObstacleLUT(g)

	entries := make([]obstacleLUTEntry, 0, len(oLUT))
	for k, v := range oLUT {
		entries = append(entries, obstacleLUTEntry{Key: k, Value: v})
	}

	data, _ := json.Marshal(entries)
	fileName := f.fileName(fObstacle, obstacles.key())
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, fileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return oLUT
}

func (f fileHandler) readNodeLUT(nodes nodes, obstacleLUT obstacleLUT, g grid) nodeLUT {
	fileName := f.fileName(fNode, nodes.key())
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, fileName))

	if errors.Is(err, os.ErrNotExist) {
		logger.log(
			infoLevel,
			setupLogic,
			"fileHandler.readNodeLUT: cannot create node lookup table file since it does not exist")
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	} else if err != nil {
		logger.log(
			errorLevel,
			setupLogic,
			fmt.Sprintf("fileHandler.readNodeLUT: file handler could not open json file: %v", err))
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(
			errorLevel,
			setupLogic,
			fmt.Sprintf("fileHandler.readNodeLUT: file handler could not read json file: %e", err))
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	}

	entries := make([]nodeLUTEntry, 0, g.basePixels)
	json.Unmarshal(bytes, &entries)

	nLUT := make(nodeLUT, g.basePixels)
	for _, e := range entries {
		nLUT[e.Key] = e.Value
	}

	if len(nLUT) != g.basePixels {
		logger.log(
			errorLevel,
			setupLogic,
			"fileHandler.readNodeLUT: create node LUT since loaded node LUT was invalid size")
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	}

	logger.log(
		infoLevel,
		setupLogic,
		"fileHandler.readNodeLUT: loaded node lookup table, skipped creating it.")

	return nLUT
}

func (f fileHandler) writeNodeLUT(nodes nodes, obstacleLUT obstacleLUT, g grid) nodeLUT {
	nLUT := nodes.setupNodeLUT(obstacleLUT, g)

	entries := make([]nodeLUTEntry, 0, len(nLUT))
	for k, v := range nLUT {
		entries = append(entries, nodeLUTEntry{Key: k, Value: v})
	}

	data, _ := json.Marshal(entries)
	fileName := f.fileName(fNode, nodes.key())
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, fileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return nLUT
}

func (f fileHandler) fileName(file file, bytes []byte) string {
	var prefix string
	switch file {
	case fObstacle:
		prefix = f.obstacleLUTFileNamePrefix
	case fNode:
		prefix = f.nodeLUTFileNamePrefix
	default:
		panic(fmt.Sprintf("fileHandler.fileName: unknown file: %d", file))
	}

	return fmt.Sprintf("%s-%d.json", prefix, crc32.ChecksumIEEE(bytes))
}

func truncate(i int) byte {
	if i > math.MaxUint8 || i < 0 {
		s := "truncate: cannot key due to byte truncation: %d"
		panic(fmt.Sprintf(s, i))
	}
	return byte(i)
}
