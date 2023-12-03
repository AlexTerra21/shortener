package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		storeURL(w, r)
	case http.MethodGet:
		getURL(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusBadRequest)
	}

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
	id := strings.TrimPrefix(r.URL.Path, "/")
	url := storage.Storage[id]
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect) // устанавливаем код 307
	_, _ = w.Write([]byte(""))
}
