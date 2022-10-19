package giteaapiclient

import (
	"testing"
	utils "xqledger/gitoperator/utils"
	. "github.com/smartystreets/goconvey/convey"
)



func TestCreateFileInRepo(t *testing.T) {
	Convey("Should create a new record in Git  ", t, func() {
		event := utils.GetNewRecordEvent()
		err := CreateFileInRepo(&event)
		So(err, ShouldBeNil)
	})
}

func TestUpdateFileInRepo(t *testing.T) {
	Convey("Should update an existing record in Git  ", t, func() {
		event := utils.GetRecordEventToUpdate()
		err := UpdateFileInRepo(&event)
		So(err, ShouldBeNil)
	})
}

func TestDeleteFileInRepo(t *testing.T) {
	Convey("Should delete an existing record in Git  ", t, func() {
		event := utils.GetRecordEventToDelete()
		err := DeleteFileInRepo(&event)
		So(err, ShouldBeNil)
	})
}