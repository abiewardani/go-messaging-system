package messaging

import (
    "log"
    "github.com/streadway/amqp"
)

type Consumer struct {
    Channel *amqp.Channel
    Queue   string
}

func NewConsumer(channel *amqp.Channel, queue string) *Consumer {
    return &Consumer{
        Channel: channel,
        Queue:   queue,
    }
}

func (c *Consumer) StartConsuming() {
    msgs, err := c.Channel.Consume(
        c.Queue,
        "",    // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %s", err)
    }

    for msg := range msgs {
        log.Printf("Received a message: %s", msg.Body)
        // Process the message here
        msg.Ack(false) // Acknowledge the message
    }
}