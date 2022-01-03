package utils



// DOMAIN         -----------> DBNAME -------> GIT REPOSITORY
// JSON STRUCTURE -----------> GROUP --------> GIT FOLDER

type RecordEventBatch struct {
	Id   string `json:"id"` // Name of the file/record in the database
	DBName string `json:"dbname"` // DB name mapped to Git repo
	OperationType string `json:"operation_type"` // Values: (new | update | delete)
	Records []RecordEvent `json:"records"` // Set of  record events
}

type RecordEvent struct {
	Id   string `json:"id"` // Name of the file/record in the database
	Session string `json:"session"` // Name of the Git branch
	Group string `json:"group"` // Name of the Git tree/folder
	DBName string `json:"dbname"` // DB name mapped to Git repo
	User string `json:"user"` // email of the individual performing the change
	OperationType string `json:"operation_type"` // Values: (new | update | delete)
	SendingTime int64 `json:"sending_time"` // Time of sending by the client
	ReceptionTime int64 `json:"reception_time"` // Time of the reception by the API
	ProcessingTime int64 `json:"processing_time"` // Time of processing by the Git Operator
	Priority string `json:"priority"`  // API can qualify an event with a priority to be considered in concurrent writing decisions (HIGH | MEDIUM | LOW)
	RecordContent string `json:"record_content"` // empty if op OperationType == delete | update
	Status string `json:"status"` // PENDING | NOTVALID | INCOMPLETE | COMPLETE
}

type ClassiffiedEventsSet struct {
	SyncEvents   []RecordEvent `json:"sync_events"`
	ParEvents   []RecordEvent `json:"par_events"`
}

