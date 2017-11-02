package errors

import (
	"bytes"
	native "errors"
	"fmt"
	"runtime"
)

type DetailedError interface {
	error
	Origin() error
	Inherit(error) DetailedError
	InheritNative(string) DetailedError
	CausedBy(error) DetailedError
}

type detailed struct {
	text   string
	file   string
	line   int
	origin error
	cause  error
}

func (e *detailed) Origin() error {
	return e.origin
}

func (e *detailed) Inherit(other error) DetailedError {
	switch err := other.(type) {
	case DetailedError:
		e.origin = err.Origin()
		e.cause = other
	default:
		e.origin = err
	}
	return e
}

func (e *detailed) InheritNative(text string) DetailedError {
	e.origin = native.New(text)
	return e
}

func (e *detailed) CausedBy(other error) DetailedError {
	e.cause = other
	return e
}

func (e *detailed) Error() string {
	buf := bytes.NewBufferString(e.text)
	fmt.Fprintf(buf, " @ %s:%d", e.file, e.line)
	if e.origin != nil && e.origin != Origin(e.cause) {
		fmt.Fprint(buf, " inhertis ", e.origin.Error())
	}
	if e.cause != nil {
		fmt.Fprint(buf, " caused by ", e.cause.Error())
	}
	return buf.String()
}

func (e *detailed) Format(s fmt.State, verb rune) {
	if verb == 'q' {
		fmt.Fprint(s, "\"")
	}
	switch verb {
	case 'v', 's', 'q':
		fmt.Fprintf(s, e.Error())
	default:
		fmt.Fprint(s, "%", string(verb), "(invalid for error ", e.Error(), ")")
	}
	if verb == 'q' {
		fmt.Fprint(s, "\"")
	}
}

func New(text string) DetailedError {
	_, file, line, _ := runtime.Caller(1)
	return &detailed{text: text, file: trimGOPATH(file), line: line}
}

func Newf(format string, a ...interface{}) DetailedError {
	return New(fmt.Sprintf(format, a...))
}

func Origin(err error) error {
	if e, ok := err.(DetailedError); ok {
		return e.Origin()
	} else {
		return err
	}
}
