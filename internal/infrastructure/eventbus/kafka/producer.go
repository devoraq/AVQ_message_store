package kafka

import (
	"github.com/segmentio/kafka-go"
)

// Создает нового продюсера, записывает в передаваемый на входе топик.
// Запуск происходит в инициализации кафки
func createWriter(address, topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return w
}
