package main

import (
	"TP2/api/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	setupRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{os.Getenv("API_IP")}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.POST("/bo1", func(c *gin.Context) { handlers.InsertData(c,os.Getenv("BO1_DB_NAME"), os.Getenv("BO1_DB_HOST"), os.Getenv("BO1_PORT"), rabbitChannel) })
	r.POST("/bo2", func(c *gin.Context) { handlers.InsertData(c,os.Getenv("BO2_DB_NAME"), os.Getenv("BO2_DB_HOST"), os.Getenv("BO2_PORT"), rabbitChannel) })
	r.GET("/bo1", func(c *gin.Context) { handlers.GetData(c,os.Getenv("BO1_DB_NAME"), os.Getenv("BO1_DB_HOST"), os.Getenv("BO1_PORT")) })
	r.GET("/bo2", func(c *gin.Context) { handlers.GetData(c,os.Getenv("BO2_DB_NAME"), os.Getenv("BO2_DB_HOST"), os.Getenv("BO2_PORT")) })
	r.GET("/ho", func(c *gin.Context) { handlers.GetData(c, os.Getenv("HO_DB_NAME"),os.Getenv("HO_DB_HOST"), os.Getenv("HO_PORT")) })

	log.Println("Starting API server on :" + os.Getenv("API_PORT"))
	if err := r.Run(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

func setupRabbitMQ() {
	var err error

	rabbitConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = rabbitChannel.QueueDeclare(
		os.Getenv("QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Println("RabbitMQ setup completed successfully")
}
