package gitactors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGitProcessNewFile(t *testing.T) {
	Convey("Put new file in git ", t, func() {
		ev := getEvent()
		err := GitProcessNewFile(&ev)
		So(err, ShouldBeNil)
	})
}