package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexTerra21/shortener/internal/app/async"
	"github.com/AlexTerra21/shortener/internal/app/storage"
)

// Структура для хранения конфигурации приложения
type Config struct {
	serverAddress   string
	baseURL         string
	logLevel        string
	fileStoragePath string
	dbConnectString string
	enableHTTPS     bool
	Storage         *storage.Storage
	DelQueue        *async.Async
}

// Инициализация конфигурации
func NewConfig() *Config {
	return &Config{}
}

// Инициализация хранилища
func (c *Config) InitStorage() (err error) {
	// logger.Log().Info(c.dbConnectString)
	c.Storage, err = storage.NewStorage(c.fileStoragePath, c.dbConnectString)
	return
}

// Инициализация процесса для удаления записей из БД
func (c *Config) InitAsync() {
	c.DelQueue = async.NewAsync(c.Storage)
}

// Получение параметра конфигурации fileStoragePath
func (c *Config) GetFileStoragePath() string {
	return c.fileStoragePath
}

// Присваивание параметра конфигурации serverAddress
func (c *Config) SetServerAddress(s string) {
	c.serverAddress = s
}

// Получение параметра конфигурации serverAddress
func (c *Config) GetServerAddress() string {
	return c.serverAddress
}

// Присваивание параметра конфигурации baseURL
func (c *Config) SetBaseURL(s string) {
	c.baseURL = s
}

// Получение параметра конфигурации baseURL
func (c *Config) GetBaseURL() string {
	return c.baseURL
}

// Получение параметра конфигурации logLevel
func (c *Config) GetLogLevel() string {
	return c.logLevel
}

// Получение признака разрешения HTTPS
func (c *Config) GetEnableHTTPS() bool {
	return c.enableHTTPS
}

// Печать основных параметров конфигурации
func (c *Config) Print() {
	fmt.Printf("Server address: %s\n", c.serverAddress)
	fmt.Printf("Base URL: %s\n", c.baseURL)
	fmt.Printf("Log level: %s\n", c.logLevel)
}

// Обработка флагов командной строки и занесение в параметры конфигурации
func (c *Config) ParseFlags() {
	serverAddress := flag.String("a", ":8080", "address and port to run server")
	baseURL := flag.String("b", "http://localhost:8080", "address and port to return")
	logLevel := flag.String("l", "info", "log level")
	fileStoragePath := flag.String("f", "", "file name for url save")
	dbConnectString := flag.String("d", "", "db connection string")
	enableHTTPS := flag.Bool("s", false, "enable HTTPS")

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
	if enableHTTPSEnv := os.Getenv("ENABLE_HTTPS"); enableHTTPSEnv != "" {
		fileStoragePath = &enableHTTPSEnv
	}
	c.serverAddress = *serverAddress
	c.baseURL = *baseURL
	c.logLevel = *logLevel
	c.fileStoragePath = *fileStoragePath
	c.dbConnectString = *dbConnectString
	c.enableHTTPS = *enableHTTPS
}
