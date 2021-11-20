package gitactors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGitProcessDeleteFile(t *testing.T) {
	Convey("Put new file in git ", t, func() {
		ev := getDeleteEvent()
		err := GitDeleteFile(&ev)
		So(err, ShouldBeNil)
	})
}