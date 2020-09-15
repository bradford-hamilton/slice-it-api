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
	Create(url SliceItURL) error
	Get(urlHash string) (string, error)
}

// SliceItURL represents a URL in in our datastore.
type SliceItURL struct {
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

// Create handles inserting a SliceItURL into the database. With no other requirements, we
// don't need to return anything but an error if it happens.
func (db *Db) Create(url SliceItURL) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Rollback is safe to call even if the tx is already closed,
	// so if the tx commits successfully, this is a no-op
	defer tx.Rollback()

	query := "INSERT INTO urls (short, long) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT unique_url_constraint DO NOTHING;"
	if _, err = tx.Exec(query, url.Short, url.Long); err != nil {
		fmt.Println(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Get takes a URL hash and finds and returns the full length original link
func (db *Db) Get(urlHash string) (string, error) {
	query := "SELECT long FROM urls WHERE short = $1;"

	var url SliceItURL
	row := db.QueryRow(query, urlHash)
	if err := row.Scan(&url.Long); err != nil {
		return "", err
	}

	return url.Long, nil
}
