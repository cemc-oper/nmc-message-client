package sender

type Sender interface {
	SendMessage([]byte) error
}
