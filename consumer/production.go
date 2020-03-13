package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type ProductionConsumer struct {
	Source      KafkaSource
	Target      ElasticSearchTarget
	WorkerCount int
	BulkSize    int
	Debug       bool
}

func (s *ProductionConsumer) ConsumeMessages() error {
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
		var message nmc_message_client.MonitorMessage
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "kafka",
				"event":     "consume",
			}).Warnf("can't parse message: %v", err)
		}
		if !isProductionGribMessage(message) {
			continue
		}

		saveProdGribMessage(message, m)
	}

	return nil
}

func isProductionGribMessage(message nmc_message_client.MonitorMessage) bool {
	source := message.Source
	if len(source) < 5 || source[:5] != "nwpc_" {
		return false
	}

	if message.MessageType != "prod_grib" {
		return false
	}
	return true
}

func saveProdGribMessage(message nmc_message_client.MonitorMessage, m kafka.Message) {
	var des nmc_message_client.ProbGribMessageDescription
	err := json.Unmarshal([]byte(message.Description), &des)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "kafka",
			"event":     "consume",
		}).Warnf("can't parse description: %v", err)
	}

	dateTime := time.Unix(message.DateTime/1000, 0)

	fmt.Printf("[%d][%s][%s][prod_grib] %s +%s \n",
		m.Offset,
		dateTime.Format("2006-01-02 15:04:05"),
		message.Source,
		des.StartTime,
		des.ForecastTime,
	)

}

func getIndexForProductionMessage(message nmc_message_client.GribProduction) string {
	messageTime := message.DateTime
	indexName := messageTime.Format("2006-01")
	return indexName
}
