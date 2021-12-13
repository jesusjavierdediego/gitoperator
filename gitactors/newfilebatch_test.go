package gitactors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGitProcessNewFileBatch(t *testing.T) {
	Convey("Put new batch in git ", t, func() {
		b := getBatch()
		err := GitProcessNewBatch(&b)
		So(err, ShouldBeNil)
	})
}