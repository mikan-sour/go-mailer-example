package service

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

type KafkaListener interface {
	SetupConsumer([]string) (sarama.Consumer, error)
	ParseMessage([]byte) (*Message, error)
}

type KafkaListenerImpl struct {
	Ports []string
}

func NewKafkaListenerService(ports []string) *KafkaListenerImpl {
	return &KafkaListenerImpl{Ports: ports}
}

func (k *KafkaListenerImpl) SetupConsumer([]string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	conn, err := sarama.NewConsumer(k.Ports, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (k *KafkaListenerImpl) ParseMessage(msg []byte) (*Message, error) {
	var message *Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %s", err.Error())
	}
	return message, nil
}
