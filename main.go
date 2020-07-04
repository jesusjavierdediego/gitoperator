package main

import (

)

func main() {
	//var content = "{\"created_at\": \"Wed Jul 03 11:00:08 +0000 2020\",\"id\": 1275753321895211003}"

	local_repo_path := "/Users/administrator/workspace/projects/go/src/me/gitoperator/localrepo"
	// remote_repo_url := "http://localhost:3000/jdediego/GoNoSQL_1.git"
	// Clone(remote_repo_url, local_repo_path)

	var id = "1275753321895211007"
	var filename = id + ".json"
	var content = "{\"created_at\": \"Wed Mar 03 11:00:08 +0000 2020\",\"id\": " + id + "}"
	var newRecordEvent RecordEvent
	newRecordEvent.Id = id
	newRecordEvent.Group = ""
	newRecordEvent.Unit = "UsFirst"
	newRecordEvent.Priority = "HIGH"
	newRecordEvent.ReceptionTime = 0
	newRecordEvent.SendingTime = 0
	newRecordEvent.OperationType = "new"
	newRecordEvent.RecordContent = content

	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id
	updateRecordEvent.Group = ""
	updateRecordEvent.Unit = "UsFirst"
	updateRecordEvent.Priority = "HIGH"
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.RecordContent = ""
	var fields map[string]string
	fields = make(map[string]string)
	fields["created_at"] = "Wed Dec 23 11:00:08 +0000 2020"

	//NewFile(local_repo_path, filename, newRecordEvent)
	UpdateFile(local_repo_path, filename, updateRecordEvent)
	
	
}