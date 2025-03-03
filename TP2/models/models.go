package models

import "time"

type ProductSale struct {
	ID       int64     `json:"id,omitempty"`
	SaleDate time.Time `json:"sale_date" binding:"required"`
	Region   string    `json:"region" binding:"required"`
	Product  string    `json:"product" binding:"required"`
	Qty      int       `json:"qty" binding:"required"`
	Cost     float64   `json:"cost" binding:"required"`
	Amt      float64   `json:"amt" binding:"required"`
	Tax      float64   `json:"tax" binding:"required"`
	Total    float64   `json:"total" binding:"required"`
	Source   string    `json:"source,omitempty"` // Used by HO service to track the source
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
