package main

import (
	"os"
	"path/filepath"
)

func appPath() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return path
}

func absPathToFile(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(appPath(), path)
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		logger.Fatalf("%s (%v)", msg, err)
	}
}
