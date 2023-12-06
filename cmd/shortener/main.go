package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func parseFlags() *config.Config {
	config := config.NewConfig()

	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "address and port to return")

	flag.Parse()
	if serverAddressEnv := os.Getenv("SERVER_ADDRESS"); serverAddressEnv != "" {
		serverAddress = &serverAddressEnv
	}
	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		baseURL = &baseURLEnv
	}
	config.SetServerAddress(*serverAddress)
	config.SetBaseURL(*baseURL)
	return config
}

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
	config := parseFlags()
	config.Print()
	utils.RandInit()
	fmt.Println("Running server on", config.ServerAddress)
	return http.ListenAndServe(config.ServerAddress, handlers.MainRouter(config))
}
