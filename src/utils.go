package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const codeAlphabet = "DN3bY2cMdP1FCfgXLQhjk06lKmRnVB5pJqr9SstTvH4wxWy7ZzG8"

func appPath() (path string) {
	path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
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

func buildCode(code int64) string {
	buf := new(bytes.Buffer)
	alphabetLen := int64(len(codeAlphabet))

	for code > 0 {
		buf.WriteByte(codeAlphabet[code%alphabetLen])
		code = code / alphabetLen
	}

	return buf.String()
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
