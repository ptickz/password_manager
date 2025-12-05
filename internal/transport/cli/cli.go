package cli

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"password_manager/internal/transport"
	"strconv"
	"syscall"
)

type Cli struct {
	NavigationCh chan int
	InputCh      chan string
	Focus        bool
}

func NewCli() *Cli {
	return &Cli{
		NavigationCh: make(chan int),
		InputCh:      make(chan string),
		Focus:        false,
	}
}

func (c *Cli) StartInputScanner() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if c.Focus {
			c.InputCh <- scanner.Text()
		} else {
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Wrong input")
			} else {
				c.NavigationCh <- n
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cli) SendMessageToUser(message string) {
	fmt.Printf(message)
}

func (c *Cli) GetPasswordHidden() (string, error) {
	passwordByte, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	passwordString := string(passwordByte)
	return passwordString, nil
}

func (c *Cli) GetChannels() *transport.Channels {
	return &transport.Channels{
		NavigationCh: c.NavigationCh,
		InputCh:      c.InputCh,
	}
}

func (c *Cli) SwitchFocus(b bool) {
	c.Focus = b
}
