package main

import (
	"bufio"
	"fmt"
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

func halt() {
	os.Exit(1)
}

func confirm(msg string) bool {
	for {
		fmt.Printf("%s [y/n]: ", msg)

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		switch {
		case answer == "y\n":
			return true
		case answer == "n\n":
			return false
		}
	}
}
