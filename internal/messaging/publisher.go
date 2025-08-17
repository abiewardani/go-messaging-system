package messaging

import (
    "github.com/streadway/amqp"
    "log"
)

type Publisher struct {
    Channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
    return &Publisher{Channel: channel}
}

func (p *Publisher) Publish(queueName string, message []byte) error {
    err := p.Channel.Publish(
        "",         // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        message,
        },
    )
    if err != nil {
        log.Printf("Failed to publish message: %s", err)
        return err
    }
    log.Printf("Message published to queue: %s", queueName)
    return nil
}