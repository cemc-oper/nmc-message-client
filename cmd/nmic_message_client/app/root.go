package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fmt.Println("NMIC Message Client")
	fmt.Println("This program is under development. Please contact Wang Dapeng(3083) if it fails.")
	fmt.Println()
}

var rootCmd = &cobra.Command{
	Use:   "nmic_message_clinet",
	Short: "A client for NMIC message.",
	Long:  "A client for NMIC message.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
