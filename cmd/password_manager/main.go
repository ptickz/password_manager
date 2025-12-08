package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"password_manager/internal/app"
	"password_manager/internal/config"
	"syscall"
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
		fmt.Println(err)
		os.Exit(1)
	}
	a, err := app.NewApp(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
