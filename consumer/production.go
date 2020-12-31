package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	nmc_message_client "github.com/nwpc-oper/nmc-message-client"
	"github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	err := s.Source.CreateConnection()
	if err != nil {
		if s.Source.Reader != nil {
			s.Source.Reader.Close()
		}
		return err
	}

	defer s.Source.Reader.Close()

	// create elasticsearch client.
	ctx := context.Background()
	// can't connect to es in docker without the last two options.
	// see https://github.com/olivere/elastic/issues/824
	client, err := elastic.NewClient(
		elastic.SetURL(s.Target.Server),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "elastic",
			"event":     "connect",
		}).Errorf("connect has error: %v", err)
		return err
	}

	// consume messages from rabbitmq
	log.WithFields(log.Fields{
		"component": "production",
		"event":     "consume",
	}).Info("start to consume...")

	messageChannel := make(chan nmc_message_client.GribProduction, 10)
	done := make(chan bool)

	go s.readFromKafka(messageChannel, done)
	go s.consumeProdGribMessageToElastic(client, ctx, messageChannel, done)

	select {}

	return nil
}

func (s *ProductionConsumer) readFromKafka(
	messageChannel chan nmc_message_client.GribProduction,
	done chan bool,
) {
	for {
		m, err := s.Source.Reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		var message nmc_message_client.MonitorMessageV2
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "production",
				"event":     "consume",
			}).Warnf("can't parse message: %v", err)
		}

		if !isProductionGribMessage(message) {
			continue
		}

		gribProduction, err := generateGribProduction(message, m)
		if err != nil {
			log.WithFields(log.Fields{
				"component": "production",
				"event":     "consume",
			}).Warnf("can't generate GribProduction: %v", err)
			continue
		}

		messageChannel <- gribProduction
	}
	done <- true
}

func (s *ProductionConsumer) consumeProdGribMessageToElastic(
	client *elastic.Client,
	ctx context.Context,
	messageChannel chan nmc_message_client.GribProduction,
	done chan bool,
) {
	log.WithFields(log.Fields{
		"component": "production",
		"event":     "consume",
	}).Info("start to consume...")

	var received []messageWithIndex
	isDone := false
	for {
		select {
		case message := <-messageChannel:
			// parse message to generate message index
			index := getIndexForProductionMessage(message)
			received = append(received, messageWithIndex{
				Index:   index,
				Id:      message.Offset,
				Message: message,
			})

			//log.WithFields(log.Fields{
			//	"component": "elastic",
			//	"event":     "message",
			//}).Infof("[%s][%s][%s][prod_grib] %s +%s",
			//	message.Offset,
			//	message.DateTime.Format("2006-01-02 15:04:05"),
			//	message.Source,
			//	message.StartTime.Format("2006010215"),
			//	message.ForecastTime,
			//)

			if len(received) > s.BulkSize {
				// send to elasticsearch
				log.WithFields(log.Fields{
					"component": "elastic",
					"event":     "push",
				}).Info("bulk size push...")
				err := pushMessages(client, received, ctx)
				if err != nil {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Warn("bulk size push...failed")
				} else {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Infof("bulk size push...done, %d", len(received))
					received = nil
				}
			}
		case <-time.After(time.Second * 1):
			if len(received) > 0 {
				// send to elasticsearch
				log.WithFields(log.Fields{
					"component": "elastic",
					"event":     "push",
				}).Info("time limit push...")
				err := pushMessages(client, received, ctx)
				if err != nil {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Warn("time limit push...failed")
				} else {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Infof("time limit push...done, %d", len(received))
					received = nil
				}
			}
		case <-done:
			if len(received) > 0 {
				// send to elasticsearch
				log.WithFields(log.Fields{
					"component": "elastic",
					"event":     "push",
				}).Info("done push...")
				err := pushMessages(client, received, ctx)
				if err != nil {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Warn("done push...failed")
				} else {
					log.WithFields(log.Fields{
						"component": "elastic",
						"event":     "push",
					}).Infof("done push...done, %d", len(received))
					received = nil
				}
			}
			isDone = true
		}
		if isDone {
			break
		}
		if len(received) >= s.BulkSize*10 {
			log.WithFields(log.Fields{
				"component": "elastic",
				"event":     "push",
			}).Fatalf("Count of received messages is larger than 10 times of bulk size: %s", len(received))
		}
	}
}

func isProductionGribMessage(message nmc_message_client.MonitorMessageV2) bool {
	return true
}

func getIndexForProductionMessage(message nmc_message_client.GribProduction) string {
	messageTime := message.DateTime
	indexName := fmt.Sprintf("nmc-product-%s", messageTime.Format("2006-01"))
	return indexName
}

func generateGribProduction(
	message nmc_message_client.MonitorMessageV2,
	m kafka.Message,
) (nmc_message_client.GribProduction, error) {
	var des nmc_message_client.ProbGribMessageDescription
	err := json.Unmarshal([]byte(message.ResultDescription), &des)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "production",
			"event":     "generateGribProduction",
		}).Warnf("can't parse description: %v", err)
		return nmc_message_client.GribProduction{}, err
	}

	startTime, err := time.Parse("2006010215", des.StartTime)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "production",
			"event":     "generateGribProduction",
		}).Warnf("can't parse start time: %v", err)
		return nmc_message_client.GribProduction{}, err
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", message.DateTime)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "production",
			"event":     "generateGribProduction",
		}).Warnf("can't parse DateTime: %v", err)
		return nmc_message_client.GribProduction{}, err
	}

	p := nmc_message_client.GribProduction{
		Offset:       strconv.Itoa(int(m.Offset)),
		Source:       message.Source,
		MessageType:  message.MessageType,
		Status:       message.Result,
		DateTime:     dateTime,
		FileName:     message.FileNames,
		StartTime:    startTime,
		ForecastTime: des.ForecastTime,
	}

	return p, nil
}
