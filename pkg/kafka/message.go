package kafka

type ReceivedMessage struct {
	Topic string
	Data  []byte
}
