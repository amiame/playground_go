package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type calls []uintptr

type Error struct {
	msg   string
	calls *calls
}

type CallFrame struct {
	File     string
	Function string
	Line     int
}

const callerDepth = 64

type SkipLevel int

const (
	skipThisFunction SkipLevel = iota
	skipGetCallsOnly
)

func New(msg string) error {
	return &Error{
		msg:   msg,
		calls: getCalls(skipThisFunction),
	}
}

func (err *Error) Error() string {
	return fmt.Sprint(err.msg)
}

func getCalls(skip SkipLevel) *calls {
	var sk int
	switch skip {
	case skipGetCallsOnly:
		sk = 2
	case skipThisFunction:
		sk = 3
	default:
		sk = 0
	}
	var pcs [callerDepth]uintptr
	length := runtime.Callers(sk, pcs[:])
	var cs calls = pcs[0:length]
	return &cs
}

func (cs *calls) getCallFrames() *[]CallFrame {
	cfs := make([]CallFrame, 0)

	frames := runtime.CallersFrames(*cs)
	for {
		frame, more := frames.Next()
		i := strings.LastIndex(frame.Function, "/")
		function := frame.Function[i+1:]
		if !strings.HasPrefix(strings.ToLower(function), "runtime.") {
			cfs = append(cfs, CallFrame{
				File:     frame.File,
				Function: function,
				Line:     frame.Line,
			})
		}
		if !more {
			break
		}
	}
	return &cfs
}

type UnpackedError struct {
	Msg   string
	Calls []CallFrame
}

func Unpack(err error) UnpackedError {
	switch err := err.(type) {
	case *Error:
		return UnpackedError{
			Msg:   err.msg,
			Calls: *err.calls.getCallFrames(),
		}
	default:
		return UnpackedError{}
	}
}
