package main

import (
	"bufio"
	"fmt"
	"os"
	"password_manager/cmd/actions"
	"password_manager/cmd/storage"
	"strconv"
	"time"
)

func main() {
	navigationInputChannel := make(chan int)
	commandInputChannel := make(chan string)
	timeout := 60 * time.Minute
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	c := actions.Command{
		Storage:      storage.GetStorage(),
		Focus:        false,
		NavigationCh: navigationInputChannel,
		InputCh:      commandInputChannel,
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
			os.Exit(0)
		case a := <-c.NavigationCh:
			c.ProcessActions(a)
		}
	}
}

func awaitInput(c *actions.Command) {
	f := os.Stdin
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if c.Focus == true {
			n := scanner.Text()
			c.InputCh <- n
		} else {
			n, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err == nil {
				c.NavigationCh <- int(n)
			} else {
				fmt.Println("Wrong input")
			}
		}
	}
}
