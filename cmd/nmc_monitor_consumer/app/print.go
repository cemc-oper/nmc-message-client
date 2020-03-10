package app

import (
	"context"
	"encoding/json"
	"fmt"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
	"time"
)

var (
	server       = ""
	topic        = ""
	offset int64 = 0
	debug        = false
	help         = false
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
		printFlagSet.Int64Var(&offset, "offset", 0, "message offset")

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
		r.SetOffset(0)

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			var message nmc_message_client.MonitorMessage
			err = json.Unmarshal(m.Value, &message)
			if err != nil {
				log.Warnf("can't parse message: %v", err)
			}
			source := message.Source
			if len(source) < 5 || source[:5] != "nwpc_" {
				continue
			}

			if message.MessageType != "prod_grib" {
				continue
			}

			var des nmc_message_client.ProbGribMessageDescription
			err = json.Unmarshal([]byte(message.Description), &des)
			if err != nil {
				log.Warnf("can't parse description: %v", err)
			}

			dateTime := time.Unix(message.DateTime/1000, 0)

			fmt.Printf("[%d][%s][%s][prod_grib] %s +%s \n",
				m.Offset,
				dateTime.Format("2006-01-02 15:04:05"),
				message.Source,
				des.StartTime,
				des.ForecastTime)
		}

		r.Close()
	},
}
