package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/auth"
	"github.com/AlexTerra21/shortener/internal/app/compress"
	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/errs"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
	"github.com/AlexTerra21/shortener/internal/app/storage/storagers"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func MainRouter(c *config.Config) chi.Router {
	r := chi.NewRouter()
	r.Post("/", auth.WithAuth(logger.WithLogging(compress.WithCompress(storeURL(c)))))
	r.Post("/api/shorten", auth.WithAuth(logger.WithLogging(compress.WithCompress(shortenURL(c)))))
	r.Post("/api/shorten/batch", auth.WithAuth(logger.WithLogging(compress.WithCompress(batch(c)))))
	r.Get("/api/user/urls", logger.WithLogging(compress.WithCompress(urls(c))))
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
		userID := r.Context().Value(auth.UserIDKey).(int)
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
		var response models.Response
		w.Header().Set("content-type", "application/json")
		if err := c.Storage.Set(r.Context(), id, request.URL, userID); err != nil {
			logger.Log().Debug("Error adding new url", zap.Error(err))
			if errors.Is(err, errs.ErrConflict) {
				w.WriteHeader(http.StatusConflict)
				db, ok := c.Storage.S.(*storagers.DB)
				if ok {
					id, _ := db.GetShortURL(r.Context(), request.URL, userID)
					response.Result = c.GetBaseURL() + "/" + id
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				response.Result = err.Error()
			}
		} else {
			w.WriteHeader(http.StatusCreated) // устанавливаем код 201
			response.Result = c.GetBaseURL() + "/" + id
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Log().Debug("error encoding response", zap.Error(err))
			return
		}
	}
}

func storeURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int)
		url, _ := io.ReadAll(r.Body)
		id := utils.RandSeq(8)
		var resp string
		w.Header().Set("content-type", "application/text")
		if err := c.Storage.Set(r.Context(), id, string(url), userID); err != nil {
			logger.Log().Debug("Error adding new url", zap.Error(err))
			if errors.Is(err, errs.ErrConflict) {
				w.WriteHeader(http.StatusConflict)
				db, ok := c.Storage.S.(*storagers.DB)
				if ok {
					id, _ := db.GetShortURL(r.Context(), string(url), userID)
					resp = c.GetBaseURL() + "/" + id
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				resp = err.Error()
			}
		} else {
			w.WriteHeader(http.StatusCreated) // устанавливаем код 201
			resp = c.GetBaseURL() + "/" + id
		}
		_, _ = w.Write([]byte(resp))
	}
}
func getURL(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		url, err := c.Storage.Get(r.Context(), id)
		if err != nil {
			logger.Log().Debug("URL not found", zap.Int("status", http.StatusNotFound), zap.String("id", id))
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

func batch(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int)
		logger.Log().Debug("decoding request")
		var request []models.BatchReq
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			logger.Log().Debug("cannot decode request JSON body", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := make([]models.BatchResp, 0)
		batchStor := make([]models.BatchStore, 0)
		for _, value := range request {
			if value.OriginalURL == "" {
				continue
			}
			id := utils.RandSeq(8)
			resp := models.BatchResp{
				CorrelationID: value.CorrelationID,
				ShortURL:      c.GetBaseURL() + "/" + id,
			}
			batch := models.BatchStore{
				OriginalURL: value.OriginalURL,
				IdxShortURL: id,
			}
			batchStor = append(batchStor, batch)
			response = append(response, resp)
		}

		if err := c.Storage.BatchSet(r.Context(), &batchStor, userID); err != nil {
			logger.Log().Debug("Error adding new url", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
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

func urls(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := auth.CheckAuth(r)
		if userID < 0 {
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}
		db, ok := c.Storage.S.(*storagers.DB)
		if !ok {
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Database not supported"))
			return
		}
		response, err := db.GetAll(r.Context(), c.GetBaseURL(), userID)
		if err != nil {
			logger.Log().Debug("error get all urls from DB", zap.Error(err))
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if response == nil {
			logger.Log().Debug("Empty DB")
			w.Header().Set("content-type", "application/text")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			logger.Log().Debug("error encoding response", zap.Error(err))
			return
		}

	}
}
