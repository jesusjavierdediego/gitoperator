package topics

import (
	"fmt"
	"strings"
	"context"
	"encoding/json"
	kafka "github.com/segmentio/kafka-go"
	configuration "me/gitpoc/configuration"
	utils "me/gitpoc/utils"
	mediator "me/gitpoc/mediator"
)

const componentConsumerMessage = "Topics Consumer Service"
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
	})
}


func StartListening() error{
	methodMsg := "StartListening"
	reader := getKafkaReader()
	defer reader.Close()
	utils.PrintLogInfo(componentConsumerMessage, methodMsg, "Start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			utils.PrintLogError(err, componentConsumerMessage, methodMsg, "")
		}
		msg := fmt.Sprintf("Message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		utils.PrintLogInfo(componentConsumerMessage, methodMsg, msg)
		event, eventErr := convertMessageToProcessable(m)
		if eventErr != nil {
			utils.PrintLogError(eventErr, componentConsumerMessage, methodMsg, fmt.Sprintf("Message convertion error - Key '%s'", m.Key))
		}
		utils.PrintLogInfo(componentConsumerMessage, methodMsg, fmt.Sprintf("Message converted to event successfully - Key '%s'", m.Key))
		mediator.ProcessIncomingMessage(event)
	}
}

func convertMessageToProcessable(msg kafka.Message) (utils.RecordEvent, error){
	methodMsg := "convertMessageToProcessable"
	var newRecordEvent utils.RecordEvent
	unmarshalErr := json.Unmarshal(msg.Value, &newRecordEvent)
	if unmarshalErr != nil {
		utils.PrintLogError(unmarshalErr, componentConsumerMessage, methodMsg, fmt.Sprintf("Error unmarshaling message content to JSON - Key '%s'", msg.Key))
		return newRecordEvent, unmarshalErr
	}
	return newRecordEvent, nil
}
