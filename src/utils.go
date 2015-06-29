package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const codeAlphabet = "DN3bY2cMdP1FCfgXLQhjk06lKmRnVB5pJqr9SstTvH4wxWy7ZzG8"

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", msg)

		answer, _ := reader.ReadString('\n')

		switch answer {
		case "y\n":
			return true
		case "n\n":
			return false
		}
	}
}
