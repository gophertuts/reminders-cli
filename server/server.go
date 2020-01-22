package server

import (
	"log"
	"os"
	"os/signal"
)

// Stopper represents a generic service that can be stopped
type Stopper interface {
	Stop() error
}

// ListenForSignals listens for OS signals and responds to the provided ones
func ListenForSignals(signals []os.Signal, apps ...Stopper) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	sig := <-c
	log.Printf("received shutdown signal: %v\n", sig.String())

	var errs []error
	for _, app := range apps {
		err := app.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}
	var exitCode int
	for _, err := range errs {
		log.Printf("could not stop service due to: %v\n", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
