package actions

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"password_manager/cmd/storage"
	"strconv"
	"strings"
	"syscall"

	"github.com/alexmullins/zip"
	"golang.org/x/term"
)

const listEntries = 1
const getEntry = 2
const addEntry = 3
const deleteEntry = 4
const changeMasterPassword = 5
const showNavigation = 9
const exit = 0

type Command struct {
	Storage storage.Storage
}

func (c Command) createEntry() error {
	//TODO implepent
	return err
}

func (c Command) getEntriesFromStorage() error {
	//TODO implement
	return err
}

func (c Command) getEntry() error {
	//TODO implement
	return err
}

func (c Command) deleteEntry() error {
	//TODO implement
	return err
}

func (c Command) changeMasterPassword() error {
	//TODO implement
	return err
}

func (c Command) showNavigation() {
	fmt.Printf("Available actions:\n\n" +
		strconv.Itoa(listEntries) + " - List entries for profile\n" +
		strconv.Itoa(addEntry) + "- Add entry for profile\n" +
		strconv.Itoa(getEntry) + " - Get entry for profile\n" +
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
		fmt.Println(err)
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
	b := c.enterStorageWithMasterPassword()
	if !b {
		fmt.Println("Wrong password")
		return
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

func (c Command) ProcessActions() {
	for exit := false; !exit; {
		fmt.Println("Select action: ")
		var err error
		exit, err = c.processAction()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c Command) processAction() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	u64Input, err := strconv.ParseUint(
		strings.ReplaceAll(input,
			"\n",
			"",
		), 10, 64)
	if err != nil {
		fmt.Println("Invalid input")
		fmt.Println(err)
		return false, nil
	}

	switch uint(u64Input) {
	case listEntries:
		err = c.getEntriesFromStorage()
		printLongSeparator()
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		return false, nil
	case addEntry:
		err = c.createEntry()
		fmt.Println("Add entry")
		return false, nil
	case getEntry:
		err = c.getEntry()
		fmt.Println("Get entry")
		return false, nil
	case deleteEntry:
		err = c.deleteEntry()
		fmt.Println("Delete entry")
		return false, nil
	case changeMasterPassword:
		err = c.changeMasterPassword()
		fmt.Println("Change master password")
		return false, nil
	case showNavigation:
		c.showNavigation()
		return false, nil
	case exit:
		fmt.Println("Exit")
		return true, nil
	default:
		fmt.Println("Unknown action")
		return false, nil
	}
}

func printLongSeparator() {
	fmt.Println("\n----------------------------------------------------")
}
