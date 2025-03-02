// TP2/api/handlers/handlers.go
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"github.com/your-username/TP2/models"
)

type Handler struct {
	rabbitmqConn *amqp.Connection
	rabbitmqChan *amqp.Channel
	bo1DB        *sql.DB
	bo2DB        *sql.DB
	hoDB         *sql.DB
}

func NewHandler() (*Handler, error) {
	// RabbitMQ connection
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare queues
	queues := []string{"bo1_queue", "bo2_queue", "ho_queue"}
	for _, queueName := range queues {
		_, err = ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare queue %s: %v", queueName, err)
		}
	}

	// Database connections
	bo1DSN := fmt.Sprintf("%s:%s@tcp(%s:3306)/bo1?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("BO1_DB_HOST"),
	)
	bo1DB, err := sql.Open("mysql", bo1DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to BO1 database: %v", err)
	}

	bo2DSN := fmt.Sprintf("%s:%s@tcp(%s:3306)/bo2?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("BO2_DB_HOST"),
	)
	bo2DB, err := sql.Open("mysql", bo2DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to BO2 database: %v", err)
	}

	hoDSN := fmt.Sprintf("%s:%s@tcp(%s:3306)/ho?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("HO_DB_HOST"),
	)
	hoDB, err := sql.Open("mysql", hoDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to HO database: %v", err)
	}

	return &Handler{
		rabbitmqConn: conn,
		rabbitmqChan: ch,
		bo1DB:        bo1DB,
		bo2DB:        bo2DB,
		hoDB:         hoDB,
	}, nil
}

func (h *Handler) InsertBO1Data(c *gin.Context) {
	var data models.DataItem
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// Send to RabbitMQ
	body, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to process data",
		})
		return
	}

	err = h.rabbitmqChan.Publish(
		"",          // exchange
		"bo1_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to queue data",
		})
		return
	}

	// Also send to HO queue
	data.Source = "BO1"
	bodyHO, _ := json.Marshal(data)

	err = h.rabbitmqChan.Publish(
		"",         // exchange
		"ho_queue", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyHO,
		})

	if err != nil {
		log.Printf("Warning: Failed to queue data for HO: %v", err)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Data queued for processing",
	})
}

func (h *Handler) InsertBO2Data(c *gin.Context) {
	var data models.DataItem
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// Send to RabbitMQ
	body, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to process data",
		})
		return
	}

	err = h.rabbitmqChan.Publish(
		"",          // exchange
		"bo2_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to queue data",
		})
		return
	}

	// Also send to HO queue
	data.Source = "BO2"
	bodyHO, _ := json.Marshal(data)

	err = h.rabbitmqChan.Publish(
		"",         // exchange
		"ho_queue", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyHO,
		})

	if err != nil {
		log.Printf("Warning: Failed to queue data for HO: %v", err)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Data queued for processing",
	})
}

func (h *Handler) GetBO1Data(c *gin.Context) {
	rows, err := h.bo1DB.Query("SELECT id, name, value, timestamp FROM data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch data from BO1",
		})
		return
	}
	defer rows.Close()

	var results []models.DataItem
	for rows.Next() {
		var item models.DataItem
		var timestamp string
		if err := rows.Scan(&item.ID, &item.Name, &item.Value, &timestamp); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		item.Timestamp = timestamp
		results = append(results, item)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    results,
	})
}

func (h *Handler) GetBO2Data(c *gin.Context) {
	rows, err := h.bo2DB.Query("SELECT id, name, value, timestamp FROM data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch data from BO2",
		})
		return
	}
	defer rows.Close()

	var results []models.DataItem
	for rows.Next() {
		var item models.DataItem
		var timestamp string
		if err := rows.Scan(&item.ID, &item.Name, &item.Value, &timestamp); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		item.Timestamp = timestamp
		results = append(results, item)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    results,
	})
}

func (h *Handler) GetHOData(c *gin.Context) {
	rows, err := h.hoDB.Query("SELECT id, source, name, value, timestamp FROM data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch data from HO",
		})
		return
	}
	defer rows.Close()

	var results []models.DataItem
	for rows.Next() {
		var item models.DataItem
		var timestamp string
		if err := rows.Scan(&item.ID, &item.Source, &item.Name, &item.Value, &timestamp); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		item.Timestamp = timestamp
		results = append(results, item)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    results,
	})
}
