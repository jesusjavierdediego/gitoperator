package utils

type RecordEvent struct {
	Id   string `json:"id"` // Name of the blob/file
	Group string `json:"group"` // Name of the tree/folder
	Unit string `json:"unit"` // Unit mapped to Git repo
	OperationType string `json:"operation_type"` // Values: (new | update | delete)
	SendingTime int64 `json:"sending_time"` // Time of sending by the client
	ReceptionTime int64 `json:"reception_time"` // Time of processing by the API
	Priority string `json:"priority"`  // API can qualify an event with a priority to be considered in concurrent writing decisions (HIGH | MEDIUM | LOW)
	Message string `json:"message"`
	RecordContent string `json:"record_content"` // empty if op OperationType == delete | update
}
