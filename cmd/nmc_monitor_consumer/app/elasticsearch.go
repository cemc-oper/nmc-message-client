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
	elasticServer = ""
	bulkSize      = 20
)

func init() {
	rootCmd.AddCommand(elasticSearchCommand)
}

var elasticSearchCommand = &cobra.Command{
	Use:                "elasticsearch",
	Short:              "Send messages to ElasticSearch",
	Long:               "Send messages to ElasticSearch",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		var flagSet = pflag.NewFlagSet("elasticsearch", pflag.ContinueOnError)
		flagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
		flagSet.SortFlags = false

		flagSet.StringVar(&server, "kafka-server", "", "kafka servers, split by ','")
		flagSet.StringVar(&topic, "kafka-topic", "nwpcproduct", "message topic")
		flagSet.Int64Var(&offset, "kafka-offset", 0, "message offset")
		flagSet.StringVar(&groupId, "kafka-group-id", "", "group id")

		flagSet.StringVar(&elasticServer,
			"elasticsearch-server", "", "elasticsearch server")
		flagSet.IntVar(&bulkSize, "bulk-size", 50, "bulk size")

		flagSet.BoolVar(&debug, "debug", false, "show debug information")

		flagSet.BoolVar(&help, "help", false,
			"show help information.")

		if err := flagSet.Parse(args); err != nil {
			cmd.Usage()
			log.Fatal(err)
		}

		// check if there are non-flag arguments in the command line
		cmds := flagSet.Args()
		if len(cmds) > 0 {
			cmd.Usage()
			log.Fatalf("unknown command: %s", cmds[0])
		}

		// short-circuit on help
		help, err := flagSet.GetBool("help")
		if err != nil {
			log.Fatal(`"help" flag is non-bool, programmer error, please correct`)
		}

		if help {
			cmd.Help()
			fmt.Printf("%s\n", flagSet.FlagUsages())
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

		if groupId == "" {
			log.Fatal("groupId option is required")
		}

		// main
		c := consumer.ProductionConsumer{
			Source: consumer.KafkaSource{
				Brokers: serverList,
				Topic:   topic,
				Offset:  offset,
				GroupId: groupId,
			},
			Target: consumer.ElasticSearchTarget{
				Server: elasticServer,
			},
			WorkerCount: 1,
			BulkSize:    bulkSize,
			Debug:       false,
		}
		c.ConsumeMessages()
	},
}
