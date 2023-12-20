package main

import (
	"fmt"
	"net/http"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

// ./cmd/shortener/shortener.exe --help
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8085
// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	config := config.NewConfig()
	config.ParseFlags()
	// config.Print()
	utils.RandInit()
	fmt.Println("Running server on", config.GetServerAddress())
	return http.ListenAndServe(config.GetServerAddress(), handlers.MainRouter(config))
}
