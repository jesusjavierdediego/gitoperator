package mediator

import (
	"errors"
	"sync"
	"encoding/json"
	"strconv"
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
	switch event.OperationType {
		case "new":
			gitErr = git.GitProcessNewFile(event)
		case "update":
			gitErr = git.GitUpdateFile(event)
		case "delete":
			gitErr = git.GitDeleteFile(event)
		default:
			gitErr = errors.New("Operation type not acceptable")
	}

	if gitErr != nil {
		utils.PrintLogError(gitErr, componentMessage, methodMessage, "Error processing in Git server - ID: "+event.Id)
		return
	}
	// Send update of record to topic
	cleanDbRecord, cleanErr := strconv.Unquote(event.RecordContent)
	if cleanErr != nil {
		utils.PrintLogError(cleanErr, componentMessage, methodMessage, "Error parsing record event payload event- ID: "+event.Id)
	}
	SendMessageToTopic(cleanDbRecord)
	utils.PrintLogInfo(componentMessage, methodMessage, fmt.Sprintf("DB Record in event succesfully sent to continuous query topic - Event ID '%s'", event.Id))

	event.Status = "COMPLETE"
	_, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		utils.PrintLogError(marshalErr, componentMessage, methodMessage, "Error parsing response event- ID: "+event.Id)
		return
	}
	
	/*
	// Once the git event has been processed properly we can send a notification to a topic
	msgBytes, err := json.Marshal(event)
    if err != nil {
		utils.PrintLogError(marshalErr, componentMessage, methodMessage, "Error serializing response event - ID: "+event.Id)
	}
	topicsender.SendMessageToTopic(string(msgBytes))
	*/
	utils.PrintLogInfo(componentMessage, methodMessage, "Event processed successfully - ID: "+event.Id)
	m.Unlock()
	wg.Done()
}