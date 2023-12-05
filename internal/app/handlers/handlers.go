package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"github.com/go-chi/chi"
)

func MainRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", storeURL)
	r.Get("/{id}", getURL)
	r.MethodNotAllowed(notAllowedHandler)
	return r
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unsupported method", http.StatusBadRequest) // В ответе код 400
}

func storeURL(w http.ResponseWriter, r *http.Request) {
	url, _ := io.ReadAll(r.Body)
	id := utils.RandSeq(8)
	storage.Storage[id] = string(url)
	resp := "http://localhost:8080/" + id
	w.Header().Set("content-type", "application/text")
	w.WriteHeader(http.StatusCreated) // устанавливаем код 201
	_, _ = w.Write([]byte(resp))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getURL")
	id := chi.URLParam(r, "id")
	url := storage.Storage[id]
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect) // устанавливаем код 307
	_, _ = w.Write([]byte(""))
}
