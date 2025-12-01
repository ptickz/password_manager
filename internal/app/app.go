package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"password_manager/internal/actions"
	"password_manager/internal/storage"
	"strconv"
	"syscall"
	"time"
)

type App struct {
	config  config.Config
	command actions.Command
	ctx     context.Context
	timeout <-chan time.Time
}

func (a *App) NewApp() *App {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	timeout := 60 * time.Minute
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	return &App{
		storage: storage.NewStorage(),
		command: a.command,
		ctx:     ctx,
		timeout: ticker.C,
	}

	c := actions.Command{
		Storage:      storage.NewStorage(),
		Focus:        false,
		ScannerState: true,
		NavigationCh: make(chan int),
		InputCh:      make(chan string),
		TickerCh:     ticker.C,
	}

	c.SetupStorage()
	c.AuthInStorage()
	c.ShowNavigation()

	go awaitInput(&c)

	for {
		select {
		case <-c.TickerCh:
			c.TimeoutMessage()
		case a := <-c.NavigationCh:
			c.ProcessActions(a)
			ticker.Reset(timeout)
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) Run() {
	a.command.SetupStorage()
	a.command.AuthInStorage()
	a.command.ShowNavigation()

	go awaitInput(&a.command)

	for {
		select {
		case <-a.command.TickerCh:
			a.command.TimeoutMessage()
		case n := <-a.command.NavigationCh:
			a.command.ProcessActions(n)
			ticker.Reset(timeout)
		case <-ctx.Done():
			return
		}
	}
}

func awaitInput(c *actions.Command) {
	f := os.Stdin
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if c.Focus {
			c.InputCh <- scanner.Text()
		} else {
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Wrong input")
			}

			c.NavigationCh <- n
		}
	}
	if err := scanner.Err(); err != nil {

	}
}
