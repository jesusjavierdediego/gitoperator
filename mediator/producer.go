package mediator

import (
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

func SendMessageToTopic(msg string) error{
	methodMsg := "SendMessageToTopic"
	broker := config.Kafka.Bootstrapserver
	topic := config.Kafka.Gitactionbacktopic
	kafkaWriter := getKafkaWriter(broker, topic)

	defer kafkaWriter.Close()

	topicContent := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: []byte(msg),
	}

	err := kafkaWriter.WriteMessages(context.Background(), topicContent)
	if err != nil {
		utils.PrintLogError(err, componentProducerMessage, methodMsg, "Error writing message to topic")
		return err
	}
	utils.PrintLogInfo(componentProducerMessage, methodMsg, "Message sent to topic successfully")
	return nil
}