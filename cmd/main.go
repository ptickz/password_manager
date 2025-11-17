package cmd

import (
	"os"
	"password_manager/cmd/actions"
	"password_manager/cmd/storage"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "password_manager",
	Short: "Password storage and generator",
	Long:  ``,
	Run:   main,
}

func main(_ *cobra.Command, _ []string) {
	s := storage.GetStorage()
	c := actions.Command{
		Storage: s,
	}
	c.SetupStorage()
	c.AuthInStorage()
	c.ShowNavigation()
	c.ProcessActions()
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
