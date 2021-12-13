package gitactors

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGitProcessUpdateFileBatch(t *testing.T) {
	Convey("Put update batch in git ", t, func() {
		b := getUpdateBatch()
		err := GitUpdateFileBatch(&b)
		So(err, ShouldBeNil)
	})
}