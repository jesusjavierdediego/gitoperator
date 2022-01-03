package mediator

import (
	"errors"
	"fmt"
	"log"
	api "xqledger/gitoperator/giteaapiclient"
	utils "xqledger/gitoperator/utils"
)


func ProcessConcurrentncomingMessageRecord(event *utils.RecordEvent) {
	utils.PrintLogInfo(componentMessage, "ProcessAsyncIncomingMessageRecord", fmt.Sprintf("Operation in session '%s'", event.Session))
	go concurrentProcessRecord(event)
}

func concurrentProcessRecord(event *utils.RecordEvent) {
	methodMessage := "concurrentProcessRecord"
	var apiErr error
	var logMsgFail = ""
	utils.PrintLogInfo(componentMessage, methodMessage, "event.OperationType: "+event.OperationType)
	switch event.OperationType {
	case "new":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("New Event with ID %s", event.Id))
		api.CreateFileInRepo(event)
	case "update":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Update Event with ID %s", event.Id))
		api.UpdateFileInRepo(event)
	case "delete":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Delete Event with ID %s", event.Id))
		api.DeleteFileInRepo(event)
	default:
		apiErr = errors.New("Operation type not acceptable")
		logMsgFail = utils.Record_operation_not_accepted
	}
	log.Println(fmt.Sprintf("Event iD: %s", event.Id))
	if len(logMsgFail) > 0 {
		utils.PrintLogError(apiErr, componentMessage, methodMessage, logMsgFail)
		return
	}
	utils.PrintLogInfo(componentMessage, methodMessage, "Operation processed")
}


