package main

import (
	"fmt"
	"log"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// go build -o cmd/shortener/shortener.exe -ldflags "-X main.buildVersion=v1.20.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git log -1 | grep commit)'" cmd/shortener/*.go

// ./cmd/shortener/shortener.exe --help
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug -f ./tmp/short-url-db.json
// ./cmd/shortener/shortener.exe -a=:8091 -b=http://localhost:8091 -l debug -d "host=localhost user=shortner password=userpassword dbname=short_urls sslmode=disable"
// ./cmd/shortener/shortener.exe -a=:443 -s -b=http://localhost:443 -l debug -d "host=localhost user=shortner password=userpassword dbname=short_urls sslmode=disable"
//
// функция main вызывается автоматически при запуске приложения
func main() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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

	server, err := server.NewServer(config)
	if err != nil {
		return err
	}
	defer server.Stop()

	config.InitAsync()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	// go http.ListenAndServe("0.0.0.0:8080", nil)
	sig := <-signalCh
	logger.Log().Sugar().Infof("Received signal: %v\n", sig)

	return nil
}
