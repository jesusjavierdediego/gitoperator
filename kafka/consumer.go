package kafka

import (
	"time"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	configuration "xqledger/gitoperator/configuration"
	mediator "xqledger/gitoperator/mediator"
	utils "xqledger/gitoperator/utils"
	kafka "github.com/segmentio/kafka-go"
)

const componentMessage = "Topics Consumer Service"
var config = configuration.GlobalConfiguration


func getKafkaReader(topic string) *kafka.Reader {
	broker := config.Kafka.Bootstrapserver
	brokers := strings.Split(broker, ",")
	groupID := config.Kafka.Groupid
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: config.Kafka.Messageminsize,
		MaxBytes: config.Kafka.Messagemaxsize,
		MaxWait: 100 * time.Millisecond,
	})
}

func StartListeningRecords() {
	methodMsg := "StartListeningRecords"
	reader := getKafkaReader(config.Kafka.Consumertopic)
	defer reader.Close()
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			utils.PrintLogError(err, componentMessage, methodMsg, fmt.Sprintf("%s - Error reading message", utils.Event_topic_received_fail))
		}
		// msg := fmt.Sprintf("Message received at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// utils.PrintLogInfo(componentMessage, methodMsg, msg)
		event, eventErr := convertMessageToProcessable(m)
		if eventErr != nil {
			utils.PrintLogError(eventErr, componentMessage, methodMsg, fmt.Sprintf("%s - Record event message convertion error - Key '%s'", utils.Event_topic_received_unacceptable, m.Key))
		} else {
			utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("%s - Message converted to record event successfully - Key '%s'", utils.Event_topic_received_ok, m.Key))
			mediator.ProcessSyncIncomingMessageRecord(&event)
		}
	}
}


func convertMessageToProcessable(msg kafka.Message) (utils.RecordEvent, error) {
	methodMsg := "convertMessageToProcessable"
	var newRecordEvent utils.RecordEvent
	unmarshalErr := json.Unmarshal(msg.Value, &newRecordEvent)
	if unmarshalErr != nil {
		utils.PrintLogWarn(unmarshalErr, componentMessage, methodMsg, fmt.Sprintf("Error unmarshaling message content to JSON - Key '%s'", msg.Key))
		return newRecordEvent, unmarshalErr
	}
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("ID '%s'", newRecordEvent.Id))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("DB Name '%s'", newRecordEvent.DBName))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("Operation Type '%s'", newRecordEvent.OperationType))
	return newRecordEvent, nil
}

func convertMessageToProcessableBatch(msg kafka.Message) (utils.RecordEventBatch, error) {
	methodMsg := "convertMessageToProcessable"
	var newRecordEventBatch utils.RecordEventBatch
	unmarshalErr := json.Unmarshal(msg.Value, &newRecordEventBatch)
	if unmarshalErr != nil {
		utils.PrintLogWarn(unmarshalErr, componentMessage, methodMsg, fmt.Sprintf("Error unmarshaling batch content to JSON - Key '%s'", msg.Key))
		return newRecordEventBatch, unmarshalErr
	}
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("ID '%s'", newRecordEventBatch.Id))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("DB Name '%s'", newRecordEventBatch.DBName))
	utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("Operation Type '%s'", newRecordEventBatch.OperationType))
	return newRecordEventBatch, nil
}