package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version   = "Unknown version"
	BuildTime = "Unknown build time"
	GitCommit = "Unknown GitCommit"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s (%s)\n", Version, GitCommit)
		fmt.Printf("Build at %s\n", BuildTime)
	},
}
