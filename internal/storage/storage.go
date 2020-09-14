// Package storage describes the interface for interacting with our datastore, and provides a PostgreSQL
// implementation for the slice-it-api.
package storage

import (
	"database/sql"
	"fmt"
	"os"

	// postgres driver
	_ "github.com/lib/pq"
)

// URLRepository describes the interface for interacting with our datastore. This can viewed
// like a plug in adapter, making testing and/or switching datastores much more trivial.
type URLRepository interface {
	Create(longURL string) (string, error)
}

// URL represents a URL in in our datastore.
type URL struct {
	ID        int    `json:"id,omitempty"`
	Short     string `json:"short,omitempty"`
	Long      string `json:"long,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Db provides a set of methods for interacting with our database.
type Db struct {
	*sql.DB
}

// NewDB creates a connection with our postgres database and returns it, otherwise an error.
func NewDB() (*Db, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("SLICE_IT_API_DB_HOST"),
		os.Getenv("SLICE_IT_API_DB_PORT"),
		os.Getenv("SLICE_IT_API_DB_USER"),
		os.Getenv("SLICE_IT_API_DB_PASSWORD"),
		os.Getenv("SLICE_IT_API_DB_NAME"),
		os.Getenv("SLICE_IT_API_SSL_MODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

// Create ...
func (db *Db) Create(longURL string) (string, error) { return "", nil }
