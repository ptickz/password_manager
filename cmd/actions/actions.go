package actions

import (
	"errors"
	"fmt"
	"os"
	"password_manager/cmd/storage"
	"strconv"
	"syscall"
	"time"

	"github.com/alexmullins/zip"
	"golang.org/x/term"
)

const listEntries int = 1
const getEntry int = 2
const addEntry int = 3
const deleteEntry int = 4
const changeMasterPassword int = 5
const showNavigation int = 9
const exit int = 0

type Command struct {
	Storage      *storage.Storage
	Focus        bool
	NavigationCh chan int
	InputCh      chan string
	TickerCh     <-chan time.Time
}

func (c *Command) createEntry() {
	fmt.Println("Creating entry...")
}

func (c *Command) getEntriesFromStorage() {
	entries, err := c.Storage.ReadEntries()
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range entries.Entries {
		fmt.Println("\nId:", k)
		fmt.Print("Service: ", v.ServiceName)
		if k+1 != len(entries.Entries) {
			printLongSeparator()
		}
	}
}

func (c *Command) getEntry() {
	fmt.Print("Input entry id: ")
	input := <-c.InputCh
	n, err := strconv.Atoi(input)
	if err == nil {
		var obj *storage.Entry
		obj, err = c.Storage.ReadEntry(n)
		if err == nil && &obj != nil {
			fmt.Println("\nService: ", obj.ServiceName)
			fmt.Println("Login: ", obj.Login)
			fmt.Println("Password: ", obj.Password)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Command) deleteEntry() {
	fmt.Println("Deleting entry...")
}

func (c *Command) changeMasterPassword() {
	fmt.Println("Changing master password...")
}

func (c *Command) ShowNavigation() {
	fmt.Printf("Available actions:\n\n" +
		strconv.Itoa(listEntries) + " - List entries\n" +
		strconv.Itoa(getEntry) + " - Get entry\n" +
		strconv.Itoa(addEntry) + " - Add entry\n" +
		strconv.Itoa(deleteEntry) + " - Delete entry\n" +
		strconv.Itoa(changeMasterPassword) + " - Change master password\n" +
		strconv.Itoa(showNavigation) + " - Show navigation\n" +
		strconv.Itoa(exit) + " - Exit\n\n")
}

func (c *Command) setupMasterPassword() {
	password := ""
	successfulPasswordSetup := false
	for successfulPasswordSetup == false {
		f := func() (string, bool) {
			fmt.Println("Setup the master password: ")
			byteInput, _ := term.ReadPassword(syscall.Stdin)
			password = string(byteInput)
			fmt.Println("Repeat password: ")
			repeatByteInput, _ := term.ReadPassword(syscall.Stdin)
			repeatPassword := string(repeatByteInput)
			if password == repeatPassword {
				fmt.Println("\n\nPassword accepted, please log in")
				printLongSeparator()
				return password, true
			} else {
				fmt.Println("Passwords not matching, try again")
				printLongSeparator()
				return password, false
			}
		}
		password, successfulPasswordSetup = f()
	}
	err := c.Storage.Init(password)
	if err != nil {
		panic(err)
	}
}

func (c *Command) SetupStorage() {
	init, err := c.Storage.CheckStorageInitiated()
	if err != nil {
		panic(err)
	}
	if !init {
		c.setupMasterPassword()
	}
}

func (c *Command) AuthInStorage() {
	for access := false; !access; {
		access = c.enterStorageWithMasterPassword()
		if !access {
			fmt.Println("Wrong password")
		}
	}
}

func (c *Command) enterStorageWithMasterPassword() bool {
	fmt.Println("Enter password. Ctrl+C to exit")
	byteInput, err := term.ReadPassword(syscall.Stdin)
	password := string(byteInput)
	_, err = c.Storage.CheckAccess(password)
	if err != nil {
		if !errors.Is(err, zip.ErrPassword) {
			fmt.Println(err)
		}
		return false
	}
	c.Storage.MasterPassword = &password
	return true
}

func (c *Command) ProcessActions(input int) {
	if c.Focus == true {
		return
	}
	switch input {
	case listEntries:
		c.getEntriesFromStorage()
		printLongSeparator()
		await()
	case addEntry:
		c.Focus = true
		c.createEntry()
		printLongSeparator()
		await()
		c.Focus = false
	case getEntry:
		c.Focus = true
		c.getEntry()
		printLongSeparator()
		await()
		c.Focus = false
	case deleteEntry:
		c.Focus = true
		c.deleteEntry()
		printLongSeparator()
		await()
		c.Focus = false
	case changeMasterPassword:
		c.Focus = true
		c.changeMasterPassword()
		printLongSeparator()
		await()
		c.Focus = false
	case showNavigation:
		c.ShowNavigation()
		printLongSeparator()
		await()
	case exit:
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		fmt.Println("Unknown action")
		printLongSeparator()
	}
}

func await() {
	fmt.Println("\nAwaiting input...")
}

func (c *Command) TimeoutMessage() {
	fmt.Println("Timeout reached, bye!")
}

func printLongSeparator() {
	fmt.Print("\n----------------------------------------------------")
}
