// TP2/models/models.go
package models

// DataItem represents the data structure to be stored in the database
type DataItem struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Timestamp string `json:"timestamp,omitempty"`
	Source    string `json:"source,omitempty"` // Used by HO service to track the source
}

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
