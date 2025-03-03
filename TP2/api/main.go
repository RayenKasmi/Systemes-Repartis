package main

import (
	"TP2/api/config"
	"TP2/api/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	setupRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.POST("/bo1", func(c *gin.Context) { handlers.InsertData(c, "bo1", "3306", rabbitChannel) })
	r.POST("/bo2", func(c *gin.Context) { handlers.InsertData(c, "bo2", "3307", rabbitChannel) })
	r.GET("/bo1", func(c *gin.Context) { handlers.GetData(c, "bo1", "3306") })
	r.GET("/bo2", func(c *gin.Context) { handlers.GetData(c, "bo2", "3307") })
	r.GET("/ho", func(c *gin.Context) { handlers.GetData(c, "ho", "3308") })

	log.Println("Starting API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

func setupRabbitMQ() {
	var err error

	rabbitConn, err = amqp.Dial(config.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = rabbitChannel.QueueDeclare(
		config.QueueName,
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
