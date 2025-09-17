package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// Consumer provides access to the kafka reader.
type Consumer struct {
	reader *kafka.Reader
	output chan ReceivedMessage
	notify chan error
}

// NewConsumerGroup returns new instance of Consumer which configured as consumer group.
func NewConsumerGroup(address []string, groupID, topic string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: address,
		GroupID: groupID,
		Topic:   topic,
	})

	return &Consumer{
		reader: r,
		notify: make(chan error),
		output: make(chan ReceivedMessage),
	}
}

// Consume reads data from kafka.
func (c *Consumer) Consume(ctx context.Context) {
	defer close(c.notify)
	defer close(c.output)

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.notify <- err
		}

		receivedMessage := ReceivedMessage{
			Topic: msg.Topic,
			Data:  msg.Value,
		}
		c.output <- receivedMessage
	}
}

// Notify - notifies about kafka reader errors.
func (c *Consumer) Notify() <-chan error {
	return c.notify
}

// Output - provides access to the received data from kafka.
func (c *Consumer) Output() <-chan ReceivedMessage {
	return c.output
}

// Close - closes connection to kafka reader.
func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrKafkaReaderClose, err)
	}

	return nil
}
