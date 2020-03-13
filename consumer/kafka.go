package consumer

import (
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaSource struct {
	Brokers []string
	Topic   string
	Offset  int64

	Reader *kafka.Reader
}

func (source *KafkaSource) CreateConnection() error {
	log.WithFields(log.Fields{
		"component": "kafka",
		"event":     "connect",
	}).Infof("create kafka reader...%s", source.Brokers)
	source.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   source.Brokers,
		Topic:     source.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	return nil
}
