package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/logger"
)

// Сервер HTTP
type HTTPServer struct {
	http *http.Server
}

// Конструктор для создания HTTP сервера
func NewHTTPServer(config *config.Config) (*HTTPServer, error) {
	server := &HTTPServer{
		http: &http.Server{
			Addr:    config.ServerAddress,
			Handler: handlers.MainRouter(config),
		},
	}
	return server, nil
}

// Старт HTTP сервера
func (s *HTTPServer) Start() error {
	logger.Log().Info("Running HTTP server", zap.String("address", s.http.Addr))
	return s.http.ListenAndServe()
}

// Остановка HTTP сервера
func (s *HTTPServer) Stop() error {
	logger.Log().Info("Stopping HTTP server")
	return s.http.Shutdown(context.Background())
}
