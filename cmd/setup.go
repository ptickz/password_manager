package cmd

import (
	"bufio"
	"fmt"
	"os"
	"password_manager/cmd/utils"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/spf13/cobra"
)

const ListEntries = 1
const AddEntry = 2
const GetEntry = 3
const DeleteEntry = 4
const ChangeMasterPassword = 5
const Exit = 0

var rootCmd = &cobra.Command{
	Use:   "password_manager",
	Short: "Password storage and generator",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.CheckStorageEntryExists() {
			setupMasterPassword()
		}
		err := enterStorageWithMasterPassword()
		if err != nil {
			fmt.Println(err)
		}
		showNavigation()
		for exit := false; !exit; {
			fmt.Println("Select action: ")
			exit, err = awaitAction()
		}
	},
}

func awaitAction() (bool, error) {
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
		fmt.Println("List entries")
		return false, nil
	case AddEntry:
		fmt.Println("Add entry")
		return false, nil
	case GetEntry:
		fmt.Println("Get entry")
		return false, nil
	case DeleteEntry:
		fmt.Println("Delete entry")
		return false, nil
	case ChangeMasterPassword:
		fmt.Println("Change master password")
		return false, nil
	case Exit:
		fmt.Println("Exit")
		return true, nil
	default:
		fmt.Println("Unknown action")
		return false, nil
	}
}

func enterStorageWithMasterPassword() error {
	fmt.Println("Enter password. Ctrl+C to exit")

	return nil
}

func showNavigation() {
	fmt.Printf("Available actions:\n\n" + "1. List entries\n" + "2. Add entry\n" + "3. Get entry\n" +
		"4. Delete entry\n" + "5. Change master password\n" + "0. Exit\n\n\n")
}

func setupMasterPassword() {
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
				return password, true
			} else {
				fmt.Println("Passwords not matching, try again\n----------------------------------------------------")
				return password, false
			}
		}
		password, successfulPasswordSetup = f()
	}
	err := utils.WriteStorageData(map[string]string{
		"master": password,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
