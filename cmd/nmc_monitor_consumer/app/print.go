package app

import (
	"fmt"
	consumer "github.com/nwpc-oper/nmc-message-client/consumer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

var (
	server        = ""
	topic         = ""
	offset  int64 = 0
	groupId       = ""
	debug         = false
	help          = false
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

		printFlagSet.StringVar(&server, "kafka-server", "", "kafka servers, split by ','")
		printFlagSet.StringVar(&topic, "kafka-topic", "nwpcproduct", "message topic")
		printFlagSet.Int64Var(&offset, "kafka-offset", 0, "message offset")
		printFlagSet.StringVar(&groupId, "kafka-group-id", "", "group id")

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
		c := consumer.PrinterConsumer{
			Source: consumer.KafkaSource{
				Brokers: serverList,
				Topic:   topic,
				Offset:  offset,
				GroupId: groupId,
			},
			WorkerCount:  0,
			ConsumerName: "printer",
			Debug:        false,
		}
		c.ConsumeMessages()
	},
}
