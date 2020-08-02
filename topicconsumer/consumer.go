package topics

import (
	"time"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	configuration "me/gitoperator/configuration"
	mediator "me/gitoperator/mediator"
	utils "me/gitoperator/utils"
	kafka "github.com/segmentio/kafka-go"
	//_ "github.com/segmentio/kafka-go/snappy"
)

const componentMessage = "Topics Consumer Service"
var config = configuration.GlobalConfiguration


func getKafkaReader() *kafka.Reader {
	broker := config.Kafka.Bootstrapserver
	brokers := strings.Split(broker, ",")
	groupID := config.Kafka.Groupid
	topic := config.Kafka.Consumertopic
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: config.Kafka.Messageminsize,
		MaxBytes: config.Kafka.Messagemaxsize,
		MaxWait: 100 * time.Millisecond,
	})
}

var ActiveReposList []string

func StartListening() {
	methodMsg := "StartListening"
	reader := getKafkaReader()
	ActiveReposList = make([]string, 3)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			utils.PrintLogError(err, componentMessage, methodMsg, "Error reading message")
		}
		msg := fmt.Sprintf("Message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		utils.PrintLogInfo(componentMessage, methodMsg, msg)
		event, eventErr := convertMessageToProcessable(m)
		if eventErr != nil {
			utils.PrintLogError(eventErr, componentMessage, methodMsg, fmt.Sprintf("Message convertion error - Key '%s'", m.Key))
			// send alert about not valid request
		} else {
			utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("Message converted to event successfully - Key '%s'", m.Key))
			mediator.ProcessSyncIncomingMessage(&event)
			/*
 			if !utils.Contains(activeReposList, event.DBName) {
				fmt.Println("FREE repo (not in list!!!" + event.DBName)
				ActiveReposList = append(activeReposList, event.DBName)
				go mediator.NotSynchronizedProcess(&event)
				activeReposList = utils.RemoveElementFromSlice(activeReposList, event.DBName)
			} else {
				mediator.ProcessSyncIncomingMessage(&event)
			} 
			*/
		}
	}
}


func convertMessageToProcessable(msg kafka.Message) (utils.RecordEvent, error) {
	methodMsg := "convertMessageToProcessable"
	var newRecordEvent utils.RecordEvent
	unmarshalErr := json.Unmarshal(msg.Value, &newRecordEvent)
	if unmarshalErr != nil {
		//utils.PrintLogError(unmarshalErr, componentMessage, methodMsg, fmt.Sprintf("Error unmarshaling message content to JSON - Key '%s'", msg.Key))
		return newRecordEvent, unmarshalErr
	}
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("ID '%s'", newRecordEvent.Id))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("Group '%s'", newRecordEvent.Group))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("OperationType '%s'", newRecordEvent.OperationType))
	return newRecordEvent, nil
}