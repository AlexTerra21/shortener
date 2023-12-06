package config

import (
	"fmt"

	"github.com/AlexTerra21/shortener/internal/app/storage"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	Storage       storage.Storage
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: "",
		BaseURL:       "",
		Storage:       *storage.NewStorage(),
	}
}

func (c *Config) SetServerAddress(s string) {
	c.ServerAddress = s
}

func (c *Config) SetBaseURL(s string) {
	c.BaseURL = s
}

func (c *Config) Print() {
	fmt.Println(c.ServerAddress)
	fmt.Println(c.BaseURL)
}
