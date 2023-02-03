package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendMessages(count int) {
	if count <= 0 {
		log.Panicf("sendMessages must have a count of more than 0")
	}

	log.Printf("Connecting\n") // Connect to RabbitMQ with plain auth
	conn, err := amqp.Dial("amqp://user:pass@172.17.0.2:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	log.Printf("Creating Channel\n") // Create Channel (lightweight TCP connection) to RabbitMQ
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Printf("Creating Queue\n") // Create Queue (stack to have messages pushed and popped)
	queue, err := ch.QueueDeclare(
		"perms", // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Create Timeout Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Sample Messages (who|what) or (who|what|where)
	body := [10]string{
		"c68b1eb3-75e5-446e-b1fe-7f26fc5588e4|public.post.create",
		"9da1c4e6-b674-4536-9b04-7a7e781c26cf|public.post.create",
		"c248cb3f-0d7b-487c-8ed9-f90bdbe6d19a|public.post.create",
		"f9b91a8e-2280-4d1e-9987-7a1918a3eba7|public.post.create",
		"821e60f0-aeda-4533-bfd1-e853deb8edc3|public.post.create",
		"ef70a148-d046-4a88-815a-d8142f15839b|public.post.create",
		"c248cb3f-0d7b-487c-8ed9-f90bdbe6d19a|public.post.delete",
		"f9b91a8e-2280-4d1e-9987-7a1918a3eba7|public.post.delete",
		"821e60f0-aeda-4533-bfd1-e853deb8edc3|public.post.delete",
		"ef70a148-d046-4a88-815a-d8142f15839b|public.post.delete",
	}

	for i := 0; i < count; i++ {
		err = ch.PublishWithContext(ctx,
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(body[i*69%10]), // pick randomish message
				ReplyTo:       queue.Name,
				CorrelationId: "rand",
			})
		failOnError(err, "Failed to publish a message")
		// log.Printf("[x] Sent %s\n", body)
		// time.Sleep(time.Duration(i) * time.Millisecond)
		// time.Sleep(1 * time.Second)
	}
}
