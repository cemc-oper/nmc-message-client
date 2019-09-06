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

type KafkaSender struct {
	Target KafkaTarget
	Debug  bool
}

func (s *KafkaSender) SendMessage(message []byte) error {

	if s.Debug {
		fmt.Println("create writer...")
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      s.Target.Brokers,
		Topic:        s.Target.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: s.Target.WriteTimeout,
	})

	if s.Debug {
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

	if s.Debug {
		fmt.Println("close writer...")
	}
	w.Close()
	if s.Debug {
		fmt.Println("close writer...done")
	}

	return nil
}
