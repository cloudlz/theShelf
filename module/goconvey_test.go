//For the unit test
//1. file name must be the end of the _test.go, so that
// the corresponding code will be executed when the go test is executed
//2. you have to import testing this bag
//3. all test case functions must be Test at the beginning,
// and the test sessions will be executed sequentially in the order of source code
//4. test format: func TestXxx (t *testing.T),
// Xxx part can be any combination of letters and numbers,
// but the initial letter can not be lowercase letter [a-z],
// for example, Testintdiv is the wrong function name.
//5. function by calling testing.T Error, Errorf, FailNow, Fatal, FatalIf method,
// indicating that the test does not pass, call the Log method used to record the test information.

package main

//GoConvey is a unit testing framework
//It is a awesome Go Testing
//Official document https://github.com/smartystreets/goconvey
import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

//GoConvey simple demo
func TestSimpleDemo(t *testing.T) {
	//function by calling testing.T Error, Errorf, FailNow, Fatal, FatalIf method,
	// indicating that the test does not pass, call the Log method used to record the test information.

	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				//normal value to judge convey items
				//so(getvalue,assertions,expected)
				So(x, ShouldEqual, 2)
			})
		})
	})
}

//function use to TestAdvanceDemo
func diveAndHandle(a int, b int) (int, error) {
	if b == 0 {
		return 000, fmt.Errorf("Divisor can not be zero")
	}
	return a / b, nil
}

//GoConvey advance demo invoking function
// use common assertions as much as possible in assertion lib
func TestAdvanceDemo(t *testing.T) {
	//function by calling testing.T Error, Errorf, FailNow, Fatal, FatalIf method,
	// indicating that the test does not pass, call the Log method used to record the test information.

	resultInt, err := diveAndHandle(4, 2)
	Convey("we will get a result if dived number if dive nummber is true", t, func() {
		//use assertion lib to judge result and expect
		So(resultInt, ShouldNotBeNil)
		So(err, ShouldBeNil)
		Convey("the result that we got is true", func() {
			//normal value to judge convey items
			//so(getvalue,assertions,expected)
			So(resultInt, ShouldEqual, 2)
		})
	})
}
