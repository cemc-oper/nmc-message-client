package consumer

type Consumer interface {
	ConsumeMessages() error
}
