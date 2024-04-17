package signalhandler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tgs-automation/internal/log"
)

// StartListening starts listening for OS signals to gracefully shutdown the application.
func StartListening() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v", sig)) // Remove the second argument 'sig'
	os.Exit(0)
}
