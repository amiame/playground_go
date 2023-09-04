package main

import (
	"fmt"
	"myPkg/errors"
	"os"
	"strings"
)

var workingDir string

func init() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
	workingDir = dir
}

func funcA() error {
	return errors.New("hey you!")
}

func removeCurrentDirectory(path string) string {
	return strings.Replace(path, workingDir+"/", "", -1)
}

func main() {
	err := funcA()
	errors.Wrap(err, "I'm wrapped")
	errors.Wrap(err, "I'm wrapped2")

	upErr := errors.Unpack(err)
	for i, call := range upErr.Calls {
		if call.Msg != "" {
			fmt.Printf("%d) %s[%d]: %s: %s\n", i, removeCurrentDirectory(call.File), call.Line, call.Function, call.Msg)
		} else {
			fmt.Printf("%d) %s[%d]: %s\n", i, removeCurrentDirectory(call.File), call.Line, call.Function)
		}
	}
}
