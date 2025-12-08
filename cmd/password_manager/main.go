package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"password_manager/internal/app"
	"password_manager/internal/config"
	"syscall"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	a, err := app.NewApp(c)
	if err != nil {
		panic(err)
	}
	go listenSignals(stop)
	go a.Run(stop)

	<-ctx.Done()
	err = a.GracefulShutdown()
	if err != nil {
		panic(err)
	}
}

func listenSignals(stop context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	sig := <-ch
	log.Printf("received signal: %s", sig)
	stop()
}
