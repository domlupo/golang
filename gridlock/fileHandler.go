// Contains a file handler for reading and writing various data to files.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type fileHandler struct {
	directoryName       string
	obstacleLUTFileName string
	nodeLUTFileName     string
}

func newFileHandler() fileHandler {
	return fileHandler{
		directoryName:       "output",
		obstacleLUTFileName: "obstacleLUT.json",
		nodeLUTFileName:     "nodeLUT.json",
	}
}

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
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.obstacleLUTFileName))
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
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, f.obstacleLUTFileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return oLUT
}

func (f fileHandler) readNodeLUT(nodes nodes, obstacleLUT obstacleLUT, g grid) nodeLUT {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.nodeLUTFileName))
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
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, f.nodeLUTFileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return nLUT
}

func (f fileHandler) key(obstacles obstacles, hc hexConfig) {
	// TODO: Add CRC32 key
}
