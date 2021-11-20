package gitactors

import (
	"log"
	"testing"
	utils "xqledger/gitoperator/utils"

	. "github.com/smartystreets/goconvey/convey"
)
 
const repo = "GitOperatorTestRepo"
const id = "123456789123456789123456"
const email = "testorchestrator@gmail.com"
const recordTime = int64(1636570869)

func getEvent()utils.RecordEvent{
	record := utils.RecordEvent{}
	record.Id = id
	record.Group = ""
	record.DBName = repo
	record.User = email
	record.OperationType = "new"
	record.SendingTime = recordTime
	record.ReceptionTime = recordTime
	record.ProcessingTime = recordTime
	record.Priority = "MEDIUM"
	record.RecordContent = "{\"browsers\":{\"firefox\":{\"name\":\"Firefox\",\"pref_url\":\"about:config\",\"releases\":{\"1\":{\"release_date\":\"2004-11-09\",\"status\":\"retired\",\"engine\":\"Gecko\",\"engine_version\":\"1.7\"}}}}}"
	record.Status = "PENDING"
	return record
}

func getUpdateEvent()utils.RecordEvent{
	record := utils.RecordEvent{}
	record.Id = id
	record.Group = ""
	record.DBName = repo
	record.User = email
	record.OperationType = "update"
	record.SendingTime = recordTime
	record.ReceptionTime = recordTime
	record.ProcessingTime = recordTime
	record.Priority = "MEDIUM"
	record.RecordContent = "{\"browsers\":{\"firefox\":{\"name\":\"Firefox\",\"pref_url\":\"about:config\",\"releases\":{\"1\":{\"release_date\":\"2004-12-23\",\"status\":\"retired\",\"engine\":\"Gecko\",\"engine_version\":\"1.8\"}}}}}"
	record.Status = "PENDING"
	return record
}

func getDeleteEvent()utils.RecordEvent{
	record := utils.RecordEvent{}
	record.Id = id
	record.Group = ""
	record.DBName = repo
	record.User = email
	record.OperationType = "delete"
	record.SendingTime = recordTime
	record.ReceptionTime = recordTime
	record.ProcessingTime = recordTime
	record.Priority = "MEDIUM"
	record.Status = "PENDING"
	return record
}

func TestGetLocalRepoPath(t *testing.T) {
	Convey("Get local repo path ", t, func() {
		ev := getEvent()
		path, err := GetLocalRepoPath(&ev)
		log.Println(path)
		So(err, ShouldBeNil)
		So(len(path), ShouldBeGreaterThan, 0)
	})
}