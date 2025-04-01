package queue

import (
	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQService(amqpURI, queueName string) (*RabbitMQService, error) {
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

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (r *RabbitMQService) GetChannel() (*amqp.Channel, error) {
	return r.channel, nil
}

func (r *RabbitMQService) GetQueueName() string {
	return r.queue
}

func (r *RabbitMQService) ConsumeMessages() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		r.queue, "", false, false, false, false, nil,
	)
}

func (r *RabbitMQService) SendMessage(body []byte) error {
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

func (r *RabbitMQService) AckMessage(msg amqp.Delivery) error {
	return msg.Ack(false)
}

func (r *RabbitMQService) Close() {
	r.channel.Close()
	r.conn.Close()
}
