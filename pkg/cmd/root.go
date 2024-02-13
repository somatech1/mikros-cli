package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mikros-cli",
		Short: "A \"swiss army knife\" for dealing with mikros framework tasks.",
		Long: `mikros-cli is a command to help the developer use the mikros
framework to create new services.`,
	}
)

// Execute puts the CLI to execute.
func Execute() {
	loadCommands()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadCommands is where all CLI options are loaded and prepared to be
// executed after.
func loadCommands() {
	serviceCmdInit()
}
