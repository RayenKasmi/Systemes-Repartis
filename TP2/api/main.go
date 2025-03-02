// TP2/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/your-username/TP2/api/handlers"
)

func main() {
	r := gin.Default()

	// Set up CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	h, err := handlers.NewHandler()
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Define API endpoints
	r.POST("/bo1", h.InsertBO1Data)
	r.POST("/bo2", h.InsertBO2Data)
	r.GET("/bo1", h.GetBO1Data)
	r.GET("/bo2", h.GetBO2Data)
	r.GET("/ho", h.GetHOData)

	// Run server
	log.Println("Starting API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
