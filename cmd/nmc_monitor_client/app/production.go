package app

import (
	"fmt"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/nwpc-oper/nmc-message-client/sender"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
	"time"
)

var (
	sourceIP        string
	fileSize        int64
	statusNo        int8
	datetimeString  string
	productInterval int32
)

func init() {
	rootCmd.AddCommand(productionCmd)
}

var productionCmd = &cobra.Command{
	Use:                "production",
	Short:              "Send message to NMC Monitor",
	Long:               "Send message to NMC Monitor",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		var cstZone = time.FixedZone("CST", 8*3600) // beijing
		currentTime := time.Now().In(cstZone)

		var flagSet = pflag.NewFlagSet("send", pflag.ContinueOnError)
		flagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
		flagSet.SortFlags = false

		flagSet.StringVar(&target, "target", "", "send targets, split by ','")
		flagSet.StringVar(&topic, "topic", "nwpcproduct", "message topic")

		flagSet.StringVar(&source, "source", "", "message source")
		flagSet.StringVar(&sourceIP, "source-ip", "", "source IP")
		flagSet.StringVar(&messageType, "type", "", "message type")
		flagSet.Int8Var(&statusNo, "status", 0, "status")
		//flagSet.StringVar(&datetimeString, "datetime", currentTime.Format("2006-01-02 15:04:05"), "datetime, default is current time.")
		flagSet.StringVar(&fileName, "file-name", "", "file name")
		flagSet.StringVar(&absoluteDataName, "absolute-data-name", "", "absolute data name")
		flagSet.Int64Var(&fileSize, "file-size", -1, "file size in bytes")
		flagSet.Int32Var(&productInterval, "product-interval", 1, "product interval, unit hour")

		flagSet.BoolVar(&debug, "debug", false, "show debug information")
		flagSet.BoolVar(&disableSend, "disable-send", false, "disable message send.")
		flagSet.BoolVar(&ignoreError, "ignore-error", false,
			"ignore error. Should be open in operation systems.")

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
		if target == "" {
			log.Fatal("target option is required")
		}
		targetBrokers := strings.Split(target, ",")
		fmt.Printf("brokers: %s\n", targetBrokers)

		if source == "" {
			log.Fatal("source option is required")
		}
		if messageType == "" {
			log.Fatal("messageType option is required")
		}

		if debug {
			fmt.Printf("Version %s (%s)\n", Version, GitCommit)
			fmt.Printf("Build at %s\n", BuildTime)
		}

		// create message blob
		var messageBlob []byte

		if len(messageType) >= 4 && messageType[len(messageType)-4:len(messageType)] == "Prod" {
			messageBlob, err = generateProductMessage(
				args,
				currentTime,
			)
		} else {
			log.Fatalf("message type is not supported: %s", messageType)
		}

		if err != nil {
			f := os.Stderr
			returnCode := 2
			if ignoreError {
				f = os.Stdout
				returnCode = 0
			}
			fmt.Fprintf(f, "create message error: %s\n", err)
			os.Exit(returnCode)
		}

		if debug {
			fmt.Printf("message:\n")
			fmt.Printf("%s\n", messageBlob)
		}

		if disableSend {
			if debug {
				fmt.Printf("disable send.\n")
				fmt.Printf("Bye.\n")
			}
			return
		}

		var s sender.Sender
		s = &sender.KafkaSender{
			Target: sender.KafkaTarget{
				Brokers:      targetBrokers,
				Topic:        topic,
				WriteTimeout: 10 * time.Second,
			},
			Debug: debug,
		}

		err = s.SendMessage(messageBlob)

		if err != nil {
			f := os.Stderr
			returnCode := 4
			if ignoreError {
				f = os.Stdout
				returnCode = 0
			}
			fmt.Fprintf(f, "send message failed: %s\n", err)
			os.Exit(returnCode)
		}
	},
}

func generateProductMessage(args []string, currentTime time.Time) ([]byte, error) {
	var startTime = ""
	var forecastTime = ""

	var gribFlagSet = pflag.NewFlagSet("prod_grid", pflag.ContinueOnError)
	gribFlagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	gribFlagSet.SortFlags = false

	gribFlagSet.StringVar(&startTime, "start-time", "", "start time, such as 2019062400")
	gribFlagSet.StringVar(&forecastTime, "forecast-time", "", "forecast time, such as 000")
	if err := gribFlagSet.Parse(args); err != nil {
		log.Fatalf("argument parse fail: %s", err)
	}

	messageBlob, err := nmc_message_client.CreateProductMessage(
		topic,
		source,
		sourceIP,
		messageType,
		currentTime,
		fileName,
		absoluteDataName,
		fileSize,
		productInterval,
		statusNo,
		startTime,
		forecastTime)

	return messageBlob, err
}
