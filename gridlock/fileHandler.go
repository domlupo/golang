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

// Golang encoding/json does not support direct marshall/unmarshall for custom types
// Multiple intermediate entry types are defined.
type obstacleLUTEntry struct {
	Key   PixelPosition `json:"key"`
	Value obstacleInfo  `json:"value"`
}
type nodeLUTEntry struct {
	Key   PixelPosition `json:"key"`
	Value nodeInfo      `json:"value"`
}

func (f fileHandler) readObstacleLUT(obstacles obstacles, g grid) obstacleLUT {
	wc := newWindowConfig()

	fileName := f.fileName(fObstacle, obstacles.key())
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, fileName))

	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "readObstacleLUT: create look up table file since it does not exist")
		return f.writeObstacleLUT(obstacles, g)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readObstacleLUT: file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readObstacleLUT: file handler could not read json file: %e", err))
		return f.writeObstacleLUT(obstacles, g)
	}

	entries := make([]obstacleLUTEntry, wc.maxPixels())
	json.Unmarshal(bytes, &entries)

	oLUT := make(obstacleLUT, wc.maxPixels())
	for _, e := range entries {
		oLUT[e.Key] = e.Value
	}

	if len(oLUT) != wc.maxPixels() {
		logger.log(errorLevel, setupLogic, "readObstacleLUT: cannot read obstacle look up table since loaded obstacle look up was invalid size")
		return f.writeObstacleLUT(obstacles, g)
	}

	logger.log(infoLevel, setupLogic, "readObstacleLUT: loaded obstacle look up table, skipped creating it.")

	return oLUT
}

func (f fileHandler) writeObstacleLUT(obstacles obstacles, g grid) obstacleLUT {
	oLUT := obstacles.setupObstacleLUT(g)

	entries := make([]obstacleLUTEntry, newWindowConfig().maxPixels())
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
		logger.log(infoLevel, setupLogic, "readNodeLUT: cannot create node lookup table file since it does not exist")
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readNodeLUT: file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readNodeLUT: file handler could not read json file: %e", err))
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	}

	wc := newWindowConfig()
	entries := make([]nodeLUTEntry, wc.maxPixels())
	json.Unmarshal(bytes, &entries)

	nLUT := make(nodeLUT, wc.maxPixels())
	for _, e := range entries {
		nLUT[e.Key] = e.Value
	}

	if len(nLUT) != wc.maxPixels() {
		logger.log(errorLevel, setupLogic, "readNodeLUT: create node lookup table since loaded node lookup table was invalid size")
		return f.writeNodeLUT(nodes, obstacleLUT, g)
	}

	logger.log(infoLevel, setupLogic, "readNodeLUT: loaded node lookup table, skipped creating it.")

	return nLUT
}

func (f fileHandler) writeNodeLUT(nodes nodes, obstacleLUT obstacleLUT, g grid) nodeLUT {
	nLUT := nodes.setupNodeLUT(obstacleLUT, g)

	entries := make([]nodeLUTEntry, newWindowConfig().maxPixels())
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
		panic(fmt.Sprintf("Cannot get file name since unkonw file: %d", file))
	}

	return fmt.Sprintf("%s-%d.json", prefix, crc32.ChecksumIEEE(bytes))
}

func truncate(i int) byte {
	if i > math.MaxUint8 || i < 0 {
		panic(fmt.Sprintf("Cannot key obstacles due to byte truncation: %d", i))
	}
	return byte(i)
}
