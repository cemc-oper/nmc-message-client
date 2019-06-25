package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	fmt.Println("NWPC Monitor Client")
	fmt.Println("This program is under development. Please contact Wang Dapeng(3083) if it fails.")
	fmt.Println()
}

var rootCmd = &cobra.Command{
	Use:   "nmc_monitor_client",
	Short: "A client for NMC Monitor.",
	Long:  "A client for NMC Monitor.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
