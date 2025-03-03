package handlers

import (
	"TP2/api/config"
	"TP2/api/database"
	"TP2/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type ProductSale models.ProductSale

func InsertData(c *gin.Context, dbName, port string, rabbitChannel *amqp.Channel) {
	db, err := database.ConnectDB(dbName, port)
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
	if err = publishMessage(sale, dbName, rabbitChannel); err != nil {
		log.Printf("Warning: Failed to publish to RabbitMQ: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}

func GetData(c *gin.Context, dbName, port string) {
	db, err := database.ConnectDB(dbName, port)
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

func publishMessage(sale ProductSale, source string, rabbitChannel *amqp.Channel) error {
	sale.Source = source

	body, err := json.Marshal(sale)
	if err != nil {
		return fmt.Errorf("error marshalling sale: %v", err)
	}

	err = rabbitChannel.Publish(
		"",
		config.QueueName,
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

	log.Printf("Message published to queue: %s", "product_sales_task_queue")
	return nil
}
