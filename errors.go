package errors

import (
	"bytes"
	native "errors"
	"fmt"
	"io"
	"runtime"
)

type Detailed interface {
	error
	WithAnnotation(string) Detailed
	CausedBy(error) Detailed
	Origin() error
}

type detailed struct {
	origin     error
	annotation string
	file       string
	line       int
	fn         string
	list       []error
}

func (err *detailed) WithAnnotation(text string) Detailed {
	err.annotation = text
	return err
}

func (err *detailed) CausedBy(cause error) Detailed {
	if err.list == nil {
		err.list = make([]error, 0, 1)
	}
	err.list = append(err.list, cause)
	return err
}

func (err *detailed) Origin() error {
	return err.origin
}

func (err *detailed) Error() string {
	buf := new(bytes.Buffer)
	if len(err.annotation) > 0 {
		buf.WriteString(err.annotation)
		buf.WriteString(": ")
	}
	buf.WriteString(err.origin.Error())
	for _, cause := range err.list {
		buf.WriteString(" caused by ")
		buf.WriteString(cause.Error())
	}
	return buf.String()
}

func (err *detailed) Format(s fmt.State, verb rune) {
	if verb == 'q' {
		fmt.Fprint(s, "\"")
	}
	if len(err.annotation) > 0 {
		io.WriteString(s, err.annotation)
		io.WriteString(s, ": ")
	}
	io.WriteString(s, err.origin.Error())

	switch verb {
	case 'v', 's', 'q':
		fmt.Fprintf(s, " @ %s:%d", err.file, err.line)
		for _, cause := range err.list {
			fmt.Fprintf(s, " caused by %v", cause)
		}
	default:
		fmt.Fprintf(s, "%%%c(invalid for %T)", verb, err)
	}
	if verb == 'q' {
		fmt.Fprint(s, "\"")
	}
}

func withDetails(err error) *detailed {
	pc, file, line, _ := runtime.Caller(2)
	return &detailed{origin: err, file: trimGOPATH(file), line: line, fn: runtime.FuncForPC(pc).Name()}
}

func Wrap(err error) Detailed {
	switch e := err.(type) {
	case Detailed:
		return e
	default:
		return withDetails(err)
	}
}

func New(text string) Detailed {
	return withDetails(native.New(text))
}

func Newf(format string, a ...interface{}) Detailed {
	return withDetails(fmt.Errorf(format, a...))
}

func Origin(err error) error {
	if det, ok := err.(Detailed); ok {
		return Origin(det.Origin())
	} else {
		return err
	}
}
