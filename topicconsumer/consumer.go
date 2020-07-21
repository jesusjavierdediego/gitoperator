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
	_ "github.com/segmentio/kafka-go/snappy"
)

const componentMessage = "Topics Consumer Service"
var config = configuration.GlobalConfiguration

func getBatchReader() (*kafka.Batch, error) {
	partition := 0
	conn, connErr := kafka.DialLeader(context.Background(), "tcp", config.Kafka.Bootstrapserver, config.Kafka.Consumertopic, partition)
	if connErr != nil {
		return nil, connErr
	}
	conn.SetReadDeadline(time.Now().Add(5*time.Second))
	return conn.ReadBatch(1024, 8192), nil
}

func StartListeningBatches() error{
	methodMsg := "StartListeningBatches"
	batch, readerErr := getBatchReader()
	if readerErr != nil {
		utils.PrintLogError(readerErr, componentMessage, methodMsg, "Error connecting to Kafka cluster")
		return readerErr
	}
	defer batch.Close()

	for {
		//Iterate the btach, store all messages in a slice, send to classify and process
		m, readErr := batch.ReadMessage()
		if readErr != nil {
			utils.PrintLogError(readErr, componentMessage, methodMsg, "EOB")
			batch.Close()
		}
		msg := fmt.Sprintf("Message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		utils.PrintLogInfo(componentMessage, methodMsg, msg)
		event, eventErr := convertMessageToProcessable(m)
		if eventErr != nil {
			utils.PrintLogError(eventErr, componentMessage, methodMsg, fmt.Sprintf("Message convertion error - Key '%s'", m.Key))
			return eventErr
			// send alert about not valid request
		}
		utils.PrintLogInfo(componentMessage, methodMsg, fmt.Sprintf("Message converted to event successfully - Key '%s'", m.Key))
		mediator.ProcessIncomingMessage(&event)
			//EventsQueue = append(EventsQueue, event)
		
		/*
		m := make([]byte, 10e3)
		_, err := batchReader.Read(b)
		if err == nil {
			m := string(b)
			cleanM := strings.Replace(m, "\x00", "", -1)
			if len(cleanM) > 0 {
				var event utils.RecordEvent
				//message := strings.Replace(m, "\\", "", -1)
				fmt.Printf("%#v\n", cleanM)
				utils.PrintLogInfo(componentMessage, methodMsg, cleanM)
				if unmarshalErr := json.Unmarshal(b, &event); unmarshalErr != nil {
					utils.PrintLogError(unmarshalErr, componentMessage, methodMsg, "Message convertion error")
				} else {
					utils.PrintLogInfo(componentMessage, methodMsg, "Message unmarshaled to event successfully - Event id: " + event.Id)
					mediator.ProcessIncomingMessage(&event)
				}
			}
		}*/
	}
	return nil
}

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

func StartListening() {
	methodMsg := "StartListening"
	reader := getKafkaReader()

	defer reader.Close()
	//utils.PrintLogInfo(componentMessage, methodMsg, "Start consuming ... !!")
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
			mediator.ProcessIncomingMessage(&event)
			//EventsQueue = append(EventsQueue, event)
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