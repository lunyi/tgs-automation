package signalhandler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tgs-automation/internal/log"
)

type initFunc func()

// StartListening starts listening for OS signals to gracefully shutdown the application.
func StartListening(init initFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go init()
	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v", sig)) // Remove the second argument 'sig'
	os.Exit(0)
}
