// Copyright © 2018 Chedder <dev@javelang.com>
//
// MIT License

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OrcaLLC/jave/src/lexer"
)

func main() {
	// Fetch arguments to get the filename
	runArgs := os.Args[1:]
	filePath := runArgs[0]

	// Jave can only run one file at a time, dingus
	if len(runArgs) > 1 {
		log.Printf("Additional arguments ignored. Using: [%s] but found %v\n", filePath, runArgs)
	}

	// Load what we need then load either the file handler or input getter
	thing, err := lexer.NewJaveFile(filePath)
	if err != nil {
		log.Fatalf("[JAVE FAIL]: %v\n", err)
	}
	fmt.Printf("%v", thing)
}
