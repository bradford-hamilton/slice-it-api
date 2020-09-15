package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bradford-hamilton/slice-it-api/internal/server"
	"github.com/bradford-hamilton/slice-it-api/internal/storage"
)

func main() {
	db, err := storage.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := server.New(db)
	port := os.Getenv("SLICE_IT_API_SERVER_PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Printf("Serving application on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, s.Mux))
}
