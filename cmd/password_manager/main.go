package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"password_manager/internal/actions"
	"password_manager/internal/app"
	"password_manager/internal/storage"
	"strconv"
	"syscall"
	"time"
)

func main() {
	app.Build()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	timeout := 60 * time.Minute
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

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

	go awaitInput(&c, ctx)

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
