package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type PrinterConsumer struct {
	Source       KafkaSource
	WorkerCount  int
	ConsumerName string
	Debug        bool
}

func (s *PrinterConsumer) ConsumeMessages() error {
	// create connection to rabbitmq
	err := s.Source.CreateConnection()
	if err != nil {
		if s.Source.Reader != nil {
			s.Source.Reader.Close()
		}
		return err
	}

	defer s.Source.Reader.Close()

	// consume messages from rabbitmq
	log.WithFields(log.Fields{
		"component": "kafka",
		"event":     "consume",
	}).Info("start to consume...")
	for {
		m, err := s.Source.Reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		var message nmc_message_client.MonitorMessageV2
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "kafka",
				"event":     "consume",
			}).Warnf("can't parse message: %v", err)
		}
		if !isProductionGribMessageV2(message) {
			continue
		}

		printProdGribMessage(message, m)
	}

	return nil
}

func printProdGribMessage(message nmc_message_client.MonitorMessageV2, m kafka.Message) {
	var des nmc_message_client.ProbGribMessageDescription
	err := json.Unmarshal([]byte(message.ResultDescription), &des)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "kafka",
			"event":     "consume",
		}).Warnf("can't parse description: %v", err)
	}

	fmt.Printf("[%d][%s][%s][prod_grib] %s +%s \n",
		m.Offset,
		message.DateTime,
		message.Source,
		des.StartTime,
		des.ForecastTime,
	)
}

func isProductionGribMessageV2(message nmc_message_client.MonitorMessageV2) bool {
	return true
}
