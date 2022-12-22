package rabbit

import (
	"encoding/json"
	"go_demo/db"

	"github.com/streadway/amqp"
)

const (
	// RabbitMQURL is the URL of the RabbitMQ server
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

// Connection is a RabbitMQ connection
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConnection creates a new RabbitMQ connection
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn:    conn,
		channel: channel,
	}, nil
}

// Close closes the RabbitMQ connection
func (c *Connection) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

// Publish publishes a message to a queue
func (c *Connection) Publish(queue string, message db.Message) error {
	// Marshal the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// publish message into the queue
	return c.channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
}

// Consume reads messages from a queue into a channel
func (c *Connection) Consume(queue string) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// CreateQueue creates a new queue
func (c *Connection) CreateQueue(queue string) error {
	_, err := c.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}
