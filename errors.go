package errors

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
)

type DetailedError struct {
	base    error
	msg     string
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
		if e.base == nil {
			return fmt.Sprintf("%s @ %s:%d", e.msg, e.file, e.line)
		}
		return fmt.Sprintf("%s @ %s:%d: %v", e.msg, e.file, e.line, e.base)
	}
	return fmt.Sprintf("%s @ %s:%d: %v | %s", e.msg, e.file, e.line, e.base, e.relatedString())
}

func (e *DetailedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		plusFlag := s.Flag('+')
		sharpFlag := s.Flag('#')
		if plusFlag || sharpFlag {
			if e.base == nil {
				fmt.Fprintf(s, "%s @ %s:%d", e.msg, e.file, e.line)
			} else {
				fmt.Fprintf(s, "%s @ %s:%d: %s", e.msg, e.file, e.line, e.base)
			}
			if len(e.related) > 0 {
				if sharpFlag {
					fmt.Fprint(s, "\nRelated errors:")
					for _, err := range e.related {
						fmt.Fprint(s, "\n  ", err.Error())
					}
				} else {
					fmt.Fprint(s, " | ", e.relatedString())
				}
			}
			if sharpFlag {
				fmt.Fprint(s, "\nCall stack:")
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
			return
		}
		fallthrough
	case 's', 'q':
		if verb == 'q' {
			fmt.Fprintf(s, "%q", e.Error())
		} else {
			io.WriteString(s, e.Error())
		}
	}
}

func (e *DetailedError) AppendRelated(err ...error) *DetailedError {
	e.related = append(e.related, err...)
	return e
}

func newError(msg string) *DetailedError {
	var file string
	var line int
	pc := make([]uintptr, 7)
	n := runtime.Callers(3, pc)
	if n > 0 {
		stack := runtime.CallersFrames(pc[:n])
		f, _ := stack.Next()
		file = trimGOPATH(f.File)
		line = f.Line
	}
	return &DetailedError{msg: msg, file: file, line: line, stack: pc[:n]}
}

func New(msg string) error {
	return newError(msg)
}

func Newf(format string, a ...interface{}) error {
	return newError(fmt.Sprintf(format, a...))
}

func Wrap(err error, msg string) *DetailedError {
	e := newError(msg)
	e.base = err
	return e
}

func Base(err error) error {
	for err != nil {
		detailed, ok := err.(*DetailedError)
		if !ok || detailed.base == nil {
			break
		}
		err = detailed.base
	}
	return err
}
