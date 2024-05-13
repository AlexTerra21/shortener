package server

import (
	"context"
	"crypto/tls"
	"net/http"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/AlexTerra21/shortener/internal/app/handlers"
	"github.com/AlexTerra21/shortener/internal/app/logger"
)

// Сервер HTTPS
type HTTPSServer struct {
	https *http.Server
}

// Конструктор для создания HTTPS сервера
func NewHTTPSServer(config *config.Config) (*HTTPSServer, error) {

	certificate, err := GenerateCertificate()
	if err != nil {
		return nil, err
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
	}

	server := &HTTPSServer{
		https: &http.Server{
			Addr:      config.ServerAddress,
			Handler:   handlers.MainRouter(config),
			TLSConfig: &tlsConfig,
		},
	}
	return server, nil
}

// Старт HTTPS сервера
func (s *HTTPSServer) Start() error {
	logger.Log().Info("Running HTTPS server", zap.String("address", s.https.Addr))
	return s.https.ListenAndServeTLS("", "")
}

// Остановка HTTPS сервера
func (s *HTTPSServer) Stop() error {
	logger.Log().Info("Stopping HTTPS server")
	return s.https.Shutdown(context.Background())
}
