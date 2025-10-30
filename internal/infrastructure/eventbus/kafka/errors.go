package kafka

import "errors"

var (
	// ErrEnsureConnection описывает ошибку установления соединения с брокером.
	ErrEnsureConnection = errors.New("kafka: ensure connection failed")
	// ErrWriteMessage сигнализирует о сбое при записи сообщения.
	ErrWriteMessage = errors.New("kafka: write message failed")
	// ErrFetchMessage сообщает о неудачном чтении сообщения.
	ErrFetchMessage = errors.New("kafka: fetch message failed")
	// ErrCommitMessage означает ошибку подтверждения оффсета.
	ErrCommitMessage = errors.New("kafka: commit message failed")
)
