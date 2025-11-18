package actions

import (
	"errors"
	"fmt"
	"os"
	"password_manager/cmd/storage"
	"strconv"
	"syscall"

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
	Storage storage.Storage
}

func (c Command) createEntry() {
	fmt.Println("Creating entry...")
}

func (c Command) getEntriesFromStorage() {
	fmt.Println("Getting entries from storage...")
}

func (c Command) getEntry() {
	fmt.Println("Getting entry...")
}

func (c Command) deleteEntry() {
	fmt.Println("Deleting entry...")
}

func (c Command) changeMasterPassword() {
	fmt.Println("Changing master password...")
}

func (c Command) ShowNavigation() {
	fmt.Printf("Available actions:\n\n" +
		strconv.Itoa(listEntries) + " - List entries for profile\n" +
		strconv.Itoa(getEntry) + " - Get entry for profile\n" +
		strconv.Itoa(addEntry) + "- Add entry for profile\n" +
		strconv.Itoa(deleteEntry) + " - Delete entry for profile\n" +
		strconv.Itoa(changeMasterPassword) + " - Change master password\n" +
		strconv.Itoa(showNavigation) + " - Show navigation\n" +
		strconv.Itoa(exit) + " - Exit\n\n")
}

func (c Command) setupMasterPassword() {
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
				fmt.Println("Password accepted")
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

func (c Command) SetupStorage() {
	init, err := c.Storage.CheckStorageInitiated()
	if err != nil {
		panic(err)
	}
	if !init {
		c.setupMasterPassword()
	}
}

func (c Command) AuthInStorage() {
	b := c.enterStorageWithMasterPassword()
	if !b {
		fmt.Println("Wrong password")
		return
	}
}

func (c Command) enterStorageWithMasterPassword() bool {
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
	return true
}

func (c Command) ProcessActions(input int) {
	switch input {
	case listEntries:
		c.getEntriesFromStorage()
		printLongSeparator()
	case addEntry:
		c.createEntry()
		printLongSeparator()
	case getEntry:
		c.getEntry()
		printLongSeparator()
	case deleteEntry:
		c.deleteEntry()
		printLongSeparator()
	case changeMasterPassword:
		c.changeMasterPassword()
		printLongSeparator()
	case showNavigation:
		c.ShowNavigation()
		printLongSeparator()
	case exit:
		exitApp()
	default:
		fmt.Println("Unknown action")
		printLongSeparator()
	}
}

func (c Command) TimeoutMessage() {
	fmt.Println("Timeout reached, bye!")
}

func exitApp() {
	fmt.Println("Bye!")
	os.Exit(0)
}

func printLongSeparator() {
	fmt.Println("\n----------------------------------------------------")
}
