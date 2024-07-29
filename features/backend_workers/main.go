package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tgs-automation/internal/opentelemetry"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/telegram"
	"time"

	"github.com/nats-io/nats.go"
)

type MessageData struct {
	ChatId  string `json:"chatid"`
	Message string `json:"message"`
}

func main() {
	config := util.GetConfig()
	ctx := context.Background()
	tp := opentelemetry.InitTracerProvider(ctx, config.JaegerCollectorUrl, "domain-api", "0.1.0", "prod")
	defer tp.Shutdown(ctx)

	nc, err := nats.Connect(config.NatsUrl)

	if err != nil {
		fmt.Println("nats connect error: %w", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	// 定義一個 Stream
	streamName := "telegram"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"telegram.*"},
	})
	if err != nil {
		log.Fatalf("Error adding stream: %v", err)
	}

	// 訂閱消息並啟用手動確認模式
	sub, err := js.SubscribeSync("telegram.*", nats.AckExplicit())
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	// 接收消息
	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		log.Fatalf("Error receiving message: %v", err)
	}
	fmt.Printf("Received message: %s\n", msg.Data)

	// 確認消息
	if err := msg.Ack(); err != nil {
		log.Fatalf("Error acknowledging message: %v", err)
	}
	fmt.Println("Message acknowledged")

	traceID := msg.Header.Get("trace-id")
	log.Printf("Received message with Trace ID: %s", traceID)

	// 解析 JSON 数据
	var receivedData MessageData
	if err := json.Unmarshal(msg.Data, &receivedData); err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}
	fmt.Printf("Parsed message data: %+v\n", receivedData)

	// 檢查未確認的消息
	pending, bytes, err := sub.Pending()
	if err != nil {
		log.Fatalf("Error checking pending messages: %v", err)
	}
	fmt.Printf("Pending messages: %d, Pending bytes: %d\n", pending, bytes)

	err = telegram.SendMessageWithChatIdAndContext(ctx, receivedData.Message, receivedData.ChatId)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}
}
