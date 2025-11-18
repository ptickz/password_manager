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
	c := actions.Command{
		Storage: storage.GetStorage(),
	}

	c.SetupStorage()
	c.AuthInStorage()
	c.ShowNavigation()

	timeout := 10 * time.Second
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	ch := make(chan int)
	go awaitInput(ch)

	for {
		select {
		case <-ticker.C:
			c.TimeoutMessage()
			os.Exit(0)
		case a := <-ch:
			c.ProcessActions(a)
		}
	}
}

func awaitInput(ch chan int) {
	f := os.Stdin
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err == nil {
			ch <- int(n)
		} else {
			fmt.Println("Wrong input")
		}
	}
}
