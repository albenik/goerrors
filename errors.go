package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type DetailedError struct {
	base    error
	file    string
	line    int
	related []error
}

func (e *DetailedError) Error() string {
	if len(e.related) == 0 {
		return fmt.Sprintf("%s @ %s:%d", e.base.Error(), e.file, e.line)
	}
	s := make([]string, 0, len(e.related))
	for _, e2 := range e.related {
		s = append(s, e2.Error())
	}
	return fmt.Sprintf("%s @ %s:%d with errors: [%s]", e.base.Error(), e.file, e.line, strings.Join(s, "; "))
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
	_, file, line, _ := runtime.Caller(1)
	return &DetailedError{base: err, file: trimGOPATH(file), line: line}
}

func Base(err error) error {
	if e, ok := err.(*DetailedError); ok {
		return e.base
	} else {
		return err
	}
}

func NativeBase(err error) error {
	for {
		if detailed, ok := err.(*DetailedError); ok {
			err = detailed.base
			continue
		}
		return err
	}
}
