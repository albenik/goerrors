package errors_test

import (
	"fmt"
	"testing"

	"github.com/albenik/goerrors"
	"github.com/stretchr/testify/assert"
)

func TestDetailedError_Error(t *testing.T) {
	err1 := errors.New("error 1")
	err1w := errors.Wrap(err1, "error 1w")

	assert.Equal(t, err1, errors.Base(err1))
	assert.Equal(t, err1, errors.Base(err1w))

	fmt.Println(err1)
	fmt.Println(err1w)

	err2 := errors.New("error 2")
	err2w := errors.Wrap(err2, "error 2w")
	err2w.AppendRelated(errors.New("suberror 1"), errors.New("suberror 2"))

	fmt.Println(err2w)

	err3 := errors.New("error 3")
	err3w := errors.Wrap(err3, "error 3w").AppendRelated(errors.New("suberror 3.1"))
	err3w = errors.Wrap(err3w, "error 3w1").AppendRelated(errors.New("suberror 3.2"))
	err3w = errors.Wrap(err3w, "error 3w2").AppendRelated(errors.New("suberror 3.3"))
	err3w = errors.Wrap(err3w, "error 3w3").AppendRelated(errors.New("suberror 3.4"))

	assert.Equal(t, err3, errors.Base(err3w))

	fmt.Println(err3w)
}
