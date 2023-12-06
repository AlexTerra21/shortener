package handlers

import (
	"io"
	"net/http"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"github.com/go-chi/chi"
)

func MainRouter(c *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Post("/", storeURL(c))
	r.Get("/{id}", getURL(c))
	r.MethodNotAllowed(notAllowedHandler)
	return r
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported method", http.StatusBadRequest) // В ответе код 400
}

func storeURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, _ := io.ReadAll(r.Body)
		id := utils.RandSeq(8)
		c.Storage.Set(id, string(url))
		resp := c.BaseURL + "/" + id
		w.Header().Set("content-type", "application/text")
		w.WriteHeader(http.StatusCreated) // устанавливаем код 201
		_, _ = w.Write([]byte(resp))
	}
}
func getURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		url := c.Storage.Get(id)
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect) // устанавливаем код 307
		_, _ = w.Write([]byte(""))
	}
}
