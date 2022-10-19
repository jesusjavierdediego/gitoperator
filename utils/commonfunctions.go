package utils

const tenant = "GitOperatorTestRepo"
const email = "TestOrchestrator@gmail.com"
const group = "tree1"
const priority = "HIGH"
const id1 = "personalrecord01"
const id2 = "personalrecord02"



func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func indexOf(element string, data []string) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
 }

func RemoveElementFromSlice(series []string, value string) []string{
	i := indexOf(value, series)
	if i != -1 {
		series[i] = series[len(series)-1] // Copy last element to index i.
		series[len(series)-1] = ""   // Erase last element (write zero value).
		series = series[:len(series)-1]   // Truncate slice.
	}
	return series
}

func GetNewRecordEvent() RecordEvent {
	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id1
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = ""
	updateRecordEvent.User = email
	updateRecordEvent.Priority = priority
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "new"
	updateRecordEvent.RecordContent = `{"name": "John", "age": 32, "city": "Cincinatti"}`
	return updateRecordEvent
}

func GetRecordEventToUpdate() RecordEvent {
	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id1
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = ""
	updateRecordEvent.User = email
	updateRecordEvent.Priority = priority
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.RecordContent = `{"name": "John", "age": 54, "city": "Cincinatti"}`
	return updateRecordEvent
}

func GetRecordEventToDelete() RecordEvent {
	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id1
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = ""
	updateRecordEvent.User = email
	updateRecordEvent.Priority = priority
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "delete"
	updateRecordEvent.RecordContent = ``
	return updateRecordEvent
}



func GetNewRecordEventLevel1() RecordEvent {
	var newRecordEvent RecordEvent
	newRecordEvent.Id = id2
	newRecordEvent.DBName = tenant
	newRecordEvent.Group = group
	newRecordEvent.User = email
	newRecordEvent.Priority = priority
	newRecordEvent.ReceptionTime = 0
	newRecordEvent.SendingTime = 0
	newRecordEvent.ProcessingTime = 0
	newRecordEvent.OperationType = "new"
	newRecordEvent.RecordContent = `{"name": "John", "age": 31, "city": "New York"}`
	return newRecordEvent
}


func getRecordEventLevel1ToUpdate() RecordEvent {
	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id2
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = group
	updateRecordEvent.User = email
	updateRecordEvent.Priority = priority
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "update"
	updateRecordEvent.RecordContent = `{"name": "John", "age": 32, "city": "Cincinatti"}`
	return updateRecordEvent
}

func GetRecordEventLevel1ToDelete() RecordEvent {
	var updateRecordEvent RecordEvent
	updateRecordEvent.Id = id2
	updateRecordEvent.DBName = tenant
	updateRecordEvent.Group = ""
	updateRecordEvent.User = email
	updateRecordEvent.Priority = priority
	updateRecordEvent.ReceptionTime = 0
	updateRecordEvent.SendingTime = 0
	updateRecordEvent.ProcessingTime = 0
	updateRecordEvent.OperationType = "delete"
	updateRecordEvent.RecordContent = ``
	return updateRecordEvent
}