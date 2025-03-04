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
	defer conn.Close() // Close the connection when the program ends

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close() // Close the channel when the program ends

	// Declare the same queue as the producer
	q, err := ch.QueueDeclare(
		"test_queue", // name
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// set up a go channel to receive messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-acknowledge
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	// Continuously listen for messages
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	for msg := range msgs {
		log.Printf(" [x] Received: %s", msg.Body)
	}
}
