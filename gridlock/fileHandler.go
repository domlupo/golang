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
	nodeLookupFileName  string
}

func newFileHandler() fileHandler {
	return fileHandler{
		directoryName:       "output",
		obstacleLUTFileName: "obstacleLUT.json",
		nodeLookupFileName:  "nodeLookup.json",
	}
}

// Golang encoding/json does not support direct marshall/unmarshall for custom types
// Multiple intermediate entry types are defined.
type obstacleLUTEntry struct {
	Key   PixelPosition `json:"key"`
	Value obstacleInfo  `json:"value"`
}
type nodeLookupEntry struct {
	Key   PixelPosition `json:"key"`
	Value NearestNodes  `json:"value"`
}

func (f fileHandler) readObstacleLUT(obstacles obstacles, hc *hexConfig) obstacleLUT {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.obstacleLUTFileName))
	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "readObstacleLUT: create obstacle look up table file since it does not exist")
		return f.writeObstacleLUT(obstacles, hc)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readObstacleLUT: file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("readObstacleLUT: file handler could not read json file: %e", err))
		return f.writeObstacleLUT(obstacles, hc)
	}

	entries := make([]obstacleLUTEntry, hc.windowWidth*hc.windowHeight)
	json.Unmarshal(bytes, &entries)

	oLUT := make(obstacleLUT, hc.windowWidth*hc.windowHeight)
	for _, e := range entries {
		oLUT[e.Key] = e.Value
	}

	if len(oLUT) != hc.windowWidth*hc.windowHeight {
		logger.log(errorLevel, setupLogic, "readObstacleLUT: cannot create obstacle look up table since loaded obstacle look up was invalid size")
		return f.writeObstacleLUT(obstacles, hc)
	}

	logger.log(infoLevel, setupLogic, "loaded obstacle look up table, skipped creating it.")

	return oLUT
}

func (f fileHandler) writeObstacleLUT(obstacles obstacles, hc *hexConfig) obstacleLUT {
	oLUT := obstacles.setupObstacleLUT(hc)

	entries := make([]obstacleLUTEntry, hc.windowWidth*hc.windowHeight)
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

func (f fileHandler) readNodeLookup(nodes nodes, obstacleLUT obstacleLUT, hc *hexConfig) nodeLookup {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.nodeLookupFileName))
	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "create node lookup file since it does not exist")
		return f.writeNodeLookup(nodes, obstacleLUT, hc)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %e", err))
		return f.writeNodeLookup(nodes, obstacleLUT, hc)
	}

	entries := make([]nodeLookupEntry, hc.windowWidth*hc.windowHeight)
	json.Unmarshal(bytes, &entries)

	nodeLookup := make(nodeLookup, hc.windowWidth*hc.windowHeight)
	for _, e := range entries {
		nodeLookup[e.Key] = e.Value
	}

	if len(nodeLookup) != hc.windowWidth*hc.windowHeight {
		logger.log(errorLevel, setupLogic, "create node lookup since loaded noode lookup was invalid size")
		return f.writeNodeLookup(nodes, obstacleLUT, hc)
	}

	logger.log(infoLevel, setupLogic, "loaded node lookup, skipped creating it.")

	return nodeLookup
}

func (f fileHandler) writeNodeLookup(nodes nodes, obstacleLUT obstacleLUT, hc *hexConfig) nodeLookup {
	nodeLookup := nodes.setupLookup(obstacleLUT, hc)

	entries := make([]nodeLookupEntry, hc.windowWidth*hc.windowHeight)
	for k, v := range nodeLookup {
		entries = append(entries, nodeLookupEntry{Key: k, Value: v})
	}

	data, _ := json.Marshal(entries)
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, f.nodeLookupFileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return nodeLookup
}

func (f fileHandler) key(obstacles obstacles, hc hexConfig) {
	// TODO: Add CRC32 key
}
