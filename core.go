package errors

import (
	goerrors "errors"
	"fmt"
	"io"
	"runtime"
)

type Detailed interface {
	error
	Origin() error
}

type DetailedWithoutCause interface {
	Detailed
	CausedBy(cause error) DetailedWithCause
}

type DetailedWithCause interface {
	Detailed
	Cause() error
}

type detailed struct {
	origin error
	cause  error
	file   string
	line   int
	fn     string
}

func (err *detailed) CausedBy(cause error) DetailedWithCause {
	err.cause = cause
	return err
}

func (err *detailed) Origin() error {
	return err.origin
}

func (err *detailed) Cause() error {
	return err.cause
}

func (err *detailed) Error() string {
	if err.cause != nil {
		return err.origin.Error() + " caused by " + err.cause.Error()
	}
	return err.origin.Error()
}

func (err *detailed) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "%+v @ %s:%d", err.origin, err.file, err.line)
			if err.cause != nil {
				fmt.Fprintf(s, " caused by %+v", err.cause)
			}
			return
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v @ %s in %s:%d", err.origin, err.fn, err.file, err.line)
			if err.cause != nil {
				fmt.Fprintf(s, " caused by %#v", err.cause)
			}
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, err.Error())
	}
}

func wrap(origin error) DetailedWithoutCause {
	switch err := origin.(type) {
	case DetailedWithoutCause:
		return err
	default:
		pc, file, line, _ := runtime.Caller(2)
		return &detailed{origin: err, file: trimGOPATH(file), line: line, fn: runtime.FuncForPC(pc).Name()}
	}
}

func Wrap(origin error) DetailedWithoutCause {
	return wrap(origin)
}

func New(text string) DetailedWithoutCause {
	return wrap(goerrors.New(text))
}

func Unwrap(err error) error {
	switch t := err.(type) {
	case Detailed:
		return t.Origin()
	case DetailedWithCause:
		return t.Origin()
	}
	return err
}
