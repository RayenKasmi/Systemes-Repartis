// TP2/api/main.go
package main

import (
	"TP2/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type ProductSale models.ProductSale

func main() {
	setupRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	r := gin.Default()

	r.POST("/bo1", func(c *gin.Context) { InsertData(c, "bo1", "3306") })
	r.POST("/bo2", func(c *gin.Context) { InsertData(c, "bo2", "3307") })
	r.GET("/bo1", func(c *gin.Context) { GetData(c, "bo1", "3306") })
	r.GET("/bo2", func(c *gin.Context) { GetData(c, "bo2", "3307") })
	r.GET("/ho", func(c *gin.Context) { GetData(c, "ho", "3308") })

	log.Println("Starting API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

const (
	queueName   = "product_sales_task_queue"
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

func setupRabbitMQ() {
	var err error

	rabbitConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = rabbitChannel.QueueDeclare(
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

	log.Println("RabbitMQ setup completed successfully")
}

func publishMessage(sale ProductSale, source string) error {
	sale.Source = source

	body, err := json.Marshal(sale)
	if err != nil {
		return fmt.Errorf("error marshalling sale: %v", err)
	}

	err = rabbitChannel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("Message published to queue: %s", queueName)
	return nil
}

func connectDB(dbName, port string) (*sql.DB, error) {
	dsn := fmt.Sprintf("user:password@tcp(localhost:%s)/%s?parseTime=true", port, dbName)
	fmt.Println("Connecting to:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func InsertData(c *gin.Context, dbName, port string) {
	db, err := connectDB(dbName, port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	var sale ProductSale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO ProductSales (sale_date, region, product, qty, cost, amt, tax, total) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, sale.SaleDate, sale.Region, sale.Product, sale.Qty, sale.Cost, sale.Amt, sale.Tax, sale.Total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	sale.Source = dbName
	if err = publishMessage(sale, dbName); err != nil {
		log.Printf("Warning: Failed to publish to RabbitMQ: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}

func GetData(c *gin.Context, dbName, port string) {
	db, err := connectDB(dbName, port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, sale_date, region, product, qty, cost, amt, tax, total FROM ProductSales")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer rows.Close()

	var sales []ProductSale
	for rows.Next() {
		var sale ProductSale
		if err := rows.Scan(&sale.ID, &sale.SaleDate, &sale.Region, &sale.Product, &sale.Qty, &sale.Cost, &sale.Amt, &sale.Tax, &sale.Total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning data"})
			return
		}
		sales = append(sales, sale)
	}

	c.JSON(http.StatusOK, sales)
}
