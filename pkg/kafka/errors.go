package kafka

import "errors"

var (
	ErrKafkaReaderClose = errors.New("error closing kafka reader")
)
