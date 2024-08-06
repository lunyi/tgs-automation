package util

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsPublisherService interface {
	Publish(chatId string, message string) error
	Close()
}

type NatsPublisher struct {
	natsURL string
	nc      *nats.Conn
	js      nats.JetStreamContext
}

type NatsData struct {
	ChatId  string `json:"chat_id"`
	Message string `json:"message"`
}

func NewNatsPublisher(natsURL string) (NatsPublisherService, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS server: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("error creating JetStream context: %w", err)
	}

	return &NatsPublisher{
		natsURL: natsURL,
		nc:      nc,
		js:      js,
	}, nil
}

func (mp *NatsPublisher) Close() {
	mp.nc.Close()
}

func (mp *NatsPublisher) Publish(chatId string, message string) error {
	msgContent := NatsData{
		ChatId:  chatId,
		Message: message,
	}

	msgBytes, err := json.Marshal(msgContent)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	_, err = mp.js.Publish("telegram.messages", msgBytes)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	fmt.Printf("Message published: %+v\n", msgContent)
	return nil
}
