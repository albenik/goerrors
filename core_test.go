package errors

import (
	"fmt"
	"testing"
)

func TestDetailed_CausedBy(t *testing.T) {
	err := New("error")
	fmt.Println(err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n\n", err)

	err = New("error").WithAnnotation("annotation")
	fmt.Println(err)
	fmt.Printf("%+v\n\n", err)

	err = New("E1").CausedBy(New("E2")).CausedBy(New("E3"))
	fmt.Println(err)
	fmt.Printf("%%s: %s\n", err)
	fmt.Printf("%%q: %q\n", err)
	fmt.Printf("%%v: %v\n", err)
	fmt.Printf("%%+v: %+v\n", err)
	fmt.Printf("%%#v: %#v\n", err)
}
