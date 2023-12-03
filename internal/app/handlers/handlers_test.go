package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers_storeURL_getURL(t *testing.T) {
	// Инициализация сервисов
	utils.RandInit()
	storage.Storage = make(map[string]string)
	// Данные для теста
	requestedURL := "https://practicum.yandex.ru/"
	postCode := http.StatusCreated
	postResponse := "http://localhost:8080/"
	postContentType := "application/text"
	getCode := http.StatusTemporaryRedirect
	testName := "complex test store and get url #1"
	t.Run(testName, func(t *testing.T) {
		bodyReader := strings.NewReader(requestedURL)
		requestPost := httptest.NewRequest(http.MethodPost, "/", bodyReader)
		wPost := httptest.NewRecorder()

		storeURL(wPost, requestPost)

		resPost := wPost.Result()
		assert.Equal(t, postCode, resPost.StatusCode)
		defer resPost.Body.Close()
		resBody, err := io.ReadAll(resPost.Body)
		require.NoError(t, err)
		assert.Contains(t, string(resBody), postResponse)
		assert.Equal(t, postContentType, resPost.Header.Get("Content-Type"))

		id := strings.TrimPrefix(string(resBody), postResponse)

		t.Logf("ID = %s", id)

		requestGet := httptest.NewRequest(http.MethodGet, "/"+id, nil)
		wGet := httptest.NewRecorder()

		getURL(wGet, requestGet)

		resGet := wGet.Result()
		defer resGet.Body.Close()
		assert.Equal(t, getCode, resGet.StatusCode)
		assert.Equal(t, requestedURL, resGet.Header.Get("Location"))

	})
}

func TestHandlers_MainHandler(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "negative PUT request test #1",
			code: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPut, "/", nil)
			w := httptest.NewRecorder()
			MainHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, tt.code, res.StatusCode)

		})
	}
}
