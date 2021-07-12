package mediator

import (
	"errors"
	"sync"
	"encoding/json"
	"fmt"
	git "xqledger/gitoperator/gitactors"
	utils "xqledger/gitoperator/utils"
)

const componentMessage = "Processor"


func ProcessSyncIncomingMessage(event *utils.RecordEvent) {
	var w sync.WaitGroup
	var m sync.Mutex
	w.Add(1)
	go synchronizedProcess(&w, &m, event)
	w.Wait()
} 

func synchronizedProcess(wg *sync.WaitGroup, m *sync.Mutex, event *utils.RecordEvent) {
	methodMessage := "synchronizedProcess"
	m.Lock()
	var gitErr error
	utils.PrintLogInfo(componentMessage, methodMessage, "event.OperationType: "+event.OperationType)
	var logMsgOk = ""
	var logMsgFail = ""
	switch event.OperationType {
		case "new":
			gitErr = git.GitProcessNewFile(event)
			logMsgOk = utils.Record_new_git_written_ok
			logMsgFail = utils.Record_new_git_written_fail
		case "update":
			gitErr = git.GitUpdateFile(event)
			logMsgOk = utils.Record_update_git_written_ok
			logMsgFail = utils.Record_update_git_written_fail
		case "delete":
			gitErr = git.GitDeleteFile(event)
			logMsgOk = utils.Record_delete_git_written_ok
			logMsgFail = utils.Record_delete_git_written_fail
		default:
			gitErr = errors.New("Operation type not acceptable")
			logMsgOk = utils.Record_new_git_written_ok
			logMsgFail = utils.Record_new_git_written_fail
	}

	if gitErr != nil {
		utils.PrintLogError(gitErr, componentMessage, methodMessage, fmt.Sprintf("% - Error processing in Git server - ID: %s", logMsgFail, event.Id))
		return
	}
	utils.PrintLogInfo(componentMessage, methodMessage, logMsgOk)

	// Send update of record to topic
	// cleanDbRecord, cleanErr := strconv.Unquote(event.RecordContent)
	// if cleanErr != nil {
	// 	utils.PrintLogError(cleanErr, componentMessage, methodMessage, "Error parsing record event payload event- ID: "+event.Id)
	// }
	// sendErr := SendMessageToTopic(cleanDbRecord, config.Kafka.Gitactionbacktopic)
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