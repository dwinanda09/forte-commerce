package domain

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueuePublisher interface {
	Publish(ctx context.Context, body []byte) error
}

type QueueConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}
