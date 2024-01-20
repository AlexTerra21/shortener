package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/compress"
	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func MainRouter(c *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Post("/", logger.WithLogging(compress.WithCompress(storeURL(c))))
	r.Post("/api/shorten", logger.WithLogging(compress.WithCompress(shortenURL(c))))
	r.Get("/{id}", logger.WithLogging(getURL(c)))
	r.Get("/ping", logger.WithLogging(ping(c)))
	r.MethodNotAllowed(notAllowedHandler)
	return r
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported method", http.StatusBadRequest) // В ответе код 400
}

func shortenURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log().Debug("decoding request")
		var request models.Request
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			logger.Log().Debug("cannot decode request JSON body", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id := utils.RandSeq(8)
		if err := c.Storage.Set(r.Context(), id, request.URL); err != nil {
			logger.Log().Debug("Error adding new url", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		response := models.Response{
			Result: c.GetBaseURL() + "/" + id,
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated) // устанавливаем код 201
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Log().Debug("error encoding response", zap.Error(err))
			return
		}
	}
}

func storeURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, _ := io.ReadAll(r.Body)
		id := utils.RandSeq(8)
		if err := c.Storage.Set(r.Context(), id, string(url)); err != nil {
			logger.Log().Debug("Error adding new url", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusConflict)
			return
		}
		resp := c.GetBaseURL() + "/" + id
		w.Header().Set("content-type", "application/text")
		w.WriteHeader(http.StatusCreated) // устанавливаем код 201
		_, _ = w.Write([]byte(resp))
	}
}
func getURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		url, err := c.Storage.Get(r.Context(), id)
		if err != nil {
			logger.Log().Error("URL not found", zap.Int("status", http.StatusNotFound), zap.String("id", id))
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect) // устанавливаем код 307
	}
}

func ping(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := c.Storage.S.(*storagers.DB)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Database not supported"))
			return
		}
		if err := db.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
