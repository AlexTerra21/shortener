package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		storeUrl(w, r)
	case http.MethodGet:
		getUrl(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusBadRequest)
	}

}

func storeUrl(w http.ResponseWriter, r *http.Request) {
	url, _ := io.ReadAll(r.Body)
	id := utils.RandSeq(8)
	Storage[id] = string(url)
	resp := "http://localhost:8080/" + id
	fmt.Println(resp)
	w.Header().Set("content-type", "application/text")
	w.WriteHeader(http.StatusCreated) // устанавливаем код 201
	w.Write([]byte(resp))
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	url := Storage[id]
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect) // устанавливаем код 307
	w.Write([]byte(""))
}
