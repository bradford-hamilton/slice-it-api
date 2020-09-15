// Package server provides a New func for spinning up a new API server.
package server

import (
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
	db  storage.URLRepository
	Mux *chi.Mux
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
	api := &API{db: db, Mux: r}
	api.initializeRoutes()
	return api
}

func (a *API) initializeRoutes() {
	a.Mux.Get("/ping", a.ping)
	a.Mux.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/urls", func(r chi.Router) {
				r.Post("/", a.createShortURL)
			})
		})
	})
}
