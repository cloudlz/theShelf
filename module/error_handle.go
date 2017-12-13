//errors is a error handle package
//all the error can get the stack info
// if only use the method in errors
package main

import (
	"github.com/pkg/errors"
	"fmt"
	"io"
)

//we can new an error only has message
func newError(msg string) error {
	newErr := errors.New(msg)
	fmt.Println("--------------------newError")
	fmt.Printf("%+v", newErr)
	fmt.Println()
	return newErr
}

//we can wrap error by add some message
func wrapError(err error) error {
	wrapedErr := errors.Wrap(err, "It is wraped")
	fmt.Println("--------------------wrapError")
	fmt.Printf("%+v", wrapedErr)
	fmt.Println()
	return wrapedErr
}

//we can give the stack info
// by use it
func errorWithStack(err error) error {
	withStackErr := errors.WithStack(err)
	fmt.Println("--------------------errorWithStack")
	fmt.Printf("%+v", withStackErr)
	fmt.Println()
	return withStackErr
}

func main() {
	err := io.EOF
	str := "new error"
	//new an error only by message
	newError(str)
	//add some message to an error
	wrapError(err)
	//give the error its stack info
	errorWithStack(err)
}
