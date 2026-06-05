package queue

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "forte.checkout"
	QueueName    = "forte.checkout.queue"
	RoutingKey   = "checkout"
)

type RabbitMQ struct {
	url     string
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	r := &RabbitMQ{url: url}
	if err := r.reconnect(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RabbitMQ) reconnect() error {
	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}

	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	r.conn = conn
	r.channel = channel
	return r.setup()
}

func (r *RabbitMQ) setup() error {
	if err := r.channel.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if _, err := r.channel.QueueDeclare(
		QueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := r.channel.QueueBind(
		QueueName,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, body []byte) error {
	return r.channel.PublishWithContext(
		ctx,
		ExchangeName,
		RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	if err := r.channel.Qos(1, 0, false); err != nil {
		// channel is dead — reconnect and retry once
		if reconnErr := r.reconnect(); reconnErr != nil {
			return nil, fmt.Errorf("failed to set QoS: %w", err)
		}
		if err := r.channel.Qos(1, 0, false); err != nil {
			return nil, fmt.Errorf("failed to set QoS after reconnect: %w", err)
		}
	}

	return r.channel.Consume(
		QueueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
