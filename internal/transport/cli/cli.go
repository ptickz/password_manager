package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/term"

	"password_manager/internal/transport"
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
				c.NavigationCh <- 1000
			} else {
				c.NavigationCh <- n

				if n == 9 {
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cli) SendMessageToUser(message string) {
	fmt.Printf("%s", message)
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
