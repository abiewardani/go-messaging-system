package messaging

import (
    "github.com/streadway/amqp"
    "log"
)

type RabbitMQ struct {
    Connection *amqp.Connection
    Channel    *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &RabbitMQ{
        Connection: conn,
        Channel:    ch,
    }, nil
}

func (r *RabbitMQ) Close() {
    if err := r.Channel.Close(); err != nil {
        log.Fatalf("Failed to close channel: %s", err)
    }
    if err := r.Connection.Close(); err != nil {
        log.Fatalf("Failed to close connection: %s", err)
    }
}