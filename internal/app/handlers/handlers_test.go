package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_storeURL_getURL(t *testing.T) {
	// Инициализация сервисов
	utils.RandInit()
	config := config.NewConfig()
	config.SetServerAddress(":8080")
	config.SetBaseURL("http://localhost:8080")
	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(MainRouter(config))
	// останавливаем сервер после завершения теста
	defer srv.Close()

	// Данные для теста
	requestedURL := "https://practicum.yandex.ru/"
	postCode := http.StatusCreated
	postContentType := "application/text"
	getCode := http.StatusTemporaryRedirect
	testName := "complex test store and get url #1"
	serverURL := srv.URL
	t.Run(testName, func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().
			SetBody(requestedURL).
			Post(serverURL)
		assert.NoError(t, err)
		assert.Equal(t, postCode, resp.StatusCode()) // 201
		assert.Equal(t, postContentType, resp.Header().Get("Content-Type"))
		// Получить ID из ответа
		parseID := strings.Split(string(resp.Body()), "/")
		id := parseID[len(parseID)-1]
		// Отключаем авто редирект, что бы поймать ответ метода getURL
		client.SetRedirectPolicy(resty.NoRedirectPolicy())
		resp, err = client.R().
			SetPathParams(map[string]string{
				"ID": id,
			}).
			Get(serverURL + "/{ID}")
		assert.Error(t, err)                        // Ошибка auto redirect is disabled
		assert.Equal(t, getCode, resp.StatusCode()) // 307
		assert.Equal(t, requestedURL, resp.Header().Get("Location"))
	})
}

func TestHandlers_MainHandler(t *testing.T) {
	config := config.NewConfig()
	config.SetServerAddress(":8080")
	config.SetBaseURL("http://localhost:8080")
	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(MainRouter(config))
	// останавливаем сервер после завершения теста
	defer srv.Close()
	tests := []struct {
		name   string
		method string
		code   int
		body   string
	}{
		{
			name:   "negative PUT request test #1",
			method: http.MethodPut,
			code:   http.StatusBadRequest,
			body:   "Unsupported method\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = srv.URL
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.code, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.body, string(resp.Body()))
		})
	}
}
