package server

import (
	"github.com/AlexTerra21/shortener/internal/app/config"
)

// Интерфейс HTTP/HTTPS сервера
type Server interface {
	Start() error
	Stop() error
}

// Конструктор для создания сервера
func NewServer(config *config.Config) (Server, error) {
	if config.EnableHTTPS {
		return NewHTTPSServer(config)
	} else {
		return NewHTTPServer(config)
	}
}
