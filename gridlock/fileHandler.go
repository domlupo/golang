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
	directoryName      string
	pixelPathFileName  string
	nodeLookupFileName string
}

func newFileHandler() fileHandler {
	return fileHandler{
		directoryName:      "output",
		pixelPathFileName:  "pixelPath.json",
		nodeLookupFileName: "nodeLookup.json",
	}
}

// Golang encoding/json does not support direct marshall/unmarshall for custom types
// Multiple intermediate entry types are defined.
type pixelPathEntry struct {
	Key   PixelPosition `json:"key"`
	Value Path          `json:"value"`
}
type nodeLookupEntry struct {
	Key   PixelPosition `json:"key"`
	Value NearestNodes  `json:"value"`
}

func (f fileHandler) readPixelPath(obstacles obstacles, hc *hexConfig) pixelPath {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.pixelPathFileName))
	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "create pixel path file since it does not exist")
		return f.writePixelPath(obstacles, hc)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %e", err))
		return f.writePixelPath(obstacles, hc)
	}

	entries := make([]pixelPathEntry, hc.windowWidth*hc.windowHeight)
	json.Unmarshal(bytes, &entries)

	pixelPath := make(pixelPath, hc.windowWidth*hc.windowHeight)
	for _, e := range entries {
		pixelPath[e.Key] = e.Value
	}

	if len(pixelPath) != hc.windowWidth*hc.windowHeight {
		logger.log(errorLevel, setupLogic, "create pixel path since loaded pixel path was invalid size")
		return f.writePixelPath(obstacles, hc)
	}

	logger.log(infoLevel, setupLogic, "loaded pixel path, skipped creating it.")

	return pixelPath
}

func (f fileHandler) writePixelPath(obstacles obstacles, hc *hexConfig) pixelPath {
	pixelPath := obstacles.setupPixelPath(hc)

	entries := make([]pixelPathEntry, hc.windowWidth*hc.windowHeight)
	for k, v := range pixelPath {
		entries = append(entries, pixelPathEntry{Key: k, Value: v})
	}

	data, _ := json.Marshal(entries)
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, f.pixelPathFileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return pixelPath
}

func (f fileHandler) readNodeLookup(nodes nodes, pixelPath pixelPath, hc *hexConfig) nodeLookup {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.nodeLookupFileName))
	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "create node lookup file since it does not exist")
		return f.writeNodeLookup(nodes, pixelPath, hc)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %e", err))
		return f.writeNodeLookup(nodes, pixelPath, hc)
	}

	entries := make([]nodeLookupEntry, hc.windowWidth*hc.windowHeight)
	json.Unmarshal(bytes, &entries)

	nodeLookup := make(nodeLookup, hc.windowWidth*hc.windowHeight)
	for _, e := range entries {
		nodeLookup[e.Key] = e.Value
	}

	if len(nodeLookup) != hc.windowWidth*hc.windowHeight {
		logger.log(errorLevel, setupLogic, "create node lookup since loaded pixel path was invalid size")
		return f.writeNodeLookup(nodes, pixelPath, hc)
	}

	logger.log(infoLevel, setupLogic, "loaded node lookup, skipped creating it.")

	return nodeLookup
}

func (f fileHandler) writeNodeLookup(nodes nodes, pixelPath pixelPath, hc *hexConfig) nodeLookup {
	nodeLookup := nodes.setupLookup(pixelPath, hc)

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
