package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bradford-hamilton/slice-it-api/internal/storage"
	"github.com/bradford-hamilton/slice-it-api/internal/urls"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (a *API) ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("pong"))
}

type createReq struct {
	LongURL string `json:"longURL"`
}

func (a *API) createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req createReq
	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.LongURL == "" {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	url := storage.SliceItURL{
		Short: a.baseURL + "/" + urls.Shorten(req.LongURL),
		Long:  req.LongURL,
	}
	if err := a.db.Create(url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, fmt.Sprintf(`{ "shortURL": %s }`, url.Short))
}

func (a *API) redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	urlHash := chi.URLParam(r, "urlHash")
	longURL, err := a.db.Get(a.baseURL + "/" + urlHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
