package app

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"strings"
)

var (
	server = ""
	topic  = ""
	debug  = false
	help   = false
)

func init() {
	rootCmd.AddCommand(printCommand)
}

var printCommand = &cobra.Command{
	Use:                "print",
	Short:              "Print messages from NMC Monitor",
	Long:               "Print messages from NMC Monitor",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		var printFlagSet = pflag.NewFlagSet("print", pflag.ContinueOnError)
		printFlagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
		printFlagSet.SortFlags = false

		printFlagSet.StringVar(&server, "server", "", "kafka servers, split by ','")
		printFlagSet.StringVar(&topic, "topic", "monitor", "message topic")

		printFlagSet.BoolVar(&debug, "debug", false, "show debug information")

		printFlagSet.BoolVar(&help, "help", false,
			"show help information.")

		if err := printFlagSet.Parse(args); err != nil {
			cmd.Usage()
			log.Fatal(err)
		}

		// check if there are non-flag arguments in the command line
		cmds := printFlagSet.Args()
		if len(cmds) > 0 {
			cmd.Usage()
			log.Fatalf("unknown command: %s", cmds[0])
		}

		// short-circuit on help
		help, err := printFlagSet.GetBool("help")
		if err != nil {
			log.Fatal(`"help" flag is non-bool, programmer error, please correct`)
		}

		if help {
			cmd.Help()
			fmt.Printf("%s\n", printFlagSet.FlagUsages())
			return
		}

		// check required flags
		if server == "" {
			log.Fatal("server option is required")
		}
		serverList := strings.Split(server, ",")

		if topic == "" {
			log.Fatal("source option is required")
		}

		// main
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   serverList,
			Topic:     topic,
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  10e6, // 10MB
		})

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}

		r.Close()
	},
}
