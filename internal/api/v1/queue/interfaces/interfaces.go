package interfaces

import "github.com/streadway/amqp"

type MessageQueue interface {
	GetChannel() (*amqp.Channel, error)
	GetQueueName() string
	ConsumeMessages() (<-chan amqp.Delivery, error)
	SendMessage(body []byte) error
	AckMessage(msg amqp.Delivery) error
	Close()
}
