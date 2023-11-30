// пакеты исполняемых приложений должны называться main
package main

import (
	"net/http"

	. "github.com/AlexTerra21/shortener/internal/app/handlers"
	. "github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	utils.RandInit()
	Storage = make(map[string]string)
	return http.ListenAndServe(`:8080`, http.HandlerFunc(MainHandler))
}
