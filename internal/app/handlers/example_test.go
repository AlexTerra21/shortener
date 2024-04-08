package handlers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Хэндлер получает оригинальный URL и возвращает сокращенный
// Отдает куку авторизации с уникальным ID пользователя
func Example_shortenURL() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/api/shorten"
	body := `{ "url": "https://practicum.yandex.ru/" }`

	resp, err := client.R().
		SetHeader("Content-Type", "application/text").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(body).
		Post(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

// Хэндлер получает оригинальный URL и возвращает сокращенный
// Устаревший метод
// Отдает куку авторизации с уникальным ID пользователя
func Example_storeURL() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/"
	body := "https://practicum.yandex.ru"

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(body).
		Post(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

// Хэндлер получает сокращенный URL и возвращает оригинальный
// Используется кука авторизации из shorten
func Example_getURL() {
	client := resty.New()

	id := "qwerty"

	endpoint := "http://localhost:8080"
	handler := fmt.Sprintf("\\%s", id)

	resp, err := client.R().
		Get(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

// Хэндлер проверяет доступность БД
func Example_ping() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/ping"

	resp, err := client.R().
		Get(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

func Example_batch() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/api/shorten/batch"
	body := `[{
				"correlation_id": "1",
				"original_url": "https://practicum.yandex.ru/"
			},
			{
				"correlation_id": "2",
				"original_url": "https://rambler.ru/"
			},
		 	{
				"correlation_id": "3",
				"original_url": "https://mail.ru/"
			}]`

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(body).
		Post(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

// Хэндлер отдает все записи по текущему пользователю
// Используется кука авторизации из shorten
// ID пользователя храниться в куке авторизации
func Example_delete() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/api/user/urls"

	resp, err := client.R().
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		Get(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}

// Хэндлер получает массив идентификаторов и ставит метки удаления
// Фильтрация по текущему ID пользователя
// Используется кука авторизации из shorten
// ID пользователя храниться в куке авторизации
func Example_urls() {
	client := resty.New()

	endpoint := "http://localhost:8080"
	handler := "/api/user/urls"
	body := `["ktjiPtll","wmWmpxpZ", "kYfcjdJE"]`

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(body).
		Delete(endpoint + handler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}
