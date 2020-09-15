package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradford-hamilton/slice-it-api/internal/storage"
)

// mockRepo implements the TodoRepository interface, and is used for testing.
type mockRepo struct{}

func (m *mockRepo) Create(url storage.SliceItURL) error { return nil }
func (m *mockRepo) Get(urlHash string) (string, error)  { return "", nil }

func TestAPI_ping(t *testing.T) {
	api := New(&mockRepo{})
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ping)
	handler.ServeHTTP(rr, req)

	// Check that the status code and content are what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := []byte("pong")
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.Bytes(), expected)
	}
}

type body struct {
	ShortURL string `json:"shortURL"`
}

func TestAPI_create(t *testing.T) {
	api := New(&mockRepo{})

	req, err := http.NewRequest("POST", "/new", bytes.NewBuffer([]byte(`{ "longURL": "https://www.reddit.com/r/golang/comments/is81vi/just_because_you_can_doesnt_mean_you_should" }`)))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.createShortURL)
	handler.ServeHTTP(rr, req)

	// Check that the status code and content are what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// `{ "shortURL\": "043fdd58" }`
	expected := []byte{34, 123, 32, 92, 34, 115, 104, 111, 114, 116, 85, 82, 76, 92, 34, 58, 32, 104, 116, 116, 112, 58, 47, 47, 108, 111, 99, 97, 108, 104, 111, 115, 116, 58, 52, 48, 48, 48, 47, 48, 52, 51, 102, 100, 100, 53, 56, 32, 125, 34, 10}
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.Bytes(), expected)
	}
}
