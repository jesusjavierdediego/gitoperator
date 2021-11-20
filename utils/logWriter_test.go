package utils

import (
	"errors"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const component = "test"
const phase = "fake phase"
const message = "fake message"


func TestPrintLogError(t *testing.T) {
	Convey("Sends a formatted error log  ", t, func() {
		err := errors.New("fake error")
		result := PrintLogError(err, component, phase, message)
		So(result, ShouldBeTrue)
	})
}

func TestPrintLogWarn(t *testing.T) {
	Convey("Sends a formatted warn log  ", t, func() {
		err := errors.New("fake error")
		result := PrintLogWarn(err, component, phase, message)
		So(result, ShouldBeTrue)
	})
}

func TestPrintLogInfo(t *testing.T) {
	Convey("Sends a formatted info log  ", t, func() {
		result := PrintLogInfo(component, phase, message)
		So(result, ShouldBeTrue)
	})
}