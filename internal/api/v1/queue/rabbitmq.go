package queue

import (
	"backend-go/internal/api/v1/interfaces"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQ(amqpURI, queueName string) (interfaces.MessageQueue, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (r *RabbitMQ) GetChannel() (*amqp.Channel, error) {
	return r.channel, nil
}

func (r *RabbitMQ) GetQueueName() string {
	return r.queue
}

func (r *RabbitMQ) ConsumeMessages() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		r.queue, "", false, false, false, false, nil,
	)
}

func (r *RabbitMQ) SendMessage(body []byte) error {
	return r.channel.Publish(
		"",
		r.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) AckMessage(msg amqp.Delivery) error {
	return msg.Ack(false)
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
