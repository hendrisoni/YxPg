package models

import "time"

// Connection represents a PostgreSQL connection configuration
type Connection struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Database  string    `json:"database"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	SSLMode   string    `json:"ssl_mode"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// ConnectionTestResult holds the result of a connection test
type ConnectionTestResult struct {
	OK      bool   `json:"ok"`
	Latency int64  `json:"latency_ms"`
	Message string `json:"message"`
}
