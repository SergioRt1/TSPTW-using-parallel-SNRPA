package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetFileAsBytes(path string) ([]byte, error) {
	baseInsideBinary := filepath.Dir(pathInsideBinary())
	file, err := readRelativeOrAbsolutePath(baseInsideBinary + "/" + path)
	if err != nil {
		file, err = readRelativeOrAbsolutePath(path)
	}

	return file, err
}

func readRelativeOrAbsolutePath(path string) ([]byte, error) {
	file, err := readFile(path)
	if err != nil {
		var currentDirectory string
		currentDirectory, err = os.Getwd()
		if err == nil {
			currentDirectory, err = filepath.Abs(currentDirectory)
			if err == nil {
				return readFile(currentDirectory + "/" + path)
			}
		}
	}
	return file, err
}

func pathInsideBinary() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		f, ok := frames.Next()
		if !ok {
			break
		}
		if strings.Contains(f.Function, "/utils") {
			return filepath.Dir(f.File)
		}
	}
	return ""
}

func readFile(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err == nil {
		switch mode := info.Mode(); {
		case mode.IsDir():
			return nil, errors.New("path '" + path + "' in a directory")
		case mode.IsRegular():
			return ioutil.ReadFile(path)
		}
	}
	return nil, err
}
