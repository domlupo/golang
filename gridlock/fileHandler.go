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
	directoryName string
	fileName      string
}

func newFileHandler() fileHandler {
	return fileHandler{
		directoryName: "output",
		fileName:      "pixelPath.json",
	}
}

// Golang encoding/json does not support direct marshall/unmarshall for custom types
// An intermediate entry type is defined.
type pixelPathEntry struct {
	Key   PixelPosition `json:"key"`
	Value Path          `json:"value"`
}

func (f fileHandler) read(obstacles obstacles, hc *hexConfig) pixelPath {
	file, err := os.Open(fmt.Sprintf("%s/%s", f.directoryName, f.fileName))
	if errors.Is(err, os.ErrNotExist) {
		logger.log(infoLevel, setupLogic, "create pixel path file since it does not exist")
		return f.write(obstacles, hc)
	} else if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %v", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		logger.log(errorLevel, setupLogic, fmt.Sprintf("file handler could not read json file: %e", err))
		return f.write(obstacles, hc)
	}

	entries := make([]pixelPathEntry, hc.windowWidth*hc.windowHeight)
	json.Unmarshal(bytes, &entries)

	pixelPath := make(pixelPath, hc.windowWidth*hc.windowHeight)
	for _, e := range entries {
		pixelPath[e.Key] = e.Value
	}

	if len(pixelPath) != hc.windowWidth*hc.windowHeight {
		logger.log(errorLevel, setupLogic, "create pixel path since loaded pixel path was invalid size")
		return f.write(obstacles, hc)
	}

	logger.log(infoLevel, setupLogic, "loaded pixel path, skipped creating it.")

	return pixelPath
}

func (f fileHandler) write(obstacles obstacles, hc *hexConfig) pixelPath {
	pixelPath := obstacles.setupPixelPath(hc)

	entries := make([]pixelPathEntry, hc.windowWidth*hc.windowHeight)
	for k, v := range pixelPath {
		entries = append(entries, pixelPathEntry{Key: k, Value: v})
	}

	data, _ := json.Marshal(entries)
	err := os.WriteFile(fmt.Sprintf("%s/%s", f.directoryName, f.fileName), data, 0644)
	if err != nil {
		panic(err)
	}

	return pixelPath
}

func (f fileHandler) key(obstacles obstacles, hc hexConfig) {
	// TODO: Add CRC32 key
}
