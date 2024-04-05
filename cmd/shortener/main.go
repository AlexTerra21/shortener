package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/logger"
)

// ./cmd/shortener/shortener.exe --help
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug -f ./tmp/short-url-db.json
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug -d "host=localhost user=shortner password=userpassword dbname=short_urls sslmode=disable"
//
// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() (err error) {
	config := config.NewConfig()
	config.ParseFlags()
	config.Print()
	if err = logger.Initialize(config.GetLogLevel()); err != nil {
		return err
	}
	if err = config.InitStorage(); err != nil {
		return err
	}
	defer config.Storage.S.Close()

	config.InitAsync()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Log().Info("Running server", zap.String("address", config.GetServerAddress()))
		err := http.ListenAndServe(config.GetServerAddress(), handlers.MainRouter(config))
		if err != nil {
			log.Fatal(err)
		}
	}()
	go http.ListenAndServe("0.0.0.0:8080", nil)
	sig := <-signalCh
	logger.Log().Sugar().Infof("Received signal: %v\n", sig)

	return nil
}
