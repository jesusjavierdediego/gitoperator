package mediator

import (
	"fmt"
	"context"
	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	utils "xqledger/gitoperator/utils"
	configuration "xqledger/gitoperator/configuration"
)

var config = configuration.GlobalConfiguration


const componentProducerMessage = "Topics Producer Service"

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}


func SendMessageToTopic(msg string, topic string) error{
	methodMsg := "SendMessageToTopic"
	broker := config.Kafka.Bootstrapserver
	kafkaWriter := getKafkaWriter(broker, topic)
	utils.PrintLogInfo(componentProducerMessage, methodMsg, fmt.Sprintf("Message sent to topic '%s'", topic))
	defer kafkaWriter.Close()


	topicContent := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: []byte(msg),

	}

	err := kafkaWriter.WriteMessages(context.Background(), topicContent)
	if err != nil {
		utils.PrintLogError(err, componentProducerMessage, methodMsg, fmt.Sprintf("Error writing message to topic '%s'", topic))
		return err
	}
	utils.PrintLogInfo(componentProducerMessage, methodMsg, fmt.Sprintf("Message sent to topic '%s' successfully", topic))
	return nil
}
