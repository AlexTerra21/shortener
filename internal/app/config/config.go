package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexTerra21/shortener/internal/app/storage"
)

type Config struct {
	serverAddress   string
	baseURL         string
	logLevel        string
	fileStoragePath string
	dbConnectString string
	Storage         *storage.Storage
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) InitStorage() (err error) {
	// logger.Log().Info(c.dbConnectString)
	c.Storage, err = storage.NewStorage(c.fileStoragePath, c.dbConnectString)
	return err
}

func (c *Config) GetFileStoragePath() string {
	return c.fileStoragePath
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

func (c *Config) GetLogLevel() string {
	return c.logLevel
}

func (c *Config) Print() {
	fmt.Printf("Server address: %s\n", c.serverAddress)
	fmt.Printf("Base URL: %s\n", c.baseURL)
	fmt.Printf("Log level: %s\n", c.logLevel)
}

func (c *Config) ParseFlags() {
	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "address and port to return")
	logLevel := flag.String("l", "info", "log level")
	fileStoragePath := flag.String("f", "", "file name for url save")
	dbConnectString := flag.String("d", "", "db connection string")

	flag.Parse()
	if serverAddressEnv := os.Getenv("SERVER_ADDRESS"); serverAddressEnv != "" {
		serverAddress = &serverAddressEnv
	}
	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		baseURL = &baseURLEnv
	}
	if logLevelEnv := os.Getenv("LOG_LEVEL"); logLevelEnv != "" {
		logLevel = &logLevelEnv
	}
	if fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH"); fileStoragePathEnv != "" {
		fileStoragePath = &fileStoragePathEnv
	}
	if dbConnectStringEnv := os.Getenv("DATABASE_DSN"); dbConnectStringEnv != "" {
		fileStoragePath = &dbConnectStringEnv
	}
	c.serverAddress = *serverAddress
	c.baseURL = *baseURL
	c.logLevel = *logLevel
	c.fileStoragePath = *fileStoragePath
	c.dbConnectString = *dbConnectString
}
