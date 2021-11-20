package gitactors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGitProcessUpdatedFile(t *testing.T) {
	Convey("Put new file in git ", t, func() {
		ev := getUpdateEvent()
		err := GitUpdateFile(&ev)
		So(err, ShouldBeNil)
	})
}