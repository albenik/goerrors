package errors_test

import (
	"fmt"
	"testing"

	"github.com/albenik/goerrors"
	"github.com/stretchr/testify/assert"
)

func TestDetailedError_Error(t *testing.T) {
	err1 := errors.New("error 1")
	err1d := errors.WithDetails(err1)

	assert.Equal(t, err1, errors.Base(err1))

	assert.Equal(t, err1, errors.Base(err1d))

	fmt.Println(err1)
	fmt.Println(err1d)

	err2 := errors.New("error 2")
	err2d := errors.WithDetails(err2)
	err2d.AppendRelated(errors.WithDetails(errors.New("suberror 1")), errors.New("suberror 2"))

	fmt.Println(err2d)

	err3 := errors.New("error 3")
	err3d := errors.WithDetails(err3).AppendRelated(errors.WithDetails(errors.New("suberror 3.1")))
	err3d = errors.WithDetails(err3d).AppendRelated(errors.WithDetails(errors.New("suberror 3.2")))
	err3d = errors.WithDetails(err3d).AppendRelated(errors.WithDetails(errors.New("suberror 3.3")))
	err3d = errors.WithDetails(err3d).AppendRelated(errors.WithDetails(errors.New("suberror 3.4")))

	assert.Equal(t, err3, errors.Base(err3d))

	fmt.Println(err3d)
}
