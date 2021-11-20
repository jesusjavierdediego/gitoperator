package gitactors

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCloneRepo(t *testing.T) {
	Convey("Get local repo path ", t, func() {
		remote_repo_url := "http://localhost:3000/TestOrchestrator/GitOperatorTestRepo"
		local_repo_path := "/var/git/repos/GitOperatorTestRepo"
		err := Clone(remote_repo_url, local_repo_path)
		So(err, ShouldBeNil)
	})
}