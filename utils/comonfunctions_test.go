package utils


import (	
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const value1 = "value1"
var array = []string {value1, "value2"}
func TestContains(t *testing.T) {
	Convey("Check if an array contains a string ", t, func() {
		result := Contains(array, value1)
		So(result, ShouldBeTrue)
	})
	Convey("Check if an array does not contain a string ", t, func() {
		result := Contains(array, "value3")
		So(result, ShouldBeFalse)
	})
}

func TestRemoveElementFromSlice(t *testing.T) {
	Convey("Check removing a string from an slice ", t, func() {
		result := RemoveElementFromSlice(array, value1)
		So(len(result), ShouldEqual, 1)
	})
	Convey("Check removing a string from an slice ", t, func() {
		result := RemoveElementFromSlice(array, "value3")
		So(len(result), ShouldEqual, 2)
	})
}
