package main

import (
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Read RabbitMQ URL from environment variable
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	// Connect to RabbitMQ
	conn, err := amqp091.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"test_queue", // name of the queue
		false,        // not durable
		false,        // auto-delete when unused
		false,        // not exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Publish a test message
	body := "Hello from Dockerized Producer!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
