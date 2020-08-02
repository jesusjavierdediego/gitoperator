package gitactors

import (
	utils "me/gitoperator/utils"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
	RULES
	Content of record is well-formed && valid && minified JSON
*/

func getNewRecordEventLevel0() utils.RecordEvent {
	var id = "personalrecord11"
	var newRecordEvent utils.RecordEvent
	newRecordEvent.Id = id
	newRecordEvent.Group = ""
	newRecordEvent.Unit = "UsFirst"
	newRecordEvent.Priority = "HIGH"
	newRecordEvent.ReceptionTime = 0
	newRecordEvent.SendingTime = 0
	newRecordEvent.ProcessingTime = 0
	newRecordEvent.OperationType = "new"
	newRecordEvent.Message = "New file test level 0"
	newRecordEvent.RecordContent = `{"id":11,"name":"Irshad","department":"IT","status":"PENDING", "rate":14.098,"designation":"Secretary","address":{"city":"Delhi","state":"Delhi","country":"India"},"sons":["Omar","Rajiv","Fatimah"],"partners":[{"name":"Robert","surname":"Young","company":"CIA"}]}`
	return newRecordEvent
}

func getNewRecordEventLevel1() utils.RecordEvent {
	var id = "personalrecord01"
	var newRecordEvent utils.RecordEvent
	newRecordEvent.Id = id
	newRecordEvent.Group = "tree1"
	newRecordEvent.Unit = "UsFirst"
	newRecordEvent.Priority = "HIGH"
	newRecordEvent.ReceptionTime = 0
	newRecordEvent.SendingTime = 0
	newRecordEvent.ProcessingTime = 0
	newRecordEvent.OperationType = "new"
	newRecordEvent.Message = "New file test level 1"
	newRecordEvent.RecordContent = `{"name": "John", "age": 31, "city": "New York"}`
	return newRecordEvent
}

func getExistingRecordEventLevel0() utils.RecordEvent {
	var updateRecordEvent utils.RecordEvent
	var id = "personalrecord11"
	updateRecordEvent.Id = id
	updateRecordEvent.Group = ""
	updateRecordEvent.Unit = "UsFirst"
	updateRecordEvent.Priority = "HIGH"
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.Message = "Update test level 0"
	updateRecordEvent.RecordContent = `{"id":11,"name":"Irshad","department":"Supplies","status":"PENDING", "rate":15.26,"designation":"Product Manager","address":{"city":"Delhi","state":"Delhi","country":"India"},"sons":["Omar","Rajiv","Fatimah"],"partners":[{"name":"Robert","surname":"Young","company":"CIA"}]}`
	return updateRecordEvent
}

func getExistingRecordEventLevel1() utils.RecordEvent {
	var updateRecordEvent utils.RecordEvent
	var id = "personalrecord01"
	updateRecordEvent.Id = id
	updateRecordEvent.Group = "tree1"
	updateRecordEvent.Unit = "UsFirst"
	updateRecordEvent.Priority = "HIGH"
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.Message = "Update test level 1"
	updateRecordEvent.RecordContent = `{"name": "John", "age": 32, "city": "Cincinatti"}`
	return updateRecordEvent
}

func TestNewFile(t *testing.T) {
	/* Convey("Creating new file in Git level 0", t, func() {
		err := GitProcessNewFile(getNewRecordEventLevel0())
		So(err, ShouldBeNil)
	}) */
	Convey("Creating new file in Git level 1", t, func() {
		event := getNewRecordEventLevel1()
		err := GitProcessNewFile(&event)
		So(err, ShouldBeNil)
	})
}

func TestUpdateFile(t *testing.T) {
	/* Convey("Update file in Git level 0", t, func() {
		err := GitUpdateFile(getExistingRecordEventLevel0())
		So(err, ShouldBeNil)
	}) */
	Convey("Update file in Git level 1", t, func() {
		event := getExistingRecordEventLevel1()
		err := GitUpdateFile(&event)
		So(err, ShouldBeNil)
	})
}
