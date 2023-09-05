package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type calls []call
type call struct {
	ptr   uintptr
	frame CallFrame
}

type myError struct {
	calls *calls
}

type CallFrame struct {
	Msg      string `json:"message"`
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
}

const callerDepth = 64

type SkipLevel int

const (
	skipThisFunctionAndOneAbove SkipLevel = iota
	skipThisFunction
	skipGetCallsOnly
)

func New(msg string) error {
	cs := getCalls(skipThisFunction)
	cs[0].frame.Msg = msg
	return &myError{
		calls: &cs,
	}
}

func (err *myError) Error() string {
	cs := *err.calls
	return fmt.Sprint(cs[0].frame.Msg)
}

func getCalls(skip SkipLevel) calls {
	var sk int
	switch skip {
	case skipGetCallsOnly:
		sk = 2
	case skipThisFunction:
		sk = 3
	case skipThisFunctionAndOneAbove:
		sk = 4
	default:
		sk = 0
	}
	var pcs [callerDepth]uintptr
	length := runtime.Callers(sk, pcs[:])
	frames := runtime.CallersFrames(pcs[0:length])
	c := make([]call, 0)
	for {
		frame, more := frames.Next()
		i := strings.LastIndex(frame.Function, "/")
		function := frame.Function[i+1:]
		if !strings.HasPrefix(strings.ToLower(function), "runtime.") {
			c = append(c, call{
				ptr: frame.PC,
				frame: CallFrame{
					File:     frame.File,
					Function: function,
					Line:     frame.Line,
				},
			})
		}
		if !more {
			break
		}
	}
	return c
}

type UnpackedError struct {
	Calls []CallFrame `json:"calls"`
}

func Unpack(err error) UnpackedError {
	switch err := err.(type) {
	case *myError:
		cfs := make([]CallFrame, 0)
		for _, c := range *err.calls {
			cfs = append(cfs, c.frame)
		}
		return UnpackedError{
			Calls: cfs,
		}
	default:
		return UnpackedError{}
	}
}

func Wrap(err error, msg string) error {
	return wrap(err, fmt.Sprint(msg))
}

func Wrapf(err error, format string, args ...interface{}) error {
	return wrap(err, fmt.Sprintf(format, args...))
}

func wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	cs := getCalls(skipThisFunctionAndOneAbove)
	switch e := err.(type) {
	case *myError:
		e.calls.insertMsg(cs, msg)
	}

	return err
}

func (cs *calls) insertMsg(newCs calls, msg string) {
	if len(newCs) == 0 {
		return
	}
	if len(newCs) == 1 {
		c := newCs[0]
		c.frame.Msg = msg
		*cs = append(*cs, c)
		return
	}
	for at, c := range *cs {
		if c.ptr == newCs[0].ptr {
			// break if the stack already contains the pc
			break
		} else if c.ptr == newCs[1].ptr {
			// insert the first call into the call stack if the second pc is found
			// this inserts the new call by breaking the call stack into two slices (cs[:at] and cs[at:])
			c := newCs[0]
			c.frame.Msg = msg
			*cs = append((*cs)[:at], append([]call{c}, (*cs)[at:]...)...)
			break
		}
	}
}
