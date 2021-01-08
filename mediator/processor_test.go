package mediator

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	utils "xqledger/gitoperator/utils"
	//configuration "me/gitoperator/configuration"
)

const tenant = "testTenant"

func getNewRecordEventLevel0() utils.RecordEvent {
	var id = "personalrecord11"
	var newRecordEvent utils.RecordEvent
	newRecordEvent.Id = id
	newRecordEvent.DBName = tenant
	newRecordEvent.Group = ""
	newRecordEvent.User = "test@test.com"
	newRecordEvent.Priority = "HIGH"
	newRecordEvent.ReceptionTime =  0
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
	newRecordEvent.DBName = tenant
	newRecordEvent.Group = "tree1"
	newRecordEvent.User = "test@test.com"
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
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = ""
	updateRecordEvent.User = "test@test.com"
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
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = "tree1"
	updateRecordEvent.User = "test@test.com"
	updateRecordEvent.Priority = "HIGH"
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.Message = "Update test level 1"
	updateRecordEvent.RecordContent = `{"name": "John", "age": 32, "city": "Cincinatti"}`
	return updateRecordEvent
}

func TestProcessSyncNewIncomingMessage(t *testing.T) {

	Convey("Processes a new incoming message ", t, func() {
		//TODO
	})
}

func TestProcessSyncIUpdateIncomingMessage(t *testing.T) {

	Convey("Processes an update incoming message ", t, func() {
		//TODO
	})
}

func TestProcessSyncDeleteIncomingMessage(t *testing.T) {

	Convey("Processes a delete incoming message ", t, func() {
		//TODO
	})
}