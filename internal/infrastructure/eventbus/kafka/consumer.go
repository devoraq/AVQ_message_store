// Package kafka содержит вспомогательные функции для адаптеров Kafka.
package kafka

import (
	"github.com/segmentio/kafka-go"
)

// createReader возвращает подготовленный kafka.Reader для заданных адреса, топика и группы.
func createReader(address, topic, groupID string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{address},
		Topic:   topic,
		GroupID: groupID,
	})
	return r
}
