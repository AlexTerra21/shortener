package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/server"
	pb "github.com/AlexTerra21/shortener/proto"
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
// ./cmd/shortener/shortener.exe -c ./config/config.json
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
	config, err := config.NewConfig()
	if err != nil {
		// logger.Log().Error("Read config error: %v", err)
		return err
	}
	config.Print()
	if err = logger.Initialize(config.LogLevel); err != nil {
		return err
	}
	if err = config.InitStorage(); err != nil {
		return err
	}
	defer config.Storage.S.Close()

	restServer, err := server.NewServer(config)
	if err != nil {
		return err
	}
	defer restServer.Stop()

	grpcServer, err := pb.NewGRPCServer(config)
	if err != nil {
		return err
	}
	defer grpcServer.Stop()

	config.InitAsync()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := restServer.Start(); err != http.ErrServerClosed && err != nil {
			logger.Log().Sugar().Errorf("Server error: %v", err)
			signalCh <- syscall.SIGTERM
		}
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Log().Sugar().Errorf("Server error: %v", err)
			signalCh <- syscall.SIGTERM
		}
	}()
	// go http.ListenAndServe("0.0.0.0:8080", nil)
	sig := <-signalCh
	logger.Log().Sugar().Infof("Received signal: %v", sig)

	return nil
}
