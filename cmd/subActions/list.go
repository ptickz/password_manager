package subActions

import (
	"bufio"
	"fmt"
	"os"
	"password_manager/cmd"

	"github.com/spf13/cobra"
)

func init() {
	cmd.rootCmd.AddCommand(savedServices)
}

var savedServices = &cobra.Command{
	Use:   "list",
	Short: "Print all saved services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter master password: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == "123\n" {
			fmt.Println("Entered")
		}
		fmt.Println("Wrong password")
	},
}
