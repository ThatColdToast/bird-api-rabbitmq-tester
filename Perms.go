package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PermissionsManager struct {
	uri      string
	port     uint
	username string
	password string

	connection *amqp.Connection
	channel    *amqp.Channel
}

func makePermissionsManager(
	uri string,
	port uint,
	username string,
	password string,
) PermissionsManager {
	log.Printf("Connecting\n")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, uri, port))
	failOnError(err, "Failed to connect to RabbitMQ")

	log.Printf("Creating Channel\n")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return PermissionsManager{
		uri:      uri,
		port:     port,
		username: username,
		password: password,

		connection: conn,
		channel:    ch,
	}
}

func (x PermissionsManager) Close() {
	x.connection.Close()
	x.channel.Close()
}

func (x PermissionsManager) check(body string) {
	// log.Printf("Creating Callback Queue\n")
	queue, err := x.channel.QueueDeclare(
		"",    // anonymous name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// log.Printf("Consuming on %s\n", queue.Name)
	msgs, err := x.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		true,       // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = x.channel.PublishWithContext(ctx,
		"",              // exchange
		"perms_request", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId.String(),
			ReplyTo:       queue.Name,
			Body:          []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	for msg := range msgs {
		// log.Printf("Received: '%s' with '%s' - '%s'", msg.Body, msg.CorrelationId, corrId)
		if corrId.String() == msg.CorrelationId {
			log.Printf("Correlated: '%s'", msg.Body)

			// msg.Ack(true)
			break
		} else {
			// msg.Ack(false)
			// break
		}
	}
}
