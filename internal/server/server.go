// Package server provides a New func for spinning up a new API server.
package server

import (
	"os"
	"time"

	"github.com/bradford-hamilton/slice-it-api/internal/storage"
	tollboothChi "github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// API is a structure that orchestrates our http layer, database,
// and the communication between the them.
type API struct {
	baseURL string
	db      storage.URLRepository
	Mux     *chi.Mux
}

// New takes a storage.Repository and set's up an API server, using that store.
func New(db storage.URLRepository) *API {
	r := chi.NewRouter()
	r.Use(
		corsMiddleware().Handler,
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.StripSlashes,                        // strip slashes to no slash URL versions
		middleware.Recoverer,                           // recover from panics without crashing server
		middleware.Timeout(30*time.Second),             // start with a pretty standard timeout
		tollboothChi.LimitHandler(limiterMiddleware()), // set up a basic rate limiter by ip
	)

	baseURL := "http://localhost:4000"
	if os.Getenv("SLICE_IT_ENVIRONMENT") == "production" {
		baseURL = "http://slice-it-api-load-balancer-475088201.us-west-2.elb.amazonaws.com"
	}

	api := &API{db: db, Mux: r, baseURL: baseURL}
	api.initializeRoutes()

	return api
}

func (a *API) initializeRoutes() {
	a.Mux.Get("/ping", a.ping)
	a.Mux.Post("/new", a.createShortURL)
	a.Mux.Get("/{urlHash}", a.redirectToLongURL)
	a.Mux.Get("/stats/{urlHash}", a.getURLStats)
}
