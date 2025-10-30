package mongodb

import "errors"

var (
	// ErrConnect описывает ошибку установления соединения с MongoDB.
	ErrConnect = errors.New("mongodb: connect failed")
	// ErrPing сигнализирует о сбое при проверке доступности MongoDB.
	ErrPing = errors.New("mongodb: ping failed")
	// ErrDisconnect сообщает о неудачном завершении соединения с MongoDB.
	ErrDisconnect = errors.New("mongodb: disconnect failed")
)
