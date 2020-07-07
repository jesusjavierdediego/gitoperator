package main

import (

)

func main() {
	//local_repo_path := "/Users/administrator/workspace/projects/go/src/me/gitoperator/localrepo"
	//remote_repo_url := "http://localhost:3000/jdediego/GoNoSQL_1.git"
	local_repo_path := "/home/administrator/go/src/me/gitpoc/gitlocal"
	//Clone(remote_repo_url, local_repo_path)

	var id = "personalrecord11"
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
	/*
	RULES
	Content of record is well formed valid JSON
	It is minified.
	*/
	updateRecordEvent.RecordContent = `{"id":11,"name":"Irshad","department":"Supplies","rate":15.26,"designation":"Product Manager","address":{"city":"Delhi","state":"Delhi","country":"India"},"sons":["Omar","Rajiv","Fatimah"],"partners":[{"name":"Robert","surname":"Young","company":"CIA"}]}`
	

	//NewFile(local_repo_path, filename, generalRecord)
	

	// new record complete, JSON compliant, pretty printed
	// var prettyJSON bytes.Buffer
	// json.Indent(&prettyJSON, []byte(newRecord), "", "\t")
	// var prettyNewRecord = string(prettyJSON.Bytes())
	// fmt.Println(local_repo_path)
	// fmt.Println(filename)
	// fmt.Println(prettyNewRecord)
	//UpdateFile(local_repo_path, filename, string(prettyJSON.Bytes()))
	//updateRecordEvent.RecordContent = string(prettyJSON.Bytes())
	var fileName = id + ".json"
	UpdateFile(local_repo_path, fileName, updateRecordEvent)
	
	
}