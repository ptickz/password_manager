package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"password_manager/internal/app"
	"password_manager/internal/config"
)

/*
* TODO check for wrong password, immediate error
* TODO change master password
* TODO tests
 */
func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	c, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Unable to build config, error: %s\n", err)
		stop()
	}

	a, err := app.NewApp(c)
	if err != nil {
		fmt.Printf("Unable to build app, error: %s\n", err)
		stop()
	}

	go listenSignals(stop)
	go a.Run(stop)

	<-ctx.Done()

	err = a.GracefulShutdown()
	if err != nil {
		fmt.Printf("Unable to gracefully shutdown application, error: %s\n", err)
		os.Exit(1)
	}
}

func listenSignals(stop context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	sig := <-ch
	log.Printf("received signal: %s", sig)
	stop()
}
