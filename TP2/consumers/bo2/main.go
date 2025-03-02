// TP2/consumers/bo2/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"github.com/your-username/TP2/models"
)

func main() {
	// Connect to RabbitMQ
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	// Keep trying to connect to RabbitMQ
	var conn *amqp.Connection
	var err error
	for {
		conn, err = amqp.Dial(rabbitmqURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bo2_queue", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	// Keep trying to connect to DB
	var db *sql.DB
	for {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}
	defer db.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("BO2 Consumer started")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data models.DataItem
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			// Insert into database
			_, err := db.Exec("INSERT INTO data (name, value) VALUES (?, ?)",
				data.Name, data.Value)
			if err != nil {
				log.Printf("Error inserting data: %v", err)
			} else {
				log.Printf("Inserted data: %s = %s", data.Name, data.Value)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
