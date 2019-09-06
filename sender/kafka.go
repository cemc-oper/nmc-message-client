package sender

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaTarget struct {
	Brokers      []string
	Topic        string
	WriteTimeout time.Duration
}

func SendMessageToKafka(kafkaTarget KafkaTarget, message []byte, debug bool) error {
	if debug {
		fmt.Println("create writer...")
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      kafkaTarget.Brokers,
		Topic:        kafkaTarget.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: kafkaTarget.WriteTimeout,
	})

	if debug {
		fmt.Println("create writer...done")
		fmt.Println("send message...")
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		return fmt.Errorf("send message failed: %s", err)
	}

	fmt.Printf("send message successful\n")

	if debug {
		fmt.Println("close writer...")
	}
	w.Close()
	if debug {
		fmt.Println("close writer...done")
	}

	return nil
}
