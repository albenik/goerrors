package errors

import (
	"fmt"
	"testing"
)

func p(err error) {
	fmt.Println(err)
	fmt.Printf("%s\n", err)
	fmt.Printf("%q\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n\n", err)
}

func TestDetailed_CausedBy(t *testing.T) {
	err1 := New("error 1")
	p(err1)

	err2 := New("error 2").WithAnnotation("annotation 2")
	p(err2)

	err3 := New("error 3").WithAnnotation("annotation 3").CausedBy(err2)
	p(err3)

	err4 := New("error 4").WithAnnotation("annotaion 4").CausedBy(err3).CausedBy(err1)
	p(err4)
}
