package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

var (
	target           = ""
	topic            = ""
	source           = ""
	messageType      = ""
	status           = "0"
	datetime         int64
	fileName         = ""
	absoluteDataName = ""
	startTime        = ""
	forecastTime     = ""
	debug            = false
	disableSend      = false
)

func init() {
	rootCmd.AddCommand(sendCmd)

	currentTimeStamp := makeTimestamp()

	sendCmd.Flags().StringVar(&target, "target", "", "send target")
	sendCmd.Flags().StringVar(&topic, "topic", "monitor", "message topic")
	sendCmd.Flags().StringVar(&source, "source", "", "message source")
	sendCmd.Flags().StringVar(&messageType, "type", "", "message type")
	sendCmd.Flags().StringVar(&status, "status", "0", "status")
	sendCmd.Flags().Int64Var(&datetime, "datetime", currentTimeStamp, "datetime, default is current time.")
	sendCmd.Flags().StringVar(&fileName, "file-name", "", "file name")
	sendCmd.Flags().StringVar(&absoluteDataName, "absolute-data-name", "", "absolute data name")
	sendCmd.Flags().StringVar(&startTime, "start-time", "", "start time, such as 2019062400")
	sendCmd.Flags().StringVar(&forecastTime, "forecast-time", "", "forecast time, such as 000")
	sendCmd.Flags().BoolVar(&debug, "debug", false, "show debug information")
	sendCmd.Flags().BoolVar(&disableSend, "disable-send", false, "disable message send.")

	sendCmd.MarkFlagRequired("target")
	sendCmd.MarkFlagRequired("source")
	sendCmd.MarkFlagRequired("type")
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send message to NMC Monitor",
	Long:  "Send message to NMC Monitor",
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			fmt.Printf("Version %s (%s)\n", Version, GitCommit)
			fmt.Printf("Build at %s\n", BuildTime)
		}

		description := MessageDescription{
			StartTime:    startTime,
			ForecastTime: forecastTime,
		}
		descriptionBlob, err := json.Marshal(description)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create description error: %s\n", err)
			os.Exit(2)
		}

		monitorMessage := MonitorMessage{
			Source:           source,
			MessageType:      messageType,
			Status:           status,
			DateTime:         datetime,
			FileName:         fileName,
			AbsoluteDataName: absoluteDataName,
			Description:      string(descriptionBlob),
		}

		monitorMessageBlob, err := json.MarshalIndent(monitorMessage, "", "  ")

		if err != nil {
			fmt.Fprintf(os.Stderr, "create message error: %s\n", err)
			os.Exit(2)
		}

		if debug {
			fmt.Printf("message:\n")
			fmt.Printf("%s\n", monitorMessageBlob)
		}

		if disableSend {
			if debug {
				fmt.Printf("disable send.\n")
				fmt.Printf("Bye.\n")
			}
			return
		}

		if debug {
			fmt.Println("create writer...")
		}

		w := kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{target},
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
		})

		if debug {
			fmt.Println("create writer...done")
		}

		if debug {
			fmt.Println("send message...")
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Value: monitorMessageBlob,
			},
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "send message failed: %s\n", err)
			os.Exit(4)
		}

		fmt.Printf("send message successful\n")

		if debug {
			fmt.Println("close writer...")
		}
		w.Close()
		if debug {
			fmt.Println("close writer...done")
		}
	},
}

type MonitorMessage struct {
	Source           string `json:"source"`
	MessageType      string `json:"type"`
	Status           string `json:"status"`
	DateTime         int64  `json:"datetime,omitempty"`
	FileName         string `json:"fileName"`
	AbsoluteDataName string `json:"absoluteDataName,omitempty"`
	Description      string `json:"desc,omitempty"`
}

type MessageDescription struct {
	StartTime    string `json:"startTime,omitempty"`
	ForecastTime string `json:"forecastTime,omitempty"`
}
