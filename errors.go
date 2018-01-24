package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
)

type DetailedError struct {
	base    error
	file    string
	line    int
	related []error
	stack   []uintptr
}

func (e *DetailedError) relatedString() string {
	buf := new(bytes.Buffer)
	if len(e.related) > 0 {
		fmt.Fprint(buf, "related errors:")
		for _, err := range e.related {
			fmt.Fprint(buf, " {", err.Error(), "}")
		}
	}
	return buf.String()
}

func (e *DetailedError) Error() string {
	if len(e.related) == 0 {
		return fmt.Sprintf("%s @ %s:%d", e.base.Error(), e.file, e.line)
	}
	return fmt.Sprintf("%s @ %s:%d | %s", e.base.Error(), e.file, e.line, e.relatedString())
}

func (e *DetailedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		plusFlag := s.Flag('+')
		sharpFlag := s.Flag('#')
		if plusFlag || sharpFlag {
			fmt.Fprintf(s, "%s @ %s:%d", e.base.Error(), e.file, e.line)
			if sharpFlag {
				if len(e.related) > 0 {
					fmt.Fprint(s, "\nRelated errors:\n")
					for _, err := range e.related {
						fmt.Fprint(s, "  ", err.Error(), "\n")
					}
				}
			} else {
				fmt.Fprint(s, " | ", e.relatedString())
			}
			if len(e.stack) > 0 {
				if sharpFlag {
					fmt.Fprint(s, "Call stack:")
				} else {
					fmt.Fprint(s, " | call stack:")
				}
				frames := runtime.CallersFrames(e.stack)
				for {
					fr, more := frames.Next()
					if !more {
						break
					}
					if sharpFlag {
						fmt.Fprint(s, "\n  ")
					} else {
						fmt.Fprint(s, " {")
					}
					fmt.Fprint(s, fr.Function, " @ ", trimGOPATH(fr.File), ":", fr.Line)
					if !sharpFlag {
						fmt.Fprint(s, "}")
					}
				}
			}
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}

func (e *DetailedError) AppendRelated(err ...error) *DetailedError {
	e.related = append(e.related, err...)
	return e
}

func New(s string) error {
	return errors.New(s)
}

func Newf(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func WithDetails(err error) *DetailedError {
	var file string
	var line int
	pc := make([]uintptr, 7)
	n := runtime.Callers(2, pc)
	if n > 0 {
		stack := runtime.CallersFrames(pc[:n])
		f, _ := stack.Next()
		file = trimGOPATH(f.File)
		line = f.Line
	}
	return &DetailedError{base: err, file: trimGOPATH(file), line: line, stack: pc[:n]}
}

func Base(err error) error {
	for err != nil {
		detailed, ok := err.(*DetailedError)
		if !ok {
			break
		}
		err = detailed.base
	}
	return err
}
