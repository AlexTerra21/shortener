package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/storage"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func parseFlags() *config.Config {
	config := config.NewConfig()

	start := flag.String("a", ":8080", "address and port to run server")
	ret := flag.String("b", "http://localhost:8080/", "address and port to return")

	flag.Parse()
	config.SetServerStartURL(*start)
	config.SetReturnURL(*ret)
	return config
}

// ./cmd/shortener/shortener.exe --help
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8085/
// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	config := parseFlags()
	config.Print()
	utils.RandInit()
	storage.Storage = make(map[string]string)
	fmt.Println("Running server on", config.ServerStartURL)
	return http.ListenAndServe(config.ServerStartURL, handlers.MainRouter(config))
}
