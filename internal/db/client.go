package db

import (
	"fmt"

	ladybug "github.com/gimigkk/Alfred-Memory/internal/ladybug"
)

type Client struct {
	DB *ladybug.Database
}

func NewClient(dbPath string) (*Client, error) {
	db, err := ladybug.NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ladybug database: %w", err)
	}
	return &Client{DB: db}, nil
}

func (c *Client) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}

// GetConnection returns a new Connection that MUST be explicitly closed by the caller.
// defer conn.Close() MUST be called immediately after calling this function.
func (c *Client) GetConnection() (*ladybug.Connection, error) {
	conn, err := ladybug.NewConnection(c.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return conn, nil
}
