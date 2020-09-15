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
	GetViewCount(urlHash string) (int, error)
}

// SliceItURL represents a URL in in our datastore.
type SliceItURL struct {
	ID        int    `json:"id,omitempty"`
	Short     string `json:"short,omitempty"`
	Long      string `json:"long,omitempty"`
	ViewCount int    `json:"view_count,omitempty"`
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

	// Thought here: would we actually want this decoupled and have the caller use it after fetching a URL?
	// We may not want every "Get" call here to do this. This is okay for now, explain thoughts to Jim.
	if err := db.incrementViewCount(urlHash); err != nil {
		return "", err
	}

	return url.Long, nil
}

// GetViewCount takes a short URL hash and finds and returns the view count stats from that URL
func (db *Db) GetViewCount(urlHash string) (int, error) {
	query := "SELECT view_count FROM urls WHERE short = $1;"

	var url SliceItURL
	row := db.QueryRow(query, urlHash)
	if err := row.Scan(&url.ViewCount); err != nil {
		return 0, err
	}

	return url.ViewCount, nil
}

func (db *Db) incrementViewCount(urlHash string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "UPDATE urls SET view_count = view_count + 1 WHERE short = $1;"
	if _, err = tx.Exec(query, urlHash); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
