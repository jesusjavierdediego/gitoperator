package mediator

import (
	"sync"
	"encoding/json"
	git "me/gitoperator/git"
	utils "me/gitoperator/utils"
	topicsender "me/gitoperator/topicsender"
)

const componentMessage = "Processor"

// https://golangbot.com/mutex/
/*
func ProcessIncomingMessage(event *utils.RecordEvent) {
	
	switch event.OperationType {
		case "new":
			go git.GitProcessNewFile(event)
		case "update":
			var w sync.WaitGroup
			var m sync.Mutex
			w.Add(1)
			go synchronizedProcess(&w, &m, event)
			w.Wait()
		case "delete":
			go git.GitDeleteFile(event)
	}
	
	
}

func synchronizedProcess(wg *sync.WaitGroup, m *sync.Mutex, event *utils.RecordEvent) {
	methodMessage := "synchronizedProcess"
	m.Lock()
	var gitErr error
	utils.PrintLogInfo(componentMessage, methodMessage, "event.OperationType: "+event.OperationType)
	gitErr = git.GitUpdateFile(event)
	m.Unlock()
	wg.Done()

	if gitErr != nil {
		utils.PrintLogError(gitErr, componentMessage, methodMessage, "Error processing in Git server - ID: "+event.Id)
		return
	}

	event.Status = "COMPLETE"
	_, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		utils.PrintLogError(marshalErr, componentMessage, methodMessage, "Error parsing response event- ID: "+event.Id)
		return
	}
	utils.PrintLogInfo(componentMessage, methodMessage, "Event processed successfully - ID: "+event.Id)

} */

/* func ProcessIncomingMessages(events *[]utils.RecordEvent) {
	for _, event := range events {
		go git.GitProcessNewFile(&event)
	}
}  */

/* func ProcessMicroBatch(events []utils.RecordEvent) {
	methodMessage := "ProcessMicroBatch"
	for _, event := range events {
		utils.PrintLogInfo(componentMessage, methodMessage, "Record eevent received - Operation: "+event.OperationType)
		switch event.OperationType {
			case "new":
				go git.GitProcessNewFile(&event)
			case "update":
				go git.GitUpdateFile(&event)
			case "delete":
				go git.GitDeleteFile(&event)
		}
	}
} */


func ProcessIncomingMessage(event *utils.RecordEvent) {
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
	}

	if gitErr != nil {
		utils.PrintLogError(gitErr, componentMessage, methodMessage, "Error processing in Git server - ID: "+event.Id)
		return
	}

	event.Status = "COMPLETE"
	_, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		utils.PrintLogError(marshalErr, componentMessage, methodMessage, "Error parsing response event- ID: "+event.Id)
		return
	}
	
	msgBytes, err := json.Marshal(event)
    if err != nil {
		utils.PrintLogError(marshalErr, componentMessage, methodMessage, "Error serializing response event - ID: "+event.Id)
	}
	// Once the git event has been processes properly is sent to the topic to update the RDB
	topicsender.SendMessageToTopic(string(msgBytes))
	utils.PrintLogInfo(componentMessage, methodMessage, "Event processed nad returned as response event successfully - ID: "+event.Id)
	m.Unlock()
	wg.Done()
}
