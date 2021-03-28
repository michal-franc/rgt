package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Version - variable used for ldflags to inject Version number to binary
var Version = "local-development"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the go pomodoro app version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
