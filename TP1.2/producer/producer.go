package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	// Declare exchange
	err = ch.ExchangeDeclare(
		os.Getenv("EXCHANGE_NAME"), // name
		"direct",                   // type of exchange
		true,                       // durable
		false,                      // auto-delete
		false,                      // internal
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Routing keys
	routingKeys := strings.Split(os.Getenv("ROUTING_KEYS"), ",")

	// Send messages for each routing key via the exchange
	for _, key := range routingKeys {
		body := fmt.Sprintf("Message with routing key: %s", key)
		err = ch.Publish(
			os.Getenv("EXCHANGE_NAME"), // exchange
			key,                       // routing key
			false,                     // mandatory
			false,                     // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		failOnError(err, fmt.Sprintf("Failed to publish message with routing key %s", key))

		log.Printf(" [x] Sent %s", body)
	}
}
