package util

import (
	"os"
	"os/signal"
	"syscall"
	"tgs-automation/internal/log"
)

func HandleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.LogInfo("Shutting down due to signal")
	os.Exit(0)
}
