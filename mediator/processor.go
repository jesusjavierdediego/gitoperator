package mediator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	git "xqledger/gitoperator/gitactors"
	utils "xqledger/gitoperator/utils"
)

const componentMessage = "Processor"

func ProcessSyncIncomingMessageRecord(event *utils.RecordEvent) {
	var w sync.WaitGroup
	var m sync.Mutex
	w.Add(1)
	go synchronizedProcessRecord(&w, &m, event)
	w.Wait()
}

func ProcessSyncIncomingMessageBatch(batch *utils.RecordEventBatch) {
	var w sync.WaitGroup
	var m sync.Mutex
	w.Add(1)
	go synchronizedProcessBatch(&w, &m, batch)
	w.Wait()
}

func synchronizedProcessBatch(wg *sync.WaitGroup, m *sync.Mutex, batch *utils.RecordEventBatch) {
	methodMessage := "synchronizedProcessBatch"
	m.Lock()
	var batchErr error
	utils.PrintLogInfo(componentMessage, methodMessage, "batch.OperationType: "+batch.OperationType)
	var logMsgOk = ""
	var logMsgFail = ""
	switch batch.OperationType {
	case "new":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("New batch with ID %s", batch.Id))
		batchErr = git.GitProcessNewBatch(batch)
		logMsgOk = utils.Record_new_git_written_ok
		logMsgFail = utils.Record_new_git_written_fail
	case "update":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Update batch with ID %s", batch.Id))
		batchErr = git.GitUpdateFileBatch(batch)
		logMsgOk = utils.Record_update_git_written_ok
		logMsgFail = utils.Record_update_git_written_fail
	// case "delete":
	// 	utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Delete batch with ID %s", batch.Id))
	// 	batchErr = git.GitDeleteFileBatch(batch)
	// 	logMsgOk = utils.Record_delete_git_written_ok
	// 	logMsgFail = utils.Record_delete_git_written_fail
	default:
		batchErr = errors.New("Operation type not acceptable")
		logMsgFail = utils.Record_operation_not_accepted
	}
	if batchErr != nil {
		utils.PrintLogError(batchErr, componentMessage, methodMessage, logMsgFail)
		return
	}
	utils.PrintLogInfo(componentMessage, methodMessage, logMsgOk)

	
	var batchRecordsToRDBErr []error
	for _, event := range batch.Records {
		eventAsJSON, err := json.Marshal(event)
		if err != nil {
			utils.PrintLogError(err, componentMessage, methodMessage, fmt.Sprintf("Event in batch cannot be marshaled to be sent to the RDB - Event ID '%s'", event.Id))
		}
		sendErr := SendMessageToTopic(string(eventAsJSON), config.Kafka.Gitactionbacktopic)
		if sendErr != nil {
			utils.PrintLogError(sendErr, componentMessage, methodMessage, utils.Event_written_record_topic_send_fail)
			batchRecordsToRDBErr = append(batchRecordsToRDBErr, sendErr)
	
		}
		utils.PrintLogInfo(componentMessage, methodMessage, utils.Event_written_record_topic_send_ok)
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("DB Record in event succesfully sent to continuous query topic - batch ID '%s'", batch.Id))
	}
	
	if len(batchRecordsToRDBErr) > 0 {
		utils.PrintLogInfo(componentMessage, methodMessage, "Batch processed with errors sending records to RDB - ID: "+batch.Id)
	}
	m.Unlock()
	wg.Done()
}

func synchronizedProcessRecord(wg *sync.WaitGroup, m *sync.Mutex, event *utils.RecordEvent) {
	methodMessage := "synchronizedProcessRecord"
	m.Lock()

	var gitErr error
	utils.PrintLogInfo(componentMessage, methodMessage, "event.OperationType: "+event.OperationType)
	var logMsgOk = ""
	var logMsgFail = ""
	switch event.OperationType {
	case "new":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("New Event with ID %s", event.Id))
		gitErr = git.GitProcessNewFile(event)
		logMsgOk = utils.Record_new_git_written_ok
		logMsgFail = utils.Record_new_git_written_fail
	case "update":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Update Event with ID %s", event.Id))
		gitErr = git.GitUpdateFile(event)
		logMsgOk = utils.Record_update_git_written_ok
		logMsgFail = utils.Record_update_git_written_fail
	case "delete":
		utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("Delete Event with ID %s", event.Id))
		gitErr = git.GitDeleteFile(event)
		logMsgOk = utils.Record_delete_git_written_ok
		logMsgFail = utils.Record_delete_git_written_fail
	default:
		gitErr = errors.New("Operation type not acceptable")
		logMsgFail = utils.Record_operation_not_accepted
	}
	log.Println(fmt.Sprintf("Event iD: %s", event.Id))
	if gitErr != nil {
		utils.PrintLogError(gitErr, componentMessage, methodMessage, logMsgFail)
		return
	}
	utils.PrintLogInfo(componentMessage, methodMessage, logMsgOk)

	eventAsJSON, err := json.Marshal(event)
	if err != nil {
		utils.PrintLogError(err, componentMessage, methodMessage, "Event cannot be marshaled properly after written to Git")
	}
	sendErr := SendMessageToTopic(string(eventAsJSON), config.Kafka.Gitactionbacktopic)
	if sendErr != nil {
		utils.PrintLogError(sendErr, componentMessage, methodMessage, utils.Event_written_record_topic_send_fail)

	}
	utils.PrintLogInfo(componentMessage, methodMessage, utils.Event_written_record_topic_send_ok)
	utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("DB Record in event succesfully sent to continuous query topic - Event ID '%s'", event.Id))
	utils.PrintLogInfo(componentMessage, methodMessage, "Event processed successfully - ID: "+event.Id)
	m.Unlock()
	wg.Done()
}
