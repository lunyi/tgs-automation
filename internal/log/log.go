package log

import (
	"fmt"
	"os"
	"tgs-automation/pkg/telegram"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("cdn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to log to file, using default stderr")
	}

	log.SetOutput(file)
}

func LogInfo(msg string) {
	fmt.Println("Info", msg)
	log.Info(msg)
	err := telegram.SendMessage("Info: " + msg)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}
}

func LogTrace(msg string) {
	fmt.Println("Trace", msg)
	log.Trace(msg)
}

func LogError(msg string) {
	fmt.Println("Error", msg)
	log.Error(msg)
	err := telegram.SendMessage("Error: " + msg)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}
}

func LogFatal(msg string) {
	fmt.Println("Fatal", msg)
	log.Fatal(msg)
	err := telegram.SendMessage("Fatal: " + msg)
	if err != nil {
		fmt.Println("Failed to send Telegram message:", err)
	}
}
