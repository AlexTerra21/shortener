package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexTerra21/shortener/internal/app/storage"
)

type Config struct {
	serverAddress string
	baseURL       string
	Storage       storage.Storage
}

func NewConfig() *Config {
	return &Config{
		serverAddress: "",
		baseURL:       "",
		Storage:       *storage.NewStorage(),
	}
}

func (c *Config) SetServerAddress(s string) {
	c.serverAddress = s
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
}

func (c *Config) SetBaseURL(s string) {
	c.baseURL = s
}

func (c *Config) GetBaseURL() string {
	return c.baseURL
}

func (c *Config) Print() {
	fmt.Println(c.serverAddress)
	fmt.Println(c.baseURL)
}

func (c *Config) ParseFlags() {
	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "address and port to return")

	flag.Parse()
	if serverAddressEnv := os.Getenv("SERVER_ADDRESS"); serverAddressEnv != "" {
		serverAddress = &serverAddressEnv
	}
	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		baseURL = &baseURLEnv
	}
	c.serverAddress = *serverAddress
	c.baseURL = *baseURL
}
