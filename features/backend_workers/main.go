package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/telegram"

	"github.com/nats-io/nats.go"
)

type TelegramMessage struct {
	ChatId  string `json:"chat_id"`
	Message string `json:"message"`
}

func main() {
	config := util.GetConfig()
	bot, err := telegram.GetTelegramBotInstance(config.Telegram.TelegramBotToken)
	if err != nil {
		log.Fatalln("Error not getting telegram token")
	}
	// cleanup := initTracer(config.NatsUrl)
	// defer cleanup()

	fmt.Println("config.NatsUrl:", config.NatsUrl)
	// 连接 NATS
	nc, err := nats.Connect(config.NatsUrl)

	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()
	fmt.Println("nats connected")

	// 創建 JetStream 上下文
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	log.Println("creating JetStream context OK")
	// 確保 Stream 存在
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "TELEGRAM",
		Subjects: []string{"telegram.messages"},
	})

	if err != nil {
		log.Fatalf("Error adding stream: %v", err)
	}

	log.Println("Add JetStream OK")

	// 訂閱消息
	sub, err := js.Subscribe("telegram.messages",
		func(m *nats.Msg) {
			log.Printf("Received a message: %s", string(m.Data))
			err := handleTelegramMessages(bot, m)
			if err != nil {
				log.Printf("Error handling message: %v", err)
				// 可以在這裡添加重試邏輯
			}
			// 確認消息
			if err := m.Ack(); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}
		},
		nats.Durable("telegram-durable"),
	)
	if err != nil {
		log.Fatalf("Error subscribing: %v", err)
	}
	defer sub.Unsubscribe()

	// 設置信號處理
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	fmt.Println("Subscriber is running. Press Ctrl+C to stop.")

	// 等待終止信號
	<-stop

	fmt.Println("Shutting down...")
}

func handleTelegramMessages(bot *telegram.TelegramBot, m *nats.Msg) error {
	// tracer := otel.Tracer("example.com/trace")
	// ctx, span := tracer.Start(ctx, "handleTelegramMessages")
	// defer span.End()
	fmt.Println("handleTelegramMessages entered")
	var msg TelegramMessage
	err := json.Unmarshal(m.Data, &msg)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	err = bot.SendMessage(msg.ChatId, msg.Message)

	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
		return err
	}
	return nil
}
