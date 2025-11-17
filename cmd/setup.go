package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"password_manager/cmd/storage"
	"password_manager/cmd/subActions"
	"strconv"
	"strings"
	"syscall"

	"github.com/alexmullins/zip"
	"golang.org/x/term"

	"github.com/spf13/cobra"
)

const ListEntries = 1
const GetEntry = 2
const AddEntry = 3
const DeleteEntry = 4
const ChangeMasterPassword = 5
const ShowNavigation = 9
const Exit = 0

type Command struct {
	Storage storage.Storage
}

var RootCmd = &cobra.Command{
	Use:   "password_manager",
	Short: "Password storage and generator",
	Long:  ``,
	Run:   main,
}

func main(cmd *cobra.Command, args []string) {
	s := storage.GetStorage()
	c := Command{
		Storage: s,
	}
	init, err := s.CheckStorageInitiated()
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
	c.showNavigation()
	for exit := false; !exit; {
		fmt.Println("Select action: ")
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
	case ListEntries:
		err = subActions.GetListFromStorage()
		printLongSeparator()
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		return false, nil
	case AddEntry:
		err = subActions.CreateEntry()
		fmt.Println("Add entry")
		return false, nil
	case GetEntry:
		//err = subActions.GetEntryForProfile()
		fmt.Println("Get entry")
		return false, nil
	case DeleteEntry:
		//err = subActions.DeleteEntryForProfile()
		fmt.Println("Delete entry")
		return false, nil
	case ChangeMasterPassword:
		//err = subActions.ChangeMasterPassword()
		fmt.Println("Change master password")
		return false, nil
	case ShowNavigation:
		c.showNavigation()
		return false, nil
	case Exit:
		fmt.Println("Exit")
		return true, nil
	default:
		fmt.Println("Unknown action")
		return false, nil
	}
}

func (c Command) enterStorageWithMasterPassword() bool {
	fmt.Println("Enter password. Ctrl+C to exit")
	byteInput, err := term.ReadPassword(syscall.Stdin)
	password := string(byteInput)
	_, err = storage.CheckPasswordValid(password)
	if err != nil {
		if !errors.Is(err, zip.ErrPassword) {
			fmt.Println(err)
		}
		return false
	}
	return true
}

func (c Command) showNavigation() {
	fmt.Printf("Available actions:\n\n" +
		strconv.Itoa(ListEntries) + " - List entries for profile\n" +
		strconv.Itoa(AddEntry) + "- Add entry for profile\n" +
		strconv.Itoa(GetEntry) + " - Get entry for profile\n" +
		strconv.Itoa(DeleteEntry) + " - Delete entry for profile\n" +
		strconv.Itoa(ChangeMasterPassword) + " - Change master password\n" +
		strconv.Itoa(ShowNavigation) + " - Show navigation\n" +
		strconv.Itoa(Exit) + " - Exit\n\n")
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
	err := storage.WriteSetupStorageData(password)
	if err != nil {
		fmt.Println(err)
	}
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func printLongSeparator() {
	fmt.Println("\n----------------------------------------------------")
}
