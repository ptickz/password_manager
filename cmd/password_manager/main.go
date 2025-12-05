package main

import (
	"context"
	"os/signal"
	"password_manager/internal/app"
	"password_manager/internal/config"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	a, err := app.NewApp(c)
	if err != nil {
		panic(err)
	}
	go func() {
		err := a.Run()
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	err = a.GracefulShutdown()
	if err != nil {
		panic(err)
	}
}
