// TP2/consumer/ho/main.go
package main

import (
	"TP2/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	queueName     string
	rabbitMQURL   string
	hoDatabaseDSN string
)

type ProductSale models.ProductSale

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting HO database consumer...")

	queueName = os.Getenv("QUEUE_NAME")
	rabbitMQURL = os.Getenv("RABBITMQ_URL")
	hoDatabaseDSN = os.Getenv("HO_DB_DSN")

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	/*
		If there are multiple workers, to ensure fair dispatch
		we need to set the prefetch count to 1
	*/
	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Consumer started. Waiting for messages on queue: %s", queueName)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			var sale ProductSale
			if err := json.Unmarshal(msg.Body, &sale); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				msg.Nack(false, false) //rejects message with no requeue
				continue
			}

			err = insertIntoHO(sale)
			if err != nil {
				log.Printf("Error inserting into HO: %v. Requeuing message.", err)
				time.Sleep(2 * time.Second)
				msg.Nack(false, true) //rejects message with requeue
				continue
			}

			msg.Ack(false)
			log.Printf("Successfully processed message from %s and inserted into HO", sale.Source)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-stopChan
	log.Println("Shutting down consumer...")
}

func connectHODB() (*sql.DB, error) {
	db, err := sql.Open("mysql", hoDatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func insertIntoHO(sale ProductSale) error {
	db, err := connectHODB()
	if err != nil {
		log.Fatalf("Failed to connect to HO database: %v", err)
		return fmt.Errorf("failed to connect to HO database: %v", err)
	}
	defer db.Close()

	query := "INSERT INTO ProductSales (sale_date, region, product, qty, cost, amt, tax, total) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, sale.SaleDate, sale.Region, sale.Product, sale.Qty, sale.Cost, sale.Amt, sale.Tax, sale.Total)
	if err != nil {
		log.Printf("Failed to insert data into HO: %v", err)
		return fmt.Errorf("failed to insert data into HO: %v", err)
	}

	log.Println("Record from %s inserted into HO database successfully", sale.Source)
	return nil
}
