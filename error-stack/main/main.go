package main

import (
	"encoding/json"
	std_errors "errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"

	"playground/package1"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

var workingDir string

func init() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err))
	}
	workingDir = dir
}

type Error struct {
	Message string   `json:"message"`
	Wraps   []string `json:"wrapping_messages,omitempty"`
	Calls   []string `json:"stack_calls,omitempty"`
}

func printStdErrors(err error) {
	fmt.Println("==============================")
	fmt.Printf("Using standard errors package\n\n")
	for err != nil {
		fmt.Println(err)
		err = std_errors.Unwrap(err)
	}
}

func printPkgErrors(err error) {
	fmt.Println("==============================")
	fmt.Printf("Using pkg/errors package\n\n")
	if sterr, ok := err.(stackTracer); ok {
		fmt.Printf("Stack trace:")
		for n, f := range sterr.StackTrace() {
			fmt.Printf("%d: %s %n:%d\n", n, f, f, f)
		}
	}
}

func printErisErrors(err error) {
	fmt.Println("==============================")
	fmt.Printf("Using rotisserie/eris package\n\n")
	fmt.Printf("1. Using default string\n\n")
	fmt.Println(eris.ToString(err, true))
	fmt.Printf("2. Using custom formatted string\n\n")
	fmt.Println(eris.ToCustomString(err, eris.StringFormat{
		Options: eris.FormatOptions{
			InvertOutput: true,
			WithTrace:    true,
			InvertTrace:  true,
			WithExternal: true,
		},
		MsgStackSep:  "\n\n",
		PreStackSep:  "  ",
		StackElemSep: ":",
		ErrorSep:     "\n",
	}))
	fmt.Printf("3. Using default JSON\n\n")
	fmt.Println(eris.ToJSON(err, true))
	u, _ := json.MarshalIndent(eris.ToJSON(err, true), "", "  ")
	fmt.Printf("%v\n", string(u))
	fmt.Printf("4. Using custom formatted JSON\n\n")
	format2 := eris.NewDefaultJSONFormat(eris.FormatOptions{
		InvertOutput: true,
		WithTrace:    true,
		InvertTrace:  true,
		WithExternal: true,
	})
	u, _ = json.MarshalIndent(eris.ToCustomJSON(err, format2), "", "  ")
	fmt.Printf("%v\n", string(u))
	fmt.Printf("5. Using eris.UnpackedError\n\n")

	unpackedErr := eris.Unpack(err)
	fmt.Printf("ErrExternal: %+v\n", unpackedErr.ErrExternal)
	fmt.Println("ErrRoot:")
	fmt.Printf("  Msg: %s\n", unpackedErr.ErrRoot.Msg)
	fmt.Println("  StackFrames:")
	for i := 0; i < len(unpackedErr.ErrRoot.Stack); i++ {
		fmt.Printf("    [%d]:\n", i)
		stack := unpackedErr.ErrRoot.Stack[i]
		fmt.Printf("      Name: %s\n", stack.Name)
		fmt.Printf("      File: %s\n", removeCurrentDirectory(stack.File))
		fmt.Printf("      Line: %d\n", stack.Line)
	}
	fmt.Println("ErrChain:")
	for i := 0; i < len(unpackedErr.ErrChain); i++ {
		fmt.Printf("    [%d]:\n", i)
		link := unpackedErr.ErrChain[i]
		fmt.Printf("      Msg: %s\n", link.Msg)
		fmt.Println("      StackFrame:")
		fmt.Printf("        Name: %s\n", link.Frame.Name)
		fmt.Printf("        File: %s\n", removeCurrentDirectory(link.Frame.File))
		fmt.Printf("        Line: %d\n", link.Frame.Line)
	}

	fmt.Printf("6. Using eris.UnpackedError to create custom message\n")
	customError := Error{Message: unpackedErr.ErrRoot.Msg}
	if len(unpackedErr.ErrRoot.Stack) > 0 {
		firstStack := unpackedErr.ErrRoot.Stack[0]
		errorLocation := fmt.Sprintf("%s[L%d]: %s", removeCurrentDirectory(firstStack.File), firstStack.Line, firstStack.Name)
		customError.Message = errorLocation + ": " + unpackedErr.ErrRoot.Msg
	}
	for _, l := range unpackedErr.ErrChain {
		customError.Wraps = append(customError.Wraps, fmt.Sprintf("%s[L%d]: %s: %s", removeCurrentDirectory(l.Frame.File), l.Frame.Line, l.Frame.Name, l.Msg))
	}
	for _, s := range unpackedErr.ErrRoot.Stack {
		customError.Calls = append(customError.Calls, fmt.Sprintf("%s[L%d]: %s", removeCurrentDirectory(s.File), s.Line, s.Name))
	}
	u, _ = json.MarshalIndent(customError, "", "  ")
	fmt.Printf("%v\n", string(u))
}

func printRuntimeDebugErrors() {
	fmt.Println("==============================")
	fmt.Printf("Using runtime/debug package\n\n")
	debug.PrintStack()
}

func main() {
	err1, err2, err3 := package1.Package1Function1(0)
	err1 = fmt.Errorf("this error is written in main: %w", err1)
	err2 = errors.Wrap(err2, "this error is written in main")

	printStdErrors(err1)
	printRuntimeDebugErrors()
	printPkgErrors(err2)
	printErisErrors(err3)
}

func removeCurrentDirectory(path string) string {
	return strings.Replace(path, workingDir+"/", "", -1)
}
