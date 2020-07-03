package main

import (

)

func main() {
	//var content = "{\"created_at\": \"Wed Jul 03 11:00:08 +0000 2020\",\"id\": 1275753321895211003}"

	//remote_repo_url := "http://localhost:3000/jdediego/GitAsNoSQL_1.git"
	local_repo_path := "/home/administrator/go/src/me/gitpoc/gitlocal"
	//local_file_path := "record4.json"
	//Clone(remote_repo_url, local_repo_path)
	//NewFile(local_repo_path, local_file_path)
	//UpdateFile(local_repo_path, local_file_path, content)

	var id = "1275753321895211007"
	var filename = id + ".json"
	var content = "{\"created_at\": \"Wed Mar 03 11:00:08 +0000 2020\",\"id\": " + id + "}"
	var event RecordEvent
	event.Id = id
	event.Group = ""
	event.Unit = "UsFirst"
	event.Priority = "HIGH"
	event.ReceptionTime = 0
	event.SendingTime = 0
	event.OperationType = "update"
	event.Content = content
	//NewFile(local_repo_path, filename, content)
	UpdateFile(local_repo_path, filename, content)
	
}