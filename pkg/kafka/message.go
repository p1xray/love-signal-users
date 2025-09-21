package kafka

type ReceivedMessage struct {
	Topic     string
	Partition int
	Offset    int64
	Data      []byte
}
